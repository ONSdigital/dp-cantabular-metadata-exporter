// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"context"
	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/handler"
	"sync"
)

// Ensure, that DatasetAPIClientMock does implement handler.DatasetAPIClient.
// If this is not the case, regenerate this file with moq.
var _ handler.DatasetAPIClient = &DatasetAPIClientMock{}

// DatasetAPIClientMock is a mock implementation of handler.DatasetAPIClient.
//
// 	func TestSomethingThatUsesDatasetAPIClient(t *testing.T) {
//
// 		// make and configure a mocked handler.DatasetAPIClient
// 		mockedDatasetAPIClient := &DatasetAPIClientMock{
// 			GetMetadataURLFunc: func(id string, edition string, version string) string {
// 				panic("mock out the GetMetadataURL method")
// 			},
// 			GetOptionsInBatchesFunc: func(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, id string, edition string, version string, dimension string, batchSize int, maxWorkers int) (dataset.Options, error) {
// 				panic("mock out the GetOptionsInBatches method")
// 			},
// 			GetVersionFunc: func(ctx context.Context, userAuthToken string, serviceAuthToken string, downloadServiceAuthToken string, collectionID string, datasetID string, edition string, version string) (dataset.Version, error) {
// 				panic("mock out the GetVersion method")
// 			},
// 			GetVersionDimensionsFunc: func(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, id string, edition string, version string) (dataset.VersionDimensions, error) {
// 				panic("mock out the GetVersionDimensions method")
// 			},
// 			GetVersionMetadataFunc: func(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, id string, edition string, version string) (dataset.Metadata, error) {
// 				panic("mock out the GetVersionMetadata method")
// 			},
// 			PutVersionFunc: func(ctx context.Context, usrAuthToken string, svcAuthToken string, collectionID string, datasetID string, edition string, ver string, v dataset.Version) error {
// 				panic("mock out the PutVersion method")
// 			},
// 		}
//
// 		// use mockedDatasetAPIClient in code that requires handler.DatasetAPIClient
// 		// and then make assertions.
//
// 	}
type DatasetAPIClientMock struct {
	// GetMetadataURLFunc mocks the GetMetadataURL method.
	GetMetadataURLFunc func(id string, edition string, version string) string

	// GetOptionsInBatchesFunc mocks the GetOptionsInBatches method.
	GetOptionsInBatchesFunc func(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, id string, edition string, version string, dimension string, batchSize int, maxWorkers int) (dataset.Options, error)

	// GetVersionFunc mocks the GetVersion method.
	GetVersionFunc func(ctx context.Context, userAuthToken string, serviceAuthToken string, downloadServiceAuthToken string, collectionID string, datasetID string, edition string, version string) (dataset.Version, error)

	// GetVersionDimensionsFunc mocks the GetVersionDimensions method.
	GetVersionDimensionsFunc func(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, id string, edition string, version string) (dataset.VersionDimensions, error)

	// GetVersionMetadataFunc mocks the GetVersionMetadata method.
	GetVersionMetadataFunc func(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, id string, edition string, version string) (dataset.Metadata, error)

	// PutVersionFunc mocks the PutVersion method.
	PutVersionFunc func(ctx context.Context, usrAuthToken string, svcAuthToken string, collectionID string, datasetID string, edition string, ver string, v dataset.Version) error

	// calls tracks calls to the methods.
	calls struct {
		// GetMetadataURL holds details about calls to the GetMetadataURL method.
		GetMetadataURL []struct {
			// ID is the id argument value.
			ID string
			// Edition is the edition argument value.
			Edition string
			// Version is the version argument value.
			Version string
		}
		// GetOptionsInBatches holds details about calls to the GetOptionsInBatches method.
		GetOptionsInBatches []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// UserAuthToken is the userAuthToken argument value.
			UserAuthToken string
			// ServiceAuthToken is the serviceAuthToken argument value.
			ServiceAuthToken string
			// CollectionID is the collectionID argument value.
			CollectionID string
			// ID is the id argument value.
			ID string
			// Edition is the edition argument value.
			Edition string
			// Version is the version argument value.
			Version string
			// Dimension is the dimension argument value.
			Dimension string
			// BatchSize is the batchSize argument value.
			BatchSize int
			// MaxWorkers is the maxWorkers argument value.
			MaxWorkers int
		}
		// GetVersion holds details about calls to the GetVersion method.
		GetVersion []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// UserAuthToken is the userAuthToken argument value.
			UserAuthToken string
			// ServiceAuthToken is the serviceAuthToken argument value.
			ServiceAuthToken string
			// DownloadServiceAuthToken is the downloadServiceAuthToken argument value.
			DownloadServiceAuthToken string
			// CollectionID is the collectionID argument value.
			CollectionID string
			// DatasetID is the datasetID argument value.
			DatasetID string
			// Edition is the edition argument value.
			Edition string
			// Version is the version argument value.
			Version string
		}
		// GetVersionDimensions holds details about calls to the GetVersionDimensions method.
		GetVersionDimensions []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// UserAuthToken is the userAuthToken argument value.
			UserAuthToken string
			// ServiceAuthToken is the serviceAuthToken argument value.
			ServiceAuthToken string
			// CollectionID is the collectionID argument value.
			CollectionID string
			// ID is the id argument value.
			ID string
			// Edition is the edition argument value.
			Edition string
			// Version is the version argument value.
			Version string
		}
		// GetVersionMetadata holds details about calls to the GetVersionMetadata method.
		GetVersionMetadata []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// UserAuthToken is the userAuthToken argument value.
			UserAuthToken string
			// ServiceAuthToken is the serviceAuthToken argument value.
			ServiceAuthToken string
			// CollectionID is the collectionID argument value.
			CollectionID string
			// ID is the id argument value.
			ID string
			// Edition is the edition argument value.
			Edition string
			// Version is the version argument value.
			Version string
		}
		// PutVersion holds details about calls to the PutVersion method.
		PutVersion []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// UsrAuthToken is the usrAuthToken argument value.
			UsrAuthToken string
			// SvcAuthToken is the svcAuthToken argument value.
			SvcAuthToken string
			// CollectionID is the collectionID argument value.
			CollectionID string
			// DatasetID is the datasetID argument value.
			DatasetID string
			// Edition is the edition argument value.
			Edition string
			// Ver is the ver argument value.
			Ver string
			// V is the v argument value.
			V dataset.Version
		}
	}
	lockGetMetadataURL       sync.RWMutex
	lockGetOptionsInBatches  sync.RWMutex
	lockGetVersion           sync.RWMutex
	lockGetVersionDimensions sync.RWMutex
	lockGetVersionMetadata   sync.RWMutex
	lockPutVersion           sync.RWMutex
}

