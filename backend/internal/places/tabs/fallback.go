// File: internal/places/tabs/fallback.go

package places

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func notImplemented(w http.ResponseWriter, methodName string) {
	writeJSON(w, http.StatusNotImplemented, map[string]any{
		"status":  "not_implemented",
		"method":  methodName,
		"message": methodName + " is not implemented yet",
	})
}

// ❓ Fallback
func GetPlaceDetailsFallback(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	placeID := requireParam(ps, "placeid")
	if placeID == "" {
		writeError(w, http.StatusBadRequest, "missing placeid")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "resource": "place_details", "placeId": placeID})
}

func GetDetailsFallback(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	placeID := requireParam(ps, "placeid")
	if placeID == "" {
		writeError(w, http.StatusBadRequest, "missing placeid")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "resource": "details", "placeId": placeID})
}
