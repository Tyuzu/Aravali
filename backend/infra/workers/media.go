package workers

import (
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"os"
	"path/filepath"
)

// MediaJob is a light-weight compatibility stub for media processing.
type MediaJob struct {
	JobID         string
	Type          string
	SavedPath     string
	UploadDir     string
	PosterDir     string
	ThumbnailPath string
	UniqueID      string
	Filename      string
	Ext           string
	ThumbWidth    int
	UserID        string
}

func ProcessVideo(savedPath, uploadDir, uniqueID, posterDir, thumbPath string) ([]int, []string, error) {
	return nil, []string{savedPath}, nil
}

func ProcessAudio(savedPath, uploadDir, uniqueID string) ([]int, []string) {
	return nil, []string{savedPath}
}

func EncodeJPEG(img image.Image, dst string, quality int) error {
	if quality <= 0 {
		quality = 80
	}
	if quality > 100 {
		quality = 100
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o750); err != nil {
		return err
	}
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()
	return jpeg.Encode(f, img, &jpeg.Options{Quality: quality})
}

func SaveAtomically(dst string, writeFn func(*os.File) error) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o750); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(dst), filepath.Base(dst)+".tmp-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer func() {
		_ = os.Remove(tmpName)
	}()
	if err := writeFn(tmp); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpName, dst)
}

func ProcessImage(srcPath, thumbDir, filename, ext string, thumbWidth int) (string, string, error) {
	_ = thumbWidth
	if err := os.MkdirAll(thumbDir, 0o750); err != nil {
		return "", "", err
	}
	if ext == "" {
		ext = filepath.Ext(srcPath)
	}
	if filename == "" {
		filename = filepath.Base(srcPath)
	}
	outPath := filepath.Join(thumbDir, filename+ext)
	in, err := os.Open(srcPath)
	if err != nil {
		return "", "", err
	}
	defer in.Close()
	out, err := os.Create(outPath)
	if err != nil {
		return "", "", err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return "", "", err
	}
	return outPath, ext, nil
}

func (m MediaJob) String() string {
	return fmt.Sprintf("job=%s type=%s uniqueID=%s", m.JobID, m.Type, m.UniqueID)
}
