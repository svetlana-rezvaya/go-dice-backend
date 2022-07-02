package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

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

		fmt.Fprintf(writer, "throws: %d\nfaces: %d", throws, faces)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
