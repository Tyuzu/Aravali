// File: internal/beats/ads/adsSQLDB.go

package ads

import (
	"context"
	"fmt"
	"time"

	"scav/config"
	"scav/infra"
)

var (
	adsTable   = config.Tables.AdsTable
	postsTable = config.Tables.FarmsTable
)

func FetchActiveAdsFromDB(ctx context.Context, app *infra.Deps) ([]Ad, error) {
	var dbAds []Ad

	return dbAds, nil
}

func ListAdsFromDB(ctx context.Context, app *infra.Deps) ([]Ad, error) {
	var ads []Ad

	return ads, nil
}

func CreateAdInDB(ctx context.Context, app *infra.Deps, ad *Ad) error {
	return nil
}

// PromotePost creates an Ad entry sourced directly from an existing post.
func PromotePostInDB(ctx context.Context, app *infra.Deps, postID, page, position, category string) (*Ad, error) {
	// Fetch target post to derive ad details
	var post struct {
		ID       string `db:"id"`
		Title    string `db:"title"`
		Summary  string `db:"summary"`
		CoverImg string `db:"cover_image"`
		Category string `db:"category"`
	}

	if category == "" {
		category = post.Category
	}

	ad := &Ad{
		Type:        TypePost,
		PostID:      postID,
		Title:       post.Title,
		Description: post.Summary,
		Image:       post.CoverImg,
		Link:        fmt.Sprintf("/posts/%s", postID), // Internal client routing path
		Category:    category,
		Page:        page,
		Position:    position,
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return ad, nil
}

func GetAdByIDFromDB(ctx context.Context, app *infra.Deps, id string) (*Ad, error) {
	var ad Ad

	return &ad, nil
}

func UpdateAdInDB(ctx context.Context, app *infra.Deps, id string, updateData map[string]any) error {
	return nil
}

func DeleteAdInDB(ctx context.Context, app *infra.Deps, id string) error {
	return nil
}
