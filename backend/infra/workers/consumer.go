package workers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	mq "scav/infra/mq"
	log "scav/utils/logger"
)

// MediaJob is the job schema for media processing.
type MediaJob struct {
	JobID         string `json:"job_id"`
	Type          string `json:"type"` // "video" | "audio" | "image"
	SavedPath     string `json:"saved_path"`
	UploadDir     string `json:"upload_dir"`
	UniqueID      string `json:"unique_id"`
	PosterDir     string `json:"poster_dir"`
	ThumbnailPath string `json:"thumbnail_path"`
	Filename      string `json:"filename"`
	Ext           string `json:"ext"`
	ThumbWidth    int    `json:"thumb_width"`
	UserID        string `json:"user_id"`
}

// StartMediaWorker subscribes to the media job subject and processes messages.
// It returns the subscription so the caller can manage lifecycle (unsubscribe on shutdown).
func StartMediaWorker(ctx context.Context, mqClient mq.MQ) (mq.Subscription, error) {
	if mqClient == nil {
		return nil, fmt.Errorf("mq client is nil")
	}

	handler := func(hCtx context.Context, msg mq.Message) error {
		var job MediaJob
		// Support both enveloped messages (via PublishWithMeta) and raw job JSON.
		if env, err := mq.UnpackEnvelope(msg.Data); err == nil {
			var payloadBytes []byte
			switch p := env.Payload.(type) {
			case string:
				if decoded, derr := base64.StdEncoding.DecodeString(p); derr == nil {
					payloadBytes = decoded
				} else {
					payloadBytes = []byte(p)
				}
			case []byte:
				payloadBytes = p
			default:
				if b, merr := json.Marshal(p); merr == nil {
					payloadBytes = b
				}
			}

			if err := json.Unmarshal(payloadBytes, &job); err != nil {
				log.Printf("media worker: invalid envelope payload: %v", err)
				return fmt.Errorf("invalid envelope payload: %w", err)
			}
		} else {
			if err := json.Unmarshal(msg.Data, &job); err != nil {
				log.Printf("media worker: invalid job payload: %v", err)
				return fmt.Errorf("invalid job payload: %w", err)
			}
		}

		log.Printf("media worker: processing job %s type=%s", job.JobID, job.Type)

		switch job.Type {
		case "video":
			_, _, err := ProcessVideo(job.SavedPath, job.UploadDir, job.UniqueID, job.PosterDir, job.ThumbnailPath)
			if err != nil {
				log.Printf("media worker: video job %s failed: %v", job.JobID, err)
				return err
			}
			return nil
		case "audio":
			_, _ = ProcessAudio(job.SavedPath, job.UploadDir, job.UniqueID)
			return nil
		case "image":
			_, _, err := ProcessImage(job.SavedPath, job.UploadDir, job.Filename, job.Ext, job.ThumbWidth)
			if err != nil {
				log.Printf("media worker: image job %s failed: %v", job.JobID, err)
				return err
			}
			return nil
		case "image_metadata":
			// Open image file and decode
			f, err := os.Open(job.SavedPath)
			if err != nil {
				log.Printf("media worker: cannot open image for metadata job %s: %v", job.JobID, err)
				return err
			}
			defer f.Close()
			img, _, err := image.Decode(f)
			if err != nil {
				log.Printf("media worker: decode failed for metadata job %s: %v", job.JobID, err)
				return err
			}
			if err := ExtractImageMetadata(img, job.UniqueID); err != nil {
				log.Printf("media worker: ExtractImageMetadata failed for job %s: %v", job.JobID, err)
				return err
			}
			return nil
		default:
			log.Printf("media worker: unknown job type %q", job.Type)
			return fmt.Errorf("unknown job type: %s", job.Type)
		}
	}

	// Use a queue subscription so multiple worker instances can share work.
	sub, err := mqClient.QueueSubscribe(ctx, "media.jobs", "media-workers", handler)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to media.jobs: %w", err)
	}

	log.Printf("media worker: subscribed to media.jobs")
	return sub, nil
}
