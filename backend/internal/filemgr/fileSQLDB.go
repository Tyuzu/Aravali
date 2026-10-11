// File: internal/filemgr/fileSQLDB.go

package filemgr

import (
	"scav/infra"
)

func updateEntityMediaInDB(app *infra.Deps, Table, idField, entityID string, update map[string]any) (int64, error) {
	return 0, nil
}
