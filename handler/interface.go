package handler

import (
	"context"
	"io"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
)

type FileManager interface {
	Upload(body io.Reader, filename string) (string, error)
	UploadPrivate(body io.Reader, filename, vaultPath string) (string, error)
}

type S3Uploader interface{}

type DatasetAPIClient interface {
	GetVersion(ctx context.Context, userAuthToken, serviceAuthToken, downloadServiceAuthToken, collectionID, datasetID, edition, version string) (dataset.Version, error)
	GetVersionMetadata(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, id, edition, version string) (dataset.Metadata, error)
	GetVersionDimensions(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, id, edition, version string) (dataset.VersionDimensions, error)
	GetOptions(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, id, edition, version, dimension string, q *dataset.QueryParams) (dataset.Options, error)
	PutVersion(ctx context.Context, usrAuthToken, svcAuthToken, collectionID, datasetID, edition, ver string, v dataset.Version) error
}
