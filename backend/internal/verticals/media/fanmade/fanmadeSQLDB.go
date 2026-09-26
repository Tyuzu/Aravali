package fanmade

import (
	"context"

	"scav/config"
	"scav/infra"
	"scav/infra/sqldb"
	"scav/internal/verticals/media"
)

var fanmadeMediaCollection = config.Collections.MediaCollection

func SQLinsertFanMedia(ctx context.Context, app *infra.Deps, media media.Media) error {
	return app.SQLDB.Insert(ctx, fanmadeMediaCollection, media)
}

func SQLgetFanMediaByID(ctx context.Context, app *infra.Deps, entityType, entityID, mediaID string) (media.Media, error) {
	var media media.Media
	query := "entityid = $1 AND entitytype = $2 AND mediaid = $3"
	args := []any{entityID, entityType, mediaID}

	err := app.SQLDB.FindOne(ctx, fanmadeMediaCollection, query, args, &media)
	return media, err
}

func SQLlistFanMediasByEntity(ctx context.Context, app *infra.Deps, entityType, entityID string) ([]media.Media, error) {
	var medias []media.Media
	query := "entityid = $1 AND entitytype = $2"
	args := []any{entityID, entityType}

	opts := sqldb.FindManyOptions{}
	err := app.SQLDB.FindManyWithOptions(ctx, fanmadeMediaCollection, query, args, opts, &medias)
	return medias, err
}

func SQLlistFanMediaGroupsByEntity(ctx context.Context, app *infra.Deps, entityType, entityID string) ([]map[string]any, error) {
	medias, err := SQLlistFanMediasByEntity(ctx, app, entityType, entityID)
	if err != nil {
		return nil, err
	}

	mediaMap := make(map[string][]media.Media)
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

func SQLupdateFanMediaGroup(ctx context.Context, app *infra.Deps, mediaGroupID string, update map[string]any) ([]media.Media, error) {
	query := "mediagroupid = $1"
	args := []any{mediaGroupID}

	if _, err := app.SQLDB.UpdateMany(ctx, fanmadeMediaCollection, query, args, update); err != nil {
		return nil, err
	}

	var updatedMedias []media.Media
	opts := sqldb.FindManyOptions{}
	err := app.SQLDB.FindManyWithOptions(ctx, fanmadeMediaCollection, query, args, opts, &updatedMedias)
	return updatedMedias, err
}

func SQLdeleteFanMediaByID(ctx context.Context, app *infra.Deps, mediaID string) (int64, error) {
	query := "mediaid = $1"
	args := []any{mediaID}

	return app.SQLDB.DeleteOne(ctx, fanmadeMediaCollection, query, args)
}
