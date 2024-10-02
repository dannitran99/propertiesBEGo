package handler

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"propertiesGo/pkg/dto"
	"propertiesGo/pkg/utils"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
		primitive.E{Key: "pinned", Value: ""},
		primitive.E{Key: "createdAt", Value: post.CreatedAt},
    }
    result, _ := utils.MongoConnect("Enterprises").InsertOne(ctx, doc)
	json.NewEncoder(writer).Encode(result)
}

func GetAllEnterprise(writer http.ResponseWriter, request *http.Request){
	writer.Header().Set("content-type", "application/json")
	keywordSearch := request.URL.Query().Get("k")
	filterType := request.URL.Query().Get("type")
	filterCity := request.URL.Query().Get("city")
	filterDistrict := request.URL.Query().Get("district")
	pageQuery := request.URL.Query().Get("p")
	limitQuery := request.URL.Query().Get("l")
	page, err := strconv.Atoi(pageQuery)
	if err != nil {
		utils.StatusBadRequest(writer) 
		return
	}
	pageSize, err := strconv.Atoi(limitQuery)
	if err != nil {
		utils.StatusBadRequest(writer) 
		return
	}
	skip := (page - 1) * pageSize
	var output []dto.Enterprise
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{}
	if keywordSearch != "" {
		keywordFilter := bson.E{ Key:"name", Value: bson.M{"$regex": keywordSearch, "$options": "i"}}
		filter = append(filter, keywordFilter)
	}
	if filterType != "" {
		typeFilter := bson.E{Key: "businessField", Value: filterType}
		filter = append(filter, typeFilter)
	}
	if filterCity != "" {
		cityFilter := bson.E{Key: "city", Value: filterCity}
		filter = append(filter, cityFilter)
	}
	if filterDistrict != "" {
		districtFilter := bson.E{Key: "district", Value: filterDistrict}
		filter = append(filter, districtFilter)
	}
	sortQuery := bson.D{{Key: "_id", Value: -1}}
	count, err := utils.MongoConnect("Enterprises").CountDocuments(ctx, filter)
	if err != nil {
		utils.StatusInternalServerError(writer)
	}
	matchStage := bson.D{{Key: "$match", Value: filter}}
	sortStage := bson.D{{Key: "$sort", Value: sortQuery}}
	limitStage := bson.D{{Key: "$limit", Value: pageSize}}
	skipStage := bson.D{{Key: "$skip", Value: skip}}
	cursor, err := utils.MongoConnect("Enterprises").Aggregate(ctx, mongo.Pipeline{matchStage, sortStage, skipStage, limitStage})
	if err != nil {
		utils.StatusInternalServerError(writer)
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var enterprise dto.Enterprise
		cursor.Decode(&enterprise)
		output = append(output, enterprise)
	}
	if err := cursor.Err(); err != nil {
		utils.StatusInternalServerError(writer)
		return
	}
	
	responseData := dto.ResponseData{
		Data:  output,
		Total: count,
	}
	json.NewEncoder(writer).Encode(responseData)
}

func SetPinnedEnterprise(writer http.ResponseWriter, request *http.Request){
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
	var post dto.PinnedEnterprise
    err = json.Unmarshal(body, &post)
    if err != nil {
		utils.StatusBadRequest(writer) 
		return
    }
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
	collection := utils.MongoConnect("Enterprises")
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "pinned", Value: post.Pinned},
		}},
	}
    result, err := collection.UpdateOne(ctx, bson.D{{Key: "_id", Value: post.ID}}, update)
    if err != nil {
		utils.StatusInternalServerError(writer)
		return
    }
    json.NewEncoder(writer).Encode(result)
}