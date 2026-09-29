package sqldb

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotImplemented = errors.New("postgres database adapter: operation not implemented for the current relational schema")

// SortItem defines column sorting parameters for queries.
type SortItem struct {
	Column     string
	Descending bool
}

// PostgresDatabase provides a pure PostgreSQL driver implementation using pgxpool.
type PostgresDatabase struct {
	db      *pgxpool.Pool
	limiter chan struct{}
}

func NewPostgresDatabase(db *pgxpool.Pool, maxConcurrent int) *PostgresDatabase {
	if maxConcurrent <= 0 {
		maxConcurrent = 50
	}
	return &PostgresDatabase{db: db, limiter: make(chan struct{}, maxConcurrent)}
}

func (p *PostgresDatabase) Ping(ctx context.Context) error {
	return p.db.Ping(ctx)
}

func (p *PostgresDatabase) WithDB(ctx context.Context, op func(ctx context.Context) error) error {
	p.limiter <- struct{}{}
	defer func() { <-p.limiter }()

	for i := 0; i < 2; i++ {
		c, cancel := context.WithTimeout(ctx, 5*time.Second)
		err := op(c)
		cancel()
		if err == nil {
			return nil
		}
		if isRetryablePostgres(err) {
			time.Sleep(200 * time.Millisecond)
			continue
		}
		return err
	}
	return ErrNotImplemented
}

func (p *PostgresDatabase) RunTransaction(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := fn(tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

/* ---------------- Create ---------------- */

func (p *PostgresDatabase) Insert(ctx context.Context, table string, record any) error {
	return p.InsertOne(ctx, table, record)
}

func (p *PostgresDatabase) InsertOne(ctx context.Context, table string, record any) error {
	cols, vals, placeholders, err := extractColumnsAndValues(record)
	if err != nil {
		return err
	}
	if len(cols) == 0 {
		return fmt.Errorf("no columns found to insert into %s", table)
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		quoteIdent(table),
		strings.Join(cols, ", "),
		strings.Join(placeholders, ", "),
	)

	_, err = p.db.Exec(ctx, query, vals...)
	return err
}

func (p *PostgresDatabase) InsertMany(ctx context.Context, table string, records []any) error {
	if len(records) == 0 {
		return nil
	}
	return p.RunTransaction(ctx, func(tx pgx.Tx) error {
		for _, rec := range records {
			cols, vals, placeholders, err := extractColumnsAndValues(rec)
			if err != nil {
				return err
			}
			query := fmt.Sprintf(
				"INSERT INTO %s (%s) VALUES (%s)",
				quoteIdent(table),
				strings.Join(cols, ", "),
				strings.Join(placeholders, ", "),
			)
			if _, err := tx.Exec(ctx, query, vals...); err != nil {
				return err
			}
		}
		return nil
	})
}

func (p *PostgresDatabase) BulkWrite(ctx context.Context, table string, operations []any) error {
	return ErrNotImplemented
}

/* ---------------- Read ---------------- */

func (p *PostgresDatabase) FindOne(ctx context.Context, table string, whereClause string, args []any, result any) error {
	return p.FindOneWithProjection(ctx, table, nil, whereClause, args, result)
}

func (p *PostgresDatabase) FindOneWithProjection(ctx context.Context, table string, columns []string, whereClause string, args []any, result any) error {
	cols := "*"
	if len(columns) > 0 {
		cols = strings.Join(quoteIdents(columns), ", ")
	}

	where := "TRUE"
	if strings.TrimSpace(whereClause) != "" {
		where = whereClause
	}

	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s LIMIT 1", cols, quoteIdent(table), where)
	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		if rows.Err() != nil {
			return rows.Err()
		}
		return pgx.ErrNoRows
	}

	vals, err := rows.Values()
	if err != nil {
		return err
	}

	fds := rows.FieldDescriptions()
	colsNames := make([]string, len(fds))
	for i, fd := range fds {
		colsNames[i] = string(fd.Name)
	}

	return mapRowToDest(result, colsNames, vals)
}

func (p *PostgresDatabase) FindMany(ctx context.Context, table string, whereClause string, args []any, result any) error {
	return p.FindManyWithOptions(ctx, table, whereClause, args, FindManyOptions{}, result)
}

func (p *PostgresDatabase) FindManyWithOptions(ctx context.Context, table string, whereClause string, args []any, opts FindManyOptions, result any) error {
	return p.FindManyWithProjection(ctx, table, whereClause, args, opts.Columns, opts, result)
}

func (p *PostgresDatabase) FindManyWithProjection(
	ctx context.Context,
	table string,
	whereClause string,
	args []any,
	columns []string,
	opts FindManyOptions,
	result any,
) error {
	cols := "*"
	if len(columns) > 0 {
		cols = strings.Join(quoteIdents(columns), ", ")
	}

	where := "TRUE"
	if strings.TrimSpace(whereClause) != "" {
		where = whereClause
	}

	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s", cols, quoteIdent(table), where)

	if len(opts.Sort) > 0 {
		var sortParts []string
		for _, s := range opts.Sort {
			dir := "ASC"
			if s.Descending {
				dir = "DESC"
			}
			sortParts = append(sortParts, fmt.Sprintf("%s %s", quoteIdent(s.Column), dir))
		}
		query += " ORDER BY " + strings.Join(sortParts, ", ")
	}

	if opts.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", opts.Limit)
	}
	if opts.Offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", opts.Offset)
	}

	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	fds := rows.FieldDescriptions()
	colsNames := make([]string, len(fds))
	for i, fd := range fds {
		colsNames[i] = string(fd.Name)
	}

	rv := reflect.ValueOf(result)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("result argument must be a pointer to a slice, got %T", result)
	}
	sv := rv.Elem()
	elemType := sv.Type().Elem()

	for rows.Next() {
		vals, err := rows.Values()
		if err != nil {
			return err
		}

		newElem := reflect.New(elemType).Interface()
		if err := mapRowToDest(newElem, colsNames, vals); err != nil {
			return err
		}

		var toAppend reflect.Value
		if elemType.Kind() == reflect.Ptr {
			toAppend = reflect.ValueOf(newElem)
		} else {
			toAppend = reflect.ValueOf(newElem).Elem()
		}
		sv.Set(reflect.Append(sv, toAppend))
	}

	return rows.Err()
}

