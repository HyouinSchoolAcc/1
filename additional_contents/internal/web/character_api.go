package web

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// handleGetUserCharacters returns all characters (optionally filtered by creator)
func (a *App) handleGetUserCharacters(w http.ResponseWriter, r *http.Request) {
	sessionData := a.getCurrentUser(r)
	
	var characters []*UserCharacter
	var err error
	
	// If user is logged in, get their characters
	if sessionData != nil {
		characters, err = a.characterStore.GetCharactersByCreator(sessionData.UserID)
	} else {
		// If not logged in, get all characters (for display purposes)
		characters, err = a.characterStore.GetAllCharacters()
	}
	
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	
	if characters == nil {
		characters = []*UserCharacter{}
	}
	
	writeJSON(w, http.StatusOK, map[string]any{
		"success":    true,
		"characters": characters,
	})
}

// handleGetUserCharacter returns a specific character by ID
func (a *App) handleGetUserCharacter(w http.ResponseWriter, r *http.Request) {
	characterID := chi.URLParam(r, "id")
	if characterID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "character ID required",
		})
		return
	}
	
	character, err := a.characterStore.GetCharacterByID(characterID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]any{
			"success": false,
			"error":   "character not found",
		})
		return
	}
	
	writeJSON(w, http.StatusOK, map[string]any{
		"success":   true,
		"character": character,
	})
}

// handleCreateUserCharacter creates a new character
func (a *App) handleCreateUserCharacter(w http.ResponseWriter, r *http.Request) {
	// Require authentication
	sessionData := a.requireAuth(w, r)
	if sessionData == nil {
		return
	}
	
	var req struct {
		Name        string `json:"name"`
		Values      string `json:"values"`
		Experiences string `json:"experiences"`
		Judgements  string `json:"judgements"`
		Abilities   string `json:"abilities"`
		Story       string `json:"story"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "invalid request body",
		})
		return
	}
	
	// Validate required fields
	if req.Name == "" || req.Values == "" || req.Experiences == "" || 
	   req.Judgements == "" || req.Abilities == "" || req.Story == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "all fields are required",
		})
		return
	}
	
	// Create character
	character := &UserCharacter{
		ID:          generateID(),
		CreatorID:   sessionData.UserID,
		CreatorName: sessionData.Username,
		Name:        req.Name,
		Values:      req.Values,
		Experiences: req.Experiences,
		Judgements:  req.Judgements,
		Abilities:   req.Abilities,
		Story:       req.Story,
	}
	
	if err := a.characterStore.CreateCharacter(character); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   "failed to create character",
		})
		return
	}
	
	writeJSON(w, http.StatusOK, map[string]any{
		"success":   true,
		"character": character,
	})
}

// handleUpdateUserCharacter updates an existing character
func (a *App) handleUpdateUserCharacter(w http.ResponseWriter, r *http.Request) {
	// Require authentication
	sessionData := a.requireAuth(w, r)
	if sessionData == nil {
		return
	}
	
	characterID := chi.URLParam(r, "id")
	if characterID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "character ID required",
		})
		return
	}
	
	var req struct {
		Name        string `json:"name"`
		Values      string `json:"values"`
		Experiences string `json:"experiences"`
		Judgements  string `json:"judgements"`
		Abilities   string `json:"abilities"`
		Story       string `json:"story"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "invalid request body",
		})
		return
	}
	
	// Get existing character to verify ownership
	existing, err := a.characterStore.GetCharacterByID(characterID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]any{
			"success": false,
			"error":   "character not found",
		})
		return
	}
	
	// Verify user owns this character OR is an editor
	if existing.CreatorID != sessionData.UserID && sessionData.Role != RoleEditor {
		writeJSON(w, http.StatusForbidden, map[string]any{
			"success": false,
			"error":   "you can only edit your own characters",
		})
		return
	}
	
	// Update character - preserve original creator info when editor edits
	character := &UserCharacter{
		ID:          characterID,
		CreatorID:   existing.CreatorID,
		CreatorName: existing.CreatorName,
		Name:        req.Name,
		Values:      req.Values,
		Experiences: req.Experiences,
		Judgements:  req.Judgements,
		Abilities:   req.Abilities,
		Story:       req.Story,
	}
	
	if err := a.characterStore.UpdateCharacter(character); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   "failed to update character",
		})
		return
	}
	
	writeJSON(w, http.StatusOK, map[string]any{
		"success":   true,
		"character": character,
	})
}

// handleDeleteUserCharacter deletes a character
func (a *App) handleDeleteUserCharacter(w http.ResponseWriter, r *http.Request) {
	// Require authentication
	sessionData := a.requireAuth(w, r)
	if sessionData == nil {
		return
	}
	
	characterID := chi.URLParam(r, "id")
	if characterID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "character ID required",
		})
		return
	}
	
	if err := a.characterStore.DeleteCharacter(characterID, sessionData.UserID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	
	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
	})
}

// handlePostNewCharacterMessage posts a message to newcharacter discussion
func (a *App) handlePostNewCharacterMessage(w http.ResponseWriter, r *http.Request) {
	// Require authentication
	sessionData := a.requireAuth(w, r)
	if sessionData == nil {
		return
	}
	
	var req struct {
		Section string `json:"section"`
		Content string `json:"content"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "invalid request body",
		})
		return
	}
	
	if req.Section == "" || req.Content == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "section and content are required",
		})
		return
	}
	
	// Save message to database
	messageID := generateID()
	query := `INSERT INTO newcharacter_messages (id, user_id, username, section, content) VALUES (?, ?, ?, ?, ?)`
	
	_, err := a.db.Exec(query, messageID, sessionData.UserID, sessionData.Username, req.Section, req.Content)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   "failed to save message",
		})
		return
	}
	
	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"message_id": messageID,
	})
}

// handleGetNewCharacterMessages retrieves messages for a section
func (a *App) handleGetNewCharacterMessages(w http.ResponseWriter, r *http.Request) {
	section := r.URL.Query().Get("section")
	if section == "" {
		section = "general"
	}
	
	query := `
		SELECT id, user_id, username, section, content, created_at
		FROM newcharacter_messages
		WHERE section = ?
		ORDER BY created_at ASC
	`
	
	rows, err := a.db.Query(query, section)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	defer rows.Close()
	
	type Message struct {
		ID        string `json:"id"`
		UserID    string `json:"user_id"`
		Username  string `json:"username"`
		Section   string `json:"section"`
		Content   string `json:"content"`
		CreatedAt string `json:"created_at"`
	}
	
	var messages []Message
	for rows.Next() {
		var msg Message
		err := rows.Scan(&msg.ID, &msg.UserID, &msg.Username, &msg.Section, &msg.Content, &msg.CreatedAt)
		if err != nil {
			continue
		}
		messages = append(messages, msg)
	}
	
	if messages == nil {
		messages = []Message{}
	}
	
	writeJSON(w, http.StatusOK, map[string]any{
		"success":  true,
		"messages": messages,
	})
}

