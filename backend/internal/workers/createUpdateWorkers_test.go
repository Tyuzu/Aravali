package workers

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParseWorkerFormAcceptsJSONPayload(t *testing.T) {
	req := httptest.NewRequest("POST", "/baitos/profile", strings.NewReader(`{"name":"Demo worker","skills":["farming","delivery"]}`))
	req.Header.Set("Content-Type", "application/json")

	worker, _, err := parseWorkerForm(req, false)
	if err != nil {
		t.Fatalf("expected JSON worker payload to parse without error, got: %v", err)
	}
	if worker.Name != "Demo worker" {
		t.Fatalf("expected name to be preserved, got %q", worker.Name)
	}
	if len(worker.Preferred) != 0 {
		t.Fatalf("expected preferred roles to remain empty when not provided, got %#v", worker.Preferred)
	}
}
