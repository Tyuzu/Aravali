package ads

import (
	"context"
	"fmt"
	"time"

	"scav/config"
	"scav/infra"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	adsTable   = config.Tables.AdsTable
	postsTable = config.Tables.FarmsTable
)

func SQLFetchActiveAdsFromDB(ctx context.Context, app *infra.Deps) ([]Ad, error) {
	if app == nil || app.DB == nil {
		return nil, nil
	}

	var dbAds []Ad
	filter := map[string]any{"status": "active"}

	err := app.DB.FindMany(ctx, adsTable, filter, &dbAds)
	if err != nil {
		return nil, err
	}

	return dbAds, nil
}

func SQLListAdsFromDB(ctx context.Context, app *infra.Deps) ([]Ad, error) {
	var ads []Ad
	if err := app.DB.FindMany(ctx, adsTable, map[string]any{}, &ads); err != nil {
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

	return app.DB.InsertOne(ctx, adsTable, ad)
}

// PromotePost creates an Ad entry sourced directly from an existing post.
func SQLPromotePostInDB(ctx context.Context, app *infra.Deps, postID, page, position, category string) (*Ad, error) {
	objID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return nil, fmt.Errorf("invalid post ID format")
	}

	// Fetch target post to derive ad details
	var post struct {
		ID       string `bson:"_id"`
		Title    string `bson:"title"`
		Summary  string `bson:"summary"`
		CoverImg string `bson:"coverImage"`
		Category string `bson:"category"`
	}

	err = app.DB.FindOne(ctx, postsTable, map[string]any{"_id": objID}, &post)
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

	err = app.DB.InsertOne(ctx, adsTable, ad)
	if err != nil {
		return nil, err
	}

	return ad, nil
}

func SQLGetAdByIDFromDB(ctx context.Context, app *infra.Deps, id string) (*Ad, error) {
	var ad Ad
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = app.DB.FindOne(ctx, adsTable, map[string]any{"_id": objID}, &ad)
	if err != nil {
		return nil, err
	}

	return &ad, nil
}

func SQLUpdateAdInDB(ctx context.Context, app *infra.Deps, id string, updateData map[string]interface{}) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	updateData["updatedAt"] = time.Now()
	update := map[string]any{"$set": updateData}
	_, err = app.DB.UpdateOne(ctx, adsTable, map[string]any{"_id": objID}, update)
	return err
}

func SQLDeleteAdInDB(ctx context.Context, app *infra.Deps, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = app.DB.DeleteOne(ctx, adsTable, map[string]any{"_id": objID})

	return err
}