func (p *PostgresDatabase) Distinct(ctx context.Context, table string, column string, whereClause string, args []any, result any) error {
	where := "TRUE"
	if strings.TrimSpace(whereClause) != "" {
		where = whereClause
	}

	query := fmt.Sprintf("SELECT DISTINCT %s FROM %s WHERE %s", quoteIdent(column), quoteIdent(table), where)
	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	rv := reflect.ValueOf(result)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("result argument must be a pointer to a slice, got %T", result)
	}
	sv := rv.Elem()
	elemType := sv.Type().Elem()

	for rows.Next() {
		vals, err := rows.Values()
		if err != nil || len(vals) == 0 {
			continue
		}

		v := vals[0]
		var toAppend reflect.Value
		if v == nil {
			toAppend = reflect.Zero(elemType)
		} else {
			valrv := reflect.ValueOf(v)
			if valrv.Type().AssignableTo(elemType) {
				toAppend = valrv
			} else if valrv.Type().ConvertibleTo(elemType) {
				toAppend = valrv.Convert(elemType)
			} else {
				toAppend = reflect.Zero(elemType)
			}
		}
		sv.Set(reflect.Append(sv, toAppend))
	}

	return rows.Err()
}

/* ---------------- Update ---------------- */

func (p *PostgresDatabase) Update(ctx context.Context, table string, whereClause string, args []any, updateValues map[string]any) (int64, error) {
	return p.UpdateMany(ctx, table, whereClause, args, updateValues)
}

func (p *PostgresDatabase) UpdateOne(ctx context.Context, table string, whereClause string, args []any, updateValues map[string]any) (int64, error) {
	return p.UpdateMany(ctx, table, whereClause, args, updateValues)
}

