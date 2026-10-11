// File: internal/baito/jobs/jobsSQLDB.go

package jobs

import (
	"context"
	"fmt"
	"time"

	"scav/config"
	"scav/infra"
	"scav/internal/baito"
)

var baitosTable = config.Tables.BaitoTable

func InsertBaitoForEntity(ctx context.Context, app *infra.Deps, baito baito.Baito) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	if baito.CreatedAt.IsZero() {
		baito.CreatedAt = time.Now().UTC()
	}
	if baito.UpdatedAt.IsZero() {
		baito.UpdatedAt = baito.CreatedAt
	}
	_, err := app.SQLDB.Exec(ctx,
		fmt.Sprintf("INSERT INTO %s (baitoid, entitytype, entityid, title, description, category, subcategory, location, wage, phone, requirements, banner, images, workhours, benefits, email, tags, duration, lastdate, createdat, updatedat, ownerid, applicationcount) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23)", baitosTable),
		baito.BaitoId, baito.EntityType, baito.EntityID, baito.Title, baito.Description, baito.Category, baito.SubCategory, baito.Location, baito.Wage, baito.Phone, baito.Requirements, baito.Banner, baito.Images, baito.WorkHours, baito.Benefits, baito.Email, baito.Tags, baito.Duration, baito.LastDateToApply, baito.CreatedAt, baito.UpdatedAt, baito.OwnerID, baito.ApplicationCount,
	)
	return err
}

func FindJobsForEntity(ctx context.Context, app *infra.Deps, entityType, entityID string) ([]baito.BaitosResponse, error) {
	if app == nil || app.SQLDB == nil {
		return nil, nil
	}
	if entityType == "" || entityID == "" {
		return []baito.BaitosResponse{}, nil
	}
	rows, err := app.SQLDB.Query(ctx, fmt.Sprintf("SELECT baitoid, title, description, category, subcategory, location, wage, requirements, banner, workhours, duration, lastdate, createdat, ownerid FROM %s WHERE entitytype = $1 AND entityid = $2 ORDER BY createdat DESC", baitosTable), entityType, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return []baito.BaitosResponse{}, nil
}
