package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type Properties struct {
	ID        		primitive.ObjectID  `bson:"_id"`
	Type      		string              `json:"type"`
    PropertyType    string    			`json:"propertyType"`    
	City    		string    			`json:"city"`    
	District    	string    			`json:"district"`   
	Ward    		string    			`json:"ward"`  
	Street    		string    			`json:"street"` 
	Project    		string    			`json:"project"` 
	MoreInfo    	string    			`json:"moreInfo"` 
	Title    		string    			`json:"title"` 
	Description    	string    			`json:"description"` 
	Area    		int32    			`json:"area"` 
	Price    		int64    			`json:"price"` 
	PriceAvg    	float32    			`json:"priceAvg"` 
	PriceType    	string    			`json:"priceType"` 
	Images    		[]Image  			`json:"images"` 
	Url    			string    			`json:"url"` 
	Name    		string    			`json:"name"` 
	PhoneNumber    	string    			`json:"phoneNumber"` 
	Email    		string    			`json:"email"` 
	User			string				`json:"user"` 
	CreatedAt 		string  			`json:"createdAt"`
}

type PropertiesInfo struct {
	ID        		primitive.ObjectID  `bson:"_id"`
	Type      		string              `json:"type"`
    PropertyType    string    			`json:"propertyType"`    
	City    		string    			`json:"city"`    
	District    	string    			`json:"district"`   
	Ward    		string    			`json:"ward"`  
	Street    		string    			`json:"street"` 
	Project    		string    			`json:"project"` 
	MoreInfo    	string    			`json:"moreInfo"` 
	Title    		string    			`json:"title"` 
	Description    	string    			`json:"description"` 
	Area    		int32    			`json:"area"` 
	Price    		int64    			`json:"price"` 
	PriceAvg    	float32    			`json:"priceAvg"` 
	PriceType    	string    			`json:"priceType"` 
	Images    		[]Image  			`json:"images"` 
	Url    			string    			`json:"url"` 
	Name    		string    			`json:"name"`
	Avatar			string    			`json:"avatar"`
	PhoneNumber    	string    			`json:"phoneNumber"` 
	Email    		string    			`json:"email"` 
	User			string				`json:"user"` 
	CreatedAt 		string  			`json:"createdAt"`
}

type Image struct {
	Description      string              `json:"description"`
	Name     		 string              `json:"name"`
	Url      		 string              `json:"url"`
}