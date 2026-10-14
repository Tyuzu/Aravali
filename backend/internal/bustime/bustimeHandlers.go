// File: internal/bustime/bustimeHandlers.go

package bustime

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"scav/infra"
	"scav/utils"
)

func resolveBusID(r *http.Request, name string) string {
	if value := utils.GetParam(r, name); value != "" {
		return value
	}
	if value := r.URL.Query().Get(name); value != "" {
		return value
	}

	path := strings.TrimSpace(r.URL.Path)
	prefixes := []string{"/api/v1/bus/time/", "/api/v1/bus/live/", "/api/v1/bus/depot/", "/api/v1/bustime/"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(path, prefix) {
			return strings.TrimPrefix(path, prefix)
		}
	}
	return ""
}

func findBusRouteByID(routeID string) (BusRoute, bool) {
	for _, route := range dummyBusRoutes {
		if route.RouteID == routeID {
			return route, true
		}
	}
	return BusRoute{}, false
}

func findBusDepotByID(depotID string) (BusDepot, bool) {
	for _, depot := range dummyDepots {
		if depot.DepotID == depotID {
			var matchingRoutes []BusRoute
			for _, route := range dummyBusRoutes {
				if route.DepotID == depot.DepotID {
					matchingRoutes = append(matchingRoutes, route)
				}
			}
			depot.Routes = matchingRoutes
			return depot, true
		}
	}
	return BusDepot{}, false
}

func filterRoutes(routes []BusRoute, query string, status string) []BusRoute {
	query = strings.ToLower(strings.TrimSpace(query))
	status = strings.TrimSpace(status)

	if query == "" && status == "" {
		return routes
	}

	var filtered []BusRoute
	for _, route := range routes {
		matchesQuery := query == "" ||
			strings.Contains(strings.ToLower(route.Name), query) ||
			strings.Contains(strings.ToLower(route.Alias), query) ||
			strings.Contains(strings.ToLower(route.DepotName), query)

		matchesStatus := status == "" || strings.EqualFold(route.Status, status)

		if matchesQuery && matchesStatus {
			filtered = append(filtered, route)
		}
	}
	return filtered
}

func GetBusOverview(_ *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		status := r.URL.Query().Get("status")

		routes := filterRoutes(dummyBusRoutes, q, status)

		utils.RespondWithJSON(w, http.StatusOK, BusOverviewResponse{
			Routes:      routes,
			TotalRoutes: len(routes),
			GeneratedAt: time.Now().UTC(),
		})
	}
}

func GetBusTimetable(_ *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		status := r.URL.Query().Get("status")

		routes := filterRoutes(dummyBusRoutes, q, status)

		utils.RespondWithJSON(w, http.StatusOK, BusTimingResponse{
			Routes:    routes,
			Generated: time.Now().UTC(),
		})
	}
}

func GetBusMap(_ *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allStops := make([]BusStop, 0)
		for _, route := range dummyBusRoutes {
			allStops = append(allStops, route.Stops...)
		}

		utils.RespondWithJSON(w, http.StatusOK, BusMapResponse{
			CenterLat: 28.044,
			CenterLng: 76.107,
			Zoom:      10.5,
			Stops:     allStops,
			Routes:    dummyBusRoutes,
		})
	}
}

func GetBusRoute(_ *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		routeID := resolveBusID(r, "busid")
		if routeID == "" {
			routeID = strings.TrimPrefix(r.URL.Path, "/api/v1/bus/time/")
		}
		if routeID == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Route id is required")
			return
		}

		route, ok := findBusRouteByID(routeID)
		if !ok {
			utils.RespondWithError(w, http.StatusNotFound, "Bus route not found")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, route)
	}
}

func GetBusLiveTracking(_ *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		routeID := resolveBusID(r, "busid")
		if routeID == "" {
			routeID = strings.TrimPrefix(r.URL.Path, "/api/v1/bus/live/")
		}
		if routeID == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Route id is required")
			return
		}

		route, ok := findBusRouteByID(routeID)
		if !ok {
			utils.RespondWithError(w, http.StatusNotFound, "Bus route not found")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, BusLiveResponse{
			RouteID:   route.RouteID,
			Name:      route.Name,
			Vehicles:  route.Vehicles,
			UpdatedAt: time.Now().UTC(),
		})
	}
}

func GetBusDepot(_ *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		depotID := resolveBusID(r, "depotid")
		if depotID == "" {
			depotID = strings.TrimPrefix(r.URL.Path, "/api/v1/bus/depot/")
		}
		if depotID == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Depot id is required")
			return
		}

		depot, ok := findBusDepotByID(depotID)
		if !ok {
			utils.RespondWithError(w, http.StatusNotFound, "Bus depot not found")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, depot)
	}
}

func PostArrivalReport(_ *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ArrivalReportRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON request payload")
			return
		}

		if req.RouteID == "" || req.StopID == "" {
			utils.RespondWithError(w, http.StatusUnprocessableEntity, "route_id and stop_id are required")
			return
		}

		reportID := fmt.Sprintf("rep-%d", time.Now().UnixNano())

		utils.RespondWithJSON(w, http.StatusCreated, ArrivalReportResponse{
			ReportID:  reportID,
			Status:    "Report received. Thank you for contributing!",
			Timestamp: time.Now().UTC(),
		})
	}
}

func GetBusRouteByIDValue(routeID string) (BusRoute, bool) {
	if routeID == "" {
		return BusRoute{}, false
	}
	parsed, err := strconv.Atoi(routeID)
	if err == nil && parsed == 0 {
		return BusRoute{}, false
	}
	return findBusRouteByID(routeID)
}
