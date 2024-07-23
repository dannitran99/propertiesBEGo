package handler

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"propertiesGo/pkg/dto"
	"propertiesGo/pkg/utils"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllNews(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("content-type", "application/json")
	var new []dto.News
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := utils.MongoConnect("News").Find(ctx, bson.M{})
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
	json.NewEncoder(writer).Encode(new)
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