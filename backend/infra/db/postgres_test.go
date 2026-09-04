package db

import (
	"strings"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestBuildFilterQuerySupportsInOperator(t *testing.T) {
	filter := bson.M{
		"userid": bson.M{"$in": []string{"u1", "u2"}},
		"status": "paid",
	}

	where, args, err := buildFilterQuery(filter)
	if err != nil {
		t.Fatalf("buildFilterQuery returned error: %v", err)
	}
	if !strings.Contains(where, "= ANY") {
		t.Fatalf("expected IN clause in query, got: %s", where)
	}
	if len(args) != 2 {
		t.Fatalf("expected 2 bind args, got %d", len(args))
	}
}

func TestBuildFilterQuerySupportsOrOperator(t *testing.T) {
	filter := bson.M{
		"$or": []any{
			bson.M{"status": "pending"},
			bson.M{"status": "approved"},
		},
	}

	where, _, err := buildFilterQuery(filter)
	if err != nil {
		t.Fatalf("buildFilterQuery returned error: %v", err)
	}
	if !strings.Contains(where, "OR") {
		t.Fatalf("expected OR clause in query, got: %s", where)
	}
}
