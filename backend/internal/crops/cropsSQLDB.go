package crops

import (
	"context"

	"scav/config"
	"scav/infra/sqldb"
	"scav/internal/farms"
)

var (
	cropsTable      = config.Tables.CropsTable
	cropsAboutTable = config.Tables.CropsAboutTable
	catalogueTable  = config.Tables.CatalogueTable
)

func SQLinsertCrop(ctx context.Context, database sqldb.Database, crop farms.Crop) error {
	return database.InsertOne(ctx, cropsTable, crop)
}

func SQLupdateCrop(ctx context.Context, database sqldb.Database, cropID string, update map[string]any) error {
	query := "cropid = $1"
	args := []any{cropID}

	_, err := database.UpdateOne(ctx, cropsTable, query, args, update)
	return err
}

func SQLfindFilteredCrops(ctx context.Context, database sqldb.Database, query string, args []any) ([]farms.Crop, error) {
	var crops []farms.Crop
	if err := database.FindMany(ctx, cropsTable, query, args, &crops); err != nil {
		return nil, err
	}
	if crops == nil {
		return []farms.Crop{}, nil
	}
	return crops, nil
}

func SQLfindCatalogueItems(ctx context.Context, database sqldb.Database, query string, args []any) ([]farms.CropCatalogueItem, error) {
	var items []farms.CropCatalogueItem
	if err := database.FindMany(ctx, catalogueTable, query, args, &items); err != nil {
		return nil, err
	}
	if items == nil {
		return []farms.CropCatalogueItem{}, nil
	}
	return items, nil
}

func SQLgetAllCrops(ctx context.Context, database sqldb.Database) ([]farms.Crop, error) {
	var crops []farms.Crop
	if err := database.FindMany(ctx, cropsTable, "", nil, &crops); err != nil {
		return nil, err
	}
	if crops == nil {
		return []farms.Crop{}, nil
	}
	return crops, nil
}

func SQLcreateCropAbout(ctx context.Context, database sqldb.Database, crop *CropAbout) error {
	return database.InsertOne(ctx, cropsAboutTable, crop)
}

func SQLgetCropAboutByID(ctx context.Context, database sqldb.Database, cropID string) (*CropAbout, error) {
	var crop CropAbout
	query := "id = $1"
	args := []any{cropID}

	if err := database.FindOne(ctx, cropsAboutTable, query, args, &crop); err != nil {
		return nil, err
	}
	return &crop, nil
}

func SQLgetAllCropAbouts(ctx context.Context, database sqldb.Database) ([]CropAbout, error) {
	var crops []CropAbout
	if err := database.FindMany(ctx, cropsAboutTable, "", nil, &crops); err != nil {
		return nil, err
	}
	if crops == nil {
		return []CropAbout{}, nil
	}
	return crops, nil
}
func SQLupdateCropAbout(ctx context.Context, database sqldb.Database, cropID string, crop *CropAbout) (int64, error) {
	if crop == nil {
		return 0, nil
	}

	query := "id = $1"
	args := []any{cropID}

	updateValues := map[string]any{
		"commonName":         crop.CommonName,
		"scientificName":     crop.ScientificName,
		"image":              crop.Image,
		"imageAlt":           crop.ImageAlt,
		"description":        crop.Description,
		"nutritionalValues":  crop.NutritionalValues,
		"growingConditions":  crop.GrowingConditions,
		"plantingHarvesting": crop.PlantingHarvesting,
		"careTips":           crop.CareTips,
		"varieties":          crop.Varieties,
		"usage":              crop.Usage,
		"funFacts":           crop.FunFacts,
	}

	return database.UpdateOne(ctx, cropsAboutTable, query, args, updateValues)
}

func SQLdeleteCropAbout(ctx context.Context, database sqldb.Database, cropID string) error {
	query := "id = $1"
	args := []any{cropID}

	_, err := database.DeleteOne(ctx, cropsAboutTable, query, args)
	return err
}
