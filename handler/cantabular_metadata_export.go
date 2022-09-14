package handler

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/config"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/csvw"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/event"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/schema"
	"github.com/pkg/errors"

	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/log.go/v2/log"
)

const (
	batchSize  = 1000
	maxWorkers = 10
)

type downloadInfo struct {
	size        int
	publicURL   string
	privateURL  string
	downloadURL string
}

// CantabularMetadataExport is the event handler for the CantabularMetadataExport event
type CantabularMetadataExport struct {
	cfg               config.Config
	dataset           DatasetAPIClient
	filter            FilterAPIClient
	file              FileManager
	producer          kafka.IProducer
	generate          Generator
	csvwPrefix        string
	metadataExtension string
	apiDomainURL      string
}

// NewCantabularMetadataExport creates a new CantabularMetadataExportHandler
func NewCantabularMetadataExport(cfg config.Config, d DatasetAPIClient, f FilterAPIClient, fm FileManager, p kafka.IProducer, g Generator) *CantabularMetadataExport {
	return &CantabularMetadataExport{
		cfg:      cfg,
		dataset:  d,
		filter:   f,
		file:     fm,
		producer: p,
		generate: g,
	}
}

// Handle takes a single event.
func (h *CantabularMetadataExport) Handle(ctx context.Context, workerID int, msg kafka.Message) error {
	e := &event.CSVCreated{}
	s := schema.CSVCreated

	if err := s.Unmarshal(msg.GetData(), e); err != nil {
		return &Error{
			err: fmt.Errorf("failed to unmarshal event: %w", err),
			logData: map[string]interface{}{
				"msg_data": msg.GetData(),
			},
		}
	}

	logData := log.Data{
		"event": e,
	}

	isFilterJob := len(e.FilterOutputID) != 0

	req := dataset.GetVersionMetadataSelectionInput{
		UserAuthToken:    "",
		ServiceAuthToken: h.cfg.ServiceAuthToken,
		CollectionID:     "",
		DatasetID:        e.DatasetID,
		Edition:          e.Edition,
		Version:          e.Version,
		Dimensions:       e.Dimensions,
	}

	m, err := h.dataset.GetVersionMetadataSelection(ctx, req)
	if err != nil {
		return Error{
			err:     fmt.Errorf("failed to get version metadata: %w", err),
			logData: logData,
		}
	}

	isPublished, err := h.isVersionPublished(ctx, e)
	if err != nil {
		return Error{
			err:     fmt.Errorf("failed to determine published state: %w", err),
			logData: logData,
		}
	}

	logData["isPublished"] = isPublished

	csvwDownload, err := h.exportCSVW(ctx, e, *m, isPublished)
	if err != nil {
		return Error{
			err:     fmt.Errorf("failed to export csvw: %w", err),
			logData: logData,
		}
	}

	txtDownload, err := h.exportTXTFile(ctx, e, *m, isPublished)
	if err != nil {
		return Error{
			err:     fmt.Errorf("failed to export metadata text file: %w", err),
			logData: logData,
		}
	}

	if isFilterJob {
		if err := h.UpdateFilterOutput(ctx, e, csvwDownload, txtDownload); err != nil {
			return Error{
				err:     errors.Wrap(err, "failed to update filter output"),
				logData: logData,
			}
		}
	} else {
		if err := h.UpdateInstance(ctx, e, csvwDownload, txtDownload, m.Version.CollectionID); err != nil {
			return Error{
				err:     errors.Wrap(err, "failed to update instance"),
				logData: logData,
			}
		}
	}

	if err := h.produceOutputMessage(e); err != nil {
		return Error{
			err:     fmt.Errorf("failed to producer output kafka message: %w", err),
			logData: logData,
		}
	}

	return nil
}

