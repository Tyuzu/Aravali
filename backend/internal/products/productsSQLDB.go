// File: internal/products/productsSQLDB.go

package products

import (
	"context"
	"strconv"
	"strings"

	"scav/config"
	"scav/infra"
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
}

func FindProductsWithOptions(ctx context.Context, app *infra.Deps, query string, args []any, opts map[string]any, out *[]farms.Product) error {
}

func buildProductSearchQuery(itemType, category, search string) (string, []any) {
	query := "1 = 1"
	args := []any{}

	if itemType != "" {
		query += " AND type = $" + strconv.Itoa(len(args)+1)
		args = append(args, itemType)
	}
	if category != "" {
		query += " AND category = $" + strconv.Itoa(len(args)+1)
		args = append(args, category)
	}
	if search != "" {
		query += " AND LOWER(name) LIKE $" + strconv.Itoa(len(args)+1)
		args = append(args, "%"+strings.ToLower(search)+"%")
	}

	return query, args
}

func buildProductListOptions(skip, limit int, sort string) map[string]any {

}

func getProductsPage(ctx context.Context, app *infra.Deps, itemType, category, search string, skip, limit int, sort string) ([]farms.Product, int64, error) {
	query, args := buildProductSearchQuery(itemType, category, search)
	opts := buildProductListOptions(skip, limit, sort)

	var items []farms.Product
	if err := FindProductsWithOptions(ctx, app, query, args, opts, &items); err != nil {
		return nil, 0, err
	}

	total, err := CountProducts(ctx, app, query, args)
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func CountProducts(ctx context.Context, app *infra.Deps, query string, args []any) (int64, error) {
}

func InsertProduct(ctx context.Context, app *infra.Deps, item farms.Product) error {
}

func UpdateProductByID(ctx context.Context, app *infra.Deps, id string, update map[string]any) (int64, error) {
}

func DeleteProductByID(ctx context.Context, app *infra.Deps, id string) (int64, error) {
}

// Orders
func GetFarmOrderByID(ctx context.Context, app *infra.Deps, id string, out *cart.FarmOrder) error {
}

func FindFarmOrders(ctx context.Context, app *infra.Deps, query string, args []any, out *[]cart.FarmOrder) error {
}

func UpdateFarmOrderByID(ctx context.Context, app *infra.Deps, id string, update map[string]any) (int64, error) {
}

// Farms, users, crops
func FindFarmsByFilter(ctx context.Context, app *infra.Deps, query string, args []any, out *[]farms.Farm) error {
}

func GetFarmByID(ctx context.Context, app *infra.Deps, id string, out *farms.Farm) error {
}

func GetUserByID(ctx context.Context, app *infra.Deps, id string, out *auth.User) error {
}

func GetCropByID(ctx context.Context, app *infra.Deps, id string, out *farms.Crop) error {
}

func FindTransactions(ctx context.Context, app *infra.Deps, query string, args []any, out any) error {
}

// Crop atomic update
func FindOneAndUpdateCrop(ctx context.Context, app *infra.Deps, query string, args []any, update map[string]any, out any) error {
}

func UpdateCropByFilter(ctx context.Context, app *infra.Deps, query string, args []any, update map[string]any) (int64, error) {
}
