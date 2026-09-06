package reports

import (
	"context"
	"time"

	"scav/config"
	"scav/infra"
)

/* -------------------------
   Collections
------------------------- */

var (
	reportsCollection       = config.Collections.ReportsCollection
	appealsCollection       = config.Collections.AppealsCollection
	moderatorAppsCollection = config.Collections.ModeratorApplications
)

/* -------------------------
   DB Wrappers
------------------------- */

func FindReportByFilter(ctx context.Context, app *infra.Deps, filter map[string]any, out *Report) error {
	return app.DB.FindOne(ctx, reportsCollection, filter, out)
}

func InsertReport(ctx context.Context, app *infra.Deps, report Report) error {
	return app.DB.Insert(ctx, reportsCollection, report)
}

func FindReports(ctx context.Context, app *infra.Deps, filter map[string]any, out *[]Report) error {
	return app.DB.FindMany(ctx, reportsCollection, filter, out)
}

func UpdateReportByID(ctx context.Context, app *infra.Deps, reportID string, update map[string]any) (any, error) {
	return app.DB.Update(ctx, reportsCollection, map[string]any{"reportid": reportID}, update)
}

// Appeals
func FindAppealByFilter(ctx context.Context, app *infra.Deps, filter map[string]any, out any) error {
	return app.DB.FindOne(ctx, appealsCollection, filter, out)
}

func InsertAppeal(ctx context.Context, app *infra.Deps, appeal any) error {
	return app.DB.Insert(ctx, appealsCollection, appeal)
}

func FindAppeals(ctx context.Context, app *infra.Deps, filter map[string]any, out any) error {
	return app.DB.FindMany(ctx, appealsCollection, filter, out)
}

func GetAppealByID(ctx context.Context, app *infra.Deps, appealID string, out any) error {
	return app.DB.FindOne(ctx, appealsCollection, map[string]any{"appealid": appealID}, out)
}

func UpdateAppealByID(ctx context.Context, app *infra.Deps, appealID string, update map[string]any) (any, error) {
	return app.DB.Update(ctx, appealsCollection, map[string]any{"appealid": appealID}, update)
}

func setEntityDeletedFlagInDB(ctx context.Context, app *infra.Deps, collection, idField, id string, deleted bool, by string) error {
	now := time.Now().UTC()
	deletedAtVal := interface{}("")
	if deleted {
		deletedAtVal = now
	}

	_, err := app.DB.Update(
		ctx,
		collection,
		map[string]any{idField: id},
		map[string]any{
			"deleted":   deleted,
			"deletedBy": by,
			"deletedAt": deletedAtVal,
		},
	)
	return err
}
