package baito

import (
	"testing"
	"time"
)

func TestBuildMyApplicationsResult(t *testing.T) {
	applications := []map[string]any{
		{
			"_id":         "app-1",
			"baitoid":     "job-1",
			"pitch":       "Interested in this role",
			"submittedAt": time.Date(2024, 1, 10, 12, 0, 0, 0, time.UTC),
		},
	}

	jobs := []Baito{
		{
			BaitoId:  "job-1",
			Title:    "Warehouse assistant",
			Location: "Nairobi",
			Wage:     "KSh 20,000",
		},
	}

	got := buildMyApplicationsResult(applications, jobs)
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}

	if got[0]["jobId"] != "job-1" {
		t.Fatalf("expected jobId job-1, got %#v", got[0]["jobId"])
	}
	if got[0]["title"] != "Warehouse assistant" {
		t.Fatalf("expected title Warehouse assistant, got %#v", got[0]["title"])
	}
	if got[0]["location"] != "Nairobi" {
		t.Fatalf("expected location Nairobi, got %#v", got[0]["location"])
	}
	if got[0]["wage"] != "KSh 20,000" {
		t.Fatalf("expected wage KSh 20,000, got %#v", got[0]["wage"])
	}
}
