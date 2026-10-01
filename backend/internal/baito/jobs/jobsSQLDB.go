// File: internal/baito/jobs/jobsSQLDB.go

package jobs

import (
	"context"

	"scav/config"
	"scav/infra"
	"scav/internal/baito"
)

var baitosTable = config.Tables.BaitoTable

func InsertBaitoForEntity(ctx context.Context, app *infra.Deps, baito baito.Baito) error {
	return app.SQLDB.Insert(ctx, baitosTable, baito)
}

func FindJobsForEntity(ctx context.Context, app *infra.Deps, entityType, entityID string) ([]baito.BaitosResponse, error) {
	var jobs []baito.BaitosResponse
	query := "entitytype = $1 AND entityid = $2"
	args := []any{entityType, entityID}
	err := app.SQLDB.FindMany(ctx, baitosTable, query, args, &jobs)
	return jobs, err
}
