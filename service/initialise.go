package service

import (
	"fmt"
	"net/http"
	"context"

	"github.com/ONSdigital/dp-cantabular-metadata-exporter/config"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/event"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/generator"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/filemanager"
	
	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	kafka "github.com/ONSdigital/dp-kafka/v2"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dphttp "github.com/ONSdigital/dp-net/http"
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
	s := dphttp.NewServer(bindAddr, router)
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

var GetGenerator = func() Generator{
	return generator.New()
}
var GetVaultClient = func(cfg *config.Config) (VaultClient, error) {
	return vault.CreateClient(cfg.VaultToken, cfg.VaultAddress, VaultRetries)
}

// GetKafkaConsumer creates a Kafka consumer
var GetKafkaConsumer = func(ctx context.Context, cfg *config.Config) (kafka.IConsumerGroup, error) {
	cgChannels := kafka.CreateConsumerGroupChannels(cfg.Kafka.NumWorkers)

	kafkaOffset := kafka.OffsetNewest
	if cfg.Kafka.OffsetOldest {
		kafkaOffset = kafka.OffsetOldest
	}

	return kafka.NewConsumerGroup(
		ctx,
		cfg.Kafka.Addr,
		cfg.Kafka.CantabularMetadataExportTopic,
		cfg.Kafka.CantabularMetadataExportGroup,
		cgChannels,
		&kafka.ConsumerGroupConfig{
			KafkaVersion: &cfg.Kafka.Version,
			Offset:       &kafkaOffset,
		},
	)
}

// GetKafkaProducer creates a Kafka producer
var GetKafkaProducer = func(ctx context.Context, cfg *config.Config) (kafka.IProducer, error) {
	pChannels := kafka.CreateProducerChannels()
	return kafka.NewProducer(
		ctx,
		cfg.Kafka.Addr,
		cfg.Kafka.CantabularMetadataExportTopic,
		pChannels,
		&kafka.ProducerConfig{},
	)
}

// GetProcessor gets and initialises the event Processor
var GetProcessor = func(cfg *config.Config) Processor {
	return event.NewProcessor(*cfg)
}

// GetDatasetAPIClient gets and initialises the DatasetAPI Client
var GetDatasetAPIClient = func(cfg *config.Config) DatasetAPIClient {
 	return dataset.NewAPIClient(cfg.DatasetAPIURL)
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
