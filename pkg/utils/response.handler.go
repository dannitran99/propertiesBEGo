package utils

import "net/http"

func StatusNotFound(writer http.ResponseWriter) {
	http.Error(writer, "404", http.StatusNotFound)
}

func StatusBadRequest(writer http.ResponseWriter) {
	http.Error(writer, "400", http.StatusBadRequest)
}

func StatusForbidden(writer http.ResponseWriter) {
	http.Error(writer, "403", http.StatusForbidden)
}

func StatusInternalServerError(writer http.ResponseWriter) {
	http.Error(writer, "500", http.StatusInternalServerError)
}

func StatusNotImplemented(writer http.ResponseWriter) {
	http.Error(writer, "501", http.StatusNotImplemented)
}