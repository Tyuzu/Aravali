package filemgr

import (
	"context"
	"fmt"

	"scav/infra"
)

func SQLupdateEntityMediaInDB(app *infra.Deps, collection, idField, entityID string, update map[string]any) (int64, error) {
	query := fmt.Sprintf("%s = $1", idField)
	args := []any{entityID}

	return app.SQLDB.UpdateOne(context.Background(), collection, query, args, update)
}
