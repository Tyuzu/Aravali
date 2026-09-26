package notices

import (
	"context"

	"scav/config"
	"scav/infra"
	"scav/infra/sqldb"
)

var noticesTable = config.Tables.NoticesTable

func SQLcreateNotice(ctx context.Context, app *infra.Deps, notice Notice) error {
	return app.SQLDB.Insert(ctx, noticesTable, notice)
}

func SQLfindNoticeByID(ctx context.Context, app *infra.Deps, noticeID string) (Notice, error) {
	var notice Notice
	query := "noticeid = $1"
	args := []any{noticeID}
	err := app.SQLDB.FindOne(ctx, noticesTable, query, args, &notice)
	return notice, err
}

func SQLlistNoticesWithOptions(ctx context.Context, app *infra.Deps, query string, args []any, opts sqldb.FindManyOptions, out *[]Notice) error {
	return app.SQLDB.FindManyWithOptions(ctx, noticesTable, query, args, opts, out)
}

func SQLupdateNoticeByID(ctx context.Context, app *infra.Deps, noticeID string, update map[string]any) error {
	query := "noticeid = $1"
	args := []any{noticeID}
	_, err := app.SQLDB.Update(ctx, noticesTable, query, args, update)
	return err
}

func SQLdeleteNoticeByID(ctx context.Context, app *infra.Deps, noticeID string) error {
	query := "noticeid = $1"
	args := []any{noticeID}
	_, err := app.SQLDB.Delete(ctx, noticesTable, query, args)
	return err
}
