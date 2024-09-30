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

func CreateEnterprise(writer http.ResponseWriter, request *http.Request){
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
	var post dto.Enterprise
    err = json.Unmarshal(body, &post)
    if err != nil {
		utils.StatusBadRequest(writer) 
		return
    }
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
	doc := bson.D{
        primitive.E{Key: "logo", Value: post.Logo},
        primitive.E{Key: "banner", Value: post.Banner},
		primitive.E{Key: "name", Value: post.Name},
		primitive.E{Key: "city", Value: post.City},
        primitive.E{Key: "district", Value: post.District},
        primitive.E{Key: "ward", Value: post.Ward},
        primitive.E{Key: "street", Value: post.Street},
        primitive.E{Key: "businessField", Value: post.BusinessField},
		primitive.E{Key: "subBusiness", Value: post.SubBusiness},
        primitive.E{Key: "description", Value: post.Description},
        primitive.E{Key: "phoneNumber", Value: post.PhoneNumber},
        primitive.E{Key: "email", Value: post.Email},
		primitive.E{Key: "website", Value: post.Website},
		primitive.E{Key: "createdAt", Value: post.CreatedAt},
    }
    result, _ := utils.MongoConnect("Enterprises").InsertOne(ctx, doc)
	json.NewEncoder(writer).Encode(result)
}