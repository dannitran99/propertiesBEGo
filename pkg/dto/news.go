package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type News struct {
	ID        	primitive.ObjectID `bson:"_id"`
	Category    string             `json:"category"`
	Tags   		[]string           `json:"tags"`
	Title 		string             `json:"title"`
	Description string             `json:"description"`
	Thumbnail 	string             `json:"thumbnail"`
	Content 	[]interface{}      `json:"content"`
	User		string			   `json:"user"` 
	CreatedAt 	string  		   `json:"createdAt"`
}