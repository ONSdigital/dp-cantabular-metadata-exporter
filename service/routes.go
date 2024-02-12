package service

import (
	"github.com/go-chi/chi/v5"
	"github.com/riandyrn/otelchi"
)

// BuildRoutes builds the routing for the API
func (svc *Service) BuildRoutesOtel(otServiceName string) {
	r := chi.NewRouter()
	r.Use(otelchi.Middleware(otServiceName))

	// Healthcheck
	r.HandleFunc("/health", svc.HealthCheck.Handler)
	svc.router = r
}

func (svc *Service) BuildRoutes() {
	r := chi.NewRouter()

	// Healthcheck
	r.HandleFunc("/health", svc.HealthCheck.Handler)
	svc.router = r
}
