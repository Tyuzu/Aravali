// File: internal/places/placedb/placeSQLDB.go

package placedb

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"scav/config"
	"scav/infra"

	"github.com/jackc/pgx/v5"
)

var placesTable = config.Tables.PlacesTable
var eventsTable = config.Tables.EventsTable
var productsTable = config.Tables.ProductTable
var membershipsTable = config.Tables.MembershipsTable

func buildUpdateSet(update map[string]any) (string, []any, error) {
	if len(update) == 0 {
		return "", nil, nil
	}
	parts := make([]string, 0, len(update))
	args := make([]any, 0, len(update))
	for key, value := range update {
		if key == "" {
			continue
		}
		parts = append(parts, fmt.Sprintf("%s = $%d", key, len(args)+1))
		args = append(args, value)
	}
	return strings.Join(parts, ", "), args, nil
}

func normalizeDBName(name string) string {
	return strings.ToLower(strings.NewReplacer("-", "", "_", "", " ", "").Replace(strings.TrimSpace(name)))
}

func columnFieldIndex(t reflect.Type, column string) (int, bool) {
	norm := normalizeDBName(column)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		for _, key := range []string{"db", "json"} {
			tag := field.Tag.Get(key)
			if tag == "" {
				continue
			}
			name := strings.Split(tag, ",")[0]
			if name != "" && normalizeDBName(name) == norm {
				return i, true
			}
		}
		if normalizeDBName(field.Name) == norm {
			return i, true
		}
	}
	return -1, false
}

func scanRowsIntoSlice(rows pgx.Rows, out any) error {
	if rows == nil {
		return nil
	}
	defer rows.Close()

	v := reflect.ValueOf(out)
	if v.Kind() != reflect.Ptr || v.IsNil() || v.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("out must be a non-nil pointer to a slice")
	}

	itemType := v.Elem().Type().Elem()
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return err
		}
		elem := reflect.New(itemType).Elem()
		for i, val := range values {
			colName := rows.FieldDescriptions()[i].Name
			idx, ok := columnFieldIndex(itemType, colName)
			if !ok || !elem.Field(idx).CanSet() || val == nil {
				continue
			}
			field := elem.Field(idx)
			rv := reflect.ValueOf(val)
			if !rv.Type().AssignableTo(field.Type()) {
				if rv.Type().ConvertibleTo(field.Type()) {
					rv = rv.Convert(field.Type())
				} else {
					continue
				}
			}
			field.Set(rv)
		}
		v.Elem().Set(reflect.Append(v.Elem(), elem))
	}
	return rows.Err()
}

func scanRowIntoStruct(rows pgx.Rows, out any) error {
	if rows == nil {
		return nil
	}
	v := reflect.ValueOf(out)
	if v.Kind() != reflect.Ptr || v.IsNil() || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("out must be a non-nil pointer to a struct")
	}
	if !rows.Next() {
		return pgx.ErrNoRows
	}
	defer rows.Close()

	values, err := rows.Values()
	if err != nil {
		return err
	}
	item := v.Elem()
	for i, val := range values {
		colName := rows.FieldDescriptions()[i].Name
		idx, ok := columnFieldIndex(item.Type(), colName)
		if !ok || !item.Field(idx).CanSet() || val == nil {
			continue
		}
		field := item.Field(idx)
		rv := reflect.ValueOf(val)
		if !rv.Type().AssignableTo(field.Type()) {
			if rv.Type().ConvertibleTo(field.Type()) {
				rv = rv.Convert(field.Type())
			} else {
				continue
			}
		}
		field.Set(rv)
	}
	return nil
}

// Places
func FindPlaces(ctx context.Context, app *infra.Deps, query string, args []any, out any) error {
	if app == nil || app.SQLDB == nil || out == nil {
		return nil
	}
	if query == "" {
		query = "1 = 1"
	}
	rows, err := app.SQLDB.Query(ctx, `SELECT * FROM `+placesTable+` WHERE `+query, args...)
	if err != nil {
		return err
	}
	return scanRowsIntoSlice(rows, out)
}

func FindOnePlace(ctx context.Context, app *infra.Deps, query string, args []any, out any) error {
	if app == nil || app.SQLDB == nil || out == nil {
		return nil
	}
	if query == "" {
		query = "1 = 1"
	}
	rows, err := app.SQLDB.Query(ctx, `SELECT * FROM `+placesTable+` WHERE `+query+` LIMIT 1`, args...)
	if err != nil {
		return err
	}
	return scanRowIntoStruct(rows, out)
}

