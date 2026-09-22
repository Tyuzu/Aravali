package filemgr

import (
	"context"
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"

	"scav/config/mqevent"
	"scav/infra"
	"scav/infra/mq"
	mediaworker "scav/infra/workers"
	"scav/utils"
	log "scav/utils/logger"
)

func FiledropHandler(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if err := validateUploadRequest(w, r); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Enforce maximum upload body size
		r.Body = http.MaxBytesReader(w, r.Body, maxUploadBytes)

		if err := r.ParseMultipartForm(maxUploadBytes); err != nil { // #nosec G120
			utils.RespondWithError(w, http.StatusBadRequest, "failed to parse multipart form: "+err.Error())
			return
		}

		// Clean up temporary files when handler finishes
		if r.MultipartForm != nil {
			defer func() {
				if err := r.MultipartForm.RemoveAll(); err != nil {
					log.Printf("[Filedrop] failed to remove multipart temp files: %v", err)
				}
			}()
		}

		entityType := strings.ToLower(strings.TrimSpace(r.FormValue("entityType")))
		entityId := strings.TrimSpace(r.FormValue("entityId"))
		remoteURL := strings.TrimSpace(r.FormValue("remoteUrl"))
		remoteKey := strings.TrimSpace(r.FormValue("remoteKey"))

		if entityType == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "entityType is required")
			return
		}

		if _, ok := validEntities[entityType]; !ok {
			utils.RespondWithError(w, http.StatusBadRequest, "invalid entityType")
			return
		}

		log.Printf("[Filedrop] entityType=%s entityId=%s", entityType, entityId) // #nosec G706

		fileService := NewFileService()
		userid := utils.GetUserIDFromRequest(r)

		var (
			attachments []Attachment
			err         error
		)

		if remoteURL != "" {
			remoteKey = string(normalizePictureKey(remoteKey))
			if remoteKey == "" {
				utils.RespondWithError(w, http.StatusBadRequest, "remoteKey is required")
				return
			}
			if _, ok := AllowedExtensions[PictureType(remoteKey)]; !ok {
				utils.RespondWithError(w, http.StatusBadRequest, "invalid remoteKey")
				return
			}
			attachments, err = fileService.ProcessRemoteFile(remoteURL, remoteKey, entityType, entityId, userid)
		} else {
			if r.MultipartForm == nil || len(r.MultipartForm.File) == 0 {
				utils.RespondWithError(w, http.StatusBadRequest, "no files uploaded")
				return
			}
			attachments, err = fileService.ProcessUploadedFiles(app, r, entityType, entityId, userid)
		}

		if err != nil {
			log.Printf("[Filedrop] processing error: %v", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "failed to process files: "+err.Error())
			return
		}

		if entityId != "" {
			if _, err := updateEntityMedia(app, entityType, entityId, attachments); err != nil {
				log.Printf("[Filedrop] failed updating entity media: %v", err)
				utils.RespondWithError(w, http.StatusInternalServerError, "failed to update entity media: "+err.Error())
				return
			}
		}

		// Publish FileCreated event
		payload := mqevent.FileCreatedPayload{
			UserID:     userid,
			EntityType: entityType,
			EntityID:   entityId,
			Count:      len(attachments),
		}

		if mqpayload, err := json.Marshal(payload); err == nil {
			if err := mq.PublishWithMeta(ctx, app.MQ, mqevent.FileCreatedEvent, mqpayload); err != nil {
				log.Printf("[Filedrop] failed to publish FileCreatedEvent: %v", err)
			}
		}

		// Enqueue media jobs in a single loop
		for _, att := range attachments {
			picType := PictureType(att.Key)
			savedPath := filepath.Join(ResolvePath(EntityType(entityType), picType), att.Filename+att.Extension)

			if isImageType(picType) {
				// Metadata Extraction Job
				metaJob := mediaworker.MediaJob{
					JobID:      generateUniqueID(),
					Type:       "image_metadata",
					SavedPath:  savedPath,
					UploadDir:  ResolvePath(EntityType(entityType), picType),
					UniqueID:   att.Filename,
					Filename:   att.Filename,
					Ext:        att.Extension,
					ThumbWidth: defaultThumbWidth,
					UserID:     userid,
				}
				publishJob(ctx, app, metaJob)

				// Thumbnail Job
				thumbJob := mediaworker.MediaJob{
					JobID:      generateUniqueID(),
					Type:       "image",
					SavedPath:  savedPath,
					UploadDir:  ResolvePath(EntityType(entityType), PicThumb),
					UniqueID:   att.Filename,
					Filename:   att.Filename,
					Ext:        att.Extension,
					ThumbWidth: defaultThumbWidth,
					UserID:     userid,
				}
				publishJob(ctx, app, thumbJob)
			} else if picType == PicVideo {
				// Video Poster Job
				videoJob := mediaworker.MediaJob{
					JobID:      generateUniqueID(),
					Type:       "video",
					SavedPath:  savedPath,
					UploadDir:  ResolvePath(EntityType(entityType), picType),
					PosterDir:  ResolvePath(EntityType(entityType), PicThumb),
					UniqueID:   att.Filename,
					Filename:   att.Filename,
					Ext:        att.Extension,
					ThumbWidth: defaultThumbWidth,
					UserID:     userid,
				}
				publishJob(ctx, app, videoJob)
			}
		}

		utils.RespondWithJSON(w, http.StatusOK, convertToAttachments(attachments))
	}
}

// Helper to marshal and publish media jobs safely
func publishJob(ctx context.Context, app *infra.Deps, job mediaworker.MediaJob) {
	data, err := json.Marshal(job)
	if err != nil {
		log.Printf("[Filedrop] failed to marshal job %s: %v", job.JobID, err)
		return
	}
	if err := mq.PublishWithMeta(ctx, app.MQ, "media.jobs", data); err != nil {
		log.Printf("[Filedrop] failed to publish media job %s: %v", job.JobID, err)
	}
}
