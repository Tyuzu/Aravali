// File: internal/menu/menuSQLDB.go

package menu

import (
	"context"
	"fmt"
	"strings"
	"time"

	"scav/config"
	"scav/infra"
)

var menuTable = config.Tables.MenuTable

func findMenuByPlaceAndID(ctx context.Context, app *infra.Deps, placeID, menuID string) (Menu, error) {
	if app == nil || app.SQLDB == nil {
		return Menu{}, nil
	}
	if strings.TrimSpace(placeID) == "" || strings.TrimSpace(menuID) == "" {
		return Menu{}, fmt.Errorf("place id and menu id are required")
	}

	var m Menu
	row := app.SQLDB.QueryRow(
		ctx,
		`SELECT menuid, placeid, name, price, discount, stock, menu_pic, description, userid, created_at, updated_at FROM `+menuTable+` WHERE placeid = $1 AND menuid = $2 LIMIT 1`,
		placeID, menuID,
	)
	if err := row.Scan(
		&m.MenuID, &m.PlaceID, &m.Name, &m.Price, &m.Discount, &m.Stock,
		&m.MenuPhoto, &m.Description, &m.UserID, &m.CreatedAt, &m.UpdatedAt,
	); err != nil {
		return Menu{}, err
	}
	return m, nil
}

func findMenusByPlace(ctx context.Context, app *infra.Deps, placeID string) ([]Menu, error) {
	if app == nil || app.SQLDB == nil {
		return []Menu{}, nil
	}
	if strings.TrimSpace(placeID) == "" {
		return nil, fmt.Errorf("place id is required")
	}

	rows, err := app.SQLDB.Query(
		ctx,
		`SELECT menuid, placeid, name, price, discount, stock, menu_pic, description, userid, created_at, updated_at FROM `+menuTable+` WHERE placeid = $1 ORDER BY created_at DESC`,
		placeID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Menu, 0)
	for rows.Next() {
		var m Menu
		if err := rows.Scan(
			&m.MenuID, &m.PlaceID, &m.Name, &m.Price, &m.Discount, &m.Stock,
			&m.MenuPhoto, &m.Description, &m.UserID, &m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, m)
	}
	return items, rows.Err()
}

func insertMenu(ctx context.Context, app *infra.Deps, menu Menu) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	if strings.TrimSpace(menu.MenuID) == "" {
		menu.MenuID = fmt.Sprintf("menu_%d", time.Now().UnixNano())
	}
	if menu.CreatedAt.IsZero() {
		menu.CreatedAt = time.Now()
	}
	if menu.UpdatedAt.IsZero() {
		menu.UpdatedAt = menu.CreatedAt
	}

	_, err := app.SQLDB.Exec(
		ctx,
		`INSERT INTO `+menuTable+` (menuid, placeid, name, price, discount, stock, menu_pic, description, userid, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		menu.MenuID, menu.PlaceID, menu.Name, menu.Price, menu.Discount, menu.Stock, menu.MenuPhoto, menu.Description, menu.UserID, menu.CreatedAt, menu.UpdatedAt,
	)
	return err
}

func updateMenuFields(ctx context.Context, app *infra.Deps, placeID, menuID string, updateFields map[string]any) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	if strings.TrimSpace(placeID) == "" || strings.TrimSpace(menuID) == "" {
		return fmt.Errorf("place id and menu id are required")
	}
	if len(updateFields) == 0 {
		return nil
	}

	parts := make([]string, 0, len(updateFields))
	args := make([]any, 0, len(updateFields)+2)
	for key, value := range updateFields {
		if key == "" || value == nil {
			continue
		}
		parts = append(parts, fmt.Sprintf("%s = $%d", key, len(args)+1))
		args = append(args, value)
	}
	if len(parts) == 0 {
		return nil
	}
	args = append(args, placeID, menuID)
	_, err := app.SQLDB.Exec(
		ctx,
		`UPDATE `+menuTable+` SET `+strings.Join(parts, ", ")+` WHERE placeid = $`+fmt.Sprintf("%d", len(args)-1)+` AND menuid = $`+fmt.Sprintf("%d", len(args)),
		args...,
	)
	return err
}

func deleteMenuByID(ctx context.Context, app *infra.Deps, placeID, menuID string) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	if strings.TrimSpace(placeID) == "" || strings.TrimSpace(menuID) == "" {
		return fmt.Errorf("place id and menu id are required")
	}
	_, err := app.SQLDB.Exec(ctx, `DELETE FROM `+menuTable+` WHERE placeid = $1 AND menuid = $2`, placeID, menuID)
	return err
}

func decrementMenuStock(ctx context.Context, app *infra.Deps, placeID, menuID string, quantity int) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	if strings.TrimSpace(placeID) == "" || strings.TrimSpace(menuID) == "" {
		return fmt.Errorf("place id and menu id are required")
	}
	if quantity <= 0 {
		return nil
	}
	_, err := app.SQLDB.Exec(
		ctx,
		`UPDATE `+menuTable+` SET stock = stock - $1, updated_at = NOW() WHERE placeid = $2 AND menuid = $3`,
		quantity, placeID, menuID,
	)
	return err
}
