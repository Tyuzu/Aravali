package places

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload == nil {
		_, _ = w.Write([]byte("null"))
		return
	}
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]any{
		"status":  "error",
		"message": message,
	})
}

func requireParam(ps httprouter.Params, name string) string {
	return strings.TrimSpace(ps.ByName(name))
}

func resourcePayload(resourceType, placeID, resourceID string) map[string]any {
	payload := map[string]any{
		"status":   "ok",
		"resource": resourceType,
		"placeId":  placeID,
	}
	if resourceID != "" {
		payload["id"] = resourceID
	}
	return payload
}

func decodeSimpleJSON(r *http.Request, dst any) error {
	if r == nil || r.Body == nil {
		return nil
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(dst)
}
