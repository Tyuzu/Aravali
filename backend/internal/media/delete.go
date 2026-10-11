// File: internal/media/delete.go

package media

import (
	"net/http"
	"scav/infra"
)

// ---------------------- Delete Media ----------------------
func DeleteMedia(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
