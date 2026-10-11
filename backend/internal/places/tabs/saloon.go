// File: internal/places/tabs/saloon.go

package places

import (
	"net/http"
	"scav/utils"

	"github.com/julienschmidt/httprouter"
)

func GetSaloonSlots(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	placeID := requireParam(ps, "placeid")
	if placeID == "" {
		writeError(w, http.StatusBadRequest, "missing placeid")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "resource": "saloon_slots", "placeId": placeID, "items": []any{}})
}

func PostSaloonSlot(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	placeID := requireParam(ps, "placeid")
	if placeID == "" {
		writeError(w, http.StatusBadRequest, "missing placeid")
		return
	}
	payload := map[string]any{"status": "created", "resource": "saloon_slot", "placeId": placeID, "id": utils.GenerateRandomString(16)}
	if err := decodeSimpleJSON(r, &payload); err != nil && err.Error() != "EOF" {
		writeError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}
	writeJSON(w, http.StatusCreated, payload)
}

func PutSaloonSlot(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	placeID := requireParam(ps, "placeid")
	slotID := requireParam(ps, "slotId")
	if placeID == "" || slotID == "" {
		writeError(w, http.StatusBadRequest, "missing placeid or slotId")
		return
	}
	payload := map[string]any{"status": "updated", "resource": "saloon_slot", "placeId": placeID, "id": slotID}
	if err := decodeSimpleJSON(r, &payload); err != nil && err.Error() != "EOF" {
		writeError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}
	writeJSON(w, http.StatusOK, payload)
}

func DeleteSaloonSlot(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	placeID := requireParam(ps, "placeid")
	slotID := requireParam(ps, "slotId")
	if placeID == "" || slotID == "" {
		writeError(w, http.StatusBadRequest, "missing placeid or slotId")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "deleted", "resource": "saloon_slot", "placeId": placeID, "id": slotID})
}

func BookSaloonSlot(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	placeID := requireParam(ps, "placeid")
	slotID := requireParam(ps, "slotId")
	if placeID == "" || slotID == "" {
		writeError(w, http.StatusBadRequest, "missing placeid or slotId")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "booked", "resource": "saloon_slot_booking", "placeId": placeID, "slotId": slotID})
}
