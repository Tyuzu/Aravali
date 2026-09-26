package metadata

import (
	"context"

	"scav/config"
	"scav/infra"
)

var usersTable = config.Tables.UserTable

// SQLFindUsersByIDs returns minimal user docs for given ids
func SQLFindUsersByIDs(ctx context.Context, app *infra.Deps, ids []string, out any) error {
	if len(ids) == 0 {
		return nil
	}

	query := "userid = ANY($1)"
	args := []any{ids}

	return app.SQLDB.FindMany(ctx, usersTable, query, args, out)
}
