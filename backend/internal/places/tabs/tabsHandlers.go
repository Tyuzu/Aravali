// File: internal/places/tabs/tabsHandlers.go

package places

import (
	"scav/infra"

	"github.com/julienschmidt/httprouter"
)

// RegisterRoutes sets up HTTP routes for the tabs package.
// The parent places package owns route registration, but this keeps the tabs package
// self-documenting and available for direct binding when needed.
func RegisterRoutes(router *httprouter.Router, app *infra.Deps) {
	if router == nil {
		return
	}
	_ = app
}
