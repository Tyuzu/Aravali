// File: internal/bustime/bustimeModels.go

package bustime

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type BusStop struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Code    string  `json:"code,omitempty"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
	Order   int     `json:"order"`
	Arrival string  `json:"arrival,omitempty"`
}

type BusVehicle struct {
	ID       string  `json:"id"`
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
	Status   string  `json:"status"`
	Capacity int     `json:"capacity"`
}

type BusRoute struct {
	RouteID         string       `json:"route_id"`
	Name            string       `json:"name"`
	Alias           string       `json:"alias,omitempty"`
	DepotID         string       `json:"depot_id"`
	DepotName       string       `json:"depot_name"`
	Status          string       `json:"status"`
	NextDeparture   string       `json:"next_departure,omitempty"`
	HeadwayMinutes  int          `json:"headway_minutes"`
	LastUpdated     time.Time    `json:"last_updated"`
	Stops           []BusStop    `json:"stops"`
	Vehicles        []BusVehicle `json:"vehicles,omitempty"`
	DistanceMinutes int          `json:"distance_minutes,omitempty"`
}

type BusDepot struct {
	DepotID     string     `json:"depot_id"`
	Name        string     `json:"name"`
	Location    string     `json:"location"`
	Phone       string     `json:"phone,omitempty"`
	Vehicles    int        `json:"vehicles"`
	Latitude    float64    `json:"latitude"`
	Longitude   float64    `json:"longitude"`
	OpenHours   string     `json:"open_hours"`
	Description string     `json:"description,omitempty"`
	Routes      []BusRoute `json:"routes,omitempty"`
}

type BusOverviewResponse struct {
	Routes      []BusRoute `json:"routes"`
	TotalRoutes int        `json:"total_routes"`
	GeneratedAt time.Time  `json:"generated_at"`
}

type BusMapResponse struct {
	CenterLat float64    `json:"center_lat"`
	CenterLng float64    `json:"center_lng"`
	Zoom      float64    `json:"zoom"`
	Stops     []BusStop  `json:"stops"`
	Routes    []BusRoute `json:"routes"`
}

type BusTimingResponse struct {
	Routes    []BusRoute `json:"routes"`
	Generated time.Time  `json:"generated"`
}

type BusLiveResponse struct {
	RouteID   string       `json:"route_id"`
	Name      string       `json:"name"`
	Vehicles  []BusVehicle `json:"vehicles"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type ArrivalReportRequest struct {
	RouteID     string `json:"route_id"`
	StopID      string `json:"stop_id"`
	VehicleID   string `json:"vehicle_id,omitempty"`
	DelayMins   int    `json:"delay_minutes"`
	UserComment string `json:"user_comment,omitempty"`
}

