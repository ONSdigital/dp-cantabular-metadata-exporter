package service

import (
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-cantabular-metadata-exporter/config"

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
var GetKafkaConsumer = func(ctx context.Context, cfg *config.Config) (dpkafka.IConsumerGroup, error) {
	cgChannels := dpkafka.CreateConsumerGroupChannels(1)

	kafkaOffset := dpkafka.OffsetNewest
	if cfg.KafkaOffsetOldest {
		kafkaOffset = dpkafka.OffsetOldest
	}

	return dpkafka.NewConsumerGroup(
		ctx,
		cfg.KafkaAddr,
		cfg.InstanceCompleteTopic,
		cfg.InstanceCompleteGroup,
		cgChannels,
		&dpkafka.ConsumerGroupConfig{
			KafkaVersion: &cfg.KafkaVersion,
			Offset:       &kafkaOffset,
		},
	)
}

// GetKafkaProducer creates a Kafka producer
var GetKafkaProducer = func(ctx context.Context, cfg *config.Config) (dpkafka.IProducer, error) {
	pChannels := dpkafka.CreateProducerChannels()
	return dpkafka.NewProducer(
		ctx,
		cfg.KafkaAddr,
		cfg.CommonOutputCreatedTopic,
		pChannels,
		&kafka.ProducerConfig{},
	)
}
