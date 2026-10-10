package products

import (
	"reflect"
	"testing"
)

func TestBuildProductSearchQuery(t *testing.T) {
	query, args := buildProductSearchQuery("tool", "Cutting Tools", "cori")

	if query != "1 = 1 AND type = $1 AND category = $2 AND LOWER(name) LIKE $3" {
		t.Fatalf("unexpected query: %s", query)
	}

	wantArgs := []any{"tool", "Cutting Tools", "%cori%"}
	if !reflect.DeepEqual(args, wantArgs) {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func TestBuildProductListOptions(t *testing.T) {
	opts := buildProductListOptions(25, 0, "price_desc")

	if opts.Offset != 25 || opts.Limit != 10 || opts.OrderBy != "price DESC" {
		t.Fatalf("unexpected options: %#v", opts)
	}
}
