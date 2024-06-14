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

func GetRequestAgency(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("content-type", "application/json")
	role := request.Context().Value("role")
	if role != "admin" {
        writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "Không có quyền truy cập" }`))
		return
    }
	var contacts []dto.Contacts
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := utils.MongoConnect("Contacts").Find(ctx, bson.D{{Key: "status", Value: "pending"}})
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var contact dto.Contacts
		cursor.Decode(&contact)
		contacts = append(contacts, contact)
	}
	if err := cursor.Err(); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(writer).Encode(contacts)
}

func GetRequestDisableAccount(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("content-type", "application/json")
	role := request.Context().Value("role")
	if role != "admin" {
        writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "Không có quyền truy cập" }`))
		return
    }
	var users []dto.UserInfo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := utils.MongoConnect("Users").Find(ctx, bson.D{{Key: "status", Value: "delete-pending"}})
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var user dto.UserInfo
		cursor.Decode(&user)
		users = append(users, user)
	}
	if err := cursor.Err(); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(writer).Encode(users)
}

func AcceptRequestAgency(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("content-type", "application/json")
	username := request.Context().Value("username")
	role := request.Context().Value("role")
	if role != "admin" {
        writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "Không có quyền truy cập" }`))
		return
    }
	body, err := ioutil.ReadAll(request.Body)
    if err != nil {
        http.Error(writer, "Lỗi đọc nội dung request body", http.StatusBadRequest)
        return
    }
    defer request.Body.Close()
	var userId dto.ID
    err = json.Unmarshal(body, &userId)
    if err != nil {
        http.Error(writer, "Lỗi giải mã nội dung request body", http.StatusBadRequest)
        return
    }
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    collectionUser := utils.MongoConnect("Users")
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "role", Value: "agency"}}}}
    _, err = collectionUser.UpdateOne(ctx, bson.D{{Key: "username", Value: username}}, update)
    if err != nil {
        writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "Đổi role không thành công" }`))
		return
    }
	collectionContact := utils.MongoConnect("Contacts")
	updateContact := bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: "active"}}}}
    results, err := collectionContact.UpdateOne(ctx, bson.D{{Key: "username", Value: username}}, updateContact)
    if err != nil {
        writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "Đổi role không thành công" }`))
		return
    }
	json.NewEncoder(writer).Encode(results)
}