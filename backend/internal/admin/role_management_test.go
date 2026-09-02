package admin

import "testing"

func TestNormalizeRoleName(t *testing.T) {
	if got := NormalizeRoleName("  Farmer "); got != "farmer" {
		t.Fatalf("NormalizeRoleName() = %q, want %q", got, "farmer")
	}

	if got := NormalizeRoleName("ADMIN"); got != "admin" {
		t.Fatalf("NormalizeRoleName() = %q, want %q", got, "admin")
	}
}

func TestMergeRoleList(t *testing.T) {
	roles := MergeRoleList([]string{"user", "farmer", "worker"}, "worker", "admin", "farmer")
	want := []string{"user", "farmer", "worker", "admin"}
	if len(roles) != len(want) {
		t.Fatalf("MergeRoleList() length = %d, want %d", len(roles), len(want))
	}
	for i, role := range want {
		if roles[i] != role {
			t.Fatalf("MergeRoleList()[%d] = %q, want %q", i, roles[i], role)
		}
	}
}
