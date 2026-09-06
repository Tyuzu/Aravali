package metadata

import (
	"context"
	"scav/config"
	"scav/infra"
)

var usersCollection = config.Collections.UserCollection

// FindUsersByIDs returns minimal user docs for given ids
func FindUsersByIDs(ctx context.Context, app *infra.Deps, ids []string, out interface{}) error {
	filter := map[string]any{"userid": map[string]any{"$in": ids}}
	return app.DB.FindMany(ctx, usersCollection, filter, out)
}
