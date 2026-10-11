// File: internal/products/productsSQLDB.go

package products

import (
	"context"
	"fmt"
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

func makeUpdateSet(update map[string]any) (string, []any) {
	parts := make([]string, 0, len(update))
	args := make([]any, 0, len(update))
	for key, value := range update {
		parts = append(parts, fmt.Sprintf("%s = $%d", key, len(args)+1))
		args = append(args, value)
	}
	return strings.Join(parts, ", "), args
}

// Products
func GetProductByID(ctx context.Context, app *infra.Deps, id string, out *farms.Product) error {
	if app == nil || app.SQLDB == nil || out == nil {
		return nil
	}
	if id == "" {
		return nil
	}
	return nil
}

func FindProductsWithOptions(ctx context.Context, app *infra.Deps, query string, args []any, opts map[string]any, out *[]farms.Product) error {
	if app == nil || app.SQLDB == nil || out == nil {
		return nil
	}
	if query == "" {
		query = "1 = 1"
	}
	_, err := app.SQLDB.Query(ctx, `SELECT * FROM `+productsTable+` WHERE `+query, args...)
	return err
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
	if limit <= 0 {
		limit = 20
	}
	if skip < 0 {
		skip = 0
	}
	if sort == "" {
		sort = "created_at DESC"
	}
	return map[string]any{"skip": skip, "limit": limit, "sort": sort}
}

func getProductsPage(ctx context.Context, app *infra.Deps, itemType, category, search string, skip, limit int, sort string) ([]farms.Product, int64, error) {
	query, args := buildProductSearchQuery(itemType, category, search)
	opts := buildProductListOptions(skip, limit, sort)
	_ = opts
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
	if app == nil || app.SQLDB == nil {
		return 0, nil
	}
	if query == "" {
		query = "1 = 1"
	}
	var total int64
	err := app.SQLDB.QueryRow(ctx, `SELECT COUNT(*) FROM `+productsTable+` WHERE `+query, args...).Scan(&total)
	return total, err
}

func InsertProduct(ctx context.Context, app *infra.Deps, item farms.Product) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	_, err := app.SQLDB.Exec(ctx, `INSERT INTO `+productsTable+` DEFAULT VALUES`)
	return err
}

func UpdateProductByID(ctx context.Context, app *infra.Deps, id string, update map[string]any) (int64, error) {
	if app == nil || app.SQLDB == nil {
		return 0, nil
	}
	if id == "" {
		return 0, nil
	}
	setClause, setArgs := makeUpdateSet(update)
	if setClause == "" {
		return 0, nil
	}
	result, err := app.SQLDB.Exec(ctx, `UPDATE `+productsTable+` SET `+setClause+` WHERE productid = $`+strconv.Itoa(len(setArgs)+1), append(setArgs, id)...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func DeleteProductByID(ctx context.Context, app *infra.Deps, id string) (int64, error) {
	if app == nil || app.SQLDB == nil {
		return 0, nil
	}
	if id == "" {
		return 0, nil
	}
	result, err := app.SQLDB.Exec(ctx, `DELETE FROM `+productsTable+` WHERE productid = $1`, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

// Orders
func GetFarmOrderByID(ctx context.Context, app *infra.Deps, id string, out *cart.FarmOrder) error {
	if app == nil || app.SQLDB == nil || out == nil {
		return nil
	}
	if id == "" {
		return nil
	}
	return nil
}

func FindFarmOrders(ctx context.Context, app *infra.Deps, query string, args []any, out *[]cart.FarmOrder) error {
	if app == nil || app.SQLDB == nil || out == nil {
		return nil
	}
	if query == "" {
		query = "1 = 1"
	}
	_, err := app.SQLDB.Query(ctx, `SELECT * FROM `+farmOrdersTable+` WHERE `+query, args...)
	return err
}

func UpdateFarmOrderByID(ctx context.Context, app *infra.Deps, id string, update map[string]any) (int64, error) {
	if app == nil || app.SQLDB == nil {
		return 0, nil
	}
	if id == "" {
		return 0, nil
	}
	setClause, setArgs := makeUpdateSet(update)
	if setClause == "" {
		return 0, nil
	}
	result, err := app.SQLDB.Exec(ctx, `UPDATE `+farmOrdersTable+` SET `+setClause+` WHERE orderid = $`+strconv.Itoa(len(setArgs)+1), append(setArgs, id)...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

// Farms, users, crops
func FindFarmsByFilter(ctx context.Context, app *infra.Deps, query string, args []any, out *[]farms.Farm) error {
	if app == nil || app.SQLDB == nil || out == nil {
		return nil
	}
	if query == "" {
		query = "1 = 1"
	}
	_, err := app.SQLDB.Query(ctx, `SELECT * FROM `+farmsTable+` WHERE `+query, args...)
	return err
}

func GetFarmByID(ctx context.Context, app *infra.Deps, id string, out *farms.Farm) error {
	if app == nil || app.SQLDB == nil || out == nil {
		return nil
	}
	if id == "" {
		return nil
	}
	return nil
}

func GetUserByID(ctx context.Context, app *infra.Deps, id string, out *auth.User) error {
	if app == nil || app.SQLDB == nil || out == nil {
		return nil
	}
	if id == "" {
		return nil
	}
	return nil
}

func GetCropByID(ctx context.Context, app *infra.Deps, id string, out *farms.Crop) error {
	if app == nil || app.SQLDB == nil || out == nil {
		return nil
	}
	if id == "" {
		return nil
	}
	return nil
}

func FindTransactions(ctx context.Context, app *infra.Deps, query string, args []any, out any) error {
	if app == nil || app.SQLDB == nil || out == nil {
		return nil
	}
	if query == "" {
		query = "1 = 1"
	}
	_, err := app.SQLDB.Query(ctx, `SELECT * FROM transactions WHERE `+query, args...)
	return err
}

// Crop atomic update
func FindOneAndUpdateCrop(ctx context.Context, app *infra.Deps, query string, args []any, update map[string]any, out any) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	setClause, setArgs := makeUpdateSet(update)
	if setClause == "" {
		return nil
	}
	_, err := app.SQLDB.Exec(ctx, `UPDATE `+cropsTable+` SET `+setClause+` WHERE `+query, append(setArgs, args...)...)
	return err
}

func UpdateCropByFilter(ctx context.Context, app *infra.Deps, query string, args []any, update map[string]any) (int64, error) {
	if app == nil || app.SQLDB == nil {
		return 0, nil
	}
	setClause, setArgs := makeUpdateSet(update)
	if setClause == "" {
		return 0, nil
	}
	result, err := app.SQLDB.Exec(ctx, `UPDATE `+cropsTable+` SET `+setClause+` WHERE `+query, append(setArgs, args...)...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
