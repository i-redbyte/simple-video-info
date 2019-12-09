package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func Respond(w http.ResponseWriter, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(data)
}

func ErrorMessage(message string) map[string]interface{} {
	return map[string]interface{}{"message": message}
}

func ReturnError(w http.ResponseWriter, err error, code int) {
	log.Println(err)
	w.WriteHeader(code)
	Respond(w, ErrorMessage(err.Error()))
}
