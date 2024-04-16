package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        	primitive.ObjectID 	`bson:"_id"`
	Username  	string             	`json:"username"`
	FullName	string				`json:"fullname"`
	Avatar		string				`json:"avatar"`
	Password  	string             	`json:"password"`
	Email 	  	string			   	`json:"email"`
	Active 	  	bool			   	`json:"active"`
	PhoneNumber string				`json:"phoneNumber"`
}

type ChangePassword struct {
	User				string          `json:"user"`
	CurrentPassword  	string          `json:"currentPassword"`
	NewPassword  		string          `json:"newPassword"`
}

type ChangeAvatar struct {
	User				string          `json:"user"`
	Avatar  			string          `json:"avatar"`
}

type UserInfo struct {
	Name				string          `json:"name"`
	PhoneNumber  		string          `json:"phoneNumber"`
	Email 	  			string			`json:"email"`
}

type UserInfoPost struct {
	User				string          `json:"user"`
	Name				string          `json:"name"`
	PhoneNumber  		string          `json:"phoneNumber"`
	Email 	  			string			`json:"email"`
}