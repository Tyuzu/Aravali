package notices

import (
	"context"

	"scav/config"
	"scav/infra"
	"scav/infra/db"
)

var noticesTable = config.Tables.NoticesTable

func SQLcreateNotice(ctx context.Context, app *infra.Deps, notice Notice) error {
	return app.DB.Insert(ctx, noticesTable, notice)
}

func SQLfindNoticeByID(ctx context.Context, app *infra.Deps, noticeID string) (Notice, error) {
	var notice Notice
	err := app.DB.FindOne(ctx, noticesTable, map[string]any{"noticeid": noticeID}, &notice)
	return notice, err
}

func SQLlistNoticesWithOptions(ctx context.Context, app *infra.Deps, filter map[string]any, opts db.FindManyOptions, out *[]Notice) error {
	return app.DB.FindManyWithOptions(ctx, noticesTable, filter, opts, out)
}

func SQLupdateNoticeByID(ctx context.Context, app *infra.Deps, noticeID string, update map[string]any) error {
	_, err := app.DB.Update(ctx, noticesTable, map[string]any{"noticeid": noticeID}, update)
	return err
}

func SQLdeleteNoticeByID(ctx context.Context, app *infra.Deps, noticeID string) error {
	_, err := app.DB.Delete(ctx, noticesTable, map[string]any{"noticeid": noticeID})
	return err
}
