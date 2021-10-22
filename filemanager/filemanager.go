package filemanager

import (
	"encoding/hex"
	"io"
	"fmt"
	"net/url"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type FileManager struct{
	cfg       Config 
	s3        S3Uploader
	vault     VaultClient
	generator Generator
}

func New(cfg Config, s3 S3Uploader, v VaultClient, g Generator) *FileManager {
	return &FileManager{
		cfg:       cfg,
		s3:        s3,
		vault:     v,
		generator: g,
	}
}

func (f *FileManager) Upload(body io.Reader, bucket, filename string) (string, error) {
	result, err := f.s3.Upload(&s3manager.UploadInput{
		Body:   body,
		Bucket: &bucket,
		Key:    &filename,
	})
	if err != nil{
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	return url.PathUnescape(result.Location)
}

func (f *FileManager) UploadEncrypted(body io.Reader, bucket, filename, vaultPath string) (string, error) {
	psk, err := f.generator.NewPSK()
	if err != nil {
		return "", fmt.Errorf("failed to generate PSK: %w", err)
	}

	// vaultPath := fmt.Sprintf("%s/%s.csv", vaultPathRoot, instanceID, h.vaultPath, instanceID)
	// ^^ leaving comment here as reminder of how to implement in handler ^^

	if err := f.vault.WriteKey(vaultPath, f.cfg.VaultKey, hex.EncodeToString(psk)); err != nil {
		return "", fmt.Errorf("failed to write key to vault: %w", err)
	}

	result, err := f.s3.UploadWithPSK(&s3manager.UploadInput{
		Body:   body,
		Bucket: &bucket,
		Key:    &filename,
	}, psk)
	if err != nil {
		return "", fmt.Errorf("failed to upload encrypted file to S3: %w", err)
	}

	return url.PathUnescape(result.Location)
}
