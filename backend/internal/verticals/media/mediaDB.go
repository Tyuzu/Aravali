package media

import (
	"context"

	"scav/config"
	"scav/infra"
)

var mediaCollection = config.Collections.MediaCollection

func insertMedia(ctx context.Context, app *infra.Deps, media Media) error {
	return app.DB.Insert(ctx, mediaCollection, media)
}

func getMediaByID(ctx context.Context, app *infra.Deps, entityType, entityID, mediaID string) (Media, error) {
	var media Media
	err := app.DB.FindOne(ctx, mediaCollection, map[string]any{
		"entityid":   entityID,
		"entitytype": entityType,
		"mediaid":    mediaID,
	}, &media)
	return media, err
}

func listMediaByEntity(ctx context.Context, app *infra.Deps, entityType, entityID string) ([]Media, error) {
	filter := map[string]any{
		"entityid":   entityID,
		"entitytype": entityType,
	}

	var medias []Media
	err := app.DB.FindMany(ctx, mediaCollection, filter, &medias)
	return medias, err
}

func getMediaGroupsByEntity(ctx context.Context, app *infra.Deps, entityType, entityID string) ([]map[string]any, error) {
	medias, err := listMediaByEntity(ctx, app, entityType, entityID)
	if err != nil {
		return nil, err
	}

	mediaMap := make(map[string][]Media)
	for _, media := range medias {
		mediaMap[media.MediaGroupID] = append(mediaMap[media.MediaGroupID], media)
	}

	groups := make([]map[string]any, 0, len(mediaMap))
	for groupID, files := range mediaMap {
		groups = append(groups, map[string]any{
			"groupId": groupID,
			"files":   files,
		})
	}

	return groups, nil
}

func updateMediaGroup(ctx context.Context, app *infra.Deps, mediaGroupID string, updateFields map[string]any) ([]Media, error) {
	if _, err := app.DB.UpdateMany(ctx, mediaCollection, map[string]any{"mediaGroupId": mediaGroupID}, map[string]any{"$set": updateFields}); err != nil {
		return nil, err
	}

	var updatedMedias []Media
	err := app.DB.FindMany(ctx, mediaCollection, map[string]any{"mediaGroupId": mediaGroupID}, &updatedMedias)
	return updatedMedias, err
}
