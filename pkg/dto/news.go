package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type News struct {
	ID        	primitive.ObjectID `bson:"_id"`
	Category    string             `json:"category"`
	Tags   		[]string           `json:"tags"`
	Title 		string             `json:"title"`
	Description string             `json:"description"`
	Thumbnail 	string             `json:"thumbnail"`
	Source 		string             `json:"source"`
	Content 	[]interface{}      `json:"content"`
	User		string			   `json:"user"` 
	CreatedAt 	string  		   `json:"createdAt"`
}

type NewsDetail struct {
	ID        	primitive.ObjectID `bson:"_id"`
	Category    string             `json:"category"`
	Tags   		[]string           `json:"tags"`
	Title 		string             `json:"title"`
	Description string             `json:"description"`
	Thumbnail 	string             `json:"thumbnail"`
	Source 		string             `json:"source"`
	Content 	[]ContentNews      `json:"content"`
	User		string			   `json:"user"` 
	RelatedUser	[]UserRelated      `json:"relatedUser"`
	CreatedAt 	string  		   `json:"createdAt"`
}

type ContentNews struct {
	ID        	string		 `json:"id"`
	Content 	string		 `json:"content,omitempty"`
	IsBold		bool		 `json:"isBold,omitempty"`
	IsItalic	bool		 `json:"isItalic,omitempty"`
	Size 		string		 `json:"size,omitempty"`
	Image	 	string		 `json:"image,omitempty"`
	Description string		 `json:"description,omitempty"`
	Url			string		 `json:"url,omitempty"`
	HasHeader	bool		 `json:"hasHeader,omitempty"`
	TableRow	[][]string	 `json:"tableRow,omitempty"`
	TableHead	[]string	 `json:"tableHead,omitempty"`
	List		[]string	 `json:"list,omitempty"`
	Article		[]Article	 `json:"article,omitempty"`
}

type Article struct {
	Title 		string             `json:"title"`
	Url		 	string             `json:"url"`
}