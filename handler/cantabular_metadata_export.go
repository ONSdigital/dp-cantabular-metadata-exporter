package handler

import (
	"context"
	"fmt"
	"bytes"

	"github.com/ONSdigital/dp-cantabular-metadata-exporter/config"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/csvw"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/event"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/log.go/v2/log"
)

// CantabularMetadataExport is the event handler for the CantabularMetadataExport event
type CantabularMetadataExport struct {
	cfg               config.Config
	dataset           DatasetAPIClient
	file              FileManager
	csvwPrefix        string
	metadataExtension string
	apiDomainURL      string
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
	if err := h.exportCSVW(ctx, e); err != nil{
		return fmt.Errorf("failed to export csvw: %w", err)
	}

	if err := h.exportTXTFile(e); err != nil{
		return fmt.Errorf("failed to export metadata text file: %w", err)
	}

	return nil
}

func (h *CantabularMetadataExport) exportTXTFile(e *event.CantabularMetadataExport) error {
	return nil
}

func (h *CantabularMetadataExport) exportCSVW(ctx context.Context, e *event.CantabularMetadataExport) error{
	ver := fmt.Sprintf("%d", e.Version)
	filename := fmt.Sprintf(
		"%s%s-%s-v%d.csvw",
		h.csvwPrefix,
		e.DatasetID,
		e.Edition,
		e.Version,
	)

	downloadURL := fmt.Sprintf(
		"%s/downloads/datasets/%s/editions/%s/versions/%d.csvw",
		h.cfg.DownloadServiceURL,
		e.DatasetID,
		e.Edition,
		e.Version,
	) // Get downloadURL from somewhere else?

	m, err := h.dataset.GetVersionMetadata(ctx, "", h.cfg.ServiceAuthToken, "", e.DatasetID, e.Edition, ver)
	if err != nil {
		return fmt.Errorf("failed to get version metadata: %w", err)
	}

	aboutURL := h.dataset.GetMetadataURL(e.DatasetID, e.Edition, ver)

	f, err := csvw.Generate(ctx, &m, downloadURL, aboutURL, h.apiDomainURL)
	if err != nil {
		return fmt.Errorf("failed to generate csvw: %w", err)
	}

	isPublished, err := h.isVersionPublished(ctx, e)
	if err != nil{
		return fmt.Errorf("failed to determine published state: %w", err)
	}

	var url string
	if isPublished{
		url, err = h.file.Upload(bytes.NewReader(f), filename)
	} else {
		url, err = h.file.UploadPrivate(bytes.NewReader(f), filename, h.generateVaultPath(e.DatasetID))
	}
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
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

	log.Info(ctx, "updating dataset api with download link", log.Data{
		"isPublished":      true, //isPublished,
		"metadataDownload": download,
	})

	v := dataset.Version{
		Downloads: map[string]dataset.Download{
			"CSVW": *download,
		},
	}

	if err := h.dataset.PutVersion(
		ctx,
		"", 
		h.cfg.ServiceAuthToken,
		"", 
		e.DatasetID,
		e.Edition,
		ver,
		v,
	); err != nil {
		return fmt.Errorf("failed to update version: %w", err)
	}

	return nil
}

func generateCSVWFilename(e *event.CantabularMetadataExport) string {
	return fmt.Sprintf("%s-%s-%d.csvw", e.DatasetID, e.Edition, e.Version)
}

// generateVaultPathForFile generates the vault path for the provided root and filename
func (h *CantabularMetadataExport) generateVaultPath(instanceID string) string {
	return fmt.Sprintf("%s/%s", h.cfg.VaultPath, instanceID)
}

func (h *CantabularMetadataExport) isVersionPublished(ctx context.Context, e *event.CantabularMetadataExport) (bool, error){
	version, err := h.dataset.GetVersion(
		ctx, 
		"",
		h.cfg.ServiceAuthToken,
		"",
		"",
		e.DatasetID,
		e.Edition,
		fmt.Sprintf("%d", e.Version),
	)
	if err != nil {
		return false, fmt.Errorf("failed to get version: %w", err)
	}

	return version.State == dataset.StatePublished.String(), nil
}
