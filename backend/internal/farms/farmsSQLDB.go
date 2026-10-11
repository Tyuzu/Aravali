// File: internal/farms/farmsSQLDB.go

package farms

import (
	"context"

	"scav/config"
	"scav/infra"
	"scav/internal/cart"
)

var (
	cropsTable      = config.Tables.CropsTable
	farmsTable      = config.Tables.FarmsTable
	farmOrdersTable = config.Tables.FarmOrdersTable
)

func insertFarm(ctx context.Context, app infra.Deps, farm Farm) error {

}

func getFarmByID(ctx context.Context, app infra.Deps, farmID string) (Farm, error) {
	var farm Farm

	return farm, err
}

func getFarmByCreatedBy(ctx context.Context, app infra.Deps, userID string) (Farm, error) {
	var farm Farm

	return farm, err
}

func getCropsByFarmID(ctx context.Context, app infra.Deps, farmID string) ([]Crop, error) {
	var crops []Crop

	return crops, nil
}

func getFarmOrdersByFarmID(ctx context.Context, app infra.Deps, farmID string) ([]cart.FarmOrder, error) {
	var orders []cart.FarmOrder

	return orders, nil
}

func getCropsByCropID(ctx context.Context, app infra.Deps, cropID string) ([]Crop, error) {
	var crops []Crop

	return crops, nil
}

func getFarmsByIDs(ctx context.Context, app infra.Deps, farmIDs []string) ([]Farm, error) {
	var farms []Farm

	return farms, nil
}

func getCropsByNameFilter(ctx context.Context, app infra.Deps, cropName string) ([]Crop, error) {

	var crops []Crop

	return crops, nil
}

func getPaginatedFarms(ctx context.Context, app infra.Deps, search string, offset, limit int) ([]Farm, int64, error) {

	var farms []Farm

	return farms, total, nil
}

func getMyFarmsPage(ctx context.Context, app infra.Deps, userID string, offset, limit int) ([]Farm, int64, error) {

	var farms []Farm

	return farms, total, nil
}

func updateOwnedFarm(ctx context.Context, app infra.Deps, farmID, userID string, update map[string]any) (int64, error) {

}

func deleteFarmByID(ctx context.Context, app infra.Deps, farmID string) (int64, error) {

}
