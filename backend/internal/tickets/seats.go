// File: internal/tickets/seats.go

package tickets

import (
	"context"
	"encoding/json"
	"net/http"
	"scav/config/mqevent"
	"scav/infra"
	log "scav/infra/logger"
	"scav/infra/mq"
	"scav/utils"
	"time"
)

// LockSeats locks specific seats for a user
func LockSeats(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventID := utils.GetParam(r, "eventid")

		// SECURITY: Get userID from authenticated request context, not from client
		userID := utils.GetUserIDFromRequest(r)
		if userID == "" {
			http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
			return
		}

		var req struct {
			Seats []string `json:"seats"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"Invalid request"}`, http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		update := map[string]any{
			"status": "locked",
			"userid": userID,
		}

		if _, err := UpdateTicketDB(ctx, app, "eventid = $1", []any{eventID}, update); err != nil {
			http.Error(w, `{"error":"Failed to lock seats"}`, http.StatusInternalServerError)
			return
		}

		if err := mq.PublishWithMeta(ctx, app.MQ, mqevent.SeatsLockedEvent, mqevent.SeatsLockedPayload{}); err != nil {
			log.Printf("failed to publish seats locked event: %v", err)
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{"success": true, "message": "Seats locked successfully"})
	}
}

// UnlockSeats releases locked seats
func UnlockSeats(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventID := utils.GetParam(r, "eventid")

		// SECURITY: Get userID from authenticated request context, not from client
		userID := utils.GetUserIDFromRequest(r)
		if userID == "" {
			http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
			return
		}

		var req struct {
			Seats []string `json:"seats"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"Invalid request"}`, http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		update := map[string]any{
			"status": "available",
			"userid": nil,
		}

		if _, err := UpdateTicketDB(ctx, app, "eventid = $1 AND userid = $2", []any{eventID, userID}, update); err != nil {
			http.Error(w, `{"error":"Failed to unlock seats"}`, http.StatusInternalServerError)
			return
		}

		if err := mq.PublishWithMeta(ctx, app.MQ, mqevent.SeatsUnlockedEvent, mqevent.SeatsUnlockedPayload{}); err != nil {
			log.Printf("failed to publish seats unlocked event: %v", err)
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{"success": true, "message": "Seats unlocked successfully"})
	}
}

// ConfirmSeatPurchase marks locked seats as booked
func ConfirmSeatPurchase(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventID := utils.GetParam(r, "eventid")
		ticketID := utils.GetParam(r, "ticketid")

		// SECURITY: Get userID from authenticated request context, not from client
		userID := utils.GetUserIDFromRequest(r)
		if userID == "" {
			http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
			return
		}

		var req struct {
			Seats []string `json:"seats"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"Invalid request"}`, http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		ticketPtr, err := FindTicketByID(ctx, app, eventID, ticketID)
		if err != nil || ticketPtr == nil {
			http.Error(w, `{"error":"Ticket not found"}`, http.StatusNotFound)
			return
		}
		ticket := *ticketPtr

		// Verify all requested seats are locked to this user
		for _, seat := range ticket.Seats {
			for _, id := range req.Seats {
				if seat.SeatID == id && (seat.Status != "locked" || seat.UserID != userID) {
					http.Error(w, `{"error":"Some seats are not properly locked or not locked to you"}`, http.StatusConflict)
					return
				}
			}
		}

		update := map[string]any{"status": "booked"}
		if _, err := UpdateTicketDB(ctx, app, "ticketid = $1 AND eventid = $2", []any{ticketID, eventID}, update); err != nil {
			http.Error(w, `{"error":"Failed to confirm purchase"}`, http.StatusInternalServerError)
			return
		}

		if err := mq.PublishWithMeta(ctx, app.MQ, mqevent.SeatPurchaseConfirmedEvent, mqevent.SeatPurchaseConfirmedPayload{}); err != nil {
			log.Printf("failed to publish seat purchase confirmed event: %v", err)
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{"success": true, "message": "Ticket purchased successfully"})
	}
}
