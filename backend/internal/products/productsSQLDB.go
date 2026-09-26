package products

import (
	"context"

	"scav/config"
	"scav/infra"
	"scav/infra/sqldb"
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
	query := "productid = $1"
	args := []any{id}

	return app.SQLDB.FindOne(ctx, productsTable, query, args, out)
}

func SQLFindProductsWithOptions(ctx context.Context, app *infra.Deps, query string, args []any, opts sqldb.FindManyOptions, out *[]farms.Product) error {
	return app.SQLDB.FindManyWithOptions(ctx, productsTable, query, args, opts, out)
}

func SQLCountProducts(ctx context.Context, app *infra.Deps, query string, args []any) (int64, error) {
	return app.SQLDB.CountDocuments(ctx, productsTable, query, args)
}

func SQLInsertProduct(ctx context.Context, app *infra.Deps, item farms.Product) error {
	return app.SQLDB.InsertOne(ctx, productsTable, item)
}

func SQLUpdateProductByID(ctx context.Context, app *infra.Deps, id string, update map[string]any) (int64, error) {
	query := "productid = $1"
	args := []any{id}

	return app.SQLDB.UpdateOne(ctx, productsTable, query, args, update)
}

func SQLDeleteProductByID(ctx context.Context, app *infra.Deps, id string) (int64, error) {
	query := "productid = $1"
	args := []any{id}

	return app.SQLDB.DeleteOne(ctx, productsTable, query, args)
}

// Orders
func SQLGetFarmOrderByID(ctx context.Context, app *infra.Deps, id string, out *cart.FarmOrder) error {
	query := "orderid = $1"
	args := []any{id}

	return app.SQLDB.FindOne(ctx, farmOrdersTable, query, args, out)
}

func SQLFindFarmOrders(ctx context.Context, app *infra.Deps, query string, args []any, out *[]cart.FarmOrder) error {
	return app.SQLDB.FindMany(ctx, farmOrdersTable, query, args, out)
}

func SQLUpdateFarmOrderByID(ctx context.Context, app *infra.Deps, id string, update map[string]any) (int64, error) {
	query := "orderid = $1"
	args := []any{id}

	return app.SQLDB.UpdateOne(ctx, farmOrdersTable, query, args, update)
}

// Farms, users, crops
func SQLFindFarmsByFilter(ctx context.Context, app *infra.Deps, query string, args []any, out *[]farms.Farm) error {
	return app.SQLDB.FindMany(ctx, farmsTable, query, args, out)
}

func SQLGetFarmByID(ctx context.Context, app *infra.Deps, id string, out *farms.Farm) error {
	query := "farmid = $1"
	args := []any{id}

	return app.SQLDB.FindOne(ctx, farmsTable, query, args, out)
}

func SQLGetUserByID(ctx context.Context, app *infra.Deps, id string, out *auth.User) error {
	query := "userid = $1"
	args := []any{id}

	return app.SQLDB.FindOne(ctx, usersTable, query, args, out)
}

func SQLGetCropByID(ctx context.Context, app *infra.Deps, id string, out *farms.Crop) error {
	query := "cropid = $1"
	args := []any{id}

	return app.SQLDB.FindOne(ctx, cropsTable, query, args, out)
}

func SQLFindTransactions(ctx context.Context, app *infra.Deps, query string, args []any, out any) error {
	return app.SQLDB.FindMany(ctx, "transactions", query, args, out)
}

// Crop atomic update
func SQLFindOneAndUpdateCrop(ctx context.Context, app *infra.Deps, query string, args []any, update map[string]any, out any) error {
	return app.SQLDB.FindOneAndUpdate(ctx, cropsTable, query, args, update, out)
}

func SQLUpdateCropByFilter(ctx context.Context, app *infra.Deps, query string, args []any, update map[string]any) (int64, error) {
	return app.SQLDB.UpdateOne(ctx, cropsTable, query, args, update)
}
