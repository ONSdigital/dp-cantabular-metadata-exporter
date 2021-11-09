package filemanager

import (
	"encoding/hex"
	"fmt"
	"io"
	"net/url"

	dps3 "github.com/ONSdigital/dp-s3"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type FileManager struct {
	cfg       Config
	s3public  S3Uploader
	s3private S3Uploader
	vault     VaultClient
	generator Generator
	vaultKey  string
	publicURL string
}

func New(cfg Config, s *session.Session, v VaultClient, g Generator) *FileManager {
	s3pub := dps3.NewUploaderWithSession(cfg.PublicBucket, s)
	s3priv := dps3.NewUploaderWithSession(cfg.PrivateBucket, s)

	return &FileManager{
		s3public:  s3pub,
		s3private: s3priv,
		vault:     v,
		generator: g,
		vaultKey:  cfg.VaultKey,
		publicURL: cfg.PublicURL,
	}
}

func (f *FileManager) Upload(body io.Reader, filename string) (string, error) {
	bucket := f.s3public.BucketName()
	result, err := f.s3public.Upload(&s3manager.UploadInput{
		Body:   body,
		Bucket: &bucket,
		Key:    &filename,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	return url.PathUnescape(result.Location)
}

func (f *FileManager) UploadPrivate(body io.Reader, filename, vaultPath string) (string, error) {
	psk, err := f.generator.NewPSK()
	if err != nil {
		return "", fmt.Errorf("failed to generate PSK: %w", err)
	}

	// vaultPath := fmt.Sprintf("%s/%s.csv", vaultPathRoot, instanceID, h.vaultPath, instanceID)
	// ^^ leaving comment here as reminder of how to implement in handler ^^

	if err := f.vault.WriteKey(vaultPath, f.vaultKey, hex.EncodeToString(psk)); err != nil {
		return "", fmt.Errorf("failed to write key to vault: %w", err)
	}

	bucket := f.s3private.BucketName()
	result, err := f.s3private.UploadWithPSK(&s3manager.UploadInput{
		Body:   body,
		Bucket: &bucket,
		Key:    &filename,
	}, psk)
	if err != nil {
		return "", fmt.Errorf("failed to upload encrypted file to S3: %w", err)
	}

	return url.PathUnescape(result.Location)
}
