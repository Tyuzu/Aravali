// File: internal/places/tabs/services.go

package places

import (
	"net/http"
	"scav/utils"

	"github.com/julienschmidt/httprouter"
)

// 🏢 Services
func GetService(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	placeID := requireParam(ps, "placeid")
	serviceID := requireParam(ps, "serviceId")
	if placeID == "" || serviceID == "" {
		writeError(w, http.StatusBadRequest, "missing placeid or serviceId")
		return
	}
	writeJSON(w, http.StatusOK, resourcePayload("service", placeID, serviceID))
}

func PostService(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	placeID := requireParam(ps, "placeid")
	if placeID == "" {
		writeError(w, http.StatusBadRequest, "missing placeid")
		return
	}
	payload := map[string]any{"status": "created", "resource": "service", "placeId": placeID, "id": utils.GenerateRandomString(16)}
	if err := decodeSimpleJSON(r, &payload); err != nil && err.Error() != "EOF" {
		writeError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}
	writeJSON(w, http.StatusCreated, payload)
}

func PutService(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	placeID := requireParam(ps, "placeid")
	serviceID := requireParam(ps, "serviceId")
	if placeID == "" || serviceID == "" {
		writeError(w, http.StatusBadRequest, "missing placeid or serviceId")
		return
	}
	payload := map[string]any{"status": "updated", "resource": "service", "placeId": placeID, "id": serviceID}
	if err := decodeSimpleJSON(r, &payload); err != nil && err.Error() != "EOF" {
		writeError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}
	writeJSON(w, http.StatusOK, payload)
}

func DeleteService(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	placeID := requireParam(ps, "placeid")
	serviceID := requireParam(ps, "serviceId")
	if placeID == "" || serviceID == "" {
		writeError(w, http.StatusBadRequest, "missing placeid or serviceId")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "deleted", "resource": "service", "placeId": placeID, "id": serviceID})
}

func GetServices(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	placeID := requireParam(ps, "placeid")
	if placeID == "" {
		writeError(w, http.StatusBadRequest, "missing placeid")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "resource": "services", "placeId": placeID, "items": []any{}})
}
