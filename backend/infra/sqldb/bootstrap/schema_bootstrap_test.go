package bootstrap

import (
	"context"
	"strings"
	"testing"
)

func TestEnsureModuleSchemasRejectsNilPool(t *testing.T) {
	err := EnsureModuleSchemas(context.Background(), nil)
	if err == nil {
		t.Fatal("EnsureModuleSchemas(nil) should return an error")
	}
	if !strings.Contains(err.Error(), "nil postgres pool") {
		t.Fatalf("unexpected error: %v", err)
	}
}
