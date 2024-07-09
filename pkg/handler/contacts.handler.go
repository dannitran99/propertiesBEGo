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

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


func RegisterAgency(writer http.ResponseWriter, request *http.Request) {
	username := request.Context().Value("username")
	body, err := ioutil.ReadAll(request.Body)
    if err != nil {
        http.Error(writer, "Lỗi đọc nội dung request body", http.StatusBadRequest)
        return
    }
    defer request.Body.Close()
    var contact dto.Contacts
    err = json.Unmarshal(body, &contact)
    if err != nil {
        http.Error(writer, "Lỗi giải mã nội dung request body", http.StatusBadRequest)
        return
    }
	var contactDb dto.Contacts
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    collection := utils.MongoConnect("Contacts")
	err = collection.FindOne(ctx, bson.D{{Key: "username", Value: username}}).Decode(&contactDb)

	if err == nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "Không thể tạo thêm request" }`))
		return
	}
	scope := []dto.Scope{{
		TypeProperty:"",
		Type:"",
		City:"",
		District:"",
	}}
	doc := bson.D{
        primitive.E{Key: "username", Value: username}, 
        primitive.E{Key: "type", Value: "ca-nhan"}, 
        primitive.E{Key: "name", Value: contact.Name},
        primitive.E{Key: "avatar", Value: contact.Avatar},
        primitive.E{Key: "phoneNumber", Value: contact.PhoneNumber},
        primitive.E{Key: "city", Value: ""},
        primitive.E{Key: "district", Value: ""},
        primitive.E{Key: "ward", Value: ""},
        primitive.E{Key: "street", Value: ""},
        primitive.E{Key: "description", Value: ""},
        primitive.E{Key: "scope", Value: scope},
        primitive.E{Key: "status", Value: "pending"},
		primitive.E{Key: "createdAt", Value: contact.CreatedAt},
    }
    result, _ := collection.InsertOne(ctx, doc)
	json.NewEncoder(writer).Encode(result)
}

func UpdateAgency(writer http.ResponseWriter, request *http.Request) {
	username := request.Context().Value("username")
	body, err := ioutil.ReadAll(request.Body)
    if err != nil {
        http.Error(writer, "Lỗi đọc nội dung request body", http.StatusBadRequest)
        return
    }
    defer request.Body.Close()
    var contact dto.Contacts
    err = json.Unmarshal(body, &contact)
    if err != nil {
        http.Error(writer, "Lỗi giải mã nội dung request body", http.StatusBadRequest)
        return
    }
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    collection := utils.MongoConnect("Contacts")
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "avatar", Value: contact.Avatar},
			{Key: "name", Value: contact.Name},
			{Key: "phoneNumber", Value: contact.PhoneNumber},
			{Key: "city", Value: contact.City},
			{Key: "district", Value: contact.District},
			{Key: "ward", Value: contact.Ward},
			{Key: "street", Value: contact.Street},
			{Key: "description", Value: contact.Description},
			{Key: "scope", Value: contact.Scope},
		}},
	}
    result, err := collection.UpdateOne(ctx, bson.D{{Key: "username", Value: username}}, update)
    if err != nil {
        writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "Sửa thông tin không thành công" }`))
		return
    }
    json.NewEncoder(writer).Encode(result)
}

func GetContactUser(writer http.ResponseWriter, request *http.Request) {
	username := request.Context().Value("username")
	var contactDb dto.Contacts
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    collection := utils.MongoConnect("Contacts")
	err := collection.FindOne(ctx, bson.D{{Key: "username", Value: username}}).Decode(&contactDb)
	if err != nil {
		writer.Write([]byte(`{ "message": "no data" }`))
		return
	}
	json.NewEncoder(writer).Encode(contactDb)
}

