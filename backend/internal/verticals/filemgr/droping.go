// File: internal/verticals/filemgr/droping.go

package filemgr

import (
	"fmt"
	"net/http"
	log "scav/utils/logger"
	"strings"

	"scav/config"
	"scav/infra"
)

const maxUploadBytes = 200 << 20 // 200 MB

type Attachment struct {
	Filename    string `json:"filename"`
	Extension   string `json:"extension"`
	Key         string `json:"key"`
	Resolutions []int  `json:"resolutions,omitempty"`
}

var validEntities = map[string]EntityType{
	"artist":       EntityArtist,
	"baito":        EntityBaito,
	"baito_worker": EntityWorker,
	"blogpost":     EntityBlogPost,
	"chat":         EntityChat,
	"crop":         EntityCrop,
	"event":        EntityEvent,
	"farm":         EntityFarm,
	"feedpost":     EntityFeed,
	"live":         EntityLive,
	"media":        EntityMedia,
	"menu":         EntityMenu,
	"merch":        EntityMerch,
	"music":        EntityMusic,
	"place":        EntityPlace,
	"product":      EntityProduct,
	"recipe":       EntityRecipe,
	"report":       EntityReport,
	"review":       EntityReview,
	"song":         EntitySong,
	"tool":         EntityProduct,
	"user":         EntityUser,
	"vendor":       EntityVendor,
	"worker":       EntityWorker,
}

type EntityMeta struct {
	Table   string
	IDField string
}

var entityMeta = map[string]EntityMeta{
	"artist":   {Table: config.Tables.ArtistsTable, IDField: config.IDField.ArtistId},
	"baito":    {Table: config.Tables.BaitoTable, IDField: config.IDField.BaitoId},
	"blogpost": {Table: config.Tables.BlogPostsTable, IDField: config.IDField.BlogPostId},
	"chat":     {Table: config.Tables.ChatsTable, IDField: config.IDField.ChatId},
	"crop":     {Table: config.Tables.CropsTable, IDField: config.IDField.CropId},
	"event":    {Table: config.Tables.EventsTable, IDField: config.IDField.EventId},
	"farm":     {Table: config.Tables.FarmsTable, IDField: config.IDField.FarmId},
	"feedpost": {Table: config.Tables.FeedPostsTable, IDField: config.IDField.FeedPostId},
	"live":     {Table: "vlive", IDField: "eventid"},
	"media":    {Table: config.Tables.MediaTable, IDField: config.IDField.MediaId},
	"menu":     {Table: config.Tables.MenuTable, IDField: config.IDField.MenuId},
	"merch":    {Table: config.Tables.MerchTable, IDField: config.IDField.MerchId},
	"music":    {Table: config.Tables.AlbumsTable, IDField: config.IDField.AlbumId},
	"place":    {Table: config.Tables.PlacesTable, IDField: config.IDField.PlaceId},
	"product":  {Table: config.Tables.ProductTable, IDField: config.IDField.ProductId},
	"recipe":   {Table: config.Tables.RecipeTable, IDField: config.IDField.RecipeId},
	"report":   {Table: config.Tables.ReportsTable, IDField: config.IDField.ReportId},
	"review":   {Table: config.Tables.ReviewsTable, IDField: config.IDField.ReviewId},
	"song":     {Table: config.Tables.SongsTable, IDField: config.IDField.SongId},
	"user":     {Table: config.Tables.UserTable, IDField: config.IDField.UserId},
	"vendor":   {Table: config.Tables.VendorTable, IDField: config.IDField.VendorId},
	"worker":   {Table: config.Tables.BaitoWorkerTable, IDField: config.IDField.BaitoWorkerId},
}

func validateUploadRequest(w http.ResponseWriter, r *http.Request) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadBytes)
	if r.Method != http.MethodPost {
		return fmt.Errorf("method must be POST")
	}
	contentType := r.Header.Get("Content-Type")
	remoteURL := strings.TrimSpace(r.FormValue("remoteUrl"))
	if remoteURL == "" && !strings.HasPrefix(contentType, "multipart/") {
		return fmt.Errorf("content-type must be multipart/form-data")
	}
	return nil
}

func convertToAttachments(serviceAttachments []Attachment) []Attachment {
	return append([]Attachment(nil), serviceAttachments...)
}

func updateEntityMedia(app *infra.Deps, entityType string, entityId string, attachments []Attachment) (any, error) {
	log.Println("updateEntityMedia:", entityType, entityId) // #nosec G706
	log.Println("updateEntityMedia:", attachments)          // #nosec G706
	meta, ok := entityMeta[entityType]
	if !ok {
		return nil, fmt.Errorf("unsupported entity type: %s", entityType)
	}

	setFields := map[string]any{}
	var images []string

	for _, attachment := range attachments {
		key := PictureType(strings.ToLower(strings.TrimSpace(attachment.Key)))
		switch key {
		case PicBanner:
			setFields["banner"] = attachment.Filename
		case PicMember:
			setFields["member"] = attachment.Filename
		case PicPoster:
			setFields["poster"] = attachment.Filename
		case PicThumb:
			setFields["thumb"] = attachment.Filename
		case PicSeating:
			setFields["seating"] = attachment.Filename
		case PicPhoto:
			setFields["photo"] = attachment.Filename
		case PicImage:
			images = append(images, attachment.Filename)
		case PicVideo:
			setFields["video"] = attachment.Filename
		case PicAudio:
			setFields["audio"] = attachment.Filename
		case PicSong:
			setFields["song"] = attachment.Filename
		case PicDocument:
			setFields["document"] = attachment.Filename
		case PicFile:
			setFields["file"] = attachment.Filename
		}
	}

	update := map[string]any{}
	if len(setFields) > 0 {
		update["$set"] = setFields
	}
	if len(images) > 0 {
		update["$push"] = map[string]any{"images": map[string]any{"$each": images}}
	}
	if len(update) == 0 {
		return nil, nil
	}

	return updateEntityMediaInDB(app, meta.Table, meta.IDField, entityId, update)
}
