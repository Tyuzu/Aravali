// File: internal/bustime/bustimeModels.go

package bustime

import "time"

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

var dummyBusRoutes = []BusRoute{
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

var dummyDepots = []BusDepot{
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
