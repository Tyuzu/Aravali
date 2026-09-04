package comments

import (
	"context"

	"scav/config"
	db "scav/infra/db"
)

var commentsCollection = config.Collections.CommentsCollection

func insertComment(ctx context.Context, database db.Database, comment Comment) error {
	return database.Insert(ctx, commentsCollection, comment)
}

func findCommentByID(ctx context.Context, database db.Database, commentID string, comment *Comment) error {
	return database.FindOne(ctx, commentsCollection, map[string]any{"commentid": commentID}, comment)
}

func updateCommentContent(ctx context.Context, database db.Database, commentID string, update map[string]any) (any, error) {
	return database.UpdateOne(ctx, commentsCollection, map[string]any{"commentid": commentID}, update)
}

func deleteComment(ctx context.Context, database db.Database, commentID, userID string) (int64, error) {
	return database.Delete(ctx, commentsCollection, map[string]any{"commentid": commentID, "createdby": userID})
}

func findCommentsByEntity(
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
	return database.FindManyWithOptions(ctx, commentsCollection, filter, opts, comments)
}
