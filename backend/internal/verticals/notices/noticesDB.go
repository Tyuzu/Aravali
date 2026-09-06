package notices

import (
	"context"

	"scav/config"
	"scav/infra"
	"scav/infra/db"
)

var noticesCollection = config.Collections.NoticesCollection

func createNotice(ctx context.Context, app *infra.Deps, notice Notice) error {
	return app.DB.Insert(ctx, noticesCollection, notice)
}

func findNoticeByID(ctx context.Context, app *infra.Deps, noticeID string) (Notice, error) {
	var notice Notice
	err := app.DB.FindOne(ctx, noticesCollection, map[string]any{"noticeid": noticeID}, &notice)
	return notice, err
}

func listNoticesWithOptions(ctx context.Context, app *infra.Deps, filter map[string]any, opts db.FindManyOptions, out *[]Notice) error {
	return app.DB.FindManyWithOptions(ctx, noticesCollection, filter, opts, out)
}

func updateNoticeByID(ctx context.Context, app *infra.Deps, noticeID string, update map[string]any) error {
	_, err := app.DB.Update(ctx, noticesCollection, map[string]any{"noticeid": noticeID}, update)
	return err
}

func deleteNoticeByID(ctx context.Context, app *infra.Deps, noticeID string) error {
	_, err := app.DB.Delete(ctx, noticesCollection, map[string]any{"noticeid": noticeID})
	return err
}
