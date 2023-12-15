package steps

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/ONSdigital/dp-cantabular-metadata-exporter/config"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/features/mock"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/service"

	cmptest "github.com/ONSdigital/dp-component-test"
	kafka "github.com/ONSdigital/dp-kafka/v4"
	"github.com/ONSdigital/log.go/v2/log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/maxcnunes/httpfake"
)

const (
	ComponentTestGroup    = "component-test" // kafka group name for the component test consumer
	DrainTopicTimeout     = 5 * time.Second  // maximum time to wait for a topic to be drained
	DrainTopicMaxMessages = 1000             // maximum number of messages that will be drained from a topic
	MinioCheckRetries     = 5                // maximum number of retires to validate that a file is present in minio
	WaitEventTimeout      = 15 * time.Second // maximum time that the component test consumer will wait for a kafka event
)

var (
	BuildTime = "1625046891"
	GitCommit = "7434fe334d9f51b7239f978094ea29d10ac33b16"
	Version   = ""
)

type Component struct {
	cmptest.ErrorFeature
	DatasetAPI       *httpfake.HTTPFake
	FilterAPI        *httpfake.HTTPFake
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
	ctx              context.Context
}

func NewComponent(t *testing.T) *Component {
	return &Component{
		errorChan:        make(chan error),
		svcStarted:       make(chan bool, 1),
		DatasetAPI:       httpfake.New(httpfake.WithTesting(t)),
		FilterAPI:        httpfake.New(httpfake.WithTesting(t)),
		wg:               &sync.WaitGroup{},
		waitEventTimeout: WaitEventTimeout,
		minioRetries:     MinioCheckRetries,
		ctx:              context.Background(),
	}
}

// initService initialises the server, the mocks and waits for the dependencies to be ready
func (c *Component) initService(ctx context.Context) error {
	// register interrupt signals
	c.signals = make(chan os.Signal, 1)
	signal.Notify(c.signals, os.Interrupt, syscall.SIGTERM)

	var err error
	c.cfg, err = config.Get()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}
	log.Info(ctx, "config used by component tests", log.Data{"cfg": c.cfg})

	c.cfg.DatasetAPIURL = c.DatasetAPI.ResolveURL("")
	c.cfg.FilterAPIURL = c.FilterAPI.ResolveURL("")

	service.GetGenerator = func() service.Generator {
		return &mock.Generator{}
	}

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(c.cfg.MinioAccessKey, c.cfg.MinioSecretKey, ""),
		Endpoint:         aws.String(c.cfg.LocalObjectStore),
		Region:           aws.String(c.cfg.AWSRegion),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}

	s, err := session.NewSession(s3Config)
	if err != nil {
		return fmt.Errorf("failed to create new aws session: %w", err)
	}
	c.S3Downloader = s3manager.NewDownloader(s)

	// producer for triggering test events that will be consumed by the service
	if c.producer, err = kafka.NewProducer(
		ctx,
		&kafka.ProducerConfig{
			BrokerAddrs:       c.cfg.Kafka.Addr,
			Topic:             c.cfg.Kafka.CantabularCSVCreatedTopic,
			MinBrokersHealthy: &c.cfg.Kafka.ProducerMinBrokersHealthy,
			KafkaVersion:      &c.cfg.Kafka.Version,
			MaxMessageBytes:   &c.cfg.Kafka.MaxBytes,
		},
	); err != nil {
		return fmt.Errorf("error creating kafka producer: %w", err)
	}

	// consumer for receiving cantabular-output-created events, produced by the service
	// (expected to be generated by the service under test)
	// use kafkaOldest to make sure we consume all the messages
	kafkaOffset := kafka.OffsetOldest
	if c.consumer, err = kafka.NewConsumerGroup(
		ctx,
		&kafka.ConsumerGroupConfig{
			BrokerAddrs:       c.cfg.Kafka.Addr,
			Topic:             c.cfg.Kafka.CantabularCSVWCreatedTopic,
			GroupName:         ComponentTestGroup,
			MinBrokersHealthy: &c.cfg.Kafka.ConsumerMinBrokersHealthy,
			KafkaVersion:      &c.cfg.Kafka.Version,
			Offset:            &kafkaOffset,
		},
	); err != nil {
		return fmt.Errorf("error creating kafka consumer: %w", err)
	}

	// start consumer group
	if err := c.consumer.Start(); err != nil {
		return fmt.Errorf("error starting kafka consumer: %w", err)
	}

	// start kafka logging go-routines
	c.producer.LogErrors(ctx)
	c.consumer.LogErrors(ctx)

	// Create service and initialise it
	c.svc = service.New()
	if err = c.svc.Init(ctx, c.cfg, BuildTime, GitCommit, Version); err != nil {
		return fmt.Errorf("failed to init service: %w", err)
	}

	// wait for component-test producer to be initialised and consumer to be in consuming state
	<-c.producer.Channels().Initialised
	log.Info(ctx, "component-test kafka producer initialised")
	c.consumer.StateWait(kafka.Consuming)
	log.Info(ctx, "component-test kafka consumer is in consuming state")

	return nil
}

func (c *Component) startService(ctx context.Context) error {
	if err := c.svc.Start(ctx, c.errorChan); err != nil {
		return fmt.Errorf("unexpected error while starting service: %w", err)
	}

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()

		// blocks until an os interrupt or a fatal error occurs
		select {
		case err := <-c.errorChan:
			if errClose := c.svc.Close(ctx); errClose != nil {
				log.Warn(ctx, "error closing server during error handing", log.Data{"close_error": errClose})
			}
			panic(fmt.Errorf("unexpected error received from errorChan: %w", err))
		case sig := <-c.signals:
			log.Info(ctx, "os signal received", log.Data{"signal": sig})
		}
		if err := c.svc.Close(ctx); err != nil {
			panic(fmt.Errorf("unexpected error during service graceful shutdown: %w", err))
		}
	}()

	return nil
}

