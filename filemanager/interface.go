package filemanager

import (
	"context"
	"io"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

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

type Generator interface {
	NewPSK() ([]byte, error)
}
