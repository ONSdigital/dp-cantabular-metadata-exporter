package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/ONSdigital/dp-cantabular-metadata-exporter/config"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/service"
	"github.com/ONSdigital/log.go/v2/log"
)

const serviceName = "dp-cantabular-metadata-exporter"

var (
	// BuildTime represents the time in which the service was built
	BuildTime string
	// GitCommit represents the commit (SHA-1) hash of the service that is running
	GitCommit string
	// Version represents the version of the service that is running
	Version string

// TODO: remove below explainer before commiting
/* NOTE: replace the above with the below to run code with for example vscode debugger.
BuildTime string = "1601119818"
GitCommit string = "6584b786caac36b6214ffe04bf62f058d4021538"
Version   string = "v0.1.0"
*/
)

func main() {
	log.Namespace = serviceName
	ctx := context.Background()

	if err := run(ctx); err != nil {
		log.Fatal(context.TODO(), "fatal runtime error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)

	// Run the service, providing an error channel for fatal errors
	svcErrors := make(chan error, 1)

	log.Info(ctx, "dp-cantabular-metadata-exporter version", log.Data{"version": Version})

	// Read config
	cfg, err := config.Get()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	// Create, initialise and start service
	svc := service.New()

	if err := svc.Init(ctx, cfg, BuildTime, GitCommit, Version); err != nil {
		return fmt.Errorf("failed to initialise service: %w", err)
	}

	if err := svc.Start(ctx, svcErrors); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}

	// blocks until an os interrupt or a fatal error occurs
	select {
	case err := <-svcErrors:
		// TODO: call svc.Close(ctx) (or something specific)
		//  if there are any service connections like Kafka that you need to shut down
		return fmt.Errorf("service error: %w", err)
	case sig := <-signals:
		log.Info(ctx, "os signal received", log.Data{"signal": sig})
	}
	return svc.Close(ctx)
}
