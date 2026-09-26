package vendors

import (
	"context"
	"time"

	"scav/config"
	"scav/infra"
)

var (
	vendorTable = config.Tables.VendorTable
	hiringTable = config.Tables.HiringTable
)

// DB layer helpers: centralize all direct DB calls here so other files don't
// access app.DB directly.

// InsertVendor inserts a vendor document into the vendor table.
func SQLInsertVendor(ctx context.Context, app *infra.Deps, vendor *Vendor) error {
	return app.SQLDB.InsertOne(ctx, vendorTable, vendor)
}

// FindVendorByID returns a vendor by vendorID. Returns ErrVendorNotFound if not found.
func SQLFindVendorByID(ctx context.Context, app *infra.Deps, vendorID string) (*Vendor, error) {
	var vendor Vendor
	query := "vendorid = $1 AND available = $2"
	args := []any{vendorID, true}

	err := app.SQLDB.FindOne(ctx, vendorTable, query, args, &vendor)
	if err != nil {
		return nil, ErrVendorNotFound
	}
	return &vendor, nil
}

// FindVendorByUserID returns a vendor by userID. Returns nil,err when FindOne fails.
func SQLFindVendorByUserID(ctx context.Context, app *infra.Deps, userID string) (*Vendor, error) {
	var vendor Vendor
	query := "userid = $1 AND available = $2"
	args := []any{userID, true}

	err := app.SQLDB.FindOne(ctx, vendorTable, query, args, &vendor)
	if err != nil {
		return nil, err
	}
	return &vendor, nil
}

// FindVendors finds many vendors using provided query and args.
func SQLFindVendors(ctx context.Context, app *infra.Deps, query string, args []any, out *[]Vendor) error {
	return app.SQLDB.FindMany(ctx, vendorTable, query, args, out)
}

// UpdateVendorDB updates vendor documents matching query with update map.
func SQLUpdateVendorDB(ctx context.Context, app *infra.Deps, query string, args []any, update map[string]any) (int64, error) {
	return app.SQLDB.Update(ctx, vendorTable, query, args, update)
}

// DeleteVendorDB marks a vendor as unavailable.
func SQLDeleteVendorDB(ctx context.Context, app *infra.Deps, vendorID string) (int64, error) {
	query := "vendorid = $1"
	args := []any{vendorID}

	update := map[string]any{
		"available":  false,
		"updated_at": time.Now(),
	}

	return app.SQLDB.Update(ctx, vendorTable, query, args, update)
}

// --- Hiring related DB helpers ---

func SQLFindHiringByID(ctx context.Context, app *infra.Deps, hiringID string) (*VendorHiring, error) {
	var h VendorHiring
	query := "hiringid = $1"
	args := []any{hiringID}

	err := app.SQLDB.FindOne(ctx, hiringTable, query, args, &h)
	if err != nil {
		return nil, ErrVendorNotFound
	}
	return &h, nil
}

func SQLFindHiringByEventAndVendor(ctx context.Context, app *infra.Deps, eventID, vendorID string) (*VendorHiring, error) {
	var h VendorHiring
	query := "eventid = $1 AND vendorid = $2 AND status != $3"
	args := []any{eventID, vendorID, "rejected"}

	err := app.SQLDB.FindOne(ctx, hiringTable, query, args, &h)
	if err != nil {
		return nil, ErrVendorNotInEvent
	}
	return &h, nil
}

func SQLInsertHiring(ctx context.Context, app *infra.Deps, hiring *VendorHiring) error {
	return app.SQLDB.InsertOne(ctx, hiringTable, hiring)
}

func SQLFindHiringsByEvent(ctx context.Context, app *infra.Deps, eventID string, out *[]VendorHiring) error {
	query := "eventid = $1 AND status != $2"
	args := []any{eventID, "rejected"}

	return app.SQLDB.FindMany(ctx, hiringTable, query, args, out)
}

func SQLFindHiringsByVendorID(ctx context.Context, app *infra.Deps, vendorID string, out *[]VendorHiring) error {
	query := "vendorid = $1 AND status != $2"
	args := []any{vendorID, "rejected"}

	return app.SQLDB.FindMany(ctx, hiringTable, query, args, out)
}

func SQLUpdateHiringDB(ctx context.Context, app *infra.Deps, query string, args []any, update map[string]any) (int64, error) {
	return app.SQLDB.Update(ctx, hiringTable, query, args, update)
}

// --- Availability related DB helpers ---

func SQLFindAvailabilitySlots(ctx context.Context, app *infra.Deps, vendorID string) ([]AvailabilitySlot, error) {
	var slots []AvailabilitySlot
	query := "vendorid = $1"
	args := []any{vendorID}

	err := app.SQLDB.FindMany(ctx, config.Tables.VendorAvailabilityTable, query, args, &slots)
	if err != nil {
		return nil, err
	}
	if slots == nil {
		slots = []AvailabilitySlot{}
	}
	return slots, nil
}

func SQLInsertAvailabilitySlotDB(ctx context.Context, app *infra.Deps, slot AvailabilitySlot) error {
	return app.SQLDB.InsertOne(ctx, config.Tables.VendorAvailabilityTable, slot)
}

func SQLFindAvailabilitySlotByID(ctx context.Context, app *infra.Deps, slotID, vendorID string) (*AvailabilitySlot, error) {
	var slot AvailabilitySlot
	query := "slotid = $1 AND vendorid = $2"
	args := []any{slotID, vendorID}

	err := app.SQLDB.FindOne(ctx, config.Tables.VendorAvailabilityTable, query, args, &slot)
	if err != nil {
		return nil, err
	}
	return &slot, nil
}

func SQLDeleteAvailabilitySlotDB(ctx context.Context, app *infra.Deps, slotID string) (int64, error) {
	query := "slotid = $1"
	args := []any{slotID}

	return app.SQLDB.DeleteOne(ctx, config.Tables.VendorAvailabilityTable, query, args)
}
