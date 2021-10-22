package handler

import (
	"context"
	"fmt"

	"github.com/ONSdigital/dp-cantabular-metadata-exporter/config"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/event"
)

// CantabularMetadataExport is the event handler for the CantabularMetadataExport event
type CantabularMetadataExport struct {
	cfg         config.Config
	dataset     DatasetAPIClient
	file        FileManager
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
	if err := h.exportCSVW(e); err != nil{
		return fmt.Errorf("failed to export csvw: %w", err)
	}

	if err := h.exportTXTFile(e); err != nil{
		return fmt.Errorf("failed to export metadata text file: %w", err)
	}

	return nil
}

func (h *CantabularMetadataExport) exportCSVW(e *event.CantabularMetadataExport) error {
	return nil
}

func (h *CantabularMetadataExport) exportTXTFile(e *event.CantabularMetadataExport) error {
	return nil
}
