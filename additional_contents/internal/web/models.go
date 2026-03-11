package web

import (
	"time"

	"data_labler_ui_go/internal/database"
)

// Re-export types from database package for backward compatibility
type UserRole = database.UserRole

const (
	RoleNewUser = database.RoleNewUser
	RoleViewer  = database.RoleViewer
	RoleWriter  = database.RoleWriter
	RoleEditor  = database.RoleEditor
)

// Re-export tutorial types
type TutorialProgress = database.TutorialProgress
type CharacterCertification = database.CharacterCertification

type User = database.User
type UserStore = database.UserStore
type CharacterStore = database.CharacterStore
type UserCharacter = database.UserCharacter

// SessionData holds session information
type SessionData struct {
	UserID   string
	Username string
	Role     UserRole
}

// StarredFile represents a file marked as starred for new users
type StarredFile struct {
	FilePath string    `json:"file_path"`
	AddedAt  time.Time `json:"added_at"`
}

// generateID generates a simple ID (in production, use UUID)
func generateID() string {
	return time.Now().Format("20060102150405")
}

