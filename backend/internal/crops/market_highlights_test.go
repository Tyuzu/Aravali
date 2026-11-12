package crops

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetMarketHighlights(t *testing.T) {
	handler := GetMarketHighlights(nil)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/market/highlights", nil)
	res := httptest.NewRecorder()

	handler(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d; body=%s", res.Code, res.Body.String())
	}

	var payload map[string][]struct {
		ID   string `json:"id"`
		Text string `json:"text"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &payload); err != nil {
		t.Fatalf("expected valid JSON payload: %v; body=%s", err, res.Body.String())
	}

	for _, key := range []string{"activeDeals", "seasonalPicks", "cropTrends", "announcements", "showcase"} {
		if _, ok := payload[key]; !ok {
			t.Fatalf("missing key %q in payload: %#v", key, payload)
		}
		if len(payload[key]) == 0 {
			t.Fatalf("key %q should contain at least one highlight item", key)
		}
	}
}
