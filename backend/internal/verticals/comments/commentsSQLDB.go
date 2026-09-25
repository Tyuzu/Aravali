package comments

import (
	"context"

	"scav/config"
	db "scav/infra/db"
)

var commentsTable = config.Tables.CommentsTable

func SQLinsertComment(ctx context.Context, database db.Database, comment Comment) error {
	return database.Insert(ctx, commentsTable, comment)
}

func SQLfindCommentByID(ctx context.Context, database db.Database, commentID string, comment *Comment) error {
	return database.FindOne(ctx, commentsTable, map[string]any{"commentid": commentID}, comment)
}

func SQLupdateCommentContent(ctx context.Context, database db.Database, commentID string, update map[string]any) (any, error) {
	return database.UpdateOne(ctx, commentsTable, map[string]any{"commentid": commentID}, update)
}

func SQLdeleteComment(ctx context.Context, database db.Database, commentID, userID string) (int64, error) {
	return database.Delete(ctx, commentsTable, map[string]any{"commentid": commentID, "createdby": userID})
}

func SQLfindCommentsByEntity(
	ctx context.Context,
	database db.Database,
	entityType string,
	entityID string,
	opts db.FindManyOptions,
	comments *[]Comment,
) error {
	filter := map[string]any{
		"entity_type": entityType,
		"entity_id":   entityID,
	}
	return database.FindManyWithOptions(ctx, commentsTable, filter, opts, comments)
}
