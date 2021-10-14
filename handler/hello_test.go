package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"io/ioutil"
	"io"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/handler"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/go-chi/chi/v5"
)

var ctx = context.Background()

// TODO: remove hello world handler test
func TestHello(t *testing.T) {

	Convey("Given a Hello handler routed to /hello/{name}", t, func() {
		hello := handler.NewHello("Hello, World!")

		r := chi.NewRouter()
		r.Get("/hello/{name}", hello.Get)

		ts := httptest.NewServer(r)
		defer ts.Close()

		Convey("when a request is made to /hello/james", func() {
			resp, body := testDoRequest(t, ts, http.MethodGet, "/hello/james", nil)

			Convey("status code should equal 200 and have expected response body", func() {
				So(resp.StatusCode, ShouldEqual, http.StatusOK)
				So(string(body), ShouldResemble, `{"message":"james says: Hello, World!"}`)
			})
			
		})

	})
}

func testDoRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, []byte) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatalf("failed to create request: %s", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to do request: %s", err)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read request body: %s", err)
	}

	defer resp.Body.Close()

	return resp, respBody
}
