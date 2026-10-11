// File: internal/admin/adminSQLDB.go

package admin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"scav/infra"

	"github.com/jackc/pgx/v5"
)

func ensureAdminDB(app *infra.Deps) error {
	if app == nil || app.SQLDB == nil {
		return errors.New("database not initialized")
	}
	return nil
}

func rowHasColumn(ctx context.Context, app *infra.Deps, tableName, columnName string) (bool, error) {
	if err := ensureAdminDB(app); err != nil {
		return false, err
	}

	query := `SELECT EXISTS (
		SELECT 1
		FROM information_schema.columns
		WHERE table_schema = current_schema()
		  AND table_name = $1
		  AND column_name = $2
	)`
	var exists bool
	if err := app.SQLDB.QueryRow(ctx, query, strings.TrimSpace(tableName), strings.TrimSpace(columnName)).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func scanRoleApplicationRow(row pgx.Row, dst *RoleApplication) error {
	if dst == nil {
		return errors.New("nil result")
	}

	var metadata []byte
	var createdAt, updatedAt time.Time
	if err := row.Scan(&dst.ID, &dst.UserID, &dst.Role, &dst.Status, &createdAt, &updatedAt, &metadata); err != nil {
		return err
	}
	if dst.Role == "" {
		dst.Role = ""
	}
	dst.CreatedAt = createdAt
	dst.UpdatedAt = updatedAt
	dst.Metadata = metadata
	return nil
}

func scanModeratorApplicationRow(row pgx.Row, dst *ModeratorApplication) error {
	if dst == nil {
		return errors.New("nil result")
	}

	var reason string
	var createdAt, updatedAt time.Time
	if err := row.Scan(&dst.ID, &dst.UserID, &reason, &dst.Status, &createdAt, &updatedAt); err != nil {
		return err
	}
	dst.Reason = reason
	dst.CreatedAt = createdAt
	dst.UpdatedAt = updatedAt
	return nil
}

// Role applications
func FindPendingRoleApplication(ctx context.Context, app *infra.Deps, userID, role string, result *RoleApplication) error {
	if err := ensureAdminDB(app); err != nil {
		return err
	}
	query := `SELECT roleapplicationid, userid, role, status, created_at, updated_at, metadata
		FROM role_applications
		WHERE userid = $1 AND role = $2 AND status = $3
		ORDER BY created_at DESC
		LIMIT 1`
	row := app.SQLDB.QueryRow(ctx, query, userID, role, "pending")
	return scanRoleApplicationRow(row, result)
}

func InsertRoleApplication(ctx context.Context, app *infra.Deps, application RoleApplication) error {
	if err := ensureAdminDB(app); err != nil {
		return err
	}

	if application.Metadata == nil {
		metadata, err := json.Marshal(map[string]any{"reason": application.Reason})
		if err != nil {
			return err
		}
		application.Metadata = metadata
	}

	query := `INSERT INTO role_applications (
		roleapplicationid, userid, role, status, created_at, updated_at, metadata
	) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := app.SQLDB.Exec(ctx, query,
		application.ID,
		application.UserID,
		application.Role,
		application.Status,
		application.CreatedAt,
		application.UpdatedAt,
		application.Metadata,
	)
	return err
}

func FindRoleApplicationsByUser(ctx context.Context, app *infra.Deps, userID string, result *[]RoleApplication) error {
	if err := ensureAdminDB(app); err != nil {
		return err
	}
	query := `SELECT roleapplicationid, userid, role, status, created_at, updated_at, metadata
		FROM role_applications
		WHERE userid = $1
		ORDER BY created_at DESC`
	rows, err := app.SQLDB.Query(ctx, query, userID)
	if err != nil {
		return err
	}
	defer rows.Close()

	items := make([]RoleApplication, 0)
	for rows.Next() {
		item := RoleApplication{}
		if err := scanRoleApplicationRow(rows, &item); err != nil {
			return err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*result = items
	return nil
}

func ListRoleApplicationsDB(ctx context.Context, app *infra.Deps, query string, args []any, result *[]RoleApplication) error {
	if err := ensureAdminDB(app); err != nil {
		return err
	}
	if query == "" {
		query = "1 = 1"
	}

	statement := `SELECT roleapplicationid, userid, role, status, created_at, updated_at, metadata
		FROM role_applications
		WHERE ` + query + ` ORDER BY created_at DESC`
	rows, err := app.SQLDB.Query(ctx, statement, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	items := make([]RoleApplication, 0)
	for rows.Next() {
		item := RoleApplication{}
		if err := scanRoleApplicationRow(rows, &item); err != nil {
			return err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*result = items
	return nil
}

func GetRoleApplicationByID(ctx context.Context, app *infra.Deps, id string, result *RoleApplication) error {
	if err := ensureAdminDB(app); err != nil {
		return err
	}
	query := `SELECT roleapplicationid, userid, role, status, created_at, updated_at, metadata
		FROM role_applications
		WHERE roleapplicationid = $1`
	row := app.SQLDB.QueryRow(ctx, query, id)
	return scanRoleApplicationRow(row, result)
}

func GetUserRoles(ctx context.Context, app *infra.Deps, userID string, result any) error {
	if err := ensureAdminDB(app); err != nil {
		return err
	}
	if result == nil {
		return errors.New("nil result")
	}

	query := `SELECT role FROM users WHERE userid = $1`
	var roles []string
	if err := app.SQLDB.QueryRow(ctx, query, userID).Scan(&roles); err != nil {
		return err
	}

	value := reflect.ValueOf(result)
	if value.Kind() != reflect.Ptr || value.IsNil() {
		return errors.New("result must be a non-nil pointer")
	}

	elem := value.Elem()
	if elem.Kind() == reflect.Struct {
		field := elem.FieldByName("Role")
		if field.IsValid() && field.CanSet() && field.Type() == reflect.TypeOf([]string{}) {
			field.Set(reflect.ValueOf(roles))
			return nil
		}
	}
	if elem.Kind() == reflect.Map {
		elem.SetMapIndex(reflect.ValueOf("role"), reflect.ValueOf(roles))
		return nil
	}
	return fmt.Errorf("unsupported role result type %T", result)
}

func UpdateUserRoles(ctx context.Context, app *infra.Deps, userID string, roles []string) (int64, error) {
	if err := ensureAdminDB(app); err != nil {
		return 0, err
	}
	query := `UPDATE users SET role = $1, updated_at = NOW() WHERE userid = $2`
	tag, err := app.SQLDB.Exec(ctx, query, roles, userID)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

func UpdateRoleApplicationStatus(ctx context.Context, app *infra.Deps, appID, status string) (int64, error) {
	if err := ensureAdminDB(app); err != nil {
		return 0, err
	}
	query := `UPDATE role_applications SET status = $1, updated_at = NOW() WHERE roleapplicationid = $2`
	tag, err := app.SQLDB.Exec(ctx, query, status, appID)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

// Moderator applications
func FindModeratorApplicationByUser(ctx context.Context, app *infra.Deps, userID string, result *ModeratorApplication) error {
	if err := ensureAdminDB(app); err != nil {
		return err
	}

	hasReasonCol, err := rowHasColumn(ctx, app, moderatorApplicationsTable, "reason")
	if err != nil {
		return err
	}

	selectExpr := "reason"
	if !hasReasonCol {
		selectExpr = "metadata->>'reason'"
	}

	query := fmt.Sprintf(`SELECT appid, userid, %s, status, created_at, updated_at
		FROM %s
		WHERE userid = $1
		ORDER BY created_at DESC
		LIMIT 1`, selectExpr, moderatorApplicationsTable)
	row := app.SQLDB.QueryRow(ctx, query, userID)
	return scanModeratorApplicationRow(row, result)
}

func InsertModeratorApplication(ctx context.Context, app *infra.Deps, application ModeratorApplication) error {
	if err := ensureAdminDB(app); err != nil {
		return err
	}

	hasReasonCol, err := rowHasColumn(ctx, app, moderatorApplicationsTable, "reason")
	if err != nil {
		return err
	}

	metadata := map[string]any{"reason": application.Reason}
	payload, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	if hasReasonCol {
		query := `INSERT INTO modapps (appid, userid, reason, status, created_at, updated_at, metadata)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`
		_, err = app.SQLDB.Exec(ctx, query, application.ID, application.UserID, application.Reason, application.Status, application.CreatedAt, application.UpdatedAt, payload)
		return err
	}

	query := `INSERT INTO modapps (appid, userid, status, created_at, updated_at, metadata)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = app.SQLDB.Exec(ctx, query, application.ID, application.UserID, application.Status, application.CreatedAt, application.UpdatedAt, payload)
	return err
}

func ListModeratorApplicationsDB(ctx context.Context, app *infra.Deps, query string, args []any, result *[]ModeratorApplication) error {
	if err := ensureAdminDB(app); err != nil {
		return err
	}
	if query == "" {
		query = "1 = 1"
	}

	hasReasonCol, err := rowHasColumn(ctx, app, moderatorApplicationsTable, "reason")
	if err != nil {
		return err
	}
	selectExpr := "reason"
	if !hasReasonCol {
		selectExpr = "metadata->>'reason'"
	}

	statement := fmt.Sprintf(`SELECT appid, userid, %s, status, created_at, updated_at
		FROM %s
		WHERE %s
		ORDER BY created_at DESC`, selectExpr, moderatorApplicationsTable, query)
	rows, err := app.SQLDB.Query(ctx, statement, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	items := make([]ModeratorApplication, 0)
	for rows.Next() {
		item := ModeratorApplication{}
		if err := scanModeratorApplicationRow(rows, &item); err != nil {
			return err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*result = items
	return nil
}

func UpdateModeratorApplicationStatus(ctx context.Context, app *infra.Deps, id, status string) (any, error) {
	if err := ensureAdminDB(app); err != nil {
		return nil, err
	}
	query := `UPDATE modapps SET status = $1, updated_at = NOW() WHERE appid = $2`
	tag, err := app.SQLDB.Exec(ctx, query, status, id)
	if err != nil {
		return nil, err
	}
	return tag.RowsAffected(), nil
}
