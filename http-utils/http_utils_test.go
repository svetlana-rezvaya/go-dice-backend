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
	type args struct {
		status    int
		format    string
		arguments []interface{}
	}

	tests := []struct {
		name             string
		args             args
		wantedLogMessage string
		wantedResponse   *http.Response
	}{
		{
			name: "succes",
			args: args{
				status:    http.StatusNotFound,
				format:    "test: %d %s",
				arguments: []interface{}{23, "one"},
			},
			wantedLogMessage: "test: 23 one\n",
			wantedResponse: &http.Response{
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
			},
		},
	}
	for _, testData := range tests {
		previousLogOutput := log.Writer()
		defer func() { log.SetOutput(previousLogOutput) }()

		var logBuffer bytes.Buffer
		log.SetOutput(&logBuffer)

		responseRecorder := httptest.NewRecorder()
		HandleError(
			responseRecorder,
			testData.args.status,
			testData.args.format,
			testData.args.arguments...,
		)

		logMessage := logBuffer.String()
		if !strings.HasSuffix(logMessage, testData.wantedLogMessage) {
			test.Logf(
				"failed %q:\n  expected: %+v\n  actual: %+v",
				testData.name,
				testData.wantedLogMessage,
				logMessage,
			)
			test.Fail()
		}

		response := responseRecorder.Result()
		if !reflect.DeepEqual(response, testData.wantedResponse) {
			test.Logf(
				"failed %q:\n  expected: %+v\n  actual: %+v",
				testData.name,
				testData.wantedResponse,
				response,
			)
			test.Fail()
		}
	}
}
