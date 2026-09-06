package products

import (
	"context"

	"scav/config"
	"scav/infra"
	"scav/infra/db"
	"scav/internal/auth"
	"scav/internal/cart"
	"scav/internal/farms"
)

var (
	farmOrdersCollection = config.Collections.FarmOrdersCollection
	productsCollection   = config.Collections.ProductCollection
	cropsCollection      = config.Collections.CropsCollection
	usersCollection      = config.Collections.UserCollection
	farmsCollection      = config.Collections.FarmsCollection
)

// Products
func GetProductByID(ctx context.Context, app *infra.Deps, id string, out *farms.Product) error {
	return app.DB.FindOne(ctx, productsCollection, map[string]any{"productid": id}, out)
}

func FindProductsWithOptions(ctx context.Context, app *infra.Deps, filter map[string]any, opts db.FindManyOptions, out *[]farms.Product) error {
	return app.DB.FindManyWithOptions(ctx, productsCollection, filter, opts, out)
}

func CountProducts(ctx context.Context, app *infra.Deps, filter map[string]any) (int64, error) {
	return app.DB.CountDocuments(ctx, productsCollection, filter)
}

func InsertProduct(ctx context.Context, app *infra.Deps, item farms.Product) error {
	return app.DB.InsertOne(ctx, productsCollection, item)
}

func UpdateProductByID(ctx context.Context, app *infra.Deps, id string, update any) (any, error) {
	return app.DB.UpdateOne(ctx, productsCollection, map[string]any{"productid": id}, update)
}

func DeleteProductByID(ctx context.Context, app *infra.Deps, id string) (int64, error) {
	return app.DB.DeleteOne(ctx, productsCollection, map[string]any{"productid": id})
}

// Orders
func GetFarmOrderByID(ctx context.Context, app *infra.Deps, id string, out *cart.FarmOrder) error {
	return app.DB.FindOne(ctx, farmOrdersCollection, map[string]any{"orderid": id}, out)
}

func FindFarmOrders(ctx context.Context, app *infra.Deps, filter map[string]any, out *[]cart.FarmOrder) error {
	return app.DB.FindMany(ctx, farmOrdersCollection, filter, out)
}

func UpdateFarmOrderByID(ctx context.Context, app *infra.Deps, id string, update any) (any, error) {
	return app.DB.UpdateOne(ctx, farmOrdersCollection, map[string]any{"orderid": id}, update)
}

// Farms, users, crops
func FindFarmsByFilter(ctx context.Context, app *infra.Deps, filter map[string]any, out *[]farms.Farm) error {
	return app.DB.FindMany(ctx, farmsCollection, filter, out)
}

func GetFarmByID(ctx context.Context, app *infra.Deps, id string, out *farms.Farm) error {
	return app.DB.FindOne(ctx, farmsCollection, map[string]any{"farmid": id}, out)
}

func GetUserByID(ctx context.Context, app *infra.Deps, id string, out *auth.User) error {
	return app.DB.FindOne(ctx, usersCollection, map[string]any{"userid": id}, out)
}

func GetCropByID(ctx context.Context, app *infra.Deps, id string, out *farms.Crop) error {
	return app.DB.FindOne(ctx, cropsCollection, map[string]any{"cropid": id}, out)
}

func FindTransactions(ctx context.Context, app *infra.Deps, filter map[string]any, out any) error {
	return app.DB.FindMany(ctx, "transactions", filter, out)
}

// Crop atomic update
func FindOneAndUpdateCrop(ctx context.Context, app *infra.Deps, filter any, update any, out any) error {
	return app.DB.FindOneAndUpdate(ctx, cropsCollection, filter, update, out)
}

func UpdateCropByFilter(ctx context.Context, app *infra.Deps, filter any, update any) (any, error) {
	return app.DB.UpdateOne(ctx, cropsCollection, filter, update)
}
