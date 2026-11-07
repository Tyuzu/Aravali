package auth

import (
	"context"
	"testing"

	"scav/infra"
)

func TestBuildUserInitializesEmptySocialLinks(t *testing.T) {
	user, err := BuildUser(SignUpRequest{
		Username: "bob",
		Password: "ValidPass!123",
		Email:    "bob@bob.com",
	})
	if err != nil {
		t.Fatalf("BuildUser() error = %v", err)
	}
	if user.SocialLinks == nil {
		t.Fatal("BuildUser() social_links should be initialized to an empty map")
	}
}

func TestQQQCreateUserRequiresDatabase(t *testing.T) {
	if err := QQQCreateUser(context.Background(), &infra.Deps{}, User{UserID: "u1", Username: "alice"}); err == nil {
		t.Fatal("QQQCreateUser() should fail when no SQL database is configured")
	}
}

func TestQQQFindUserByUsernameRequiresDatabase(t *testing.T) {
	if _, err := QQQFindUserByUsername(context.Background(), &infra.Deps{}, "alice"); err == nil {
		t.Fatal("QQQFindUserByUsername() should fail when no SQL database is configured")
	}
}

func TestPGTextArrayRoleEncodesAsPostgresArray(t *testing.T) {
	role := pgTextArray([]string{"user"})
	if role != "ARRAY['user']::text[]" {
		t.Fatalf("pgTextArray() = %q; want %q", role, "ARRAY['user']::text[]")
	}
}
