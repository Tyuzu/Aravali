// File: internal/bustime/bustime_handlers_test.go

package bustime

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/julienschmidt/httprouter"
)

func TestLoadBusRoutesFromCSVFile(t *testing.T) {
	tmpDir := t.TempDir()
	csvPath := filepath.Join(tmpDir, "SRC_to_DEST_bus_timetable.csv")
	content := "route_id,name,alias,depot_id,depot_name,status,next_departure,headway_minutes,stops\n" +
		"HR-901,Ambala - Chandigarh,AMB-CHD,dep-ambala,Ambala Depot,On Time,07:45,15,\"a1|Ambala Stand|AMB|30.378|76.776|1|07:45;a2|Chandigarh Sector 17|CHD|30.733|76.788|2|08:15\"\n"
	if err := os.WriteFile(csvPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write csv: %v", err)
	}

	routes, err := loadBusRoutesFromCSV(csvPath)
	if err != nil {
		t.Fatalf("load csv: %v", err)
	}
	if len(routes) != 1 {
		t.Fatalf("expected 1 route, got %d", len(routes))
	}
	if routes[0].Name != "Ambala - Chandigarh" {
		t.Fatalf("unexpected route name: %s", routes[0].Name)
	}
	if len(routes[0].Stops) != 2 {
		t.Fatalf("expected 2 stops, got %d", len(routes[0].Stops))
	}
}

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

func TestBusRoutesRemainPublicWithoutAuth(t *testing.T) {
	router := httprouter.New()
	AddBustimeRoutes(router, nil, nil)

	routeID := dummyBusRoutes[0].RouteID
	depotID := dummyDepots[0].DepotID

	tests := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/v1/bus"},
		{http.MethodGet, "/api/v1/bus/time"},
		{http.MethodGet, "/api/v1/bus/map"},
		{http.MethodGet, "/api/v1/bus/time/" + routeID},
		{http.MethodGet, "/api/v1/bus/live/" + routeID},
		{http.MethodGet, "/api/v1/bus/depot/" + depotID},
	}

	for _, tc := range tests {
		req := httptest.NewRequest(tc.method, tc.path, nil)
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		if res.Code == http.StatusUnauthorized || res.Code == http.StatusForbidden {
			t.Fatalf("route %s rejected without auth: status=%d", tc.path, res.Code)
		}
		if res.Code != http.StatusOK {
			t.Fatalf("route %s expected status %d, got %d", tc.path, http.StatusOK, res.Code)
		}
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
