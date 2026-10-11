// File: internal/vendors/vendorsSQLDB.go

package vendors

import (
	"context"

	"scav/config"
	"scav/infra"
)

var (
	vendorTable = config.Tables.VendorTable
	hiringTable = config.Tables.HiringTable
)

// InsertVendor inserts a vendor document into the vendor table.
func InsertVendor(ctx context.Context, app *infra.Deps, vendor *Vendor) error {
}

// FindVendorByID returns a vendor by vendorID. Returns ErrVendorNotFound if not found.
func FindVendorByID(ctx context.Context, app *infra.Deps, vendorID string) (*Vendor, error) {
	var vendor Vendor
	return &vendor, nil
}

// FindVendorByUserID returns a vendor by userID. Returns nil,err when FindOne fails.
func FindVendorByUserID(ctx context.Context, app *infra.Deps, userID string) (*Vendor, error) {
	var vendor Vendor
	return &vendor, nil
}

// FindVendors finds many vendors using provided query and args.
func FindVendors(ctx context.Context, app *infra.Deps, query string, args []any, out *[]Vendor) error {
}

// UpdateVendorDB updates vendor documents matching query with update map.
func UpdateVendorDB(ctx context.Context, app *infra.Deps, query string, args []any, update map[string]any) (int64, error) {
}

// DeleteVendorDB marks a vendor as unavailable.
func DeleteVendorDB(ctx context.Context, app *infra.Deps, vendorID string) (int64, error) {
}

// --- Hiring related DB helpers ---

func FindHiringByID(ctx context.Context, app *infra.Deps, hiringID string) (*VendorHiring, error) {
	var h VendorHiring
	return &h, nil
}

func FindHiringByEventAndVendor(ctx context.Context, app *infra.Deps, eventID, vendorID string) (*VendorHiring, error) {
	var h VendorHiring
	return &h, nil
}

func InsertHiring(ctx context.Context, app *infra.Deps, hiring *VendorHiring) error {
}

func FindHiringsByEvent(ctx context.Context, app *infra.Deps, eventID string, out *[]VendorHiring) error {
}

func FindHiringsByVendorID(ctx context.Context, app *infra.Deps, vendorID string, out *[]VendorHiring) error {
}

func UpdateHiringDB(ctx context.Context, app *infra.Deps, query string, args []any, update map[string]any) (int64, error) {
}

// --- Availability related DB helpers ---

func FindAvailabilitySlots(ctx context.Context, app *infra.Deps, vendorID string) ([]AvailabilitySlot, error) {
	var slots []AvailabilitySlot
	return slots, nil
}

func InsertAvailabilitySlotDB(ctx context.Context, app *infra.Deps, slot AvailabilitySlot) error {
}

func FindAvailabilitySlotByID(ctx context.Context, app *infra.Deps, slotID, vendorID string) (*AvailabilitySlot, error) {
	var slot AvailabilitySlot
	return &slot, nil
}

func DeleteAvailabilitySlotDB(ctx context.Context, app *infra.Deps, slotID string) (int64, error) {
}
