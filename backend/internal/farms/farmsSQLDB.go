package farms

import (
	"context"
	"regexp"

	"scav/config"
	"scav/infra/db"
	"scav/internal/cart"
	"scav/utils"
)

var (
	cropsTable      = config.Tables.CropsTable
	farmsTable      = config.Tables.FarmsTable
	farmOrdersTable = config.Tables.FarmOrdersTable
	// usersTable and productsTable removed because they were unused
)

func SQLinsertFarm(ctx context.Context, database db.Database, farm Farm) error {
	return database.InsertOne(ctx, farmsTable, farm)
}

func SQLgetFarmByID(ctx context.Context, database db.Database, farmID string) (Farm, error) {
	var farm Farm
	err := database.FindOne(ctx, farmsTable, map[string]any{"farmid": farmID}, &farm)
	return farm, err
}

func SQLgetFarmByCreatedBy(ctx context.Context, database db.Database, userID string) (Farm, error) {
	var farm Farm
	err := database.FindOne(ctx, farmsTable, map[string]any{"createdBy": userID}, &farm)
	return farm, err
}

func SQLgetCropsByFarmID(ctx context.Context, database db.Database, farmID string) ([]Crop, error) {
	var crops []Crop
	if err := database.FindMany(ctx, cropsTable, map[string]any{"farmid": farmID}, &crops); err != nil {
		return nil, err
	}
	if crops == nil {
		return []Crop{}, nil
	}
	return crops, nil
}

func SQLgetFarmOrdersByFarmID(ctx context.Context, database db.Database, farmID string) ([]cart.FarmOrder, error) {
	var orders []cart.FarmOrder
	if err := database.FindMany(ctx, farmOrdersTable, map[string]any{"farmid": farmID}, &orders); err != nil {
		return nil, err
	}
	if orders == nil {
		return []cart.FarmOrder{}, nil
	}
	return orders, nil
}

func SQLgetCropsByCropID(ctx context.Context, database db.Database, cropID string) ([]Crop, error) {
	var crops []Crop
	if err := database.FindMany(ctx, cropsTable, map[string]any{"cropid": cropID}, &crops); err != nil {
		return nil, err
	}
	if crops == nil {
		return []Crop{}, nil
	}
	return crops, nil
}

func SQLgetFarmsByIDs(ctx context.Context, database db.Database, farmIDs []string) ([]Farm, error) {
	var farms []Farm
	if err := database.FindMany(ctx, farmsTable, map[string]any{"farmid": map[string]any{"$in": farmIDs}}, &farms); err != nil {
		return nil, err
	}
	if farms == nil {
		return []Farm{}, nil
	}
	return farms, nil
}

func SQLgetCropsByNameFilter(ctx context.Context, database db.Database, cropName string) ([]Crop, error) {
	filter := map[string]any{
		"name": map[string]any{
			"$regex":   "^" + regexp.QuoteMeta(cropName) + "$",
			"$options": "i",
		},
	}
	var crops []Crop
	if err := database.FindMany(ctx, cropsTable, filter, &crops); err != nil {
		return nil, err
	}
	if crops == nil {
		return []Crop{}, nil
	}
	return crops, nil
}

func SQLgetPaginatedFarms(ctx context.Context, database db.Database, search string, skip, limit int) ([]Farm, int64, error) {
	pipeline := make([]any, 0)
	if search != "" {
		pipeline = append(pipeline, map[string]any{
			"$match": map[string]any{
				"$or": []map[string]any{
					utils.RegexFilter("name", search),
					utils.RegexFilter("location", search),
					utils.RegexFilter("owner", search),
				},
			},
		})
	}
	pipeline = append(
		pipeline,
		map[string]any{"$sort": map[string]any{"createdAt": -1}},
		map[string]any{"$lookup": map[string]any{
			"from":         "crops",
			"localField":   "farmid",
			"foreignField": "farmid",
			"as":           "crops",
		}},
		map[string]any{"$skip": skip},
		map[string]any{"$limit": limit},
	)

	var farms []Farm
	if err := database.Aggregate(ctx, farmsTable, pipeline, &farms); err != nil {
		return nil, 0, err
	}

	total, err := database.CountDocuments(ctx, farmsTable, map[string]any{})
	if err != nil {
		if search == "" {
			return farms, 0, nil
		}
		return farms, 0, err
	}
	return farms, total, nil
}

func SQLgetMyFarmsPage(ctx context.Context, database db.Database, userID string, skip, limit int) ([]Farm, int64, error) {
	pipeline := []any{
		map[string]any{"$match": map[string]any{"createdBy": userID}},
		map[string]any{"$sort": map[string]any{"createdAt": -1}},
		map[string]any{"$lookup": map[string]any{
			"from":         "crops",
			"localField":   "farmid",
			"foreignField": "farmid",
			"as":           "crops",
		}},
		map[string]any{"$skip": skip},
		map[string]any{"$limit": limit},
	}

	var farms []Farm
	if err := database.Aggregate(ctx, farmsTable, pipeline, &farms); err != nil {
		return nil, 0, err
	}

	total, err := database.CountDocuments(ctx, farmsTable, map[string]any{"createdBy": userID})
	if err != nil {
		return farms, 0, nil
	}
	return farms, total, nil
}

func SQLupdateOwnedFarm(ctx context.Context, database db.Database, farmID, userID string, update any) (any, error) {
	// The owner field on Farm is stored as "createdBy" (see farmModels.go).
	// Use that field to ensure the update only affects farms owned by the user.
	return database.UpdateOne(ctx, farmsTable, map[string]any{"farmid": farmID, "createdBy": userID}, update)
}

func SQLdeleteFarmByID(ctx context.Context, database db.Database, farmID string) (int64, error) {
	return database.DeleteOne(ctx, farmsTable, map[string]any{"farmid": farmID})
}
