package service

import (
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/handler"

	"github.com/go-chi/chi/v5"
)

// BuildRoutes builds the routing for the API
func (svc *Service) BuildRoutes() {
	r := chi.NewRouter()

	metadata := handler.NewMetadata()
	r.Post("/metadata", metadata.Post)

	// Healthcheck
	r.HandleFunc("/health", svc.healthCheck.Handler)

	svc.router = r
}
