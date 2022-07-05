package httputils

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

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
