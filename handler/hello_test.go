package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/handler"

	. "github.com/smartystreets/goconvey/convey"
)

var ctx = context.Background()

// TODO: remove hello world handler test
func TestHello(t *testing.T) {

	Convey("Given a Hello handler ", t, func() {
		hello := handler.NewHello("Hello, World!")

		Convey("when a good response is returned", func() {
			req := httptest.NewRequest("GET", "http://localhost:8080/hello", nil)
			resp := httptest.NewRecorder()

			hello.Get(resp, req)

			So(resp.Code, ShouldEqual, http.StatusOK)
			So(resp.Body.String(), ShouldResemble, `{"message":" says: Hello, World!"}`)
		})

	})
}
