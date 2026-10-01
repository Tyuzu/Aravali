// File: internal/verticals/menu/menuSQLDB.go

package menu

import (
	"context"
	"time"

	"scav/config"
	"scav/infra"
)

var menuTable = config.Tables.MenuTable

func findMenuByPlaceAndID(ctx context.Context, app *infra.Deps, placeID, menuID string) (Menu, error) {
	var menu Menu
	query := "placeid = $1 AND menuid = $2"
	args := []any{placeID, menuID}

	err := app.SQLDB.FindOne(ctx, menuTable, query, args, &menu)
	return menu, err
}

func findMenusByPlace(ctx context.Context, app *infra.Deps, placeID string) ([]Menu, error) {
	var menus []Menu
	query := "placeid = $1"
	args := []any{placeID}

	if err := app.SQLDB.FindMany(ctx, menuTable, query, args, &menus); err != nil {
		return nil, err
	}
	if menus == nil {
		return []Menu{}, nil
	}
	return menus, nil
}

func insertMenu(ctx context.Context, app *infra.Deps, menu Menu) error {
	return app.SQLDB.Insert(ctx, menuTable, menu)
}

func updateMenuFields(ctx context.Context, app *infra.Deps, placeID, menuID string, updateFields map[string]any) error {
	query := "placeid = $1 AND menuid = $2"
	args := []any{placeID, menuID}

	_, err := app.SQLDB.UpdateOne(ctx, menuTable, query, args, updateFields)
	return err
}

func deleteMenuByID(ctx context.Context, app *infra.Deps, placeID, menuID string) error {
	query := "placeid = $1 AND menuid = $2"
	args := []any{placeID, menuID}

	_, err := app.SQLDB.DeleteOne(ctx, menuTable, query, args)
	return err
}

func decrementMenuStock(ctx context.Context, app *infra.Deps, placeID, menuID string, quantity int) error {
	query := "placeid = $1 AND menuid = $2"
	args := []any{placeID, menuID}

	// Decrement stock atomically
	if err := app.SQLDB.Inc(ctx, menuTable, query, args, "stock", int64(-quantity)); err != nil {
		return err
	}

	// Update timestamp
	_, err := app.SQLDB.UpdateOne(ctx, menuTable, query, args, map[string]any{
		"updated_at": time.Now(),
	})
	return err
}
