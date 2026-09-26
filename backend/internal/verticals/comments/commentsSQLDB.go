package comments

import (
	"context"

	"scav/config"
	"scav/infra"
	"scav/infra/sqldb"
)

var commentsTable = config.Tables.CommentsTable

func SQLinsertComment(ctx context.Context, app *infra.Deps, comment Comment) error {
	return app.SQLDB.Insert(ctx, commentsTable, comment)
}

func SQLfindCommentByID(ctx context.Context, app *infra.Deps, commentID string, comment *Comment) error {
	query := "commentid = $1"
	args := []any{commentID}

	return app.SQLDB.FindOne(ctx, commentsTable, query, args, comment)
}

func SQLupdateCommentContent(ctx context.Context, app *infra.Deps, commentID string, update map[string]any) (int64, error) {
	query := "commentid = $1"
	args := []any{commentID}

	return app.SQLDB.UpdateOne(ctx, commentsTable, query, args, update)
}

func SQLdeleteComment(ctx context.Context, app *infra.Deps, commentID, userID string) (int64, error) {
	query := "commentid = $1 AND createdby = $2"
	args := []any{commentID, userID}

	return app.SQLDB.Delete(ctx, commentsTable, query, args)
}

func SQLfindCommentsByEntity(
	ctx context.Context,
	app *infra.Deps,
	entityType string,
	entityID string,
	opts sqldb.FindManyOptions,
	comments *[]Comment,
) error {
	query := "entity_type = $1 AND entity_id = $2"
	args := []any{entityType, entityID}

	return app.SQLDB.FindManyWithOptions(ctx, commentsTable, query, args, opts, comments)
}
