package models

import "github.com/golang-jwt/jwt/v5"

// Standard Response API Envelope
type APIResponse struct {
	Success   bool        `json:"success"`
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data,omitempty"`
	Message   string      `json:"message,omitempty"`
}

// -----------------------------------------------------------------------------
// Data Contracts
// -----------------------------------------------------------------------------

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  string `json:"user"`
}

type FeedData struct {
	Title string   `json:"title"`
	Body  string   `json:"body"`
	Tags  []string `json:"tags"`
}

type ProfileData struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Status string `json:"status"`
}

type AnalyticsData struct {
	DailyActiveUsers int     `json:"dailyActiveUsers"`
	TotalRevenue     float64 `json:"totalRevenue"`
	ServerUptime     string  `json:"serverUptime"`
	ConversionRate   string  `json:"conversionRate"`
}

type SettingsData struct {
	NotificationsEnabled bool   `json:"notificationsEnabled"`
	DarkMode             bool   `json:"darkMode"`
	ApiVersion           string `json:"apiVersion"`
	Theme                string `json:"theme"`
}

type AuthRequest struct {
	Name     string `json:"name,omitempty"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthData struct {
	Token string `json:"token"`
	User  string `json:"user"`
}

// Custom JWT Claims
type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}
