package handler

import (
	"io"
	"context"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
)

type FileManager interface {
	Upload(body io.Reader, bucket, filename string) (string, error)
	UploadEncrypted(body io.Reader, bucket, filename, vaultPath string) (string, error)
}

type S3Uploader interface {}

type DatasetAPIClient interface {
	GetVersionMetadata(ctx context.Context, usrAuthToken, svcAuthToken, collectionID, id, edition, ver string) (dataset.Metadata, error)
	PutVersion(ctx context.Context, usrAuthToken, svcAuthToken, collectionID, datasetID, edition, ver string, v dataset.Version) error

	GetMetadataURL(id, edition, version string) (url string)
}