// GetMetadataURL calls GetMetadataURLFunc.
func (mock *DatasetAPIClientMock) GetMetadataURL(id string, edition string, version string) string {
	if mock.GetMetadataURLFunc == nil {
		panic("DatasetAPIClientMock.GetMetadataURLFunc: method is nil but DatasetAPIClient.GetMetadataURL was just called")
	}
	callInfo := struct {
		ID      string
		Edition string
		Version string
	}{
		ID:      id,
		Edition: edition,
		Version: version,
	}
	mock.lockGetMetadataURL.Lock()
	mock.calls.GetMetadataURL = append(mock.calls.GetMetadataURL, callInfo)
	mock.lockGetMetadataURL.Unlock()
	return mock.GetMetadataURLFunc(id, edition, version)
}

// GetMetadataURLCalls gets all the calls that were made to GetMetadataURL.
// Check the length with:
//     len(mockedDatasetAPIClient.GetMetadataURLCalls())
func (mock *DatasetAPIClientMock) GetMetadataURLCalls() []struct {
	ID      string
	Edition string
	Version string
} {
	var calls []struct {
		ID      string
		Edition string
		Version string
	}
	mock.lockGetMetadataURL.RLock()
	calls = mock.calls.GetMetadataURL
	mock.lockGetMetadataURL.RUnlock()
	return calls
}

// GetOptionsInBatches calls GetOptionsInBatchesFunc.
func (mock *DatasetAPIClientMock) GetOptionsInBatches(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, id string, edition string, version string, dimension string, batchSize int, maxWorkers int) (dataset.Options, error) {
	if mock.GetOptionsInBatchesFunc == nil {
		panic("DatasetAPIClientMock.GetOptionsInBatchesFunc: method is nil but DatasetAPIClient.GetOptionsInBatches was just called")
	}
	callInfo := struct {
		Ctx              context.Context
		UserAuthToken    string
		ServiceAuthToken string
		CollectionID     string
		ID               string
		Edition          string
		Version          string
		Dimension        string
		BatchSize        int
		MaxWorkers       int
	}{
		Ctx:              ctx,
		UserAuthToken:    userAuthToken,
		ServiceAuthToken: serviceAuthToken,
		CollectionID:     collectionID,
		ID:               id,
		Edition:          edition,
		Version:          version,
		Dimension:        dimension,
		BatchSize:        batchSize,
		MaxWorkers:       maxWorkers,
	}
	mock.lockGetOptionsInBatches.Lock()
	mock.calls.GetOptionsInBatches = append(mock.calls.GetOptionsInBatches, callInfo)
	mock.lockGetOptionsInBatches.Unlock()
	return mock.GetOptionsInBatchesFunc(ctx, userAuthToken, serviceAuthToken, collectionID, id, edition, version, dimension, batchSize, maxWorkers)
}

