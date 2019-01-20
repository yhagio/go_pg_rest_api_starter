package controllers

import (
	"encoding/json"
	"net/http"
)

func setSuccessStatus(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	return
}

func sendErrorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorMessage{Message: message})
	return
}
