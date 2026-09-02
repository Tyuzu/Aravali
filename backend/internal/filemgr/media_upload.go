package filemgr

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"scav/infra"
	"scav/infra/mq"
	mediaworker "scav/infra/workers"
	log "scav/utils/logger"
	"strings"
)

type MediaType string

const (
	Video MediaType = "video"
	Audio MediaType = "audio"
)

type MediaResult struct {
	Resolutions []int
	Paths       []string
	IDs         []string
}

type mediaProcessor func(r *http.Request, savedPath, uploadDir, uniqueID string, entity EntityType) ([]int, []string, error)

var mediaPicTypes = map[MediaType]PictureType{
	Video: PicVideo,
	Audio: PicAudio,
}

var mediaProcessors = map[MediaType]mediaProcessor{
	Video: func(r *http.Request, savedPath, uploadDir, uniqueID string, entity EntityType) ([]int, []string, error) {
		var tmpThumbPath string
		thumbnailFile, _, thumbErr := r.FormFile("thumbnail")
		if thumbErr == nil {
			defer thumbnailFile.Close()
			tmpThumb, err := os.CreateTemp("", uniqueID+"_thumb-*")
			if err == nil {
				if _, err := io.Copy(tmpThumb, thumbnailFile); err == nil {
					tmpThumbPath = tmpThumb.Name()
				}
				_ = tmpThumb.Close()
			}
		}

		posterDir := ResolvePath(entity, PicThumb)
		res, paths, err := mediaworker.ProcessVideo(savedPath, uploadDir, uniqueID, posterDir, tmpThumbPath)
		if tmpThumbPath != "" {
			_ = os.Remove(tmpThumbPath)
		}
		return res, paths, err
	},
	Audio: func(r *http.Request, savedPath, uploadDir, uniqueID string, entity EntityType) ([]int, []string, error) {
		res, paths := mediaworker.ProcessAudio(savedPath, uploadDir, uniqueID)
		return res, paths, nil
	},
}

func ProcessMediaUpload(app *infra.Deps, r *http.Request, formKey string, mediaType MediaType, entity EntityType, userid string) (*MediaResult, error) {
	file, err := getUploadedFile(r, formKey)
	if err != nil || file == nil {
		return nil, fmt.Errorf("no file uploaded: %w", err)
	}

	picType, ok := mediaPicTypes[mediaType]
	if !ok {
		return nil, fmt.Errorf("unsupported media type: %s", mediaType)
	}

	log.Println("ProcessMediaUpload :", picType)

	savedPath, uniqueID, _, err := SaveUploadedFile(file, entity, picType, userid)
	if err != nil {
		return nil, err
	}

	processor, ok := mediaProcessors[mediaType]
	if !ok {
		return nil, fmt.Errorf("no processor for media type: %s", mediaType)
	}

	// If MQ is available and the processor is the video processor, enqueue a job.
	if app != nil && app.MQ != nil && mediaType == Video {
		// reuse the existing inline video processor behavior to extract optional thumbnail
		var tmpThumbPath string
		thumbnailFile, _, thumbErr := r.FormFile("thumbnail")
		if thumbErr == nil {
			defer thumbnailFile.Close()
			tmpThumb, err := os.CreateTemp("", uniqueID+"_thumb-*")
			if err == nil {
				if _, err := io.Copy(tmpThumb, thumbnailFile); err == nil {
					tmpThumbPath = tmpThumb.Name()
				}
				_ = tmpThumb.Close()
			}
		}

		job := mediaworker.MediaJob{
			JobID:         generateUniqueID(),
			Type:          "video",
			SavedPath:     savedPath,
			UploadDir:     ResolvePath(entity, picType),
			PosterDir:     ResolvePath(entity, PicThumb),
			ThumbnailPath: tmpThumbPath,
			UniqueID:      uniqueID,
			Filename:      uniqueID,
			Ext:           ".mp4",
			ThumbWidth:    defaultThumbWidth,
			UserID:        userid,
		}
		_ = mq.PublishWithMeta(r.Context(), app.MQ, "media.jobs", job)

		return &MediaResult{Resolutions: nil, Paths: nil, IDs: []string{uniqueID}}, nil
	}

	res, paths, err := processor(r, savedPath, ResolvePath(entity, picType), uniqueID, entity)
	if err != nil {
		return nil, err
	}

	return &MediaResult{
		Resolutions: res,
		Paths:       paths,
		IDs:         []string{uniqueID},
	}, nil
}

func getUploadedFile(r *http.Request, formKey string) (*multipart.FileHeader, error) {
	if r.MultipartForm == nil {
		r.Body = http.MaxBytesReader(nil, r.Body, 32<<20)
		if err := r.ParseMultipartForm(32 << 20); err != nil { // #nosec G120
			return nil, fmt.Errorf("failed to parse form: %w", err)
		}
	}
	files := r.MultipartForm.File[formKey]
	if len(files) == 0 {
		return nil, nil
	}
	return files[0], nil
}

func SaveUploadedFile(file *multipart.FileHeader, entity EntityType, picType PictureType, userid string) (string, string, string, error) {
	src, err := file.Open()
	if err != nil {
		return "", "", "", fmt.Errorf("cannot open uploaded file: %w", err)
	}
	defer src.Close()

	log.Println("SaveUploadedFile :", picType)
	savedName, ext, err := SaveFileForEntity(src, file, entity, picType, userid)
	if err != nil {
		return "", "", "", fmt.Errorf("file save failed: %w", err)
	}

	savedPath := filepath.Join(ResolvePath(entity, picType), savedName+ext)
	uniqueID := strings.TrimSuffix(savedName, ext)
	return savedPath, uniqueID, ext, nil
}
