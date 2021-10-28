package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config represents service configuration for dp-cantabular-metadata-exporter
type Config struct {
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	DatasetAPIURL              string        `envconfig:"DATASET_API_URL"`
	AWSRegion                  string        `envconfig:"AWS_REGION"`
	PublicBucket               string        `envconfig:"PUBLIC_BUCKET"`
	PrivateBucket              string        `envconfig:"PRIVATE_BUCKET"`
	S3BucketURL                string        `envconfig:"PUBLIC_URL"`
	LocalObjectStore           string        `envconfig:"LOCAL_OBJECT_STORE"`
	MinioAccessKey             string        `envconfig:"MINIO_ACCESS_KEY"`
	MinioSecretKey             string        `envconfig:"MINIO_SECRET_KEY"`
	VaultToken                 string        `envconfig:"VAULT_TOKEN"                   json:"-"`
	VaultAddress               string        `envconfig:"VAULT_ADDR"`
	VaultPath                  string        `envconfig:"VAULT_PATH"`
	EncryptionDisabled         bool          `envconfig:"ENCRYPTION_DISABLED"`
	ServiceAuthToken           string        `envconfig:"SERVICE_AUTH_TOKEN"`
	Kafka                      KafkaConfig
}

// KafkaConfig contains the config required to connect to Kafka
type KafkaConfig struct {
	Addr                          []string `envconfig:"KAFKA_ADDR"                            json:"-"`
	Version                       string   `envconfig:"KAFKA_VERSION"`
	OffsetOldest                  bool     `envconfig:"KAFKA_OFFSET_OLDEST"`
	NumWorkers                    int      `envconfig:"KAFKA_NUM_WORKERS"`
	MaxBytes                      int      `envconfig:"KAFKA_MAX_BYTES"`
	SecProtocol                   string   `envconfig:"KAFKA_SEC_PROTO"`
	SecCACerts                    string   `envconfig:"KAFKA_SEC_CA_CERTS"`
	SecClientKey                  string   `envconfig:"KAFKA_SEC_CLIENT_KEY"                  json:"-"`
	SecClientCert                 string   `envconfig:"KAFKA_SEC_CLIENT_CERT"`
	SecSkipVerify                 bool     `envconfig:"KAFKA_SEC_SKIP_VERIFY"`
	CantabularMetadataExportTopic string   `envconfig:"CANTABULAR_METADATA_EXPORT_TOPIC"`
	CantabularMetadataExportGroup string   `envconfig:"CANTABULAR_METADATA_EXPORT_GROUP"`
}

var cfg *Config

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		BindAddr:                   ":26700",
		GracefulShutdownTimeout:    5 * time.Second,
		HealthCheckInterval:        30 * time.Second,
		HealthCheckCriticalTimeout: 90 * time.Second,
		VaultPath:                  "secret/shared/psk",
		VaultAddress:               "http://localhost:8200",
		VaultToken:                 "",
		PublicBucket:               "cantabular-metadata-export",
		PrivateBucket:              "cantabular-metadata-export",
		S3BucketURL:                "",
		EncryptionDisabled:         false,
		Kafka: KafkaConfig{
			Addr:                          []string{"localhost:9092"},
			Version:                       "1.0.2",
			OffsetOldest:                  true,
			NumWorkers:                    1,
			MaxBytes:                      2000000,
			SecProtocol:                   "",
			SecCACerts:                    "",
			SecClientKey:                  "",
			SecClientCert:                 "",
			SecSkipVerify:                 false,
			CantabularMetadataExportGroup: "cantabular-metadata-export",
			CantabularMetadataExportTopic: "cantabular-metadata-export",
		},
	}

	return cfg, envconfig.Process("", cfg)
}
