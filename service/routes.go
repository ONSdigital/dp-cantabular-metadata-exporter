package service

import (
	"fmt"

	"github.com/ONSdigital/dp-cantabular-metadata-exporter/handler"

	"github.com/go-chi/chi/v5"
)

// BuildRoutes builds the routing for the API
func (svc *Service) BuildRoutes() {
	r := chi.NewRouter()

	hello := handler.NewHello(fmt.Sprintf("Hello, I am listening on %s!", svc.config.BindAddr))

	r.Get("/hello/{name}", hello.Get)

	// Healthcheck
	r.HandleFunc("/health", svc.healthCheck.Handler)

	svc.router = r
}
