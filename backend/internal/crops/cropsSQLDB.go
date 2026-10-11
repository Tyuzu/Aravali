// File: internal/crops/cropsSQLDB.go

package crops

import (
	"context"
	"fmt"
	"strings"

	"scav/config"
	"scav/infra"
	"scav/internal/farms"
)

var (
	cropsTable      = config.Tables.CropsTable
	cropsAboutTable = config.Tables.CropsAboutTable
	catalogueTable  = config.Tables.CatalogueTable
)

func insertCrop(ctx context.Context, app infra.Deps, crop farms.Crop) error {

}

func updateCrop(ctx context.Context, app infra.Deps, cropID string, update map[string]any) error {

}

func findFilteredCrops(ctx context.Context, app infra.Deps, query string, args []any) ([]farms.Crop, error) {
	var crops []farms.Crop
	return crops, nil
}

func findCatalogueItems(ctx context.Context, app infra.Deps, query string, args []any) ([]farms.CropCatalogueItem, error) {
	var items []farms.CropCatalogueItem

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

func getAllCrops(ctx context.Context, app infra.Deps) ([]farms.Crop, error) {
	var crops []farms.Crop

	return crops, nil
}

func createCropAbout(ctx context.Context, app infra.Deps, crop *CropAbout) error {

}

func getCropAboutByID(ctx context.Context, app infra.Deps, cropID string) (*CropAbout, error) {
	var crop CropAbout

	return &crop, nil
}

func getAllCropAbouts(ctx context.Context, app infra.Deps) ([]CropAbout, error) {
	var crops []CropAbout

	return crops, nil
}

func updateCropAbout(ctx context.Context, app infra.Deps, cropID string, crop *CropAbout) (int64, error) {

}

func deleteCropAbout(ctx context.Context, app infra.Deps, cropID string) error {

}
