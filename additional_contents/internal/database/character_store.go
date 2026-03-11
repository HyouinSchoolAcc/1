package database

import (
	"database/sql"
	"fmt"
	"time"
)

// UserCharacter represents a character created by a user
type UserCharacter struct {
	ID          string    `json:"id"`
	CreatorID   string    `json:"creator_id"`
	CreatorName string    `json:"creator_name"`
	Name        string    `json:"name"`
	Values      string    `json:"values"`
	Experiences string    `json:"experiences"`
	Judgements  string    `json:"judgements"`
	Abilities   string    `json:"abilities"`
	Story       string    `json:"story"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CharacterStore manages user-created characters
type CharacterStore struct {
	db *sql.DB
}

// NewCharacterStore creates a new character store
func NewCharacterStore(db *DB) *CharacterStore {
	return &CharacterStore{db: db.DB}
}

// CreateCharacter creates a new user character
func (cs *CharacterStore) CreateCharacter(character *UserCharacter) error {
	query := `
		INSERT INTO user_characters (id, creator_id, creator_name, name, character_values, experiences, judgements, abilities, story, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	character.CreatedAt = now
	character.UpdatedAt = now

	_, err := cs.db.Exec(query,
		character.ID,
		character.CreatorID,
		character.CreatorName,
		character.Name,
		character.Values,
		character.Experiences,
		character.Judgements,
		character.Abilities,
		character.Story,
		now,
		now,
	)

	return err
}

// GetCharacterByID retrieves a character by ID
func (cs *CharacterStore) GetCharacterByID(id string) (*UserCharacter, error) {
	query := `
		SELECT id, creator_id, creator_name, name, character_values, experiences, judgements, abilities, story, created_at, updated_at
		FROM user_characters
		WHERE id = ?
	`

	character := &UserCharacter{}
	err := cs.db.QueryRow(query, id).Scan(
		&character.ID,
		&character.CreatorID,
		&character.CreatorName,
		&character.Name,
		&character.Values,
		&character.Experiences,
		&character.Judgements,
		&character.Abilities,
		&character.Story,
		&character.CreatedAt,
		&character.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("character not found")
	}
	if err != nil {
		return nil, err
	}

	return character, nil
}

// GetCharactersByCreator retrieves all characters created by a user
func (cs *CharacterStore) GetCharactersByCreator(creatorID string) ([]*UserCharacter, error) {
	query := `
		SELECT id, creator_id, creator_name, name, character_values, experiences, judgements, abilities, story, created_at, updated_at
		FROM user_characters
		WHERE creator_id = ?
		ORDER BY created_at DESC
	`

	rows, err := cs.db.Query(query, creatorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var characters []*UserCharacter
	for rows.Next() {
		character := &UserCharacter{}
		err := rows.Scan(
			&character.ID,
			&character.CreatorID,
			&character.CreatorName,
			&character.Name,
			&character.Values,
			&character.Experiences,
			&character.Judgements,
			&character.Abilities,
			&character.Story,
			&character.CreatedAt,
			&character.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		characters = append(characters, character)
	}

	return characters, nil
}

// GetAllCharacters retrieves all characters
func (cs *CharacterStore) GetAllCharacters() ([]*UserCharacter, error) {
	query := `
		SELECT id, creator_id, creator_name, name, character_values, experiences, judgements, abilities, story, created_at, updated_at
		FROM user_characters
		ORDER BY created_at DESC
	`

	rows, err := cs.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var characters []*UserCharacter
	for rows.Next() {
		character := &UserCharacter{}
		err := rows.Scan(
			&character.ID,
			&character.CreatorID,
			&character.CreatorName,
			&character.Name,
			&character.Values,
			&character.Experiences,
			&character.Judgements,
			&character.Abilities,
			&character.Story,
			&character.CreatedAt,
			&character.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		characters = append(characters, character)
	}

	return characters, nil
}

// UpdateCharacter updates an existing character
func (cs *CharacterStore) UpdateCharacter(character *UserCharacter) error {
	query := `
		UPDATE user_characters
		SET name = ?, character_values = ?, experiences = ?, judgements = ?, abilities = ?, story = ?, updated_at = ?
		WHERE id = ? AND creator_id = ?
	`

	character.UpdatedAt = time.Now()

	result, err := cs.db.Exec(query,
		character.Name,
		character.Values,
		character.Experiences,
		character.Judgements,
		character.Abilities,
		character.Story,
		character.UpdatedAt,
		character.ID,
		character.CreatorID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("character not found or unauthorized")
	}

	return nil
}

// DeleteCharacter deletes a character (only by creator)
func (cs *CharacterStore) DeleteCharacter(id, creatorID string) error {
	query := `DELETE FROM user_characters WHERE id = ? AND creator_id = ?`

	result, err := cs.db.Exec(query, id, creatorID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("character not found or unauthorized")
	}

	return nil
}

