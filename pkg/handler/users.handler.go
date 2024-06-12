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


func CheckVerifyToken(writer http.ResponseWriter, request *http.Request) {
	username := request.Context().Value("username").(string)
    token, err := utils.CreateToken(username)
    if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "Lỗi tạo token" }`))
		return
	}
	json.NewEncoder(writer).Encode(token)
}

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
    if !userDb.Active {
        writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "Tài khoản bị vô hiệu hóa" }`))
		return
    }
    token, err := utils.CreateToken(userDb.Username)
    if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "Lỗi tạo token" }`))
		return
	}
    var reponse dto.LoginResponse
    reponse.Token = token
    reponse.Username = userDb.Username
    reponse.Avatar = userDb.Avatar
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
        primitive.E{Key: "active", Value: true},
        primitive.E{Key: "fullname", Value: ""},
        primitive.E{Key: "avatar", Value: ""},
        primitive.E{Key: "phoneNumber", Value: ""},
        primitive.E{Key: "role", Value: "user"},
    }
    result, _ := collection.InsertOne(ctx, doc)
	json.NewEncoder(writer).Encode(result)
}

func ChangePassword(writer http.ResponseWriter, request *http.Request) {
	username := request.Context().Value("username")
    body, err := ioutil.ReadAll(request.Body)
    if err != nil {
        http.Error(writer, "Lỗi đọc nội dung request body", http.StatusBadRequest)
        return
    }
    defer request.Body.Close()
    var newPass dto.ChangePassword
    err = json.Unmarshal(body, &newPass)
    if err != nil {
        http.Error(writer, "Lỗi giải mã nội dung request body", http.StatusBadRequest)
        return
    }
    var userDb dto.User
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    collection := utils.MongoConnect("Users")
	err = collection.FindOne(ctx, bson.D{{Key: "username", Value: username}}).Decode(&userDb)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "Tài khoản không tồn tại" }`))
		return
	}
    if utils.SHA1(newPass.CurrentPassword) != userDb.Password {
        writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "Mật khẩu cũ không đúng" }`))
		return
    }
    update := bson.D{{Key: "$set", Value: bson.D{{Key: "password", Value: utils.SHA1(newPass.NewPassword)}}}}
    result, err := collection.UpdateOne(ctx, bson.D{{Key: "_id", Value: userDb.ID}}, update)
    if err != nil {
        writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "Đổi không thành công" }`))
		return
    }
    json.NewEncoder(writer).Encode(result)
}

func DisableAccount(writer http.ResponseWriter, request *http.Request) {
	username := request.Context().Value("username")
    body, err := ioutil.ReadAll(request.Body)
    if err != nil {
        http.Error(writer, "Lỗi đọc nội dung request body", http.StatusBadRequest)
        return
    }
    defer request.Body.Close()
    var newPass dto.ChangePassword
    err = json.Unmarshal(body, &newPass)
    if err != nil {
        http.Error(writer, "Lỗi giải mã nội dung request body", http.StatusBadRequest)
        return
    }
    var userDb dto.User
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    collection := utils.MongoConnect("Users")
	err = collection.FindOne(ctx, bson.D{{Key: "username", Value: username}}).Decode(&userDb)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "Tài khoản không tồn tại" }`))
		return
	}
    if utils.SHA1(newPass.CurrentPassword) != userDb.Password {
        writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "Mật khẩu cũ không đúng" }`))
		return
    }
    update := bson.D{{Key: "$set", Value: bson.D{{Key: "active", Value: false}}}}
    result, err := collection.UpdateOne(ctx, bson.D{{Key: "_id", Value: userDb.ID}}, update)
    if err != nil {
        writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "Vô hiệu hóa không thành công" }`))
		return
    }
    json.NewEncoder(writer).Encode(result)
}

func DeleteAccount(writer http.ResponseWriter, request *http.Request) {
	username := request.Context().Value("username")
    var userDb dto.User
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    collection := utils.MongoConnect("Users")
	err := collection.FindOne(ctx, bson.D{{Key: "username", Value: username}}).Decode(&userDb)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "Tài khoản không tồn tại" }`))
		return
	}
    deleteResult, _ := collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: userDb.ID}})
    if deleteResult.DeletedCount == 0 {
        writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "Xóa không thành công" }`))
		return
    }
    json.NewEncoder(writer).Encode(deleteResult)
}

func ChangeAvatar(writer http.ResponseWriter, request *http.Request) {
	username := request.Context().Value("username")
    body, err := ioutil.ReadAll(request.Body)
    if err != nil {
        http.Error(writer, "Lỗi đọc nội dung request body", http.StatusBadRequest)
        return
    }
    defer request.Body.Close()
    var avatar dto.ChangeAvatar
    err = json.Unmarshal(body, &avatar)
    if err != nil {
        http.Error(writer, "Lỗi giải mã nội dung request body", http.StatusBadRequest)
        return
    }
    var userDb dto.User
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    collection := utils.MongoConnect("Users")
	err = collection.FindOne(ctx, bson.D{{Key: "username", Value: username}}).Decode(&userDb)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "Tài khoản không tồn tại" }`))
		return
	}
    update := bson.D{{Key: "$set", Value: bson.D{{Key: "avatar", Value: avatar.Avatar}}}}
    _, err = collection.UpdateOne(ctx, bson.D{{Key: "_id", Value: userDb.ID}}, update)
    if err != nil {
        writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "Update Avatar thất bại" }`))
		return
    }
    json.NewEncoder(writer).Encode(avatar)
}

func GetInfoUser(writer http.ResponseWriter, request *http.Request) {
	username := request.Context().Value("username")
    var userDb dto.User
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    collection := utils.MongoConnect("Users")
	err := collection.FindOne(ctx, bson.D{{Key: "username", Value: username}}).Decode(&userDb)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "Tài khoản không tồn tại" }`))
		return
	}
    var userInfo dto.UserInfo
    userInfo.Name = userDb.FullName
    userInfo.PhoneNumber = userDb.PhoneNumber
    userInfo.Email = userDb.Email
    
	json.NewEncoder(writer).Encode(userInfo)
}

func ChangeInfo(writer http.ResponseWriter, request *http.Request) {
	username := request.Context().Value("username")
    body, err := ioutil.ReadAll(request.Body)
    if err != nil {
        http.Error(writer, "Lỗi đọc nội dung request body", http.StatusBadRequest)
        return
    }
    defer request.Body.Close()
    var userPost dto.UserInfo
    err = json.Unmarshal(body, &userPost)
    if err != nil {
        http.Error(writer, "Lỗi giải mã nội dung request body", http.StatusBadRequest)
        return
    }
    var userDb dto.User
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    collection := utils.MongoConnect("Users")
	err = collection.FindOne(ctx, bson.D{{Key: "username", Value: username}}).Decode(&userDb)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "Tài khoản không tồn tại" }`))
		return
	}
    update := bson.D{{Key: "$set", Value: bson.D{{Key: "fullname", Value: userPost.Name},{Key: "phoneNumber", Value: userPost.PhoneNumber},{Key: "email", Value: userPost.Email}}}}
    result, err := collection.UpdateOne(ctx, bson.D{{Key: "_id", Value: userDb.ID}}, update)
    if err != nil {
        writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "Sửa thông tin không thành công" }`))
		return
    }
    json.NewEncoder(writer).Encode(result)
}