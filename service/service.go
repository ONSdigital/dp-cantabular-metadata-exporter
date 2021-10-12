package service

import (
	"context"
	"fmt"

	"github.com/ONSdigital/dp-cantabular-metadata-exporter/api"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/config"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
)

// Service contains all the configs, server and clients to run the dp-topic-api API
type Service struct {
	config      *config.Config
	server      HTTPServer
	router      *mux.Router
	api         *api.API
	healthCheck HealthChecker
}

// New returns a new Service
func New() *Service{
	return &Service{}
}

// Init initialises the service
func (svc *Service) Init(ctx context.Context, cfg *config.Config, buildT, commit, ver string) error {
	log.Event(ctx, "initialising service with config", log.Data{"config": cfg}, log.INFO)

	svc.config = cfg

	var err error

	svc.router = mux.NewRouter()
	svc.api    = api.Setup(ctx, svc.router)
	svc.server = GetHTTPServer(cfg.BindAddr, svc.router)
	// TODO: Add other(s) to serviceList here

	if svc.healthCheck, err = GetHealthCheck(cfg, buildT, commit, ver); err != nil {
		return fmt.Errorf("could not get healtcheck: %w", err)
	}

	if err := svc.registerCheckers(); err != nil {
		return fmt.Errorf("unable to register checkers: %w", err)
	}

	svc.router.StrictSlash(true).Path("/health").HandlerFunc(svc.healthCheck.Handler)

	return nil
}

// Start starts the service
func (svc *Service) Start(ctx context.Context, svcErrors chan error){
	log.Event(ctx, "starting service", log.Data{}, log.INFO)

	svc.healthCheck.Start(ctx)

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
	log.Event(ctx, "commencing graceful shutdown", log.Data{"graceful_shutdown_timeout": timeout}, log.INFO)
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

	log.Event(ctx, "graceful shutdown was successful", log.INFO)
	return nil
}

func (svc *Service) registerCheckers() error {
	// TODO: add other health checks here, as per dp-upload-service
	return nil
}
