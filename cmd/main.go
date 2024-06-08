package main

import (
	"fmt"
	"log"
	"net/http"
	"propertiesGo/pkg/handler"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)


func homeLink(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "Welcome home!")
}

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)

	router.HandleFunc("/api/login", handler.Login).Methods(http.MethodPost)
	router.HandleFunc("/api/register", handler.Register).Methods(http.MethodPost)
	router.HandleFunc("/api/changePassword", handler.ChangePassword).Methods(http.MethodPost)
	router.HandleFunc("/api/disableAccount", handler.DisableAccount).Methods(http.MethodPost)
	router.HandleFunc("/api/deleteAccount", handler.DeleteAccount).Methods(http.MethodPost)
	router.HandleFunc("/api/changeAvatar", handler.ChangeAvatar).Methods(http.MethodPost)
	router.HandleFunc("/api/getInfoUser", handler.GetInfoUser).Methods(http.MethodGet)
	router.HandleFunc("/api/changeInfo", handler.ChangeInfo).Methods(http.MethodPost)

	router.HandleFunc("/api/news", handler.GetAllNews).Methods(http.MethodGet)
	router.HandleFunc("/api/news/{id}", handler.GetNewsByID).Methods(http.MethodGet)

	router.HandleFunc("/api/properties", handler.GetAllProperties).Methods(http.MethodGet)
	router.HandleFunc("/api/propertiesMain", handler.GetAllPropertiesHome).Methods(http.MethodGet)
	router.HandleFunc("/api/properties/{id}", handler.GetPropertiesDetail).Methods(http.MethodGet)
	router.HandleFunc("/api/postProperties", handler.PostProperties).Methods(http.MethodPost)
	router.HandleFunc("/api/getPostedProperty", handler.GetPostedProperty).Methods(http.MethodGet)

	c := cors.New(cors.Options{
        AllowedOrigins: []string{"http://localhost:8080"},
        AllowCredentials: true,
    })

    handler := c.Handler(router)

	log.Fatal(http.ListenAndServe(":5000", handler))
}