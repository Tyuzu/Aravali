// File: internal/auth/types.go

package auth

import "time"

type SignUpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type SignUpResponse struct {
	Message string `json:"message"`
	UserID  string `json:"userid"`
}

const (
	// Keep the short-lived access token from expiring too quickly while the 7-day
	// refresh token rotates in the background. This prevents unnecessary re-login
	// loops in browsers when the refresh cookie is healthy but the client is idle.
	AccessTokenTTL    = 1 * time.Hour
	maxFailedAttempts = 5
	lockoutDuration   = 15 * time.Minute
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Token   string `json:"token"`
	UserID  string `json:"userid"`
}

// Structural Data Transfers
type RequestOTPInput struct {
	Email string `json:"email"`
}

type VerifyOTPInput struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

// RefreshResult communicates intended cookie side-effects and tokens.
type RefreshResult struct {
	UserID      string
	AccessToken string
	NewRefresh  string // non-empty => set this new refresh in cookie
	ClearCookie bool   // true => clear cookie on response
}

type User struct {
	UserID       string    `json:"userid" db:"userid"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	Password     string    `json:"-" db:"password"`
	PasswordHash string    `json:"password_hash" db:"password_hash"`
	Role         []string  `json:"role" db:"role"`
	Name         string    `json:"name,omitempty" db:"name,omitempty"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	Bio          string    `json:"bio,omitempty" db:"bio,omitempty"`
	Online       bool      `json:"online"`
	LastLogin    time.Time `json:"last_login" db:"last_login"`
	Avatar       string    `json:"avatar" db:"avatar"`
	Banner       string    `json:"banner" db:"banner"`
	ProfileViews int       `json:"profile_views,omitempty" db:"profile_views,omitempty"`
	PhoneNumber  string    `json:"phone_number,omitempty" db:"phone_number,omitempty"`
	Address      string    `json:"address,omitempty" db:"address,omitempty"`
	// DateOfBirth    time.Time         `json:"dob" db:"dob"`
	SocialLinks    map[string]string `json:"social_links,omitempty" db:"social_links,omitempty"`
	IsVerified     bool              `json:"is_verified" db:"is_verified"`
	EmailVerified  bool              `json:"email_verified" db:"email_verified"`
	FollowersCount int               `json:"followerscount" db:"followerscount"`
	FollowingCount int               `json:"followscount" db:"followscount"`
	WalletBalance  float64           `db:"wallet_balance" json:"wallet_balance"`
	RefreshToken   string            `json:"-" db:"refresh_token,omitempty"`
	RefreshExpiry  time.Time         `json:"-" db:"refresh_expiry,omitempty"`
	RefreshUA      string            `db:"refresh_ua,omitempty"`
	RefreshIP      string            `db:"refresh_ip,omitempty"`
	RefreshPrev    string            `db:"refresh_prev,omitempty"`
}
