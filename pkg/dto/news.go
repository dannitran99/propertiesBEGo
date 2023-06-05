package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type News struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string             `json:"name"`
	Content   string             `json:"content"`
	Thumbnail string             `json:"thumbnail"`
	CreatedAt string  			 `json:"createdAt"`
}