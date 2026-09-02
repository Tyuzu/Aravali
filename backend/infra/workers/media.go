package workers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	log "scav/utils/logger"

	"github.com/disintegration/imaging"
)

const (
	ffprobeTimeout    = 30 * time.Second
	transcodeTimeout  = 10 * time.Minute
	posterTimeout     = 45 * time.Second
	audioTimeout      = 3 * time.Minute
	defaultThumbWidth = 500
	defaultQuality    = 85
)

type Runner interface {
	Run(timeout time.Duration, name string, args ...string) (stdout string, stderr string, err error)
}

type realRunner struct{}

func (realRunner) Run(timeout time.Duration, name string, args ...string) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...) // #nosec G702 G204
	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb

	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return out.String(), errb.String(), fmt.Errorf("%s timed out after %s", name, timeout)
	}
	return out.String(), errb.String(), err
}

var cmdRunner Runner = realRunner{}

func ProcessAudio(savedPath, uploadDir, uniqueID string) ([]int, []string) {
	resolutions, outputPath := processAudioResolutions(savedPath, uploadDir, uniqueID)
	var paths []string
	if outputPath != "" {
		paths = []string{normalizePath(outputPath)}
	}
	return resolutions, paths
}

func ProcessVideo(savedPath, uploadDir, uniqueID, posterDir, thumbnailPath string) ([]int, []string, error) {
	width, height, err := getVideoDimensions(savedPath)
	if err != nil {
		_ = os.Remove(savedPath) // #nosec G703
		return nil, nil, fmt.Errorf("failed to get video dimensions: %w", err)
	}

	resolutions, outputPaths := processVideoResolutionsParallel(savedPath, uploadDir, uniqueID, width, height, 3)
	if len(outputPaths) == 0 {
		_ = os.Remove(savedPath) // #nosec G703
		return nil, nil, fmt.Errorf("video transcoding failed")
	}

	if err := os.MkdirAll(posterDir, 0o750); err != nil {
		cleanupPaths(outputPaths)
		_ = os.Remove(savedPath) // #nosec G703
		return nil, nil, fmt.Errorf("failed to create poster directory: %w", err)
	}
	thumbPath := filepath.Join(posterDir, uniqueID+".jpg")
	if thumbnailPath != "" {
		args := []string{
			"-y",
			"-i", thumbnailPath,
			"-vf", "scale=w=iw*min(1280/iw\\,720/ih):h=ih*min(1280/iw\\,720/ih),pad=1280:720:(1280-iw*min(1280/iw\\,720/ih))/2:(720-ih*min(1280/iw\\,720/ih))/2:black",
			thumbPath,
		}
		stdout, stderr, err := cmdRunner.Run(time.Minute, "ffmpeg", args...)
		if err != nil {
			cleanupPaths(outputPaths)
			_ = os.Remove(savedPath) // #nosec G703
			return nil, nil, fmt.Errorf("failed to process thumbnail: %w (stdout=%s, stderr=%s)", err, stdout, stderr)
		}
	} else {
		if err := CreatePoster(savedPath, thumbPath); err != nil {
			cleanupPaths(outputPaths)
			_ = os.Remove(savedPath) // #nosec G703
			return nil, nil, fmt.Errorf("poster creation failed: %w", err)
		}
	}

	return resolutions, outputPaths, nil
}

func ProcessImage(fullPath, thumbDir, filename, ext string, thumbWidth int) (string, string, error) {
	img, _, err := openImage(fullPath)
	if err != nil {
		return fullPath, ext, fmt.Errorf("open image %s: %w", fullPath, err)
	}

	if err := ValidateImageDimensions(img, 12000, 12000); err != nil {
		return fullPath, ext, err
	}

	finalPath := fullPath
	finalExt := ext
	if strings.ToLower(ext) != ".png" {
		finalPath, err = normalizeImageFormat(fullPath, ext, img)
		if err != nil {
			return fullPath, ext, err
		}
		finalExt = filepath.Ext(finalPath)
	}

	if thumbWidth <= 0 {
		return finalPath, finalExt, fmt.Errorf("invalid thumbnail width: %d", thumbWidth)
	}

	imgCopy := imaging.Clone(img)
	if err := generateThumbnail(imgCopy, thumbDir, filename+".jpg", thumbWidth); err != nil {
		return finalPath, finalExt, fmt.Errorf("generate thumbnail %s: %w", filename, err)
	}
	return finalPath, finalExt, nil
}

