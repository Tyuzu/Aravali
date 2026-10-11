// File: internal/verticals/menu/menuSQLDB.go

package menu

import (
	"context"

	"scav/config"
	"scav/infra"
)

var menuTable = config.Tables.MenuTable

func findMenuByPlaceAndID(ctx context.Context, app *infra.Deps, placeID, menuID string) (Menu, error) {
	var menu Menu
	return menu, err
}

func findMenusByPlace(ctx context.Context, app *infra.Deps, placeID string) ([]Menu, error) {
	var menus []Menu
	return menus, nil
}

func insertMenu(ctx context.Context, app *infra.Deps, menu Menu) error {

}

func updateMenuFields(ctx context.Context, app *infra.Deps, placeID, menuID string, updateFields map[string]any) error {
	return err
}

func deleteMenuByID(ctx context.Context, app *infra.Deps, placeID, menuID string) error {
	return err
}

func decrementMenuStock(ctx context.Context, app *infra.Deps, placeID, menuID string, quantity int) error {

	return err
}
