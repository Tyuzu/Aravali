package tickets

import (
	"context"
	"net/http"
	"scav/infra"
	"scav/utils"
	"time"
)

// GetAvailableSeats returns the list of available seats for an event
func GetAvailableSeats(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventID := utils.GetParam(r, "eventid")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var tickets []Ticket
		if err := FindTicketsByEvent(ctx, app, eventID, &tickets); err != nil || len(tickets) == 0 {
			http.Error(w, `{"error": "No tickets found for this event"}`, http.StatusNotFound)
			return
		}

		// Use first ticket document's seats
		availableSeats := make([]string, 0, len(tickets[0].Seats))
		for _, seat := range tickets[0].Seats {
			if seat.Status == "available" {
				availableSeats = append(availableSeats, seat.SeatID)
			}
		}

		// Ensure empty slice instead of nil
		if availableSeats == nil {
			availableSeats = []string{}
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{"seats": availableSeats})
	}
}

// GetTicketSeats returns the seats for a specific ticket
func GetTicketSeats(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventID := utils.GetParam(r, "eventid")
		ticketID := utils.GetParam(r, "ticketid")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		ticketPtr, err := FindTicketByID(ctx, app, eventID, ticketID)
		if err != nil || ticketPtr == nil {
			http.Error(w, "Ticket not found", http.StatusNotFound)
			return
		}
		ticket := *ticketPtr

		// Ensure empty slice if Seats is nil
		if ticket.Seats == nil {
			ticket.Seats = []Seat{}
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"seats":   ticket.Seats,
		})
	}
}
