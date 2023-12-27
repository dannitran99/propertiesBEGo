package dto

type LoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}