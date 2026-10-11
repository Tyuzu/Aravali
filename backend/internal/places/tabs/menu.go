// File: internal/places/tabs/menu.go

package places

import (
	"net/http"
	"scav/utils"

	"github.com/julienschmidt/httprouter"
)

func GetMenuTab(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	placeID := requireParam(ps, "placeid")
	if placeID == "" {
		writeError(w, http.StatusBadRequest, "missing placeid")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "resource": "menu", "placeId": placeID, "items": []any{}})
}

func PostMenuTab(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	placeID := requireParam(ps, "placeid")
	if placeID == "" {
		writeError(w, http.StatusBadRequest, "missing placeid")
		return
	}
	payload := map[string]any{"status": "created", "resource": "menu", "placeId": placeID, "id": utils.GenerateRandomString(16)}
	if err := decodeSimpleJSON(r, &payload); err != nil && err.Error() != "EOF" {
		writeError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}
	writeJSON(w, http.StatusCreated, payload)
}

func PutMenuTab(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	placeID := requireParam(ps, "placeid")
	itemID := requireParam(ps, "itemId")
	if placeID == "" || itemID == "" {
		writeError(w, http.StatusBadRequest, "missing placeid or itemId")
		return
	}
	payload := map[string]any{"status": "updated", "resource": "menu", "placeId": placeID, "id": itemID}
	if err := decodeSimpleJSON(r, &payload); err != nil && err.Error() != "EOF" {
		writeError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}
	writeJSON(w, http.StatusOK, payload)
}

func DeleteMenuTab(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	placeID := requireParam(ps, "placeid")
	itemID := requireParam(ps, "itemId")
	if placeID == "" || itemID == "" {
		writeError(w, http.StatusBadRequest, "missing placeid or itemId")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "deleted", "resource": "menu", "placeId": placeID, "id": itemID})
}

func PostPlaceMenuOrder(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	placeID := requireParam(ps, "placeid")
	itemID := requireParam(ps, "itemId")
	if placeID == "" || itemID == "" {
		writeError(w, http.StatusBadRequest, "missing placeid or itemId")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "ordered", "resource": "menu_order", "placeId": placeID, "itemId": itemID})
}

func PostMenuOrder(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	placeID := requireParam(ps, "placeid")
	if placeID == "" {
		writeError(w, http.StatusBadRequest, "missing placeid")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "ordered", "resource": "menu_order", "placeId": placeID})
}
