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

func GetAllProperties(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("content-type", "application/json")
	var output []dto.Properties
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := utils.MongoConnect("Properties").Find(ctx, bson.M{})
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var property dto.Properties
		cursor.Decode(&property)
		output = append(output, property)
	}
	if err := cursor.Err(); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(writer).Encode(output)
}

func PostProperties(writer http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
    if err != nil {
        http.Error(writer, "Lỗi đọc nội dung request body", http.StatusBadRequest)
        return
    }
    defer request.Body.Close()
	var post dto.Properties
    err = json.Unmarshal(body, &post)
    if err != nil {
        http.Error(writer, "Lỗi giải mã nội dung request body", http.StatusBadRequest)
        return
    }
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    collection := utils.MongoConnect("Properties")
	doc := bson.D{
        primitive.E{Key: "type", Value: post.Type},
		primitive.E{Key: "propertyType", Value: post.PropertyType},
        primitive.E{Key: "city", Value: post.City},
		primitive.E{Key: "district", Value: post.District},
        primitive.E{Key: "ward", Value: post.Ward},
		primitive.E{Key: "street", Value: post.Street},
        primitive.E{Key: "project", Value: post.Project},
		primitive.E{Key: "moreInfo", Value: post.MoreInfo},
        primitive.E{Key: "title", Value: post.Title},
		primitive.E{Key: "description", Value: post.Description},
        primitive.E{Key: "area", Value: post.Area},
		primitive.E{Key: "price", Value: post.Price},
        primitive.E{Key: "priceType", Value: post.PriceType},
		primitive.E{Key: "images", Value: post.Images},
        primitive.E{Key: "name", Value: post.Name},
		primitive.E{Key: "phoneNumber", Value: post.PhoneNumber},
		primitive.E{Key: "email", Value: post.Email},
        primitive.E{Key: "user", Value: post.User},
		primitive.E{Key: "createdAt", Value: post.CreatedAt},
    }
    result, _ := collection.InsertOne(ctx, doc)
	json.NewEncoder(writer).Encode(result)
}