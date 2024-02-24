package handler

import (
	"encoding/json"
	"net/http"
)

func sendJSONResponse(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return err
	}
	return nil
}

func sendJSONErrorResponse(w http.ResponseWriter, errorMessage string, status int) error {
	errorData := map[string]string{"error": errorMessage}
	if err := sendJSONResponse(w, status, errorData); err != nil {
		return err
	}
	return nil
}