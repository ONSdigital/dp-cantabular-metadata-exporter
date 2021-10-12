package api

import (
	"github.com/go-chi/chi/v5"
)

// Setup function sets up the api and returns an api
func Init() chi.Router {
	r := chi.NewRouter()

	r.Get("/hello", HelloHandler)

	return r
}
