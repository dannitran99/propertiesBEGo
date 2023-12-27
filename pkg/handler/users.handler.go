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

func Login(writer http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
    if err != nil {
        http.Error(writer, "Lỗi đọc nội dung request body", http.StatusBadRequest)
        return
    }
    defer request.Body.Close()
    // Giải mã nội dung của request body thành một struct User
    var user dto.User
    err = json.Unmarshal(body, &user)
    if err != nil {
        http.Error(writer, "Lỗi giải mã nội dung request body", http.StatusBadRequest)
        return
    }
    var userDb dto.User
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    collection := utils.MongoConnect("Users")
	err = collection.FindOne(ctx, bson.D{{Key: "username", Value: user.Username}}).Decode(&userDb)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "Tài khoản không tồn tại" }`))
		return
	}
    if utils.SHA1(user.Password) != userDb.Password {
        writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "Sai mật khẩu" }`))
		return
    }
    token, err := utils.CreateToken(user.Username)
    if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "Lỗi tạo token" }`))
		return
	}
    var reponse dto.LoginResponse
    reponse.Token = token
    reponse.Username = user.Username
	json.NewEncoder(writer).Encode(reponse)
}

func Register(writer http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
    if err != nil {
        http.Error(writer, "Lỗi đọc nội dung request body", http.StatusBadRequest)
        return
    }
    defer request.Body.Close()
    var user dto.User
    err = json.Unmarshal(body, &user)
    if err != nil {
        http.Error(writer, "Lỗi giải mã nội dung request body", http.StatusBadRequest)
        return
    }
	 var userDb dto.User
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    collection := utils.MongoConnect("Users")
	err = collection.FindOne(ctx, bson.D{{Key: "username", Value: user.Username}}).Decode(&userDb)

	if err == nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "Tài khoản đã tồn tại" }`))
		return
	}

    err = collection.FindOne(ctx, bson.D{{Key: "email", Value: user.Email}}).Decode(&userDb)
	if err == nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "Email đã tồn tại" }`))
		return
	}
    doc := bson.D{
        primitive.E{Key: "username", Value: user.Username}, 
        primitive.E{Key: "password", Value: utils.SHA1(user.Password)},
        primitive.E{Key: "email", Value: user.Email},
    }
    result, _ := collection.InsertOne(ctx, doc)
	json.NewEncoder(writer).Encode(result)
}