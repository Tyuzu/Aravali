package booking

import (
	"context"

	"scav/config"
	"scav/infra/db"
)

var (
	slotsTable    = config.Tables.SlotTable
	bookingsTable = config.Tables.BookingsTable
	dateCapsTable = config.Tables.DateCapsTable
	tiersTable    = config.Tables.TiersTable
)

// DB helper wrappers for booking package. Keeps all Mongo interactions here.

func SQLFindSlots(ctx context.Context, d db.Database, filter any, result any) error {
	return d.FindMany(ctx, slotsTable, filter, result)
}

func SQLFindBookings(ctx context.Context, d db.Database, filter any, result any) error {
	return d.FindMany(ctx, bookingsTable, filter, result)
}

func SQLFindTiers(ctx context.Context, d db.Database, filter any, result any) error {
	return d.FindMany(ctx, tiersTable, filter, result)
}

func SQLCountBookings(ctx context.Context, d db.Database, filter any) (int64, error) {
	return d.CountDocuments(ctx, bookingsTable, filter)
}

func SQLFindSlotByID(ctx context.Context, d db.Database, id string, out *Slot) error {
	return d.FindOne(ctx, slotsTable, map[string]any{"id": id}, out)
}

func SQLFindTierByID(ctx context.Context, d db.Database, id string, out *Tier) error {
	return d.FindOne(ctx, tiersTable, map[string]any{"id": id}, out)
}

func SQLFindDateCap(ctx context.Context, d db.Database, entityType, entityId, date string, out *DateCap) error {
	return d.FindOne(ctx, dateCapsTable, map[string]any{"entityType": entityType, "entityId": entityId, "date": date}, out)
}

func SQLFindDateBookings(ctx context.Context, d db.Database, entityType, entityId, date string, out any) error {
	return d.FindMany(ctx, bookingsTable, map[string]any{
		"entityType": entityType,
		"entityId":   entityId,
		"date":       date,
		"status":     map[string]any{"$ne": StatusCancelled},
	}, out)
}

func SQLFindVendorAvailability(ctx context.Context, d db.Database, vendorId string, date string, out any) error {
	return d.FindMany(ctx, config.Tables.VendorAvailabilityTable, map[string]any{"vendorid": vendorId, "start_date": map[string]any{"$lte": date}, "end_date": map[string]any{"$gte": date}}, out)
}

func SQLInsertBooking(ctx context.Context, d db.Database, b Booking) error {
	return d.InsertOne(ctx, bookingsTable, b)
}

func SQLUpdateBookingStatusByID(ctx context.Context, d db.Database, bookingID string, update any, out *Booking) error {
	return d.FindOneAndUpdate(ctx, bookingsTable, map[string]any{"id": bookingID}, update, out)
}

func SQLUpdateDateCapacity(ctx context.Context, d db.Database, entityType, entityId, date string, payload any) (any, error) {
	return d.UpdateOne(ctx, dateCapsTable, map[string]any{"entityType": entityType, "entityId": entityId, "date": date}, payload)
}

func SQLDeleteSlotByID(ctx context.Context, d db.Database, slotID string) (int64, error) {
	return d.DeleteOne(ctx, slotsTable, map[string]any{"id": slotID})
}

func SQLDeleteBookingsBySlot(ctx context.Context, d db.Database, slotID string) error {
	return d.DeleteMany(ctx, bookingsTable, map[string]any{"slotId": slotID})
}

func SQLInsertSlotsMany(ctx context.Context, d db.Database, docs []any) error {
	return d.InsertMany(ctx, slotsTable, docs)
}

func SQLInsertTier(ctx context.Context, d db.Database, t Tier) error {
	return d.InsertOne(ctx, tiersTable, t)
}

func SQLDeleteTierByID(ctx context.Context, d db.Database, tierID string) (int64, error) {
	return d.DeleteOne(ctx, tiersTable, map[string]any{"id": tierID})
}

func SQLInsertSlot(ctx context.Context, d db.Database, s Slot) error {
	return d.InsertOne(ctx, slotsTable, s)
}
