package handler_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-cantabular-metadata-exporter/handler"

	"github.com/go-chi/chi/v5"
	. "github.com/smartystreets/goconvey/convey"
)

var ctx = context.Background()

func TestMetadata(t *testing.T) {

	Convey("Given a Metadata handler routed to /metadata", t, func() {
		metadata := handler.NewMetadata()

		r := chi.NewRouter()
		r.Post("/metadata", metadata.Post)

		ts := httptest.NewServer(r)
		defer ts.Close()

		Convey("when a request is made to /metadata", func() {
			resp := testDoRequest(t, ts, http.MethodPost, "/metadata", bytes.NewReader([]byte("{}")))

			Convey("status code should equal 202 Status Accepted", func() {
				So(resp.StatusCode, ShouldEqual, http.StatusAccepted)
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
