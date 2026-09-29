// File: infra/workers/consumer.go

package workers

import (
	"context"
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
	Type          string `json:"type"` // "video" | "audio" | "image" | "image_metadata"
	SavedPath     string `json:"saved_path"`
	UploadDir     string `json:"upload_dir"`
	UniqueID      string `json:"unique_id"`
	PosterDir     string `json:"poster_dir"`
	ThumbnailPath string `json:"thumbnail_path"`
	Filename      string `json:"filename"`
	Ext           string `json:"ext"`
	ThumbWidth    int    `json:"thumb_width"`
	UserID        string `json:"userid"`
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
			// Re-inject trace & source metadata into context if present
			if env.TraceID != "" {
				hCtx = mq.WithTraceID(hCtx, env.TraceID)
			}
			if env.Source != "" {
				hCtx = mq.WithServiceName(hCtx, env.Source)
			}

			if err := json.Unmarshal(env.Payload, &job); err != nil {
				log.L.Sugar().Errorw("media worker: invalid envelope payload", "error", err)
				return fmt.Errorf("invalid envelope payload: %w", err)
			}
		} else {
			if err := json.Unmarshal(msg.Data, &job); err != nil {
				log.L.Sugar().Errorw("media worker: invalid job payload", "error", err)
				return fmt.Errorf("invalid job payload: %w", err)
			}
		}

		traceID := mq.TraceIDFromContext(hCtx)
		log.L.Sugar().Infow("media worker: processing job", "job_id", job.JobID, "type", job.Type, "trace_id", traceID)

		switch job.Type {
		case "video":
			_, _, err := ProcessVideo(job.SavedPath, job.UploadDir, job.UniqueID, job.PosterDir, job.ThumbnailPath)
			if err != nil {
				log.L.Sugar().Errorw("media worker: video job failed", "job_id", job.JobID, "error", err)
				return err
			}
			return nil

		case "audio":
			_, _ = ProcessAudio(job.SavedPath, job.UploadDir, job.UniqueID)
			return nil

		case "image":
			_, _, err := ProcessImage(job.SavedPath, job.UploadDir, job.Filename, job.Ext, job.ThumbWidth)
			if err != nil {
				log.L.Sugar().Errorw("media worker: image job failed", "job_id", job.JobID, "error", err)
				return err
			}
			return nil

		case "image_metadata":
			f, err := os.Open(job.SavedPath)
			if err != nil {
				log.L.Sugar().Errorw("media worker: cannot open image for metadata job", "job_id", job.JobID, "error", err)
				return err
			}
			defer f.Close()

			img, _, err := image.Decode(f)
			if err != nil {
				log.L.Sugar().Errorw("media worker: decode failed for metadata job", "job_id", job.JobID, "error", err)
				return err
			}

			if err := ExtractImageMetadata(img, job.UniqueID); err != nil {
				log.L.Sugar().Errorw("media worker: ExtractImageMetadata failed", "job_id", job.JobID, "error", err)
				return err
			}
			return nil

		default:
			log.L.Sugar().Warnw("media worker: unknown job type", "job_type", job.Type)
			return fmt.Errorf("unknown job type: %s", job.Type)
		}
	}

	// Use a queue subscription so multiple worker instances can share work.
	sub, err := mqClient.QueueSubscribe(ctx, "media.jobs", "media-workers", handler)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to media.jobs: %w", err)
	}

	// Automatically unsubscribe when the application lifecycle context stops
	go func() {
		<-ctx.Done()
		if err := sub.Unsubscribe(); err != nil {
			log.L.Sugar().Warnw("media worker: failed to unsubscribe on context cancellation", "error", err)
		}
	}()

	log.L.Sugar().Infow("media worker: subscribed to media.jobs")
	return sub, nil
}
