package handler

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"propertiesGo/pkg/dto"
	"propertiesGo/pkg/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

const contentType = "content-type"
const accepted = "application/json"

func GetRequestAgency(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set(contentType, accepted)
	role := request.Context().Value("role")
	if role != "admin" {
        utils.StatusForbidden(writer)
		return
    }
	var contacts []dto.Contacts
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := utils.MongoConnect("Contacts").Find(ctx, bson.D{{Key: "status", Value: "pending"}})
	if err != nil {
		utils.StatusInternalServerError(writer)
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var contact dto.Contacts
		cursor.Decode(&contact)
		contacts = append(contacts, contact)
	}
	if err := cursor.Err(); err != nil {
		utils.StatusInternalServerError(writer)
		return
	}
	json.NewEncoder(writer).Encode(contacts)
}

func GetRequestDisableAccount(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set(contentType, accepted)
	role := request.Context().Value("role")
	if role != "admin" {
        utils.StatusForbidden(writer)
		return
    }
	var users []dto.UserGet
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := utils.MongoConnect("Users").Find(ctx, bson.D{{Key: "status", Value: "delete-pending"}})
	if err != nil {
		utils.StatusInternalServerError(writer)
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var user dto.UserGet
		cursor.Decode(&user)
		users = append(users, user)
	}
	if err := cursor.Err(); err != nil {
		utils.StatusInternalServerError(writer)
		return
	}
	json.NewEncoder(writer).Encode(users)
}

func ResponseRequestAgency(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set(contentType, accepted)
	role := request.Context().Value("role")
	if role != "admin" {
        utils.StatusForbidden(writer)
		return
    }
	body, err := ioutil.ReadAll(request.Body)
    if err != nil {
		utils.StatusBadRequest(writer) 
		return
    }
    defer request.Body.Close()
	var userId dto.ID
    err = json.Unmarshal(body, &userId)
    if err != nil {
		utils.StatusBadRequest(writer)
		return 
    }
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
	if userId.Action == "active" {
		collectionUser := utils.MongoConnect("Users")
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "role", Value: userId.Role}}}}
		_, err = collectionUser.UpdateOne(ctx, bson.D{{Key: "username", Value: userId.Username}}, update)
		if err != nil {
			utils.StatusInternalServerError(writer)
			return
		}
	}
	collectionContact := utils.MongoConnect("Contacts")
	updateContact := bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: userId.Action}}}}
    results, err := collectionContact.UpdateOne(ctx, bson.D{{Key: "username", Value: userId.Username}}, updateContact)
    if err != nil {
		utils.StatusInternalServerError(writer)
		return
    }
	json.NewEncoder(writer).Encode(results)
}

func AdminDeleteAccount(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set(contentType, accepted)
	role := request.Context().Value("role")
	if role != "admin" {
        utils.StatusForbidden(writer)
		return
    }
	body, err := ioutil.ReadAll(request.Body)
    if err != nil {
		utils.StatusBadRequest(writer)
		return 
    }
    defer request.Body.Close()
	var userId dto.ID
    err = json.Unmarshal(body, &userId)
    if err != nil {
		utils.StatusBadRequest(writer) 
		return
    }
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
	collection := utils.MongoConnect("Users")
    deleteResult, _ := collection.DeleteOne(ctx, bson.D{{Key: "username", Value: userId.Username}})
    if deleteResult.DeletedCount == 0 {
		utils.StatusInternalServerError(writer)
		return
    }
    json.NewEncoder(writer).Encode(deleteResult)
}

func CancelDeleteAccount(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set(contentType, accepted)
	role := request.Context().Value("role")
	if role != "admin" {
        utils.StatusForbidden(writer)
		return
    }
	body, err := ioutil.ReadAll(request.Body)
    if err != nil {
		utils.StatusBadRequest(writer) 
		return
    }
    defer request.Body.Close()
	var userId dto.ID
    err = json.Unmarshal(body, &userId)
    if err != nil {
		utils.StatusBadRequest(writer) 
		return
    }
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
	collection := utils.MongoConnect("Users")
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: "active"}}}}
    result, err := collection.UpdateOne(ctx, bson.D{{Key: "username", Value: userId.Username}}, update)
    if err != nil {
		utils.StatusInternalServerError(writer)
		return
    }
    json.NewEncoder(writer).Encode(result)
}