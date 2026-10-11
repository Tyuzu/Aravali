// File: internal/reviews/reviewsSQLDB.go

package reviews

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"scav/config"
	"scav/infra"
)

/* -------------------------
   Table
------------------------- */

var reviewsTable = config.Tables.ReviewsTable

func scanReviewRow(row map[string]any) Review {
	review := Review{}
	if reviewID, ok := row["reviewid"].(string); ok {
		review.ReviewID = reviewID
	}
	if userID, ok := row["userid"].(string); ok {
		review.UserID = userID
	}
	if entityType, ok := row["entity_type"].(string); ok {
		review.EntityType = entityType
	}
	if entityID, ok := row["entity_id"].(string); ok {
		review.EntityID = entityID
	}
	if rating, ok := row["rating"].(int64); ok {
		review.Rating = int(rating)
	}
	if rating, ok := row["rating"].(int); ok {
		review.Rating = rating
	}
	if comment, ok := row["comment"].(string); ok {
		review.Comment = comment
	}
	if createdAt, ok := row["created_at"].(time.Time); ok {
		review.CreatedAt = createdAt
	}
	if updatedAt, ok := row["updated_at"].(time.Time); ok {
		review.UpdatedAt = updatedAt
	}
	return review
}

/* -------------------------
   Database Helpers
------------------------- */

func GetUserReviewForEntity(ctx context.Context, app *infra.Deps, userID, entityType, entityID string) (*Review, error) {
	if app == nil || app.SQLDB == nil {
		return nil, nil
	}
	if userID == "" || entityType == "" || entityID == "" {
		return nil, errors.New("review lookup requires user, entity type, and entity id")
	}
	var reviewID, fetchedUserID, fetchedEntityType, fetchedEntityID, comment string
	var rating int
	var createdAt, updatedAt time.Time
	if err := app.SQLDB.QueryRow(ctx, `SELECT reviewid, userid, entity_type, entity_id, rating, comment, created_at, updated_at FROM reviews WHERE userid = $1 AND entity_type = $2 AND entity_id = $3 LIMIT 1`, userID, entityType, entityID).Scan(&reviewID, &fetchedUserID, &fetchedEntityType, &fetchedEntityID, &rating, &comment, &createdAt, &updatedAt); err != nil {
		return nil, err
	}
	review := &Review{
		ReviewID:   reviewID,
		UserID:     fetchedUserID,
		EntityType: fetchedEntityType,
		EntityID:   fetchedEntityID,
		Rating:     rating,
		Comment:    comment,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}
	return review, nil
}

func GetReviewByID(ctx context.Context, app *infra.Deps, reviewID string) (*Review, error) {
	if app == nil || app.SQLDB == nil {
		return nil, nil
	}
	if strings.TrimSpace(reviewID) == "" {
		return nil, errors.New("review id required")
	}
	var fetchedReviewID, userID, entityType, entityID, comment string
	var rating int
	var createdAt, updatedAt time.Time
	if err := app.SQLDB.QueryRow(ctx, `SELECT reviewid, userid, entity_type, entity_id, rating, comment, created_at, updated_at FROM reviews WHERE reviewid = $1 LIMIT 1`, reviewID).Scan(&fetchedReviewID, &userID, &entityType, &entityID, &rating, &comment, &createdAt, &updatedAt); err != nil {
		return nil, err
	}
	return &Review{ReviewID: fetchedReviewID, UserID: userID, EntityType: entityType, EntityID: entityID, Rating: rating, Comment: comment, CreatedAt: createdAt, UpdatedAt: updatedAt}, nil
}

