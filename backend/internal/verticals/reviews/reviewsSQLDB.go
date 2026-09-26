package reviews

import (
	"context"
	"scav/config"
	"scav/infra"
)

/* -------------------------
   Table
------------------------- */

var reviewsTable = config.Tables.ReviewsTable

/* -------------------------
   Database Helpers
------------------------- */

func SQLGetUserReviewForEntity(ctx context.Context, app *infra.Deps, userID, entityType, entityID string) (*Review, error) {
	query := "userid = $1 AND entityType = $2 AND entityId = $3"
	args := []any{userID, entityType, entityID}

	var review Review
	if err := app.SQLDB.FindOne(ctx, reviewsTable, query, args, &review); err != nil {
		return nil, err
	}
	return &review, nil
}

func SQLGetReviewByID(ctx context.Context, app *infra.Deps, reviewID string) (*Review, error) {
	query := "reviewid = $1"
	args := []any{reviewID}

	var review Review
	if err := app.SQLDB.FindOne(ctx, reviewsTable, query, args, &review); err != nil {
		return nil, err
	}
	return &review, nil
}

func SQLInsertReview(ctx context.Context, app *infra.Deps, review Review) error {
	return app.SQLDB.Insert(ctx, reviewsTable, review)
}

func SQLUpdateReviewByID(ctx context.Context, app *infra.Deps, reviewID string, update map[string]any) error {
	query := "reviewid = $1"
	args := []any{reviewID}

	_, err := app.SQLDB.Update(ctx, reviewsTable, query, args, update)
	return err
}

func SQLDeleteReviewByID(ctx context.Context, app *infra.Deps, reviewID string) error {
	query := "reviewid = $1"
	args := []any{reviewID}

	_, err := app.SQLDB.Delete(ctx, reviewsTable, query, args)
	return err
}

/* -------------------------
   Fetch Helpers
------------------------- */

// FetchReviewsByEntity queries reviews for a specific entity type and ID.
func SQLFetchReviewsByEntity(ctx context.Context, app *infra.Deps, entityType, entityID string) ([]Review, error) {
	query := "entityType = $1 AND entityId = $2"
	args := []any{entityType, entityID}

	var reviews []Review
	if err := app.SQLDB.FindMany(ctx, reviewsTable, query, args, &reviews); err != nil {
		return nil, err
	}

	return reviews, nil
}

// FetchReviewByID retrieves a single review by its ID.
func SQLFetchReviewByID(ctx context.Context, app *infra.Deps, reviewID string) (*Review, error) {
	query := "reviewid = $1"
	args := []any{reviewID}

	var review Review
	if err := app.SQLDB.FindOne(ctx, reviewsTable, query, args, &review); err != nil {
		return nil, err
	}
	return &review, nil
}
