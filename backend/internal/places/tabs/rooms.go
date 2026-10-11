// File: internal/places/tabs/rooms.go

package places

import (
	"net/http"
	"scav/utils"

	"github.com/julienschmidt/httprouter"
)

func GetRoom(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	placeID := requireParam(ps, "placeid")
	roomID := requireParam(ps, "roomId")
	if placeID == "" || roomID == "" {
		writeError(w, http.StatusBadRequest, "missing placeid or roomId")
		return
	}
	writeJSON(w, http.StatusOK, resourcePayload("room", placeID, roomID))
}

func PostRoom(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	placeID := requireParam(ps, "placeid")
	if placeID == "" {
		writeError(w, http.StatusBadRequest, "missing placeid")
		return
	}
	payload := map[string]any{"status": "created", "resource": "room", "placeId": placeID, "id": utils.GenerateRandomString(16)}
	if err := decodeSimpleJSON(r, &payload); err != nil && err.Error() != "EOF" {
		writeError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}
	writeJSON(w, http.StatusCreated, payload)
}

func PutRoom(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	placeID := requireParam(ps, "placeid")
	roomID := requireParam(ps, "roomId")
	if placeID == "" || roomID == "" {
		writeError(w, http.StatusBadRequest, "missing placeid or roomId")
		return
	}
	payload := map[string]any{"status": "updated", "resource": "room", "placeId": placeID, "id": roomID}
	if err := decodeSimpleJSON(r, &payload); err != nil && err.Error() != "EOF" {
		writeError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}
	writeJSON(w, http.StatusOK, payload)
}

func DeleteRoom(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	placeID := requireParam(ps, "placeid")
	roomID := requireParam(ps, "roomId")
	if placeID == "" || roomID == "" {
		writeError(w, http.StatusBadRequest, "missing placeid or roomId")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "deleted", "resource": "room", "placeId": placeID, "id": roomID})
}

func GetRooms(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	placeID := requireParam(ps, "placeid")
	if placeID == "" {
		writeError(w, http.StatusBadRequest, "missing placeid")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "resource": "rooms", "placeId": placeID, "items": []any{}})
}
