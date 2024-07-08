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

type ResponseData struct {
	Data  interface{}
	Total int64
}

type ResponseContactData struct {
	Data           interface{} `json:"data"`
	PropertiesData interface{} `json:"propertiesData"`
	Total          int64       `json:"total"`
}