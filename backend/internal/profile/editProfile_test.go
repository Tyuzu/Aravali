package profile

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEditProfileRejectsMissingJWT(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/profile", strings.NewReader(""))
	res := httptest.NewRecorder()

	EditProfile(nil).ServeHTTP(res, req)

	if res.Code != http.StatusUnauthorized {
		t.Fatalf("EditProfile() status = %d; want %d", res.Code, http.StatusUnauthorized)
	}
}
