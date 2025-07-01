package api

import (
	"encoding/json"
	"log"
	"net/http"
	"wallet/internal/repository"
	"wallet/internal/service"
)

func HandleError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	var statusCode int
	switch err {
	case repository.ErrWalletNotFound:
		statusCode = http.StatusNotFound
	case service.ErrInsufficientFunds, service.ErrInvalidAmount:
		statusCode = http.StatusBadRequest
	default:
		statusCode = http.StatusInternalServerError
	}

	w.WriteHeader(statusCode)

	jsonResponse := struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	}

	if err := json.NewEncoder(w).Encode(jsonResponse); err != nil {
		log.Printf("Failed to encode error response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
