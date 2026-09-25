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
	return app.DB.InsertOne(ctx, vendorTable, vendor)
}

// FindVendorByID returns a vendor by vendorID. Returns ErrVendorNotFound if not found.
func SQLFindVendorByID(ctx context.Context, app *infra.Deps, vendorID string) (*Vendor, error) {
	var vendor Vendor
	err := app.DB.FindOne(ctx, vendorTable, map[string]any{"vendorid": vendorID, "available": true}, &vendor)
	if err != nil {
		return nil, ErrVendorNotFound
	}
	return &vendor, nil
}

// FindVendorByUserID returns a vendor by userID. Returns nil,err when FindOne fails.
func SQLFindVendorByUserID(ctx context.Context, app *infra.Deps, userID string) (*Vendor, error) {
	var vendor Vendor
	err := app.DB.FindOne(ctx, vendorTable, map[string]any{"userid": userID, "available": true}, &vendor)
	if err != nil {
		return nil, err
	}
	return &vendor, nil
}

// FindVendors finds many vendors using provided filter.
func SQLFindVendors(ctx context.Context, app *infra.Deps, filter map[string]any, out *[]Vendor) error {
	return app.DB.FindMany(ctx, vendorTable, filter, out)
}

// UpdateVendorDB updates vendor documents matching filter with update doc.
func SQLUpdateVendorDB(ctx context.Context, app *infra.Deps, filter map[string]any, update map[string]any) (any, error) {
	return app.DB.Update(ctx, vendorTable, filter, update)
}

// DeleteVendorDB marks a vendor as unavailable.
func SQLDeleteVendorDB(ctx context.Context, app *infra.Deps, vendorID string) (any, error) {
	return app.DB.Update(ctx, vendorTable, map[string]any{"vendorid": vendorID}, map[string]any{"$set": map[string]any{"available": false, "updated_at": time.Now()}})
}

// --- Hiring related DB helpers ---
func SQLFindHiringByID(ctx context.Context, app *infra.Deps, hiringID string) (*VendorHiring, error) {
	var h VendorHiring
	err := app.DB.FindOne(ctx, hiringTable, map[string]any{"hiringid": hiringID}, &h)
	if err != nil {
		return nil, ErrVendorNotFound
	}
	return &h, nil
}

func SQLFindHiringByEventAndVendor(ctx context.Context, app *infra.Deps, eventID, vendorID string) (*VendorHiring, error) {
	var h VendorHiring
	err := app.DB.FindOne(ctx, hiringTable, map[string]any{"eventid": eventID, "vendorid": vendorID, "status": map[string]any{"$ne": "rejected"}}, &h)
	if err != nil {
		return nil, ErrVendorNotInEvent
	}
	return &h, nil
}

func SQLInsertHiring(ctx context.Context, app *infra.Deps, hiring *VendorHiring) error {
	return app.DB.InsertOne(ctx, hiringTable, hiring)
}

func SQLFindHiringsByEvent(ctx context.Context, app *infra.Deps, eventID string, out *[]VendorHiring) error {
	return app.DB.FindMany(ctx, hiringTable, map[string]any{"eventid": eventID, "status": map[string]any{"$ne": "rejected"}}, out)
}

func SQLFindHiringsByVendorID(ctx context.Context, app *infra.Deps, vendorID string, out *[]VendorHiring) error {
	return app.DB.FindMany(ctx, hiringTable, map[string]any{"vendorid": vendorID, "status": map[string]any{"$ne": "rejected"}}, out)
}

func SQLUpdateHiringDB(ctx context.Context, app *infra.Deps, filter map[string]any, update map[string]any) (any, error) {
	return app.DB.Update(ctx, hiringTable, filter, update)
}

// --- Availability related DB helpers ---
func SQLFindAvailabilitySlots(ctx context.Context, app *infra.Deps, vendorID string) ([]AvailabilitySlot, error) {
	var slots []AvailabilitySlot
	err := app.DB.FindMany(ctx, config.Tables.VendorAvailabilityTable, map[string]any{"vendorid": vendorID}, &slots)
	if err != nil {
		return nil, err
	}
	if slots == nil {
		slots = []AvailabilitySlot{}
	}
	return slots, nil
}

func SQLInsertAvailabilitySlotDB(ctx context.Context, app *infra.Deps, slot AvailabilitySlot) error {
	return app.DB.InsertOne(ctx, config.Tables.VendorAvailabilityTable, slot)
}

func SQLFindAvailabilitySlotByID(ctx context.Context, app *infra.Deps, slotID, vendorID string) (*AvailabilitySlot, error) {
	var slot AvailabilitySlot
	err := app.DB.FindOne(ctx, config.Tables.VendorAvailabilityTable, map[string]any{"slotid": slotID, "vendorid": vendorID}, &slot)
	if err != nil {
		return nil, err
	}
	return &slot, nil
}

func SQLDeleteAvailabilitySlotDB(ctx context.Context, app *infra.Deps, slotID string) (any, error) {
	return app.DB.DeleteOne(ctx, config.Tables.VendorAvailabilityTable, map[string]any{"slotid": slotID})
}
