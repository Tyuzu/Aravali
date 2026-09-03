package auth

import (
	"testing"
	"time"
)

func TestAccessTokenTTLIsLongerThanShortSession(t *testing.T) {
	if AccessTokenTTL <= 30*time.Minute {
		t.Fatalf("AccessTokenTTL = %s; expected a longer session window to avoid forced re-login", AccessTokenTTL)
	}
}

func TestCookieSecurityModeForLocalhost(t *testing.T) {
	tests := []struct {
		scheme string
		host   string
		want   bool
	}{
		{scheme: "https", host: "app.example.com", want: true},
		{scheme: "https", host: "localhost:5173", want: true},
		{scheme: "http", host: "localhost:5173", want: false},
		{scheme: "http", host: "127.0.0.1:5173", want: false},
		{scheme: "http", host: "app.example.com", want: true},
		{scheme: "", host: "localhost:5173", want: false},
	}

	for _, tc := range tests {
		got := isSecureCookieForHost(tc.scheme, tc.host)
		if got != tc.want {
			t.Fatalf("isSecureCookieForHost(%q, %q) = %v; want %v", tc.scheme, tc.host, got, tc.want)
		}
	}
}
