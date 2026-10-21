// File: routes/routes.go

package routes

import (
	"net/http"
	"scav/infra"
	"scav/middleware"
	"scav/utils"

	"github.com/julienschmidt/httprouter"
)

// func AddStaticRoutes(router *httprouter.Router) {
// 	router.ServeFiles("/static/uploads/*filepath", http.Dir("static/uploads"))
// }

// Routes registration

func AddUtilityRoutes(router *httprouter.Router, app *infra.Deps, rateLimiter *middleware.RateLimiter) {
	authmidware := middleware.Authenticate(app)
	router.HandlerFunc(http.MethodGet, "/api/v1/csrf", rateLimiter.Limit(authmidware(utils.CSRF)))
}
