package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-cantabular-metadata-exporter/config"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/ONSdigital/dp-cantabular-metadata-exporter/filemanager"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/generator"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	kafka "github.com/ONSdigital/dp-kafka/v4"
	dphttp "github.com/ONSdigital/dp-net/v2/http"
	vault "github.com/ONSdigital/dp-vault"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

const (
	VaultRetries = 3
)

// GetHTTPServer creates an http server
var GetHTTPServer = func(bindAddr string, router http.Handler) HTTPServer {
	otelHandler := otelhttp.NewHandler(router, "/")
	s := dphttp.NewServer(bindAddr, otelHandler)
	s.HandleOSSignals = false
	return s
}

// GetHealthCheck creates a healthcheck with versionInfo and sets teh HealthCheck flag to true
var GetHealthCheck = func(cfg *config.Config, buildT, commit, ver string) (HealthChecker, error) {
	versionInfo, err := healthcheck.NewVersionInfo(buildT, commit, ver)
	if err != nil {
		return nil, fmt.Errorf("failed to get version info: %w", err)
	}

	hc := healthcheck.New(versionInfo, cfg.HealthCheckCriticalTimeout, cfg.HealthCheckInterval)
	return &hc, nil
}

var GetGenerator = func() Generator {
	return generator.New()
}
var GetVaultClient = func(cfg *config.Config) (VaultClient, error) {
	return vault.CreateClient(cfg.VaultToken, cfg.VaultAddress, VaultRetries)
}

// GetKafkaConsumer creates a Kafka consumer
var GetKafkaConsumer = func(ctx context.Context, cfg *config.Config) (kafka.IConsumerGroup, error) {
	kafkaOffset := kafka.OffsetNewest
	if cfg.Kafka.OffsetOldest {
		kafkaOffset = kafka.OffsetOldest
	}

	cgConfig := &kafka.ConsumerGroupConfig{
		BrokerAddrs:       cfg.Kafka.Addr,
		Topic:             cfg.Kafka.CantabularCSVCreatedTopic,
		GroupName:         cfg.Kafka.CantabularMetadataExportGroup,
		KafkaVersion:      &cfg.Kafka.Version,
		Offset:            &kafkaOffset,
		NumWorkers:        &cfg.Kafka.NumWorkers,
		MinBrokersHealthy: &cfg.Kafka.ConsumerMinBrokersHealthy,
	}

	if cfg.Kafka.SecProtocol == config.KafkaTLSProtocolFlag {
		cgConfig.SecurityConfig = kafka.GetSecurityConfig(
			cfg.Kafka.SecCACerts,
			cfg.Kafka.SecClientCert,
			cfg.Kafka.SecClientKey,
			cfg.Kafka.SecSkipVerify,
		)
	}

	return kafka.NewConsumerGroup(ctx, cgConfig)
}

// GetKafkaProducer creates a Kafka producer
var GetKafkaProducer = func(ctx context.Context, cfg *config.Config) (kafka.IProducer, error) {
	pConfig := &kafka.ProducerConfig{
		BrokerAddrs:       cfg.Kafka.Addr,
		Topic:             cfg.Kafka.CantabularCSVWCreatedTopic,
		KafkaVersion:      &cfg.Kafka.Version,
		MaxMessageBytes:   &cfg.Kafka.MaxBytes,
		MinBrokersHealthy: &cfg.Kafka.ProducerMinBrokersHealthy,
	}

	if cfg.Kafka.SecProtocol == config.KafkaTLSProtocolFlag {
		pConfig.SecurityConfig = kafka.GetSecurityConfig(
			cfg.Kafka.SecCACerts,
			cfg.Kafka.SecClientCert,
			cfg.Kafka.SecClientKey,
			cfg.Kafka.SecSkipVerify,
		)
	}

	return kafka.NewProducer(ctx, pConfig)
}

// GetDatasetAPIClient gets and initialises the DatasetAPI Client
var GetDatasetAPIClient = func(cfg *config.Config) DatasetAPIClient {
	return dataset.NewAPIClient(cfg.DatasetAPIURL)
}

// GetFilterAPIClient gets and initialises the FilterAPI Client
var GetFilterAPIClient = func(cfg *config.Config) FilterAPIClient {
	return filter.New(cfg.FilterAPIURL)
}

// GetPopulationTypesAPIClient gets and initialises the PopulationTypesAPI Client
var GetPopulationTypesAPIClient = func(cfg *config.Config) (PopulationTypesAPIClient, error) {
	return population.NewClient(cfg.PopulationTypesAPIURL)
}

// GetCantabularClient gets and initialises the Cantabular Client
var GetCantabularClient = func(cfg *config.Config) CantabularClient {
	return cantabular.NewClient(
		cantabular.Config{
			Host:           cfg.CantabularURL,
			ExtApiHost:     cfg.CantabularExtURL,
			GraphQLTimeout: cfg.DefaultRequestTimeout,
		},
		dphttp.NewClient(),
		nil,
	)
}

// GetFileManager instantiates teh service FileManager
var GetFileManager = func(cfg *config.Config, vault VaultClient, generator Generator) (FileManager, error) {
	awscfg := &aws.Config{
		Region: aws.String(cfg.AWSRegion),
	}

	if cfg.LocalObjectStore != "" {
		awscfg = &aws.Config{
			Credentials:      credentials.NewStaticCredentials(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
			Endpoint:         aws.String(cfg.LocalObjectStore),
			Region:           aws.String(cfg.AWSRegion),
			DisableSSL:       aws.Bool(true),
			S3ForcePathStyle: aws.Bool(true),
		}
	}

	sess, err := session.NewSession(awscfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create aws session: %w", err)
	}

	return filemanager.New(
		filemanager.Config{
			VaultKey:      "key",
			PublicBucket:  cfg.PublicBucket,
			PrivateBucket: cfg.PrivateBucket,
			PublicURL:     cfg.S3BucketURL,
		},
		sess,
		vault,
		generator,
	), nil
}
