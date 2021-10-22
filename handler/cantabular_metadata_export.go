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
	var userAccessToken string

	if err := h.exportCSVW(e); err != nil {
		return fmt.Errorf("failed to export csvw: %w", err)
	}

	if err := h.exportTXTFile(ctx, userAccessToken, e); err != nil {
		return fmt.Errorf("failed to export metadata text file: %w", err)
	}

	return nil
}

func (h *CantabularMetadataExport) exportCSVW(e *event.CantabularMetadataExport) error {
	return nil
}

func (h *CantabularMetadataExport) exportTXTFile(ctx context.Context, userAccessToken string, e *event.CantabularMetadataExport) error {
	metadata, err := h.dataset.GetVersionMetadata(ctx, userAccessToken, "", e.CollectionID, e.DatasetID, e.Edition, e.Version)
	if err != nil {
		return fmt.Errorf("failed to get version metadata: %w", err)
	}

	dimensions, err := h.dataset.GetVersionDimensions(ctx, userAccessToken, "", e.CollectionID, e.DatasetID, e.Edition, e.Version)
	if err != nil {
		return fmt.Errorf("failed to get version dimensions: %w", err)
	}

	b, err := h.getText(ctx, userAccessToken, metadata, dimensions, e)
	if err != nil {
		return fmt.Errorf("failed to get text bytes: %w", err)
	}

	log.Info(ctx, "got text", log.Data{"text": string(b)})

	return nil
}

// getText gets a byte array containing the metadata content, based on options returned by dataset API.
// If a dimension has more than maxMetadataOptions, an error will be returned
func (h *CantabularMetadataExport) getText(ctx context.Context, userAccessToken string, metadata dataset.Metadata, dimensions dataset.VersionDimensions, e *event.CantabularMetadataExport) ([]byte, error) {
	var b bytes.Buffer

	b.WriteString(metadata.ToString())
	b.WriteString("Dimensions:\n")

	for _, dimension := range dimensions.Items {
		q := dataset.QueryParams{Offset: 0, Limit: maxMetadataOptions}
		options, err := h.dataset.GetOptions(ctx, userAccessToken, "", e.CollectionID, e.DatasetID, e.Edition, e.Version, dimension.Name, &q)
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