// GetOptionsInBatchesCalls gets all the calls that were made to GetOptionsInBatches.
// Check the length with:
//     len(mockedDatasetAPIClient.GetOptionsInBatchesCalls())
func (mock *DatasetAPIClientMock) GetOptionsInBatchesCalls() []struct {
	Ctx              context.Context
	UserAuthToken    string
	ServiceAuthToken string
	CollectionID     string
	ID               string
	Edition          string
	Version          string
	Dimension        string
	BatchSize        int
	MaxWorkers       int
} {
	var calls []struct {
		Ctx              context.Context
		UserAuthToken    string
		ServiceAuthToken string
		CollectionID     string
		ID               string
		Edition          string
		Version          string
		Dimension        string
		BatchSize        int
		MaxWorkers       int
	}
	mock.lockGetOptionsInBatches.RLock()
	calls = mock.calls.GetOptionsInBatches
	mock.lockGetOptionsInBatches.RUnlock()
	return calls
}

// GetVersion calls GetVersionFunc.
func (mock *DatasetAPIClientMock) GetVersion(ctx context.Context, userAuthToken string, serviceAuthToken string, downloadServiceAuthToken string, collectionID string, datasetID string, edition string, version string) (dataset.Version, error) {
	if mock.GetVersionFunc == nil {
		panic("DatasetAPIClientMock.GetVersionFunc: method is nil but DatasetAPIClient.GetVersion was just called")
	}
	callInfo := struct {
		Ctx                      context.Context
		UserAuthToken            string
		ServiceAuthToken         string
		DownloadServiceAuthToken string
		CollectionID             string
		DatasetID                string
		Edition                  string
		Version                  string
	}{
		Ctx:                      ctx,
		UserAuthToken:            userAuthToken,
		ServiceAuthToken:         serviceAuthToken,
		DownloadServiceAuthToken: downloadServiceAuthToken,
		CollectionID:             collectionID,
		DatasetID:                datasetID,
		Edition:                  edition,
		Version:                  version,
	}
	mock.lockGetVersion.Lock()
	mock.calls.GetVersion = append(mock.calls.GetVersion, callInfo)
	mock.lockGetVersion.Unlock()
	return mock.GetVersionFunc(ctx, userAuthToken, serviceAuthToken, downloadServiceAuthToken, collectionID, datasetID, edition, version)
}

// GetVersionCalls gets all the calls that were made to GetVersion.
// Check the length with:
//     len(mockedDatasetAPIClient.GetVersionCalls())
func (mock *DatasetAPIClientMock) GetVersionCalls() []struct {
	Ctx                      context.Context
	UserAuthToken            string
	ServiceAuthToken         string
	DownloadServiceAuthToken string
	CollectionID             string
	DatasetID                string
	Edition                  string
	Version                  string
} {
	var calls []struct {
		Ctx                      context.Context
		UserAuthToken            string
		ServiceAuthToken         string
		DownloadServiceAuthToken string
		CollectionID             string
		DatasetID                string
		Edition                  string
		Version                  string
	}
	mock.lockGetVersion.RLock()
	calls = mock.calls.GetVersion
	mock.lockGetVersion.RUnlock()
	return calls
}

// GetVersionDimensions calls GetVersionDimensionsFunc.
func (mock *DatasetAPIClientMock) GetVersionDimensions(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, id string, edition string, version string) (dataset.VersionDimensions, error) {
	if mock.GetVersionDimensionsFunc == nil {
		panic("DatasetAPIClientMock.GetVersionDimensionsFunc: method is nil but DatasetAPIClient.GetVersionDimensions was just called")
	}
	callInfo := struct {
		Ctx              context.Context
		UserAuthToken    string
		ServiceAuthToken string
		CollectionID     string
		ID               string
		Edition          string
		Version          string
	}{
		Ctx:              ctx,
		UserAuthToken:    userAuthToken,
		ServiceAuthToken: serviceAuthToken,
		CollectionID:     collectionID,
		ID:               id,
		Edition:          edition,
		Version:          version,
	}
	mock.lockGetVersionDimensions.Lock()
	mock.calls.GetVersionDimensions = append(mock.calls.GetVersionDimensions, callInfo)
	mock.lockGetVersionDimensions.Unlock()
	return mock.GetVersionDimensionsFunc(ctx, userAuthToken, serviceAuthToken, collectionID, id, edition, version)
}

