package analytics

import (
	"net/http"
	"scav/utils"
)

// RegisterRoutes sets up HTTP routes for the analytics package.

// --- Delegator Handler ---
func GetEntityAnalytics(w http.ResponseWriter, r *http.Request) {
	entityType := utils.GetParam(r, "entityType")
	entityID := utils.GetParam(r, "entityId")

	var analytics Analytics

	switch entityType {
	case "events":
		analytics = getEventAnalytics(entityID)
	case "places":
		analytics = getPlaceAnalytics(entityID)
	case "products":
		analytics = getProductAnalytics(entityID)
	default:
		http.Error(w, "Invalid entity type", http.StatusBadRequest)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, analytics)
}
