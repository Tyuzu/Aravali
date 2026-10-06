package farms

import (
	"reflect"
	"testing"
)

func TestMergeFarmerRole(t *testing.T) {
	tests := []struct {
		name     string
		existing []string
		want     []string
	}{
		{
			name:     "adds farmer when missing",
			existing: []string{"user", "worker"},
			want:     []string{"user", "worker", "farmer"},
		},
		{
			name:     "keeps farmer once present",
			existing: []string{"user", "farmer", "worker"},
			want:     []string{"user", "farmer", "worker"},
		},
		{
			name:     "normalizes casing and blanks",
			existing: []string{" User ", "FARMER", "", "admin"},
			want:     []string{"user", "farmer", "admin"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mergeFarmerRole(tt.existing); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("mergeFarmerRole() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
