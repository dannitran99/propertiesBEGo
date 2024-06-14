package dto

type LoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

type ID struct {
	Username string `json:"username"`
	Action   string `json:"action"`
}