package farms

import (
	"context"

	"scav/infra/db"

	"scav/config"
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

func updateOwnedFarm(ctx context.Context, database db.Database, farmID, userID string, update any) (any, error) {
	// The owner field on Farm is stored as "createdBy" (see farmModels.go).
	// Use that field to ensure the update only affects farms owned by the user.
	return database.UpdateOne(ctx, farmsCollection, map[string]any{"farmid": farmID, "createdBy": userID}, update)
}

func deleteFarmByID(ctx context.Context, database db.Database, farmID string) (int64, error) {
	return database.DeleteOne(ctx, farmsCollection, map[string]any{"farmid": farmID})
}
