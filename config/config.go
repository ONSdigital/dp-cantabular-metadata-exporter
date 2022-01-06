package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

const KafkaTLSProtocolFlag = "TLS"

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
	VaultToken                 string        `envconfig:"VAULT_TOKEN"`
	VaultAddress               string        `envconfig:"VAULT_ADDR"`
	VaultPath                  string        `envconfig:"VAULT_PATH"`
	EncryptionDisabled         bool          `envconfig:"ENCRYPTION_DISABLED"`
	ComponentTestUseLogFile    bool          `envconfig:"COMPONENT_TEST_USE_LOG_FILE"`
	ServiceAuthToken           string        `envconfig:"SERVICE_AUTH_TOKEN"`
	DownloadServiceURL         string        `envconfig:"DOWNLOAD_SERVICE_URL"`
	StopConsumingOnUnhealthy   bool          `envconfig:"STOP_CONSUMING_ON_UNHEALTHY"`
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
	ConsumerMinBrokersHealthy     int      `envconfig:"KAFKA_CONSUMER_MIN_BROKERS_HEALTHY"`
	ProducerMinBrokersHealthy     int      `envconfig:"KAFKA_PRODUCER_MIN_BROKERS_HEALTHY"`
	CantabularCSVCreatedTopic     string   `envconfig:"KAFKA_TOPIC_CANTABULAR_CSV_CREATED"`
	CantabularCSVWCreatedTopic    string   `envconfig:"KAFKA_TOPIC_CANTABULAR_CSVW_CREATED"`
	CantabularMetadataExportGroup string   `envconfig:"KAFKA_GROUP_CANTABULAR_METADATA_EXPORT"`
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
		AWSRegion:                  "eu-west-1",
		VaultPath:                  "secret/shared/psk",
		VaultAddress:               "http://localhost:8200",
		VaultToken:                 "",
		PublicBucket:               "dp-cantabular-metadata-exporter",
		PrivateBucket:              "dp-cantabular-metadata-exporter",
		S3BucketURL:                "",
		DatasetAPIURL:              "http://localhost:22000",
		ComponentTestUseLogFile:    false,
		DownloadServiceURL:         "http://localhost:23600",
		StopConsumingOnUnhealthy:   true,
		Kafka: KafkaConfig{
			Addr:                          []string{"localhost:9092", "localhost:9093", "localhost:9094"},
			Version:                       "1.0.2",
			OffsetOldest:                  true,
			NumWorkers:                    1,
			MaxBytes:                      2000000,
			SecProtocol:                   "",
			SecCACerts:                    "",
			SecClientKey:                  "",
			SecClientCert:                 "",
			SecSkipVerify:                 false,
			ConsumerMinBrokersHealthy:     1,
			ProducerMinBrokersHealthy:     2,
			CantabularMetadataExportGroup: "cantabular-metadata-export",
			CantabularCSVCreatedTopic:     "cantabular-csv-created",
			CantabularCSVWCreatedTopic:    "cantabular-csvw-created",
		},
	}

	return cfg, envconfig.Process("", cfg)
}
