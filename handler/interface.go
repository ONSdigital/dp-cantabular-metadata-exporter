package handler

import (
	"io"
)

type FileManager interface {
	Upload(body io.Reader, bucket, filename string) (string, error)
	UploadEncrypted(body io.Reader, bucket, filename, vaultPath string) (string, error)
}

type DatasetAPIClient interface {}
