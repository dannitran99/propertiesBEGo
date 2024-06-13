package main

import (
	"fmt"
	"log"
	"net/http"
	"propertiesGo/pkg/handler"
	"propertiesGo/pkg/middleware"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)


func homeLink(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "Welcome home!")
}

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	
	router.HandleFunc("/api/checkVerifyToken", middleware.VerifyJWT(handler.CheckVerifyToken)).Methods(http.MethodPost)

	router.HandleFunc("/api/login", handler.Login).Methods(http.MethodPost)
	router.HandleFunc("/api/register", handler.Register).Methods(http.MethodPost)
	router.HandleFunc("/api/changePassword", middleware.VerifyJWT(handler.ChangePassword)).Methods(http.MethodPost)
	router.HandleFunc("/api/disableAccount", middleware.VerifyJWT(handler.DisableAccount)).Methods(http.MethodPost)
	router.HandleFunc("/api/deleteAccount", middleware.VerifyJWT(handler.DeleteAccount)).Methods(http.MethodPost)
	router.HandleFunc("/api/changeAvatar", middleware.VerifyJWT(handler.ChangeAvatar)).Methods(http.MethodPost)
	router.HandleFunc("/api/getInfoUser", middleware.VerifyJWT(handler.GetInfoUser)).Methods(http.MethodGet)
	router.HandleFunc("/api/changeInfo", middleware.VerifyJWT(handler.ChangeInfo)).Methods(http.MethodPost)

	router.HandleFunc("/api/news", handler.GetAllNews).Methods(http.MethodGet)
	router.HandleFunc("/api/news/{id}", handler.GetNewsByID).Methods(http.MethodGet)

	router.HandleFunc("/api/properties", handler.GetAllProperties).Methods(http.MethodGet)
	router.HandleFunc("/api/propertiesMain", handler.GetAllPropertiesHome).Methods(http.MethodGet)
	router.HandleFunc("/api/properties/{id}", handler.GetPropertiesDetail).Methods(http.MethodGet)
	router.HandleFunc("/api/postProperties", middleware.VerifyJWT(handler.PostProperties)).Methods(http.MethodPost)
	router.HandleFunc("/api/getPostedProperty", middleware.VerifyJWT(handler.GetPostedProperty)).Methods(http.MethodGet)

	router.HandleFunc("/api/registerAgency", middleware.VerifyJWT(handler.RegisterAgency)).Methods(http.MethodPost)
	router.HandleFunc("/api/getContactUser", middleware.VerifyJWT(handler.GetContactUser)).Methods(http.MethodGet)
	
	router.HandleFunc("/api/admin/getRequestAgency", middleware.VerifyJWT(handler.GetRequestAgency)).Methods(http.MethodGet)
	router.HandleFunc("/api/admin/getRequestDisableAccount", middleware.VerifyJWT(handler.GetRequestDisableAccount)).Methods(http.MethodGet)
	
	c := cors.New(cors.Options{
        AllowedOrigins: []string{"http://localhost:8080"},
        AllowCredentials: true,
		AllowedMethods: []string{"POST", "GET", "DELETE"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
    })

    handler := c.Handler(router)

	log.Fatal(http.ListenAndServe(":5000", handler))
}