func GenerateVideoPoster(videoPath, thumbDir, baseFilename string) (string, error) {
	thumbName := strings.TrimSuffix(baseFilename, filepath.Ext(baseFilename)) + ".jpg"
	thumbPath := filepath.Join(thumbDir, thumbName)
	if err := os.MkdirAll(thumbDir, 0o750); err != nil {
		return "", fmt.Errorf("mkdir %s: %w", thumbDir, err)
	}
	if err := CreatePoster(videoPath, thumbPath); err != nil {
		return "", err
	}
	return thumbName, nil
}

func openImage(path string) (image.Image, string, error) {
	f, err := os.Open(path) // #nosec G703 G304
	if err != nil {
		return nil, "", fmt.Errorf("open image: %w", err)
	}
	defer f.Close()
	img, format, err := image.Decode(f)
	return img, format, err
}

func generateThumbnail(img image.Image, thumbDir, baseFilename string, thumbWidth int) error {
	if img == nil {
		return fmt.Errorf("nil image")
	}
	if thumbWidth <= 0 {
		return fmt.Errorf("invalid thumbnail width: %d", thumbWidth)
	}

	resized := imaging.Resize(img, thumbWidth, 0, imaging.Lanczos)
	name := strings.TrimSuffix(baseFilename, filepath.Ext(baseFilename)) + ".jpg"
	path := filepath.Join(thumbDir, name)

	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil { // #nosec G703
		return fmt.Errorf("mkdir %s: %w", filepath.Dir(path), err)
	}

	out, err := os.Create(path) // #nosec G703 G304
	if err != nil {
		return fmt.Errorf("create thumbnail %s: %w", path, err)
	}
	defer out.Close()

	if err := jpeg.Encode(out, resized, &jpeg.Options{Quality: defaultQuality}); err != nil {
		_ = os.Remove(path) // #nosec G703
		return fmt.Errorf("encode thumbnail %s: %w", path, err)
	}
	if err := out.Sync(); err != nil {
		return fmt.Errorf("sync thumbnail %s: %w", path, err)
	}
	return nil
}

func normalizeImageFormat(fullPath, ext string, img image.Image) (string, error) {
	if strings.EqualFold(ext, ".png") {
		return fullPath, nil
	}
	pngPath := strings.TrimSuffix(fullPath, ext) + ".png"
	out, err := os.Create(pngPath) // #nosec G703 G304
	if err != nil {
		return fullPath, fmt.Errorf("create png %s: %w", pngPath, err)
	}
	defer out.Close()

	if err := png.Encode(out, img); err != nil {
		_ = os.Remove(pngPath) // #nosec G703
		return fullPath, fmt.Errorf("encode png: %w", err)
	}
	_ = os.Remove(fullPath) // #nosec G703
	return pngPath, nil
}

func ValidateImageDimensions(img image.Image, maxWidth, maxHeight int) error {
	if img == nil {
		return fmt.Errorf("validate dimensions: nil image")
	}
	bounds := img.Bounds()
	if bounds.Dx() > maxWidth || bounds.Dy() > maxHeight {
		return fmt.Errorf("image dimensions %dx%d exceed max %dx%d", bounds.Dx(), bounds.Dy(), maxWidth, maxHeight)
	}
	return nil
}

func StripEXIF(img image.Image) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: 90}); err != nil {
		return nil, fmt.Errorf("strip exif: encode failed: %w", err)
	}
	return buf, nil
}

func ExtractImageMetadata(img image.Image, uid string) error {
	if img == nil {
		return fmt.Errorf("extract metadata: nil image")
	}

	bounds := img.Bounds()
	buf, err := StripEXIF(img)
	if err != nil {
		return fmt.Errorf("extract metadata: encoding failed: %w", err)
	}

	size := buf.Len()
	msg := fmt.Sprintf("metadata uid=%s width=%d height=%d size=%d", uid, bounds.Dx(), bounds.Dy(), size)
	log.Println(msg)
	return nil
}

func getVideoDimensions(videoPath string) (int, int, error) {
	args := []string{
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=width,height",
		"-of", "csv=p=0",
		videoPath,
	}
	stdout, stderr, err := cmdRunner.Run(ffprobeTimeout, "ffprobe", args...)
	if err != nil {
		return 0, 0, fmt.Errorf("ffprobe getVideoDimensions(%s) failed: %w (stderr=%s)", videoPath, err, stderr)
	}

	parts := strings.Split(strings.TrimSpace(stdout), ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("ffprobe getVideoDimensions unexpected output for %s: %q", videoPath, stdout)
	}

	width, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, fmt.Errorf("ffprobe parse width for %s: %w", videoPath, err)
	}
	height, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, fmt.Errorf("ffprobe parse height for %s: %w", videoPath, err)
	}
	return width, height, nil
}

