package auth

import (
	"errors"
	"fmt"
	"time"

	"blog2/models"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// JWTConfig holds configuration for JWT tokens
type JWTConfig struct {
	SecretKey     string
	TokenDuration time.Duration
}

// DefaultJWTConfig returns a default JWT configuration
func DefaultJWTConfig() JWTConfig {
	return JWTConfig{
		SecretKey:     "your-secret-key-change-in-production", // Should be set from environment variable in production
		TokenDuration: 24 * time.Hour,                         // 24 hours
	}
}

// GenerateToken creates a new JWT token for a user
func GenerateToken(user models.User, config JWTConfig) (string, error) {
	// Create the claims
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(config.TokenDuration).Unix(),
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken checks if a token is valid and returns the claims
func ValidateToken(tokenString string, config JWTConfig) (models.TokenClaims, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.SecretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return models.TokenClaims{}, ErrExpiredToken
		}
		return models.TokenClaims{}, ErrInvalidToken
	}

	// Check if the token is valid
	if !token.Valid {
		return models.TokenClaims{}, ErrInvalidToken
	}

	// Extract the claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return models.TokenClaims{}, ErrInvalidToken
	}

	// Convert claims to our TokenClaims struct
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return models.TokenClaims{}, ErrInvalidToken
	}

	username, ok := claims["username"].(string)
	if !ok {
		return models.TokenClaims{}, ErrInvalidToken
	}

	return models.TokenClaims{
		UserID:   int(userID),
		Username: username,
	}, nil
}
