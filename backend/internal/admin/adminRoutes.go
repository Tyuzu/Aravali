package admin

import (
	"net/http"

	"scav/infra"
	"scav/internal/reports"
	"scav/middleware"

	"github.com/julienschmidt/httprouter"
)

func AddAdminRoutes(router *httprouter.Router, app *infra.Deps, rateLimiter *middleware.RateLimiter) {
	authMid := middleware.Authenticate(app)
	modOnly := middleware.Chain(authMid, middleware.RequireRoles("moderator"))
	adminOnly := middleware.Chain(authMid, middleware.RequireRoles("admin"))

	router.HandlerFunc(http.MethodPost, "/api/v1/admin/role/request", middleware.Chain(rateLimiter.Limit, authMid)(ApplyForRole(app)))
	router.HandlerFunc(http.MethodGet, "/api/v1/admin/role/requests/me", authMid(GetMyRoleRequests(app)))
	router.HandlerFunc(http.MethodGet, "/api/v1/admin/role/requests", adminOnly(ListRoleRequests(app)))
	router.HandlerFunc(http.MethodPut, "/api/v1/admin/role/requests/:id/approve", adminOnly(ApproveRoleRequest(app)))
	router.HandlerFunc(http.MethodPut, "/api/v1/admin/role/requests/:id/reject", adminOnly(RejectRoleRequest(app)))

	// Moderator application submission and review live with admin-related concerns.
	router.HandlerFunc(http.MethodPost, "/api/v1/moderator/apply", middleware.Chain(rateLimiter.Limit, authMid)(ApplyModerator(app)))
	router.HandlerFunc(http.MethodGet, "/api/v1/moderator/applications", adminOnly(ListModeratorApplications(app)))
	router.HandlerFunc(http.MethodPut, "/api/v1/moderator/approve/:id", adminOnly(ApproveModerator(app)))
	router.HandlerFunc(http.MethodPut, "/api/v1/moderator/reject/:id", adminOnly(RejectModerator(app)))

	// Moderator-level report review and content enforcement.
	router.HandlerFunc(http.MethodPut, "/api/v1/report/:id", modOnly(reports.UpdateReport(app)))
	router.HandlerFunc(http.MethodGet, "/api/v1/moderator/reports", modOnly(reports.GetReportsForMod(app)))
	router.HandlerFunc(http.MethodPut, "/api/v1/moderator/delete/:type/:id", modOnly(reports.SoftDeleteEntity(app)))
	router.HandlerFunc(http.MethodGet, "/api/v1/moderator/appeals", modOnly(reports.GetAppeals(app)))
	router.HandlerFunc(http.MethodPut, "/api/v1/moderator/appeals/:id", modOnly(reports.UpdateAppeal(app)))
}
