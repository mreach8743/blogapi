package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"blog2/auth"
	"blog2/db"
	"blog2/models"
)

// UsersHandler handles all user-related HTTP requests
type UsersHandler struct {
	DB        *db.DB
	Validator *validator.Validate
	JWTConfig auth.JWTConfig
}

// NewUsersHandler creates a new UsersHandler
func NewUsersHandler(db *db.DB, jwtConfig auth.JWTConfig) *UsersHandler {
	return &UsersHandler{
		DB:        db,
		Validator: validator.New(),
		JWTConfig: jwtConfig,
	}
}

// ServeHTTP handles all HTTP requests for users
func (h *UsersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Route based on HTTP method and path
	switch {
	case r.Method == http.MethodPost && r.URL.Path == "/users/register":
		h.registerUser(w, r)
	case r.Method == http.MethodPost && r.URL.Path == "/users/login":
		h.loginUser(w, r)
	case r.Method == http.MethodGet && r.URL.Path == "/users/me":
		h.getCurrentUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// registerUser handles user registration
func (h *UsersHandler) registerUser(w http.ResponseWriter, r *http.Request) {
	var newUser models.NewUser
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the input
	if err := h.Validator.Struct(newUser); err != nil {
		http.Error(w, "Validation error: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Create the user
	user, err := h.DB.CreateUser(newUser)
	if err != nil {
		if errors.Is(err, db.ErrUserAlreadyExists) {
			http.Error(w, "Username or email already exists", http.StatusConflict)
		} else {
			http.Error(w, "Error creating user: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Generate a token for the new user
	token, err := auth.GenerateToken(user, h.JWTConfig)
	if err != nil {
		http.Error(w, "Error generating token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the user and token
	response := models.LoginResponse{
		Token: token,
		User:  user,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// loginUser handles user login
func (h *UsersHandler) loginUser(w http.ResponseWriter, r *http.Request) {
	var loginRequest models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the input
	if err := h.Validator.Struct(loginRequest); err != nil {
		http.Error(w, "Validation error: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Authenticate the user
	user, err := h.DB.AuthenticateUser(loginRequest)
	if err != nil {
		if errors.Is(err, db.ErrInvalidCredentials) {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		} else {
			http.Error(w, "Error authenticating user: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Generate a token
	token, err := auth.GenerateToken(user, h.JWTConfig)
	if err != nil {
		http.Error(w, "Error generating token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the user and token
	response := models.LoginResponse{
		Token: token,
		User:  user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// getCurrentUser returns the current authenticated user
func (h *UsersHandler) getCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Get the user claims from the context
	claims, ok := auth.GetUserClaims(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get the user from the database
	user, err := h.DB.GetUserByID(claims.UserID)
	if err != nil {
		if errors.Is(err, db.ErrUserNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error retrieving user: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
