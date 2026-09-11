package utils

import (
	"encoding/json"
	"net/http"
	"scav/config"
)

// ----------------------
// HTTP JSON Helpers
// ----------------------

func RespondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func RespondWithError(
	w http.ResponseWriter,
	code int,
	msg string,
) {
	RespondWithJSON(
		w,
		code,
		map[string]any{
			"message": msg,
			"error":   msg,
		},
	)
}

type M map[string]interface{}

func ToJSON(v interface{}) []byte {
	data, _ := json.Marshal(v)
	return data
}

// ----------------------
// User Context Helpers
// ----------------------

func GetUserIDFromRequest(r *http.Request) string {
	ctx := r.Context()
	userID, ok := ctx.Value(config.UserIDKey).(string)
	if !ok || userID == "" {
		return ""
	}
	return userID
}

func GetUsernameFromRequest(r *http.Request) string {
	ctx := r.Context()
	UserName, ok := ctx.Value(config.UserNameKey).(string)
	if !ok || UserName == "" {
		return ""
	}
	return UserName
}
