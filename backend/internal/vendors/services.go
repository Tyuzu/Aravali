package vendors

import (
	"context"
	"regexp"
	"strings"
	"time"

	"scav/infra"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func vendorBaseFilter() map[string]any {
	return map[string]any{
		"available": true,
	}
}

// RegisterVendor creates a new vendor profile.
func RegisterVendor(
	ctx context.Context,
	app *infra.Deps,
	userID,
	name,
	category,
	description,
	email,
	phone,
	location string,
) (*Vendor, error) {
	existing, _ := GetVendorByUserID(ctx, app, userID)
	if existing != nil {
		return nil, ErrVendorAlreadyExists
	}

	now := time.Now()

	vendor := &Vendor{
		VendorID:    primitive.NewObjectID().Hex(),
		UserID:      userID,
		Name:        name,
		Category:    category,
		Description: description,
		Email:       email,
		Phone:       phone,
		Location:    location,
		Rating:      0,
		RatingCount: 0,
		Verified:    false,
		Available:   true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := app.DB.InsertOne(ctx, vendorCollection, vendor); err != nil {
		return nil, err
	}

	return vendor, nil
}

// GetVendorByID retrieves a vendor by vendor ID.
func GetVendorByID(ctx context.Context, app *infra.Deps, vendorID string) (*Vendor, error) {
	var vendor Vendor
	err := app.DB.FindOne(
		ctx,
		vendorCollection,
		map[string]any{
			"vendorid":  vendorID,
			"available": true,
		},
		&vendor,
	)
	if err != nil {
		return nil, ErrVendorNotFound
	}

	return &vendor, nil
}

// GetVendorByUserID retrieves the active vendor profile for a specific user.
func GetVendorByUserID(ctx context.Context, app *infra.Deps, userID string) (*Vendor, error) {
	var vendor Vendor
	err := app.DB.FindOne(
		ctx,
		vendorCollection,
		map[string]any{
			"userid":    userID,
			"available": true,
		},
		&vendor,
	)
	if err != nil {
		return nil, nil
	}

	return &vendor, nil
}

// GetVendorsByCategory retrieves all vendors in a specific category.
func GetVendorsByCategory(ctx context.Context, app *infra.Deps, category string) ([]Vendor, error) {
	var vendors []Vendor
	err := app.DB.FindMany(ctx, vendorCollection, map[string]any{
		"available": true,
		"category":  category,
	}, &vendors)
	if err != nil {
		return nil, err
	}

	if vendors == nil {
		vendors = []Vendor{}
	}

	return vendors, nil
}

// GetAllVendors retrieves all available vendors, optionally filtered by search/category.
func GetAllVendors(ctx context.Context, app *infra.Deps, search string, category string) ([]Vendor, error) {
	filter := vendorBaseFilter()

	if category != "" {
		filter["category"] = category
	}

	if search != "" {
		escaped := regexp.QuoteMeta(strings.TrimSpace(search))
		filter["$or"] = []map[string]any{
			{"name": map[string]any{"$regex": escaped, "$options": "i"}},
			{"category": map[string]any{"$regex": escaped, "$options": "i"}},
			{"description": map[string]any{"$regex": escaped, "$options": "i"}},
			{"location": map[string]any{"$regex": escaped, "$options": "i"}},
		}
	}

	var vendors []Vendor
	err := app.DB.FindMany(ctx, vendorCollection, filter, &vendors)
	if err != nil {
		return nil, err
	}

	if vendors == nil {
		vendors = []Vendor{}
	}

	return vendors, nil
}

// UpdateVendor updates vendor information.
func UpdateVendor(ctx context.Context, app *infra.Deps, vendorID string, updates map[string]any) (any, error) {
	if updates == nil {
		updates = map[string]any{}
	}

	updates["updated_at"] = time.Now()

	return app.DB.Update(
		ctx,
		vendorCollection,
		map[string]any{"vendorid": vendorID, "available": true},
		map[string]any{"$set": updates},
	)
}

// DeleteVendor soft-deletes a vendor by setting available to false.
func DeleteVendor(ctx context.Context, app *infra.Deps, vendorID string) (any, error) {
	return app.DB.Update(
		ctx,
		vendorCollection,
		map[string]any{"vendorid": vendorID},
		map[string]any{
			"$set": map[string]any{
				"available":  false,
				"updated_at": time.Now(),
			},
		},
	)
}

// GetVendorHiringByID retrieves a hiring record by hiring ID.
func GetVendorHiringByID(ctx context.Context, app *infra.Deps, hiringID string) (*VendorHiring, error) {
	var hiring VendorHiring
	err := app.DB.FindOne(
		ctx,
		hiringCollection,
		map[string]any{"hiringid": hiringID},
		&hiring,
	)
	if err != nil {
		return nil, ErrVendorNotFound
	}

	return &hiring, nil
}

// GetVendorHiringByEventAndVendor retrieves a hiring record for a specific event/vendor pair.
func GetVendorHiringByEventAndVendor(ctx context.Context, app *infra.Deps, eventID, vendorID string) (*VendorHiring, error) {
	var hiring VendorHiring
	err := app.DB.FindOne(
		ctx,
		hiringCollection,
		map[string]any{
			"eventid":  eventID,
			"vendorid": vendorID,
			"status":   map[string]any{"$ne": "rejected"},
		},
		&hiring,
	)
	if err != nil {
		return nil, ErrVendorNotInEvent
	}

	return &hiring, nil
}

// HireVendor creates a vendor hiring record for an event.
func HireVendor(ctx context.Context, app *infra.Deps, eventID, vendorID, vendorName, vendorCategory, hiredBy string) (*VendorHiring, error) {
	existing, err := GetVendorHiringByEventAndVendor(ctx, app, eventID, vendorID)
	if err == nil && existing != nil {
		return nil, ErrVendorAlreadyHired
	}

	now := time.Now()

	hiring := &VendorHiring{
		HiringID:       primitive.NewObjectID().Hex(),
		EventID:        eventID,
		VendorID:       vendorID,
		VendorName:     vendorName,
		VendorCategory: vendorCategory,
		HiredBy:        hiredBy,
		Status:         "pending",
		HiredAt:        now,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := app.DB.InsertOne(ctx, hiringCollection, hiring); err != nil {
		return nil, err
	}

	return hiring, nil
}

// GetEventVendors retrieves all vendors hired for an event.
func GetEventVendors(ctx context.Context, app *infra.Deps, eventID string) ([]VendorHiring, error) {
	var hirings []VendorHiring
	err := app.DB.FindMany(ctx, hiringCollection, map[string]any{
		"eventid": eventID,
		"status":  map[string]any{"$ne": "rejected"},
	}, &hirings)
	if err != nil {
		return nil, err
	}

	if hirings == nil {
		hirings = []VendorHiring{}
	}

	return hirings, nil
}

// GetVendorHiringsByVendorID retrieves vendor hiring records for a specific vendor.
func GetVendorHiringsByVendorID(ctx context.Context, app *infra.Deps, vendorID string) ([]VendorHiring, error) {
	var hirings []VendorHiring
	err := app.DB.FindMany(ctx, hiringCollection, map[string]any{
		"vendorid": vendorID,
		"status":   map[string]any{"$ne": "rejected"},
	}, &hirings)
	if err != nil {
		return nil, err
	}

	if hirings == nil {
		hirings = []VendorHiring{}
	}

	return hirings, nil
}

// RemoveVendorFromEvent removes a vendor from an event.
func RemoveVendorFromEvent(ctx context.Context, app *infra.Deps, eventID, vendorID string) (any, error) {
	var existing VendorHiring
	err := app.DB.FindOne(ctx, hiringCollection, map[string]any{
		"eventid":  eventID,
		"vendorid": vendorID,
		"status":   map[string]any{"$ne": "rejected"},
	}, &existing)
	if err != nil {
		return nil, ErrVendorNotInEvent
	}

	return app.DB.Update(
		ctx,
		hiringCollection,
		map[string]any{
			"eventid":  eventID,
			"vendorid": vendorID,
		},
		map[string]any{
			"$set": map[string]any{
				"status":     "rejected",
				"updated_at": time.Now(),
			},
		},
	)
}

// UpdateVendorStatus updates the status of a vendor hiring.
func UpdateVendorStatus(ctx context.Context, app *infra.Deps, hiringID, status string) (any, error) {
	return app.DB.Update(
		ctx,
		hiringCollection,
		map[string]any{"hiringid": hiringID},
		map[string]any{
			"$set": map[string]any{
				"status":     status,
				"updated_at": time.Now(),
			},
		},
	)
}

// GetVendorsByEvent retrieves detailed vendor info for an event.
func GetVendorsByEvent(ctx context.Context, app *infra.Deps, eventID string) ([]VendorResponse, error) {
	hirings, err := GetEventVendors(ctx, app, eventID)
	if err != nil {
		return nil, err
	}

	responses := make([]VendorResponse, 0, len(hirings))

	for _, hiring := range hirings {
		vendor, err := GetVendorByID(ctx, app, hiring.VendorID)
		if err != nil || vendor == nil {
			continue
		}

		responses = append(responses, VendorResponse{
			VendorID:     vendor.VendorID,
			Name:         vendor.Name,
			Category:     vendor.Category,
			Description:  vendor.Description,
			Email:        vendor.Email,
			Phone:        vendor.Phone,
			Location:     vendor.Location,
			Rating:       vendor.Rating,
			RatingCount:  vendor.RatingCount,
			ProfileImage: vendor.ProfileImage,
			Portfolio:    vendor.Portfolio,
			Verified:     vendor.Verified,
			Status:       hiring.Status,
			HiringID:     hiring.HiringID,
			HiredAt:      hiring.HiredAt,
		})
	}

	return responses, nil
}
