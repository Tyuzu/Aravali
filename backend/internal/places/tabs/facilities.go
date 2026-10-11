// File: internal/places/tabs/facilities.go

package places

import (
	"net/http"
	"scav/utils"

	"github.com/julienschmidt/httprouter"
)

// 🌳 Facilities
func GetFacility(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	placeID := requireParam(ps, "placeid")
	facilityID := requireParam(ps, "facilityId")
	if placeID == "" || facilityID == "" {
		writeError(w, http.StatusBadRequest, "missing placeid or facilityId")
		return
	}
	writeJSON(w, http.StatusOK, resourcePayload("facility", placeID, facilityID))
}

func PostFacility(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	placeID := requireParam(ps, "placeid")
	if placeID == "" {
		writeError(w, http.StatusBadRequest, "missing placeid")
		return
	}
	payload := map[string]any{"status": "created", "resource": "facility", "placeId": placeID, "id": utils.GenerateRandomString(16)}
	if err := decodeSimpleJSON(r, &payload); err != nil && err.Error() != "EOF" {
		writeError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}
	writeJSON(w, http.StatusCreated, payload)
}

func PutFacility(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	placeID := requireParam(ps, "placeid")
	facilityID := requireParam(ps, "facilityId")
	if placeID == "" || facilityID == "" {
		writeError(w, http.StatusBadRequest, "missing placeid or facilityId")
		return
	}
	payload := map[string]any{"status": "updated", "resource": "facility", "placeId": placeID, "id": facilityID}
	if err := decodeSimpleJSON(r, &payload); err != nil && err.Error() != "EOF" {
		writeError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}
	writeJSON(w, http.StatusOK, payload)
}

func DeleteFacility(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	placeID := requireParam(ps, "placeid")
	facilityID := requireParam(ps, "facilityId")
	if placeID == "" || facilityID == "" {
		writeError(w, http.StatusBadRequest, "missing placeid or facilityId")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "deleted", "resource": "facility", "placeId": placeID, "id": facilityID})
}

func GetFacilities(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	placeID := requireParam(ps, "placeid")
	if placeID == "" {
		writeError(w, http.StatusBadRequest, "missing placeid")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "resource": "facilities", "placeId": placeID, "items": []any{}})
}
