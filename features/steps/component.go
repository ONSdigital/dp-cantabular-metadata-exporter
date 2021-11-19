package steps

import (
	"context"
	"fmt"
	"testing"
	"sync"
	"time"
	"os"
	"os/signal"
	"syscall"

	"github.com/ONSdigital/dp-cantabular-metadata-exporter/config"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/service"

	kafka "github.com/ONSdigital/dp-kafka/v2"
	cmptest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/log.go/v2/log"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/maxcnunes/httpfake"
)

var (
	BuildTime string = "1625046891"
	GitCommit string = "7434fe334d9f51b7239f978094ea29d10ac33b16"
	Version   string = ""
)

type Component struct {
	cmptest.ErrorFeature
	apiFeature       *cmptest.APIFeature
	DatasetAPI       *httpfake.HTTPFake
	S3Downloader     *s3manager.Downloader
	producer         kafka.IProducer
	errorChan        chan error
	svc              *service.Service
	cfg              *config.Config
	wg               *sync.WaitGroup
	signals          chan os.Signal
	waitEventTimeout time.Duration
}

func NewComponent(t *testing.T) *Component {
	return &Component{
		errorChan:        make(chan error),
		DatasetAPI:       httpfake.New(httpfake.WithTesting(t)),
		wg:               &sync.WaitGroup{},
		waitEventTimeout: time.Second * 5,
	}
}

// initService initialises the server, the mocks and waits for the dependencies to be ready
func (c *Component) initService(ctx context.Context) error {
	// register interrupt signals
	c.signals = make(chan os.Signal, 1)
	signal.Notify(c.signals, os.Interrupt, syscall.SIGTERM)

	// Read config
	cfg, err := config.Get()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	log.Info(ctx, "config read", log.Data{"cfg": cfg})

	cfg.DatasetAPIURL = c.DatasetAPI.ResolveURL("")

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Endpoint:         aws.String(cfg.LocalObjectStore),
		Region:           aws.String(cfg.AWSRegion),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}

	s := session.New(s3Config)
	c.S3Downloader = s3manager.NewDownloader(s)

	// producer for triggering test events
	if c.producer, err = kafka.NewProducer(
		ctx,
		cfg.Kafka.Addr,
		cfg.Kafka.CantabularCSVCreatedTopic,
		kafka.CreateProducerChannels(),
		&kafka.ProducerConfig{
			KafkaVersion:    &cfg.Kafka.Version,
			MaxMessageBytes: &cfg.Kafka.MaxBytes,
		},
	); err != nil {
		return fmt.Errorf("error creating kafka producer: %w", err)
	}

	// start kafka logging go-routines
	c.producer.Channels().LogErrors(ctx, "component producer")

	// Create service and initialise it
	c.svc = service.New()
	if err = c.svc.Init(ctx, cfg, BuildTime, GitCommit, Version); err != nil {
		return fmt.Errorf("failed to init service: %w", err)
	}

	c.cfg = cfg

	// wait for producer and consumer to be ready
	<-c.producer.Channels().Ready
	log.Info(ctx, "component-test kafka producer ready")

	return nil
}

func (c *Component) startService(ctx context.Context) {
	defer c.wg.Done()
	c.svc.Start(context.Background(), c.errorChan)

	// blocks until an os interrupt or a fatal error occurs
	select {
	case err := <-c.errorChan:
		err = fmt.Errorf("service error received: %w", err)
		c.svc.Close(ctx)
		panic(fmt.Errorf("unexpected error received from errorChan: %w", err))
	case sig := <-c.signals:
		log.Info(ctx, "os signal received", log.Data{"signal": sig})
	}
	if err := c.svc.Close(ctx); err != nil {
		panic(fmt.Errorf("unexpected error during service graceful shutdown: %w", err))
	}
}

// Close kills the application under test, and then it shuts down the testing consumer and producer.
func (c *Component) Close() {
	ctx := context.Background()

	// close producer
	if err := c.producer.Close(ctx); err != nil {
		log.Error(ctx, "error closing kafka producer", err)
	}

	// kill application
	c.signals <- os.Interrupt

	// wait for graceful shutdown to finish (or timeout)
	c.wg.Wait()
}

// Reset initialises the service under test, the api mocks and then starts the service
func (c *Component) Reset() error {
	ctx := context.Background()

	if err := c.initService(ctx); err != nil {
		return fmt.Errorf("failed to initialise service: %w", err)
	}

	c.DatasetAPI.Reset()

	// run application in separate goroutine
	c.wg.Add(1)
	go c.startService(ctx)

	return nil
}
