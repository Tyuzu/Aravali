// File: internal/media/delete.go

package media

import (
	"net/http"
	"scav/config"
	"scav/config/mqevent"
	"scav/infra"
	log "scav/infra/logger"
	"scav/infra/mq"
	"scav/internal/beats/userdata"
	"scav/utils"
)

// DeleteMedia deletes a single media item if the requesting user is the creator
func DeleteMedia(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		entityType := utils.GetParam(r, "entitytype")
		entityID := utils.GetParam(r, "entityid")
		mediaID := utils.GetParam(r, "id")

		requestingUserID, ok := ctx.Value(config.UserIDKey).(string)
		if !ok || requestingUserID == "" {
			http.Error(w, "Invalid user", http.StatusUnauthorized)
			return
		}

		mediaItem, err := getMediaByID(ctx, app, entityType, entityID, mediaID)
		if err != nil {
			http.Error(w, "Media not found", http.StatusNotFound)
			return
		}

		if mediaItem.CreatorID != requestingUserID {
			http.Error(w, "Not authorized to delete this media", http.StatusForbidden)
			return
		}

		if _, err := deleteMediaByID(ctx, app, mediaID); err != nil {
			http.Error(w, "Failed to delete media", http.StatusInternalServerError)
			return
		}

		userdata.DelUserData("media", mediaID, requestingUserID, app)

		if err := mq.PublishWithMeta(ctx, app.MQ, mqevent.MediaDeletedEvent, mqevent.MediaDeletedPayload{}); err != nil {
			log.Printf("failed to publish media deleted event: %v", err)
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"message": "Media deleted successfully",
		})
	}
}
