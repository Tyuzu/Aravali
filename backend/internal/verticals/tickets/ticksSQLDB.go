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
func SQLFindTicketsByEvent(ctx context.Context, app *infra.Deps, eventID string, out *[]Ticket) error {
	query := "eventid = $1"
	args := []any{eventID}
	return app.SQLDB.FindMany(ctx, ticketsTable, query, args, out)
}

func SQLFindTicketByID(ctx context.Context, app *infra.Deps, eventID, ticketID string) (*Ticket, error) {
	var t Ticket
	query := "eventid = $1 AND ticketid = $2"
	args := []any{eventID, ticketID}
	if err := app.SQLDB.FindOne(ctx, ticketsTable, query, args, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

func SQLInsertTicketDB(ctx context.Context, app *infra.Deps, t Ticket) error {
	return app.SQLDB.Insert(ctx, ticketsTable, t)
}

func SQLUpdateTicketDB(ctx context.Context, app *infra.Deps, query string, args []any, update map[string]any) (int64, error) {
	return app.SQLDB.UpdateOne(ctx, ticketsTable, query, args, update)
}

func SQLDeleteTicketDB(ctx context.Context, app *infra.Deps, query string, args []any) (int64, error) {
	return app.SQLDB.DeleteOne(ctx, ticketsTable, query, args)
}

// Purchased tickets helpers
func SQLFindPurchasedTickets(ctx context.Context, app *infra.Deps, query string, args []any, out *[]PurchasedTicket) error {
	return app.SQLDB.FindMany(ctx, purchasedTicketsTable, query, args, out)
}

func SQLFindPurchasedTicketByUnique(ctx context.Context, app *infra.Deps, eventID, uniqueCode string) (*PurchasedTicket, error) {
	var p PurchasedTicket
	query := "eventid = $1 AND uniquecode = $2"
	args := []any{eventID, uniqueCode}
	if err := app.SQLDB.FindOne(ctx, purchasedTicketsTable, query, args, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func SQLInsertPurchasedTicket(ctx context.Context, app *infra.Deps, p PurchasedTicket) error {
	return app.SQLDB.InsertOne(ctx, purchasedTicketsTable, p)
}

func SQLUpdatePurchasedTicket(ctx context.Context, app *infra.Deps, query string, args []any, update map[string]any) (int64, error) {
	return app.SQLDB.UpdateOne(ctx, purchasedTicketsTable, query, args, update)
}

// Refunds
func SQLFindRefunds(ctx context.Context, app *infra.Deps, query string, args []any, out *[]RefundRequest) error {
	return app.SQLDB.FindMany(ctx, refundsTable, query, args, out)
}

func SQLInsertRefund(ctx context.Context, app *infra.Deps, refund RefundRequest) error {
	return app.SQLDB.Insert(ctx, refundsTable, refund)
}

// Events helper (events table is named literal "events")
func SQLFindEventByID(ctx context.Context, app *infra.Deps, eventID string, out interface{}) error {
	query := "eventid = $1"
	args := []any{eventID}
	return app.SQLDB.FindOne(ctx, "events", query, args, out)
}
