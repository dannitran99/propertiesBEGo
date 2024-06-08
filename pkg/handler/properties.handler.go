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
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ResponseData struct {
	Data  interface{}
	Total int64
}

func GetAllProperties(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("content-type", "application/json")
	typeProperties := request.URL.Query().Get("type")
	categoryProperties := request.URL.Query().Get("category")
	keywordProperties := request.URL.Query().Get("k")
	cityQuery := request.URL.Query().Get("city")
	districtQuery := request.URL.Query().Get("district")
	minPriceQuery := request.URL.Query().Get("minp")
	maxPriceQuery := request.URL.Query().Get("maxp")
	minSquareQuery := request.URL.Query().Get("mins")
	maxSquareQuery := request.URL.Query().Get("maxs")
	pageQuery := request.URL.Query().Get("p")
	limitQuery := request.URL.Query().Get("l")
	page, err := strconv.Atoi(pageQuery)
	if err != nil {
		panic(err) 
	}
	pageSize, err := strconv.Atoi(limitQuery)
	if err != nil {
		panic(err) 
	}
	skip := (page - 1) * pageSize
	var output []dto.PropertiesInfo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{Key: "type", Value: typeProperties}}
	if categoryProperties != "" {
		categoryfilter := bson.E{
			Key: "propertyType",
			Value: bson.M{
				"$in": strings.Split(categoryProperties, ","),
			},
		}
		filter = append(filter, categoryfilter)
	}
	if keywordProperties != "" {
		keywordFilter := bson.E{ Key:"title", Value: bson.M{"$regex": keywordProperties, "$options": "i"}}
		filter = append(filter, keywordFilter)
	}
	if cityQuery != "" {
		cityFilter := bson.E{Key: "city", Value: cityQuery}
		filter = append(filter, cityFilter)
	}
	if districtQuery != "" {
		districtfilter := bson.E{
			Key: "district",
			Value: bson.M{
				"$in": strings.Split(districtQuery, ","),
			},
		}
		filter = append(filter, districtfilter)
	}
	if minPriceQuery != "" {
		i, err := strconv.Atoi(minPriceQuery)
		if err != nil {
			panic(err) // Handle error
		}
		minPriceFilter := bson.E{Key: "price",Value: bson.M{"$gte": i*1000000}}
		filter = append(filter, minPriceFilter)
	}
	if maxPriceQuery != "" {
		i, err := strconv.Atoi(maxPriceQuery)
		if err != nil {
			panic(err) // Handle error
		}
		maxPriceFilter := bson.E{Key: "price",Value: bson.M{"$lte": i*1000000}}
		filter = append(filter, maxPriceFilter)
	}
	if minSquareQuery != "" {
		i, err := strconv.Atoi(minSquareQuery)
		if err != nil {
			panic(err) // Handle error
		}
		minSquareFilter := bson.E{Key: "area",Value: bson.M{"$gte": i}}
		filter = append(filter, minSquareFilter)
	}
	if maxSquareQuery != "" {
		i, err := strconv.Atoi(maxSquareQuery)
		if err != nil {
			panic(err) // Handle error
		}
		maxSquareFilter := bson.E{Key: "area",Value: bson.M{"$lte": i}}
		filter = append(filter, maxSquareFilter)
	}
	findOptions := options.Find().SetSort(bson.D{{Key: "_id", Value: -1}}).SetSkip(int64(skip)).SetLimit(int64(pageSize))
	count, err := utils.MongoConnect("Properties").CountDocuments(ctx, filter)
	if err != nil {
		panic(err)
	}
	cursor, err := utils.MongoConnect("Properties").Find(ctx, filter, findOptions)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var property dto.PropertiesInfo
		var user dto.User
		cursor.Decode(&property)
    	collection := utils.MongoConnect("Users")
		_ = collection.FindOne(ctx, bson.D{{Key: "username", Value: property.User}}).Decode(&user)
		property.Avatar = user.Avatar
		output = append(output, property)
	}
	if err := cursor.Err(); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	responseData := ResponseData{
		Data:  output,
		Total: count,
	}
	json.NewEncoder(writer).Encode(responseData)
}

func GetAllPropertiesHome(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("content-type", "application/json")
	pageQuery := request.URL.Query().Get("p")
	limitQuery := request.URL.Query().Get("l")
	page, err := strconv.Atoi(pageQuery)
	if err != nil {
		panic(err) 
	}
	pageSize, err := strconv.Atoi(limitQuery)
	if err != nil {
		panic(err) 
	}
	skip := (page - 1) * pageSize
	var output []dto.Properties
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	findOptions := options.Find().SetSort(bson.D{{Key: "_id", Value: -1}}).SetSkip(int64(skip)).SetLimit(int64(pageSize))
	if err != nil {
		panic(err)
	}
	cursor, err := utils.MongoConnect("Properties").Find(ctx, bson.D{}, findOptions)
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
		primitive.E{Key: "url", Value: post.Url},
		primitive.E{Key: "email", Value: post.Email},
        primitive.E{Key: "user", Value: post.User},
		primitive.E{Key: "createdAt", Value: post.CreatedAt},
    }
    result, _ := collection.InsertOne(ctx, doc)
	json.NewEncoder(writer).Encode(result)
}

func GetPostedProperty(writer http.ResponseWriter, request *http.Request) {
	userName := request.URL.Query().Get("name")
	var output []dto.PropertiesInfo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := utils.MongoConnect("Properties").Find(ctx, bson.D{{Key: "user", Value: userName}})
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var property dto.PropertiesInfo
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

func GetPropertiesDetail(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("content-type", "application/json")
	id, _ := mux.Vars(request)["id"]
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "Nguồn không tồn tại" }`))
		return
	}
	var property dto.PropertiesInfo
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    collection := utils.MongoConnect("Properties")
	err = collection.FindOne(ctx, bson.D{{Key: "_id", Value: objID}}).Decode(&property)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "Bài đăng không tồn tại" }`))
		return
	}
	var user dto.User
    collectionUser := utils.MongoConnect("Users")
	_ = collectionUser.FindOne(ctx, bson.D{{Key: "username", Value: property.User}}).Decode(&user)
	property.Avatar = user.Avatar
	json.NewEncoder(writer).Encode(property)
}
	