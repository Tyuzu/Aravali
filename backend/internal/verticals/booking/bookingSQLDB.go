package booking

import (
	"context"

	"scav/config"
	"scav/infra"
)

var (
	slotsTable    = config.Tables.SlotTable
	bookingsTable = config.Tables.BookingsTable
	dateCapsTable = config.Tables.DateCapsTable
	tiersTable    = config.Tables.TiersTable
)

// DB helper wrappers for booking package.

func SQLFindSlots(ctx context.Context, app *infra.Deps, query string, args []any, result any) error {
	return app.SQLDB.FindMany(ctx, slotsTable, query, args, result)
}

func SQLFindBookings(ctx context.Context, app *infra.Deps, query string, args []any, result any) error {
	return app.SQLDB.FindMany(ctx, bookingsTable, query, args, result)
}

func SQLFindTiers(ctx context.Context, app *infra.Deps, query string, args []any, result any) error {
	return app.SQLDB.FindMany(ctx, tiersTable, query, args, result)
}

func SQLCountBookings(ctx context.Context, app *infra.Deps, query string, args []any) (int64, error) {
	return app.SQLDB.CountDocuments(ctx, bookingsTable, query, args)
}

func SQLFindSlotByID(ctx context.Context, app *infra.Deps, id string, out *Slot) error {
	query := "id = $1"
	args := []any{id}

	return app.SQLDB.FindOne(ctx, slotsTable, query, args, out)
}

func SQLFindTierByID(ctx context.Context, app *infra.Deps, id string, out *Tier) error {
	query := "id = $1"
	args := []any{id}

	return app.SQLDB.FindOne(ctx, tiersTable, query, args, out)
}

func SQLFindDateCap(ctx context.Context, app *infra.Deps, entityType, entityId, date string, out *DateCap) error {
	query := "entityType = $1 AND entityId = $2 AND date = $3"
	args := []any{entityType, entityId, date}

	return app.SQLDB.FindOne(ctx, dateCapsTable, query, args, out)
}

func SQLFindDateBookings(ctx context.Context, app *infra.Deps, entityType, entityId, date string, out any) error {
	query := "entityType = $1 AND entityId = $2 AND date = $3 AND status != $4"
	args := []any{entityType, entityId, date, StatusCancelled}

	return app.SQLDB.FindMany(ctx, bookingsTable, query, args, out)
}

func SQLFindVendorAvailability(ctx context.Context, app *infra.Deps, vendorId string, date string, out any) error {
	query := "vendorid = $1 AND start_date <= $2 AND end_date >= $3"
	args := []any{vendorId, date, date}

	return app.SQLDB.FindMany(ctx, config.Tables.VendorAvailabilityTable, query, args, out)
}

func SQLInsertBooking(ctx context.Context, app *infra.Deps, b Booking) error {
	return app.SQLDB.InsertOne(ctx, bookingsTable, b)
}

func SQLUpdateBookingStatusByID(ctx context.Context, app *infra.Deps, bookingID string, update map[string]any, out *Booking) error {
	query := "id = $1"
	args := []any{bookingID}

	if _, err := app.SQLDB.UpdateOne(ctx, bookingsTable, query, args, update); err != nil {
		return err
	}

	return app.SQLDB.FindOne(ctx, bookingsTable, query, args, out)
}

func SQLUpdateDateCapacity(ctx context.Context, app *infra.Deps, entityType, entityId, date string, payload map[string]any) (int64, error) {
	query := "entityType = $1 AND entityId = $2 AND date = $3"
	args := []any{entityType, entityId, date}

	return app.SQLDB.UpdateOne(ctx, dateCapsTable, query, args, payload)
}

func SQLDeleteSlotByID(ctx context.Context, app *infra.Deps, slotID string) (int64, error) {
	query := "id = $1"
	args := []any{slotID}

	return app.SQLDB.DeleteOne(ctx, slotsTable, query, args)
}

func SQLDeleteBookingsBySlot(ctx context.Context, app *infra.Deps, slotID string) (int64, error) {
	query := "slotId = $1"
	args := []any{slotID}

	return app.SQLDB.DeleteMany(ctx, bookingsTable, query, args)
}

func SQLInsertSlotsMany(ctx context.Context, app *infra.Deps, docs []any) error {
	return app.SQLDB.InsertMany(ctx, slotsTable, docs)
}

func SQLInsertTier(ctx context.Context, app *infra.Deps, t Tier) error {
	return app.SQLDB.InsertOne(ctx, tiersTable, t)
}

func SQLDeleteTierByID(ctx context.Context, app *infra.Deps, tierID string) (int64, error) {
	query := "id = $1"
	args := []any{tierID}

	return app.SQLDB.DeleteOne(ctx, tiersTable, query, args)
}

func SQLInsertSlot(ctx context.Context, app *infra.Deps, s Slot) error {
	return app.SQLDB.InsertOne(ctx, slotsTable, s)
}
