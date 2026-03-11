package database

import (
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// UserRole represents the type of user in the system
type UserRole string

const (
	RoleNewUser UserRole = "new_user" // Registered but not yet a writer (viewing only)
	RoleViewer  UserRole = "viewer"   // Explicitly view-only role (if needed)
	RoleWriter  UserRole = "writer"   // Passed tutorial, can write
	RoleEditor  UserRole = "editor"   // Full access, can approve/QC
)

// User represents a user account
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"` // bcrypt hashed
	Role      UserRole  `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// UserStore handles user data persistence in SQL
type UserStore struct {
	db *DB
}

// NewUserStore creates a new SQL-based user store
func NewUserStore(db *DB) *UserStore {
	return &UserStore{db: db}
}

// CreateUser creates a new user account
func (us *UserStore) CreateUser(username, email, password string, role UserRole) error {
	// Check if username already exists
	var exists bool
	err := us.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", username).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("username already exists")
	}

	// Check if email already exists
	err = us.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", email).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := User{
		ID:        generateID(),
		Username:  username,
		Email:     email,
		Password:  string(hashedPassword),
		Role:      role,
		CreatedAt: time.Now(),
	}

	_, err = us.db.Exec(
		"INSERT INTO users (id, username, email, password, role, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		user.ID, user.Username, user.Email, user.Password, string(user.Role), user.CreatedAt,
	)
	return err
}

// ValidateUser validates a user's credentials
func (us *UserStore) ValidateUser(username, password string) (*User, error) {
	user := &User{}
	var roleStr string
	
	err := us.db.QueryRow(
		"SELECT id, username, email, password, role, created_at FROM users WHERE username = ?",
		username,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &roleStr, &user.CreatedAt)
	
	if err == sql.ErrNoRows {
		return nil, errors.New("invalid credentials")
	}
	if err != nil {
		return nil, err
	}

	user.Role = UserRole(roleStr)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

// GetUser retrieves a user by username
func (us *UserStore) GetUser(username string) (*User, bool) {
	user := &User{}
	var roleStr string
	
	err := us.db.QueryRow(
		"SELECT id, username, email, password, role, created_at FROM users WHERE username = ?",
		username,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &roleStr, &user.CreatedAt)
	
	if err == sql.ErrNoRows {
		return nil, false
	}
	if err != nil {
		return nil, false
	}

	user.Role = UserRole(roleStr)
	return user, true
}

// GetUserByID retrieves a user by ID
func (us *UserStore) GetUserByID(id string) (*User, error) {
	user := &User{}
	var roleStr string
	
	err := us.db.QueryRow(
		"SELECT id, username, email, password, role, created_at FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &roleStr, &user.CreatedAt)
	
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	user.Role = UserRole(roleStr)
	return user, nil
}

// ListUsers returns all users
func (us *UserStore) ListUsers() ([]*User, error) {
	rows, err := us.db.Query("SELECT id, username, email, password, role, created_at FROM users ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*User{}
	for rows.Next() {
		user := &User{}
		var roleStr string
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &roleStr, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		user.Role = UserRole(roleStr)
		users = append(users, user)
	}

	return users, rows.Err()
}

// UpdateUser updates user information
func (us *UserStore) UpdateUser(user *User) error {
	_, err := us.db.Exec(
		"UPDATE users SET username = ?, email = ?, role = ? WHERE id = ?",
		user.Username, user.Email, string(user.Role), user.ID,
	)
	return err
}

// DeleteUser deletes a user
func (us *UserStore) DeleteUser(id string) error {
	_, err := us.db.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}

// GetUserByEmail retrieves a user by email address
func (us *UserStore) GetUserByEmail(email string) (*User, error) {
	user := &User{}
	var roleStr string
	
	err := us.db.QueryRow(
		"SELECT id, username, email, password, role, created_at FROM users WHERE email = ?",
		email,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &roleStr, &user.CreatedAt)
	
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	user.Role = UserRole(roleStr)
	return user, nil
}

// ResetPasswordByEmail resets a user's password by email and returns the new password
func (us *UserStore) ResetPasswordByEmail(email string) (string, string, error) {
	user, err := us.GetUserByEmail(email)
	if err != nil {
		return "", "", err
	}
	
	// Reset to default password
	newPassword := "password123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}
	
	_, err = us.db.Exec("UPDATE users SET password = ? WHERE id = ?", string(hashedPassword), user.ID)
	if err != nil {
		return "", "", err
	}
	
	return user.Username, newPassword, nil
}

// generateID generates a simple ID based on timestamp
func generateID() string {
	return time.Now().Format("20060102150405")
}

