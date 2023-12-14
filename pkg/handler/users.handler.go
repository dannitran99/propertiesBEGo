package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"propertiesGo/pkg/dto"
)

// import (
// 	"net/http"
// )

func Login(writer http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
    if err != nil {
        http.Error(writer, "Lỗi đọc nội dung request body", http.StatusBadRequest)
        return
    }
    defer request.Body.Close()

    // Giải mã nội dung của request body thành một struct User
    var user dto.User
    err = json.Unmarshal(body, &user)
    if err != nil {
        http.Error(writer, "Lỗi giải mã nội dung request body", http.StatusBadRequest)
        return
    }

    // In thông tin của user ra màn hình
    fmt.Println("Tên của user là: ", user.Username)
    fmt.Println("Tuổi của user là: ", user.Password)
}