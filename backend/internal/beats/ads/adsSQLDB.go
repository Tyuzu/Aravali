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

func SQLFetchActiveAdsFromDB(ctx context.Context, app *infra.Deps) ([]Ad, error) {
	var dbAds []Ad
	where := "status = $1"
	args := []any{"active"}

	err := app.SQLDB.FindMany(ctx, adsTable, where, args, &dbAds)
	if err != nil {
		return nil, err
	}

	return dbAds, nil
}

func SQLListAdsFromDB(ctx context.Context, app *infra.Deps) ([]Ad, error) {
	var ads []Ad
	where := "1=1"
	args := []any{}

	if err := app.SQLDB.FindMany(ctx, adsTable, where, args, &ads); err != nil {
		return nil, err
	}
	if ads == nil {
		ads = []Ad{}
	}
	return ads, nil
}

func SQLCreateAdInDB(ctx context.Context, app *infra.Deps, ad *Ad) error {
	ad.CreatedAt = time.Now()
	ad.UpdatedAt = time.Now()

	if ad.Status == "" {
		ad.Status = "active"
	}
	if ad.Type == "" {
		ad.Type = TypeExternal
	}

	return app.SQLDB.InsertOne(ctx, adsTable, ad)
}

// PromotePost creates an Ad entry sourced directly from an existing post.
func SQLPromotePostInDB(ctx context.Context, app *infra.Deps, postID, page, position, category string) (*Ad, error) {
	// Fetch target post to derive ad details
	var post struct {
		ID       string `db:"id"`
		Title    string `db:"title"`
		Summary  string `db:"summary"`
		CoverImg string `db:"cover_image"`
		Category string `db:"category"`
	}

	where := "id = $1"
	args := []any{postID}

	err := app.SQLDB.FindOne(ctx, postsTable, where, args, &post)
	if err != nil {
		return nil, fmt.Errorf("post not found: %w", err)
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

	err = app.SQLDB.InsertOne(ctx, adsTable, ad)
	if err != nil {
		return nil, err
	}

	return ad, nil
}

func SQLGetAdByIDFromDB(ctx context.Context, app *infra.Deps, id string) (*Ad, error) {
	var ad Ad
	where := "id = $1"
	args := []any{id}

	err := app.SQLDB.FindOne(ctx, adsTable, where, args, &ad)
	if err != nil {
		return nil, err
	}

	return &ad, nil
}

func SQLUpdateAdInDB(ctx context.Context, app *infra.Deps, id string, updateData map[string]any) error {
	where := "id = $1"
	args := []any{id}

	updateData["updated_at"] = time.Now()

	_, err := app.SQLDB.UpdateOne(ctx, adsTable, where, args, updateData)
	return err
}

func SQLDeleteAdInDB(ctx context.Context, app *infra.Deps, id string) error {
	where := "id = $1"
	args := []any{id}

	_, err := app.SQLDB.DeleteOne(ctx, adsTable, where, args)
	return err
}