// drainTopic drains the provided topic and group of any residual messages between scenarios.
// Prevents future tests failing if previous tests fail unexpectedly and
// leave messages in the queue.
//
// A temporary batch consumer is used, that is created and closed within this func
// A maximum of DrainTopicMaxMessages messages will be drained from the provided topic and group.
//
// This method accepts a waitGroup pionter. If it is not nil, it will wait for the topic to be drained
// in a new go-routine, which will be added to the waitgroup. If it is nil, execution will be blocked
// until the topic is drained (or time out expires)
func (c *Component) drainTopic(ctx context.Context, topic, group string, wg *sync.WaitGroup) error {
	msgs := []kafka.Message{}

	kafkaOffset := kafka.OffsetOldest
	batchSize := DrainTopicMaxMessages
	batchWaitTime := DrainTopicTimeout
	drainer, err := kafka.NewConsumerGroup(
		ctx,
		&kafka.ConsumerGroupConfig{
			BrokerAddrs:   c.cfg.Kafka.Addr,
			Topic:         topic,
			GroupName:     group,
			KafkaVersion:  &c.cfg.Kafka.Version,
			Offset:        &kafkaOffset,
			BatchSize:     &batchSize,
			BatchWaitTime: &batchWaitTime,
		},
	)
	if err != nil {
		return fmt.Errorf("error creating kafka consumer to drain topic: %w", err)
	}

	// register batch handler with 'drained channel'
	drained := make(chan struct{})
	if err := drainer.RegisterBatchHandler(
		ctx,
		func(ctx context.Context, batch []kafka.Message) error {
			defer close(drained)
			msgs = append(msgs, batch...)
			return nil
		},
	); err != nil {
		return fmt.Errorf("error creating kafka drainer: %w", err)
	}

	// start drainer consumer group
	if err := drainer.Start(); err != nil {
		log.Error(ctx, "error starting kafka drainer", err)
	}

	// start kafka logging go-routines
	drainer.LogErrors(ctx)

	// waitUntilDrained is a func that will wait until the batch is consumed or the timeout expires
	// (with 100 ms of extra time to allow any in-flight drain)
	waitUntilDrained := func() {
		drainer.StateWait(kafka.Consuming)
		log.Info(ctx, "drainer is consuming", log.Data{"topic": topic, "group": group})

		select {
		case <-time.After(DrainTopicTimeout + 100*time.Millisecond):
			log.Info(ctx, "drain timeout has expired (no messages drained)")
		case <-drained:
			log.Info(ctx, "message(s) have been drained")
		}

		defer func() {
			log.Info(ctx, "drained topic", log.Data{
				"len":      len(msgs),
				"messages": msgs,
				"topic":    topic,
				"group":    group,
			})
		}()

		if err := drainer.Close(ctx); err != nil {
			log.Warn(ctx, "error closing drain consumer", log.Data{"err": err})
		}

		<-drainer.Channels().Closed
		log.Info(ctx, "drainer is closed")
	}

	// sync wait if wg is not provided
	if wg == nil {
		waitUntilDrained()
		return nil
	}

	// async wait if wg is provided
	wg.Add(1)
	go func() {
		defer wg.Done()
		waitUntilDrained()
	}()
	return nil
}

// Close kills the application under test, and then it shuts down the testing consumer and producer.
func (c *Component) Close() {
	// kill application
	c.signals <- os.Interrupt

	// wait for graceful shutdown to finish (or timeout)
	c.wg.Wait()

	// stop listening to consumer, waiting for any in-flight message to be committed
	if err := c.consumer.StopAndWait(); err != nil {
		log.Error(c.ctx, "error stopping kafka consumer", err)
	}

	// close producer
	if err := c.producer.Close(c.ctx); err != nil {
		log.Error(c.ctx, "error closing kafka producer", err)
	}

	// close consumer
	if err := c.consumer.Close(c.ctx); err != nil {
		log.Error(c.ctx, "error closing kafka consumer", err)
	}

	// drain topics in parallel
	wg := &sync.WaitGroup{}
	if err := c.drainTopic(c.ctx, c.cfg.Kafka.CantabularCSVCreatedTopic, ComponentTestGroup, wg); err != nil {
		log.Error(c.ctx, "error draining topic", err, log.Data{
			"topic": c.cfg.Kafka.CantabularCSVCreatedTopic,
			"group": ComponentTestGroup,
		})
	}
	if err := c.drainTopic(c.ctx, c.cfg.Kafka.CantabularCSVWCreatedTopic, c.cfg.Kafka.CantabularMetadataExportGroup, wg); err != nil {
		log.Error(c.ctx, "error draining topic", err, log.Data{
			"topic": c.cfg.Kafka.CantabularCSVWCreatedTopic,
			"group": c.cfg.Kafka.CantabularMetadataExportGroup,
		})
	}
	wg.Wait()
}

// Reset initialises the service under test, the api mocks and then starts the service
func (c *Component) Reset() error {
	if err := c.initService(c.ctx); err != nil {
		return fmt.Errorf("failed to initialise service: %w", err)
	}

	c.DatasetAPI.Reset()
	c.FilterAPI.Reset()

	return nil
}
