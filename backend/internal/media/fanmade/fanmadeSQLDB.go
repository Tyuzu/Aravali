// File: internal/media/fanmade/fanmadeSQLDB.go

package fanmade

import (
	"context"
	"fmt"
	"strings"
	"time"

	"scav/config"
	"scav/infra"
	"scav/internal/media"
)

var fanmadeMediaTable = config.Tables.MediaTable

func insertFanMedia(ctx context.Context, app *infra.Deps, media media.Media) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	if media.CreatedAt.IsZero() {
		media.CreatedAt = time.Now().UTC()
	}
	_, err := app.SQLDB.Exec(ctx,
		fmt.Sprintf("INSERT INTO %s (mediaid, mediagroupid, type, url, thumbnailurl, caption, description, creatorid, likescount, commentscount, visibility, tags, duration, filesize, mimetype, isfeatured, entityid, entitytype, createdat, updatedat, userid, extn, captionlang) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23)", fanmadeMediaTable),
		media.MediaID, media.MediaGroupID, media.Type, media.URL, media.ThumbnailURL, media.Caption, media.Description, media.CreatorID, media.LikesCount, media.CommentsCount, media.Visibility, media.Tags, media.Duration, media.FileSize, media.MimeType, media.IsFeatured, media.EntityID, media.EntityType, media.CreatedAt, media.UpdatedAt, media.UserID, media.Extn, media.CaptionLang,
	)
	return err
}

func getFanMediaByID(ctx context.Context, app *infra.Deps, entityType, entityID, mediaID string) (media.Media, error) {
	if app == nil || app.SQLDB == nil {
		return media.Media{}, nil
	}
	if entityType == "" || entityID == "" || mediaID == "" {
		return media.Media{}, nil
	}
	row := app.SQLDB.QueryRow(ctx, fmt.Sprintf("SELECT * FROM %s WHERE entitytype = $1 AND entityid = $2 AND mediaid = $3 LIMIT 1", fanmadeMediaTable), entityType, entityID, mediaID)
	if row == nil {
		return media.Media{}, nil
	}
	return media.Media{}, nil
}

func listFanMediasByEntity(ctx context.Context, app *infra.Deps, entityType, entityID string) ([]media.Media, error) {
	if app == nil || app.SQLDB == nil {
		return nil, nil
	}
	if entityType == "" || entityID == "" {
		return []media.Media{}, nil
	}
	rows, err := app.SQLDB.Query(ctx, fmt.Sprintf("SELECT * FROM %s WHERE entitytype = $1 AND entityid = $2 ORDER BY createdat DESC", fanmadeMediaTable), entityType, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return []media.Media{}, nil
}

func listFanMediaGroupsByEntity(ctx context.Context, app *infra.Deps, entityType, entityID string) ([]map[string]any, error) {
	medias, err := listFanMediasByEntity(ctx, app, entityType, entityID)
	if err != nil {
		return nil, err
	}

	mediaMap := make(map[string][]media.Media)
	for _, item := range medias {
		mediaMap[item.MediaGroupID] = append(mediaMap[item.MediaGroupID], item)
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
	if app == nil || app.SQLDB == nil {
		return nil, nil
	}
	if mediaGroupID == "" || len(update) == 0 {
		return []media.Media{}, nil
	}
	setParts := make([]string, 0, len(update))
	args := make([]any, 0, len(update))
	idx := 1
	for key, value := range update {
		setParts = append(setParts, fmt.Sprintf("%s = $%d", strings.ToLower(key), idx))
		args = append(args, value)
		idx++
	}
	_, err := app.SQLDB.Exec(ctx, fmt.Sprintf("UPDATE %s SET %s WHERE mediagroupid = $%d", fanmadeMediaTable, strings.Join(setParts, ", "), len(args)+1), append(args, mediaGroupID)...)
	if err != nil {
		return nil, err
	}
	return []media.Media{}, nil
}

func deleteFanMediaByID(ctx context.Context, app *infra.Deps, mediaID string) (int64, error) {
	if app == nil || app.SQLDB == nil {
		return 0, nil
	}
	if mediaID == "" {
		return 0, nil
	}
	res, err := app.SQLDB.Exec(ctx, fmt.Sprintf("DELETE FROM %s WHERE mediaid = $1", fanmadeMediaTable), mediaID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}
