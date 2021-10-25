package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ONSdigital/dp-cantabular-metadata-exporter/api"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/schema"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/event"

	"github.com/ONSdigital/log.go/v2/log"
	kafka "github.com/ONSdigital/dp-kafka/v2"
)

// Metadata is the handler struct which holds dependencies for requests to /metadata
type Metadata struct{
	producer kafka.IProducer
}

// NewMetadata returns a new Metadata Handler
func NewMetadata(p kafka.IProducer) *Metadata {
	return &Metadata{
		producer: p,
	}
}

// Post handles HTTP requests for POST /metadata.
func (h *Metadata) Post(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var logData log.Data
	defer r.Body.Close()

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

	b, err = schema.CantabularMetadataExport.Marshal(&event.CantabularMetadataExport{
		DatasetID: req.DatasetID,
		Edition:   req.Edition,
		Version:   req.Version,
	})
	if err != nil {
		api.Error(ctx, w, Error{
			err:        fmt.Errorf("error marshalling avro event: %w", err),
			statusCode: http.StatusInternalServerError,
			logData:    logData,
		})
	}

	// Send bytes to kafka producer output channel
	h.producer.Channels().Output <- b

	w.WriteHeader(http.StatusAccepted)
}
