package api_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-cantabular-metadata-exporter/api"

	. "github.com/smartystreets/goconvey/convey"
)

type testResponse struct {
	Message string `json:"message"`
}

type testError struct {
	err  error
	resp string
	code int
}

func (e testError) Error() string {
	if e.err == nil {
		return ""
	}
	return e.err.Error()
}

func (e testError) Code() int {
	return e.code
}

func (e testError) Response() string {
	return e.resp
}

func TestError(t *testing.T) {

	Convey("Given a valid context and response writer", t, func() {
		ctx := context.Background()
		w := httptest.NewRecorder()

		Convey("Given a standard Go error", func() {
			err := errors.New("test error")

			Convey("when RespondError is called", func() {
				api.Error(ctx, w, err)

				Convey("the response writer should record status code 500 and appropriate error response body", func() {
					expectedCode := http.StatusInternalServerError
					expectedBody := `{"errors":["test error"]}`

					So(w.Code, ShouldEqual, expectedCode)
					So(w.Body.String(), ShouldResemble, expectedBody)
				})

			})
		})

		Convey("Given an error that satisfies interfaces providing Code() and Response() functions", func() {
			err := testError{
				err:  errors.New("test error"),
				resp: "test response",
				code: http.StatusUnauthorized,
			}

			Convey("when RespondError is called", func() {
				api.Error(ctx, w, err)

				Convey("the response writer should record the appropriate status code and response message", func() {
					expectedCode := http.StatusUnauthorized
					expectedBody := `{"errors":["test response"]}`

					So(w.Code, ShouldEqual, expectedCode)
					So(w.Body.String(), ShouldResemble, expectedBody)
				})
			})
		})
	})
}
