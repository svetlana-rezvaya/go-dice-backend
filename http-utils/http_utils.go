package httputils

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func LoggingMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(
		writer http.ResponseWriter,
		request *http.Request,
	) {
		handler(writer, request)
		log.Printf("%s %s", request.Method, request.URL)
	}
}

// GetIntFormValue ...
func GetIntFormValue(request *http.Request, key string) (int, error) {
	valueAsString := request.FormValue(key)
	if valueAsString == "" {
		return 0, errors.New("form value is missing")
	}

	value, err := strconv.Atoi(valueAsString)
	if err != nil {
		return 0, fmt.Errorf("incorrect integer value: %w", err)
	}

	return value, nil
}

// HandleError ...
func HandleError(
	writer http.ResponseWriter,
	status int,
	format string,
	arguments ...any,
) {
	errMessage := fmt.Sprintf(format, arguments...)
	log.Print(errMessage)
	http.Error(writer, errMessage, status)
}
