package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type Properties struct {
	ID        primitive.ObjectID `bson:"_id"`
	Picture   []string 			 `json:"picture"`
	Name      string             `json:"name"`
	Price     int64              `json:"price"`
	Area      int32              `json:"area"`
	Address   string  			 `json:"address"`
	CreatedAt string  			 `json:"createdAt"`
}