package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ONSdigital/log.go/v2/log"
)

// Error responds with a single error, formatted to fit in ONS's desired error
// response structure (essentially an array of errors)
func Error(ctx context.Context, w http.ResponseWriter, err error) {
	log.Error(ctx, "error responding to HTTP request", err, unwrapLogData(err))

	status := statusCode(err)
	msg := errorResponse(err)

	resp := ErrorResponse{
		Errors: []string{msg},
	}

	logData := log.Data{
		"error":       err.Error(),
		"response":    msg,
		"status_code": status,
	}

	b, err := json.Marshal(resp)
	if err != nil {
		log.Error(ctx, "badly formed error response", err, logData)
		http.Error(w, msg, status)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if _, err := w.Write(b); err != nil {
		log.Error(ctx, "failed to write error response", err, logData)
		return
	}

	log.Info(ctx, "returned error response", logData)
}
