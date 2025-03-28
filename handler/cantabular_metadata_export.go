package handler

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/config"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/csvw"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/event"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/schema"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/text"

	"github.com/pkg/errors"

	kafka "github.com/ONSdigital/dp-kafka/v4"
	"github.com/ONSdigital/log.go/v2/log"
)

const (
	batchSize      = 1000
	custom         = "custom"
	datasetsPrefix = "datasets/"
	maxWorkers     = 10
	flexible       = "flexible"
	multivariate   = "multivariate"
)

type downloadInfo struct {
	size        int
	publicURL   string
	privateURL  string
	downloadURL string
}

// CantabularMetadataExport is the event handler for the CantabularMetadataExport event
type CantabularMetadataExport struct {
	cfg             config.Config
	dataset         DatasetAPIClient
	filter          FilterAPIClient
	populationTypes PopulationTypesAPIClient
	ctblr           CantabularClient
	file            FileManager
	producer        kafka.IProducer
	generate        Generator
	csvwPrefix      string
	apiDomainURL    string
}

// NewCantabularMetadataExport creates a new CantabularMetadataExportHandler
func NewCantabularMetadataExport(cfg config.Config, d DatasetAPIClient, f FilterAPIClient, t PopulationTypesAPIClient, c CantabularClient, fm FileManager, p kafka.IProducer, g Generator) *CantabularMetadataExport {
	return &CantabularMetadataExport{
		cfg:             cfg,
		dataset:         d,
		filter:          f,
		populationTypes: t,
		ctblr:           c,
		file:            fm,
		producer:        p,
		generate:        g,
	}
}

// Handle takes a single event
//
//nolint:gocyclo //cyclomatic complexity 21
func (h *CantabularMetadataExport) Handle(ctx context.Context, _ int, msg kafka.Message) error {
	var err error
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

	isFilterJob := e.FilterOutputID != ""

	req := dataset.GetVersionMetadataSelectionInput{
		UserAuthToken:    "",
		ServiceAuthToken: h.cfg.ServiceAuthToken,
		CollectionID:     "",
		DatasetID:        e.DatasetID,
		Edition:          e.Edition,
		Version:          e.Version,
		Dimensions:       e.Dimensions,
	}

	isFilter := e.FilterOutputID != ""
	var m *dataset.Metadata
	var filterOutput filter.Model

	if isFilter {
		filterOutput, err = h.filter.GetOutput(ctx, "", h.cfg.ServiceAuthToken, "", "", e.FilterOutputID)
		if err != nil {
			return errors.Wrap(err, "failed to get filter output")
		}
	}

	m, err = h.dataset.GetVersionMetadataSelection(ctx, req)
	if err != nil {
		return &Error{
			err:     errors.Wrap(err, "failed to get version metadata"),
			logData: logData,
		}
	}

	if filterOutput.Type == multivariate {
		m.Title += " - customised"
	}
	if filterOutput.Custom == nil {
		falseFlag := false
		filterOutput.Custom = &falseFlag
	}

	var dims []dataset.VersionDimension

	if isFilter {
		dims, err = h.GetFilterDimensions(ctx, filterOutput)
		if err != nil {
			return &Error{
				err:     errors.Wrap(err, "failed to get filter dimensions"),
				logData: logData,
			}
		}
		for i := range dims {
			if dims[i].IsAreaType != nil {
				if filterOutput.Type == multivariate && !*dims[i].IsAreaType {
					nonAreaTypeDimension := dataset.VersionDimension{
						Name:                 dims[i].Name,
						URL:                  fmt.Sprintf("%s/code-lists/%s", h.cfg.ExternalPrefixURL, strings.ToLower(dims[i].Name)),
						Label:                dims[i].Label,
						Description:          dims[i].Description,
						ID:                   dims[i].ID,
						NumberOfOptions:      dims[i].NumberOfOptions,
						QualityStatementText: dims[i].QualityStatementText,
					}
					m.Version.Dimensions = append(m.Version.Dimensions, nonAreaTypeDimension)
					m.CSVHeader = append(m.CSVHeader, dims[i].Name)
				}
				if *dims[i].IsAreaType {
					areaTypeDimension := dataset.VersionDimension{
						Name:                 dims[i].Name,
						URL:                  fmt.Sprintf("%s/code-lists/%s", h.cfg.ExternalPrefixURL, strings.ToLower(dims[i].Name)),
						Label:                dims[i].Label,
						Description:          dims[i].Description,
						ID:                   dims[i].ID,
						NumberOfOptions:      dims[i].NumberOfOptions,
						QualityStatementText: dims[i].QualityStatementText,
					}
					m.Version.Dimensions = append(m.Version.Dimensions, areaTypeDimension)
					m.CSVHeader = append(m.CSVHeader, dims[i].Name)
				}
			}
		}
		m.Version.Dimensions = removeDuplicateDimensions(m.Version.Dimensions)
	}

	isPublished, err := h.isVersionPublished(ctx, e)
	if err != nil {
		return Error{
			err:     fmt.Errorf("failed to determine published state: %w", err),
			logData: logData,
		}
	}

	logData["isPublished"] = isPublished

	csvwDownload, err := h.exportCSVW(ctx, e, *m, isPublished, *filterOutput.Custom)
	if err != nil {
		return Error{
			err:     fmt.Errorf("failed to export csvw: %w", err),
			logData: logData,
		}
	}

	txtDownload, err := h.exportTXTFile(ctx, e, *m, isPublished, *filterOutput.Custom)
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

	if err := h.produceOutputMessage(ctx, e); err != nil {
		return Error{
			err:     fmt.Errorf("failed to producer output kafka message: %w", err),
			logData: logData,
		}
	}

	return nil
}

