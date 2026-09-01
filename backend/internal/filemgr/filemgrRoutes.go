package filemgr

import (
	"net/http"
	"scav/infra"
	"scav/middleware"

	"github.com/julienschmidt/httprouter"
)

// RegisterRoutes sets up HTTP routes for the filemgr package.

func AddFiledropRoutes(
	router *httprouter.Router,
	app *infra.Deps,
	rateLimiter *middleware.RateLimiter,
) {
	authMid := middleware.Authenticate(app)

	// Combine middleware: Rate Limit -> Auth
	filedropChain := middleware.Chain(
		rateLimiter.Limit,
		authMid,
	)

	// Main filedrop upload route
	router.HandlerFunc(
		http.MethodPost,
		"/api/v1/filedrop",
		filedropChain(FiledropHandler(app)),
	)

	// CORS Preflight handler
	router.HandlerFunc(
		http.MethodOptions,
		"/api/v1/filedrop",
		OptionsHandler,
	)
	// Proxy handler for external media supporting both query parameter & wildcard paths
	router.HandlerFunc(http.MethodGet, "/static/proxy/*url", ProxyHandler)
	router.HandlerFunc(http.MethodGet, "/static/proxy", ProxyHandler)
}
