package api

import ()

// ErrorResponse is the generic ONS error response for HTTP errors
type ErrorResponse struct {
	Errors []string `json:"errors"`
}

// ExportMetadataRequest is the request for POST /metadata
type ExportMetadataRequest struct {}
