package main

import (
	"fmt"
	"log"
	"net/http"
	"propertiesGo/pkg/handler"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)


func homeLink(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "Welcome home!")
}

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/api/news", handler.GetAllNews).Methods(http.MethodGet)
	router.HandleFunc("/api/news/{id}", handler.GetNewsByID).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(":5000", handlers.CORS()(router)))
}