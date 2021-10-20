package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ONSdigital/dp-cantabular-metadata-exporter/api"
	"github.com/ONSdigital/log.go/v2/log"
)

// Metadata is the handler struct which holds dependencies for requests to /metadata
type Metadata struct{}

// NewMetadata returns a new Metadata Handler
func NewMetadata() *Metadata {
	return &Metadata{}
}

// Post handles HTTP requests for POST /metadata.
func (h *Metadata) Post(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()
	logData := log.Data{}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.Error(ctx, w, Error{
			err:        fmt.Errorf("failed to read request body: %w", err),
			statusCode: http.StatusBadRequest,
		})
		return
	}

	logData["request_body"] = string(b)

	var req api.ExportMetadataRequest
	if err := json.Unmarshal(b, &req); err != nil {
		api.Error(ctx, w, Error{
			err:        fmt.Errorf("failed to unmarshal request body: %w", err),
			statusCode: http.StatusBadRequest,
			logData:    logData,
		})
		return
	}

	// Add export job to queue

	w.WriteHeader(http.StatusAccepted)
}
