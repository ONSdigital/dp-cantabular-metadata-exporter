package service

import (
	"context"
	"io"
	"net/http"

	"github.com/ONSdigital/dp-cantabular-metadata-exporter/event"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	kafka "github.com/ONSdigital/dp-kafka/v2"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

//go:generate moq -out mock/server.go -pkg mock . HTTPServer
//go:generate moq -out mock/healthCheck.go -pkg mock . HealthChecker
//go:generate moq -out mock/processor.go -pkg mock . Processor

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
	AddCheck(name string, checker healthcheck.Checker) error
}

// Processor defines the required methods for the Processor object
type Processor interface {
	Consume(context.Context, kafka.IConsumerGroup, event.Handler)
}

type DatasetAPIClient interface {
	GetInstance(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, instanceID, ifMatch string) (m dataset.Instance, eTag string, err error)
	PutInstance(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, instanceID string, instanceUpdate dataset.UpdateInstance, ifMatch string) (eTag string, err error)
	GetVersion(ctx context.Context, userAuthToken, serviceAuthToken, downloadServiceAuthToken, collectionID, datasetID, edition, version string) (dataset.Version, error)
	GetVersionMetadata(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, id, edition, version string) (dataset.Metadata, error)
	GetVersionDimensions(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, id, edition, version string) (dataset.VersionDimensions, error)
	GetOptions(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, id, edition, version, dimension string, q *dataset.QueryParams) (dataset.Options, error)
	Checker(context.Context, *healthcheck.CheckState) error
	PutVersion(ctx context.Context, usrAuthToken, svcAuthToken, collectionID, datasetID, edition, ver string, v dataset.Version) error
}

type S3Uploader interface {
	Get(key string) (io.ReadCloser, *int64, error)
	Upload(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)
	UploadWithPSK(input *s3manager.UploadInput, psk []byte) (*s3manager.UploadOutput, error)
	BucketName() string
	Checker(context.Context, *healthcheck.CheckState) error
}

type VaultClient interface {
	WriteKey(path, key, value string) error
	Checker(context.Context, *healthcheck.CheckState) error
}

type FileManager interface {
	Upload(body io.Reader, filename string) (string, error)
	UploadPrivate(body io.Reader, filename, vaultPath string) (string, error)
}

// Generator contains methods for dynamically required strings and tokens
// e.g. UUIDs, PSKs.
type Generator interface {
	NewPSK() ([]byte, error)
}
