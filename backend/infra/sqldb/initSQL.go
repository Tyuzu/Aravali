package sqldb

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// EnsureMigrations executes every SQL migration file in order so the relational
// schema exists before the app writes its first records. It is intentionally
// defensive: if a table already exists with a partial legacy schema, the bootstrap
// will still create missing tables and skip incompatible index statements instead
// of crashing the app on startup.
func EnsureMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	if pool == nil {
		return fmt.Errorf("nil postgres pool")
	}

	paths, err := findMigrationFiles()
	if err != nil {
		return err
	}
	if len(paths) == 0 {
		return fmt.Errorf("no migration files found")
	}

	for _, path := range paths {
		body, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", path, err)
		}

		for _, stmt := range splitSQLStatements(string(body)) {
			if shouldSkipStatement(ctx, pool, stmt) {
				continue
			}
			if _, err := pool.Exec(ctx, stmt); err != nil {
				return fmt.Errorf("execute migration %s: %w", filepath.Base(path), err)
			}
		}
	}

	return nil
}

func shouldSkipStatement(ctx context.Context, pool *pgxpool.Pool, stmt string) bool {
	normalized := strings.TrimSpace(stmt)
	if normalized == "" {
		return true
	}

	upper := strings.ToUpper(normalized)
	if !strings.HasPrefix(upper, "CREATE INDEX") {
		return false
	}

	tableName, colNames, ok := parseCreateIndexTarget(normalized)
	if !ok {
		return false
	}
	if tableName == "" {
		return true
	}

	exists, err := tableExists(ctx, pool, tableName)
	if err != nil || !exists {
		return true
	}

	for _, col := range colNames {
		if col == "" {
			continue
		}
		if exists, err := columnExists(ctx, pool, tableName, col); err != nil || !exists {
			return true
		}
	}

	return false
}

func tableExists(ctx context.Context, pool *pgxpool.Pool, tableName string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (
		SELECT 1
		FROM information_schema.tables
		WHERE table_schema = current_schema()
		  AND table_name = $1
	)`
	err := pool.QueryRow(ctx, query, normalizeIdentifier(tableName)).Scan(&exists)
	return exists, err
}

func columnExists(ctx context.Context, pool *pgxpool.Pool, tableName, columnName string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (
		SELECT 1
		FROM information_schema.columns
		WHERE table_schema = current_schema()
		  AND table_name = $1
		  AND column_name = $2
	)`
	err := pool.QueryRow(ctx, query, normalizeIdentifier(tableName), normalizeIdentifier(columnName)).Scan(&exists)
	return exists, err
}

func normalizeIdentifier(name string) string {
	name = strings.TrimSpace(name)
	name = strings.Trim(name, `"`)
	return name
}

func parseCreateIndexTarget(stmt string) (string, []string, bool) {
	pattern := regexp.MustCompile(`(?is)^CREATE\s+INDEX\s+(?:IF\s+NOT\s+EXISTS\s+)?(?:[A-Za-z0-9_]+\s+)?ON\s+("?[A-Za-z0-9_]+"?)\s*(?:USING\s+[A-Za-z0-9_]+)?\s*\((.*)\)\s*$`)
	matches := pattern.FindStringSubmatch(stmt)
	if len(matches) < 3 {
		return "", nil, false
	}

	tableName := normalizeIdentifier(matches[1])
	colsPart := matches[2]
	parts := strings.Split(colsPart, ",")
	cols := make([]string, 0, len(parts))
	for _, col := range parts {
		trimmed := strings.TrimSpace(strings.Trim(col, `"`))
		if trimmed == "" {
			continue
		}
		cols = append(cols, trimmed)
	}
	return tableName, cols, true
}

func findMigrationFiles() ([]string, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("get working directory: %w", err)
	}

	candidates := []string{
		filepath.Join(workDir, "migrations"),
		filepath.Join(workDir, "backend", "migrations"),
		filepath.Join(workDir, "..", "migrations"),
		filepath.Join(workDir, "..", "backend", "migrations"),
	}

	if exe, err := os.Executable(); err == nil {
		baseDir := filepath.Dir(exe)
		candidates = append(candidates,
			filepath.Join(baseDir, "migrations"),
			filepath.Join(baseDir, "backend", "migrations"),
			filepath.Join(baseDir, "..", "backend", "migrations"),
		)
	}

	seen := map[string]struct{}{}
	paths := make([]string, 0, 8)
	for _, dir := range candidates {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
				continue
			}
			path := filepath.Join(dir, entry.Name())
			if _, ok := seen[path]; ok {
				continue
			}
			seen[path] = struct{}{}
			paths = append(paths, path)
		}
	}

	if len(paths) == 0 {
		return nil, fmt.Errorf("migrations directory not found from cwd=%s", workDir)
	}

	sort.Strings(paths)
	return paths, nil
}

func splitSQLStatements(sql string) []string {
	parts := strings.Split(sql, ";")
	statements := make([]string, 0, len(parts))

	for _, part := range parts {
		stmt := strings.TrimSpace(stripSQLComments(part))
		if stmt == "" {
			continue
		}
		statements = append(statements, stmt)
	}

	return statements
}

func stripSQLComments(sql string) string {
	lines := strings.Split(sql, "\n")
	filtered := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "--") {
			continue
		}
		filtered = append(filtered, line)
	}
	return strings.Join(filtered, "\n")
}
