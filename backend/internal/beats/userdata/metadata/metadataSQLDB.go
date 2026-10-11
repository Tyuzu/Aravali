// File: internal/beats/userdata/metadata/metadataSQLDB.go

package metadata

import (
	"context"

	"scav/config"
	"scav/infra"
)

var usersTable = config.Tables.UserTable

// SQLFindUsersByIDs returns minimal user docs for given ids
func FindUsersByIDs(ctx context.Context, app *infra.Deps, ids []string, out any) error {
	if len(ids) == 0 {
		return nil
	}

}
