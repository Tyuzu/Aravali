// File: internal/merch/merchSQLDB.go

package merch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"scav/config"
	"scav/infra"
)

var merchTable = config.Tables.MerchTable

func getEntityOwner(ctx context.Context, app *infra.Deps, entityType, entityID string) (string, error) {
	switch entityType {
	case "event", "farm", "artist":
		if app == nil || app.SQLDB == nil || entityID == "" {
			return "", nil
		}
		var ownerID string
		if err := app.SQLDB.QueryRow(ctx, `SELECT userid FROM `+config.Tables.FarmsTable+` WHERE farmid = $1 LIMIT 1`, entityID).Scan(&ownerID); err == nil {
			return ownerID, nil
		}
		return "", nil
	default:
		return "", errors.New("invalid entity type")
	}
}

func findMerchByEntity(ctx context.Context, app *infra.Deps, _, entityID, merchID string) (Merch, error) {
	if app == nil || app.SQLDB == nil {
		return Merch{}, nil
	}
	var merch Merch
	row := app.SQLDB.QueryRow(ctx, `SELECT merchid, name, description, ownerid, price, created_at, updated_at, metadata FROM merch WHERE merchid = $1 AND ownerid = $2 LIMIT 1`, merchID, entityID)
	var createdAt, updatedAt time.Time
	var metadata []byte
	if err := row.Scan(&merch.MerchID, &merch.Name, &merch.Description, &merch.UserID, &merch.Price, &createdAt, &updatedAt, &metadata); err != nil {
		return Merch{}, err
	}
	merch.CreatedAt = createdAt
	merch.UpdatedAt = updatedAt
	if len(metadata) > 0 {
		_ = json.Unmarshal(metadata, &merch)
	}
	return merch, nil
}

func findMerchsByEntity(ctx context.Context, app *infra.Deps, _, entityID string) ([]Merch, error) {
	if app == nil || app.SQLDB == nil {
		return []Merch{}, nil
	}
	rows, err := app.SQLDB.Query(ctx, `SELECT merchid, name, description, ownerid, price, created_at, updated_at, metadata FROM merch WHERE ownerid = $1 ORDER BY created_at DESC`, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]Merch, 0)
	for rows.Next() {
		var merch Merch
		var createdAt, updatedAt time.Time
		var metadata []byte
		if err := rows.Scan(&merch.MerchID, &merch.Name, &merch.Description, &merch.UserID, &merch.Price, &createdAt, &updatedAt, &metadata); err != nil {
			return nil, err
		}
		merch.CreatedAt = createdAt
		merch.UpdatedAt = updatedAt
		if len(metadata) > 0 {
			_ = json.Unmarshal(metadata, &merch)
		}
		items = append(items, merch)
	}
	return items, rows.Err()
}

func findMerchByMerchID(ctx context.Context, app *infra.Deps, merchID string) (Merch, error) {
	if app == nil || app.SQLDB == nil {
		return Merch{}, nil
	}
	var merch Merch
	var createdAt, updatedAt time.Time
	var metadata []byte
	row := app.SQLDB.QueryRow(ctx, `SELECT merchid, name, description, ownerid, price, created_at, updated_at, metadata FROM merch WHERE merchid = $1 LIMIT 1`, merchID)
	if err := row.Scan(&merch.MerchID, &merch.Name, &merch.Description, &merch.UserID, &merch.Price, &createdAt, &updatedAt, &metadata); err != nil {
		return Merch{}, err
	}
	merch.CreatedAt = createdAt
	merch.UpdatedAt = updatedAt
	if len(metadata) > 0 {
		_ = json.Unmarshal(metadata, &merch)
	}
	return merch, nil
}

func findMerchForPurchase(ctx context.Context, app *infra.Deps, _, merchID string) (Merch, error) {
	return findMerchByMerchID(ctx, app, merchID)
}

