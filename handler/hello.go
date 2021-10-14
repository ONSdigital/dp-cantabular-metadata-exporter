package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-cantabular-metadata-exporter/api"
	"github.com/ONSdigital/log.go/v2/log"

	"github.com/go-chi/chi/v5"
)

// Hello is the handler struct which holds dependencies for requests to /hello
type Hello struct {
	message string
}

// NewHello returns a new Hello Handler
func NewHello(message string) *Hello {
	return &Hello{
		message: message,
	}
}

// Get handles HTTP requests for GET /hello. Includes examples
// of returning regular JSON response, regular Go error and custom
// error type with defined response status code. See api/respond_test.go
// for expected behaviour.
func (h *Hello) Get(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	name := chi.URLParam(req, "name")

	if len(h.message) == 0 {
		api.RespondError(ctx, w, errors.New("response message length 0"))
		return
	}

	if name == "dave" {
		api.RespondError(ctx, w, Error{
			err:        errors.New("inappropriate name"),
			resp:       "you're my wife now",
			statusCode: http.StatusTeapot,
			logData: log.Data{
				"hello": "dave",
			},
		})
		return
	}

	resp := api.GetHelloResponse{
		Message: fmt.Sprintf("%s says: %s", name, h.message),
	}

	api.RespondJSON(ctx, w, http.StatusOK, resp)
}
