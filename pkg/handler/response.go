package handler

import (
	"encoding/json"
	"net/http"
)

func sendJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func sendJSONErrorResponse(w http.ResponseWriter, errorMessage string, status int) {
	errorData := map[string]string{"error": errorMessage}
	sendJSONResponse(w, status, errorData)
}