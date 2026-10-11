// File: internal/verticals/media/fanmade/fanmadeSQLDB.go

package fanmade

import (
	"context"

	"scav/config"
	"scav/infra"
	"scav/internal/verticals/media"
)

var fanmadeMediaTable = config.Tables.MediaTable

func insertFanMedia(ctx context.Context, app *infra.Deps, media media.Media) error {
}

func getFanMediaByID(ctx context.Context, app *infra.Deps, entityType, entityID, mediaID string) (media.Media, error) {
	var media media.Media
	return media, err
}

func listFanMediasByEntity(ctx context.Context, app *infra.Deps, entityType, entityID string) ([]media.Media, error) {
	var medias []media.Media
	return medias, err
}

func listFanMediaGroupsByEntity(ctx context.Context, app *infra.Deps, entityType, entityID string) ([]map[string]any, error) {
	medias, err := listFanMediasByEntity(ctx, app, entityType, entityID)
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

func updateFanMediaGroup(ctx context.Context, app *infra.Deps, mediaGroupID string, update map[string]any) ([]media.Media, error) {
	var updatedMedias []media.Media
	return updatedMedias, err
}

func deleteFanMediaByID(ctx context.Context, app *infra.Deps, mediaID string) (int64, error) {
}
