// File: internal/booking/bookingSQLDB.go

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

func FindSlots(ctx context.Context, app *infra.Deps, query string, args []any, result any) error {
	return nil
}

func FindBookings(ctx context.Context, app *infra.Deps, query string, args []any, result any) error {
	return nil
}

func FindTiers(ctx context.Context, app *infra.Deps, query string, args []any, result any) error {
	return nil
}

func CountBookings(ctx context.Context, app *infra.Deps, query string, args []any) (int64, error) {
	return 0, nil
}

func FindSlotByID(ctx context.Context, app *infra.Deps, id string, out *Slot) error {
	return nil
}

func FindTierByID(ctx context.Context, app *infra.Deps, id string, out *Tier) error {
	return nil
}

func FindDateCap(ctx context.Context, app *infra.Deps, entityType, entityId, date string, out *DateCap) error {
	return nil
}

func FindDateBookings(ctx context.Context, app *infra.Deps, entityType, entityId, date string, out any) error {
	return nil
}

func FindVendorAvailability(ctx context.Context, app *infra.Deps, vendorId string, date string, out any) error {
	return nil
}

func InsertBooking(ctx context.Context, app *infra.Deps, b Booking) error {
	return nil
}

func UpdateBookingStatusByID(ctx context.Context, app *infra.Deps, bookingID string, update map[string]any, out *Booking) error {
	return nil
}

func UpdateDateCapacity(ctx context.Context, app *infra.Deps, entityType, entityId, date string, payload map[string]any) (int64, error) {
	return 0, nil
}

func DeleteSlotByID(ctx context.Context, app *infra.Deps, slotID string) (int64, error) {
	return 0, nil
}

func DeleteBookingsBySlot(ctx context.Context, app *infra.Deps, slotID string) (int64, error) {
	return 0, nil
}

func InsertSlotsMany(ctx context.Context, app *infra.Deps, docs []any) error {
	return nil
}

func InsertTier(ctx context.Context, app *infra.Deps, t Tier) error {
	return nil
}

func DeleteTierByID(ctx context.Context, app *infra.Deps, tierID string) (int64, error) {
	return 0, nil
}

func InsertSlot(ctx context.Context, app *infra.Deps, s Slot) error {
	return nil
}
