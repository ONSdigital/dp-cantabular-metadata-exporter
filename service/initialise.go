package service

import (
	"fmt"
	"net/http"
	"context"

	"github.com/ONSdigital/dp-cantabular-metadata-exporter/config"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/event"

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
