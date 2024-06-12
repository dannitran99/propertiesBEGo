package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type Contacts struct {
	ID        		primitive.ObjectID 	`bson:"_id"`
	Username  		string             	`json:"username"`
	Name  	  		string             	`json:"name"`
	Avatar 	  		string             	`json:"avatar"`
	PhoneNumber    	string    			`json:"phoneNumber"`
	Status 			string    			`json:"status"`
	CreatedAt 		string  			`json:"createdAt"`
}