package reviews

import (
	"context"
	"scav/config"
	"scav/infra"

	"go.mongodb.org/mongo-driver/bson"
)

/* -------------------------
   Table
------------------------- */

var reviewsTable = config.Tables.ReviewsTable

/* -------------------------
   Database Helpers
------------------------- */

func SQLGetUserReviewForEntity(ctx context.Context, app *infra.Deps, userID, entityType, entityID string) (*Review, error) {
	filter := bson.M{
		"userid":     userID,
		"entityType": entityType,
		"entityId":   entityID,
	}

	var review Review
	if err := app.DB.FindOne(ctx, reviewsTable, filter, &review); err != nil {
		return nil, err
	}
	return &review, nil
}

func SQLGetReviewByID(ctx context.Context, app *infra.Deps, reviewID string) (*Review, error) {
	var review Review
	if err := app.DB.FindOne(ctx, reviewsTable, bson.M{"reviewid": reviewID}, &review); err != nil {
		return nil, err
	}
	return &review, nil
}

func SQLInsertReview(ctx context.Context, app *infra.Deps, review Review) error {
	return app.DB.Insert(ctx, reviewsTable, review)
}

func SQLUpdateReviewByID(ctx context.Context, app *infra.Deps, reviewID string, update bson.M) error {
	_, err := app.DB.Update(ctx, reviewsTable, bson.M{"reviewid": reviewID}, update)
	return err
}

func SQLDeleteReviewByID(ctx context.Context, app *infra.Deps, reviewID string) error {
	_, err := app.DB.Delete(ctx, reviewsTable, bson.M{"reviewid": reviewID})
	return err
}

/* -------------------------
   Database Helpers
------------------------- */

// FetchReviewsByEntity queries reviews for a specific entity type and ID.
func SQLFetchReviewsByEntity(ctx context.Context, app *infra.Deps, entityType, entityID string) ([]Review, error) {
	filter := bson.M{
		"entityType": entityType,
		"entityId":   entityID,
	}

	var reviews []Review
	if err := app.DB.FindMany(ctx, reviewsTable, filter, &reviews); err != nil {
		return nil, err
	}

	return reviews, nil
}

// FetchReviewByID retrieves a single review by its ID.
func SQLFetchReviewByID(ctx context.Context, app *infra.Deps, reviewID string) (*Review, error) {
	var review Review
	if err := app.DB.FindOne(ctx, reviewsTable, bson.M{"reviewid": reviewID}, &review); err != nil {
		return nil, err
	}
	return &review, nil
}
