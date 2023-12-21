package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

const KafkaTLSProtocolFlag = "TLS"

// Config represents service configuration for dp-cantabular-metadata-exporter
type Config struct {
	AWSRegion                  string        `envconfig:"AWS_REGION"`
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	CantabularExtURL           string        `envconfig:"CANTABULAR_API_EXT_URL"`
	CantabularURL              string        `envconfig:"CANTABULAR_URL"`
	ComponentTestUseLogFile    bool          `envconfig:"COMPONENT_TEST_USE_LOG_FILE"`
	DatasetAPIURL              string        `envconfig:"DATASET_API_URL"`
	DefaultRequestTimeout      time.Duration `envconfig:"DEFAULT_REQUEST_TIMEOUT"`
	DownloadServiceURL         string        `envconfig:"DOWNLOAD_SERVICE_URL"`
	EncryptionDisabled         bool          `envconfig:"ENCRYPTION_DISABLED"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	ExternalPrefixURL          string        `envconfig:"EXTERNAL_PREFIX_URL"`
	FilterAPIURL               string        `envconfig:"FILTER_API_URL"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	Kafka                      KafkaConfig
	LocalObjectStore           string        `envconfig:"LOCAL_OBJECT_STORE"`
	MinioAccessKey             string        `envconfig:"MINIO_ACCESS_KEY"                           json:"-"`
	MinioSecretKey             string        `envconfig:"MINIO_SECRET_KEY"                           json:"-"`
	OTExporterOTLPEndpoint     string        `envconfig:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	OTServiceName              string        `envconfig:"OTEL_SERVICE_NAME"`
	OTBatchTimeout             time.Duration `envconfig:"OTEL_BATCH_TIMEOUT"`
	PopulationTypesAPIURL      string        `envconfig:"POPULATION_TYPES_API_URL"`
	PublicBucket               string        `envconfig:"PUBLIC_BUCKET"`
	PrivateBucket              string        `envconfig:"PRIVATE_BUCKET"`
	ServiceAuthToken           string        `envconfig:"SERVICE_AUTH_TOKEN"                         json:"-"`
	StopConsumingOnUnhealthy   bool          `envconfig:"STOP_CONSUMING_ON_UNHEALTHY"`
	S3BucketURL                string        `envconfig:"PUBLIC_URL"`
	S3PublicURL                string        `envconfig:"S3_PUBLIC_URL"`
	VaultAddress               string        `envconfig:"VAULT_ADDR"`
	VaultPath                  string        `envconfig:"VAULT_PATH"`
	VaultToken                 string        `envconfig:"VAULT_TOKEN"                                json:"-"`
}

// KafkaConfig contains the config required to connect to Kafka
type KafkaConfig struct {
	Addr                          []string `envconfig:"KAFKA_ADDR"                            json:"-"`
	CantabularCSVCreatedTopic     string   `envconfig:"KAFKA_TOPIC_CANTABULAR_CSV_CREATED"`
	CantabularCSVWCreatedTopic    string   `envconfig:"KAFKA_TOPIC_CANTABULAR_CSVW_CREATED"`
	CantabularMetadataExportGroup string   `envconfig:"KAFKA_GROUP_CANTABULAR_METADATA_EXPORT"`
	ConsumerMinBrokersHealthy     int      `envconfig:"KAFKA_CONSUMER_MIN_BROKERS_HEALTHY"`
	MaxBytes                      int      `envconfig:"KAFKA_MAX_BYTES"`
	NumWorkers                    int      `envconfig:"KAFKA_NUM_WORKERS"`
	OffsetOldest                  bool     `envconfig:"KAFKA_OFFSET_OLDEST"`
	ProducerMinBrokersHealthy     int      `envconfig:"KAFKA_PRODUCER_MIN_BROKERS_HEALTHY"`
	SecCACerts                    string   `envconfig:"KAFKA_SEC_CA_CERTS"`
	SecClientCert                 string   `envconfig:"KAFKA_SEC_CLIENT_CERT"`
	SecClientKey                  string   `envconfig:"KAFKA_SEC_CLIENT_KEY"                  json:"-"`
	SecProtocol                   string   `envconfig:"KAFKA_SEC_PROTO"`
	SecSkipVerify                 bool     `envconfig:"KAFKA_SEC_SKIP_VERIFY"`
	Version                       string   `envconfig:"KAFKA_VERSION"`
}

var cfg *Config

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		AWSRegion:                  "eu-west-2",
		BindAddr:                   ":26700",
		CantabularExtURL:           "http://localhost:8492",
		CantabularURL:              "http://localhost:8491",
		ComponentTestUseLogFile:    false,
		DatasetAPIURL:              "http://localhost:22000",
		DefaultRequestTimeout:      10 * time.Second,
		DownloadServiceURL:         "http://localhost:23600",
		EncryptionDisabled:         false,
		ExternalPrefixURL:          "http://localhost:22000",
		FilterAPIURL:               "http://localhost:22100",
		GracefulShutdownTimeout:    5 * time.Second,
		HealthCheckCriticalTimeout: 90 * time.Second,
		HealthCheckInterval:        30 * time.Second,
		Kafka: KafkaConfig{
			Addr:                          []string{"localhost:9092", "localhost:9093", "localhost:9094"},
			CantabularCSVCreatedTopic:     "cantabular-csv-created",
			CantabularCSVWCreatedTopic:    "cantabular-csvw-created",
			CantabularMetadataExportGroup: "cantabular-metadata-export",
			ConsumerMinBrokersHealthy:     1,
			MaxBytes:                      2000000,
			NumWorkers:                    1,
			OffsetOldest:                  true,
			ProducerMinBrokersHealthy:     2,
			SecCACerts:                    "",
			SecClientCert:                 "",
			SecClientKey:                  "",
			SecProtocol:                   "",
			SecSkipVerify:                 false,
			Version:                       "1.0.2",
		},
		LocalObjectStore:         "",
		MinioAccessKey:           "",
		MinioSecretKey:           "",
		OTExporterOTLPEndpoint:   "localhost:4317",
		OTServiceName:            "dp-cantabular-metadata-exporter",
		OTBatchTimeout:           5 * time.Second,
		PopulationTypesAPIURL:    "http://localhost:27300",
		PrivateBucket:            "dp-cantabular-metadata-exporter",
		PublicBucket:             "dp-cantabular-metadata-exporter",
		ServiceAuthToken:         "",
		StopConsumingOnUnhealthy: true,
		S3BucketURL:              "",
		S3PublicURL:              "http://public-bucket",
		VaultAddress:             "http://localhost:8200",
		VaultPath:                "secret/shared/psk",
		VaultToken:               "",
	}

	return cfg, envconfig.Process("", cfg)
}
