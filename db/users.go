package db

import (
	"database/sql"
	"errors"
	"time"

	"blog2/models"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("username or email already exists")
	ErrInvalidCredentials = errors.New("invalid username or password")
)

// CreateUser adds a new user to the database
func (db *DB) CreateUser(nu models.NewUser) (models.User, error) {
	// Check if username or email already exists
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM users WHERE username = $1 OR email = $2
		)
	`, nu.Username, nu.Email).Scan(&exists)

	if err != nil {
		return models.User{}, err
	}

	if exists {
		return models.User{}, ErrUserAlreadyExists
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	// Insert the new user
	var user models.User
	err = db.QueryRow(`
		INSERT INTO users (username, email, password_hash) 
		VALUES ($1, $2, $3) 
		RETURNING id, username, email, password_hash, date_created, last_login
	`, nu.Username, nu.Email, string(hashedPassword)).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.DateCreated, &user.LastLogin,
	)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// GetUserByUsername retrieves a user by username
func (db *DB) GetUserByUsername(username string) (models.User, error) {
	var user models.User
	err := db.QueryRow(`
		SELECT id, username, email, password_hash, date_created, last_login 
		FROM users 
		WHERE username = $1
	`, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.DateCreated, &user.LastLogin,
	)

	if err == sql.ErrNoRows {
		return models.User{}, ErrUserNotFound
	}

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (db *DB) GetUserByID(id int) (models.User, error) {
	var user models.User
	err := db.QueryRow(`
		SELECT id, username, email, password_hash, date_created, last_login 
		FROM users 
		WHERE id = $1
	`, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.DateCreated, &user.LastLogin,
	)

	if err == sql.ErrNoRows {
		return models.User{}, ErrUserNotFound
	}

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// UpdateLastLogin updates the last_login timestamp for a user
func (db *DB) UpdateLastLogin(userID int) error {
	now := time.Now()
	_, err := db.Exec(`
		UPDATE users 
		SET last_login = $1 
		WHERE id = $2
	`, now, userID)

	return err
}

// AuthenticateUser checks if the provided credentials are valid
func (db *DB) AuthenticateUser(login models.LoginRequest) (models.User, error) {
	user, err := db.GetUserByUsername(login.Username)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return models.User{}, ErrInvalidCredentials
		}
		return models.User{}, err
	}

	// Compare the provided password with the stored hash
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(login.Password))
	if err != nil {
		return models.User{}, ErrInvalidCredentials
	}

	// Update last login time
	err = db.UpdateLastLogin(user.ID)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
