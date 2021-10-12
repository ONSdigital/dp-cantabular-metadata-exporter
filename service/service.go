package service

import (
	"context"

	"github.com/ONSdigital/dp-cantabular-metadata-exporter/api"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/config"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// Service contains all the configs, server and clients to run the dp-topic-api API
type Service struct {
	Config      *config.Config
	Server      HTTPServer
	Router      *mux.Router
	Api         *api.API
	ServiceList *ExternalServiceList
	HealthCheck HealthChecker
}

// Run the service
func Run(ctx context.Context, cfg *config.Config, serviceList *ExternalServiceList, buildTime, gitCommit, version string, svcErrors chan error) (*Service, error) {
	log.Event(ctx, "running service", log.INFO)
	log.Event(ctx, "using service configuration", log.Data{"config": cfg}, log.INFO)

	// Get HTTP Server and ... // TODO: Add any middleware that your service requires
	r := mux.NewRouter()

	s := serviceList.GetHTTPServer(cfg.BindAddr, r)
	// TODO: Add other(s) to serviceList here

	// Setup the API
	a := api.Setup(ctx, r)

	hc, err := serviceList.GetHealthCheck(cfg, buildTime, gitCommit, version)
	if err != nil {
		return nil, fmt.Errorf("could not get healtcheck: %w", err)
	}

	if err := registerCheckers(ctx, hc); err != nil {
		return nil, fmt.Errorf("unable to register checkers: %w", err)
	}

	r.StrictSlash(true).Path("/health").HandlerFunc(hc.Handler)
	hc.Start(ctx)

	// Run the http server in a new go-routine
	go func() {
		if err := s.ListenAndServe(); err != nil {
			svcErrors <- fmt.Errorf("failed to start main http server: %w", err)
		}
	}()

	return &Service{
		Config:      cfg,
		Router:      r,
		Api:         a,
		HealthCheck: hc,
		ServiceList: serviceList,
		Server:      s,
	}, nil
}

// Close gracefully shuts the service down in the required order, with timeout
func (svc *Service) Close(ctx context.Context) error {
	timeout := svc.Config.GracefulShutdownTimeout
	log.Event(ctx, "commencing graceful shutdown", log.Data{"graceful_shutdown_timeout": timeout}, log.INFO)
	ctx, cancel := context.WithTimeout(ctx, timeout)

	// track shutown gracefully closes up
	var hasShutdownError bool

	go func() {
		defer cancel()

		// stop healthcheck, as it depends on everything else
		if svc.ServiceList.HealthCheck {
			svc.HealthCheck.Stop()
		}

		// stop any incoming requests before closing any outbound connections
		if err := svc.Server.Shutdown(ctx); err != nil {
			hasShutdownError = true
		}

		// TODO: Close other dependencies, in the expected order
	}()

	// wait for shutdown success (via cancel) or failure (timeout)
	<-ctx.Done()

	// timeout expired
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("shutdown timed out: %w", ctx.Err())
	}

	// other error
	if hasShutdownError {
		return fmt.Errorf("failed to shutdown gracefully: %w", err)
	}

	log.Event(ctx, "graceful shutdown was successful", log.INFO)
	return nil
}

func registerCheckers(ctx context.Context, hc HealthChecker) error {
	// TODO: add other health checks here, as per dp-upload-service
	return nil
}
