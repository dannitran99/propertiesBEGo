package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"propertiesGo/pkg/dto"
	"propertiesGo/pkg/utils"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func GetNewsByID(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("content-type", "application/json")
	var news dto.News
	id, _ := mux.Vars(request)["id"]
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.StatusBadRequest(writer)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := utils.MongoConnect("News")
	err = collection.FindOne(ctx, bson.D{{Key: "_id", Value: objID}}).Decode(&news)

	if err != nil {
		utils.StatusNotFound(writer)
		return
	}
	json.NewEncoder(writer).Encode(news)
}