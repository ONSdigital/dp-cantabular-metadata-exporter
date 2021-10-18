package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config represents service configuration for dp-cantabular-metadata-exporter
type Config struct {
	BindAddr                      string        `envconfig:"BIND_ADDR"`
	GracefulShutdownTimeout       time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval           time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout    time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	KafkaAddr                     []string      `envconfig:"KAFKA_ADDR"`
	CantabularMetadataExportTopic string        `envconfig:"CANTABULAR_METADATA_EXPORT_TOPIC"`
	CantabularMetadataExportGroup string        `envconfig:"CANTABULAR_METADATA_EXPORT_GROUP"`
}

var cfg *Config

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		BindAddr:                      ":26700",
		GracefulShutdownTimeout:       5 * time.Second,
		HealthCheckInterval:           30 * time.Second,
		HealthCheckCriticalTimeout:    90 * time.Second,
		KafkaAddr:                     []string{"localhost:9092"},
		CantabularMetadataExportGroup: "cantabular-metadata-export",
		CantabularMetadataExportTopic: "cantabular-metadata-export",
	}

	return cfg, envconfig.Process("", cfg)
}