func (p *PostgresDatabase) UpdateMany(ctx context.Context, table string, whereClause string, args []any, updateValues map[string]any) (int64, error) {
	if len(updateValues) == 0 {
		return 0, nil
	}

	setClauses := make([]string, 0, len(updateValues))
	queryArgs := make([]any, 0, len(args)+len(updateValues))

	argIdx := 1
	for col, val := range updateValues {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", quoteIdent(col), argIdx))
		queryArgs = append(queryArgs, val)
		argIdx++
	}

	// Safely re-index placeholders ($1, $2) inside WHERE clause without corrupting indexes > 9
	adjustedWhere := reindexPlaceholders(whereClause, argIdx-1)
	queryArgs = append(queryArgs, args...)

	where := "TRUE"
	if strings.TrimSpace(adjustedWhere) != "" {
		where = adjustedWhere
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", quoteIdent(table), strings.Join(setClauses, ", "), where)
	res, err := p.db.Exec(ctx, query, queryArgs...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}

func (p *PostgresDatabase) Upsert(ctx context.Context, table string, conflictColumn string, record any) error {
	cols, vals, placeholders, err := extractColumnsAndValues(record)
	if err != nil {
		return err
	}

	updates := make([]string, 0, len(cols))
	for _, col := range cols {
		if col == quoteIdent(conflictColumn) {
			continue
		}
		updates = append(updates, fmt.Sprintf("%s = EXCLUDED.%s", col, col))
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) ON CONFLICT (%s) DO UPDATE SET %s",
		quoteIdent(table),
		strings.Join(cols, ", "),
		strings.Join(placeholders, ", "),
		quoteIdent(conflictColumn),
		strings.Join(updates, ", "),
	)

	_, err = p.db.Exec(ctx, query, vals...)
	return err
}

func (p *PostgresDatabase) Inc(ctx context.Context, table string, whereClause string, args []any, column string, value int64) error {
	where := "TRUE"
	if strings.TrimSpace(whereClause) != "" {
		where = reindexPlaceholders(whereClause, 1)
	}

	query := fmt.Sprintf("UPDATE %s SET %s = %s + $1 WHERE %s", quoteIdent(table), quoteIdent(column), quoteIdent(column), where)
	queryArgs := append([]any{value}, args...)
	_, err := p.db.Exec(ctx, query, queryArgs...)
	return err
}

func (p *PostgresDatabase) AddToSet(ctx context.Context, table string, whereClause string, args []any, arrayColumn string, value any) error {
	where := "TRUE"
	if strings.TrimSpace(whereClause) != "" {
		where = reindexPlaceholders(whereClause, 1)
	}

	col := quoteIdent(arrayColumn)
	query := fmt.Sprintf(
		"UPDATE %s SET %s = CASE WHEN $1 = ANY(%s) THEN %s ELSE array_append(%s, $1) END WHERE %s",
		quoteIdent(table), col, col, col, col, where,
	)
	queryArgs := append([]any{value}, args...)
	_, err := p.db.Exec(ctx, query, queryArgs...)
	return err
}

/* ---------------- Delete ---------------- */

func (p *PostgresDatabase) Delete(ctx context.Context, table string, whereClause string, args []any) (int64, error) {
	return p.DeleteMany(ctx, table, whereClause, args)
}

func (p *PostgresDatabase) DeleteOne(ctx context.Context, table string, whereClause string, args []any) (int64, error) {
	return p.DeleteMany(ctx, table, whereClause, args)
}

func (p *PostgresDatabase) DeleteMany(ctx context.Context, table string, whereClause string, args []any) (int64, error) {
	where := "TRUE"
	if strings.TrimSpace(whereClause) != "" {
		where = whereClause
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE %s", quoteIdent(table), where)
	res, err := p.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}

/* ---------------- Atomic & Aggregations ---------------- */

func (p *PostgresDatabase) FindOneAndUpdate(ctx context.Context, table string, whereClause string, args []any, updateValues map[string]any, result any) error {
	return p.RunTransaction(ctx, func(tx pgx.Tx) error {
		rowCount, err := p.UpdateOne(ctx, table, whereClause, args, updateValues)
		if err != nil {
			return err
		}
		if rowCount == 0 {
			return pgx.ErrNoRows
		}
		return p.FindOne(ctx, table, whereClause, args, result)
	})
}

func (p *PostgresDatabase) QueryRaw(ctx context.Context, sqlQuery string, args []any, result any) error {
	rows, err := p.db.Query(ctx, sqlQuery, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	fds := rows.FieldDescriptions()
	colsNames := make([]string, len(fds))
	for i, fd := range fds {
		colsNames[i] = string(fd.Name)
	}

	rv := reflect.ValueOf(result)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("result must be a pointer, got %T", result)
	}
	ev := rv.Elem()

	if ev.Kind() == reflect.Slice {
		elemType := ev.Type().Elem()
		for rows.Next() {
			vals, err := rows.Values()
			if err != nil {
				return err
			}
			newElem := reflect.New(elemType).Interface()
			if err := mapRowToDest(newElem, colsNames, vals); err != nil {
				return err
			}
			var toAppend reflect.Value
			if elemType.Kind() == reflect.Ptr {
				toAppend = reflect.ValueOf(newElem)
			} else {
				toAppend = reflect.ValueOf(newElem).Elem()
			}
			ev.Set(reflect.Append(ev, toAppend))
		}
		return rows.Err()
	}

	if !rows.Next() {
		if rows.Err() != nil {
			return rows.Err()
		}
		return pgx.ErrNoRows
	}
	vals, err := rows.Values()
	if err != nil {
		return err
	}
	return mapRowToDest(result, colsNames, vals)
}

func (p *PostgresDatabase) Count(ctx context.Context, table string, whereClause string, args []any) (int64, error) {
	return p.CountDocuments(ctx, table, whereClause, args)
}

func (p *PostgresDatabase) CountDocuments(ctx context.Context, table string, whereClause string, args []any) (int64, error) {
	where := "TRUE"
	if strings.TrimSpace(whereClause) != "" {
		where = whereClause
	}

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", quoteIdent(table), where)
	var count int64
	err := p.db.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (p *PostgresDatabase) Aggregate(ctx context.Context, table string, whereClause string, args []any, out any) error {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", quoteIdent(table))
	if strings.TrimSpace(whereClause) != "" {
		query += " WHERE " + whereClause
	}
	return p.db.QueryRow(ctx, query, args...).Scan(out)
}

func (p *PostgresDatabase) EstimatedDocumentCount(ctx context.Context, table string) (int64, error) {
	var count int64
	query := "SELECT reltuples::bigint FROM pg_class WHERE relname = $1"
	err := p.db.QueryRow(ctx, query, table).Scan(&count)
	return count, err
}

/* ---------------- Helpers ---------------- */

func reindexPlaceholders(whereClause string, offset int) string {
	re := regexp.MustCompile(`\$(\d+)`)
	return re.ReplaceAllStringFunc(whereClause, func(match string) string {
		idx, _ := strconv.Atoi(match[1:])
		return fmt.Sprintf("$%d", idx+offset)
	})
}

func mapRowToDest(dest any, cols []string, vals []any) error {
	if dest == nil {
		return fmt.Errorf("nil destination")
	}
	rv := reflect.ValueOf(dest)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("destination must be a pointer, got %T", dest)
	}
	ev := rv.Elem()

	if ev.Kind() == reflect.Map {
		if ev.IsNil() {
			ev.Set(reflect.MakeMap(ev.Type()))
		}
		for i, c := range cols {
			ev.SetMapIndex(reflect.ValueOf(c), reflect.ValueOf(vals[i]))
		}
		return nil
	}

	if ev.Kind() != reflect.Struct {
		if len(vals) == 0 || vals[0] == nil {
			ev.Set(reflect.Zero(ev.Type()))
			return nil
		}
		valrv := reflect.ValueOf(vals[0])
		if valrv.Type().AssignableTo(ev.Type()) {
			ev.Set(valrv)
		} else if valrv.Type().ConvertibleTo(ev.Type()) {
			ev.Set(valrv.Convert(ev.Type()))
		}
		return nil
	}

	typ := ev.Type()
	for i, c := range cols {
		for j := 0; j < ev.NumField(); j++ {
			field := typ.Field(j)
			dbTag := field.Tag.Get("db")
			if dbTag == "" {
				dbTag = strings.ToLower(field.Name)
			}
			if dbTag == c || strings.EqualFold(field.Name, c) {
				fv := ev.Field(j)
				if !fv.CanSet() {
					break
				}
				v := vals[i]
				if v == nil {
					fv.Set(reflect.Zero(fv.Type()))
					break
				}
				valrv := reflect.ValueOf(v)
				if valrv.Type().AssignableTo(fv.Type()) {
					fv.Set(valrv)
				} else if valrv.Type().ConvertibleTo(fv.Type()) {
					fv.Set(valrv.Convert(fv.Type()))
				}
				break
			}
		}
	}
	return nil
}

func extractColumnsAndValues(record any) ([]string, []any, []string, error) {
	val := reflect.ValueOf(record)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil, nil, nil, fmt.Errorf("expected struct or struct pointer, got %T", record)
	}

	typ := val.Type()
	cols := make([]string, 0, val.NumField())
	vals := make([]any, 0, val.NumField())
	placeholders := make([]string, 0, val.NumField())

	argIdx := 1
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag == "-" {
			continue
		}
		if dbTag == "" {
			dbTag = strings.ToLower(field.Name)
		}

		cols = append(cols, quoteIdent(dbTag))
		vals = append(vals, val.Field(i).Interface())
		placeholders = append(placeholders, fmt.Sprintf("$%d", argIdx))
		argIdx++
	}

	return cols, vals, placeholders, nil
}

func quoteIdent(name string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9_]+`)
	return `"` + re.ReplaceAllString(name, "_") + `"`
}

func quoteIdents(names []string) []string {
	out := make([]string, len(names))
	for i, n := range names {
		out[i] = quoteIdent(n)
	}
	return out
}

func isRetryablePostgres(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	return strings.Contains(s, "connection refused") ||
		strings.Contains(s, "timeout") ||
		strings.Contains(s, "too many clients")
}
