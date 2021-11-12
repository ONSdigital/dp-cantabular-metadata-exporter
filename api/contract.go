package api

// ErrorResponse is the generic ONS error response for HTTP errors
type ErrorResponse struct {
	Errors []string `json:"errors"`
}

// ExportMetadataRequest is the request for POST /metadata
type ExportMetadataRequest struct {
	DatasetID    string `json:"dataset_id"`
	Edition      string `json:"edition"`
	Version      string `json:"version"`
	CollectionID string `json:"collection_id"`
}
