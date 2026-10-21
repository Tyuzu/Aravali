package auth

// RegisterRequest defines the payload for creating a new account.
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest defines the payload for a sign-in request.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// OTPRequest defines the payload for OTP-based verification flows.
type OTPRequest struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	OTP      string `json:"otp"`
}
