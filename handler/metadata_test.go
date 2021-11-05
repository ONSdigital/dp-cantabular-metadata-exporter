package handler_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-cantabular-metadata-exporter/handler"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/event"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/schema"
	"github.com/ONSdigital/dp-kafka/v2/kafkatest"

	"github.com/go-chi/chi/v5"
	. "github.com/smartystreets/goconvey/convey"
)

var ctx = context.Background()

func TestMetadata(t *testing.T) {

	producer := kafkatest.NewMessageProducer(true)

	Convey("Given a Metadata handler routed to /metadata", t, func(c C) {
		metadata := handler.NewMetadata(producer)

		r := chi.NewRouter()
		r.Post("/metadata", metadata.Post)

		ts := httptest.NewServer(r)
		defer ts.Close()

		Convey("when a request is made to /metadata", func(c C) {
			go func(c C){
				body := []byte(`{"dataset_id":"cantabular-example-1","edition":"2021","version":1}`)
				resp := testDoRequest(t, ts, http.MethodPost, "/metadata", bytes.NewReader(body))

				c.Convey("The the returned status code should equal 202 Status Accepted", func(c C) {
					c.So(resp.StatusCode, ShouldEqual, http.StatusAccepted)
				})
			}(c)

			expected := event.CantabularMetadataExport{
				DatasetID: "cantabular-example-1",
				Edition:   "2021",
				Version:   1,
			}
			Convey("And the expected message is produced", func() {
				b := <-producer.Channels().Output
				var got event.CantabularMetadataExport
				err := schema.CantabularMetadataExport.Unmarshal(b, &got)
				So(err, ShouldBeNil)
				So(got, ShouldResemble, expected)
			})
		})
	})
}

func testDoRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) *http.Response {
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatalf("failed to create request: %s", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to do request: %s", err)
	}

	return resp
}
