package steps

import (
	"context"
	"fmt"
	"errors"
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
	consumer         kafka.IConsumerGroup
	errorChan        chan error
	svcStarted       chan bool
	svc              *service.Service
	cfg              *config.Config
	wg               *sync.WaitGroup
	signals          chan os.Signal
	waitEventTimeout time.Duration
	minioRetries     int
}

func NewComponent(t *testing.T) *Component {
	return &Component{
		errorChan:        make(chan error),
		svcStarted:       make(chan bool, 1),
		DatasetAPI:       httpfake.New(httpfake.WithTesting(t)),
		wg:               &sync.WaitGroup{},
		waitEventTimeout: time.Second * 2,
		minioRetries:     5,
	}
}

// initService initialises the server, the mocks and waits for the dependencies to be ready
func (c *Component) initService(ctx context.Context) error {
	// register interrupt signals
	c.signals = make(chan os.Signal, 1)
	signal.Notify(c.signals, os.Interrupt, syscall.SIGTERM)

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

	// consumer for receiving events
	cgChannels := kafka.CreateConsumerGroupChannels(1)
	kafkaOffset := kafka.OffsetNewest
	if cfg.Kafka.OffsetOldest {
		kafkaOffset = kafka.OffsetOldest
	}

	if c.consumer, err = kafka.NewConsumerGroup(
		ctx,
		cfg.Kafka.Addr,
		cfg.Kafka.CantabularCSVWCreatedTopic,
		"csvw-created-group",
		cgChannels,
		&kafka.ConsumerGroupConfig{
			KafkaVersion: &cfg.Kafka.Version,
			Offset:       &kafkaOffset,
		},
	); err != nil {
		return fmt.Errorf("error creating kafka consumer: %w", err)
	}

	// start kafka logging go-routines
	c.producer.Channels().LogErrors(ctx, "component producer")
	c.consumer.Channels().LogErrors(ctx, "component consumer")

	// wait for consumer and producer to be ready
	<-c.consumer.Channels().Ready
	log.Info(ctx, "consumer ready")

	<-c.producer.Channels().Ready
	log.Info(ctx, "producer ready")

	// Create service and initialise it
	c.svc = service.New()
	if err = c.svc.Init(ctx, cfg, BuildTime, GitCommit, Version); err != nil {
		return fmt.Errorf("failed to init service: %w", err)
	}

	c.cfg = cfg

	return nil
}

func (c *Component) startService(ctx context.Context) {
	defer c.wg.Done()

	c.svc.Start(context.Background(), c.errorChan)
	// wait for producer to be initialised and consumer to be in consuming state
	<-c.svc.Producer().Channels().Initialised
	log.Info(ctx, "kafka producer initialised")
	<-c.svc.Consumer().Channels().State.Consuming
	log.Info(ctx, "kafka consumer is in consuming state")
//	c.svcStarted <- true

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

// drainTopic drains the topic of any residual messages between scenarios.
// Prevents future tests failing if previous tests fail unexpectedly and
// leave messages in the queue.
func (c *Component) drainTopic(ctx context.Context) error {
	var msgs []interface{}

	defer func() {
		log.Info(ctx, "drained topic", log.Data{
			"len":      len(msgs),
			"messages": msgs,
		})
	}()

	for {
		select {
		case <-time.After(time.Second * 1):
			return nil
		case msg, ok := <-c.consumer.Channels().Upstream:
			if !ok {
				return errors.New("upstream channel closed")
			}

			msgs = append(msgs, msg)
			msg.Commit()
			msg.Release()
		}
	}
}

// Close kills the application under test, and then it shuts down the testing consumer and producer.
func (c *Component) Close() {
	ctx := context.Background()

	if err := c.drainTopic(ctx); err != nil {
		log.Error(ctx, "error draining topic", err)
	}

	// close producer
	if err := c.producer.Close(ctx); err != nil {
		log.Error(ctx, "error closing kafka producer", err)
	}

	// close consumer
	if err := c.consumer.Close(ctx); err != nil {
		log.Error(ctx, "error closing kafka consumer", err)
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
//	c.wg.Add(1)
//	go c.startService(ctx)

	// don't allow scenario to start until svc fully initialised
//	<- c.svcStarted

	return nil
}
