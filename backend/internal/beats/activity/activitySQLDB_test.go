package activity

import "testing"

func TestAnalyticsEventRowFromPayload(t *testing.T) {
	event := map[string]any{
		"type": "pageview",
		"data": map[string]any{"path": "/"},
		"ts":   123,
	}

	row, ok := analyticsEventRowFromPayload(event, "user-1", "session-1", "https://example.test/", "127.0.0.1")
	if !ok {
		t.Fatal("expected metadata payload to produce a row")
	}
	if row.EventName != "pageview" {
		t.Fatalf("expected event name pageview, got %q", row.EventName)
	}
	if row.EntityType != "event" {
		t.Fatalf("expected entity type event, got %q", row.EntityType)
	}
	if row.EntityID != "session-1" {
		t.Fatalf("expected entity id session-1, got %q", row.EntityID)
	}
	if string(row.Payload) == "" {
		t.Fatal("expected payload JSON to be set")
	}
}
