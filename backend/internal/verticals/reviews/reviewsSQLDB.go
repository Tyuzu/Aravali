// File: internal/verticals/reviews/reviewsSQLDB.go

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

func GetUserReviewForEntity(ctx context.Context, app *infra.Deps, userID, entityType, entityID string) (*Review, error) {

	var review Review

	return &review, nil
}

func GetReviewByID(ctx context.Context, app *infra.Deps, reviewID string) (*Review, error) {

	var review Review

	return &review, nil
}

func InsertReview(ctx context.Context, app *infra.Deps, review Review) error {

}

func UpdateReviewByID(ctx context.Context, app *infra.Deps, reviewID string, update map[string]any) error {

	return err
}

func DeleteReviewByID(ctx context.Context, app *infra.Deps, reviewID string) error {
	return err
}

/* -------------------------
   Fetch Helpers
------------------------- */

// FetchReviewsByEntity queries reviews for a specific entity type and ID.
func FetchReviewsByEntity(ctx context.Context, app *infra.Deps, entityType, entityID string) ([]Review, error) {
	var reviews []Review
	return reviews, nil
}

// FetchReviewByID retrieves a single review by its ID.
func FetchReviewByID(ctx context.Context, app *infra.Deps, reviewID string) (*Review, error) {
	var review Review
	return &review, nil
}