func (h *CantabularMetadataExport) exportTXTFile(ctx context.Context, e *event.CSVCreated, m dataset.Metadata, isPublished, isCustom bool) (*downloadInfo, error) {
	var b []byte
	if isCustom {
		b = text.NewMetadataCustom(&m, e.FilterOutputID, h.cfg.DownloadServiceURL)
	} else {
		b = text.NewMetadata(&m)
	}
	filename := h.generateTextFilename(e, isCustom)

	var url string
	var err error
	if isPublished {
		url, err = h.file.Upload(ctx, bytes.NewReader(b), filename)
	} else {
		url, err = h.file.UploadPrivate(ctx, bytes.NewReader(b), filename, h.generateVaultPath(e, "txt"))
	}
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	d := downloadInfo{
		size:        len(b),
		downloadURL: h.generateDownloadURL(e, "txt"),
	}

	if isPublished {
		d.publicURL = fmt.Sprintf("%s/%s",
			h.cfg.S3PublicURL,
			filename,
		)
	} else {
		d.privateURL = url
	}

	return &d, nil
}

func (h *CantabularMetadataExport) exportCSVW(ctx context.Context, e *event.CSVCreated, m dataset.Metadata, isPublished, isCustom bool) (*downloadInfo, error) {
	filename := h.generateCSVWFilename(e, isCustom)
	downloadURL := h.generateDownloadURL(e, "csv-metadata.json")
	aboutURL := h.dataset.GetMetadataURL(e.DatasetID, e.Edition, e.Version)

	f, err := csvw.Generate(ctx, &m, downloadURL, aboutURL, h.apiDomainURL, h.cfg.ExternalPrefixURL, e.FilterOutputID, h.cfg.DownloadServiceURL, isCustom)
	if err != nil {
		return nil, fmt.Errorf("failed to generate csvw: %w", err)
	}

	log.Info(ctx, "uploading csvw file to s3", log.Data{
		"filename":       filename,
		"private_bucket": !isPublished,
	})

	var url string
	if isPublished {
		url, err = h.file.Upload(ctx, bytes.NewReader(f), filename)
	} else {
		url, err = h.file.UploadPrivate(ctx, bytes.NewReader(f), filename, h.generateVaultPath(e, "csvw"))
	}
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	d := downloadInfo{
		size:        len(f),
		downloadURL: downloadURL,
	}

	if isPublished {
		d.publicURL = fmt.Sprintf("%s/%s",
			h.cfg.S3PublicURL,
			filename,
		)
	} else {
		d.privateURL = url
	}

	return &d, nil
}

