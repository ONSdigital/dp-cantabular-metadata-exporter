package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/ONSdigital/dp-cantabular-metadata-exporter/config"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/handler"
	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/log.go/v2/log"

	"github.com/go-chi/chi/v5"
)

// Service contains all the configs, server and clients to run the dp-topic-api API
type Service struct {
	Config           *config.Config
	Server           HTTPServer
	router           chi.Router
	consumer         kafka.IConsumerGroup
	producer         kafka.IProducer
	datasetAPIClient DatasetAPIClient
	HealthCheck      HealthChecker
	vaultClient      VaultClient
	generator        Generator
	fileManager      FileManager
}

// New returns a new Service
func New() *Service {
	return &Service{}
}

// Producer is a getter for the kafka producer for use outside package
func (svc *Service) Producer() kafka.IProducer {
	return svc.producer
}

// Consumer is a getter for the kafka consumer for use outside package
func (svc *Service) Consumer() kafka.IConsumerGroup {
	return svc.consumer
}

// Init initialises the service
func (svc *Service) Init(ctx context.Context, cfg *config.Config, buildT, commit, ver string) error {
	log.Info(ctx, "initialising service with config", log.Data{"config": cfg})

	svc.Config = cfg

	var err error

	if svc.consumer, err = GetKafkaConsumer(ctx, cfg); err != nil {
		return fmt.Errorf("failed to create kafka consumer: %w", err)
	}
	if svc.producer, err = GetKafkaProducer(ctx, cfg); err != nil {
		return fmt.Errorf("failed to create kafka producer: %w", err)
	}

	if svc.vaultClient, err = GetVaultClient(cfg); err != nil {
		return fmt.Errorf("failed to initialise vault client: %w", err)
	}

	svc.generator = GetGenerator()

	if svc.fileManager, err = GetFileManager(cfg, svc.vaultClient, svc.generator); err != nil {
		return fmt.Errorf("failed to initialise file manager: %w", err)
	}

	svc.datasetAPIClient = GetDatasetAPIClient(cfg)

	if svc.HealthCheck, err = GetHealthCheck(cfg, buildT, commit, ver); err != nil {
		return fmt.Errorf("could not get healtcheck: %w", err)
	}

	h := handler.NewCantabularMetadataExport(
		*svc.Config,
		svc.datasetAPIClient,
		svc.fileManager,
		svc.producer,
	)
	if err := svc.consumer.RegisterHandler(ctx, h.Handle); err != nil {
		return fmt.Errorf("could not register kafka handler: %w", err)
	}

	if err := svc.registerCheckers(); err != nil {
		return fmt.Errorf("unable to register checkers: %w", err)
	}

	svc.BuildRoutes()
	svc.Server = GetHTTPServer(cfg.BindAddr, svc.router)

	return nil
}

// Start starts the service
func (svc *Service) Start(ctx context.Context, svcErrors chan error) error {
	log.Info(ctx, "starting service", log.Data{})

	svc.HealthCheck.Start(ctx)

	// Kafka error logging go-routine
	svc.consumer.LogErrors(ctx)
	svc.producer.LogErrors(ctx)

	// If start/stop on health updates is disabled, start consuming as soon as possible
	if !svc.Config.StopConsumingOnUnhealthy {
		if err := svc.consumer.Start(); err != nil {
			return fmt.Errorf("consumer failed to start: %w", err)
		}
	}

	// Run the http server in a new go-routine
	go func() {
		if err := svc.Server.ListenAndServe(); err != nil {
			svcErrors <- fmt.Errorf("failed to start main http server: %w", err)
		}
	}()

	return nil
}

// Close gracefully shuts the service down in the required order, with timeout
func (svc *Service) Close(ctx context.Context) error {
	timeout := svc.Config.GracefulShutdownTimeout
	log.Info(ctx, "commencing graceful shutdown", log.Data{"graceful_shutdown_timeout": timeout})
	ctx, cancel := context.WithTimeout(ctx, timeout)

	// track shutown gracefully closes up
	var hasShutDownErr bool

	go func() {
		defer cancel()

		// stop healthcheck, as it depends on everything else
		if svc.HealthCheck != nil {
			svc.HealthCheck.Stop()
			log.Info(ctx, "stopped health checker")
		}

		// If kafka consumer exists, stop and close it.
		if svc.consumer != nil {
			if err := svc.consumer.StopAndWait(); err != nil {
				log.Error(ctx, "failed to stop kafka consumer", err)
				hasShutDownErr = true
			} else {
				log.Info(ctx, "stopped kafka consumer")
			}
		}

		// If kafka producer exists, close it.
		if svc.producer != nil {
			if err := svc.producer.Close(ctx); err != nil {
				log.Error(ctx, "error closing kafka producer", err)
				hasShutDownErr = true
			} else {
				log.Info(ctx, "closed kafka producer")
			}
		}

		// stop any incoming requests before closing any outbound connections
		if svc.Server != nil {
			if err := svc.Server.Shutdown(ctx); err != nil {
				log.Error(ctx, "failed to shutdown http server", err)
				hasShutDownErr = true
			} else {
				log.Info(ctx, "stopped http server")
			}
		}

		// If kafka consumer exists, close it.
		if svc.consumer != nil {
			if err := svc.consumer.Close(ctx); err != nil {
				log.Error(ctx, "error closing kafka consumer", err)
				hasShutDownErr = true
			} else {
				log.Info(ctx, "closed kafka consumer")
			}
		}
	}()

	// wait for shutdown success (via cancel) or failure (timeout)
	<-ctx.Done()

	// timeout expired
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("shutdown timed out: %w", ctx.Err())
	}

	// other error
	if hasShutDownErr {
		return errors.New("failed to shutdown gracefully")
	}

	log.Info(ctx, "graceful shutdown was successful")
	return nil
}

func (svc *Service) registerCheckers() error {
	if _, err := svc.HealthCheck.AddAndGetCheck("Kafka consumer", svc.consumer.Checker); err != nil {
		return fmt.Errorf("error adding Kafka consumer health check: %w", err)
	}

	if _, err := svc.HealthCheck.AddAndGetCheck("Kafka producer", svc.producer.Checker); err != nil {
		return fmt.Errorf("error adding Kafka producer health check: %w", err)
	}

	if _, err := svc.HealthCheck.AddAndGetCheck("Vault", svc.vaultClient.Checker); err != nil {
		return fmt.Errorf("error adding vault health check: %w", err)
	}

	if _, err := svc.HealthCheck.AddAndGetCheck("S3 private uploader", svc.fileManager.PublicUploader().Checker); err != nil {
		return fmt.Errorf("error adding s3 private uploader health check: %w", err)
	}

	if _, err := svc.HealthCheck.AddAndGetCheck("S3 public uploader", svc.fileManager.PublicUploader().Checker); err != nil {
		return fmt.Errorf("error adding s3 public uploader health check: %w", err)
	}

	if _, err := svc.HealthCheck.AddAndGetCheck("Dataset API client", svc.datasetAPIClient.Checker); err != nil {
		return fmt.Errorf("error adding dataset API health check: %w", err)
	}

	if svc.Config.StopConsumingOnUnhealthy {
		svc.HealthCheck.SubscribeAll(svc.consumer)
	}

	return nil
}
