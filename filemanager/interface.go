package filemanager

import (
	"context"
	"io"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Uploader interface {
	Get(ctx context.Context, key string) (io.ReadCloser, *int64, error)
	Upload(ctx context.Context, input *s3.PutObjectInput, options ...func(*manager.Uploader)) (*manager.UploadOutput, error)
	UploadWithPSK(ctx context.Context, input *s3.PutObjectInput, psk []byte) (*manager.UploadOutput, error)
	BucketName() string
	Checker(context.Context, *healthcheck.CheckState) error
}

type VaultClient interface {
	WriteKey(path, key, value string) error
	Checker(context.Context, *healthcheck.CheckState) error
}

type Generator interface {
	NewPSK() ([]byte, error)
}
