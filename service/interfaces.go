package service

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/ONSdigital/dp-cantabular-metadata-exporter/filemanager"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

//go:generate moq -out mock/server.go -pkg mock . HTTPServer
//go:generate moq -out mock/healthCheck.go -pkg mock . HealthChecker

// HTTPServer defines the required methods from the HTTP server
type HTTPServer interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

// HealthChecker defines the required methods from Healthcheck
type HealthChecker interface {
	Handler(w http.ResponseWriter, req *http.Request)
	Start(ctx context.Context)
	Stop()
	AddAndGetCheck(name string, checker healthcheck.Checker) (*healthcheck.Check, error)
	SubscribeAll(s healthcheck.Subscriber)
}

type DatasetAPIClient interface {
	GetInstance(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, instanceID, ifMatch string) (m dataset.Instance, eTag string, err error)
	PutInstance(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, instanceID string, instanceUpdate dataset.UpdateInstance, ifMatch string) (eTag string, err error)
	GetVersion(ctx context.Context, userAuthToken, serviceAuthToken, downloadServiceAuthToken, collectionID, datasetID, edition, version string) (dataset.Version, error)
	GetVersionMetadataSelection(ctx context.Context, req dataset.GetVersionMetadataSelectionInput) (*dataset.Metadata, error)
	GetVersionMetadata(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, id, edition, version string) (dataset.Metadata, error)
	GetVersionDimensions(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, id, edition, version string) (dataset.VersionDimensions, error)
	GetOptionsInBatches(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, id, edition, version, dimension string, batchSize, maxWorkers int) (dataset.Options, error)
	GetMetadataURL(id, edition, version string) (url string)
	Checker(context.Context, *healthcheck.CheckState) error
	PutVersion(ctx context.Context, usrAuthToken, svcAuthToken, collectionID, datasetID, edition, ver string, v dataset.Version) error
}

type FilterAPIClient interface {
	UpdateFilterOutput(ctx context.Context, userAuthToken, serviceAuthToken, downloadServiceToken, filterOutputID string, m *filter.Model) error
	GetOutput(ctx context.Context, userAuthToken, serviceAuthToken, downloadServiceToken, collectionID, filterOutput string) (m filter.Model, err error)
	Checker(context.Context, *healthcheck.CheckState) error
}

type PopulationTypesAPIClient interface {
	GetAreaTypes(ctx context.Context, input population.GetAreaTypesInput) (population.GetAreaTypesResponse, error)
}

type CantabularClient interface {
	GetDimensionsByName(context.Context, cantabular.GetDimensionsByNameRequest) (*cantabular.GetDimensionsResponse, error)
}

type S3Uploader interface {
	Get(key string) (io.ReadCloser, *int64, error)
	Upload(ctx context.Context, input *s3.PutObjectInput, options ...func(*manager.Uploader)) (*manager.UploadOutput, error)
	UploadWithPSK(ctx context.Context, input *s3.PutObjectInput, psk []byte) (*manager.UploadOutput, error)
	BucketName() string
	Checker(context.Context, *healthcheck.CheckState) error
}

type VaultClient interface {
	WriteKey(path, key, value string) error
	Checker(context.Context, *healthcheck.CheckState) error
}

type FileManager interface {
	Upload(ctx context.Context, body io.Reader, filename string) (string, error)
	UploadPrivate(ctx context.Context, body io.Reader, filename, vaultPath string) (string, error)
	PrivateUploader() filemanager.S3Uploader
	PublicUploader() filemanager.S3Uploader
}

// Generator contains methods for dynamically required strings and tokens
// e.g. UUIDs, PSKs.
type Generator interface {
	NewPSK() ([]byte, error)
	Timestamp() time.Time
}
