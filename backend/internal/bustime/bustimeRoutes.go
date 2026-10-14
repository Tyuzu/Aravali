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
	applyLimit := func(next http.HandlerFunc) http.HandlerFunc {
		if rateLimiter == nil {
			return next
		}
		return rateLimiter.Limit(next)
	}

	router.HandlerFunc(http.MethodGet, "/api/v1/bus", applyLimit(GetBusOverview(app)))
	router.HandlerFunc(http.MethodGet, "/api/v1/bus/overview", applyLimit(GetBusOverview(app)))
	router.HandlerFunc(http.MethodGet, "/api/v1/bustime", applyLimit(GetBusOverview(app)))

	router.HandlerFunc(http.MethodGet, "/api/v1/bus/time", applyLimit(GetBusTimetable(app)))
	router.HandlerFunc(http.MethodGet, "/api/v1/bus/map", applyLimit(GetBusMap(app)))
	router.HandlerFunc(http.MethodGet, "/api/v1/bus/time/:busid", applyLimit(GetBusRoute(app)))
	router.HandlerFunc(http.MethodGet, "/api/v1/bus/live/:busid", applyLimit(GetBusLiveTracking(app)))
	router.HandlerFunc(http.MethodGet, "/api/v1/bus/depot/:depotid", applyLimit(GetBusDepot(app)))

	router.HandlerFunc(http.MethodPost, "/api/v1/bus/report", applyLimit(PostArrivalReport(app)))

	router.HandlerFunc(http.MethodGet, "/api/v1/bustime/:busid", applyLimit(GetBusRoute(app)))
}
