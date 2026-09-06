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
	cropsCollection      = config.Collections.CropsCollection
	farmsCollection      = config.Collections.FarmsCollection
	farmOrdersCollection = config.Collections.FarmOrdersCollection
	// usersCollection and productsCollection removed because they were unused
)

func insertFarm(ctx context.Context, database db.Database, farm Farm) error {
	return database.InsertOne(ctx, farmsCollection, farm)
}

func getFarmByID(ctx context.Context, database db.Database, farmID string) (Farm, error) {
	var farm Farm
	err := database.FindOne(ctx, farmsCollection, map[string]any{"farmid": farmID}, &farm)
	return farm, err
}

func getFarmByCreatedBy(ctx context.Context, database db.Database, userID string) (Farm, error) {
	var farm Farm
	err := database.FindOne(ctx, farmsCollection, map[string]any{"createdBy": userID}, &farm)
	return farm, err
}

func getCropsByFarmID(ctx context.Context, database db.Database, farmID string) ([]Crop, error) {
	var crops []Crop
	if err := database.FindMany(ctx, cropsCollection, map[string]any{"farmid": farmID}, &crops); err != nil {
		return nil, err
	}
	if crops == nil {
		return []Crop{}, nil
	}
	return crops, nil
}

func getFarmOrdersByFarmID(ctx context.Context, database db.Database, farmID string) ([]cart.FarmOrder, error) {
	var orders []cart.FarmOrder
	if err := database.FindMany(ctx, farmOrdersCollection, map[string]any{"farmid": farmID}, &orders); err != nil {
		return nil, err
	}
	if orders == nil {
		return []cart.FarmOrder{}, nil
	}
	return orders, nil
}

func getCropsByCropID(ctx context.Context, database db.Database, cropID string) ([]Crop, error) {
	var crops []Crop
	if err := database.FindMany(ctx, cropsCollection, map[string]any{"cropid": cropID}, &crops); err != nil {
		return nil, err
	}
	if crops == nil {
		return []Crop{}, nil
	}
	return crops, nil
}

func getFarmsByIDs(ctx context.Context, database db.Database, farmIDs []string) ([]Farm, error) {
	var farms []Farm
	if err := database.FindMany(ctx, farmsCollection, map[string]any{"farmid": map[string]any{"$in": farmIDs}}, &farms); err != nil {
		return nil, err
	}
	if farms == nil {
		return []Farm{}, nil
	}
	return farms, nil
}

func getCropsByNameFilter(ctx context.Context, database db.Database, cropName string) ([]Crop, error) {
	filter := map[string]any{
		"name": map[string]any{
			"$regex":   "^" + regexp.QuoteMeta(cropName) + "$",
			"$options": "i",
		},
	}
	var crops []Crop
	if err := database.FindMany(ctx, cropsCollection, filter, &crops); err != nil {
		return nil, err
	}
	if crops == nil {
		return []Crop{}, nil
	}
	return crops, nil
}

func getPaginatedFarms(ctx context.Context, database db.Database, search string, skip, limit int) ([]Farm, int64, error) {
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
	if err := database.Aggregate(ctx, farmsCollection, pipeline, &farms); err != nil {
		return nil, 0, err
	}

	total, err := database.CountDocuments(ctx, farmsCollection, map[string]any{})
	if err != nil {
		if search == "" {
			return farms, 0, nil
		}
		return farms, 0, err
	}
	return farms, total, nil
}

func getMyFarmsPage(ctx context.Context, database db.Database, userID string, skip, limit int) ([]Farm, int64, error) {
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
	if err := database.Aggregate(ctx, farmsCollection, pipeline, &farms); err != nil {
		return nil, 0, err
	}

	total, err := database.CountDocuments(ctx, farmsCollection, map[string]any{"createdBy": userID})
	if err != nil {
		return farms, 0, nil
	}
	return farms, total, nil
}

func updateOwnedFarm(ctx context.Context, database db.Database, farmID, userID string, update any) (any, error) {
	// The owner field on Farm is stored as "createdBy" (see farmModels.go).
	// Use that field to ensure the update only affects farms owned by the user.
	return database.UpdateOne(ctx, farmsCollection, map[string]any{"farmid": farmID, "createdBy": userID}, update)
}

func deleteFarmByID(ctx context.Context, database db.Database, farmID string) (int64, error) {
	return database.DeleteOne(ctx, farmsCollection, map[string]any{"farmid": farmID})
}
