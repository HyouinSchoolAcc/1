package database

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// EmailConfirmation represents a pending email confirmation
type EmailConfirmation struct {
	ID           string
	Email        string
	Username     string
	PasswordHash string
	Token        string
	ExpiresAt    time.Time
	CreatedAt    time.Time
	Confirmed    bool
}

// RegistrationStore handles database operations for user registration
type RegistrationStore struct {
	db *DB
}

// NewRegistrationStore creates a new RegistrationStore
func NewRegistrationStore(db *DB) *RegistrationStore {
	return &RegistrationStore{db: db}
}

// GenerateToken generates a secure random token
func GenerateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// CreateEmailConfirmation creates a new email confirmation record
func (s *RegistrationStore) CreateEmailConfirmation(email, username, passwordHash string) (*EmailConfirmation, error) {
	token, err := GenerateToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	id := uuid.New().String()
	now := time.Now()
	expiresAt := now.Add(24 * time.Hour) // Token expires in 24 hours

	_, err = s.db.Exec(
		`INSERT INTO email_confirmations (id, email, username, password_hash, token, expires_at, created_at, confirmed)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		id, email, username, passwordHash, token, expiresAt, now, false,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create email confirmation: %w", err)
	}

	return &EmailConfirmation{
		ID:           id,
		Email:        email,
		Username:     username,
		PasswordHash: passwordHash,
		Token:        token,
		ExpiresAt:    expiresAt,
		CreatedAt:    now,
		Confirmed:    false,
	}, nil
}

// GetEmailConfirmationByToken retrieves an email confirmation by token
func (s *RegistrationStore) GetEmailConfirmationByToken(token string) (*EmailConfirmation, error) {
	var conf EmailConfirmation
	err := s.db.QueryRow(
		`SELECT id, email, username, password_hash, token, expires_at, created_at, confirmed
		FROM email_confirmations WHERE token = ?`,
		token,
	).Scan(&conf.ID, &conf.Email, &conf.Username, &conf.PasswordHash, &conf.Token, &conf.ExpiresAt, &conf.CreatedAt, &conf.Confirmed)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get email confirmation: %w", err)
	}

	return &conf, nil
}

// ConfirmEmail confirms an email and creates the user account
func (s *RegistrationStore) ConfirmEmail(token string) (*User, error) {
	// Get confirmation record
	conf, err := s.GetEmailConfirmationByToken(token)
	if err != nil {
		return nil, err
	}
	if conf == nil {
		return nil, fmt.Errorf("invalid confirmation token")
	}
	if conf.Confirmed {
		return nil, fmt.Errorf("email already confirmed")
	}
	if time.Now().After(conf.ExpiresAt) {
		return nil, fmt.Errorf("confirmation token has expired")
	}

	// Create user account - new registrations get "writer" role so they can start writing immediately
	userID := uuid.New().String()
	now := time.Now()

	_, err = s.db.Exec(
		`INSERT INTO users (id, username, email, password, role, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		userID, conf.Username, conf.Email, conf.PasswordHash, "writer", now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Mark confirmation as confirmed
	_, err = s.db.Exec(
		`UPDATE email_confirmations SET confirmed = ? WHERE token = ?`,
		true, token,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update confirmation: %w", err)
	}

	return &User{
		ID:        userID,
		Username:  conf.Username,
		Email:     conf.Email,
		Password:  conf.PasswordHash,
		Role:      "writer",
		CreatedAt: now,
	}, nil
}

// CheckEmailExists checks if an email is already registered or pending confirmation
func (s *RegistrationStore) CheckEmailExists(email string) (bool, error) {
	// Check if email exists in users table
	var count int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM users WHERE email = ?`, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check email in users: %w", err)
	}
	if count > 0 {
		return true, nil
	}

	// Check if email has a pending confirmation
	err = s.db.QueryRow(
		`SELECT COUNT(*) FROM email_confirmations WHERE email = ? AND confirmed = ? AND expires_at > ?`,
		email, false, time.Now(),
	).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check email in confirmations: %w", err)
	}

	return count > 0, nil
}

// CheckUsernameExists checks if a username is already taken
func (s *RegistrationStore) CheckUsernameExists(username string) (bool, error) {
	// Check if username exists in users table
	var count int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM users WHERE username = ?`, username).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check username in users: %w", err)
	}
	if count > 0 {
		return true, nil
	}

	// Check if username has a pending confirmation
	err = s.db.QueryRow(
		`SELECT COUNT(*) FROM email_confirmations WHERE username = ? AND confirmed = ? AND expires_at > ?`,
		username, false, time.Now(),
	).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check username in confirmations: %w", err)
	}

	return count > 0, nil
}
