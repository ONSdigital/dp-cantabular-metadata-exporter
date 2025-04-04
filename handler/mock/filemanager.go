// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"context"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/handler"
	"io"
	"sync"
)

// Ensure, that FileManagerMock does implement handler.FileManager.
// If this is not the case, regenerate this file with moq.
var _ handler.FileManager = &FileManagerMock{}

// FileManagerMock is a mock implementation of handler.FileManager.
//
//	func TestSomethingThatUsesFileManager(t *testing.T) {
//
//		// make and configure a mocked handler.FileManager
//		mockedFileManager := &FileManagerMock{
//			UploadFunc: func(ctx context.Context, body io.Reader, filename string) (string, error) {
//				panic("mock out the Upload method")
//			},
//			UploadPrivateFunc: func(ctx context.Context, body io.Reader, filename string, vaultPath string) (string, error) {
//				panic("mock out the UploadPrivate method")
//			},
//		}
//
//		// use mockedFileManager in code that requires handler.FileManager
//		// and then make assertions.
//
//	}
type FileManagerMock struct {
	// UploadFunc mocks the Upload method.
	UploadFunc func(ctx context.Context, body io.Reader, filename string) (string, error)

	// UploadPrivateFunc mocks the UploadPrivate method.
	UploadPrivateFunc func(ctx context.Context, body io.Reader, filename string, vaultPath string) (string, error)

	// calls tracks calls to the methods.
	calls struct {
		// Upload holds details about calls to the Upload method.
		Upload []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Body is the body argument value.
			Body io.Reader
			// Filename is the filename argument value.
			Filename string
		}
		// UploadPrivate holds details about calls to the UploadPrivate method.
		UploadPrivate []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Body is the body argument value.
			Body io.Reader
			// Filename is the filename argument value.
			Filename string
			// VaultPath is the vaultPath argument value.
			VaultPath string
		}
	}
	lockUpload        sync.RWMutex
	lockUploadPrivate sync.RWMutex
}

// Upload calls UploadFunc.
func (mock *FileManagerMock) Upload(ctx context.Context, body io.Reader, filename string) (string, error) {
	if mock.UploadFunc == nil {
		panic("FileManagerMock.UploadFunc: method is nil but FileManager.Upload was just called")
	}
	callInfo := struct {
		Ctx      context.Context
		Body     io.Reader
		Filename string
	}{
		Ctx:      ctx,
		Body:     body,
		Filename: filename,
	}
	mock.lockUpload.Lock()
	mock.calls.Upload = append(mock.calls.Upload, callInfo)
	mock.lockUpload.Unlock()
	return mock.UploadFunc(ctx, body, filename)
}

// UploadCalls gets all the calls that were made to Upload.
// Check the length with:
//
//	len(mockedFileManager.UploadCalls())
func (mock *FileManagerMock) UploadCalls() []struct {
	Ctx      context.Context
	Body     io.Reader
	Filename string
} {
	var calls []struct {
		Ctx      context.Context
		Body     io.Reader
		Filename string
	}
	mock.lockUpload.RLock()
	calls = mock.calls.Upload
	mock.lockUpload.RUnlock()
	return calls
}

// UploadPrivate calls UploadPrivateFunc.
func (mock *FileManagerMock) UploadPrivate(ctx context.Context, body io.Reader, filename string, vaultPath string) (string, error) {
	if mock.UploadPrivateFunc == nil {
		panic("FileManagerMock.UploadPrivateFunc: method is nil but FileManager.UploadPrivate was just called")
	}
	callInfo := struct {
		Ctx       context.Context
		Body      io.Reader
		Filename  string
		VaultPath string
	}{
		Ctx:       ctx,
		Body:      body,
		Filename:  filename,
		VaultPath: vaultPath,
	}
	mock.lockUploadPrivate.Lock()
	mock.calls.UploadPrivate = append(mock.calls.UploadPrivate, callInfo)
	mock.lockUploadPrivate.Unlock()
	return mock.UploadPrivateFunc(ctx, body, filename, vaultPath)
}

// UploadPrivateCalls gets all the calls that were made to UploadPrivate.
// Check the length with:
//
//	len(mockedFileManager.UploadPrivateCalls())
func (mock *FileManagerMock) UploadPrivateCalls() []struct {
	Ctx       context.Context
	Body      io.Reader
	Filename  string
	VaultPath string
} {
	var calls []struct {
		Ctx       context.Context
		Body      io.Reader
		Filename  string
		VaultPath string
	}
	mock.lockUploadPrivate.RLock()
	calls = mock.calls.UploadPrivate
	mock.lockUploadPrivate.RUnlock()
	return calls
}
