package admin

import (
	"net/http"

	"scav/infra"
	"scav/middleware"

	"github.com/julienschmidt/httprouter"
)

func AddAdminRoutes(router *httprouter.Router, app *infra.Deps, rateLimiter *middleware.RateLimiter) {
	authMid := middleware.Authenticate(app)
	adminOnly := middleware.Chain(authMid, middleware.RequireRoles("admin"))

	router.HandlerFunc(http.MethodPost, "/api/v1/admin/role/request", middleware.Chain(rateLimiter.Limit, authMid)(ApplyForRole(app)))
	router.HandlerFunc(http.MethodGet, "/api/v1/admin/role/requests/me", authMid(GetMyRoleRequests(app)))
	router.HandlerFunc(http.MethodGet, "/api/v1/admin/role/requests", adminOnly(ListRoleRequests(app)))
	router.HandlerFunc(http.MethodPut, "/api/v1/admin/role/requests/:id/approve", adminOnly(ApproveRoleRequest(app)))
	router.HandlerFunc(http.MethodPut, "/api/v1/admin/role/requests/:id/reject", adminOnly(RejectRoleRequest(app)))
}
