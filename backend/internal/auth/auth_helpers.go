package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"

	"scav/config"
	"scav/middleware"

	"github.com/golang-jwt/jwt/v5"
)

const (
	// Keep the short-lived access token from expiring too quickly while the 7-day
	// refresh token rotates in the background. This prevents unnecessary re-login
	// loops in browsers when the refresh cookie is healthy but the client is idle.
	AccessTokenTTL  = 1 * time.Hour
	RefreshTokenTTL = 7 * 24 * time.Hour
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
	emailRegex    = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
)

/* ============================================================
   HELPERS
============================================================ */

func clientIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		return strings.TrimSpace(strings.Split(fwd, ",")[0])
	}
	if rip := r.Header.Get("X-Real-IP"); rip != "" {
		return rip
	}
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	return host
}

func ipPrefix(ip string) string {
	if strings.Contains(ip, ":") {
		return ip
	}
	parts := strings.Split(ip, ".")
	if len(parts) < 2 {
		return ip
	}
	return parts[0] + "." + parts[1]
}

func uaHash(r *http.Request) string {
	sum := sha256.Sum256([]byte(r.UserAgent()))
	return hex.EncodeToString(sum[:])
}

func hashRefreshToken(token string) string {
	mac := hmac.New(sha256.New, config.RefreshTokenSecret)
	mac.Write([]byte(token))
	return hex.EncodeToString(mac.Sum(nil))
}

func generateRefreshToken() (string, error) {
	b := make([]byte, 64)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func createAccessToken(claims *middleware.Claims) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(config.JwtSecret)
}

func isSecureCookieForHost(scheme, host string) bool {
	scheme = strings.TrimSpace(strings.ToLower(scheme))
	host = strings.TrimSpace(host)

	if scheme == "https" {
		return true
	}
	if scheme == "http" {
		if strings.HasPrefix(host, "localhost") || strings.HasPrefix(host, "127.0.0.1") || strings.HasPrefix(host, "[::1]") {
			return false
		}
		return true
	}

	if strings.HasPrefix(host, "localhost") || strings.HasPrefix(host, "127.0.0.1") || strings.HasPrefix(host, "[::1]") {
		return false
	}

	return true
}

func isSecureCookie(r *http.Request) bool {
	if r == nil {
		return true
	}

	if r.TLS != nil {
		return true
	}

	proto := strings.TrimSpace(r.Header.Get("X-Forwarded-Proto"))
	if proto != "" {
		parts := strings.Split(proto, ",")
		return isSecureCookieForHost(strings.TrimSpace(parts[0]), r.Host)
	}

	proto = strings.TrimSpace(r.Header.Get("X-Forwarded-Scheme"))
	if proto != "" {
		parts := strings.Split(proto, ",")
		return isSecureCookieForHost(strings.TrimSpace(parts[0]), r.Host)
	}

	return isSecureCookieForHost("", r.Host)
}

/* ============================================================
   COOKIE MANAGEMENT
============================================================ */

func setAccessCookie(w http.ResponseWriter, r *http.Request, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   isSecureCookie(r),
		Expires:  time.Now().Add(AccessTokenTTL),
		MaxAge:   int(AccessTokenTTL.Seconds()),
	})
}

func clearAccessCookie(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   isSecureCookie(r),
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})
}

func setRefreshCookie(w http.ResponseWriter, r *http.Request, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   isSecureCookie(r),
		Expires:  time.Now().Add(RefreshTokenTTL),
		MaxAge:   int(RefreshTokenTTL.Seconds()),
	})
}

func clearRefreshCookie(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   isSecureCookie(r),
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})
}

func SetAuthCookies(w http.ResponseWriter, r *http.Request, accessToken, refreshToken string) {
	setAccessCookie(w, r, accessToken)
	setRefreshCookie(w, r, refreshToken)
}

func ClearAuthCookies(w http.ResponseWriter, r *http.Request) {
	clearAccessCookie(w, r)
	clearRefreshCookie(w, r)
}

/* ============================================================
   VALIDATORS
============================================================ */

func validateUsername(u string) bool { return usernameRegex.MatchString(u) }
func validateEmail(e string) bool    { return emailRegex.MatchString(e) }
func validatePassword(p string) bool { return len(p) >= 6 }
