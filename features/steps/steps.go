package steps

import (
	"fmt"
	"net/http"
	"time"
	"errors"
	"context"
	"encoding/json"

	"github.com/ONSdigital/dp-cantabular-metadata-exporter/schema"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/event"

	"github.com/ONSdigital/log.go/v2/log"

	"github.com/cucumber/godog"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (c *Component) RegisterSteps(ctx *godog.ScenarioContext) {
	ctx.Step(
		`^the following metadata document with dataset id "([^"]*)", edition "([^"]*)" and version "([^"]*)" is available from dp-dataset-api:$`,
		c.theFollowingMetadataDocumentIsAvailable,
	)
	ctx.Step(
		`^the following version document with dataset id "([^"]*)", edition "([^"]*)" and version "([^"]*)" is available from dp-dataset-api:$`,
		c.theFollowingVersionDocumentIsAvailable,
	)
	ctx.Step(`^this cantabular-metadata-export event is consumed:$`, c.thisCantabularMetadataExportEventIsConsumed)
	ctx.Step(`^a file with filename "([^"]*)" can be seen in minio bucket "([^"]*)"`, c.theFollowingFileCanBeSeenInMinio)
	ctx.Step(`^the following version with dataset id "([^"]*)", edition "([^"]*)" and version "([^"]*)" is updated to dp-dataset-api:$`, c.theFollowingVersionIsUpdated)
}

// theFollowingMetadataDocumentIsAvailable generate a mocked response for dataset API
// GET /datasets/{dataset_id}/editions/{edition}/versions/{version}/metadata
func (c *Component) theFollowingMetadataDocumentIsAvailable(datasetID, edition, version string, md *godog.DocString) error {
	url := fmt.Sprintf(
		"/datasets/%s/editions/%s/versions/%s/metadata",
		datasetID,
		edition,
		version,
	)

	c.DatasetAPI.NewHandler().
		Get(url).
		Reply(http.StatusOK).
		BodyString(md.Content)

	return nil
}

// theFollowingVersionDocumentIsAvailable generate a mocked response for dataset API
// GET /datasets/{dataset_id}/editions/{edition}/versions/{version}
func (c *Component) theFollowingVersionDocumentIsAvailable(datasetID, edition, version string, v *godog.DocString) error {
	url := fmt.Sprintf(
		"/datasets/%s/editions/%s/versions/%s",
		datasetID,
		edition,
		version,
	)

	c.DatasetAPI.NewHandler().
		Get(url).
		Reply(http.StatusOK).
		BodyString(v.Content)

	return nil
}

// theFollowingVersionIsUpdated generate a mocked response for dataset API
// PUT /instances/{id} with the provided instance response
func (c *Component) theFollowingVersionIsUpdated(datasetID, edition, version string, v *godog.DocString) error {
	url := fmt.Sprintf(
		"/datasets/%s/editions/%s/versions/%s",
		datasetID,
		edition,
		version,
	)

	c.DatasetAPI.NewHandler().
		Put(url).
		AssertBody([]byte(v.Content)).
		Reply(http.StatusOK)

	return nil
}

func (c *Component) thisCantabularMetadataExportEventIsConsumed(input *godog.DocString) error {
	ctx := context.Background()

	// testing kafka message that will be produced
	var testEvent event.CantabularMetadataExport
	if err := json.Unmarshal([]byte(input.Content), &testEvent); err != nil {
		return fmt.Errorf("error unmarshaling input to event: %w body: %s", err, input.Content)
	}

	log.Info(ctx, "event to marshal: ", log.Data{
		"event": testEvent,
	})

	// marshal and send message
	b, err := schema.CantabularMetadataExport.Marshal(testEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal event from schema: %w", err)
	}

	log.Info(ctx, "marshalled event: ", log.Data{
		"event": b,
	})

	c.producer.Channels().Output <- b

	return nil
}

func (c *Component) theFollowingFileCanBeSeenInMinio(fileName string, bucket string) error {
	ctx := context.Background()

	var b []byte
	f := aws.NewWriteAtBuffer(b)

	// probe bucket with backoff to give time for event to be processed
	retries := 4
	timeout := 1
	var numBytes int64
	var err error

	for {
		if numBytes, err = c.S3Downloader.Download(f, &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(fileName),
		}); err == nil || retries <= 0 {
			break
		}

		retries--

		log.Info(ctx, "error obtaining file from minio. Retrying.", log.Data{
			"error":        err,
			"retries_left": retries,
		})

		time.Sleep(time.Second * time.Duration(timeout))
		timeout *= 2
	}
	if err != nil {
		return fmt.Errorf(
			"error obtaining file from minio. Last error: %w",
			err,
		)
	}

	if numBytes < 1 {
		return errors.New("file length zero")
	}

	log.Info(ctx, "got file contents", log.Data{
		"contents": string(f.Bytes()),
	})

	return nil
}
