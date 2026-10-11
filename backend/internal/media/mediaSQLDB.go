// File: internal/media/mediaSQLDB.go

package media

import (
	"context"
	"fmt"
	"strings"
	"time"

	"scav/config"
	"scav/infra"
)

var mediaTable = config.Tables.MediaTable

func insertMedia(ctx context.Context, app *infra.Deps, m Media) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now().UTC()
	}
	_, err := app.SQLDB.Exec(ctx,
		fmt.Sprintf("INSERT INTO %s (mediaid, mediagroupid, type, url, thumbnailurl, caption, description, creatorid, likescount, commentscount, visibility, tags, duration, filesize, mimetype, isfeatured, entityid, entitytype, createdat, updatedat, userid, extn, captionlang) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23)", mediaTable),
		m.MediaID, m.MediaGroupID, m.Type, m.URL, m.ThumbnailURL, m.Caption, m.Description, m.CreatorID, m.LikesCount, m.CommentsCount, m.Visibility, m.Tags, m.Duration, m.FileSize, m.MimeType, m.IsFeatured, m.EntityID, m.EntityType, m.CreatedAt, m.UpdatedAt, m.UserID, m.Extn, m.CaptionLang,
	)
	return err
}

func getMediaByID(ctx context.Context, app *infra.Deps, entityType, entityID, mediaID string) (Media, error) {
	if app == nil || app.SQLDB == nil {
		return Media{}, nil
	}
	if entityType == "" || entityID == "" || mediaID == "" {
		return Media{}, nil
	}
	row := app.SQLDB.QueryRow(ctx, fmt.Sprintf("SELECT * FROM %s WHERE entitytype = $1 AND entityid = $2 AND mediaid = $3 LIMIT 1", mediaTable), entityType, entityID, mediaID)
	if row == nil {
		return Media{}, nil
	}
	return Media{}, nil
}

func listMediaByEntity(ctx context.Context, app *infra.Deps, entityType, entityID string) ([]Media, error) {
	if app == nil || app.SQLDB == nil {
		return nil, nil
	}
	if entityType == "" || entityID == "" {
		return []Media{}, nil
	}
	rows, err := app.SQLDB.Query(ctx, fmt.Sprintf("SELECT * FROM %s WHERE entitytype = $1 AND entityid = $2 ORDER BY createdat DESC", mediaTable), entityType, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return []Media{}, nil
}

func getMediaGroupsByEntity(ctx context.Context, app *infra.Deps, entityType, entityID string) ([]map[string]any, error) {
	medias, err := listMediaByEntity(ctx, app, entityType, entityID)
	if err != nil {
		return nil, err
	}

	mediaMap := make(map[string][]Media)
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

func updateMediaGroup(ctx context.Context, app *infra.Deps, mediaGroupID string, update map[string]any) ([]Media, error) {
	if app == nil || app.SQLDB == nil {
		return nil, nil
	}
	if mediaGroupID == "" || len(update) == 0 {
		return []Media{}, nil
	}
	setParts := make([]string, 0, len(update))
	args := make([]any, 0, len(update))
	idx := 1
	for key, value := range update {
		setParts = append(setParts, fmt.Sprintf("%s = $%d", strings.ToLower(key), idx))
		args = append(args, value)
		idx++
	}
	_, err := app.SQLDB.Exec(ctx, fmt.Sprintf("UPDATE %s SET %s WHERE mediagroupid = $%d", mediaTable, strings.Join(setParts, ", "), len(args)+1), append(args, mediaGroupID)...)
	if err != nil {
		return nil, err
	}
	return []Media{}, nil
}

func deleteMediaByID(ctx context.Context, app *infra.Deps, mediaID string) (int64, error) {
	if app == nil || app.SQLDB == nil {
		return 0, nil
	}
	if mediaID == "" {
		return 0, nil
	}
	res, err := app.SQLDB.Exec(ctx, fmt.Sprintf("DELETE FROM %s WHERE mediaid = $1", mediaTable), mediaID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}
