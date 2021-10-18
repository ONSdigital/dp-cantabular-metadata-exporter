package handler

import (
	"context"

	"github.com/ONSdigital/dp-cantabular-metadata-exporter/config"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/event"
)

// CantabularMetadataExport is the event handler for the CantabularMetadataExport event
type CantabularMetadataExport struct {
	cfg         config.Config
}

// NewCantabularMetadataExport creates a new CantabularMetadataExportHandler
func NewCantabularMetadataExport(cfg config.Config) *CantabularMetadataExport {
	return &CantabularMetadataExport{
		cfg:         cfg,
	}
}

// Handle takes a single event.
func (h *CantabularMetadataExport) Handle(ctx context.Context, e *event.CantabularMetadataExport) error {
	return nil
}