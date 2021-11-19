package service

import (
	"github.com/go-chi/chi/v5"
)

// BuildRoutes builds the routing for the API
func (svc *Service) BuildRoutes() {
	r := chi.NewRouter()

	// Healthcheck
	r.HandleFunc("/health", svc.healthCheck.Handler)

	svc.router = r
}
