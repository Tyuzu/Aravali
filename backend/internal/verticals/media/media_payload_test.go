package media

import "testing"

func TestFilePayloadNormalizedExtn(t *testing.T) {
	payload := FilePayload{Filename: "abc123", Extn: ".jpg"}
	if got := payload.NormalizedExtn(); got != ".jpg" {
		t.Fatalf("expected .jpg, got %q", got)
	}

	payload = FilePayload{Filename: "abc123.png", Extension: ".png"}
	if got := payload.NormalizedExtn(); got != ".png" {
		t.Fatalf("expected .png, got %q", got)
	}

	payload = FilePayload{Filename: "abc123"}
	if got := payload.NormalizedExtn(); got != "" {
		t.Fatalf("expected empty extn, got %q", got)
	}
}
