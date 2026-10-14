// File: internal/bustime/bustimeRoutes.go

package bustime

import (
	"net/http"
	"scav/infra"
	"scav/middleware"

	"github.com/julienschmidt/httprouter"
)

// AddBustimeRoutes registers HTTP endpoints for the bustime module.
func AddBustimeRoutes(router *httprouter.Router, app *infra.Deps, rateLimiter *middleware.RateLimiter) {
	public := func(next http.HandlerFunc) http.HandlerFunc {
		if rateLimiter == nil {
			return next
		}
		return rateLimiter.Limit(next)
	}

	// These routes are intentionally public so guests can browse bus status,
	// route details, maps, and timetables without requiring sign-in.
	router.HandlerFunc(http.MethodGet, "/api/v1/bus", public(GetBusOverview(app)))
	router.HandlerFunc(http.MethodGet, "/api/v1/bus/overview", public(GetBusOverview(app)))
	router.HandlerFunc(http.MethodGet, "/api/v1/bustime", public(GetBusOverview(app)))

	router.HandlerFunc(http.MethodGet, "/api/v1/bus/time", public(GetBusTimetable(app)))
	router.HandlerFunc(http.MethodGet, "/api/v1/bus/map", public(GetBusMap(app)))
	router.HandlerFunc(http.MethodGet, "/api/v1/bus/time/:busid", public(GetBusRoute(app)))
	router.HandlerFunc(http.MethodGet, "/api/v1/bus/live/:busid", public(GetBusLiveTracking(app)))
	router.HandlerFunc(http.MethodGet, "/api/v1/bus/depot/:depotid", public(GetBusDepot(app)))

	router.HandlerFunc(http.MethodPost, "/api/v1/bus/report", public(PostArrivalReport(app)))

	router.HandlerFunc(http.MethodGet, "/api/v1/bustime/:busid", public(GetBusRoute(app)))
}
