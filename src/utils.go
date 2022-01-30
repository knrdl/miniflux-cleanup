package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func apiErrorResponse(writer http.ResponseWriter, err error, status int) {
	apiTxtResponse(writer, err.Error(), status)
}

func apiTxtResponse(writer http.ResponseWriter, content string, status int) {
	writer.WriteHeader(status)
	writer.Header().Set("Content-Type", "text/plain")
	if _, err := writer.Write([]byte(content)); err != nil {
		fmt.Printf(err.Error())
	}
}

func apiJsonResponse(writer http.ResponseWriter, content interface{}, status int) {
	writer.WriteHeader(status)
	writer.Header().Set("Content-Type", "application/json")
	res, _ := json.Marshal(content)
	if _, err := writer.Write(res); err != nil {
		fmt.Printf(err.Error())
	}
}

func contains(s []int64, e int64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
