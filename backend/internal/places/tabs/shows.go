// File: internal/places/tabs/shows.go

package places

import (
	"net/http"
	"scav/utils"

	"github.com/julienschmidt/httprouter"
)

// 🎭 Shows
func GetShow(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	placeID := requireParam(ps, "placeid")
	showID := requireParam(ps, "showId")
	if placeID == "" || showID == "" {
		writeError(w, http.StatusBadRequest, "missing placeid or showId")
		return
	}
	writeJSON(w, http.StatusOK, resourcePayload("show", placeID, showID))
}

func PostShow(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	placeID := requireParam(ps, "placeid")
	if placeID == "" {
		writeError(w, http.StatusBadRequest, "missing placeid")
		return
	}
	payload := map[string]any{"status": "created", "resource": "show", "placeId": placeID, "id": utils.GenerateRandomString(16)}
	if err := decodeSimpleJSON(r, &payload); err != nil && err.Error() != "EOF" {
		writeError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}
	writeJSON(w, http.StatusCreated, payload)
}

func PutShow(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	placeID := requireParam(ps, "placeid")
	showID := requireParam(ps, "showId")
	if placeID == "" || showID == "" {
		writeError(w, http.StatusBadRequest, "missing placeid or showId")
		return
	}
	payload := map[string]any{"status": "updated", "resource": "show", "placeId": placeID, "id": showID}
	if err := decodeSimpleJSON(r, &payload); err != nil && err.Error() != "EOF" {
		writeError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}
	writeJSON(w, http.StatusOK, payload)
}

func DeleteShow(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	placeID := requireParam(ps, "placeid")
	showID := requireParam(ps, "showId")
	if placeID == "" || showID == "" {
		writeError(w, http.StatusBadRequest, "missing placeid or showId")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "deleted", "resource": "show", "placeId": placeID, "id": showID})
}

func PostBookShow(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	placeID := requireParam(ps, "placeid")
	showID := requireParam(ps, "showId")
	if placeID == "" || showID == "" {
		writeError(w, http.StatusBadRequest, "missing placeid or showId")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "booked", "resource": "show_booking", "placeId": placeID, "showId": showID})
}

func GetShows(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	placeID := requireParam(ps, "placeid")
	if placeID == "" {
		writeError(w, http.StatusBadRequest, "missing placeid")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "resource": "shows", "placeId": placeID, "items": []any{}})
}
