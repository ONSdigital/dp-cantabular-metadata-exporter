package steps

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ONSdigital/dp-cantabular-metadata-exporter/event"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/schema"
	kafka "github.com/ONSdigital/dp-kafka/v4"

	"github.com/ONSdigital/log.go/v2/log"

	assistdog "github.com/ONSdigital/dp-assistdog"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/cucumber/godog"
	"github.com/google/go-cmp/cmp"
)

func (c *Component) RegisterSteps(ctx *godog.ScenarioContext) {
	ctx.Step(
		`^the service starts`,
		c.theServiceStarts,
	)
	ctx.Step(
		`^the following metadata document with dataset id "([^"]*)", edition "([^"]*)" and version "([^"]*)" is available from dp-dataset-api:$`,
		c.theFollowingMetadataDocumentIsAvailable,
	)
	ctx.Step(
		`^the following version document with dataset id "([^"]*)", edition "([^"]*)" and version "([^"]*)" is available from dp-dataset-api:$`,
		c.theFollowingVersionDocumentIsAvailable,
	)
	ctx.Step(
		`^this cantabular-metadata-export event is consumed:$`,
		c.thisCantabularMetadataExportEventIsConsumed,
	)
	ctx.Step(
		`^a file with filename "([^"]*)" can be seen in minio bucket "([^"]*)"`,
		c.theFollowingFileCanBeSeenInMinio,
	)
	ctx.Step(
		`^the following version with dataset id "([^"]*)", edition "([^"]*)" and version "([^"]*)" will be updated to dp-dataset-api:$`,
		c.theFollowingVersionWillBeUpdated,
	)
	ctx.Step(
		`^the following filter output with id "([^"]*)" will be updated to dp-filter-api:$`,
		c.theFollowingFilterOutputWillBeUpdated,
	)
	ctx.Step(
		`^the following dimensions are available from dataset "([^"]*)" edition "([^"]*)" version "([^"]*)":$`,
		c.theFollowingDimensionsAreAvailable,
	)
	ctx.Step(
		`^the following options response is available for dimension "([^"]*)" for dataset "([^"]*)" edition "([^"]*)" version "([^"]*)" with query params "([^"]*)":$`,
		c.theFollowingOptionsResponseIsAvailable,
	)
	ctx.Step(
		`^these CSVW Created events should be produced:$`,
		c.theseCSVWCreatedEventsShouldBeProduced,
	)
	ctx.Step(
		`^no CSVW Created events should be produced`,
		c.noCSVWCreatedEventsShouldBeProduced,
	)
	ctx.Step(
		`^dataset-api is healthy`,
		c.datasetAPIIsHealthy,
	)
	ctx.Step(
		`^dataset-api is unhealthy`,
		c.datasetAPIIsUnhealthy,
	)
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

func (c *Component) theFollowingOptionsResponseIsAvailable(dimension, datasetID, edition, version, params string, o *godog.DocString) error {
	url := fmt.Sprintf(
		"/datasets/%s/editions/%s/versions/%s/dimensions/%s/options?%s",
		datasetID,
		edition,
		version,
		dimension,
		params,
	)

	c.DatasetAPI.NewHandler().
		Get(url).
		Reply(http.StatusOK).
		BodyString(o.Content)

	return nil
}

// theFollowingDimensionsAreAvailable generates a mocked response for dataset API
// GET /datasets/{dataset_id}/editions/{edition}/versions/{version}/dimensions
func (c *Component) theFollowingDimensionsAreAvailable(datasetID, edition, version string, d *godog.DocString) error {
	url := fmt.Sprintf(
		"/datasets/%s/editions/%s/versions/%s/dimensions",
		datasetID,
		edition,
		version,
	)

	c.DatasetAPI.NewHandler().
		Get(url).
		Reply(http.StatusOK).
		BodyString(d.Content)

	return nil
}

// theFollowingVersionDocumentIsAvailable generates a mocked response for dataset API
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

// theFollowingVersionWillBeUpdated generates a mocked response for dataset API
// PUT /datasets/{dataset_id}/editions/{edition}/versions/{version} with the provided update in the request body
func (c *Component) theFollowingVersionWillBeUpdated(datasetID, edition, version string, v *godog.DocString) error {
	url := fmt.Sprintf(
		"/datasets/%s/editions/%s/versions/%s",
		datasetID,
		edition,
		version,
	)

	c.DatasetAPI.NewHandler().
		Put(url).
		AssertCustom(newPutVersionAssertor([]byte(v.Content))).
		Reply(http.StatusOK)

	return nil
}

func (c *Component) theFollowingFilterOutputWillBeUpdated(filterOutputID string, v *godog.DocString) error {
	url := fmt.Sprintf(
		"/filter-outputs/%s",
		filterOutputID,
	)

	c.FilterAPI.NewHandler().
		Put(url).
		AssertCustom(newPutFilterOutputAssertor([]byte(v.Content))).
		Reply(http.StatusOK)

	return nil
}

func (c *Component) thisCantabularMetadataExportEventIsConsumed(input *godog.DocString) error {
	ctx := context.Background()

	// testing kafka message that will be produced
	var testEvent event.CSVCreated
	if err := json.Unmarshal([]byte(input.Content), &testEvent); err != nil {
		return fmt.Errorf("error unmarshaling input to event: %w body: %s", err, input.Content)
	}

	log.Info(ctx, "event to marshal: ", log.Data{
		"event": testEvent,
	})

	// marshal and send message
	b, err := schema.CSVCreated.Marshal(testEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal event from schema: %w", err)
	}

	log.Info(ctx, "marshalled event: ", log.Data{
		"event": b,
	})

	c.producer.Channels().Output <- kafka.BytesMessage{Value: b, Context: ctx}
	return nil
}

func (c *Component) theFollowingFileCanBeSeenInMinio(fileName, bucket string) error {
	ctx := context.Background()

	var b []byte
	f := aws.NewWriteAtBuffer(b)

	// probe bucket with backoff to give time for event to be processed
	retries := c.minioRetries
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

func (c *Component) theseCSVWCreatedEventsShouldBeProduced(events *godog.Table) error {
	expected, err := assistdog.NewDefault().CreateSlice(new(event.CSVWCreated), events)
	if err != nil {
		return fmt.Errorf("failed to create slice from godog table: %w", err)
	}

	var got []*event.CSVWCreated
	listen := true

	for listen {
		select {
		case <-time.After(c.waitEventTimeout):
			listen = false
		case <-c.consumer.Channels().Closer:
			return errors.New("closer channel closed")
		case msg, ok := <-c.consumer.Channels().Upstream:
			if !ok {
				return errors.New("upstream channel closed")
			}

			var e event.CSVWCreated
			var s = schema.CSVWCreated

			if err := s.Unmarshal(msg.GetData(), &e); err != nil {
				msg.Commit()
				msg.Release()
				return fmt.Errorf("error unmarshalling message: %w", err)
			}

			msg.Commit()
			msg.Release()

			got = append(got, &e)
		}
	}

	if diff := cmp.Diff(got, expected); diff != "" {
		return fmt.Errorf("-got +expected)\n%s", diff)
	}

	return nil
}

func (c *Component) noCSVWCreatedEventsShouldBeProduced() error {
	listen := true

	for listen {
		select {
		case <-time.After(c.waitEventTimeout):
			listen = false
		case <-c.consumer.Channels().Closer:
			return errors.New("closer channel closed")
		case msg, ok := <-c.consumer.Channels().Upstream:
			if !ok {
				return errors.New("upstream channel closed")
			}

			err := fmt.Errorf("unexpected message receieved: %s", msg.GetData())

			msg.Commit()
			msg.Release()

			return err
		}
	}

	return nil
}

func (c *Component) datasetAPIIsHealthy() error {
	resp := `{"status": "OK"}`
	c.DatasetAPI.NewHandler().
		Get("/health").
		Reply(http.StatusOK).
		BodyString(resp)
	return nil
}

func (c *Component) datasetAPIIsUnhealthy() error {
	resp := `{"status": "CRITICAL"}`
	c.DatasetAPI.NewHandler().
		Get("/health").
		Reply(http.StatusInternalServerError).
		BodyString(resp)
	return nil
}

func (c *Component) theServiceStarts() error {
	return c.startService(c.ctx)
}
