// File: internal/admin/role_management_test.go

package admin

import (
	"net/http"
	"testing"
)

func TestModeratorApplicationHandlersExist(t *testing.T) {
	var _ http.HandlerFunc = ApplyModerator(nil)
	var _ ModeratorApplication
}

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

func TestNormalizeRoleRequestStatus(t *testing.T) {
	if got := normalizeRoleRequestStatus("PENDING"); got != "pending" {
		t.Fatalf("normalizeRoleRequestStatus() = %q, want %q", got, "pending")
	}
	if got := normalizeRoleRequestStatus(" Approved "); got != "approved" {
		t.Fatalf("normalizeRoleRequestStatus() = %q, want %q", got, "approved")
	}
	if !isFinalRoleRequestStatus("rejected") {
		t.Fatal("isFinalRoleRequestStatus(rejected) = false, want true")
	}
	if isFinalRoleRequestStatus("pending") {
		t.Fatal("isFinalRoleRequestStatus(pending) = true, want false")
	}
}

func TestRoleGrantPermissions(t *testing.T) {
	if !CanAssignRole([]string{"admin"}, "moderator") {
		t.Fatal("admin should be able to assign moderator")
	}

	if !CanAssignRole([]string{"admin"}, "farmer") {
		t.Fatal("admin should be able to assign farmer")
	}

	if CanAssignRole([]string{"moderator"}, "moderator") {
		t.Fatal("moderator should not be able to assign moderator")
	}

	if !CanAssignRole([]string{"moderator"}, "farmer") {
		t.Fatal("moderator should be able to assign farmer")
	}

	roles := AssignableRolesFor([]string{"moderator"})
	if len(roles) == 0 {
		t.Fatal("moderator should have at least one assignable role")
	}
	for _, role := range roles {
		if role == "moderator" {
			t.Fatal("moderator should not be able to assign moderator role")
		}
	}
}
