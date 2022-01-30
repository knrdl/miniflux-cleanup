package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	checkEnv()

	startCronjob()

	r := mux.NewRouter()
	api := r.PathPrefix("/api/").Subrouter()
	api.HandleFunc("/config/preview", handlePreview).Methods("POST")
	api.HandleFunc("/config", handleRetrieveConfig).Methods("GET")
	api.HandleFunc("/config", handleSaveConfig).Methods("PUT")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./www")))
	http.Handle("/", r)
	log.Println("Server started on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
