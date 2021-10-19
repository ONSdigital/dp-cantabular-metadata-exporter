package api

import (
	"errors"
	"net/http"
	"reflect"

	"github.com/ONSdigital/log.go/v2/log"
)

// statusCode is a callback function that allows you to extract
// a status code from an error, or returns 500 as a default
func statusCode(err error) int {
	var cerr coder
	if errors.As(err, &cerr) {
		if code := cerr.Code(); code != 0 {
			return code
		}
	}

	return http.StatusInternalServerError
}

// logData returns logData for an error if there is any. This is used
// to extract log.Data embedded in an error if it implements the dataLogger
// interface
func logData(err error) log.Data {
	var lderr dataLogger
	if errors.As(err, &lderr) {
		return lderr.LogData()
	}

	return nil
}

// unwrapLogData recursively unwraps logData from an error. This allows an
// error to be wrapped with log.Data at each level of the call stack, and
// then extracted and combined here as a single log.Data entry. This allows
// us to log errors only once but maintain the context provided by log.Data
// at each level.
func unwrapLogData(err error) log.Data {
	var data []log.Data

	for err != nil && errors.Unwrap(err) != nil {
		if lderr, ok := err.(dataLogger); ok {
			if d := lderr.LogData(); d != nil {
				data = append(data, d)
			}
		}

		err = errors.Unwrap(err)
	}

	// flatten []log.Data into single log.Data with slice
	// entries for duplicate keyed entries, but not for duplicate
	// key-value pairs
	logData := log.Data{}
	for _, d := range data {
		for k, v := range d {
			if val, ok := logData[k]; ok {
				if !reflect.DeepEqual(val, v) {
					if s, ok := val.([]interface{}); ok {
						s = append(s, v)
						logData[k] = s
					} else {
						logData[k] = []interface{}{val, v}
					}
				}
			} else {
				logData[k] = v
			}
		}
	}

	return logData
}

// errorResponse extracts a specified error response to be returned
// to the caller if present, otherwise returns the default error
// string
func errorResponse(err error) string {
	var rerr responser
	if errors.As(err, &rerr) {
		if resp := rerr.Response(); resp != "" {
			return resp
		}
	}

	return err.Error()
}
