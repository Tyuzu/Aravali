package utils

import (
	"encoding/json"
	"log"
	"nae/models"
	"net/http"
	"time"
)

// -----------------------------------------------------------------------------
// Helper Functions & Response Encoders
// -----------------------------------------------------------------------------

func SendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := models.APIResponse{
		Success:   true,
		Timestamp: time.Now().Unix(),
		Data:      data,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("[Error] Failed to encode JSON response: %v", err)
	}
}

func SendError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := models.APIResponse{
		Success:   false,
		Timestamp: time.Now().Unix(),
		Message:   message,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("[Error] Failed to encode JSON error: %v", err)
	}
}
