package handler

import (
	"context"
	"io"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
)

//go:generate moq -out mock/filemanager.go        -pkg mock . FileManager
//go:generate moq -out mock/dataset_api_client.go -pkg mock . DatasetAPIClient

type FileManager interface {
	Upload(body io.Reader, filename string) (string, error)
	UploadPrivate(body io.Reader, filename, vaultPath string) (string, error)
}

type DatasetAPIClient interface {
	GetVersion(ctx context.Context, userAuthToken, serviceAuthToken, downloadServiceAuthToken, collectionID, datasetID, edition, version string) (dataset.Version, error)
	GetVersionMetadataSelection(ctx context.Context, req dataset.GetVersionMetadataSelectionInput) (*dataset.Metadata, error)
	GetVersionMetadata(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, id, edition, version string) (dataset.Metadata, error)
	GetVersionDimensions(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, id, edition, version string) (dataset.VersionDimensions, error)
	GetOptionsInBatches(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, id, edition, version, dimension string, batchSize, maxWorkers int) (dataset.Options, error)
	PutVersion(ctx context.Context, usrAuthToken, svcAuthToken, collectionID, datasetID, edition, ver string, v dataset.Version) error
	GetMetadataURL(id, edition, version string) (url string)
}
