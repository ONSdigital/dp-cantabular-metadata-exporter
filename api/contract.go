package api

import ()

// ErrorResponse is the generic ONS error response for HTTP errors
type ErrorResponse struct {
	Errors []string `json:"errors"`
}

// GetHelloResponse is the response for GET /hello
type GetHelloResponse struct {
	Message string `json:"message,omitempty"`
}

// ExportMetadataRequest is the request for POST /metadata
type ExportMetadataRequest struct {}