func (h *CantabularMetadataExport) exportTXTFile(ctx context.Context, e *event.CSVCreated, m dataset.Metadata, isPublished bool) (*downloadInfo, error) {
	dimensions, err := h.dataset.GetVersionDimensions(ctx, "", h.cfg.ServiceAuthToken, "", e.DatasetID, e.Edition, e.Version)
	if err != nil {
		return nil, fmt.Errorf("failed to get version dimensions: %w", err)
	}

	b, err := h.getText(ctx, m, dimensions, e)
	if err != nil {
		return nil, fmt.Errorf("failed to get text bytes: %w", err)
	}

	var url string

	if isPublished {
		url, err = h.file.Upload(bytes.NewReader(b), h.generateTextFilename(e))
	} else {
		url, err = h.file.UploadPrivate(bytes.NewReader(b), h.generateTextFilename(e), h.generateVaultPath(e, "txt"))
	}
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	d := downloadInfo{
		size:        len(b),
		downloadURL: h.generateDownloadURL(e, "txt"),
	}

	if isPublished {
		d.publicURL = url
	} else {
		d.privateURL = url
	}

	return &d, nil
}

func (h *CantabularMetadataExport) exportCSVW(ctx context.Context, e *event.CSVCreated, m dataset.Metadata, isPublished bool) (*downloadInfo, error) {
	filename := h.generateCSVWFilename(e)
	downloadURL := h.generateDownloadURL(e, "csv-metadata.json")
	aboutURL := h.dataset.GetMetadataURL(e.DatasetID, e.Edition, e.Version)

	f, err := csvw.Generate(ctx, &m, downloadURL, aboutURL, h.apiDomainURL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate csvw: %w", err)
	}

	log.Info(ctx, "uploading csvw file to s3", log.Data{
		"filename":       filename,
		"private_bucket": !isPublished,
	})

	var url string
	if isPublished {
		url, err = h.file.Upload(bytes.NewReader(f), filename)
	} else {
		url, err = h.file.UploadPrivate(bytes.NewReader(f), filename, h.generateVaultPath(e, "csvw"))
	}
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	d := downloadInfo{
		size:        len(f),
		downloadURL: downloadURL,
	}

	if isPublished {
		d.publicURL = url
	} else {
		d.privateURL = url
	}

	return &d, nil
}

func (h *CantabularMetadataExport) generateTextFilename(e *event.CSVCreated) string {
	var prefix, suffix string

	if len(e.FilterOutputID) == 0 {
		prefix, suffix = "datasets/", ".txt"
	} else {
		prefix = fmt.Sprintf("datasets/%s/", e.FilterOutputID)
		suffix = fmt.Sprintf("-%s.txt", h.generate.Timestamp().Format(time.RFC3339))
	}

	fn := fmt.Sprintf("%s-%s-%s", e.DatasetID, e.Edition, e.Version)

	return prefix + fn + suffix
}

func (h *CantabularMetadataExport) generateCSVWFilename(e *event.CSVCreated) string {
	var prefix, suffix string

	if len(e.FilterOutputID) == 0 {
		prefix, suffix = "datasets/", ".csvw"
	} else {
		prefix = fmt.Sprintf("datasets/%s/", e.FilterOutputID)
		suffix = fmt.Sprintf("-%s.csvw", h.generate.Timestamp().Format(time.RFC3339))
	}

	fn := fmt.Sprintf(
		"%s%s-%s-%s",
		h.csvwPrefix,
		e.DatasetID,
		e.Edition,
		e.Version,
	)

	return prefix + fn + suffix
}

func (h *CantabularMetadataExport) generateDownloadURL(e *event.CSVCreated, extension string) string {
	return fmt.Sprintf(
		"%s/downloads/datasets/%s/editions/%s/versions/%s.%s",
		h.cfg.DownloadServiceURL,
		e.DatasetID,
		e.Edition,
		e.Version,
		extension,
	)
}

// generateVaultPathForFile generates the vault path for the provided root and filename
func (h *CantabularMetadataExport) generateVaultPath(e *event.CSVCreated, filetype string) string {
	return fmt.Sprintf("%s/%s-%s-%s.%s", h.cfg.VaultPath, e.DatasetID, e.Edition, e.Version, filetype)
}

