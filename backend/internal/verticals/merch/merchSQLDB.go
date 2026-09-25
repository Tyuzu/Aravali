package merch

import (
	"context"
	"errors"
	"scav/config"
	"scav/infra"
	"scav/internal/beats/userdata"
	"time"
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

	var ownerEntity map[string]any
	if err := app.DB.FindOne(ctx, table, map[string]any{idField: entityID}, &ownerEntity); err != nil {
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
	if err := app.DB.FindOne(ctx, merchTable, map[string]any{
		"entity_type": entityType,
		"entity_id":   entityID,
		"merchid":     merchID,
		"deletedAt":   map[string]any{"$exists": false},
	}, &merch); err != nil {
		return Merch{}, err
	}
	return merch, nil
}

func SQLfindMerchsByEntity(ctx context.Context, app *infra.Deps, entityType, entityID string) ([]Merch, error) {
	var list []Merch
	if err := app.DB.FindMany(ctx, merchTable, map[string]any{
		"entity_type": entityType,
		"entity_id":   entityID,
		"deletedAt":   map[string]any{"$exists": false},
	}, &list); err != nil {
		return nil, err
	}
	if list == nil {
		return []Merch{}, nil
	}
	return list, nil
}

func SQLfindMerchByMerchID(ctx context.Context, app *infra.Deps, merchID string) (Merch, error) {
	var merch Merch
	if err := app.DB.FindOne(ctx, merchTable, map[string]any{
		"merchid":   merchID,
		"deletedAt": map[string]any{"$exists": false},
	}, &merch); err != nil {
		return Merch{}, err
	}
	return merch, nil
}

func SQLfindMerchForPurchase(ctx context.Context, app *infra.Deps, eventID, merchID string) (Merch, error) {
	var merch Merch
	if err := app.DB.FindOne(ctx, merchTable, map[string]any{
		"entity_id": eventID,
		"merchid":   merchID,
	}, &merch); err != nil {
		return Merch{}, err
	}
	return merch, nil
}

func SQLinsertMerch(ctx context.Context, app *infra.Deps, merch Merch) error {
	return app.DB.Insert(ctx, merchTable, merch)
}

func SQLupdateMerchFields(ctx context.Context, app *infra.Deps, entityType, entityID, merchID string, update map[string]any) error {
	_, err := app.DB.UpdateOne(ctx, merchTable, map[string]any{
		"entity_type": entityType,
		"entity_id":   entityID,
		"merchid":     merchID,
	}, map[string]any{"$set": update})
	return err
}

func SQLsoftDeleteMerch(ctx context.Context, app *infra.Deps, entityType, entityID, merchID string, now time.Time) error {
	_, err := app.DB.UpdateOne(ctx, merchTable,
		map[string]any{
			"entity_type": entityType,
			"entity_id":   entityID,
			"merchid":     merchID,
			"deletedAt":   map[string]any{"$exists": false},
		},
		map[string]any{"$set": map[string]any{
			"deletedAt": now,
			"updatedat": now,
		}},
	)
	return err
}

func SQLconfirmMerchPurchase(ctx context.Context, app *infra.Deps, eventID, merchID string, quantity int) (Merch, error) {
	var updatedMerch Merch
	err := app.DB.FindOneAndUpdate(
		ctx,
		merchTable,
		map[string]any{
			"entity_id": eventID,
			"merchid":   merchID,
			"stock":     map[string]any{"$gte": quantity},
		},
		map[string]any{
			"$inc": map[string]any{"stock": -quantity},
		},
		&updatedMerch,
	)
	if err != nil {
		return Merch{}, err
	}
	return updatedMerch, nil
}

func SQLbuyMerchTransaction(ctx context.Context, app *infra.Deps, r context.Context, userID, entityType, entityID, merchID string, quantity int) error {
	_ = ctx
	return app.DB.WithDB(r, func(txCtx context.Context) error {
		var merch Merch
		err := app.DB.FindOne(txCtx, merchTable, map[string]any{
			"entity_type": entityType,
			"entity_id":   entityID,
			"merchid":     merchID,
		}, &merch)
		if err != nil {
			return errors.New("merch not found")
		}

		if merch.Stock < quantity {
			return errors.New("insufficient stock")
		}

		_, err = app.DB.UpdateOne(
			txCtx,
			merchTable,
			map[string]any{"merchid": merch.MerchID},
			map[string]any{"$inc": map[string]any{"stock": -quantity}},
		)
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