func (h *CantabularMetadataExport) generateTextFilename(e *event.CSVCreated, isCustom bool) string {
	var prefix, suffix, fn string

	if e.FilterOutputID == "" {
		prefix, suffix = datasetsPrefix, ".txt"
	} else {
		prefix = fmt.Sprintf("datasets/%s/", e.FilterOutputID)
		suffix = fmt.Sprintf("-%s.txt", h.generate.Timestamp().Format(time.RFC3339))
	}

	fn = fmt.Sprintf(
		"%s-%s-%s",
		e.DatasetID,
		e.Edition,
		e.Version,
	)

	if isCustom {
		fn = custom
	}

	return prefix + fn + suffix
}

func (h *CantabularMetadataExport) generateCSVWFilename(e *event.CSVCreated, isCustom bool) string {
	var prefix, suffix, fn string

	if e.FilterOutputID == "" {
		prefix, suffix = datasetsPrefix, ".csvw"
	} else {
		prefix = fmt.Sprintf("datasets/%s/", e.FilterOutputID)
		suffix = fmt.Sprintf("-%s.csvw", h.generate.Timestamp().Format(time.RFC3339))
	}

	fn = fmt.Sprintf(
		"%s%s-%s-%s",
		h.csvwPrefix,
		e.DatasetID,
		e.Edition,
		e.Version,
	)

	if isCustom {
		fn = custom
	}

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

func (h *CantabularMetadataExport) produceOutputMessage(ctx context.Context, e *event.CSVCreated) error {
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
	h.producer.Channels().Output <- kafka.BytesMessage{Value: b, Context: ctx}

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

func (h *CantabularMetadataExport) GetFilterDimensions(ctx context.Context, filterOutput filter.Model) ([]dataset.VersionDimension, error) {
	var areaType string
	for i := range filterOutput.Dimensions {
		if filterOutput.Dimensions[i].IsAreaType != nil && *filterOutput.Dimensions[i].IsAreaType {
			areaType = filterOutput.Dimensions[i].ID
		}
	}

	cReq := cantabular.GetDimensionsByNameRequest{
		Dataset: filterOutput.PopulationType,
	}
	for i := range filterOutput.Dimensions {
		cReq.DimensionNames = append(cReq.DimensionNames, filterOutput.Dimensions[i].ID)
	}
	resp, err := h.ctblr.GetDimensionsByName(ctx, cReq)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query dimensions")
	}

	dims := []dataset.VersionDimension{}
	for i := range resp.Dataset.Variables.Edges {
		isAreaType := resp.Dataset.Variables.Edges[i].Node.Name == areaType
		dim := dataset.VersionDimension{
			ID:                   resp.Dataset.Variables.Edges[i].Node.Name,
			Name:                 resp.Dataset.Variables.Edges[i].Node.Name,
			Description:          resp.Dataset.Variables.Edges[i].Node.Description,
			Label:                resp.Dataset.Variables.Edges[i].Node.Label,
			NumberOfOptions:      resp.Dataset.Variables.Edges[i].Node.Categories.TotalCount,
			IsAreaType:           &isAreaType,
			QualityStatementText: resp.Dataset.Variables.Edges[i].Node.Meta.ONSVariable.QualityStatementText,
		}
		dims = append(dims, dim)
	}

	return dims, nil
}

func removeDuplicateDimensions(vDims []dataset.VersionDimension) []dataset.VersionDimension {
	allDims := make(map[string]bool)
	dimensions := []dataset.VersionDimension{}
	for i := range vDims {
		if !allDims[vDims[i].ID] {
			allDims[vDims[i].ID] = true
			dimensions = append(dimensions, vDims[i])
		}
	}
	return dimensions
}
