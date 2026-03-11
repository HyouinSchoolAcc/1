package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"
)

// TutorialProgress represents a user's progress through the writer tutorial
type TutorialProgress struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	CurrentStep    int       `json:"current_step"`
	CompletedSteps []string  `json:"completed_steps"`
	StartedAt      time.Time `json:"started_at"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
	IsCompleted    bool      `json:"is_completed"`
}

// CharacterCertification represents a user's certification for a specific character
type CharacterCertification struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	CharacterID string    `json:"character_id"`
	CertifiedAt time.Time `json:"certified_at"`
	CertifiedBy string    `json:"certified_by,omitempty"` // empty if self-certified
}

// TutorialStore handles tutorial progress and certifications
type TutorialStore struct {
	db *DB
}

// NewTutorialStore creates a new tutorial store
func NewTutorialStore(db *DB) *TutorialStore {
	return &TutorialStore{db: db}
}

// GetTutorialProgress retrieves a user's tutorial progress
func (ts *TutorialStore) GetTutorialProgress(userID string) (*TutorialProgress, error) {
	progress := &TutorialProgress{}
	var completedStepsJSON string
	var completedAt sql.NullTime

	err := ts.db.QueryRow(
		`SELECT id, user_id, current_step, completed_steps, started_at, completed_at, is_completed 
		 FROM tutorial_progress WHERE user_id = ?`,
		userID,
	).Scan(&progress.ID, &progress.UserID, &progress.CurrentStep, &completedStepsJSON, 
		&progress.StartedAt, &completedAt, &progress.IsCompleted)

	if err == sql.ErrNoRows {
		return nil, nil // No progress yet
	}
	if err != nil {
		return nil, err
	}

	// Parse completed steps JSON
	if err := json.Unmarshal([]byte(completedStepsJSON), &progress.CompletedSteps); err != nil {
		progress.CompletedSteps = []string{}
	}

	if completedAt.Valid {
		progress.CompletedAt = &completedAt.Time
	}

	return progress, nil
}

// StartTutorial creates or resets a user's tutorial progress
func (ts *TutorialStore) StartTutorial(userID string) (*TutorialProgress, error) {
	// Check if progress already exists
	existing, err := ts.GetTutorialProgress(userID)
	if err != nil {
		return nil, err
	}

	if existing != nil && !existing.IsCompleted {
		// Return existing progress
		return existing, nil
	}

	// Create new progress
	progress := &TutorialProgress{
		ID:             generateID(),
		UserID:         userID,
		CurrentStep:    0,
		CompletedSteps: []string{},
		StartedAt:      time.Now(),
		IsCompleted:    false,
	}

	completedStepsJSON, _ := json.Marshal(progress.CompletedSteps)

	_, err = ts.db.Exec(
		`INSERT OR REPLACE INTO tutorial_progress 
		 (id, user_id, current_step, completed_steps, started_at, is_completed) 
		 VALUES (?, ?, ?, ?, ?, ?)`,
		progress.ID, progress.UserID, progress.CurrentStep, string(completedStepsJSON),
		progress.StartedAt, progress.IsCompleted,
	)

	if err != nil {
		return nil, err
	}

	return progress, nil
}

// UpdateTutorialProgress updates a user's tutorial progress
func (ts *TutorialStore) UpdateTutorialProgress(userID string, currentStep int, completedSteps []string) error {
	completedStepsJSON, _ := json.Marshal(completedSteps)

	result, err := ts.db.Exec(
		`UPDATE tutorial_progress SET current_step = ?, completed_steps = ? WHERE user_id = ?`,
		currentStep, string(completedStepsJSON), userID,
	)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("tutorial progress not found")
	}

	return nil
}

// CompleteTutorial marks a user's tutorial as completed
func (ts *TutorialStore) CompleteTutorial(userID string) error {
	now := time.Now()
	_, err := ts.db.Exec(
		`UPDATE tutorial_progress SET is_completed = 1, completed_at = ? WHERE user_id = ?`,
		now, userID,
	)
	return err
}

// GetUserCertifications retrieves all character certifications for a user
func (ts *TutorialStore) GetUserCertifications(userID string) ([]CharacterCertification, error) {
	rows, err := ts.db.Query(
		`SELECT id, user_id, character_id, certified_at, COALESCE(certified_by, '') 
		 FROM character_certifications WHERE user_id = ? ORDER BY certified_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	certs := []CharacterCertification{}
	for rows.Next() {
		cert := CharacterCertification{}
		err := rows.Scan(&cert.ID, &cert.UserID, &cert.CharacterID, &cert.CertifiedAt, &cert.CertifiedBy)
		if err != nil {
			return nil, err
		}
		certs = append(certs, cert)
	}

	return certs, rows.Err()
}

// IsCertified checks if a user is certified for a specific character
func (ts *TutorialStore) IsCertified(userID, characterID string) (bool, error) {
	var exists bool
	err := ts.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM character_certifications WHERE user_id = ? AND character_id = ?)`,
		userID, characterID,
	).Scan(&exists)
	return exists, err
}

// CertifyUser certifies a user for a specific character
func (ts *TutorialStore) CertifyUser(userID, characterID, certifiedBy string) error {
	cert := CharacterCertification{
		ID:          generateID(),
		UserID:      userID,
		CharacterID: characterID,
		CertifiedAt: time.Now(),
		CertifiedBy: certifiedBy,
	}

	_, err := ts.db.Exec(
		`INSERT OR REPLACE INTO character_certifications 
		 (id, user_id, character_id, certified_at, certified_by) 
		 VALUES (?, ?, ?, ?, ?)`,
		cert.ID, cert.UserID, cert.CharacterID, cert.CertifiedAt, cert.CertifiedBy,
	)
	return err
}

// RevokeCertification removes a user's certification for a character
func (ts *TutorialStore) RevokeCertification(userID, characterID string) error {
	_, err := ts.db.Exec(
		`DELETE FROM character_certifications WHERE user_id = ? AND character_id = ?`,
		userID, characterID,
	)
	return err
}

