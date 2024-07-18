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
	"go.mongodb.org/mongo-driver/mongo"
)


func CheckVerifyToken(writer http.ResponseWriter, request *http.Request) {
	username := request.Context().Value("username").(string)
	role := request.Context().Value("role").(string)
    token, err := utils.CreateToken(username,role)
    if err != nil {
		utils.StatusInternalServerError(writer)
		return
	}
	json.NewEncoder(writer).Encode(token)
}

func Login(writer http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
    if err != nil {
		utils.StatusBadRequest(writer) 
		return
    }
    defer request.Body.Close()
    var user dto.User
    err = json.Unmarshal(body, &user)
    if err != nil {
		utils.StatusBadRequest(writer) 
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
    if userDb.Status == "disabled" || userDb.Status == "delete-pending"{
        writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "Tài khoản bị vô hiệu hóa" }`))
		return
    }
    token, err := utils.CreateToken(userDb.Username, userDb.Role)
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
		utils.StatusBadRequest(writer) 
		return
    }
    defer request.Body.Close()
    var user dto.User
    err = json.Unmarshal(body, &user)
    if err != nil {
		utils.StatusBadRequest(writer) 
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
        primitive.E{Key: "status", Value: "active"},
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
		utils.StatusBadRequest(writer) 
		return
    }
    defer request.Body.Close()
    var newPass dto.ChangePassword
    err = json.Unmarshal(body, &newPass)
    if err != nil {
		utils.StatusBadRequest(writer) 
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
		utils.StatusBadRequest(writer) 
		return
    }
    defer request.Body.Close()
    var newPass dto.ChangePassword
    err = json.Unmarshal(body, &newPass)
    if err != nil {
		utils.StatusBadRequest(writer) 
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
    update := bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: "disabled"}}}}
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
    update := bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: "delete-pending"}}}}
    result, err := collection.UpdateOne(ctx, bson.D{{Key: "_id", Value: userDb.ID}}, update)
    if err != nil {
        writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "Yêu cầu xóa không thành công" }`))
		return
    }
    json.NewEncoder(writer).Encode(result)
}

func ChangeAvatar(writer http.ResponseWriter, request *http.Request) {
	username := request.Context().Value("username")
    body, err := ioutil.ReadAll(request.Body)
    if err != nil {
		utils.StatusBadRequest(writer) 
		return
    }
    defer request.Body.Close()
    var avatar dto.ChangeAvatar
    err = json.Unmarshal(body, &avatar)
    if err != nil {
		utils.StatusBadRequest(writer) 
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
    var userDb []dto.UserInfo
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    cursor, err := utils.MongoConnect("Users").Aggregate(ctx, mongo.Pipeline{bson.D{{Key: "$match", Value: bson.D{{Key: "username", Value: username}}}}, 
		bson.D{
			{Key: "$lookup", Value: bson.M{
				"from":         "Contacts",
                "let": bson.M{ "username" : "$username" },
                "pipeline": []bson.M{
                    {
                        "$match": bson.M{
                            "$expr": bson.M{
                                "$and": []interface{}{
                                    bson.M{"$eq": []interface{}{"$username", "$$username"}},
                                    bson.M{"$eq": []interface{}{"$status", "active"}},
                                },
                            },
                        },
                    },
                },
				"as":           "agencyInfo",
			}},
		},
	})
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "Tài khoản không tồn tại" }`))
		return
	}
    defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var user dto.UserInfo
		cursor.Decode(&user)
		userDb = append(userDb, user)
	}
	if err := cursor.Err(); err != nil {
		utils.StatusInternalServerError(writer)
		return
	}
	json.NewEncoder(writer).Encode(userDb)
}

func ChangeInfo(writer http.ResponseWriter, request *http.Request) {
	username := request.Context().Value("username")
    body, err := ioutil.ReadAll(request.Body)
    if err != nil {
		utils.StatusBadRequest(writer) 
		return
    }
    defer request.Body.Close()
    var userPost dto.UserInfoUpdate
    err = json.Unmarshal(body, &userPost)
    if err != nil {
		utils.StatusBadRequest(writer) 
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
		utils.StatusInternalServerError(writer)
		return
    }
    json.NewEncoder(writer).Encode(result)
}