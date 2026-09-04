package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var errNotImplemented = errors.New("postgres database adapter: operation not implemented for the current document schema")

// PostgresDatabase is a lightweight compatibility layer that stores each Mongo-like
// collection as a JSONB table so the application can run against PostgreSQL without
// rewriting every repository. It is intentionally conservative and focused on the
// interface used across the codebase.
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
	return errNotImplemented
}

func (p *PostgresDatabase) RunTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := fn(ctx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (p *PostgresDatabase) Insert(ctx context.Context, collection string, document any) error {
	return p.InsertOne(ctx, collection, document)
}

func (p *PostgresDatabase) InsertOne(ctx context.Context, collection string, document any) error {
	if err := p.ensureTable(ctx, collection); err != nil {
		return err
	}
	jsonDoc, err := marshalDocument(document)
	if err != nil {
		return err
	}
	id := documentID(document)
	_, err = p.db.Exec(ctx, "INSERT INTO "+quoteIdent(collection)+" (id, payload) VALUES ($1, $2::jsonb) ON CONFLICT (id) DO UPDATE SET payload = EXCLUDED.payload",
		id, string(jsonDoc),
	)
	return err
}

func (p *PostgresDatabase) InsertMany(ctx context.Context, collection string, documents []any) error {
	if len(documents) == 0 {
		return nil
	}
	for _, doc := range documents {
		if err := p.InsertOne(ctx, collection, doc); err != nil {
			return err
		}
	}
	return nil
}

func (p *PostgresDatabase) BulkWrite(ctx context.Context, collection string, operations []any) error {
	return errNotImplemented
}

func (p *PostgresDatabase) FindOne(ctx context.Context, collection string, filter any, result any) error {
	return p.findOneBase(ctx, collection, filter, nil, result)
}

func (p *PostgresDatabase) FindOneWithProjection(ctx context.Context, collection string, filter any, projection []string, result any) error {
	return p.findOneBase(ctx, collection, filter, projection, result)
}

func buildFilterQuery(filter any) (string, []any, error) {
	where, args, err := buildFilterPredicate(filter)
	if err != nil {
		return "", nil, err
	}
	if where == "" {
		where = "TRUE"
	}
	return where, args, nil
}

func buildFilterPredicate(filter any) (string, []any, error) {
	if filter == nil {
		return "TRUE", nil, nil
	}

	switch v := filter.(type) {
	case bson.M:
		return buildMapPredicate(v)
	case map[string]any:
		return buildMapPredicate(v)
	case []any:
		return buildListPredicate(v, "OR")
	default:
		jsonFilter, err := toJSONFilter(filter)
		if err != nil {
			return "", nil, err
		}
		return "payload @> $1::jsonb", []any{string(jsonFilter)}, nil
	}
}

func buildMapPredicate(m map[string]any) (string, []any, error) {
	clauses := make([]string, 0, len(m))
	args := make([]any, 0, len(m))

	for key, value := range m {
		switch strings.ToLower(key) {
		case "$and":
			clause, clauseArgs, err := buildListPredicate(value, "AND")
			if err != nil {
				return "", nil, err
			}
			if clause != "" && clause != "TRUE" {
				clauses = append(clauses, clause)
				args = append(args, clauseArgs...)
			}
		case "$or":
			clause, clauseArgs, err := buildListPredicate(value, "OR")
			if err != nil {
				return "", nil, err
			}
			if clause != "" && clause != "TRUE" {
				clauses = append(clauses, clause)
				args = append(args, clauseArgs...)
			}
		default:
			clause, clauseArgs, err := buildFieldPredicate(key, value)
			if err != nil {
				return "", nil, err
			}
			if clause != "" && clause != "TRUE" {
				clauses = append(clauses, clause)
				args = append(args, clauseArgs...)
			}
		}
	}
	if len(clauses) == 0 {
		return "TRUE", nil, nil
	}
	return strings.Join(clauses, " AND "), args, nil
}

func buildListPredicate(items any, op string) (string, []any, error) {
	list, ok := items.([]any)
	if !ok {
		return "TRUE", nil, nil
	}
	if len(list) == 0 {
		return "TRUE", nil, nil
	}

	parts := make([]string, 0, len(list))
	args := make([]any, 0, len(list))
	for _, item := range list {
		clause, clauseArgs, err := buildFilterPredicate(item)
		if err != nil {
			return "", nil, err
		}
		if clause == "" || clause == "TRUE" {
			continue
		}
		parts = append(parts, "("+clause+")")
		args = append(args, clauseArgs...)
	}
	if len(parts) == 0 {
		return "TRUE", nil, nil
	}
	return strings.Join(parts, " "+op+" "), args, nil
}

func buildFieldPredicate(field string, value any) (string, []any, error) {
	if field == "" {
		return "TRUE", nil, nil
	}

	switch v := value.(type) {
	case bson.M:
		return buildOperatorPredicate(field, v)
	case map[string]any:
		return buildOperatorPredicate(field, v)
	default:
		jsonFilter, err := json.Marshal(map[string]any{field: value})
		if err != nil {
			return "", nil, err
		}
		return "payload @> $1::jsonb", []any{string(jsonFilter)}, nil
	}
}

func buildOperatorPredicate(field string, ops map[string]any) (string, []any, error) {
	if len(ops) == 0 {
		return "TRUE", nil, nil
	}
	if len(ops) > 1 {
		jsonFilter, err := json.Marshal(map[string]any{field: ops})
		if err != nil {
			return "", nil, err
		}
		return "payload @> $1::jsonb", []any{string(jsonFilter)}, nil
	}

	for op, operand := range ops {
		switch strings.ToLower(op) {
		case "$in":
			values, err := asStringSlice(operand)
			if err != nil {
				return "", nil, err
			}
			return "COALESCE(" + jsonTextPath(field) + ", '') = ANY($1)", []any{values}, nil
		case "$nin":
			values, err := asStringSlice(operand)
			if err != nil {
				return "", nil, err
			}
			return "NOT (COALESCE(" + jsonTextPath(field) + ", '') = ANY($1))", []any{values}, nil
		case "$eq":
			jsonFilter, err := json.Marshal(map[string]any{field: operand})
			if err != nil {
				return "", nil, err
			}
			return "payload @> $1::jsonb", []any{string(jsonFilter)}, nil
		case "$ne":
			jsonFilter, err := json.Marshal(map[string]any{field: operand})
			if err != nil {
				return "", nil, err
			}
			return "NOT (payload @> $1::jsonb)", []any{string(jsonFilter)}, nil
		default:
			jsonFilter, err := json.Marshal(map[string]any{field: map[string]any{op: operand}})
			if err != nil {
				return "", nil, err
			}
			return "payload @> $1::jsonb", []any{string(jsonFilter)}, nil
		}
	}
	return "TRUE", nil, nil
}

func jsonTextPath(field string) string {
	parts := strings.Split(field, ".")
	quoted := make([]string, 0, len(parts))
	for _, part := range parts {
		quoted = append(quoted, "'"+strings.ReplaceAll(part, "'", "''")+"'")
	}
	if len(quoted) == 0 {
		return "''"
	}
	return "payload #>> ARRAY[" + strings.Join(quoted, ",") + "]"
}

func asStringSlice(v any) ([]string, error) {
	switch value := v.(type) {
	case []string:
		return value, nil
	case []any:
		out := make([]string, 0, len(value))
		for _, item := range value {
			out = append(out, fmt.Sprintf("%v", item))
		}
		return out, nil
	default:
		return nil, fmt.Errorf("unsupported $in value type %T", v)
	}
}

func (p *PostgresDatabase) FindMany(
	ctx context.Context,
	collection string,
	filter any,
	result any,
	opts ...*options.FindOptions,
) error {
	return p.findManyBase(ctx, collection, filter, nil, result)
}

func (p *PostgresDatabase) FindManyWithOptions(
	ctx context.Context,
	collection string,
	filter any,
	opts FindManyOptions,
	result any,
) error {
	return p.findManyBase(ctx, collection, filter, nil, result)
}

func (p *PostgresDatabase) FindManyWithProjection(
	ctx context.Context,
	collection string,
	filter any,
	projection []string,
	opts FindManyOptions,
	result any,
) error {
	return p.findManyBase(ctx, collection, filter, projection, result)
}

func (p *PostgresDatabase) Distinct(ctx context.Context, collection string, field string, filter any, result any) error {
	return errNotImplemented
}

func (p *PostgresDatabase) Update(ctx context.Context, collection string, filter any, update any) (any, error) {
	return p.UpdateOne(ctx, collection, filter, update)
}

func (p *PostgresDatabase) UpdateOne(ctx context.Context, collection string, filter any, update any) (any, error) {
	if err := p.ensureTable(ctx, collection); err != nil {
		return nil, err
	}
	merged, err := p.mergeUpdate(ctx, collection, filter, update)
	if err != nil {
		return nil, err
	}
	return merged, nil
}

func (p *PostgresDatabase) UpdateMany(ctx context.Context, collection string, filter any, update any) (any, error) {
	return p.UpdateOne(ctx, collection, filter, update)
}

func (p *PostgresDatabase) Upsert(ctx context.Context, collection string, filter any, document any) error {
	if err := p.ensureTable(ctx, collection); err != nil {
		return err
	}
	if count, err := p.CountDocuments(ctx, collection, filter); err == nil && count == 0 {
		return p.InsertOne(ctx, collection, document)
	}
	_, err := p.UpdateOne(ctx, collection, filter, document)
	return err
}

func (p *PostgresDatabase) Inc(ctx context.Context, collection string, filter any, field string, value int64) error {
	return errNotImplemented
}

func (p *PostgresDatabase) AddToSet(ctx context.Context, collection string, filter any, field string, value any) error {
	return errNotImplemented
}

func (p *PostgresDatabase) Delete(ctx context.Context, collection string, filter any) (int64, error) {
	return p.DeleteOne(ctx, collection, filter)
}

func (p *PostgresDatabase) DeleteOne(ctx context.Context, collection string, filter any) (int64, error) {
	if err := p.ensureTable(ctx, collection); err != nil {
		return 0, err
	}
	where, args, err := buildFilterQuery(filter)
	if err != nil {
		return 0, err
	}
	cmd := "DELETE FROM " + quoteIdent(collection) + " WHERE " + where
	res, err := p.db.Exec(ctx, cmd, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}

func (p *PostgresDatabase) DeleteMany(ctx context.Context, collection string, filter any) error {
	_, err := p.DeleteOne(ctx, collection, filter)
	return err
}

func (p *PostgresDatabase) FindOneAndUpdate(ctx context.Context, collection string, filter any, update any, result any) error {
	if err := p.ensureTable(ctx, collection); err != nil {
		return err
	}
	if err := p.FindOne(ctx, collection, filter, result); err != nil {
		return err
	}
	_, err := p.UpdateOne(ctx, collection, filter, update)
	return err
}

func (p *PostgresDatabase) Aggregate(ctx context.Context, collection string, pipeline any, result any) error {
	return errNotImplemented
}

func (p *PostgresDatabase) Count(ctx context.Context, collection string, filter any) (int64, error) {
	return p.CountDocuments(ctx, collection, filter)
}

func (p *PostgresDatabase) CountDocuments(ctx context.Context, collection string, filter any) (int64, error) {
	if err := p.ensureTable(ctx, collection); err != nil {
		return 0, err
	}
	where, args, err := buildFilterQuery(filter)
	if err != nil {
		return 0, err
	}
	var count int64
	err = p.db.QueryRow(ctx, "SELECT COUNT(*) FROM "+quoteIdent(collection)+" WHERE "+where, args...).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (p *PostgresDatabase) EstimatedDocumentCount(ctx context.Context, collection string) (int64, error) {
	var count int64
	err := p.db.QueryRow(ctx, "SELECT COUNT(*) FROM "+quoteIdent(collection)).Scan(&count)
	return count, err
}

func (p *PostgresDatabase) ensureTable(ctx context.Context, collection string) error {
	name := quoteIdent(collection)
	_, err := p.db.Exec(ctx, "CREATE TABLE IF NOT EXISTS "+name+" (id TEXT PRIMARY KEY, payload JSONB NOT NULL)")
	return err
}

func (p *PostgresDatabase) findOneBase(ctx context.Context, collection string, filter any, projection []string, result any) error {
	if err := p.ensureTable(ctx, collection); err != nil {
		return err
	}
	where, args, err := buildFilterQuery(filter)
	if err != nil {
		return err
	}
	query := "SELECT payload FROM " + quoteIdent(collection) + " WHERE " + where + " LIMIT 1"
	var payload []byte
	argsList := make([]any, len(args))
	copy(argsList, args)
	err = p.db.QueryRow(ctx, query, argsList...).Scan(&payload)
	if err != nil {
		return err
	}
	if len(projection) > 0 {
		var doc map[string]any
		if err := json.Unmarshal(payload, &doc); err != nil {
			return err
		}
		filtered := map[string]any{}
		for _, field := range projection {
			if val, ok := doc[field]; ok {
				filtered[field] = val
			}
		}
		payload, err = json.Marshal(filtered)
		if err != nil {
			return err
		}
	}
	return json.Unmarshal(payload, result)
}

func (p *PostgresDatabase) findManyBase(ctx context.Context, collection string, filter any, projection []string, result any) error {
	if err := p.ensureTable(ctx, collection); err != nil {
		return err
	}
	where, args, err := buildFilterQuery(filter)
	if err != nil {
		return err
	}
	query := "SELECT payload FROM " + quoteIdent(collection) + " WHERE " + where + " ORDER BY id"
	argsList := make([]any, len(args))
	copy(argsList, args)
	rows, err := p.db.Query(ctx, query, argsList...)
	if err != nil {
		return err
	}
	defer rows.Close()

	items := make([]json.RawMessage, 0)
	for rows.Next() {
		var payload []byte
		if err := rows.Scan(&payload); err != nil {
			return err
		}
		items = append(items, payload)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if len(items) == 0 {
		return nil
	}
	if len(projection) > 0 {
		filtered := make([]map[string]any, 0, len(items))
		for _, item := range items {
			var doc map[string]any
			if err := json.Unmarshal(item, &doc); err != nil {
				return err
			}
			out := map[string]any{}
			for _, field := range projection {
				if val, ok := doc[field]; ok {
					out[field] = val
				}
			}
			filtered = append(filtered, out)
		}
		payload, err := json.Marshal(filtered)
		if err != nil {
			return err
		}
		return json.Unmarshal(payload, result)
	}
	payload, err := json.Marshal(items)
	if err != nil {
		return err
	}
	return json.Unmarshal(payload, result)
}

func (p *PostgresDatabase) mergeUpdate(ctx context.Context, collection string, filter any, update any) (any, error) {
	var current any
	if err := p.FindOne(ctx, collection, filter, &current); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			if err := p.InsertOne(ctx, collection, update); err != nil {
				return nil, err
			}
			return update, nil
		}
		return nil, err
	}
	baseMap, err := asMap(current)
	if err != nil {
		return nil, err
	}
	incoming, err := asMap(update)
	if err != nil {
		return nil, err
	}
	for k, v := range incoming {
		baseMap[k] = v
	}
	jsonDoc, err := json.Marshal(baseMap)
	if err != nil {
		return nil, err
	}
	id := documentID(current)
	_, err = p.db.Exec(ctx, "UPDATE "+quoteIdent(collection)+" SET payload = $2::jsonb WHERE id = $1", id, string(jsonDoc))
	if err != nil {
		return nil, err
	}
	return baseMap, nil
}

func marshalDocument(document any) ([]byte, error) {
	if document == nil {
		return []byte("{}"), nil
	}
	return json.Marshal(document)
}

func toJSONFilter(filter any) ([]byte, error) {
	if filter == nil {
		return []byte("{}"), nil
	}
	if f, ok := filter.(bson.M); ok {
		return json.Marshal(map[string]any(f))
	}
	if f, ok := filter.(map[string]any); ok {
		return json.Marshal(f)
	}
	return json.Marshal(filter)
}

func asMap(v any) (map[string]any, error) {
	switch x := v.(type) {
	case map[string]any:
		return x, nil
	case bson.M:
		return map[string]any(x), nil
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		out := map[string]any{}
		if err := json.Unmarshal(b, &out); err != nil {
			return nil, err
		}
		return out, nil
	}
}

func documentID(document any) string {
	switch v := document.(type) {
	case map[string]any:
		if val, ok := v["_id"]; ok {
			return fmt.Sprintf("%v", val)
		}
		if val, ok := v["id"]; ok {
			return fmt.Sprintf("%v", val)
		}
	case bson.M:
		if val, ok := v["_id"]; ok {
			return fmt.Sprintf("%v", val)
		}
		if val, ok := v["id"]; ok {
			return fmt.Sprintf("%v", val)
		}
	}
	return uuid.NewString()
}

func quoteIdent(name string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9_]+`)
	return `"` + re.ReplaceAllString(name, "_") + `"`
}

func isRetryablePostgres(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "too many clients")
}
