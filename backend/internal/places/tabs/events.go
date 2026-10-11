// File: internal/places/tabs/events.go

package places

import (
	"context"
	"net/http"
	"scav/infra"
	"scav/utils"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

// 🏟️ Events

func GetEvent(app *infra.Deps) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		placeID := requireParam(ps, "placeid")
		eventID := requireParam(ps, "eventId")
		if placeID == "" || eventID == "" {
			writeError(w, http.StatusBadRequest, "missing placeid or eventId")
			return
		}
		writeJSON(w, http.StatusOK, resourcePayload("event", placeID, eventID))
	}
}

func PostEvent(app *infra.Deps) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		placeID := requireParam(ps, "placeid")
		if placeID == "" {
			writeError(w, http.StatusBadRequest, "missing placeid")
			return
		}
		payload := map[string]any{"status": "created", "resource": "event", "placeId": placeID}
		if err := decodeSimpleJSON(r, &payload); err == nil && len(payload) > 0 {
			if _, exists := payload["id"]; !exists {
				payload["id"] = utils.GenerateRandomString(16)
			}
		}
		writeJSON(w, http.StatusCreated, payload)
	}
}

func PutEvent(app *infra.Deps) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		placeID := requireParam(ps, "placeid")
		eventID := requireParam(ps, "eventId")
		if placeID == "" || eventID == "" {
			writeError(w, http.StatusBadRequest, "missing placeid or eventId")
			return
		}
		payload := map[string]any{"status": "updated", "resource": "event", "placeId": placeID, "id": eventID}
		if err := decodeSimpleJSON(r, &payload); err != nil && !errorsAsNoBody(err) {
			writeError(w, http.StatusBadRequest, "invalid JSON payload")
			return
		}
		writeJSON(w, http.StatusOK, payload)
	}
}

func DeleteEvent(app *infra.Deps) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		placeID := requireParam(ps, "placeid")
		eventID := requireParam(ps, "eventId")
		if placeID == "" || eventID == "" {
			writeError(w, http.StatusBadRequest, "missing placeid or eventId")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"status": "deleted", "resource": "event", "placeId": placeID, "id": eventID})
	}
}

func PostViewEventDetails(app *infra.Deps) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		placeID := requireParam(ps, "placeid")
		eventID := requireParam(ps, "eventId")
		if placeID == "" || eventID == "" {
			writeError(w, http.StatusBadRequest, "missing placeid or eventId")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"status":   "ok",
			"resource": "event_view",
			"placeId":  placeID,
			"eventId":  eventID,
			"viewedAt": time.Now().UTC().Format(time.RFC3339),
		})
	}
}

func GetEvents(app *infra.Deps) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		placeID := ps.ByName("placeid")
		if placeID == "" {
			http.Error(w, "missing required path parameter: placeid", http.StatusBadRequest)
			return
		}

		page := 1
		limit := 10

		if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && p > 0 {
			page = p
		}
		if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 {
			limit = l
		}

		now := time.Now()

		ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
		defer cancel()

		placeevents, total, err := getPlaceEventsPage(ctx, app, placeID, page, limit, now)
		if err != nil {
			http.Error(w, "Failed to fetch events", http.StatusInternalServerError)
			return
		}

		response := map[string]any{
			"events": placeevents,
			"total":  total,
			"page":   page,
			"limit":  limit,
		}

		utils.RespondWithJSON(w, http.StatusOK, response)
	}
}

func errorsAsNoBody(err error) bool {
	return err != nil && err.Error() == "EOF"
}
