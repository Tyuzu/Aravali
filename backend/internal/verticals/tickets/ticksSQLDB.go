// File: internal/verticals/tickets/ticksSQLDB.go

package tickets

import (
	"context"
	"scav/config"
	"scav/infra"
)

var ticketsTable = config.Tables.TicketsTable
var bookingsTable = config.Tables.BookingsTable
var purchasedTicketsTable = config.Tables.PurchasedTicketsTable
var refundsTable = config.Tables.RefundsTable

// Tickets DB helpers
func FindTicketsByEvent(ctx context.Context, app *infra.Deps, eventID string, out *[]Ticket) error {

}

func FindTicketByID(ctx context.Context, app *infra.Deps, eventID, ticketID string) (*Ticket, error) {
	var t Ticket
	return &t, nil
}

func InsertTicketDB(ctx context.Context, app *infra.Deps, t Ticket) error {
}

func UpdateTicketDB(ctx context.Context, app *infra.Deps, query string, args []any, update map[string]any) (int64, error) {
}

func DeleteTicketDB(ctx context.Context, app *infra.Deps, query string, args []any) (int64, error) {
}

// Purchased tickets helpers
func FindPurchasedTickets(ctx context.Context, app *infra.Deps, query string, args []any, out *[]PurchasedTicket) error {
}

func FindPurchasedTicketByUnique(ctx context.Context, app *infra.Deps, eventID, uniqueCode string) (*PurchasedTicket, error) {
	var p PurchasedTicket
	return &p, nil
}

func InsertPurchasedTicket(ctx context.Context, app *infra.Deps, p PurchasedTicket) error {
}

func UpdatePurchasedTicket(ctx context.Context, app *infra.Deps, query string, args []any, update map[string]any) (int64, error) {
}

// Refunds
func FindRefunds(ctx context.Context, app *infra.Deps, query string, args []any, out *[]RefundRequest) error {
}

func InsertRefund(ctx context.Context, app *infra.Deps, refund RefundRequest) error {
}

// Events helper (events table is named literal "events")
func FindEventByID(ctx context.Context, app *infra.Deps, eventID string, out interface{}) error {

}
