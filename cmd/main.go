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
	router.HandleFunc("/api/postNews", middleware.VerifyJWT(handler.PostNews)).Methods(http.MethodPost)
	router.HandleFunc("/api/setPinnedNews", middleware.VerifyJWT(handler.SetPinnedNews)).Methods(http.MethodPost)

	router.HandleFunc("/api/properties", handler.GetAllProperties).Methods(http.MethodGet)
	router.HandleFunc("/api/propertiesMain", handler.GetAllPropertiesHome).Methods(http.MethodGet)
	router.HandleFunc("/api/properties/{id}", handler.GetPropertiesDetail).Methods(http.MethodGet)
	router.HandleFunc("/api/postProperties", middleware.VerifyJWT(handler.PostProperties)).Methods(http.MethodPost)
	router.HandleFunc("/api/getPostedProperty", middleware.VerifyJWT(handler.GetPostedProperty)).Methods(http.MethodGet)

	router.HandleFunc("/api/registerAgency", middleware.VerifyJWT(handler.RegisterAgency)).Methods(http.MethodPost)
	router.HandleFunc("/api/registerEnterprise", middleware.VerifyJWT(handler.RegisterEnterprise)).Methods(http.MethodPost)
	router.HandleFunc("/api/updateAgency", middleware.VerifyJWT(handler.UpdateAgency)).Methods(http.MethodPost)
	router.HandleFunc("/api/getContactUser", middleware.VerifyJWT(handler.GetContactUser)).Methods(http.MethodGet)
	router.HandleFunc("/api/contact/{id}", handler.GetContactDetail).Methods(http.MethodGet)
	router.HandleFunc("/api/deleteRequestAgency", middleware.VerifyJWT(handler.DeleteRequestAgency)).Methods(http.MethodDelete)
	router.HandleFunc("/api/getAllContact", handler.GetAllContact).Methods(http.MethodGet)

	router.HandleFunc("/api/createEnterprise", middleware.VerifyJWT(handler.CreateEnterprise)).Methods(http.MethodPost)
	router.HandleFunc("/api/getAllEnterprise", handler.GetAllEnterprise).Methods(http.MethodGet)
	router.HandleFunc("/api/getPinnedEnterprise", handler.GetPinnedEnterprise).Methods(http.MethodGet)
	router.HandleFunc("/api/setPinnedEnterprise", middleware.VerifyJWT(handler.SetPinnedEnterprise)).Methods(http.MethodPost)
	router.HandleFunc("/api/enterprise/{id}", handler.GetEnterpriseDetail).Methods(http.MethodGet)
	
	router.HandleFunc("/api/admin/getRequestAgency", middleware.VerifyJWT(handler.GetRequestAgency)).Methods(http.MethodGet)
	router.HandleFunc("/api/admin/getRequestDisableAccount", middleware.VerifyJWT(handler.GetRequestDisableAccount)).Methods(http.MethodGet)
	router.HandleFunc("/api/admin/responseRequestAgency", middleware.VerifyJWT(handler.ResponseRequestAgency)).Methods(http.MethodPost)
	router.HandleFunc("/api/admin/deleteAccount", middleware.VerifyJWT(handler.AdminDeleteAccount)).Methods(http.MethodPost)
	router.HandleFunc("/api/admin/cancelDeleteAccount", middleware.VerifyJWT(handler.CancelDeleteAccount)).Methods(http.MethodPost)
	
	c := cors.New(cors.Options{
        AllowedOrigins: []string{"http://localhost:8080","http://localhost:3000","https://properties-by-vue.vercel.app","https://properties-by-knxkzi2mp-dannitran99s-projects.vercel.app"},
        AllowCredentials: true,
		AllowedMethods: []string{"POST", "GET", "DELETE"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
    })

    handler := c.Handler(router)

	log.Fatal(http.ListenAndServe(":5000", handler))
}