func processVideoResolution(inputPath, outputPath string, targetHeight int) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o750); err != nil {
		return fmt.Errorf("create output dir for %s: %w", outputPath, err)
	}

	scaleFilter := fmt.Sprintf("scale=-2:%d", targetHeight)
	args := []string{
		"-y",
		"-i", inputPath,
		"-vf", scaleFilter,
		"-c:v", "libx264",
		"-crf", "23",
		"-preset", "veryfast",
		"-tune", "zerolatency",
		"-pix_fmt", "yuv420p",
		"-max_muxing_queue_size", "9999",
		"-c:a", "aac",
		"-b:a", "128k",
		"-movflags", "+faststart",
		outputPath,
	}

	stdout, stderr, err := cmdRunner.Run(transcodeTimeout, "ffmpeg", args...)
	if err != nil {
		return fmt.Errorf("ffmpeg transcode %s -> %s (%s) failed: %w (stdout=%s, stderr=%s)", inputPath, outputPath, scaleFilter, err, stdout, stderr)
	}
	return nil
}

func CreatePoster(videoPath, posterPath string) error {
	if err := os.MkdirAll(filepath.Dir(posterPath), 0o750); err != nil { // #nosec G703
		return fmt.Errorf("failed to create poster directory for %s: %w", posterPath, err)
	}

	duration, err := getVideoDuration(videoPath)
	if err != nil || duration <= 0 {
		log.Printf("CreatePoster duration unavailable for %s: %v", videoPath, err) // #nosec G706
		duration = 3.0
	}

	t := duration * 0.25
	if t < 1.0 {
		t = 1.0
	}
	if t > duration-0.5 {
		t = math.Max(0.0, duration-0.5)
	}
	timestamp := formatTimestamp(t)

	args := []string{
		"-y",
		"-ss", timestamp,
		"-i", videoPath,
		"-vframes", "1",
		"-q:v", "2",
		"-vf", "scale=w=iw*min(1280/iw\\,720/ih):h=ih*min(1280/iw\\,720/ih),pad=1280:720:(1280-iw*min(1280/iw\\,720/ih))/2:(720-ih*min(1280/iw\\,720/ih))/2:black",
		posterPath,
	}

	stdout, stderr, err := cmdRunner.Run(posterTimeout, "ffmpeg", args...)
	if err != nil {
		return fmt.Errorf("poster creation failed for %s at %s: %w (stdout=%s, stderr=%s)", videoPath, timestamp, err, stdout, stderr)
	}
	return nil
}

func getVideoDuration(path string) (float64, error) {
	args := []string{
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "json",
		path,
	}
	stdout, stderr, err := cmdRunner.Run(ffprobeTimeout, "ffprobe", args...)
	if err != nil {
		return 0, fmt.Errorf("ffprobe getVideoDuration(%s) failed: %w (stderr=%s)", path, err, stderr)
	}

	var result struct {
		Format struct {
			Duration string `json:"duration"`
		} `json:"format"`
	}
	if err := json.Unmarshal([]byte(stdout), &result); err != nil {
		return 0, fmt.Errorf("ffprobe unmarshal duration for %s: %w (stdout=%s)", path, err, stdout)
	}
	if strings.TrimSpace(result.Format.Duration) == "" {
		return 0, fmt.Errorf("ffprobe duration not found for %s (stdout=%s)", path, stdout)
	}

	dur, err := strconv.ParseFloat(result.Format.Duration, 64)
	if err != nil {
		return 0, fmt.Errorf("parse duration for %s: %w (value=%s)", path, err, result.Format.Duration)
	}
	return dur, nil
}

func formatTimestamp(seconds float64) string {
	if seconds < 0 {
		seconds = 0
	}
	totalMs := int(seconds * 1000.0)
	h := totalMs / 3600000
	m := (totalMs % 3600000) / 60000
	s := (totalMs % 60000) / 1000
	ms := totalMs % 1000
	return fmt.Sprintf("%02d:%02d:%02d.%03d", h, m, s, ms)
}

func processAudioResolutions(originalFilePath, uploadDir, uniqueID string) ([]int, string) {
	if err := os.MkdirAll(uploadDir, 0o750); err != nil {
		log.Printf("audio: failed to create output dir %s: %v", uploadDir, err)
		return []int{}, originalFilePath
	}
	outputPath := filepath.Join(uploadDir, uniqueID+".mp3")

	inputBitrate := probeAudioBitrate(originalFilePath)

	targetKbps := 128
	if inputBitrate > 0 {
		inKbps := inputBitrate / 1000
		if inKbps < targetKbps {
			targetKbps = inKbps
		}
		if targetKbps <= 0 {
			targetKbps = 128
		}
	}

	args := []string{
		"-y",
		"-i", originalFilePath,
		"-vn",
		"-c:a", "libmp3lame",
		"-b:a", fmt.Sprintf("%dk", targetKbps),
		"-filter:a", "loudnorm",
		outputPath,
	}

	stdout, stderr, err := cmdRunner.Run(audioTimeout, "ffmpeg", args...)
	if err != nil {
		log.Printf("audio processing failed for %s -> %s: %v; stdout: %s; stderr: %s", originalFilePath, outputPath, err, stdout, stderr)
		return []int{}, originalFilePath
	}

	return []int{targetKbps}, outputPath
}