func InsertReview(ctx context.Context, app *infra.Deps, review Review) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	if strings.TrimSpace(review.ReviewID) == "" {
		review.ReviewID = fmt.Sprintf("review_%d", time.Now().UnixNano())
	}
	if review.CreatedAt.IsZero() {
		review.CreatedAt = time.Now().UTC()
	}
	if review.UpdatedAt.IsZero() {
		review.UpdatedAt = review.CreatedAt
	}
	_, err := app.SQLDB.Exec(ctx,
		`INSERT INTO reviews (reviewid, userid, entity_type, entity_id, rating, comment, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (reviewid) DO UPDATE SET
			userid = EXCLUDED.userid,
			entity_type = EXCLUDED.entity_type,
			entity_id = EXCLUDED.entity_id,
			rating = EXCLUDED.rating,
			comment = EXCLUDED.comment,
			updated_at = NOW()`,
		review.ReviewID,
		review.UserID,
		review.EntityType,
		review.EntityID,
		review.Rating,
		review.Comment,
		review.CreatedAt,
		review.UpdatedAt,
	)
	return err
}

func UpdateReviewByID(ctx context.Context, app *infra.Deps, reviewID string, update map[string]any) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	if len(update) == 0 {
		return nil
	}
	updates := make([]string, 0, len(update))
	args := make([]any, 0, len(update)+1)
	for key, value := range update {
		switch key {
		case "updatedAt":
			updates = append(updates, "updated_at = $"+fmt.Sprintf("%d", len(args)+1))
			args = append(args, value)
		case "rating":
			updates = append(updates, "rating = $"+fmt.Sprintf("%d", len(args)+1))
			args = append(args, value)
		case "comment":
			updates = append(updates, "comment = $"+fmt.Sprintf("%d", len(args)+1))
			args = append(args, value)
		case "likes":
			updates = append(updates, "metadata = jsonb_set(COALESCE(metadata, '{}'::jsonb), '{likes}', to_jsonb($"+fmt.Sprintf("%d", len(args)+1)+")::jsonb, true)")
			args = append(args, value)
		case "dislikes":
			updates = append(updates, "metadata = jsonb_set(COALESCE(metadata, '{}'::jsonb), '{dislikes}', to_jsonb($"+fmt.Sprintf("%d", len(args)+1)+")::jsonb, true)")
			args = append(args, value)
		default:
			updates = append(updates, key+" = $"+fmt.Sprintf("%d", len(args)+1))
			args = append(args, value)
		}
	}
	if len(updates) == 0 {
		return nil
	}
	updates = append(updates, "updated_at = NOW()")
	query := `UPDATE reviews SET ` + strings.Join(updates, ", ") + ` WHERE reviewid = $` + fmt.Sprintf("%d", len(args)+1)
	args = append(args, reviewID)
	_, err := app.SQLDB.Exec(ctx, query, args...)
	return err
}

func DeleteReviewByID(ctx context.Context, app *infra.Deps, reviewID string) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	_, err := app.SQLDB.Exec(ctx, `DELETE FROM reviews WHERE reviewid = $1`, reviewID)
	return err
}

/* -------------------------
   Fetch Helpers
------------------------- */

// FetchReviewsByEntity queries reviews for a specific entity type and ID.
func FetchReviewsByEntity(ctx context.Context, app *infra.Deps, entityType, entityID string) ([]Review, error) {
	if app == nil || app.SQLDB == nil {
		return nil, nil
	}
	rows, err := app.SQLDB.Query(ctx, `SELECT reviewid, userid, entity_type, entity_id, rating, comment, created_at, updated_at FROM reviews WHERE entity_type = $1 AND entity_id = $2 ORDER BY created_at DESC`, entityType, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	results := make([]Review, 0)
	for rows.Next() {
		var reviewID, userID, fetchedEntityType, fetchedEntityID, comment string
		var rating int
		var createdAt, updatedAt time.Time
		if err := rows.Scan(&reviewID, &userID, &fetchedEntityType, &fetchedEntityID, &rating, &comment, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		results = append(results, Review{ReviewID: reviewID, UserID: userID, EntityType: fetchedEntityType, EntityID: fetchedEntityID, Rating: rating, Comment: comment, CreatedAt: createdAt, UpdatedAt: updatedAt})
	}
	return results, rows.Err()
}

// FetchReviewByID retrieves a single review by its ID.
func FetchReviewByID(ctx context.Context, app *infra.Deps, reviewID string) (*Review, error) {
	return GetReviewByID(ctx, app, reviewID)
}
