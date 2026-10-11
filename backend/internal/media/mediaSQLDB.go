// File: internal/media/mediaSQLDB.go

package media

import (
	"context"

	"scav/config"
	"scav/infra"
)

var mediaTable = config.Tables.MediaTable

func insertMedia(ctx context.Context, app *infra.Deps, media Media) error {
	return nil
}

func getMediaByID(ctx context.Context, app *infra.Deps, entityType, entityID, mediaID string) (Media, error) {
	return Media{}, nil
}

func listMediaByEntity(ctx context.Context, app *infra.Deps, entityType, entityID string) ([]Media, error) {
	return nil, nil
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
	return nil, nil
}
