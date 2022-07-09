package main

import (
	"encoding/json"
	"log"
	"net/http"

	httputils "github.com/svetlana-rezvaya/go-dice-backend/http-utils"
	"github.com/svetlana-rezvaya/go-dice-cli"
	"github.com/svetlana-rezvaya/go-dice-cli/statistics"
)

type result struct {
	Throws     []int
	Statistics statistics.Statistics
}

func main() {
	http.HandleFunc("/api/v1/dice", func(writer http.ResponseWriter, request *http.Request) {
		throws, err := httputils.GetIntFormValue(request, "throws")
		if err != nil {
			httputils.HandleError(writer, http.StatusBadRequest, "unable to get throw count: %s", err)
			return
		}

		faces, err := httputils.GetIntFormValue(request, "faces")
		if err != nil {
			httputils.HandleError(writer, http.StatusBadRequest, "unable to get face count: %s", err)
			return
		}

		throwResults := dice.GenerateDiceThrows(throws, faces)
		throwStatistics := statistics.CollectStatistics(throwResults)

		throwTotalResult := result{Throws: throwResults, Statistics: throwStatistics}
		responseBytes, err := json.Marshal(throwTotalResult)
		if err != nil {
			httputils.HandleError(writer, http.StatusInternalServerError, "unable to marshal the response: %s", err)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.Write(responseBytes)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
