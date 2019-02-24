package main

import (
	"log"
	"net/http"

	"github.com/elumbantoruan/app-metadata/handlers"
	"github.com/elumbantoruan/app-metadata/repository"

	"github.com/gorilla/mux"
)

func main() {

	m, err := registerHandlers()
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/", m)

	err = http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func registerHandlers() (*mux.Router, error) {
	m := mux.NewRouter()

	// initialize in-memory metadata repository
	inMem := repository.NewInMemoryMetadataRepository()
	// initialize metadata handler and inject the implementation of repository interface
	appMd := handlers.NewMetadataHandler(inMem)

	// Register app-metadata resource
	m.HandleFunc("/app-metadata", appMd.HandlePostMetadata).Methods("POST")
	m.HandleFunc("/app-metadata/{appID}", appMd.HandlePutMetadata).Methods("PUT")
	m.HandleFunc("/app-metadata/{appID}", appMd.HandleGetMetadata).Methods("GET")
	m.HandleFunc("/app-metadata", appMd.HandleGetAllMetadata).Methods("GET")
	m.HandleFunc("/app-metadata/{appID}", appMd.HandleDeleteMetadata).Methods("DELETE")

	return m, nil
}
