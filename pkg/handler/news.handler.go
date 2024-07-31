package handler

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"propertiesGo/pkg/dto"
	"propertiesGo/pkg/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllNews(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("content-type", "application/json")
	typeNews := request.URL.Query().Get("type")
	tags := request.URL.Query().Get("tags")
	keywordSearch := request.URL.Query().Get("k")
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
	filter := bson.D{}
	if typeNews != "" {
		typeFilter := bson.E{ Key:"category", Value: bson.M{
				"$in": strings.Split(typeNews, ","),
			},}
		filter = append(filter, typeFilter)
	}
	if keywordSearch != "" {
		keywordTitleFilter := bson.M{ "title": bson.M{"$regex": keywordSearch, "$options": "i"}}
		keywordTagsFilter := bson.M{ "tags": bson.M{"$regex": keywordSearch, "$options": "i"}}
		keywordFilter := bson.E{Key: "$or", Value: []bson.M{keywordTitleFilter,keywordTagsFilter}}
		filter = append(filter, keywordFilter)
	}
	if tags != "" {
		tagsFilter := bson.E{ Key:"tags", Value: tags}
		filter = append(filter, tagsFilter)
	}
	sortQuery := bson.D{{Key: "_id", Value: -1}}
	matchStage := bson.D{{Key: "$match", Value: filter}}
	sortStage := bson.D{{Key: "$sort", Value: sortQuery}}
	limitStage := bson.D{{Key: "$limit", Value: pageSize}}
	skipStage := bson.D{{Key: "$skip", Value: skip}}
	var new []dto.News
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	count, err := utils.MongoConnect("News").CountDocuments(ctx, filter)
	if err != nil {
		utils.StatusInternalServerError(writer)
	}
	cursor, err := utils.MongoConnect("News").Aggregate(ctx, mongo.Pipeline{matchStage, sortStage, skipStage, limitStage})
	if err != nil {
		utils.StatusInternalServerError(writer)
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var newItem dto.News
		cursor.Decode(&newItem)
		new = append(new, newItem)
	}
	if err := cursor.Err(); err != nil {
		utils.StatusInternalServerError(writer)
		return
	}
	responseData := dto.ResponseData{
		Data:  new,
		Total: count,
	}
	json.NewEncoder(writer).Encode(responseData)
}

func PostNews(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("content-type", "application/json")
	role := request.Context().Value("role")
	if role != "admin" {
        utils.StatusForbidden(writer)
		return
    }
	username := request.Context().Value("username")
	body, err := ioutil.ReadAll(request.Body)
    if err != nil {
		utils.StatusBadRequest(writer) 
		return
    }
    defer request.Body.Close()
	var post dto.News
    err = json.Unmarshal(body, &post)
    if err != nil {
		utils.StatusBadRequest(writer) 
		return
    }
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
	doc := bson.D{
        primitive.E{Key: "category", Value: post.Category},
        primitive.E{Key: "tags", Value: post.Tags},
		primitive.E{Key: "title", Value: post.Title},
		primitive.E{Key: "description", Value: post.Description},
        primitive.E{Key: "thumbnail", Value: post.Thumbnail},
        primitive.E{Key: "source", Value: post.Source},
        primitive.E{Key: "content", Value: post.Content},
        primitive.E{Key: "user", Value: username},
		primitive.E{Key: "createdAt", Value: post.CreatedAt},
    }
    result, _ := utils.MongoConnect("News").InsertOne(ctx, doc)
	json.NewEncoder(writer).Encode(result)
}

func GetNewsByID(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("content-type", "application/json")
	id, _ := mux.Vars(request)["id"]
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.StatusBadRequest(writer)
		return
	}
	var news dto.NewsDetail
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := utils.MongoConnect("News")
	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: objID}}}}, 
		{
			{Key: "$lookup", Value: bson.M{
				"from":         "Users",
				"localField":   "user",
				"foreignField": "username",
				"as":           "relatedUser",
			}},
		},
	})
	if err != nil {
		utils.StatusNotFound(writer)
		return
	}
	if cursor.Next(ctx) {
		cursor.Decode(&news)
	}
	json.NewEncoder(writer).Encode(news)
}