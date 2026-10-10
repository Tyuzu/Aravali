package musicon

import "testing"

func TestBuildRecommendationFilter(t *testing.T) {
	cases := []struct {
		name     string
		basedOn  string
		wantKey  string
		wantVal  any
		wantSize int
	}{
		{name: "recently played", basedOn: "recently_played", wantKey: "plays", wantVal: 0, wantSize: 2},
		{name: "language english", basedOn: "language_en", wantKey: "language", wantVal: "en", wantSize: 2},
		{name: "genre pop", basedOn: "genre_pop", wantKey: "genre", wantVal: "Pop", wantSize: 2},
		{name: "default published", basedOn: "", wantKey: "published", wantVal: true, wantSize: 1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := buildRecommendationFilter(tc.basedOn)
			if len(got) != tc.wantSize {
				t.Fatalf("buildRecommendationFilter(%q) size = %d, want %d: %#v", tc.basedOn, len(got), tc.wantSize, got)
			}
			if got[tc.wantKey] != tc.wantVal {
				t.Fatalf("buildRecommendationFilter(%q)[%q] = %#v, want %#v", tc.basedOn, tc.wantKey, got[tc.wantKey], tc.wantVal)
			}
		})
	}
}
