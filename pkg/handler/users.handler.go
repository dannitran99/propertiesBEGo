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
	json.NewEncoder(writer).Encode(token)
}