// File: routes/routesWrapper.go

package routes

import (
	"scav/infra"
	"scav/internal/auth"
	"scav/internal/profile"
	"scav/middleware"

	"github.com/julienschmidt/httprouter"
)

func RoutesWrapper(router *httprouter.Router, app *infra.Deps, rateLimiter *middleware.RateLimiter) {
	auth.AddAuthRoutes(router, app, rateLimiter)
	profile.AddProfileRoutes(router, app, rateLimiter)

	AddUtilityRoutes(router, app, rateLimiter)
}
