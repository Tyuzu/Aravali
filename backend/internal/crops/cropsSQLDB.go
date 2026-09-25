package crops

import (
	"context"
	"scav/config"
	"scav/infra/db"
	"scav/internal/farms"
)

var (
	cropsTable      = config.Tables.CropsTable
	cropsAboutTable = config.Tables.CropsAboutTable
	catalogueTable  = config.Tables.CatalogueTable
)

func SQLinsertCrop(ctx context.Context, database db.Database, crop farms.Crop) error {
	return database.InsertOne(ctx, cropsTable, crop)
}

func SQLupdateCrop(ctx context.Context, database db.Database, cropID string, update map[string]any) error {
	_, err := database.UpdateOne(ctx, cropsTable, map[string]any{"cropid": cropID}, map[string]any{"$set": update})
	return err
}

func SQLfindFilteredCrops(ctx context.Context, database db.Database, filter map[string]any) ([]farms.Crop, error) {
	var crops []farms.Crop
	if err := database.FindMany(ctx, cropsTable, filter, &crops); err != nil {
		return nil, err
	}
	if crops == nil {
		return []farms.Crop{}, nil
	}
	return crops, nil
}

func SQLfindCatalogueItems(ctx context.Context, database db.Database, filter map[string]any) ([]farms.CropCatalogueItem, error) {
	var items []farms.CropCatalogueItem
	if err := database.FindMany(ctx, catalogueTable, filter, &items); err != nil {
		return nil, err
	}
	if items == nil {
		return []farms.CropCatalogueItem{}, nil
	}
	return items, nil
}

func SQLgetAllCrops(ctx context.Context, database db.Database) ([]farms.Crop, error) {
	var crops []farms.Crop
	if err := database.FindMany(ctx, cropsTable, map[string]any{}, &crops); err != nil {
		return nil, err
	}
	if crops == nil {
		return []farms.Crop{}, nil
	}
	return crops, nil
}

func SQLcreateCropAbout(ctx context.Context, database db.Database, crop *CropAbout) error {
	return database.InsertOne(ctx, cropsAboutTable, crop)
}

func SQLgetCropAboutByID(ctx context.Context, database db.Database, cropID string) (*CropAbout, error) {
	var crop CropAbout
	if err := database.FindOne(ctx, cropsAboutTable, map[string]any{"id": cropID}, &crop); err != nil {
		return nil, err
	}
	return &crop, nil
}

func SQLgetAllCropAbouts(ctx context.Context, database db.Database) ([]CropAbout, error) {
	var crops []CropAbout
	if err := database.FindMany(ctx, cropsAboutTable, map[string]any{}, &crops); err != nil {
		return nil, err
	}
	if crops == nil {
		return []CropAbout{}, nil
	}
	return crops, nil
}

func SQLupdateCropAbout(ctx context.Context, database db.Database, cropID string, crop *CropAbout) (any, error) {
	return database.UpdateOne(ctx, cropsAboutTable, map[string]any{"id": cropID}, map[string]any{"$set": crop})
}

func SQLdeleteCropAbout(ctx context.Context, database db.Database, cropID string) error {
	_, err := database.DeleteOne(ctx, cropsAboutTable, map[string]any{"id": cropID})
	return err
}
