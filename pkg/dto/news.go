package dto

type News struct {
	ID        int    `json:"_id"`
	Name      string `json:"name"`
	Content   string `json:"content"`
	Thumbnail string `json:"thumbnail"`
}