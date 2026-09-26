package media

import (
	"context"

	"scav/config"
	"scav/infra"
)

var mediaTable = config.Tables.MediaTable

func SQLinsertMedia(ctx context.Context, app *infra.Deps, media Media) error {
	return app.SQLDB.Insert(ctx, mediaTable, media)
}

func SQLgetMediaByID(ctx context.Context, app *infra.Deps, entityType, entityID, mediaID string) (Media, error) {
	var media Media
	query := "entityid = $1 AND entitytype = $2 AND mediaid = $3"
	args := []any{entityID, entityType, mediaID}

	err := app.SQLDB.FindOne(ctx, mediaTable, query, args, &media)
	return media, err
}

func SQLlistMediaByEntity(ctx context.Context, app *infra.Deps, entityType, entityID string) ([]Media, error) {
	query := "entityid = $1 AND entitytype = $2"
	args := []any{entityID, entityType}

	var medias []Media
	err := app.SQLDB.FindMany(ctx, mediaTable, query, args, &medias)
	return medias, err
}

func SQLgetMediaGroupsByEntity(ctx context.Context, app *infra.Deps, entityType, entityID string) ([]map[string]any, error) {
	medias, err := SQLlistMediaByEntity(ctx, app, entityType, entityID)
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

func SQLupdateMediaGroup(ctx context.Context, app *infra.Deps, mediaGroupID string, updateFields map[string]any) ([]Media, error) {
	query := "mediagroupid = $1"
	args := []any{mediaGroupID}

	if _, err := app.SQLDB.UpdateMany(ctx, mediaTable, query, args, updateFields); err != nil {
		return nil, err
	}

	var updatedMedias []Media
	err := app.SQLDB.FindMany(ctx, mediaTable, query, args, &updatedMedias)
	return updatedMedias, err
}
