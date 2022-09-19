package httputils

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestLoggingMiddleware(test *testing.T) {
	previousLogOutput := log.Writer()
	defer func() { log.SetOutput(previousLogOutput) }()

	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)

	responseRecorder := httptest.NewRecorder()
	handler := LoggingMiddleware(func(
		writer http.ResponseWriter,
		request *http.Request,
	) {
		writer.Write([]byte("Hello, world!"))
	})
	handler(
		responseRecorder,
		httptest.NewRequest(http.MethodGet, "http://example.com/test", nil),
	)

	logMessage := logBuffer.String()
	wantedLogMessage := "GET http://example.com/test"
	if !strings.Contains(logMessage, wantedLogMessage) {
		test.Logf(
			"failed:\n  expected: %+v\n  actual: %+v",
			wantedLogMessage,
			logMessage,
		)
		test.Fail()
	}

	response := responseRecorder.Result()
	wantedResponse := &http.Response{
		Status: strconv.Itoa(http.StatusOK) + " " +
			http.StatusText(http.StatusOK),
		StatusCode: http.StatusOK,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header: http.Header{
			"Content-Type": []string{"text/plain; charset=utf-8"},
		},
		Body:          ioutil.NopCloser(bytes.NewReader([]byte("Hello, world!"))),
		ContentLength: -1,
	}
	if !reflect.DeepEqual(response, wantedResponse) {
		test.Logf(
			"failed:\n  expected: %+v\n  actual: %+v",
			wantedResponse,
			response,
		)
		test.Fail()
	}
}

func TestGetIntFormValue(test *testing.T) {
	type args struct {
		request *http.Request
		key     string
	}
	type data struct {
		name         string
		args         args
		wantedValue  int
		wantedErrStr string
	}

	tests := []data{
		data{
			name: "success",
			args: args{
				request: httptest.NewRequest(http.MethodGet, "/test?key=23", nil),
				key:     "key",
			},
			wantedValue:  23,
			wantedErrStr: "",
		},
		data{
			name: "error with a missed form value",
			args: args{
				request: httptest.NewRequest(http.MethodGet, "/test", nil),
				key:     "key",
			},
			wantedValue:  0,
			wantedErrStr: "form value is missing",
		},
		data{
			name: "error with an incorrect integer value",
			args: args{
				request: httptest.NewRequest(http.MethodGet, "/test?key=value", nil),
				key:     "key",
			},
			wantedValue: 0,
			wantedErrStr: "incorrect integer value: " +
				`strconv.Atoi: parsing "value": invalid syntax`,
		},
	}
	for _, testData := range tests {
		value, err := GetIntFormValue(testData.args.request, testData.args.key)

		if value != testData.wantedValue {
			test.Logf(
				"failed %q:\n  expected: %+v\n  actual: %+v",
				testData.name,
				testData.wantedValue,
				value,
			)
			test.Fail()
		}

		wantedErr := testData.wantedErrStr != ""
		if !wantedErr && err != nil ||
			wantedErr && (err == nil || err.Error() != testData.wantedErrStr) {
			test.Logf(
				"failed %q:\n  expected: %+v\n  actual: %+v",
				testData.name,
				testData.wantedErrStr,
				err,
			)
			test.Fail()
		}
	}
}

func TestHandleError(test *testing.T) {
	previousLogOutput := log.Writer()
	defer func() { log.SetOutput(previousLogOutput) }()

	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)

	responseRecorder := httptest.NewRecorder()
	HandleError(responseRecorder, http.StatusNotFound, "test: %d %s", 23, "one")

	logMessage := logBuffer.String()
	wantedLogMessage := "test: 23 one\n"
	if !strings.HasSuffix(logMessage, wantedLogMessage) {
		test.Logf(
			"failed:\n  expected: %+v\n  actual: %+v",
			wantedLogMessage,
			logMessage,
		)
		test.Fail()
	}

	response := responseRecorder.Result()
	wantedResponse := &http.Response{
		Status: strconv.Itoa(http.StatusNotFound) + " " +
			http.StatusText(http.StatusNotFound),
		StatusCode: http.StatusNotFound,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header: http.Header{
			"Content-Type":           []string{"text/plain; charset=utf-8"},
			"X-Content-Type-Options": []string{"nosniff"},
		},
		Body:          ioutil.NopCloser(bytes.NewReader([]byte("test: 23 one\n"))),
		ContentLength: -1,
	}
	if !reflect.DeepEqual(response, wantedResponse) {
		test.Logf(
			"failed:\n  expected: %+v\n  actual: %+v",
			wantedResponse,
			response,
		)
		test.Fail()
	}
}
