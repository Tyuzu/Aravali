// File: internal/booking/GETs_tiers_slots.go

package booking

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"scav/infra"
	"scav/utils"
)

func ListSlots(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entityType := r.URL.Query().Get("entityType")
		entityID := r.URL.Query().Get("entityId")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		query := "1 = 1"
		args := []any{}
		if entityType != "" {
			query += " AND entityType = $" + strconv.Itoa(len(args)+1)
			args = append(args, entityType)
		}
		if entityID != "" {
			query += " AND entityId = $" + strconv.Itoa(len(args)+1)
			args = append(args, entityID)
		}

		var slots []Slot
		if err := FindSlots(ctx, app, query, args, &slots); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"slots": slots,
		})
	}
}

func ListBookings(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entityType := r.URL.Query().Get("entityType")
		entityID := r.URL.Query().Get("entityId")
		status := r.URL.Query().Get("status")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		query := "1 = 1"
		args := []any{}
		if entityType != "" {
			query += " AND entityType = $" + strconv.Itoa(len(args)+1)
			args = append(args, entityType)
		}
		if entityID != "" {
			query += " AND entityId = $" + strconv.Itoa(len(args)+1)
			args = append(args, entityID)
		}
		if status != "" {
			query += " AND status = $" + strconv.Itoa(len(args)+1)
			args = append(args, status)
		}

		var bookings []Booking
		if err := FindBookings(ctx, app, query, args, &bookings); err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"bookings": bookings,
		})
	}
}

func GetDateCapacity(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entityType := r.URL.Query().Get("entityType")
		entityID := r.URL.Query().Get("entityId")
		date := r.URL.Query().Get("date")

		if entityType == "" || entityID == "" || date == "" {
			http.Error(w, "missing params", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var dc DateCap
		err := FindDateCap(ctx, app, entityType, entityID, date, &dc)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusOK, map[string]any{
				"capacity": nil,
			})
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"capacity": dc.Capacity,
		})
	}
}

func ListTiers(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entityType := r.URL.Query().Get("entityType")
		entityID := r.URL.Query().Get("entityId")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		query := "1 = 1"
		args := []any{}
		if entityType != "" {
			query += " AND entityType = $" + strconv.Itoa(len(args)+1)
			args = append(args, entityType)
		}
		if entityID != "" {
			query += " AND entityId = $" + strconv.Itoa(len(args)+1)
			args = append(args, entityID)
		}

		var tiers []Tier
		if err := FindTiers(ctx, app, query, args, &tiers); err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"tiers": tiers,
		})
	}
}
