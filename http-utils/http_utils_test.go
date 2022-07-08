package httputils

import (
	"net/http"
	"net/http/httptest"
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
