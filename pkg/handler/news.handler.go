package handler

import (
	"fmt"
	"net/http"
)

func GetAllTodo(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "e!")
}