func GetContactDetail(writer http.ResponseWriter, request *http.Request) {
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
	var contact dto.Contacts
	id, _ := mux.Vars(request)["id"]
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
	panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := utils.MongoConnect("Contacts")
	err = collection.FindOne(ctx, bson.D{{Key: "_id", Value: objID}}).Decode(&contact)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	matchFilter :=  bson.D{{Key: "user", Value: contact.Username}}
	sortStage := bson.D{{Key: "$sort", Value: bson.D{{Key: "_id", Value: -1}}}}
	limitStage := bson.D{{Key: "$limit", Value: pageSize}}
	skipStage := bson.D{{Key: "$skip", Value: skip}}
	var properties []dto.RelatedProperties
	collectionProperty := utils.MongoConnect("Properties")
	count, err := collectionProperty.CountDocuments(ctx, matchFilter)
	if err != nil {
		panic(err)
	}
	cursor, err := collectionProperty.Aggregate(ctx, mongo.Pipeline{bson.D{{Key: "$match", Value: matchFilter}}, sortStage , skipStage, limitStage})
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var property dto.RelatedProperties
		cursor.Decode(&property)
		properties = append(properties, property)
	}
	if err := cursor.Err(); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	responseData := dto.ResponseContactData{
		Data:  contact,
		PropertiesData: properties,
		Total: count,
	}
	json.NewEncoder(writer).Encode(responseData)
}

func DeleteRequestAgency(writer http.ResponseWriter, request *http.Request) {
	username := request.Context().Value("username")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
	collection := utils.MongoConnect("Contacts")
    deleteResult, _ := collection.DeleteOne(ctx, bson.D{{Key: "username", Value: username}})
    if deleteResult.DeletedCount == 0 {
        writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "Xóa không thành công" }`))
		return
    }
    json.NewEncoder(writer).Encode(deleteResult)
}

func GetAllContact(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("content-type", "application/json")
	typeContact := request.URL.Query().Get("contactType")
	keywordSearch := request.URL.Query().Get("k")
	filterType := request.URL.Query().Get("type")
	filterTypeProperty := request.URL.Query().Get("typeProperty")
	filterCity := request.URL.Query().Get("city")
	filterDistrict := request.URL.Query().Get("district")
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
	var output []dto.Contacts
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{Key: "type", Value: typeContact}}
	filter = append(filter, bson.E{ Key: "status", Value: "active" })
	if keywordSearch != "" {
		keywordFilter := bson.E{ Key:"name", Value: bson.M{"$regex": keywordSearch, "$options": "i"}}
		filter = append(filter, keywordFilter)
	}
	if filterType != "" || filterTypeProperty != "" || filterCity != ""{
		matchElement := bson.D{}
		if filterType != "" {
			typeMatch := bson.E{ Key:"typeproperty", Value: filterType}
			matchElement = append(matchElement, typeMatch)
		}
		if filterTypeProperty != "" {
			typePropertyMatch := bson.E{ Key:"type", Value: filterTypeProperty}
			matchElement = append(matchElement, typePropertyMatch)
		}
		if filterCity != "" {
			cityMatch := bson.E{ Key:"city", Value: filterCity}
			matchElement = append(matchElement, cityMatch)
		}
		if filterDistrict != "" {
			districtMatch := bson.E{ Key:"district", Value: filterDistrict}
			matchElement = append(matchElement, districtMatch)
		}
		match := bson.E{Key: "scope", Value: bson.M{ "$elemMatch" : matchElement }}
		filter = append(filter, match)
    }
	sortQuery := bson.D{{Key: "_id", Value: -1}}
	count, err := utils.MongoConnect("Contacts").CountDocuments(ctx, filter)
	if err != nil {
		panic(err)
	}
	matchStage := bson.D{{Key: "$match", Value: filter}}
	sortStage := bson.D{{Key: "$sort", Value: sortQuery}}
	limitStage := bson.D{{Key: "$limit", Value: pageSize}}
	skipStage := bson.D{{Key: "$skip", Value: skip}}
	cursor, err := utils.MongoConnect("Contacts").Aggregate(ctx, mongo.Pipeline{matchStage, sortStage, skipStage, limitStage})
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var contact dto.Contacts
		cursor.Decode(&contact)
		output = append(output, contact)
	}
	if err := cursor.Err(); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	
	responseData := dto.ResponseData{
		Data:  output,
		Total: count,
	}
	json.NewEncoder(writer).Encode(responseData)
}