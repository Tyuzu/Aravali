package menu

import (
	"context"
	"scav/config"
	"scav/infra"
	"time"
)

var menuCollection = config.Collections.MenuCollection

func findMenuByPlaceAndID(ctx context.Context, app *infra.Deps, placeID, menuID string) (Menu, error) {
	var menu Menu
	err := app.DB.FindOne(ctx, menuCollection, map[string]string{
		"placeid": placeID,
		"menuid":  menuID,
	}, &menu)
	return menu, err
}

func findMenusByPlace(ctx context.Context, app *infra.Deps, placeID string) ([]Menu, error) {
	var menus []Menu
	if err := app.DB.FindMany(ctx, menuCollection, map[string]string{
		"placeid": placeID,
	}, &menus); err != nil {
		return nil, err
	}
	if menus == nil {
		return []Menu{}, nil
	}
	return menus, nil
}

func insertMenu(ctx context.Context, app *infra.Deps, menu Menu) error {
	return app.DB.Insert(ctx, menuCollection, menu)
}

func updateMenuFields(ctx context.Context, app *infra.Deps, placeID, menuID string, updateFields map[string]any) error {
	_, err := app.DB.UpdateOne(ctx, menuCollection, map[string]string{"placeid": placeID, "menuid": menuID}, map[string]any{"$set": updateFields})
	return err
}

func deleteMenuByID(ctx context.Context, app *infra.Deps, placeID, menuID string) error {
	_, err := app.DB.DeleteOne(ctx, menuCollection, map[string]string{"placeid": placeID, "menuid": menuID})
	return err
}

func decrementMenuStock(ctx context.Context, app *infra.Deps, placeID, menuID string, quantity int) error {
	_, err := app.DB.UpdateOne(ctx, menuCollection, map[string]string{"placeid": placeID, "menuid": menuID}, map[string]any{
		"$inc": map[string]int{"stock": -quantity},
		"$set": map[string]any{"updated_at": time.Now()},
	})
	return err
}
