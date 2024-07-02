package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type Contacts struct {
	ID        		primitive.ObjectID 	`bson:"_id"`
	Username  		string             	`json:"username"`
	Type			string				`json:"type"`
	Name  	  		string             	`json:"name"`
	Avatar 	  		string             	`json:"avatar"`
	PhoneNumber    	string    			`json:"phoneNumber"`
	City    		string    			`json:"city"`
	District    	string    			`json:"district"`
	Ward    		string    			`json:"ward"`
	Street    		string    			`json:"street"`
	Description    	string    			`json:"description"`
	Scope    		[]Scope   			`json:"scope"`
	Status 			string    			`json:"status"`
	CreatedAt 		string  			`json:"createdAt"`
}

type Scope struct {
	TypeProperty    string    			`json:"typeProperty"`
	Type    		string    			`json:"type"`
	City    		string    			`json:"city"`
	District    	string    			`json:"district"`
}