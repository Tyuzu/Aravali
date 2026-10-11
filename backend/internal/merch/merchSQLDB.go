// File: internal/merch/merchSQLDB.go

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
	switch entityType {
	case "event", "farm", "artist":
		return "", nil
	default:
		return "", errors.New("invalid entity type")
	}
}

func findMerchByEntity(ctx context.Context, app *infra.Deps, entityType, entityID, merchID string) (Merch, error) {
	return Merch{}, nil
}

func findMerchsByEntity(ctx context.Context, app *infra.Deps, entityType, entityID string) ([]Merch, error) {
	return nil, nil
}

func findMerchByMerchID(ctx context.Context, app *infra.Deps, merchID string) (Merch, error) {
	return Merch{}, nil
}

func findMerchForPurchase(ctx context.Context, app *infra.Deps, eventID, merchID string) (Merch, error) {
	return Merch{}, nil
}

func insertMerch(ctx context.Context, app *infra.Deps, merch Merch) error {
	return nil
}

func updateMerchFields(ctx context.Context, app *infra.Deps, entityType, entityID, merchID string, update map[string]any) error {
	return nil
}

func softDeleteMerch(ctx context.Context, app *infra.Deps, entityType, entityID, merchID string, now time.Time) error {
	return nil
}

func confirmMerchPurchase(ctx context.Context, app *infra.Deps, eventID, merchID string, quantity int) (Merch, error) {
	return Merch{}, nil
}

func buyMerchTransaction(ctx context.Context, app *infra.Deps, r context.Context, userID, entityType, entityID, merchID string, quantity int) error {
	return nil
}
