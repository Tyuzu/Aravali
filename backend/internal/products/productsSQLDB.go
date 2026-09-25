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
	farmOrdersTable = config.Tables.FarmOrdersTable
	productsTable   = config.Tables.ProductTable
	cropsTable      = config.Tables.CropsTable
	usersTable      = config.Tables.UserTable
	farmsTable      = config.Tables.FarmsTable
)

// Products
func SQLGetProductByID(ctx context.Context, app *infra.Deps, id string, out *farms.Product) error {
	return app.DB.FindOne(ctx, productsTable, map[string]any{"productid": id}, out)
}

func SQLFindProductsWithOptions(ctx context.Context, app *infra.Deps, filter map[string]any, opts db.FindManyOptions, out *[]farms.Product) error {
	return app.DB.FindManyWithOptions(ctx, productsTable, filter, opts, out)
}

func SQLCountProducts(ctx context.Context, app *infra.Deps, filter map[string]any) (int64, error) {
	return app.DB.CountDocuments(ctx, productsTable, filter)
}

func SQLInsertProduct(ctx context.Context, app *infra.Deps, item farms.Product) error {
	return app.DB.InsertOne(ctx, productsTable, item)
}

func SQLUpdateProductByID(ctx context.Context, app *infra.Deps, id string, update any) (any, error) {
	return app.DB.UpdateOne(ctx, productsTable, map[string]any{"productid": id}, update)
}

func SQLDeleteProductByID(ctx context.Context, app *infra.Deps, id string) (int64, error) {
	return app.DB.DeleteOne(ctx, productsTable, map[string]any{"productid": id})
}

// Orders
func SQLGetFarmOrderByID(ctx context.Context, app *infra.Deps, id string, out *cart.FarmOrder) error {
	return app.DB.FindOne(ctx, farmOrdersTable, map[string]any{"orderid": id}, out)
}

func SQLFindFarmOrders(ctx context.Context, app *infra.Deps, filter map[string]any, out *[]cart.FarmOrder) error {
	return app.DB.FindMany(ctx, farmOrdersTable, filter, out)
}

func SQLUpdateFarmOrderByID(ctx context.Context, app *infra.Deps, id string, update any) (any, error) {
	return app.DB.UpdateOne(ctx, farmOrdersTable, map[string]any{"orderid": id}, update)
}

// Farms, users, crops
func SQLFindFarmsByFilter(ctx context.Context, app *infra.Deps, filter map[string]any, out *[]farms.Farm) error {
	return app.DB.FindMany(ctx, farmsTable, filter, out)
}

func SQLGetFarmByID(ctx context.Context, app *infra.Deps, id string, out *farms.Farm) error {
	return app.DB.FindOne(ctx, farmsTable, map[string]any{"farmid": id}, out)
}

func SQLGetUserByID(ctx context.Context, app *infra.Deps, id string, out *auth.User) error {
	return app.DB.FindOne(ctx, usersTable, map[string]any{"userid": id}, out)
}

func SQLGetCropByID(ctx context.Context, app *infra.Deps, id string, out *farms.Crop) error {
	return app.DB.FindOne(ctx, cropsTable, map[string]any{"cropid": id}, out)
}

func SQLFindTransactions(ctx context.Context, app *infra.Deps, filter map[string]any, out any) error {
	return app.DB.FindMany(ctx, "transactions", filter, out)
}

// Crop atomic update
func SQLFindOneAndUpdateCrop(ctx context.Context, app *infra.Deps, filter any, update any, out any) error {
	return app.DB.FindOneAndUpdate(ctx, cropsTable, filter, update, out)
}

func SQLUpdateCropByFilter(ctx context.Context, app *infra.Deps, filter any, update any) (any, error) {
	return app.DB.UpdateOne(ctx, cropsTable, filter, update)
}
