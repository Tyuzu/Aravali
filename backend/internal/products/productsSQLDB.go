// File: internal/products/productsSQLDB.go

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
func GetProductByID(ctx context.Context, app *infra.Deps, id string, out *farms.Product) error {
	query := "productid = $1"
	args := []any{id}

	return app.SQLDB.FindOne(ctx, productsTable, query, args, out)
}

func FindProductsWithOptions(ctx context.Context, app *infra.Deps, query string, args []any, opts sqldb.FindManyOptions, out *[]farms.Product) error {
	return app.SQLDB.FindManyWithOptions(ctx, productsTable, query, args, opts, out)
}

func CountProducts(ctx context.Context, app *infra.Deps, query string, args []any) (int64, error) {
	return app.SQLDB.CountDocuments(ctx, productsTable, query, args)
}

func InsertProduct(ctx context.Context, app *infra.Deps, item farms.Product) error {
	return app.SQLDB.InsertOne(ctx, productsTable, item)
}

func UpdateProductByID(ctx context.Context, app *infra.Deps, id string, update map[string]any) (int64, error) {
	query := "productid = $1"
	args := []any{id}

	return app.SQLDB.UpdateOne(ctx, productsTable, query, args, update)
}

func DeleteProductByID(ctx context.Context, app *infra.Deps, id string) (int64, error) {
	query := "productid = $1"
	args := []any{id}

	return app.SQLDB.DeleteOne(ctx, productsTable, query, args)
}

// Orders
func GetFarmOrderByID(ctx context.Context, app *infra.Deps, id string, out *cart.FarmOrder) error {
	query := "orderid = $1"
	args := []any{id}

	return app.SQLDB.FindOne(ctx, farmOrdersTable, query, args, out)
}

func FindFarmOrders(ctx context.Context, app *infra.Deps, query string, args []any, out *[]cart.FarmOrder) error {
	return app.SQLDB.FindMany(ctx, farmOrdersTable, query, args, out)
}

func UpdateFarmOrderByID(ctx context.Context, app *infra.Deps, id string, update map[string]any) (int64, error) {
	query := "orderid = $1"
	args := []any{id}

	return app.SQLDB.UpdateOne(ctx, farmOrdersTable, query, args, update)
}

// Farms, users, crops
func FindFarmsByFilter(ctx context.Context, app *infra.Deps, query string, args []any, out *[]farms.Farm) error {
	return app.SQLDB.FindMany(ctx, farmsTable, query, args, out)
}

func GetFarmByID(ctx context.Context, app *infra.Deps, id string, out *farms.Farm) error {
	query := "farmid = $1"
	args := []any{id}

	return app.SQLDB.FindOne(ctx, farmsTable, query, args, out)
}

func GetUserByID(ctx context.Context, app *infra.Deps, id string, out *auth.User) error {
	query := "userid = $1"
	args := []any{id}

	return app.SQLDB.FindOne(ctx, usersTable, query, args, out)
}

func GetCropByID(ctx context.Context, app *infra.Deps, id string, out *farms.Crop) error {
	query := "cropid = $1"
	args := []any{id}

	return app.SQLDB.FindOne(ctx, cropsTable, query, args, out)
}

func FindTransactions(ctx context.Context, app *infra.Deps, query string, args []any, out any) error {
	return app.SQLDB.FindMany(ctx, "transactions", query, args, out)
}

// Crop atomic update
func FindOneAndUpdateCrop(ctx context.Context, app *infra.Deps, query string, args []any, update map[string]any, out any) error {
	return app.SQLDB.FindOneAndUpdate(ctx, cropsTable, query, args, update, out)
}

func UpdateCropByFilter(ctx context.Context, app *infra.Deps, query string, args []any, update map[string]any) (int64, error) {
	return app.SQLDB.UpdateOne(ctx, cropsTable, query, args, update)
}
