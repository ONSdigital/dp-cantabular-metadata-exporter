package handler

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/config"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/event"
	"github.com/ONSdigital/log.go/v2/log"
)

const maxMetadataOptions = 1000

// CantabularMetadataExport is the event handler for the CantabularMetadataExport event
type CantabularMetadataExport struct {
	cfg     config.Config
	dataset DatasetAPIClient
	file    FileManager
}

// NewCantabularMetadataExport creates a new CantabularMetadataExportHandler
func NewCantabularMetadataExport(cfg config.Config, d DatasetAPIClient, fm FileManager) *CantabularMetadataExport {
	return &CantabularMetadataExport{
		cfg:     cfg,
		dataset: d,
		file:    fm,
	}
}

// Handle takes a single event.
func (h *CantabularMetadataExport) Handle(ctx context.Context, e *event.CantabularMetadataExport) error {
	if err := h.exportCSVW(e); err != nil {
		return fmt.Errorf("failed to export csvw: %w", err)
	}

	if err := h.exportTXTFile(ctx, e); err != nil {
		return fmt.Errorf("failed to export metadata text file: %w", err)
	}

	return nil
}

func (h *CantabularMetadataExport) exportCSVW(e *event.CantabularMetadataExport) error {
	return nil
}

func (h *CantabularMetadataExport) exportTXTFile(ctx context.Context, e *event.CantabularMetadataExport) error {
	ver := fmt.Sprintf("%d", e.Version)
	metadata, err := h.dataset.GetVersionMetadata(ctx, "", h.cfg.ServiceAuthToken, e.CollectionID, e.DatasetID, e.Edition, ver)
	if err != nil {
		return fmt.Errorf("failed to get version metadata: %w", err)
	}

	dimensions, err := h.dataset.GetVersionDimensions(ctx, "", h.cfg.ServiceAuthToken, e.CollectionID, e.DatasetID, e.Edition, ver)
	if err != nil {
		return fmt.Errorf("failed to get version dimensions: %w", err)
	}

	b, err := h.getText(ctx, metadata, dimensions, e)
	if err != nil {
		return fmt.Errorf("failed to get text bytes: %w", err)
	}

	var url string
	// if metadata.Version.State == dataset.StatePublished.String() {
	if true {
		url, err = h.file.Upload(bytes.NewReader(b), h.generateTextFilename(e))
	} else {
		url, err = h.file.UploadPrivate(bytes.NewReader(b), h.generateTextFilename(e), h.generateVaultPath(e.DatasetID))
	}
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	log.Info(ctx, "got text", log.Data{"text": url})

	download := &dataset.Download{
		Size: fmt.Sprintf("%d", len(b)),
	}

	/* What's happening with public/private?
	if isPublished {
		download.Public = url
	} else {
		download.Private = url
	}*/
	download.Public = url
	download.Private = url
	download.URL = url

	log.Info(ctx, "updating dataset api with download link", log.Data{
		"isPublished":      true, //isPublished,
		"metadataDownload": download,
	})

	v := dataset.Version{
		Downloads: map[string]dataset.Download{
			"TXT": *download,
		},
	}

	if err := h.dataset.PutVersion(
		ctx,
		"",
		h.cfg.ServiceAuthToken,
		"", e.DatasetID,
		e.Edition,
		ver,
		v,
	); err != nil {
		return fmt.Errorf("failed to update version: %w", err)
	}

	return nil
}

func (h *CantabularMetadataExport) generateTextFilename(e *event.CantabularMetadataExport) string {
	return fmt.Sprintf("%s-%s-%d.txt", e.DatasetID, e.Edition, e.Version)
}

// generateVaultPathForFile generates the vault path for the provided root and filename
func (h *CantabularMetadataExport) generateVaultPath(instanceID string) string {
	return fmt.Sprintf("%s/%s.txt", h.cfg.VaultPath, instanceID)
}

// getText gets a byte array containing the metadata content, based on options returned by dataset API.
// If a dimension has more than maxMetadataOptions, an error will be returned
func (h *CantabularMetadataExport) getText(ctx context.Context, metadata dataset.Metadata, dimensions dataset.VersionDimensions, e *event.CantabularMetadataExport) ([]byte, error) {
	var b bytes.Buffer
	ver := fmt.Sprintf("%d", e.Version)

	b.WriteString(metadata.ToString())
	b.WriteString("Dimensions:\n")

	for _, dimension := range dimensions.Items {
		q := dataset.QueryParams{Offset: 0, Limit: maxMetadataOptions}
		options, err := h.dataset.GetOptions(ctx, "", h.cfg.ServiceAuthToken, e.CollectionID, e.DatasetID, e.Edition, ver, dimension.Name, &q)
		if err != nil {
			return nil, fmt.Errorf("failed to get dimension options: %w", err)
		}
		if options.TotalCount > maxMetadataOptions {
			return nil, errors.New("too many options in dimension")
		}

		b.WriteString(options.String())
	}

	log.Info(context.Background(), string(b.Bytes()), log.Data{})

	return b.Bytes(), nil
}
