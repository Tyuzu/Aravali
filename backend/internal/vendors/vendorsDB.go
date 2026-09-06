package vendors

import (
	"context"
	"errors"
	"time"

	"scav/config"
	"scav/infra"
)

var (
	vendorCollection = config.Collections.VendorCollection
	hiringCollection = config.Collections.HiringCollection

	ErrVendorNotFound      = errors.New("vendor not found")
	ErrVendorAlreadyExists = errors.New("vendor profile already exists")
	ErrVendorAlreadyHired  = errors.New("vendor already hired for this event")
	ErrVendorNotInEvent    = errors.New("vendor not found for this event")
	ErrUnauthorizedVendor  = errors.New("unauthorized vendor action")
)

// DB layer helpers: centralize all direct DB calls here so other files don't
// access app.DB directly.

// InsertVendor inserts a vendor document into the vendor collection.
func InsertVendor(ctx context.Context, app *infra.Deps, vendor *Vendor) error {
	return app.DB.InsertOne(ctx, vendorCollection, vendor)
}

// FindVendorByID returns a vendor by vendorID. Returns ErrVendorNotFound if not found.
func FindVendorByID(ctx context.Context, app *infra.Deps, vendorID string) (*Vendor, error) {
	var vendor Vendor
	err := app.DB.FindOne(ctx, vendorCollection, map[string]any{"vendorid": vendorID, "available": true}, &vendor)
	if err != nil {
		return nil, ErrVendorNotFound
	}
	return &vendor, nil
}

// FindVendorByUserID returns a vendor by userID. Returns nil,err when FindOne fails.
func FindVendorByUserID(ctx context.Context, app *infra.Deps, userID string) (*Vendor, error) {
	var vendor Vendor
	err := app.DB.FindOne(ctx, vendorCollection, map[string]any{"userid": userID, "available": true}, &vendor)
	if err != nil {
		return nil, err
	}
	return &vendor, nil
}

// FindVendors finds many vendors using provided filter.
func FindVendors(ctx context.Context, app *infra.Deps, filter map[string]any, out *[]Vendor) error {
	return app.DB.FindMany(ctx, vendorCollection, filter, out)
}

// UpdateVendorDB updates vendor documents matching filter with update doc.
func UpdateVendorDB(ctx context.Context, app *infra.Deps, filter map[string]any, update map[string]any) (any, error) {
	return app.DB.Update(ctx, vendorCollection, filter, update)
}

// DeleteVendorDB marks a vendor as unavailable.
func DeleteVendorDB(ctx context.Context, app *infra.Deps, vendorID string) (any, error) {
	return app.DB.Update(ctx, vendorCollection, map[string]any{"vendorid": vendorID}, map[string]any{"$set": map[string]any{"available": false, "updated_at": time.Now()}})
}

// --- Hiring related DB helpers ---
func FindHiringByID(ctx context.Context, app *infra.Deps, hiringID string) (*VendorHiring, error) {
	var h VendorHiring
	err := app.DB.FindOne(ctx, hiringCollection, map[string]any{"hiringid": hiringID}, &h)
	if err != nil {
		return nil, ErrVendorNotFound
	}
	return &h, nil
}

func FindHiringByEventAndVendor(ctx context.Context, app *infra.Deps, eventID, vendorID string) (*VendorHiring, error) {
	var h VendorHiring
	err := app.DB.FindOne(ctx, hiringCollection, map[string]any{"eventid": eventID, "vendorid": vendorID, "status": map[string]any{"$ne": "rejected"}}, &h)
	if err != nil {
		return nil, ErrVendorNotInEvent
	}
	return &h, nil
}

func InsertHiring(ctx context.Context, app *infra.Deps, hiring *VendorHiring) error {
	return app.DB.InsertOne(ctx, hiringCollection, hiring)
}

func FindHiringsByEvent(ctx context.Context, app *infra.Deps, eventID string, out *[]VendorHiring) error {
	return app.DB.FindMany(ctx, hiringCollection, map[string]any{"eventid": eventID, "status": map[string]any{"$ne": "rejected"}}, out)
}

func FindHiringsByVendorID(ctx context.Context, app *infra.Deps, vendorID string, out *[]VendorHiring) error {
	return app.DB.FindMany(ctx, hiringCollection, map[string]any{"vendorid": vendorID, "status": map[string]any{"$ne": "rejected"}}, out)
}

func UpdateHiringDB(ctx context.Context, app *infra.Deps, filter map[string]any, update map[string]any) (any, error) {
	return app.DB.Update(ctx, hiringCollection, filter, update)
}

// --- Availability related DB helpers ---
func FindAvailabilitySlots(ctx context.Context, app *infra.Deps, vendorID string) ([]AvailabilitySlot, error) {
	var slots []AvailabilitySlot
	err := app.DB.FindMany(ctx, config.Collections.VendorAvailabilityCollection, map[string]any{"vendorid": vendorID}, &slots)
	if err != nil {
		return nil, err
	}
	if slots == nil {
		slots = []AvailabilitySlot{}
	}
	return slots, nil
}

func InsertAvailabilitySlotDB(ctx context.Context, app *infra.Deps, slot AvailabilitySlot) error {
	return app.DB.InsertOne(ctx, config.Collections.VendorAvailabilityCollection, slot)
}

func FindAvailabilitySlotByID(ctx context.Context, app *infra.Deps, slotID, vendorID string) (*AvailabilitySlot, error) {
	var slot AvailabilitySlot
	err := app.DB.FindOne(ctx, config.Collections.VendorAvailabilityCollection, map[string]any{"slotid": slotID, "vendorid": vendorID}, &slot)
	if err != nil {
		return nil, err
	}
	return &slot, nil
}

func DeleteAvailabilitySlotDB(ctx context.Context, app *infra.Deps, slotID string) (any, error) {
	return app.DB.DeleteOne(ctx, config.Collections.VendorAvailabilityCollection, map[string]any{"slotid": slotID})
}
