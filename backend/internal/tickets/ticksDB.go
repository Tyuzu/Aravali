package tickets

import (
	"context"
	"scav/config"
	"scav/infra"
)

var ticketsCollection = config.Collections.TicketsCollection
var bookingsCollection = config.Collections.BookingsCollection
var purchasedTicketsCollection = config.Collections.PurchasedTicketsCollection
var refundsCollection = config.Collections.RefundsCollection

// Tickets DB helpers
func FindTicketsByEvent(ctx context.Context, app *infra.Deps, eventID string, out *[]Ticket) error {
	return app.DB.FindMany(ctx, ticketsCollection, map[string]any{"eventid": eventID}, out)
}

func FindTicketByID(ctx context.Context, app *infra.Deps, eventID, ticketID string) (*Ticket, error) {
	var t Ticket
	if err := app.DB.FindOne(ctx, ticketsCollection, map[string]any{"eventid": eventID, "ticketid": ticketID}, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

func InsertTicketDB(ctx context.Context, app *infra.Deps, t Ticket) error {
	return app.DB.Insert(ctx, ticketsCollection, t)
}

func UpdateTicketDB(ctx context.Context, app *infra.Deps, filter map[string]any, update map[string]any) (any, error) {
	return app.DB.UpdateOne(ctx, ticketsCollection, filter, update)
}

func DeleteTicketDB(ctx context.Context, app *infra.Deps, filter map[string]any) (any, error) {
	return app.DB.DeleteOne(ctx, ticketsCollection, filter)
}

// Purchased tickets helpers
func FindPurchasedTickets(ctx context.Context, app *infra.Deps, filter map[string]any, out *[]PurchasedTicket) error {
	return app.DB.FindMany(ctx, purchasedTicketsCollection, filter, out)
}

func FindPurchasedTicketByUnique(ctx context.Context, app *infra.Deps, eventID, uniqueCode string) (*PurchasedTicket, error) {
	var p PurchasedTicket
	if err := app.DB.FindOne(ctx, purchasedTicketsCollection, map[string]any{"eventid": eventID, "uniquecode": uniqueCode}, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func InsertPurchasedTicket(ctx context.Context, app *infra.Deps, p PurchasedTicket) error {
	return app.DB.InsertOne(ctx, purchasedTicketsCollection, p)
}

func UpdatePurchasedTicket(ctx context.Context, app *infra.Deps, filter map[string]any, update map[string]any) (any, error) {
	return app.DB.UpdateOne(ctx, purchasedTicketsCollection, filter, update)
}

// Refunds
func FindRefunds(ctx context.Context, app *infra.Deps, filter map[string]any, out *[]RefundRequest) error {
	return app.DB.FindMany(ctx, refundsCollection, filter, out)
}

func InsertRefund(ctx context.Context, app *infra.Deps, refund RefundRequest) error {
	return app.DB.Insert(ctx, refundsCollection, refund)
}

// Events helper (events collection is named literal "events")
func FindEventByID(ctx context.Context, app *infra.Deps, eventID string, out interface{}) error {
	return app.DB.FindOne(ctx, "events", map[string]any{"eventid": eventID}, out)
}