type ArrivalReportResponse struct {
	ReportID  string    `json:"report_id"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

func fallbackBusRoutes() []BusRoute {
	return []BusRoute{
		{
			RouteID:        "HR-101",
			Name:           "Narnaul - Gurugram Express",
			Alias:          "NRN-GUG",
			DepotID:        "dep-narnaul",
			DepotName:      "Narnaul Depot",
			Status:         "On Time",
			NextDeparture:  "07:30",
			HeadwayMinutes: 20,
			LastUpdated:    time.Now().UTC(),
			Stops: []BusStop{
				{ID: "s1", Name: "Narnaul Bus Stand", Code: "NRN", Lat: 28.044, Lng: 76.107, Order: 1, Arrival: "07:30"},
				{ID: "s2", Name: "Rewari Junction", Code: "REW", Lat: 28.183, Lng: 76.618, Order: 2, Arrival: "08:25"},
				{ID: "s3", Name: "IFFCO Chowk Gurugram", Code: "GUG", Lat: 28.472, Lng: 77.072, Order: 3, Arrival: "09:30"},
			},
			Vehicles: []BusVehicle{{ID: "HR-35-A-1234", Lat: 28.115, Lng: 76.350, Status: "On route", Capacity: 65}},
		},
		{
			RouteID:        "HR-102",
			Name:           "Rewari - Delhi ISBT",
			Alias:          "REW-DEL",
			DepotID:        "dep-rewari",
			DepotName:      "Rewari Depot",
			Status:         "Delayed",
			NextDeparture:  "08:15",
			HeadwayMinutes: 15,
			LastUpdated:    time.Now().UTC(),
			Stops: []BusStop{
				{ID: "s10", Name: "Rewari Bus Stand", Code: "REW", Lat: 28.183, Lng: 76.618, Order: 1, Arrival: "08:15"},
				{ID: "s11", Name: "Dharuhera Chowk", Code: "DHR", Lat: 28.210, Lng: 76.790, Order: 2, Arrival: "08:45"},
				{ID: "s12", Name: "Kashmere Gate ISBT", Code: "DEL", Lat: 28.667, Lng: 77.228, Order: 3, Arrival: "10:15"},
			},
			Vehicles: []BusVehicle{{ID: "HR-36-B-5678", Lat: 28.195, Lng: 76.700, Status: "Approaching stop", Capacity: 85}},
		},
		{
			RouteID:        "HR-103",
			Name:           "Narnaul - Mahendragarh Circular",
			Alias:          "NRN-MHG",
			DepotID:        "dep-narnaul",
			DepotName:      "Narnaul Depot",
			Status:         "On Time",
			NextDeparture:  "08:00",
			HeadwayMinutes: 30,
			LastUpdated:    time.Now().UTC(),
			Stops: []BusStop{
				{ID: "s20", Name: "Narnaul Bus Stand", Code: "NRN", Lat: 28.044, Lng: 76.107, Order: 1, Arrival: "08:00"},
				{ID: "s21", Name: "Ateli Mandi", Code: "ATL", Lat: 28.112, Lng: 76.262, Order: 2, Arrival: "08:25"},
				{ID: "s22", Name: "Mahendragarh Fort", Code: "MHG", Lat: 28.272, Lng: 76.151, Order: 3, Arrival: "09:00"},
			},
			Vehicles: []BusVehicle{{ID: "HR-35-C-9012", Lat: 28.080, Lng: 76.180, Status: "On route", Capacity: 40}},
		},
	}
}

func fallbackBusDepots() []BusDepot {
	return []BusDepot{
		{
			DepotID:     "dep-narnaul",
			Name:        "Narnaul Main Depot",
			Location:    "Near Mahavir Chowk, Narnaul",
			Phone:       "+91-1282-250100",
			Vehicles:    45,
			Latitude:    28.044,
			Longitude:   76.107,
			OpenHours:   "04:30 - 22:30",
			Description: "Central Haryana Roadways hub for Mahendragarh district.",
		},
		{
			DepotID:     "dep-rewari",
			Name:        "Rewari Central Depot",
			Location:    "Circular Road, Rewari",
			Phone:       "+91-1274-225200",
			Vehicles:    60,
			Latitude:    28.183,
			Longitude:   76.618,
			OpenHours:   "04:00 - 23:00",
			Description: "Major intersection depot serving NCR and inter-state corridors.",
		},
	}
}

func resolveBustimeCSVPath() string {
	paths := []string{
		filepath.Join(".", "SRC_to_DEST_bus_timetable.csv"),
		filepath.Join(".", "data", "SRC_to_DEST_bus_timetable.csv"),
		filepath.Join("..", "SRC_to_DEST_bus_timetable.csv"),
		filepath.Join("backend", "internal", "bustime", "SRC_to_DEST_bus_timetable.csv"),
		filepath.Join("backend", "data", "SRC_to_DEST_bus_timetable.csv"),
	}

	for _, candidate := range paths {
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return candidate
		}
	}

	for _, dir := range []string{".", filepath.Join("backend", "internal", "bustime"), filepath.Join("backend"), filepath.Join("..") } {
		for _, name := range []string{"SRC_to_DEST_bus_timetable.csv", filepath.Join("data", "SRC_to_DEST_bus_timetable.csv")} {
			candidate := filepath.Join(dir, name)
			if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
				return candidate
			}
		}
	}

	return ""
}

func parseBusStopToken(token string) (BusStop, error) {
	parts := strings.Split(token, "|")
	if len(parts) == 7 {
		lat, err := strconv.ParseFloat(strings.TrimSpace(parts[3]), 64)
		if err != nil {
			return BusStop{}, err
		}
		lng, err := strconv.ParseFloat(strings.TrimSpace(parts[4]), 64)
		if err != nil {
			return BusStop{}, err
		}
		order, err := strconv.Atoi(strings.TrimSpace(parts[5]))
		if err != nil {
			return BusStop{}, err
		}
		return BusStop{
			ID:      strings.TrimSpace(parts[0]),
			Name:    strings.TrimSpace(parts[1]),
			Code:    strings.TrimSpace(parts[2]),
			Lat:     lat,
			Lng:     lng,
			Order:   order,
			Arrival: strings.TrimSpace(parts[6]),
		}, nil
	}
	return BusStop{}, fmt.Errorf("invalid stop token: %q", token)
}

func parseStops(raw string) []BusStop {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	items := strings.Split(raw, ";")
	stops := make([]BusStop, 0, len(items))
	for _, item := range items {
		stop, err := parseBusStopToken(strings.TrimSpace(item))
		if err == nil && stop.Name != "" {
			stops = append(stops, stop)
		}
	}
	return stops
}

func loadBusRoutesFromCSV(csvPath string) ([]BusRoute, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) < 2 {
		return nil, fmt.Errorf("CSV has no data rows")
	}

	header := records[0]
	indices := map[string]int{}
	for idx, name := range header {
		indices[strings.ToLower(strings.TrimSpace(name))] = idx
	}
	if len(indices) == 0 {
		return nil, fmt.Errorf("CSV header is empty")
	}

	getField := func(row []string, name string) string {
		idx, ok := indices[name]
		if !ok {
			for key, value := range indices {
				if strings.Contains(key, name) || strings.Contains(name, key) {
					idx = value
					ok = true
					break
				}
			}
		}
		if !ok || idx >= len(row) {
			return ""
		}
		return strings.TrimSpace(row[idx])
	}

	routes := make([]BusRoute, 0, len(records)-1)
	for _, row := range records[1:] {
		if len(row) == 0 || strings.TrimSpace(strings.Join(row, "")) == "" {
			continue
		}

		routeID := getField(row, "route_id")
		if routeID == "" {
			routeID = getField(row, "routeid")
		}
		if routeID == "" {
			continue
		}

		name := getField(row, "name")
		source := getField(row, "source")
		destination := getField(row, "destination")
		if name == "" && source != "" && destination != "" {
			name = fmt.Sprintf("%s - %s", source, destination)
		}
		if name == "" {
			name = routeID
		}

		alias := getField(row, "alias")
		if alias == "" && source != "" && destination != "" {
			alias = fmt.Sprintf("%s-%s", strings.ToUpper(source[:minInt(3, len(source))]), strings.ToUpper(destination[:minInt(3, len(destination))]))
		}

		depotID := getField(row, "depot_id")
		if depotID == "" {
			depotID = getField(row, "depotid")
		}
		depotName := getField(row, "depot_name")
		if depotName == "" {
			depotName = getField(row, "depotname")
		}
		if depotName == "" {
			depotName = depotID
		}

		status := getField(row, "status")
		if status == "" {
			status = "On Time"
		}

		nextDeparture := getField(row, "next_departure")
		if nextDeparture == "" {
			nextDeparture = getField(row, "departure_time")
		}
		if nextDeparture == "" {
			nextDeparture = "--:--"
		}

		headway := getField(row, "headway_minutes")
		if headway == "" {
			headway = getField(row, "headway")
		}
		headwayMinutes := 20
		if headway != "" {
			if parsed, err := strconv.Atoi(headway); err == nil {
				headwayMinutes = parsed
			}
		}

		stops := parseStops(getField(row, "stops"))
		if len(stops) == 0 {
			if source != "" && destination != "" {
				stops = []BusStop{
					{ID: "src", Name: source, Code: strings.ToUpper(source[:3]), Order: 1, Arrival: nextDeparture},
					{ID: "dst", Name: destination, Code: strings.ToUpper(destination[:3]), Order: 2, Arrival: nextDeparture},
				}
			}
		}

		routes = append(routes, BusRoute{
			RouteID:        routeID,
			Name:           name,
			Alias:          alias,
			DepotID:        depotID,
			DepotName:      depotName,
			Status:         status,
			NextDeparture:  nextDeparture,
			HeadwayMinutes: headwayMinutes,
			LastUpdated:    time.Now().UTC(),
			Stops:          stops,
		})
	}

	return routes, nil
}

func loadCSVBusRoutes() []BusRoute {
	csvPath := resolveBustimeCSVPath()
	if csvPath == "" {
		return fallbackBusRoutes()
	}

	routes, err := loadBusRoutesFromCSV(csvPath)
	if err != nil || len(routes) == 0 {
		return fallbackBusRoutes()
	}
	return routes
}

func buildBusDepots(routes []BusRoute) []BusDepot {
	depotsByID := map[string]BusDepot{}
	for _, route := range routes {
		depotID := route.DepotID
		if depotID == "" {
			continue
		}
		if _, exists := depotsByID[depotID]; !exists {
			depotsByID[depotID] = BusDepot{
				DepotID:  depotID,
				Name:     route.DepotName,
				Location: route.DepotName,
				Vehicles: 0,
				Latitude: 0,
				Longitude: 0,
			}
		}
	}

	ids := make([]string, 0, len(depotsByID))
	for id := range depotsByID {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	depots := make([]BusDepot, 0, len(ids))
	for _, id := range ids {
		depots = append(depots, depotsByID[id])
	}
	return depots
}

var dummyBusRoutes = loadCSVBusRoutes()
var dummyDepots = buildBusDepots(dummyBusRoutes)

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
