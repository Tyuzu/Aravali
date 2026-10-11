// File: internal/places/tabs/exhibits.go

package places

import (
	"net/http"
	"scav/utils"

	"github.com/julienschmidt/httprouter"
)

// 🖼️ Exhibits
func GetExhibit(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	placeID := requireParam(ps, "placeid")
	exhibitID := requireParam(ps, "exhibitId")
	if placeID == "" || exhibitID == "" {
		writeError(w, http.StatusBadRequest, "missing placeid or exhibitId")
		return
	}
	writeJSON(w, http.StatusOK, resourcePayload("exhibit", placeID, exhibitID))
}

func PostExhibit(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	placeID := requireParam(ps, "placeid")
	if placeID == "" {
		writeError(w, http.StatusBadRequest, "missing placeid")
		return
	}
	payload := map[string]any{"status": "created", "resource": "exhibit", "placeId": placeID, "id": utils.GenerateRandomString(16)}
	if err := decodeSimpleJSON(r, &payload); err != nil && err.Error() != "EOF" {
		writeError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}
	writeJSON(w, http.StatusCreated, payload)
}

func PutExhibit(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	placeID := requireParam(ps, "placeid")
	exhibitID := requireParam(ps, "exhibitId")
	if placeID == "" || exhibitID == "" {
		writeError(w, http.StatusBadRequest, "missing placeid or exhibitId")
		return
	}
	payload := map[string]any{"status": "updated", "resource": "exhibit", "placeId": placeID, "id": exhibitID}
	if err := decodeSimpleJSON(r, &payload); err != nil && err.Error() != "EOF" {
		writeError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}
	writeJSON(w, http.StatusOK, payload)
}

func DeleteExhibit(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	placeID := requireParam(ps, "placeid")
	exhibitID := requireParam(ps, "exhibitId")
	if placeID == "" || exhibitID == "" {
		writeError(w, http.StatusBadRequest, "missing placeid or exhibitId")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "deleted", "resource": "exhibit", "placeId": placeID, "id": exhibitID})
}

func GetExhibits(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	placeID := requireParam(ps, "placeid")
	if placeID == "" {
		writeError(w, http.StatusBadRequest, "missing placeid")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "resource": "exhibits", "placeId": placeID, "items": []any{}})
}
