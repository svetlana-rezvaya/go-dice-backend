package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/svetlana-rezvaya/go-dice-cli"
	"github.com/svetlana-rezvaya/go-dice-cli/statistics"
)

type result struct {
	Throws     []int
	Statistics statistics.Statistics
}

func main() {
	http.HandleFunc("/api/v1/dice", func(writer http.ResponseWriter, request *http.Request) {
		throwsAsString := request.FormValue("throws")
		if throwsAsString == "" {
			const errMessage = "throw count is missing"
			log.Print(errMessage)
			http.Error(writer, errMessage, http.StatusBadRequest)

			return
		}

		throws, err := strconv.Atoi(throwsAsString)
		if err != nil {
			errMessage := fmt.Sprintf("incorrect throw count: %s", err)
			log.Print(errMessage)
			http.Error(writer, errMessage, http.StatusBadRequest)

			return
		}

		facesAsString := request.FormValue("faces")
		if facesAsString == "" {
			const errMessage = "face count is missing"
			log.Print(errMessage)
			http.Error(writer, errMessage, http.StatusBadRequest)

			return
		}

		faces, err := strconv.Atoi(facesAsString)
		if err != nil {
			errMessage := fmt.Sprintf("incorrect face count: %s", err)
			log.Print(errMessage)
			http.Error(writer, errMessage, http.StatusBadRequest)

			return
		}

		throwResults := dice.GenerateDiceThrows(throws, faces)
		throwStatistics := statistics.CollectStatistics(throwResults)

		throwTotalResult := result{Throws: throwResults, Statistics: throwStatistics}
		responseBytes, err := json.Marshal(throwTotalResult)
		if err != nil {
			errMessage := fmt.Sprintf("unable to marshal the response: %s", err)
			log.Print(errMessage)
			http.Error(writer, errMessage, http.StatusInternalServerError)

			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.Write(responseBytes)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
