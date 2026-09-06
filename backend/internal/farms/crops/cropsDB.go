package crops

import (
	"context"
	"scav/config"
	"scav/infra/db"
	"scav/internal/farms"
)

var (
	cropsCollection      = config.Collections.CropsCollection
	cropsAboutCollection = config.Collections.CropsAboutCollection
	catalogueCollection  = config.Collections.CatalogueCollection
)

func insertCrop(ctx context.Context, database db.Database, crop farms.Crop) error {
	return database.InsertOne(ctx, cropsCollection, crop)
}

func updateCrop(ctx context.Context, database db.Database, cropID string, update map[string]any) error {
	_, err := database.UpdateOne(ctx, cropsCollection, map[string]any{"cropid": cropID}, map[string]any{"$set": update})
	return err
}

func findFilteredCrops(ctx context.Context, database db.Database, filter map[string]any) ([]farms.Crop, error) {
	var crops []farms.Crop
	if err := database.FindMany(ctx, cropsCollection, filter, &crops); err != nil {
		return nil, err
	}
	if crops == nil {
		return []farms.Crop{}, nil
	}
	return crops, nil
}

func findCatalogueItems(ctx context.Context, database db.Database, filter map[string]any) ([]farms.CropCatalogueItem, error) {
	var items []farms.CropCatalogueItem
	if err := database.FindMany(ctx, catalogueCollection, filter, &items); err != nil {
		return nil, err
	}
	if items == nil {
		return []farms.CropCatalogueItem{}, nil
	}
	return items, nil
}

func getAllCrops(ctx context.Context, database db.Database) ([]farms.Crop, error) {
	var crops []farms.Crop
	if err := database.FindMany(ctx, cropsCollection, map[string]any{}, &crops); err != nil {
		return nil, err
	}
	if crops == nil {
		return []farms.Crop{}, nil
	}
	return crops, nil
}

func createCropAbout(ctx context.Context, database db.Database, crop *CropAbout) error {
	return database.InsertOne(ctx, cropsAboutCollection, crop)
}

func getCropAboutByID(ctx context.Context, database db.Database, cropID string) (*CropAbout, error) {
	var crop CropAbout
	if err := database.FindOne(ctx, cropsAboutCollection, map[string]any{"id": cropID}, &crop); err != nil {
		return nil, err
	}
	return &crop, nil
}

func getAllCropAbouts(ctx context.Context, database db.Database) ([]CropAbout, error) {
	var crops []CropAbout
	if err := database.FindMany(ctx, cropsAboutCollection, map[string]any{}, &crops); err != nil {
		return nil, err
	}
	if crops == nil {
		return []CropAbout{}, nil
	}
	return crops, nil
}

func updateCropAbout(ctx context.Context, database db.Database, cropID string, crop *CropAbout) (any, error) {
	return database.UpdateOne(ctx, cropsAboutCollection, map[string]any{"id": cropID}, map[string]any{"$set": crop})
}

func deleteCropAbout(ctx context.Context, database db.Database, cropID string) error {
	_, err := database.DeleteOne(ctx, cropsAboutCollection, map[string]any{"id": cropID})
	return err
}
