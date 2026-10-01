// File: internal/crops/cropsSQLDB.go

package crops

import (
	"context"
	"fmt"
	"strings"

	"scav/config"
	"scav/infra/sqldb"
	"scav/internal/farms"
)

var (
	cropsTable      = config.Tables.CropsTable
	cropsAboutTable = config.Tables.CropsAboutTable
	catalogueTable  = config.Tables.CatalogueTable
)

func insertCrop(ctx context.Context, database sqldb.Database, crop farms.Crop) error {
	return database.InsertOne(ctx, cropsTable, crop)
}

func updateCrop(ctx context.Context, database sqldb.Database, cropID string, update map[string]any) error {
	query := "cropid = $1"
	args := []any{cropID}

	_, err := database.UpdateOne(ctx, cropsTable, query, args, update)
	return err
}

func findFilteredCrops(ctx context.Context, database sqldb.Database, query string, args []any) ([]farms.Crop, error) {
	var crops []farms.Crop
	if err := database.FindMany(ctx, cropsTable, query, args, &crops); err != nil {
		return nil, err
	}
	if crops == nil {
		return []farms.Crop{}, nil
	}
	return crops, nil
}

func findCatalogueItems(ctx context.Context, database sqldb.Database, query string, args []any) ([]farms.CropCatalogueItem, error) {
	var items []farms.CropCatalogueItem
	if err := database.FindMany(ctx, catalogueTable, query, args, &items); err != nil {
		return nil, err
	}
	if items == nil {
		return []farms.CropCatalogueItem{}, nil
	}
	return items, nil
}

func buildCropFilterQuery(filter map[string]any) (string, []any) {
	if len(filter) == 0 {
		return "1 = 1", nil
	}
	clauses := make([]string, 0, len(filter))
	args := make([]any, 0, len(filter))
	for key, value := range filter {
		if value == nil {
			continue
		}
		switch v := value.(type) {
		case map[string]any:
			if gt, ok := v["$gt"]; ok {
				clauses = append(clauses, fmt.Sprintf("%s > $%d", key, len(args)+1))
				args = append(args, gt)
				continue
			}
			if gte, ok := v["$gte"]; ok {
				clauses = append(clauses, fmt.Sprintf("%s >= $%d", key, len(args)+1))
				args = append(args, gte)
				continue
			}
			if lte, ok := v["$lte"]; ok {
				clauses = append(clauses, fmt.Sprintf("%s <= $%d", key, len(args)+1))
				args = append(args, lte)
				continue
			}
		}
		clauses = append(clauses, fmt.Sprintf("%s = $%d", key, len(args)+1))
		args = append(args, value)
	}
	if len(clauses) == 0 {
		return "1 = 1", nil
	}
	return strings.Join(clauses, " AND "), args
}

func getAllCrops(ctx context.Context, database sqldb.Database) ([]farms.Crop, error) {
	var crops []farms.Crop
	if err := database.FindMany(ctx, cropsTable, "", nil, &crops); err != nil {
		return nil, err
	}
	if crops == nil {
		return []farms.Crop{}, nil
	}
	return crops, nil
}

func createCropAbout(ctx context.Context, database sqldb.Database, crop *CropAbout) error {
	return database.InsertOne(ctx, cropsAboutTable, crop)
}

func getCropAboutByID(ctx context.Context, database sqldb.Database, cropID string) (*CropAbout, error) {
	var crop CropAbout
	query := "id = $1"
	args := []any{cropID}

	if err := database.FindOne(ctx, cropsAboutTable, query, args, &crop); err != nil {
		return nil, err
	}
	return &crop, nil
}

func getAllCropAbouts(ctx context.Context, database sqldb.Database) ([]CropAbout, error) {
	var crops []CropAbout
	if err := database.FindMany(ctx, cropsAboutTable, "", nil, &crops); err != nil {
		return nil, err
	}
	if crops == nil {
		return []CropAbout{}, nil
	}
	return crops, nil
}
func updateCropAbout(ctx context.Context, database sqldb.Database, cropID string, crop *CropAbout) (int64, error) {
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

func deleteCropAbout(ctx context.Context, database sqldb.Database, cropID string) error {
	query := "id = $1"
	args := []any{cropID}

	_, err := database.DeleteOne(ctx, cropsAboutTable, query, args)
	return err
}
