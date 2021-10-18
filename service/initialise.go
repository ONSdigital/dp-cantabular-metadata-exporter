package service

import (
	"fmt"
	"net/http"
	"context"

	"github.com/ONSdigital/dp-cantabular-metadata-exporter/config"
	kafka "github.com/ONSdigital/dp-kafka/v2"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dphttp "github.com/ONSdigital/dp-net/http"
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

// GetKafkaConsumer creates a Kafka consumer
var GetKafkaConsumer = func(ctx context.Context, cfg *config.Config) (kafka.IConsumerGroup, error) {
	cgChannels := kafka.CreateConsumerGroupChannels(1)

	kafkaOffset := kafka.OffsetNewest
	if cfg.KafkaOffsetOldest {
		kafkaOffset = kafka.OffsetOldest
	}

	return kafka.NewConsumerGroup(
		ctx,
		cfg.KafkaAddr,
		cfg.CantabularMetadataExportTopic,
		cfg.CantabularMetadataExportGroup,
		cgChannels,
		&kafka.ConsumerGroupConfig{
			KafkaVersion: &cfg.KafkaVersion,
			Offset:       &kafkaOffset,
		},
	)
}

// GetKafkaProducer creates a Kafka producer
var GetKafkaProducer = func(ctx context.Context, cfg *config.Config) (kafka.IProducer, error) {
	pChannels := kafka.CreateProducerChannels()
	return kafka.NewProducer(
		ctx,
		cfg.KafkaAddr,
		cfg.CantabularMetadataExportTopic,
		pChannels,
		&kafka.ProducerConfig{},
	)
}
