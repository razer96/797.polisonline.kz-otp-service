package helpers

import (
	"encoding/json"
	"net/http"
)

func JsonResponse(w http.ResponseWriter, status int, data interface{}) {
	dataJson, err := json.Marshal(data)
	if err != nil {
		log.Error("Error creating JSON response")
		ErrorJsonResponse(w, http.StatusInternalServerError, "Error creating JSON response")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(dataJson)
}

func ErrorJsonResponse(w http.ResponseWriter, status int, errorMessage string) {
	errorResponse := &ErrorResponse{
		Message: errorMessage,
	}
	errorJson, err := json.Marshal(errorResponse)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		log.Error("Error creating JSON response")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(errorJson)
}
