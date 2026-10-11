// File: internal/verticals/booking/bookingSQLDB.go

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
}

func FindBookings(ctx context.Context, app *infra.Deps, query string, args []any, result any) error {
}

func FindTiers(ctx context.Context, app *infra.Deps, query string, args []any, result any) error {
}

func CountBookings(ctx context.Context, app *infra.Deps, query string, args []any) (int64, error) {
}

func FindSlotByID(ctx context.Context, app *infra.Deps, id string, out *Slot) error {
}

func FindTierByID(ctx context.Context, app *infra.Deps, id string, out *Tier) error {
}

func FindDateCap(ctx context.Context, app *infra.Deps, entityType, entityId, date string, out *DateCap) error {
}

func FindDateBookings(ctx context.Context, app *infra.Deps, entityType, entityId, date string, out any) error {
}

func FindVendorAvailability(ctx context.Context, app *infra.Deps, vendorId string, date string, out any) error {
}

func InsertBooking(ctx context.Context, app *infra.Deps, b Booking) error {
}

func UpdateBookingStatusByID(ctx context.Context, app *infra.Deps, bookingID string, update map[string]any, out *Booking) error {
}

func UpdateDateCapacity(ctx context.Context, app *infra.Deps, entityType, entityId, date string, payload map[string]any) (int64, error) {
}

func DeleteSlotByID(ctx context.Context, app *infra.Deps, slotID string) (int64, error) {
}

func DeleteBookingsBySlot(ctx context.Context, app *infra.Deps, slotID string) (int64, error) {
}

func InsertSlotsMany(ctx context.Context, app *infra.Deps, docs []any) error {
}

func InsertTier(ctx context.Context, app *infra.Deps, t Tier) error {
}

func DeleteTierByID(ctx context.Context, app *infra.Deps, tierID string) (int64, error) {
}

func InsertSlot(ctx context.Context, app *infra.Deps, s Slot) error {
}
