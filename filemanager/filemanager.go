package filemanager

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"net/url"

	dps3 "github.com/ONSdigital/dp-s3/v3"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type FileManager struct {
	s3public  *dps3.Client
	s3private *dps3.Client
	vault     VaultClient
	generator Generator
	vaultKey  string
	publicURL string
}

func New(cfg Config, v VaultClient, g Generator, s3pub, s3priv *dps3.Client) *FileManager {
	return &FileManager{
		s3public:  s3pub,
		s3private: s3priv,
		vault:     v,
		generator: g,
		vaultKey:  cfg.VaultKey,
		publicURL: cfg.PublicURL,
	}
}

// PublicUploader returns the public S3 uploader
func (f *FileManager) PublicUploader() S3Uploader {
	return f.s3public
}

// PrivateUploader returns the private S3 uploader
func (f *FileManager) PrivateUploader() S3Uploader {
	return f.s3private
}

func (f *FileManager) Upload(ctx context.Context, body io.Reader, filename string) (string, error) {
	bucket := f.s3public.BucketName()
	result, err := f.s3public.Upload(ctx, &s3.PutObjectInput{
		Body:   body,
		Bucket: &bucket,
		Key:    &filename,
	})
	if err != nil {
		return "", Error{
			err: fmt.Errorf("failed to upload to S3: %w", err),
			logData: map[string]interface{}{
				"filename":    filename,
				"bucket_name": bucket,
			},
		}
	}

	return url.PathUnescape(result.Location)
}

func (f *FileManager) UploadPrivate(ctx context.Context, body io.Reader, filename, vaultPath string) (string, error) {
	psk, err := f.generator.NewPSK()
	if err != nil {
		return "", fmt.Errorf("failed to generate PSK: %w", err)
	}

	if err := f.vault.WriteKey(vaultPath, f.vaultKey, hex.EncodeToString(psk)); err != nil {
		return "", Error{
			err: fmt.Errorf("failed to write key to vault: %w", err),
			logData: map[string]interface{}{
				"vault_path": vaultPath,
				"vault_key":  f.vaultKey,
				"psk":        hex.EncodeToString(psk),
			},
		}
	}

	bucket := f.s3private.BucketName()
	result, err := f.s3private.UploadWithPSK(ctx, &s3.PutObjectInput{
		Body:   body,
		Bucket: &bucket,
		Key:    &filename,
	}, psk)
	if err != nil {
		return "", Error{
			err: fmt.Errorf("failed to upload encrypted file to S3: %w", err),
			logData: map[string]interface{}{
				"filename":    filename,
				"bucket_name": bucket,
			},
		}
	}

	return url.PathUnescape(result.Location)
}