func insertMerch(ctx context.Context, app *infra.Deps, merch Merch) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	metadata, err := json.Marshal(merch)
	if err != nil {
		return err
	}
	_, err = app.SQLDB.Exec(ctx, `INSERT INTO merch (merchid, name, description, ownerid, price, created_at, updated_at, metadata) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT (merchid) DO UPDATE SET name = EXCLUDED.name, description = EXCLUDED.description, ownerid = EXCLUDED.ownerid, price = EXCLUDED.price, updated_at = NOW(), metadata = EXCLUDED.metadata`, merch.MerchID, merch.Name, merch.Description, merch.UserID, merch.Price, merch.CreatedAt, time.Now().UTC(), metadata)
	return err
}

func updateMerchFields(ctx context.Context, app *infra.Deps, _, entityID, merchID string, update map[string]any) error {
	if app == nil || app.SQLDB == nil || len(update) == 0 {
		return nil
	}
	setParts := make([]string, 0, len(update))
	args := make([]any, 0, len(update)+2)
	idx := 1
	for key, value := range update {
		setParts = append(setParts, key+" = $"+fmt.Sprint(idx))
		args = append(args, value)
		idx++
	}
	args = append(args, merchID)
	args = append(args, entityID)
	_, err := app.SQLDB.Exec(ctx, `UPDATE merch SET `+strings.Join(setParts, ", ")+`, updated_at = NOW() WHERE merchid = $`+fmt.Sprint(len(args)-1)+` AND ownerid = $`+fmt.Sprint(len(args)), args...)
	return err
}

func softDeleteMerch(ctx context.Context, app *infra.Deps, _, entityID, merchID string, now time.Time) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	_, err := app.SQLDB.Exec(ctx, `UPDATE merch SET metadata = COALESCE(metadata, '{}'::jsonb) || '{"deleted": true, "deleted_at": "`+now.Format(time.RFC3339)+`"}'::jsonb, updated_at = NOW() WHERE merchid = $1 AND ownerid = $2`, merchID, entityID)
	return err
}

func confirmMerchPurchase(ctx context.Context, app *infra.Deps, eventID, merchID string, quantity int) (Merch, error) {
	if app == nil || app.SQLDB == nil {
		return Merch{}, nil
	}
	if quantity <= 0 {
		return Merch{}, errors.New("quantity must be greater than zero")
	}
	current, err := findMerchByMerchID(ctx, app, merchID)
	if err != nil {
		return Merch{}, err
	}
	if current.Stock < quantity {
		return Merch{}, errors.New("not enough merch available")
	}
	if eventID != "" && current.EntityID != "" && current.EntityID != eventID {
		return Merch{}, errors.New("merch does not belong to the requested event")
	}
	updated := current
	updated.Stock = current.Stock - quantity
	updated.UpdatedAt = time.Now().UTC()
	result, err := app.SQLDB.Exec(
		ctx,
		`UPDATE merch SET stock = $1, updated_at = NOW() WHERE merchid = $2 AND stock >= $3`,
		updated.Stock,
		merchID,
		quantity,
	)
	if err != nil {
		return Merch{}, err
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return Merch{}, errors.New("not enough merch available")
	}
	return updated, nil
}

func buyMerchTransaction(ctx context.Context, app *infra.Deps, r context.Context, userID, entityType, entityID, merchID string, quantity int) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	if userID == "" || merchID == "" {
		return errors.New("invalid purchase request")
	}
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}
	_ = r
	_ = entityType
	_ = entityID
	current, err := findMerchByMerchID(ctx, app, merchID)
	if err != nil {
		return err
	}
	if current.Stock < quantity {
		return errors.New("not enough merch available")
	}
	_, err = app.SQLDB.Exec(
		ctx,
		`UPDATE merch SET stock = stock - $1, updated_at = NOW() WHERE merchid = $2 AND stock >= $1`,
		quantity,
		merchID,
	)
	if err != nil {
		return err
	}
	return nil
}
