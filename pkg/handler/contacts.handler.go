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
        primitive.E{Key: "type", Value: "ca-nhan"}, 
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
	err := collection.FindOne(ctx, bson.D{{Key: "username", Value: username}}).Decode(&contactDb)
	if err != nil {
		writer.Write([]byte(`{ "message": "no data" }`))
		return
	}
	json.NewEncoder(writer).Encode(contactDb)
}

func DeleteRequestAgency(writer http.ResponseWriter, request *http.Request) {
	username := request.Context().Value("username")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
	collection := utils.MongoConnect("Contacts")
    deleteResult, _ := collection.DeleteOne(ctx, bson.D{{Key: "username", Value: username}})
    if deleteResult.DeletedCount == 0 {
        writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "Xóa không thành công" }`))
		return
    }
    json.NewEncoder(writer).Encode(deleteResult)
}

func GetAllContact(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("content-type", "application/json")
	typeContact := request.URL.Query().Get("contactType")
	var output []dto.Contacts
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{Key: "type", Value: typeContact}}
	filter = append(filter, bson.E{ Key: "status", Value: "active" })
	cursor, err := utils.MongoConnect("Contacts").Find(ctx, filter)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var contact dto.Contacts
		cursor.Decode(&contact)
		output = append(output, contact)
	}
	if err := cursor.Err(); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(writer).Encode(output)
}