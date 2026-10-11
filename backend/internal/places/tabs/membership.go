// File: internal/places/tabs/membership.go

package places

import (
	"context"
	"encoding/json"
	"net/http"
	"scav/infra"
	placedb "scav/internal/places/placedb"
	"scav/utils"
	"time"

	"github.com/julienschmidt/httprouter"
)

type Membership struct {
	ID          string    `db:"_id,omitempty" json:"_id"`
	PlaceID     string    `db:"placeId,omitempty" json:"placeId"`
	Name        string    `db:"name" json:"name"`
	Price       float64   `db:"price" json:"price"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"createdAt" json:"createdAt"`
}

// handlers use the SQL table via placedb helpers

// GET /place/:placeId/membership/:id
func GetMembership(app *infra.Deps) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		id := ps.ByName("membershipId")
		if id == "" {
			http.Error(w, "Invalid membership ID", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		var m Membership
		if err := placedb.FindOneMembership(ctx, app, "_id = $1", []any{id}, &m); err != nil {
			http.Error(w, "Membership not found", http.StatusNotFound)
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, m)
	}
}

// POST /place/:placeId/membership
func PostMembership(app *infra.Deps) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		placeID := ps.ByName("placeid")
		if placeID == "" {
			http.Error(w, "Invalid place ID", http.StatusBadRequest)
			return
		}

		var membership Membership
		if err := json.NewDecoder(r.Body).Decode(&membership); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		membership.ID = utils.GenerateRandomString(16)
		membership.PlaceID = placeID
		membership.CreatedAt = time.Now()

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		if err := placedb.InsertMembership(ctx, app, membership); err != nil {
			http.Error(w, "Failed to create membership", http.StatusInternalServerError)
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, membership)
	}
}

// PUT /place/:placeId/membership/:id
func PutMembership(app *infra.Deps) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		id := ps.ByName("membershipId")
		if id == "" {
			http.Error(w, "Invalid membership ID", http.StatusBadRequest)
			return
		}

		var update Membership
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		if _, err := placedb.UpdateMembership(ctx, app, "_id = $1", []any{id}, map[string]any{
			"name":        update.Name,
			"price":       update.Price,
			"description": update.Description,
		}); err != nil {
			http.Error(w, "Failed to update membership", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// DELETE /place/:placeId/membership/:id
func DeleteMembership(app *infra.Deps) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		id := ps.ByName("membershipId")
		if id == "" {
			http.Error(w, "Invalid membership ID", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		if _, err := placedb.DeleteMembership(ctx, app, "_id = $1", []any{id}); err != nil {
			http.Error(w, "Failed to delete membership", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// POST /place/:placeId/membership/:id/join
func PostJoinMembership(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	placeID := ps.ByName("placeid")
	membershipID := ps.ByName("membershipId")
	if placeID == "" || membershipID == "" {
		http.Error(w, "Invalid place or membership ID", http.StatusBadRequest)
		return
	}

	userID := utils.GetUserIDFromRequest(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	payload := map[string]any{
		"status":       "joined",
		"placeId":      placeID,
		"membershipId": membershipID,
		"userId":       userID,
		"joinedAt":     time.Now().UTC().Format(time.RFC3339),
	}
	utils.RespondWithJSON(w, http.StatusOK, payload)
}

// GET /place/:placeId/memberships
func GetMemberships(app *infra.Deps) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		placeID := ps.ByName("placeid")
		if placeID == "" {
			http.Error(w, "Invalid place ID", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		var memberships []Membership
		if err := placedb.FindPlaceMemberships(ctx, app, "placeid = $1", []any{placeID}, &memberships); err != nil {
			http.Error(w, "Failed to fetch memberships", http.StatusInternalServerError)
			return
		}

		if memberships == nil {
			memberships = make([]Membership, 0)
		}

		utils.RespondWithJSON(w, http.StatusOK, memberships)
	}
}
