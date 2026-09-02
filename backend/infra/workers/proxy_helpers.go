package workers

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"time"
)

var proxyEncSem = make(chan struct{}, 6)

// SaveAtomically writes to a temporary file in the same directory and renames it.
func SaveAtomically(path string, writeFn func(*os.File) error) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, ".cache-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	cleanup := true
	defer func() {
		_ = tmp.Close()
		if cleanup {
			_ = os.Remove(tmpName)
		}
	}()

	if err := writeFn(tmp); err != nil {
		return err
	}
	if err := tmp.Sync(); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmpName, path); err != nil {
		return err
	}
	cleanup = false
	return nil
}

// EncodeJPEG encodes the provided image as a JPEG to cachePath using an atomic save.
func EncodeJPEG(img image.Image, cachePath string, quality int) error {
	select {
	case proxyEncSem <- struct{}{}:
		defer func() { <-proxyEncSem }()
	case <-time.After(5 * time.Second):
		return fmt.Errorf("encoding semaphore timeout")
	}

	return SaveAtomically(cachePath, func(f *os.File) error {
		return jpeg.Encode(f, img, &jpeg.Options{Quality: quality})
	})
}
