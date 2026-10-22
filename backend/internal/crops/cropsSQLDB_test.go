package crops

import (
	"testing"
	"time"

	"scav/internal/farms"
)

func TestCropMetadataIncludesCoreFields(t *testing.T) {
	harvest := time.Date(2026, time.October, 7, 0, 0, 0, 0, time.UTC)
	expiry := time.Date(2026, time.November, 6, 0, 0, 0, 0, time.UTC)

	crop := farms.Crop{
		Name:        "spinach",
		CropId:      "ztJZ1EVh65ONh",
		FarmID:      "HjAYMSbe37qmXR",
		FarmName:    "HjAYMSbe37qmXR",
		Category:    "vegetables",
		Price:       567,
		Discount:    5,
		Quantity:    5,
		Unit:        "kg",
		Notes:       "56",
		CreatedBy:   "user-123",
		HarvestDate: &harvest,
		ExpiryDate:  &expiry,
	}

	meta := QQQcropMetadataMap(crop)

	if meta["name"] != "spinach" {
		t.Fatalf("expected name to be preserved, got %v", meta["name"])
	}
	if meta["farmid"] != "HjAYMSbe37qmXR" {
		t.Fatalf("expected farmid to be preserved, got %v", meta["farmid"])
	}
	if meta["category"] != "vegetables" {
		t.Fatalf("expected category to be preserved, got %v", meta["category"])
	}
	if meta["price"] != 567.0 {
		t.Fatalf("expected price to be preserved, got %v", meta["price"])
	}
	if meta["unit"] != "kg" {
		t.Fatalf("expected unit to be preserved, got %v", meta["unit"])
	}
	if meta["owner"] != "user-123" {
		t.Fatalf("expected owner to be preserved, got %v", meta["owner"])
	}
}
