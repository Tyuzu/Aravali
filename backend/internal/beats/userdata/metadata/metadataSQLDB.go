package metadata

import (
	"context"
	"scav/config"
	"scav/infra"
)

var usersTable = config.Tables.UserTable

// FindUsersByIDs returns minimal user docs for given ids
func SQLFindUsersByIDs(ctx context.Context, app *infra.Deps, ids []string, out interface{}) error {
	filter := map[string]any{"userid": map[string]any{"$in": ids}}
	return app.DB.FindMany(ctx, usersTable, filter, out)
}