func UpdatePlace(ctx context.Context, app *infra.Deps, query string, args []any, update map[string]any) (int64, error) {
	if app == nil || app.SQLDB == nil {
		return 0, nil
	}
	setClause, setArgs, err := buildUpdateSet(update)
	if err != nil || setClause == "" {
		return 0, err
	}
	values := append(setArgs, args...)
	result, err := app.SQLDB.Exec(ctx, `UPDATE `+placesTable+` SET `+setClause+` WHERE `+query, values...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func InsertPlace(ctx context.Context, app *infra.Deps, place any) error {
	if app == nil || app.SQLDB == nil || place == nil {
		return nil
	}
	_, err := app.SQLDB.Exec(ctx, `INSERT INTO `+placesTable+` DEFAULT VALUES`)
	return err
}

func DeletePlace(ctx context.Context, app *infra.Deps, query string, args []any) (int64, error) {
	if app == nil || app.SQLDB == nil {
		return 0, nil
	}
	if query == "" {
		query = "1 = 1"
	}
	result, err := app.SQLDB.Exec(ctx, `DELETE FROM `+placesTable+` WHERE `+query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

// Events
func CountEvents(ctx context.Context, app *infra.Deps, query string, args []any) (int64, error) {
	if app == nil || app.SQLDB == nil {
		return 0, nil
	}
	if query == "" {
		query = "1 = 1"
	}
	var count int64
	err := app.SQLDB.QueryRow(ctx, `SELECT COUNT(*) FROM `+eventsTable+` WHERE `+query, args...).Scan(&count)
	return count, err
}

func FindEventsWithOptions(ctx context.Context, app *infra.Deps, query string, args []any, opts map[string]any, out any) error {
	if app == nil || app.SQLDB == nil || out == nil {
		return nil
	}
	if query == "" {
		query = "1 = 1"
	}
	rows, err := app.SQLDB.Query(ctx, `SELECT * FROM `+eventsTable+` WHERE `+query, args...)
	if err != nil {
		return err
	}
	return scanRowsIntoSlice(rows, out)
}

// Place products (generic)
func FindPlaceProducts(ctx context.Context, app *infra.Deps, query string, args []any, out any) error {
	if app == nil || app.SQLDB == nil || out == nil {
		return nil
	}
	if query == "" {
		query = "1 = 1"
	}
	rows, err := app.SQLDB.Query(ctx, `SELECT * FROM `+productsTable+` WHERE `+query, args...)
	if err != nil {
		return err
	}
	return scanRowsIntoSlice(rows, out)
}

func InsertPlaceProduct(ctx context.Context, app *infra.Deps, product any) error {
	if app == nil || app.SQLDB == nil || product == nil {
		return nil
	}
	_, err := app.SQLDB.Exec(ctx, `INSERT INTO `+productsTable+` DEFAULT VALUES`)
	return err
}

func UpdatePlaceProduct(ctx context.Context, app *infra.Deps, query string, args []any, update map[string]any) (int64, error) {
	if app == nil || app.SQLDB == nil {
		return 0, nil
	}
	setClause, setArgs, err := buildUpdateSet(update)
	if err != nil || setClause == "" {
		return 0, err
	}
	values := append(setArgs, args...)
	result, err := app.SQLDB.Exec(ctx, `UPDATE `+productsTable+` SET `+setClause+` WHERE `+query, values...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func DeletePlaceProduct(ctx context.Context, app *infra.Deps, query string, args []any) (int64, error) {
	if app == nil || app.SQLDB == nil {
		return 0, nil
	}
	if query == "" {
		query = "1 = 1"
	}
	result, err := app.SQLDB.Exec(ctx, `DELETE FROM `+productsTable+` WHERE `+query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

// Memberships
func FindPlaceMemberships(ctx context.Context, app *infra.Deps, query string, args []any, out any) error {
	if app == nil || app.SQLDB == nil || out == nil {
		return nil
	}
	if query == "" {
		query = "1 = 1"
	}
	rows, err := app.SQLDB.Query(ctx, `SELECT * FROM `+membershipsTable+` WHERE `+query, args...)
	if err != nil {
		return err
	}
	return scanRowsIntoSlice(rows, out)
}

func FindOneMembership(ctx context.Context, app *infra.Deps, query string, args []any, out any) error {
	if app == nil || app.SQLDB == nil || out == nil {
		return nil
	}
	if query == "" {
		query = "1 = 1"
	}
	rows, err := app.SQLDB.Query(ctx, `SELECT * FROM `+membershipsTable+` WHERE `+query+` LIMIT 1`, args...)
	if err != nil {
		return err
	}
	return scanRowIntoStruct(rows, out)
}

func InsertMembership(ctx context.Context, app *infra.Deps, membership any) error {
	if app == nil || app.SQLDB == nil || membership == nil {
		return nil
	}
	_, err := app.SQLDB.Exec(ctx, `INSERT INTO `+membershipsTable+` DEFAULT VALUES`)
	return err
}

func UpdateMembership(ctx context.Context, app *infra.Deps, query string, args []any, update map[string]any) (int64, error) {
	if app == nil || app.SQLDB == nil {
		return 0, nil
	}
	setClause, setArgs, err := buildUpdateSet(update)
	if err != nil || setClause == "" {
		return 0, err
	}
	values := append(setArgs, args...)
	result, err := app.SQLDB.Exec(ctx, `UPDATE `+membershipsTable+` SET `+setClause+` WHERE `+query, values...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func DeleteMembership(ctx context.Context, app *infra.Deps, query string, args []any) (int64, error) {
	if app == nil || app.SQLDB == nil {
		return 0, nil
	}
	if query == "" {
		query = "1 = 1"
	}
	result, err := app.SQLDB.Exec(ctx, `DELETE FROM `+membershipsTable+` WHERE `+query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
