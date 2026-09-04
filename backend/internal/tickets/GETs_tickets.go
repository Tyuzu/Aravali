package tickets

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"scav/infra"
	"scav/utils"
)

func GetTickets(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventID := utils.GetParam(r, "eventid")

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		var tickets []Ticket
		if err := app.DB.FindMany(
			ctx,
			ticketsCollection,
			map[string]any{"eventid": eventID},
			&tickets,
		); err != nil {
			http.Error(w, "Failed to fetch tickets", http.StatusInternalServerError)
			return
		}

		if tickets == nil {
			tickets = []Ticket{}
		}

		utils.RespondWithJSON(w, http.StatusOK, tickets)
	}
}

// Fetch a single ticket
func GetTicket(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventID := utils.GetParam(r, "eventid")
		ticketID := utils.GetParam(r, "ticketid")

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		var ticket Ticket
		if err := app.DB.FindOne(
			ctx,
			ticketsCollection,
			map[string]any{
				"eventid":  eventID,
				"ticketid": ticketID,
			},
			&ticket,
		); err != nil {
			http.Error(w, fmt.Sprintf("Ticket not found: %v", err), http.StatusNotFound)
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, ticket)
	}
}
