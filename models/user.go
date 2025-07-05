package models

import (
	"time"
)

// User represents a user account in the system
type User struct {
	ID           int        `json:"id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"` // Never expose password hash in JSON responses
	DateCreated  time.Time  `json:"date_created"`
	LastLogin    *time.Time `json:"last_login,omitempty"`
}

// NewUser is used when registering a new user
type NewUser struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginRequest is used for user login
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse is returned after successful login
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// TokenClaims represents the claims in a JWT token
type TokenClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
}
