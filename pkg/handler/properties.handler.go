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
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
	filterVerify := request.URL.Query().Get("f")
	sort := request.URL.Query().Get("sort")
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
			utils.StatusBadRequest(writer) 
			return
		}
		minPriceFilter := bson.E{Key: "price",Value: bson.M{"$gte": i*1000000}}
		filter = append(filter, minPriceFilter)
	}
	if maxPriceQuery != "" {
		i, err := strconv.Atoi(maxPriceQuery)
		if err != nil {
			utils.StatusBadRequest(writer) 
			return
		}
		maxPriceFilter := bson.E{Key: "price",Value: bson.M{"$lte": i*1000000}}
		filter = append(filter, maxPriceFilter)
	}
	if minSquareQuery != "" {
		i, err := strconv.Atoi(minSquareQuery)
		if err != nil {
			utils.StatusBadRequest(writer) 
			return
		}
		minSquareFilter := bson.E{Key: "area",Value: bson.M{"$gte": i}}
		filter = append(filter, minSquareFilter)
	}
	if maxSquareQuery != "" {
		i, err := strconv.Atoi(maxSquareQuery)
		if err != nil {
			utils.StatusBadRequest(writer) 
			return
		}
		maxSquareFilter := bson.E{Key: "area",Value: bson.M{"$lte": i}}
		filter = append(filter, maxSquareFilter)
	}
	if filterVerify == "agency" {
		cursor, err := utils.MongoConnect("Contacts").Find(ctx, bson.D{{ Key: "status", Value: "active" }})
		if err != nil {
			utils.StatusInternalServerError(writer)
			return
		}
		defer cursor.Close(ctx)
		var results []bson.M
		var user []string
		if err := cursor.All(context.TODO(), &results); err != nil {
			utils.StatusInternalServerError(writer)
			return
		}
		for _, result := range results {
			user = append(user, result["username"].(string))
		}
		agencyfilter := bson.E{
			Key: "user",
			Value: bson.M{
				"$in": user,
			},
		}
		filter = append(filter, agencyfilter)
	}
	sortQuery := bson.D{{Key: "_id", Value: -1}}
	if sort != "" {
		sortType, err := strconv.Atoi(sort)
		if err != nil {
			utils.StatusBadRequest(writer) 
			return
		}
		switch sortType {
            case 1:
				hasPrice := bson.E{Key: "price",Value: bson.M{"$ne": 0}}
				filter = append(filter, hasPrice)
				sortQuery = bson.D{{Key: "price", Value: 1}}
            case 2:
				hasPrice := bson.E{Key: "price",Value: bson.M{"$ne": 0}}
				filter = append(filter, hasPrice)
				sortQuery = bson.D{{Key: "price", Value: -1}}
			case 3:
				hasPriceAvg := bson.E{Key: "priceAvg",Value: bson.M{"$ne": 0}}
				filter = append(filter, hasPriceAvg)
				sortQuery = bson.D{{Key: "priceAvg", Value: 1}}
            case 4:
				hasPriceAvg := bson.E{Key: "priceAvg",Value: bson.M{"$ne": 0}}
				filter = append(filter, hasPriceAvg)
				sortQuery = bson.D{{Key: "priceAvg", Value: -1}}
            case 5:
				sortQuery = bson.D{{Key: "area", Value: 1}}
            case 6:
				sortQuery = bson.D{{Key: "area", Value: -1}}
        }
	}
	count, err := utils.MongoConnect("Properties").CountDocuments(ctx, filter)
	if err != nil {
		utils.StatusInternalServerError(writer)
		return
	}
	matchStage := bson.D{{Key: "$match", Value: filter}}
	sortStage := bson.D{{Key: "$sort", Value: sortQuery}}
	limitStage := bson.D{{Key: "$limit", Value: pageSize}}
	skipStage := bson.D{{Key: "$skip", Value: skip}}
	cursor, err := utils.MongoConnect("Properties").Aggregate(ctx, mongo.Pipeline{matchStage, sortStage, skipStage, limitStage , 
		bson.D{
			{Key: "$lookup", Value: bson.M{
				"from":         "Users",
				"localField":   "user",
				"foreignField": "username",
				"as":           "relatedUser",
			}},
		},
	})
	if err != nil {
		utils.StatusInternalServerError(writer)
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var property dto.PropertiesInfo
		cursor.Decode(&property)
		output = append(output, property)
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

func GetAllPropertiesHome(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("content-type", "application/json")
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
	var output []dto.Properties
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	findOptions := options.Find().SetSort(bson.D{{Key: "_id", Value: -1}}).SetSkip(int64(skip)).SetLimit(int64(pageSize))
	if err != nil {
		utils.StatusInternalServerError(writer)
		return
	}
	cursor, err := utils.MongoConnect("Properties").Find(ctx, bson.D{}, findOptions)
	if err != nil {
		utils.StatusInternalServerError(writer)
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var property dto.Properties
		cursor.Decode(&property)
		output = append(output, property)
	}
	if err := cursor.Err(); err != nil {
		utils.StatusInternalServerError(writer)
		return
	}
	json.NewEncoder(writer).Encode(output)
}

func PostProperties(writer http.ResponseWriter, request *http.Request) {
	username := request.Context().Value("username")
	body, err := ioutil.ReadAll(request.Body)
    if err != nil {
		utils.StatusBadRequest(writer) 
		return
    }
    defer request.Body.Close()
	var post dto.Properties
    err = json.Unmarshal(body, &post)
    if err != nil {
		utils.StatusBadRequest(writer) 
		return
    }
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    collection := utils.MongoConnect("Properties")
	var price int64
	var priceAvg float32
	switch post.PriceType {
		case "VND":
			price = post.Price
			priceAvg = float32(post.Price) / float32(post.Area)
		case "Giá / m²":
			price = int64(post.Price) * int64(post.Area)
			priceAvg = float32(post.Price)
	}
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
		primitive.E{Key: "price", Value: price},
        primitive.E{Key: "priceType", Value: post.PriceType},
        primitive.E{Key: "priceAvg", Value: priceAvg},
		primitive.E{Key: "images", Value: post.Images},
        primitive.E{Key: "name", Value: post.Name},
		primitive.E{Key: "phoneNumber", Value: post.PhoneNumber},
		primitive.E{Key: "url", Value: post.Url},
		primitive.E{Key: "email", Value: post.Email},
        primitive.E{Key: "user", Value: username},
		primitive.E{Key: "createdAt", Value: post.CreatedAt},
    }
    result, _ := collection.InsertOne(ctx, doc)
	json.NewEncoder(writer).Encode(result)
}

func GetPostedProperty(writer http.ResponseWriter, request *http.Request) {
	username := request.Context().Value("username")
	var output []dto.PropertiesInfo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := utils.MongoConnect("Properties").Find(ctx, bson.D{{Key: "user", Value: username}})
	if err != nil {
		utils.StatusInternalServerError(writer)
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var property dto.PropertiesInfo
		cursor.Decode(&property)
		output = append(output, property)
	}
	if err := cursor.Err(); err != nil {
		utils.StatusInternalServerError(writer)
		return
	}
	json.NewEncoder(writer).Encode(output)
}

func GetPropertiesDetail(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("content-type", "application/json")
	id := mux.Vars(request)["id"]
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.StatusNotFound(writer)
		return
	}
	var property dto.PropertiesInfo
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    collection := utils.MongoConnect("Properties")
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
		cursor.Decode(&property)
	}
	json.NewEncoder(writer).Encode(property)
}
	