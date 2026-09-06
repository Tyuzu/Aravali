package merch

import (
	"context"
	"errors"
	"scav/config"
	"scav/infra"
	"scav/internal/userdata"
	"time"
)

var merchCollection = config.Collections.MerchCollection

func getEntityOwner(ctx context.Context, app *infra.Deps, entityType, entityID string) (string, error) {
	collection := ""
	idField := ""
	ownerField := ""

	switch entityType {
	case "event":
		collection = "events"
		idField = "eventid"
		ownerField = "creatorid"
	case "farm":
		collection = "farms"
		idField = "farmid"
		ownerField = "createdBy"
	case "artist":
		collection = "artists"
		idField = "artistid"
		ownerField = "creatorid"
	default:
		return "", errors.New("invalid entity type")
	}

	var ownerEntity map[string]any
	if err := app.DB.FindOne(ctx, collection, map[string]any{idField: entityID}, &ownerEntity); err != nil {
		return "", err
	}

	owner, ok := ownerEntity[ownerField].(string)
	if !ok {
		return "", errors.New("cannot verify ownership")
	}
	return owner, nil
}

func findMerchByEntity(ctx context.Context, app *infra.Deps, entityType, entityID, merchID string) (Merch, error) {
	var merch Merch
	if err := app.DB.FindOne(ctx, merchCollection, map[string]any{
		"entity_type": entityType,
		"entity_id":   entityID,
		"merchid":     merchID,
		"deletedAt":   map[string]any{"$exists": false},
	}, &merch); err != nil {
		return Merch{}, err
	}
	return merch, nil
}

func findMerchsByEntity(ctx context.Context, app *infra.Deps, entityType, entityID string) ([]Merch, error) {
	var list []Merch
	if err := app.DB.FindMany(ctx, merchCollection, map[string]any{
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

func findMerchByMerchID(ctx context.Context, app *infra.Deps, merchID string) (Merch, error) {
	var merch Merch
	if err := app.DB.FindOne(ctx, merchCollection, map[string]any{
		"merchid":   merchID,
		"deletedAt": map[string]any{"$exists": false},
	}, &merch); err != nil {
		return Merch{}, err
	}
	return merch, nil
}

func findMerchForPurchase(ctx context.Context, app *infra.Deps, eventID, merchID string) (Merch, error) {
	var merch Merch
	if err := app.DB.FindOne(ctx, merchCollection, map[string]any{
		"entity_id": eventID,
		"merchid":   merchID,
	}, &merch); err != nil {
		return Merch{}, err
	}
	return merch, nil
}

func insertMerch(ctx context.Context, app *infra.Deps, merch Merch) error {
	return app.DB.Insert(ctx, merchCollection, merch)
}

func updateMerchFields(ctx context.Context, app *infra.Deps, entityType, entityID, merchID string, update map[string]any) error {
	_, err := app.DB.UpdateOne(ctx, merchCollection, map[string]any{
		"entity_type": entityType,
		"entity_id":   entityID,
		"merchid":     merchID,
	}, map[string]any{"$set": update})
	return err
}

func softDeleteMerch(ctx context.Context, app *infra.Deps, entityType, entityID, merchID string, now time.Time) error {
	_, err := app.DB.UpdateOne(ctx, merchCollection,
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

func confirmMerchPurchase(ctx context.Context, app *infra.Deps, eventID, merchID string, quantity int) (Merch, error) {
	var updatedMerch Merch
	err := app.DB.FindOneAndUpdate(
		ctx,
		merchCollection,
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

func buyMerchTransaction(ctx context.Context, app *infra.Deps, r context.Context, userID, entityType, entityID, merchID string, quantity int) error {
	return app.DB.WithDB(r, func(txCtx context.Context) error {
		var merch Merch
		err := app.DB.FindOne(txCtx, merchCollection, map[string]any{
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
			merchCollection,
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
