package auth

import "testing"

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
