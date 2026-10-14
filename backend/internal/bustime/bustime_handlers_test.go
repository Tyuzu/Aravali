// File: internal/bustime/bustime_handlers_test.go

package bustime

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetBusOverviewReturnsData(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/bus", nil)
	res := httptest.NewRecorder()

	GetBusOverview(nil)(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, res.Code)
	}

	var payload map[string]any
	if err := json.Unmarshal(res.Body.Bytes(), &payload); err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	if payload["routes"] == nil {
		t.Fatal("expected routes in payload")
	}
}

func TestGetBusOverviewQueryFiltering(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/bus?q=Gurugram", nil)
	res := httptest.NewRecorder()

	GetBusOverview(nil)(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, res.Code)
	}

	var payload BusOverviewResponse
	if err := json.Unmarshal(res.Body.Bytes(), &payload); err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	if payload.TotalRoutes != 1 {
		t.Fatalf("expected 1 route, got %d", payload.TotalRoutes)
	}
}

func TestGetBusRouteReturnsRoute(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/bus/time/HR-101", nil)
	res := httptest.NewRecorder()

	GetBusRoute(nil)(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, res.Code)
	}

	var payload map[string]any
	if err := json.Unmarshal(res.Body.Bytes(), &payload); err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	if payload["route_id"] == nil {
		t.Fatal("expected route_id in payload")
	}
}

func TestPostArrivalReportValidation(t *testing.T) {
	body, _ := json.Marshal(map[string]string{"route_id": "HR-101", "stop_id": "s1"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/bus/report", bytes.NewBuffer(body))
	res := httptest.NewRecorder()

	PostArrivalReport(nil)(res, req)

	if res.Code != http.StatusCreated {
		t.Fatalf("expected %d, got %d", http.StatusCreated, res.Code)
	}
}
