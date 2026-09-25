package jobs

import (
	"context"

	"scav/config"
	"scav/infra"
	"scav/internal/baito"
)

var baitosTable = config.Tables.BaitoTable

func SQLInsertBaitoForEntity(ctx context.Context, app *infra.Deps, baito baito.Baito) error {
	return app.DB.Insert(ctx, baitosTable, baito)
}

func SQLFindJobsForEntity(ctx context.Context, app *infra.Deps, entityType, entityID string) ([]baito.BaitosResponse, error) {
	var jobs []baito.BaitosResponse
	err := app.DB.FindMany(ctx, baitosTable, map[string]any{
		"entityType": entityType,
		"entityId":   entityID,
	}, &jobs)
	return jobs, err
}
