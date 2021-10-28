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
	GetVersionMetadata(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, id, edition, version string) (m dataset.Metadata, err error)
	GetVersionDimensions(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, id, edition, version string) (m dataset.VersionDimensions, err error)
	GetOptions(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, id, edition, version, dimension string, q *dataset.QueryParams) (m dataset.Options, err error)
}
