package merch

import (
	"context"
	"errors"
	"fmt"
	"time"

	"scav/config"
	"scav/infra"
	"scav/internal/beats/userdata"
)

var merchTable = config.Tables.MerchTable

func SQLgetEntityOwner(ctx context.Context, app *infra.Deps, entityType, entityID string) (string, error) {
	table := ""
	idField := ""
	ownerField := ""

	switch entityType {
	case "event":
		table = config.Tables.EventsTable
		idField = "eventid"
		ownerField = "creatorid"
	case "farm":
		table = config.Tables.FarmsTable
		idField = "farmid"
		ownerField = "createdBy"
	case "artist":
		table = config.Tables.ArtistsTable
		idField = "artistid"
		ownerField = "creatorid"
	default:
		return "", errors.New("invalid entity type")
	}

	query := fmt.Sprintf("%s = $1", idField)
	args := []any{entityID}

	var ownerEntity map[string]any
	if err := app.SQLDB.FindOne(ctx, table, query, args, &ownerEntity); err != nil {
		return "", err
	}

	owner, ok := ownerEntity[ownerField].(string)
	if !ok {
		return "", errors.New("cannot verify ownership")
	}
	return owner, nil
}

func SQLfindMerchByEntity(ctx context.Context, app *infra.Deps, entityType, entityID, merchID string) (Merch, error) {
	var merch Merch
	query := "entity_type = $1 AND entity_id = $2 AND merchid = $3 AND deletedat IS NULL"
	args := []any{entityType, entityID, merchID}

	if err := app.SQLDB.FindOne(ctx, merchTable, query, args, &merch); err != nil {
		return Merch{}, err
	}
	return merch, nil
}

func SQLfindMerchsByEntity(ctx context.Context, app *infra.Deps, entityType, entityID string) ([]Merch, error) {
	var list []Merch
	query := "entity_type = $1 AND entity_id = $2 AND deletedat IS NULL"
	args := []any{entityType, entityID}

	if err := app.SQLDB.FindMany(ctx, merchTable, query, args, &list); err != nil {
		return nil, err
	}
	if list == nil {
		return []Merch{}, nil
	}
	return list, nil
}

func SQLfindMerchByMerchID(ctx context.Context, app *infra.Deps, merchID string) (Merch, error) {
	var merch Merch
	query := "merchid = $1 AND deletedat IS NULL"
	args := []any{merchID}

	if err := app.SQLDB.FindOne(ctx, merchTable, query, args, &merch); err != nil {
		return Merch{}, err
	}
	return merch, nil
}

func SQLfindMerchForPurchase(ctx context.Context, app *infra.Deps, eventID, merchID string) (Merch, error) {
	var merch Merch
	query := "entity_id = $1 AND merchid = $2"
	args := []any{eventID, merchID}

	if err := app.SQLDB.FindOne(ctx, merchTable, query, args, &merch); err != nil {
		return Merch{}, err
	}
	return merch, nil
}

func SQLinsertMerch(ctx context.Context, app *infra.Deps, merch Merch) error {
	return app.SQLDB.Insert(ctx, merchTable, merch)
}

func SQLupdateMerchFields(ctx context.Context, app *infra.Deps, entityType, entityID, merchID string, update map[string]any) error {
	query := "entity_type = $1 AND entity_id = $2 AND merchid = $3"
	args := []any{entityType, entityID, merchID}

	_, err := app.SQLDB.UpdateOne(ctx, merchTable, query, args, update)
	return err
}

func SQLsoftDeleteMerch(ctx context.Context, app *infra.Deps, entityType, entityID, merchID string, now time.Time) error {
	query := "entity_type = $1 AND entity_id = $2 AND merchid = $3 AND deletedat IS NULL"
	args := []any{entityType, entityID, merchID}

	update := map[string]any{
		"deletedat": now,
		"updatedat": now,
	}

	_, err := app.SQLDB.UpdateOne(ctx, merchTable, query, args, update)
	return err
}

func SQLconfirmMerchPurchase(ctx context.Context, app *infra.Deps, eventID, merchID string, quantity int) (Merch, error) {
	query := "entity_id = $1 AND merchid = $2 AND stock >= $3"
	args := []any{eventID, merchID, quantity}

	err := app.SQLDB.Inc(ctx, merchTable, query, args, "stock", int64(-quantity))
	if err != nil {
		return Merch{}, err
	}

	var updatedMerch Merch
	findQuery := "entity_id = $1 AND merchid = $2"
	findArgs := []any{eventID, merchID}
	if err := app.SQLDB.FindOne(ctx, merchTable, findQuery, findArgs, &updatedMerch); err != nil {
		return Merch{}, err
	}
	return updatedMerch, nil
}

func SQLbuyMerchTransaction(ctx context.Context, app *infra.Deps, r context.Context, userID, entityType, entityID, merchID string, quantity int) error {
	_ = ctx
	return app.SQLDB.WithDB(r, func(txCtx context.Context) error {
		var merch Merch
		query := "entity_type = $1 AND entity_id = $2 AND merchid = $3"
		args := []any{entityType, entityID, merchID}

		err := app.SQLDB.FindOne(txCtx, merchTable, query, args, &merch)
		if err != nil {
			return errors.New("merch not found")
		}

		if merch.Stock < quantity {
			return errors.New("insufficient stock")
		}

		err = app.SQLDB.Inc(txCtx, merchTable, "merchid = $1", []any{merch.MerchID}, "stock", int64(-quantity))
		if err != nil {
			return err
		}

		userdata.SetUserData(
			"merch",
			merch.MerchID,
			userID,
			merch.EntityType,
			merch.EntityID,
			app,
		)

		return nil
	})
}
