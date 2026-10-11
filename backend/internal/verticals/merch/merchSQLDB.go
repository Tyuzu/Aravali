// File: internal/verticals/merch/merchSQLDB.go

package merch

import (
	"context"
	"errors"
	"time"

	"scav/config"
	"scav/infra"
)

var merchTable = config.Tables.MerchTable

func getEntityOwner(ctx context.Context, app *infra.Deps, entityType, entityID string) (string, error) {
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

	owner, ok := ownerEntity[ownerField].(string)
	if !ok {
		return "", errors.New("cannot verify ownership")
	}
	return owner, nil
}

func findMerchByEntity(ctx context.Context, app *infra.Deps, entityType, entityID, merchID string) (Merch, error) {
	var merch Merch

	return merch, nil
}

func findMerchsByEntity(ctx context.Context, app *infra.Deps, entityType, entityID string) ([]Merch, error) {
	var list []Merch

	return list, nil
}

func findMerchByMerchID(ctx context.Context, app *infra.Deps, merchID string) (Merch, error) {
	var merch Merch

	return merch, nil
}

func findMerchForPurchase(ctx context.Context, app *infra.Deps, eventID, merchID string) (Merch, error) {
	var merch Merch

	return merch, nil
}

func insertMerch(ctx context.Context, app *infra.Deps, merch Merch) error {

}

func updateMerchFields(ctx context.Context, app *infra.Deps, entityType, entityID, merchID string, update map[string]any) error {
	return err
}

func softDeleteMerch(ctx context.Context, app *infra.Deps, entityType, entityID, merchID string, now time.Time) error {
	return err
}

func confirmMerchPurchase(ctx context.Context, app *infra.Deps, eventID, merchID string, quantity int) (Merch, error) {

	var updatedMerch Merch

	return updatedMerch, nil
}

func buyMerchTransaction(ctx context.Context, app *infra.Deps, r context.Context, userID, entityType, entityID, merchID string, quantity int) error {

}
