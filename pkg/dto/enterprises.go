package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type Enterprise struct {
	ID        		primitive.ObjectID 	`bson:"_id"`
	Logo	  		string             	`json:"logo"`
	Banner			string				`json:"banner"`
	Name  	  		string             	`json:"name"`
	City    		string    			`json:"city"`
	District    	string    			`json:"district"`
	Ward    		string    			`json:"ward"`
	Street    		string    			`json:"street"`
	BusinessField	string				`json:"businessField"`
	SubBusiness		[]string			`json:"subBusiness"`
	Description    	string    			`json:"description"`
	PhoneNumber    	string    			`json:"phoneNumber"`
	Email	    	string    			`json:"email"`
	Website    		string    			`json:"website"`
	Pinned			string				`json:"pinned"`
	CreatedAt 		string  			`json:"createdAt"`
}

type PinnedEnterprise struct {
	ID				primitive.ObjectID  `json:"id"`
	Pinned			string             	`json:"pinned"`
}