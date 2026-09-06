package jobs

import (
	"context"

	"scav/config"
	"scav/infra"
	"scav/internal/baito"
)

var baitosCollection = config.Collections.BaitoCollection

func InsertBaitoForEntity(ctx context.Context, app *infra.Deps, baito baito.Baito) error {
	return app.DB.Insert(ctx, baitosCollection, baito)
}

func FindJobsForEntity(ctx context.Context, app *infra.Deps, entityType, entityID string) ([]baito.BaitosResponse, error) {
	var jobs []baito.BaitosResponse
	err := app.DB.FindMany(ctx, baitosCollection, map[string]any{
		"entityType": entityType,
		"entityId":   entityID,
	}, &jobs)
	return jobs, err
}
