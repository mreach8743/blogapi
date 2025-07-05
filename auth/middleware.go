package auth

import (
	"context"
	"net/http"
	"strings"

	"blog2/models"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

// Context keys
const (
	UserClaimsKey contextKey = "user_claims"
)

// AuthMiddleware creates middleware for authenticating requests
func AuthMiddleware(config JWTConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			// Check if it's a Bearer token
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Invalid authorization format, Bearer token required", http.StatusUnauthorized)
				return
			}

			// Extract the token
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// Validate the token
			claims, err := ValidateToken(tokenString, config)
			if err != nil {
				if err == ErrExpiredToken {
					http.Error(w, "Token has expired", http.StatusUnauthorized)
				} else {
					http.Error(w, "Invalid token", http.StatusUnauthorized)
				}
				return
			}

			// Add the claims to the request context
			ctx := context.WithValue(r.Context(), UserClaimsKey, claims)

			// Call the next handler with the updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserClaims extracts user claims from the request context
func GetUserClaims(r *http.Request) (models.TokenClaims, bool) {
	claims, ok := r.Context().Value(UserClaimsKey).(models.TokenClaims)
	return claims, ok
}

// RequireAuth is a middleware that requires authentication for specific routes
func RequireAuth(config JWTConfig) func(http.Handler) http.Handler {
	return AuthMiddleware(config)
}

// OptionalAuth is middleware that adds user claims to context if token is present, but doesn't require it
func OptionalAuth(config JWTConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the Authorization header
			authHeader := r.Header.Get("Authorization")

			// If no Authorization header, just continue
			if authHeader == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Check if it's a Bearer token
			if !strings.HasPrefix(authHeader, "Bearer ") {
				next.ServeHTTP(w, r)
				return
			}

			// Extract the token
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// Validate the token
			claims, err := ValidateToken(tokenString, config)
			if err == nil {
				// Add the claims to the request context
				ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
				r = r.WithContext(ctx)
			}

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}
