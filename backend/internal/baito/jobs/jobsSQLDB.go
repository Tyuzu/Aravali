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

}

func FindJobsForEntity(ctx context.Context, app *infra.Deps, entityType, entityID string) ([]baito.BaitosResponse, error) {
	var jobs []baito.BaitosResponse
	return jobs, err
}
