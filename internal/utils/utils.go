package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func LogError(msg string, err error) bool {
	if err != nil {
		log.Printf("%s: %v \n", msg, err)
		return true
	}
	return false
}

func ResponseWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": message,
		"code":  code,
	})
}

func ResponseSuccess(w http.ResponseWriter, status int, responseValue any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if responseValue != nil {
		json.NewEncoder(w).Encode(responseValue)
	}
}
