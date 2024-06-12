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
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func RegisterAgency(writer http.ResponseWriter, request *http.Request) {
	username := request.Context().Value("username")
	body, err := ioutil.ReadAll(request.Body)
    if err != nil {
        http.Error(writer, "Lỗi đọc nội dung request body", http.StatusBadRequest)
        return
    }
    defer request.Body.Close()
    var contact dto.Contacts
    err = json.Unmarshal(body, &contact)
    if err != nil {
        http.Error(writer, "Lỗi giải mã nội dung request body", http.StatusBadRequest)
        return
    }
	var contactDb dto.Contacts
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    collection := utils.MongoConnect("Contacts")
	err = collection.FindOne(ctx, bson.D{{Key: "username", Value: username}}).Decode(&contactDb)

	if err == nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "Không thể tạo thêm request" }`))
		return
	}
	doc := bson.D{
        primitive.E{Key: "username", Value: username}, 
        primitive.E{Key: "name", Value: contact.Name},
        primitive.E{Key: "avatar", Value: contact.Avatar},
        primitive.E{Key: "phoneNumber", Value: contact.PhoneNumber},
        primitive.E{Key: "status", Value: "pending"},
		primitive.E{Key: "createdAt", Value: contact.CreatedAt},
    }
    result, _ := collection.InsertOne(ctx, doc)
	json.NewEncoder(writer).Encode(result)
}

func GetContactUser(writer http.ResponseWriter, request *http.Request) {
	username := request.Context().Value("username")
	var contactDb dto.Contacts
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    collection := utils.MongoConnect("Contacts")
	collection.FindOne(ctx, bson.D{{Key: "username", Value: username}}).Decode(&contactDb)
	json.NewEncoder(writer).Encode(contactDb)
}