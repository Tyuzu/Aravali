package reviews

import (
	"context"
	"scav/config"
	"scav/infra"

	"go.mongodb.org/mongo-driver/bson"
)

/* -------------------------
   Collection
------------------------- */

var reviewsCollection = config.Collections.ReviewsCollection

/* -------------------------
   Database Helpers
------------------------- */

func GetUserReviewForEntity(ctx context.Context, app *infra.Deps, userID, entityType, entityID string) (*Review, error) {
	filter := bson.M{
		"userid":     userID,
		"entityType": entityType,
		"entityId":   entityID,
	}

	var review Review
	if err := app.DB.FindOne(ctx, reviewsCollection, filter, &review); err != nil {
		return nil, err
	}
	return &review, nil
}

func GetReviewByID(ctx context.Context, app *infra.Deps, reviewID string) (*Review, error) {
	var review Review
	if err := app.DB.FindOne(ctx, reviewsCollection, bson.M{"reviewid": reviewID}, &review); err != nil {
		return nil, err
	}
	return &review, nil
}

func InsertReview(ctx context.Context, app *infra.Deps, review Review) error {
	return app.DB.Insert(ctx, reviewsCollection, review)
}

func UpdateReviewByID(ctx context.Context, app *infra.Deps, reviewID string, update bson.M) error {
	_, err := app.DB.Update(ctx, reviewsCollection, bson.M{"reviewid": reviewID}, update)
	return err
}

func DeleteReviewByID(ctx context.Context, app *infra.Deps, reviewID string) error {
	_, err := app.DB.Delete(ctx, reviewsCollection, bson.M{"reviewid": reviewID})
	return err
}

/* -------------------------
   Database Helpers
------------------------- */

// FetchReviewsByEntity queries reviews for a specific entity type and ID.
func FetchReviewsByEntity(ctx context.Context, app *infra.Deps, entityType, entityID string) ([]Review, error) {
	filter := bson.M{
		"entityType": entityType,
		"entityId":   entityID,
	}

	var reviews []Review
	if err := app.DB.FindMany(ctx, reviewsCollection, filter, &reviews); err != nil {
		return nil, err
	}

	return reviews, nil
}

// FetchReviewByID retrieves a single review by its ID.
func FetchReviewByID(ctx context.Context, app *infra.Deps, reviewID string) (*Review, error) {
	var review Review
	if err := app.DB.FindOne(ctx, reviewsCollection, bson.M{"reviewid": reviewID}, &review); err != nil {
		return nil, err
	}
	return &review, nil
}
