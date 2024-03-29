package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	httputils "github.com/svetlana-rezvaya/go-dice-backend/http-utils"
	"github.com/svetlana-rezvaya/go-dice-cli"
	"github.com/svetlana-rezvaya/go-dice-cli/statistics"
)

//go:generate swag init --generalInfo main.go --parseDependency --output ../../docs --outputTypes json,yaml

// @title go-dice-backend API
// @version 1.0.0
// @license.name MIT
// @host localhost:8080
// @basePath /api/v1

type result struct {
	Throws     []int
	Statistics statistics.Statistics
}

// @router /dice [POST]
// @param throws query int true "throw count"
// @param faces query int true "face count"
// @produce json
// @success 200 {object} result
// @failure 400 {string} string
// @failure 500 {string} string
func main() {
	rand.Seed(time.Now().UnixNano())

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	http.HandleFunc("/api/v1/dice", httputils.LoggingMiddleware(func(
		writer http.ResponseWriter,
		request *http.Request,
	) {
		throws, err := httputils.GetIntFormValue(request, "throws")
		if err != nil {
			status, format := http.StatusBadRequest, "unable to get throw count: %s"
			httputils.HandleError(writer, status, format, err)

			return
		}

		faces, err := httputils.GetIntFormValue(request, "faces")
		if err != nil {
			status, format := http.StatusBadRequest, "unable to get face count: %s"
			httputils.HandleError(writer, status, format, err)

			return
		}

		throwResults := dice.GenerateDiceThrows(throws, faces)
		throwStatistics := statistics.CollectStatistics(throwResults)

		throwTotalResult := result{Throws: throwResults, Statistics: throwStatistics}
		responseBytes, err := json.Marshal(throwTotalResult)
		if err != nil {
			status, format :=
				http.StatusInternalServerError, "unable to marshal the response: %s"
			httputils.HandleError(writer, status, format, err)

			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.Write(responseBytes) // nolint: errcheck
	}))

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