// GetVersionDimensionsCalls gets all the calls that were made to GetVersionDimensions.
// Check the length with:
//     len(mockedDatasetAPIClient.GetVersionDimensionsCalls())
func (mock *DatasetAPIClientMock) GetVersionDimensionsCalls() []struct {
	Ctx              context.Context
	UserAuthToken    string
	ServiceAuthToken string
	CollectionID     string
	ID               string
	Edition          string
	Version          string
} {
	var calls []struct {
		Ctx              context.Context
		UserAuthToken    string
		ServiceAuthToken string
		CollectionID     string
		ID               string
		Edition          string
		Version          string
	}
	mock.lockGetVersionDimensions.RLock()
	calls = mock.calls.GetVersionDimensions
	mock.lockGetVersionDimensions.RUnlock()
	return calls
}

// GetVersionMetadata calls GetVersionMetadataFunc.
func (mock *DatasetAPIClientMock) GetVersionMetadata(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, id string, edition string, version string) (dataset.Metadata, error) {
	if mock.GetVersionMetadataFunc == nil {
		panic("DatasetAPIClientMock.GetVersionMetadataFunc: method is nil but DatasetAPIClient.GetVersionMetadata was just called")
	}
	callInfo := struct {
		Ctx              context.Context
		UserAuthToken    string
		ServiceAuthToken string
		CollectionID     string
		ID               string
		Edition          string
		Version          string
	}{
		Ctx:              ctx,
		UserAuthToken:    userAuthToken,
		ServiceAuthToken: serviceAuthToken,
		CollectionID:     collectionID,
		ID:               id,
		Edition:          edition,
		Version:          version,
	}
	mock.lockGetVersionMetadata.Lock()
	mock.calls.GetVersionMetadata = append(mock.calls.GetVersionMetadata, callInfo)
	mock.lockGetVersionMetadata.Unlock()
	return mock.GetVersionMetadataFunc(ctx, userAuthToken, serviceAuthToken, collectionID, id, edition, version)
}

// GetVersionMetadataCalls gets all the calls that were made to GetVersionMetadata.
// Check the length with:
//     len(mockedDatasetAPIClient.GetVersionMetadataCalls())
func (mock *DatasetAPIClientMock) GetVersionMetadataCalls() []struct {
	Ctx              context.Context
	UserAuthToken    string
	ServiceAuthToken string
	CollectionID     string
	ID               string
	Edition          string
	Version          string
} {
	var calls []struct {
		Ctx              context.Context
		UserAuthToken    string
		ServiceAuthToken string
		CollectionID     string
		ID               string
		Edition          string
		Version          string
	}
	mock.lockGetVersionMetadata.RLock()
	calls = mock.calls.GetVersionMetadata
	mock.lockGetVersionMetadata.RUnlock()
	return calls
}

// PutVersion calls PutVersionFunc.
func (mock *DatasetAPIClientMock) PutVersion(ctx context.Context, usrAuthToken string, svcAuthToken string, collectionID string, datasetID string, edition string, ver string, v dataset.Version) error {
	if mock.PutVersionFunc == nil {
		panic("DatasetAPIClientMock.PutVersionFunc: method is nil but DatasetAPIClient.PutVersion was just called")
	}
	callInfo := struct {
		Ctx          context.Context
		UsrAuthToken string
		SvcAuthToken string
		CollectionID string
		DatasetID    string
		Edition      string
		Ver          string
		V            dataset.Version
	}{
		Ctx:          ctx,
		UsrAuthToken: usrAuthToken,
		SvcAuthToken: svcAuthToken,
		CollectionID: collectionID,
		DatasetID:    datasetID,
		Edition:      edition,
		Ver:          ver,
		V:            v,
	}
	mock.lockPutVersion.Lock()
	mock.calls.PutVersion = append(mock.calls.PutVersion, callInfo)
	mock.lockPutVersion.Unlock()
	return mock.PutVersionFunc(ctx, usrAuthToken, svcAuthToken, collectionID, datasetID, edition, ver, v)
}

// PutVersionCalls gets all the calls that were made to PutVersion.
// Check the length with:
//     len(mockedDatasetAPIClient.PutVersionCalls())
func (mock *DatasetAPIClientMock) PutVersionCalls() []struct {
	Ctx          context.Context
	UsrAuthToken string
	SvcAuthToken string
	CollectionID string
	DatasetID    string
	Edition      string
	Ver          string
	V            dataset.Version
} {
	var calls []struct {
		Ctx          context.Context
		UsrAuthToken string
		SvcAuthToken string
		CollectionID string
		DatasetID    string
		Edition      string
		Ver          string
		V            dataset.Version
	}
	mock.lockPutVersion.RLock()
	calls = mock.calls.PutVersion
	mock.lockPutVersion.RUnlock()
	return calls
}