func probeAudioBitrate(path string) int {
	args := []string{
		"-v", "error",
		"-select_streams", "a:0",
		"-show_entries", "stream=bit_rate",
		"-of", "json",
		path,
	}
	stdout, _, err := cmdRunner.Run(ffprobeTimeout, "ffprobe", args...)
	if err != nil {
		return 0
	}

	var result struct {
		Streams []struct {
			BitRate json.Number `json:"bit_rate"`
		} `json:"streams"`
	}
	if err := json.Unmarshal([]byte(stdout), &result); err != nil || len(result.Streams) == 0 {
		return 0
	}

	brStr := strings.TrimSpace(string(result.Streams[0].BitRate))
	if brStr == "" {
		return 0
	}
	br, err := strconv.Atoi(brStr)
	if err != nil || br <= 0 {
		return 0
	}
	return br
}

type videoTask struct {
	Height     int
	OutputPath string
}

func processVideoResolutionsParallel(originalFilePath, uploadDir, uniqueID string, origWidth, origHeight int, maxParallel int) ([]int, []string) {
	_ = origWidth
	if maxParallel <= 0 {
		maxParallel = 2
	}

	ladder := []struct {
		Height int
	}{
		{4320}, {2160}, {1440}, {1080}, {720}, {480}, {360}, {240}, {144},
	}

	var tasks []videoTask
	for _, r := range ladder {
		if r.Height > origHeight {
			continue
		}
		tasks = append(tasks, videoTask{
			Height:     r.Height,
			OutputPath: generateFilePath(uploadDir, uniqueID+"-"+strconv.Itoa(r.Height), "mp4"),
		})
	}
	if len(tasks) == 0 {
		return nil, nil
	}

	workers := maxParallel
	if workers > len(tasks) {
		workers = len(tasks)
	}

	results := make(chan struct {
		height int
		path   string
		ok     bool
	}, len(tasks))

	jobs := make(chan videoTask)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for t := range jobs {
				if err := processVideoResolution(originalFilePath, t.OutputPath, t.Height); err != nil {
					log.Printf("Skipping %d due to error: %v", t.Height, err)
					results <- struct {
						height int
						path   string
						ok     bool
					}{ok: false}
					continue
				}
				results <- struct {
					height int
					path   string
					ok     bool
				}{height: t.Height, path: "/" + filepath.ToSlash(t.OutputPath), ok: true}
			}
		}()
	}

	go func() {
		for _, t := range tasks {
			jobs <- t
		}
		close(jobs)
		wg.Wait()
		close(results)
	}()

	type pair struct {
		h int
		p string
	}
	var pairs []pair
	for res := range results {
		if res.ok {
			pairs = append(pairs, pair{h: res.height, p: res.path})
		}
	}

	if len(pairs) == 0 {
		return nil, nil
	}

	sort.Slice(pairs, func(i, j int) bool { return pairs[i].h > pairs[j].h })

	heights := make([]int, 0, len(pairs))
	outputs := make([]string, 0, len(pairs))
	for _, pr := range pairs {
		heights = append(heights, pr.h)
		outputs = append(outputs, pr.p)
	}
	return heights, outputs
}

func cleanupPaths(paths []string) {
	for _, p := range paths {
		if p == "" {
			continue
		}
		_ = os.Remove(strings.TrimPrefix(filepath.FromSlash(p), string(filepath.Separator)))
	}
}

func normalizePath(p string) string {
	if !strings.HasPrefix(p, "/") {
		p = "/" + filepath.ToSlash(p)
	}
	return filepath.ToSlash(p)
}

func generateFilePath(baseDir, uniqueID, extension string) string {
	extension = strings.TrimPrefix(extension, ".")
	return filepath.Join(baseDir, uniqueID+"."+extension)
}

// GenerateFilePath returns a file path under baseDir for uniqueID with extension.
func GenerateFilePath(baseDir, uniqueID, extension string) string {
	return generateFilePath(baseDir, uniqueID, extension)
}

// NormalizePath makes a filesystem path canonical and ensures it begins with '/'
func NormalizePath(p string) string {
	return normalizePath(p)
}

// CleanupPaths removes a list of paths on disk.
func CleanupPaths(paths []string) {
	cleanupPaths(paths)
}
