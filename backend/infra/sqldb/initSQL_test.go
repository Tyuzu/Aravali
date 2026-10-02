package sqldb

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUsersMigrationIncludesAuthColumns(t *testing.T) {
	path := filepath.Join("..", "..", "migrations", "001_users_and_profiles.sql")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read migration: %v", err)
	}

	sql := strings.ToLower(string(data))
	for _, column := range []string{
		"userid",
		"password_hash",
		"role",
		"email_verified",
		"is_verified",
		"refresh_token",
		"last_login",
		"online",
		"banner",
		"phone_number",
		"address",
		"social_links",
	} {
		if !strings.Contains(sql, column) {
			t.Fatalf("migration missing required users column %q", column)
		}
	}
}

func TestExtractColumnsAndValuesSkipsNilOptionalFields(t *testing.T) {
	type sample struct {
		ID    string            `db:"id"`
		Tags  []string          `db:"tags"`
		Meta  map[string]string `db:"meta"`
		Name  string            `db:"name"`
		Count int               `db:"count"`
	}

	cols, vals, _, err := extractColumnsAndValues(sample{
		ID:    "u1",
		Name:  "bob",
		Count: 7,
	})
	if err != nil {
		t.Fatalf("extractColumnsAndValues() error = %v", err)
	}

	seen := map[string]bool{}
	for _, c := range cols {
		seen[c] = true
	}
	if !seen["\"id\""] || !seen["\"name\""] || !seen["\"count\""] {
		t.Fatalf("expected required fields to remain: got %v", cols)
	}
	if seen["\"tags\""] || seen["\"meta\""] {
		t.Fatalf("nil optional fields should be skipped: got %v", cols)
	}
	if len(vals) != 3 {
		t.Fatalf("expected 3 values, got %d: %#v", len(vals), vals)
	}
}

func TestParseCreateIndexTarget(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		table    string
		columns  []string
		shouldOK bool
	}{
		{
			name:     "simple column index",
			input:    "CREATE INDEX IF NOT EXISTS idx_users_username ON users (username)",
			table:    "users",
			columns:  []string{"username"},
			shouldOK: true,
		},
		{
			name:     "gin array index",
			input:    "CREATE INDEX IF NOT EXISTS idx_chats_users ON chats USING GIN (users)",
			table:    "chats",
			columns:  []string{"users"},
			shouldOK: true,
		},
		{
			name:     "quoted identifier",
			input:    "CREATE INDEX IF NOT EXISTS idx_messages_read_by ON \"messages\" USING GIN (read_by)",
			table:    "messages",
			columns:  []string{"read_by"},
			shouldOK: true,
		},
		{
			name:     "non-create statement",
			input:    "ALTER TABLE users ADD COLUMN status TEXT",
			shouldOK: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			table, cols, ok := parseCreateIndexTarget(tc.input)
			if ok != tc.shouldOK {
				t.Fatalf("parseCreateIndexTarget(%q) ok=%v, want %v", tc.input, ok, tc.shouldOK)
			}
			if !ok {
				return
			}
			if table != tc.table {
				t.Fatalf("table=%q, want %q", table, tc.table)
			}
			if len(cols) != len(tc.columns) {
				t.Fatalf("len(cols)=%d, want %d", len(cols), len(tc.columns))
			}
			for i := range cols {
				if cols[i] != tc.columns[i] {
					t.Fatalf("cols[%d]=%q, want %q", i, cols[i], tc.columns[i])
				}
			}
		})
	}
}