// getText gets a byte array containing the metadata content, based on options returned by dataset API.
// If a dimension has more than maxMetadataOptions, an error will be returned
func (h *CantabularMetadataExport) getText(ctx context.Context, metadata dataset.Metadata, dimensions dataset.VersionDimensions, e *event.CSVCreated) ([]byte, error) {
	var b bytes.Buffer

	b.WriteString(metadata.ToString())
	b.WriteString("Dimensions:\n")

	for _, dimension := range dimensions.Items {
		options, err := h.dataset.GetOptionsInBatches(
			ctx,
			"",
			h.cfg.ServiceAuthToken,
			"",
			e.DatasetID,
			e.Edition,
			e.Version,
			dimension.Name,
			batchSize,
			maxWorkers,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to get dimension options: %w", err)
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
		InstanceID: e.InstanceID,
		DatasetID:  e.DatasetID,
		Edition:    e.Edition,
		Version:    e.Version,
		RowCount:   e.RowCount,
	})
	if err != nil {
		return fmt.Errorf("error marshalling csvw created event: %w", err)
	}

	// Send bytes to kafka producer output channel
	h.producer.Channels().Output <- b

	return nil
}

func (h *CantabularMetadataExport) UpdateInstance(ctx context.Context, e *event.CSVCreated, csvwInfo, txtInfo *downloadInfo, collectionID string) error {
	log.Info(ctx, "updating instance with download link")

	downloads := map[string]dataset.Download{
		"CSVW": dataset.Download{
			Size:    fmt.Sprintf("%d", csvwInfo.size),
			Public:  csvwInfo.publicURL,
			Private: csvwInfo.privateURL,
			URL:     csvwInfo.downloadURL,
		},
		"TXT": dataset.Download{
			Size:    fmt.Sprintf("%d", txtInfo.size),
			Public:  txtInfo.publicURL,
			Private: txtInfo.privateURL,
			URL:     txtInfo.downloadURL,
		},
	}

	v := dataset.Version{
		Downloads: downloads,
	}

	if err := h.dataset.PutVersion(
		ctx,
		"",
		h.cfg.ServiceAuthToken,
		collectionID,
		e.DatasetID,
		e.Edition,
		e.Version,
		v,
	); err != nil {
		return Error{
			err: errors.Wrap(err, "failed to put version"),
			logData: log.Data{
				"downloads": downloads,
			},
		}
	}

	return nil
}

func (h *CantabularMetadataExport) UpdateFilterOutput(ctx context.Context, e *event.CSVCreated, csvwInfo, txtInfo *downloadInfo) error {
	log.Info(ctx, "Updating filter output with download link")

	txtDownload := filter.Download{
		URL:     fmt.Sprintf("%s/downloads/filter-outputs/%s.txt", h.cfg.DownloadServiceURL, e.FilterOutputID),
		Size:    fmt.Sprintf("%d", txtInfo.size),
		Skipped: false,
	}

	txtDownload.Public = txtInfo.publicURL
	txtDownload.Private = txtInfo.privateURL

	csvwDownload := filter.Download{
		URL:     fmt.Sprintf("%s/downloads/filter-outputs/%s.csv-metadata.json", h.cfg.DownloadServiceURL, e.FilterOutputID),
		Size:    fmt.Sprintf("%d", csvwInfo.size),
		Skipped: false,
	}

	csvwDownload.Public = csvwInfo.publicURL
	csvwDownload.Private = csvwInfo.privateURL

	m := filter.Model{
		Downloads: map[string]filter.Download{
			"CSVW": csvwDownload,
			"TXT":  txtDownload,
		},
	}

	if err := h.filter.UpdateFilterOutput(ctx, "", h.cfg.ServiceAuthToken, "", e.FilterOutputID, &m); err != nil {
		return errors.Wrap(err, "failed to update filter output")
	}

	return nil
}
