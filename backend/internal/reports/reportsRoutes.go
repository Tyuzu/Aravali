package reports

import (
	"net/http"
	"scav/infra"
	"scav/middleware"

	"github.com/julienschmidt/httprouter"
)

// RegisterRoutes sets up HTTP routes for the reports package.

func AddReportingRoutes(router *httprouter.Router, app *infra.Deps, rateLimiter *middleware.RateLimiter) {
	authMid := middleware.Authenticate(app)

	// Submit a report
	router.HandlerFunc(
		http.MethodPost,
		"/api/v1/report",
		middleware.Chain(rateLimiter.Limit, authMid)(ReportContent(app)),
	)

	// Create an appeal
	router.HandlerFunc(
		http.MethodPost,
		"/api/v1/appeals",
		middleware.Chain(rateLimiter.Limit, authMid)(CreateAppeal(app)),
	)
}
