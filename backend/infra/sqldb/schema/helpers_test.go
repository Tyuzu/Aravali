package schema

import "testing"

func TestNormalizeIdentifier(t *testing.T) {
	if got, want := NormalizeIdentifier("\"users\""), "users"; got != want {
		t.Fatalf("NormalizeIdentifier() = %q, want %q", got, want)
	}
	if got, want := NormalizeIdentifier("  farms  "), "farms"; got != want {
		t.Fatalf("NormalizeIdentifier() = %q, want %q", got, want)
	}
}

func TestQuoteIdent(t *testing.T) {
	if got, want := QuoteIdent("user profile"), "\"user_profile\""; got != want {
		t.Fatalf("QuoteIdent() = %q, want %q", got, want)
	}
	if got, want := QuoteIdent("created_at"), "\"created_at\""; got != want {
		t.Fatalf("QuoteIdent() = %q, want %q", got, want)
	}
}
