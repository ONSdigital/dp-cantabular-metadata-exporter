package handler

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/config"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/csvw"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/event"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/schema"

	"github.com/ONSdigital/log.go/v2/log"
	kafka "github.com/ONSdigital/dp-kafka/v2"
)

const maxMetadataOptions = 1000

// CantabularMetadataExport is the event handler for the CantabularMetadataExport event
type CantabularMetadataExport struct {
	cfg               config.Config
	dataset           DatasetAPIClient
	file              FileManager
	producer          kafka.IProducer
	csvwPrefix        string
	metadataExtension string
	apiDomainURL      string
}

// NewCantabularMetadataExport creates a new CantabularMetadataExportHandler
func NewCantabularMetadataExport(cfg config.Config, d DatasetAPIClient, fm FileManager, p kafka.IProducer) *CantabularMetadataExport {
	return &CantabularMetadataExport{
		cfg:      cfg,
		dataset:  d,
		file:     fm,
		producer: p,
	}
}

// Handle takes a single event.
func (h *CantabularMetadataExport) Handle(ctx context.Context, e *event.CSVCreated) error {
	logData := log.Data{
		"event": e,
	}

	m, err := h.dataset.GetVersionMetadata(ctx, "", h.cfg.ServiceAuthToken, "", e.DatasetID, e.Edition, e.Version)
	if err != nil {
		return Error{
			err:     fmt.Errorf("failed to get version metadata: %w", err),
			logData: logData,
		}
	}

	isPublished, err := h.isVersionPublished(ctx, e)
	if err != nil{
		return Error{
			err:     fmt.Errorf("failed to determin published state: %w", err),
			logData: logData,
		}
	}

	logData["isPublished"] = isPublished

	csvwDownload, err := h.exportCSVW(ctx, e, m, isPublished)
	if err != nil {
		return Error{
			err:     fmt.Errorf("failed to export csvw: %w", err),
			logData: logData,
		}
	}

	txtDownload, err := h.exportTXTFile(ctx, e, m, isPublished)
	if err != nil {
		return Error{
			err:     fmt.Errorf("failed to export metadata text file: %w", err),
			logData: logData,
		}
	}

	downloads := map[string]dataset.Download{
		"CSVW": *csvwDownload,
		"TXT":  *txtDownload,
	}

	logData["downloads"] = downloads

	log.Info(ctx, "updating dataset api with download link", logData)

	v := dataset.Version{
		Downloads: downloads,
	}

	if err := h.dataset.PutVersion(
		ctx,
		"",
		h.cfg.ServiceAuthToken,
		m.Version.CollectionID, 
		e.DatasetID,
		e.Edition,
		e.Version,
		v,
	); err != nil {
		return Error{
			err:     fmt.Errorf("failed to update version: %w", err),
			logData: logData,
		}
	}

	if err := h.produceOutputMessage(e); err != nil{
		return Error{
			err:     fmt.Errorf("failed to producer output kafka message: %w", err),
			logData: logData,
		}
	}

	return nil
}

func (h *CantabularMetadataExport) exportTXTFile(ctx context.Context, e *event.CSVCreated, m dataset.Metadata, isPublished bool) (*dataset.Download, error) {
	dimensions, err := h.dataset.GetVersionDimensions(ctx, "", h.cfg.ServiceAuthToken, "", e.DatasetID, e.Edition, e.Version)
	if err != nil {
		return nil, fmt.Errorf("failed to get version dimensions: %w", err)
	}

	b, err := h.getText(ctx, m, dimensions, e)
	if err != nil {
		return nil, fmt.Errorf("failed to get text bytes: %w", err)
	}

	var url string

	if isPublished{
		url, err = h.file.Upload(bytes.NewReader(b), h.generateTextFilename(e))
	} else {
		url, err = h.file.UploadPrivate(bytes.NewReader(b), h.generateTextFilename(e), h.generateVaultPath(e.DatasetID))
	}
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	download := &dataset.Download{
		Size: fmt.Sprintf("%d", len(b)),
	}

	if isPublished {
		download.Public = url
	} else {
		download.Private = url
	}

	download.URL = url

	return download, nil
}

