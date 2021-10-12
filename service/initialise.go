package service

import (
	"net/http"
	"fmt"

	"github.com/ONSdigital/dp-cantabular-metadata-exporter/config"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dphttp "github.com/ONSdigital/dp-net/http"
)

// GetHTTPServer creates an http server
func GetHTTPServer(bindAddr string, router http.Handler) HTTPServer {
	s := dphttp.NewServer(bindAddr, router)
	s.HandleOSSignals = false
	return s
}

// GetHealthCheck creates a healthcheck with versionInfo and sets teh HealthCheck flag to true
func GetHealthCheck(cfg *config.Config, buildT, commit, ver string) (HealthChecker, error) {
	versionInfo, err := healthcheck.NewVersionInfo(buildT, commit, ver)
	if err != nil {
		return nil, fmt.Errorf("failed to get version info: %w", err)
	}

	hc := healthcheck.New(versionInfo, cfg.HealthCheckCriticalTimeout, cfg.HealthCheckInterval)
	return &hc, nil
}
