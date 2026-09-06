package reviews

import (
	"context"

	"scav/config"
	"scav/infra"
)

/* -------------------------
   Collection
------------------------- */

var reviewsCollection = config.Collections.ReviewsCollection

/* -------------------------
   DB Wrappers
------------------------- */

func FindReviewByUserEntity(ctx context.Context, app *infra.Deps, userID, entityType, entityID string, out *Review) error {
	return app.SQLDB.FindOne(ctx, reviewsCollection, "userid = $1 AND entityType = $2 AND entityId = $3", []any{userID, entityType, entityID}, out)
}

func InsertReview(ctx context.Context, app *infra.Deps, review Review) error {
	return app.SQLDB.Insert(ctx, reviewsCollection, review)
}

func GetReviewsForEntity(ctx context.Context, app *infra.Deps, entityType, entityID string, out *[]Review) error {
	where := "entityType = $1 AND entityId = $2"
	args := []any{entityType, entityID}
	return app.SQLDB.FindMany(ctx, reviewsCollection, where, args, out)
}

func GetReviewByID(ctx context.Context, app *infra.Deps, reviewID string, out *Review) error {
	return app.SQLDB.FindOne(ctx, reviewsCollection, "reviewid = $1", []any{reviewID}, out)
}

func UpdateReviewByID(ctx context.Context, app *infra.Deps, reviewID string, update map[string]any) (int64, error) {
	return app.SQLDB.Update(ctx, reviewsCollection, "reviewid = $1", []any{reviewID}, update)
}

func DeleteReviewByID(ctx context.Context, app *infra.Deps, reviewID string) (int64, error) {
	return app.SQLDB.Delete(ctx, reviewsCollection, "reviewid = $1", []any{reviewID})
}