func (h *CantabularMetadataExport) exportCSVW(ctx context.Context, e *event.CSVCreated, m dataset.Metadata, isPublished bool) (*dataset.Download, error) {
	filename := h.generateCSVWFilename(e)
	downloadURL := h.generateDownloadURL(e)
	aboutURL := h.dataset.GetMetadataURL(e.DatasetID, e.Edition, e.Version)

	f, err := csvw.Generate(ctx, &m, downloadURL, aboutURL, h.apiDomainURL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate csvw: %w", err)
	}

	var url string
	if isPublished {
		url, err = h.file.Upload(bytes.NewReader(f), filename)
	} else {
		url, err = h.file.UploadPrivate(bytes.NewReader(f), filename, h.generateVaultPath(e.DatasetID))
	}
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	download := &dataset.Download{
		Size: fmt.Sprintf("%d", len(f)),
	}

	if isPublished {
		download.Public = url
	} else {
		download.Private = url
	}

	download.URL = downloadURL + h.metadataExtension

	return download, nil
}

func (h *CantabularMetadataExport) generateTextFilename(e *event.CSVCreated) string {
	return fmt.Sprintf("datasets/%s-%s-%s.txt", e.DatasetID, e.Edition, e.Version)
}

func (h *CantabularMetadataExport) generateCSVWFilename(e *event.CSVCreated) string {
	return fmt.Sprintf(
		"datasets/%s%s-%s-%s.csvw",
		h.csvwPrefix,
		e.DatasetID,
		e.Edition,
		e.Version,
	)
}

func (h *CantabularMetadataExport) generateDownloadURL(e *event.CSVCreated) string {
	return fmt.Sprintf(
		"%s/downloads/datasets/%s/editions/%s/versions/%s.csvw",
		h.cfg.DownloadServiceURL,
		e.DatasetID,
		e.Edition,
		e.Version,
	)
}

// generateVaultPathForFile generates the vault path for the provided root and filename
func (h *CantabularMetadataExport) generateVaultPath(datasetID string) string {
	return fmt.Sprintf("%s/%s.txt", h.cfg.VaultPath, datasetID)
}

// getText gets a byte array containing the metadata content, based on options returned by dataset API.
// If a dimension has more than maxMetadataOptions, an error will be returned
func (h *CantabularMetadataExport) getText(ctx context.Context, metadata dataset.Metadata, dimensions dataset.VersionDimensions, e *event.CSVCreated) ([]byte, error) {
	var b bytes.Buffer

	b.WriteString(metadata.ToString())
	b.WriteString("Dimensions:\n")

	for _, dimension := range dimensions.Items {
		q := dataset.QueryParams{
			Offset: 0,
			Limit: maxMetadataOptions,
		}

		options, err := h.dataset.GetOptions(ctx, "", h.cfg.ServiceAuthToken, "", e.DatasetID, e.Edition, e.Version, dimension.Name, &q)
		if err != nil {
			return nil, fmt.Errorf("failed to get dimension options: %w", err)
		}

		if options.TotalCount > maxMetadataOptions {
			return nil, errors.New("too many options in dimension")
		}

		b.WriteString(options.String())
	}

	return b.Bytes(), nil
}

func (h *CantabularMetadataExport) isVersionPublished(ctx context.Context, e *event.CSVCreated) (bool, error) {
	version, err := h.dataset.GetVersion(
		ctx,
		"",
		h.cfg.ServiceAuthToken,
		"",
		"",
		e.DatasetID,
		e.Edition,
		e.Version,
	)
	if err != nil {
		return false, fmt.Errorf("failed to get version: %w", err)
	}

	return version.State == dataset.StatePublished.String(), nil
}

func (h *CantabularMetadataExport) produceOutputMessage(e *event.CSVCreated) error {
	s := schema.CSVWCreated

	b, err := s.Marshal(&event.CSVWCreated{
		InstanceID:   e.InstanceID,
		DatasetID:    e.DatasetID,
		Edition:      e.Edition,
		Version:      e.Version,
		RowCount:     e.RowCount,
	})
	if err != nil {
		return fmt.Errorf("error marshalling csvw created event: %w", err)
	}

	// Send bytes to kafka producer output channel
	h.producer.Channels().Output <- b

	return nil
}
