package cart

import "testing"

func TestResolveLookupTypeAlias(t *testing.T) {
	tests := []struct {
		name     string
		itemType string
		category string
		want     string
	}{
		{name: "tool maps to product lookup", itemType: "tool", category: "Irrigation", want: "product"},
		{name: "product maps to product lookup", itemType: "product", category: "tools", want: "product"},
		{name: "crop maps to crop lookup", itemType: "crop", category: "crops", want: "crop"},
		{name: "unsupported item type falls back to category", itemType: "", category: "products", want: "product"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolveLookupTypeAlias(tt.itemType, tt.category); got != tt.want {
				t.Fatalf("resolveLookupTypeAlias(%q, %q) = %q; want %q", tt.itemType, tt.category, got, tt.want)
			}
		})
	}
}
