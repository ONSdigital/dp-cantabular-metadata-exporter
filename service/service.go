package service

import (
	"context"
	"fmt"

	"github.com/ONSdigital/dp-cantabular-metadata-exporter/config"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/handler"
	"github.com/ONSdigital/log.go/v2/log"
	kafka "github.com/ONSdigital/dp-kafka/v2"

	"github.com/go-chi/chi/v5"
)

// Service contains all the configs, server and clients to run the dp-topic-api API
type Service struct {
	config      *config.Config
	server      HTTPServer
	router      chi.Router
	consumer    kafka.IConsumerGroup
	producer    kafka.IProducer
	processor   Processor
	healthCheck HealthChecker
}

// New returns a new Service
func New() *Service {
	return &Service{}
}

// Init initialises the service
func (svc *Service) Init(ctx context.Context, cfg *config.Config, buildT, commit, ver string) error {
	log.Info(ctx, "initialising service with config", log.Data{"config": cfg})

	svc.config = cfg

	var err error

	if svc.consumer, err = GetKafkaConsumer(ctx, cfg); err != nil {
		return fmt.Errorf("failed to create kafka consumer: %w", err)
	}
	if svc.producer, err = GetKafkaProducer(ctx, cfg); err != nil {
		return fmt.Errorf("failed to create kafka producer: %w", err)
	}

	if svc.healthCheck, err = GetHealthCheck(cfg, buildT, commit, ver); err != nil {
		return fmt.Errorf("could not get healtcheck: %w", err)
	}

	if err := svc.registerCheckers(); err != nil {
		return fmt.Errorf("unable to register checkers: %w", err)
	}

	svc.BuildRoutes()
	svc.server = GetHTTPServer(cfg.BindAddr, svc.router)

	return nil
}

// Start starts the service
func (svc *Service) Start(ctx context.Context, svcErrors chan error) {
	log.Info(ctx, "starting service", log.Data{})

	svc.healthCheck.Start(ctx)

	// Kafka error logging go-routine
	svc.consumer.Channels().LogErrors(ctx, "kafka consumer")

	// Event Handler for Kafka Consumer
	svc.processor.Consume(
		ctx,
		svc.consumer,
		handler.NewCantabularMetadataExport(
			*svc.config,
		),
	)

	// Run the http server in a new go-routine
	go func() {
		if err := svc.server.ListenAndServe(); err != nil {
			svcErrors <- fmt.Errorf("failed to start main http server: %w", err)
		}
	}()
}

// Close gracefully shuts the service down in the required order, with timeout
func (svc *Service) Close(ctx context.Context) error {
	timeout := svc.config.GracefulShutdownTimeout
	log.Info(ctx, "commencing graceful shutdown", log.Data{"graceful_shutdown_timeout": timeout})
	ctx, cancel := context.WithTimeout(ctx, timeout)

	// track shutown gracefully closes up
	var shutDownErr error

	go func() {
		defer cancel()

		// stop healthcheck, as it depends on everything else
		if svc.healthCheck != nil {
			svc.healthCheck.Stop()
		}

		// stop any incoming requests before closing any outbound connections
		shutDownErr = svc.server.Shutdown(ctx)

		// TODO: Close other dependencies, in the expected order
	}()

	// wait for shutdown success (via cancel) or failure (timeout)
	<-ctx.Done()

	// timeout expired
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("shutdown timed out: %w", ctx.Err())
	}

	// other error
	if shutDownErr != nil {
		return fmt.Errorf("failed to shutdown gracefully: %w", shutDownErr)
	}

	log.Info(ctx, "graceful shutdown was successful")
	return nil
}

func (svc *Service) registerCheckers() error {
	// TODO: add other health checks here, as per dp-upload-service
	return nil
}
