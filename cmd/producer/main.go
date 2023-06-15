package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ONSdigital/dp-cantabular-metadata-exporter/config"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/event"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/schema"
	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/log.go/v2/log"
)

const serviceName = "dp-cantabular-metadata-exporter"

func main() {
	log.Namespace = serviceName
	ctx := context.Background()

	if err := run(ctx); err != nil {
		log.Error(ctx, "fatal runtime error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	// Get Config
	cfg, err := config.Get()
	if err != nil {
		log.Error(ctx, "error getting config", err)
		os.Exit(1)
	}

	// Create Kafka Producer
	pConfig := &kafka.ProducerConfig{
		BrokerAddrs:     cfg.Kafka.Addr,
		Topic:           cfg.Kafka.CantabularCSVCreatedTopic,
		KafkaVersion:    &cfg.Kafka.Version,
		MaxMessageBytes: &cfg.Kafka.MaxBytes,
	}
	if cfg.Kafka.SecProtocol == config.KafkaTLSProtocolFlag {
		pConfig.SecurityConfig = kafka.GetSecurityConfig(
			cfg.Kafka.SecCACerts,
			cfg.Kafka.SecClientCert,
			cfg.Kafka.SecClientKey,
			cfg.Kafka.SecSkipVerify,
		)
	}
	kafkaProducer, err := kafka.NewProducer(ctx, pConfig)
	if err != nil {
		return fmt.Errorf("failed to create kafka producer: %w", err)
	}

	// kafka error logging go-routines
	kafkaProducer.LogErrors(ctx)

	// Wait for producer to be initialised plus 500ms
	<-kafkaProducer.Channels().Initialised
	time.Sleep(500 * time.Millisecond)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		e := scanEvent(scanner)
		log.Info(ctx, "sending event", log.Data{"event": e, "topic": cfg.Kafka.CantabularCSVCreatedTopic})

		bytes, err := schema.CSVCreated.Marshal(e)
		if err != nil {
			return fmt.Errorf("failed to marshal event: %w", err)
		}

		// Send bytes to Output channel
		kafkaProducer.Channels().Output <- bytes
	}
}

// scanEvent creates a HelloCalled event according to the user input
func scanEvent(scanner *bufio.Scanner) *event.CSVCreated {
	fmt.Println("--- [Send Kafka CantabularCsvCreated] ---")

	e := &event.CSVCreated{}

	fmt.Println("Please type the instance_id")
	fmt.Printf("$ ")
	scanner.Scan()
	e.InstanceID = scanner.Text()

	fmt.Println("Please type the dataset_id")
	fmt.Printf("$ ")
	scanner.Scan()
	e.DatasetID = scanner.Text()

	fmt.Println("Please type the edition")
	fmt.Printf("$ ")
	scanner.Scan()
	e.Edition = scanner.Text()

	fmt.Println("Please type the version")
	fmt.Printf("$ ")
	scanner.Scan()
	e.Version = scanner.Text()

	for {
		fmt.Println("Please type the row_count")
		fmt.Printf("$ ")
		scanner.Scan()
		i, err := strconv.ParseInt(scanner.Text(), 10, 32)
		if err != nil {
			fmt.Println("Wrong value provided for row_count. Value must be int32")
			continue
		}
		e.RowCount = int32(i)
		break
	}

	return e
}
