// File: internal/reports/reportsSQLDB.go

package reports

import (
	"context"

	"scav/config"
	"scav/infra"
)

/* -------------------------
   Tables
------------------------- */

var (
	reportsTable       = config.Tables.ReportsTable
	appealsTable       = config.Tables.AppealsTable
	moderatorAppsTable = config.Tables.ModeratorApplications
)

/* -------------------------
   DB Wrappers
------------------------- */

func FindReportByFilter(ctx context.Context, app *infra.Deps, query string, args []any, out *Report) error {
}

func InsertReport(ctx context.Context, app *infra.Deps, report Report) error {
}

func FindReports(ctx context.Context, app *infra.Deps, query string, args []any, out *[]Report) error {
}

func UpdateReportByID(ctx context.Context, app *infra.Deps, reportID string, update map[string]any) (int64, error) {
}

// Appeals
func FindAppealByFilter(ctx context.Context, app *infra.Deps, query string, args []any, out any) error {
}

func InsertAppeal(ctx context.Context, app *infra.Deps, appeal any) error {
}

func FindAppeals(ctx context.Context, app *infra.Deps, query string, args []any, out any) error {
}

func GetAppealByID(ctx context.Context, app *infra.Deps, appealID string, out any) error {
}

func UpdateAppealByID(ctx context.Context, app *infra.Deps, appealID string, update map[string]any) (int64, error) {
}

func setEntityDeletedFlagInDB(ctx context.Context, app *infra.Deps, table, idField, id string, deleted bool, by string) error {

	return err
}
