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
	return app.DB.FindMany(ctx, ticketsTable, map[string]any{"eventid": eventID}, out)
}

func SQLFindTicketByID(ctx context.Context, app *infra.Deps, eventID, ticketID string) (*Ticket, error) {
	var t Ticket
	if err := app.DB.FindOne(ctx, ticketsTable, map[string]any{"eventid": eventID, "ticketid": ticketID}, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

func SQLInsertTicketDB(ctx context.Context, app *infra.Deps, t Ticket) error {
	return app.DB.Insert(ctx, ticketsTable, t)
}

func SQLUpdateTicketDB(ctx context.Context, app *infra.Deps, filter map[string]any, update map[string]any) (any, error) {
	return app.DB.UpdateOne(ctx, ticketsTable, filter, update)
}

func SQLDeleteTicketDB(ctx context.Context, app *infra.Deps, filter map[string]any) (any, error) {
	return app.DB.DeleteOne(ctx, ticketsTable, filter)
}

// Purchased tickets helpers
func SQLFindPurchasedTickets(ctx context.Context, app *infra.Deps, filter map[string]any, out *[]PurchasedTicket) error {
	return app.DB.FindMany(ctx, purchasedTicketsTable, filter, out)
}

func SQLFindPurchasedTicketByUnique(ctx context.Context, app *infra.Deps, eventID, uniqueCode string) (*PurchasedTicket, error) {
	var p PurchasedTicket
	if err := app.DB.FindOne(ctx, purchasedTicketsTable, map[string]any{"eventid": eventID, "uniquecode": uniqueCode}, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func SQLInsertPurchasedTicket(ctx context.Context, app *infra.Deps, p PurchasedTicket) error {
	return app.DB.InsertOne(ctx, purchasedTicketsTable, p)
}

func SQLUpdatePurchasedTicket(ctx context.Context, app *infra.Deps, filter map[string]any, update map[string]any) (any, error) {
	return app.DB.UpdateOne(ctx, purchasedTicketsTable, filter, update)
}

// Refunds
func SQLFindRefunds(ctx context.Context, app *infra.Deps, filter map[string]any, out *[]RefundRequest) error {
	return app.DB.FindMany(ctx, refundsTable, filter, out)
}

func SQLInsertRefund(ctx context.Context, app *infra.Deps, refund RefundRequest) error {
	return app.DB.Insert(ctx, refundsTable, refund)
}

// Events helper (events table is named literal "events")
func SQLFindEventByID(ctx context.Context, app *infra.Deps, eventID string, out interface{}) error {
	return app.DB.FindOne(ctx, "events", map[string]any{"eventid": eventID}, out)
}
