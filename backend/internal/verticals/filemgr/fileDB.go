package filemgr

import (
	"context"
	"scav/infra"
)

func updateEntityMediaInDB(app *infra.Deps, collection, idField, entityID string, update map[string]any) (any, error) {
	return app.DB.UpdateOne(context.Background(), collection, map[string]any{idField: entityID}, update)
}
