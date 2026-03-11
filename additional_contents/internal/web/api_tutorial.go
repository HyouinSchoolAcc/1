package web

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// RegisterTutorialAPI wires all tutorial-related endpoints
func (a *App) RegisterTutorialAPI(r chi.Router) {
	// Tutorial progress endpoints
	r.Get("/api/tutorial/progress", a.handleGetTutorialProgress)
	r.Post("/api/tutorial/start", a.handleStartTutorial)
	r.Post("/api/tutorial/update", a.handleUpdateTutorialProgress)
	r.Post("/api/tutorial/complete", a.handleCompleteTutorial)

	// Character certification endpoints
	r.Get("/api/certifications", a.handleGetCertifications)
	r.Get("/api/certifications/{characterId}", a.handleCheckCertification)
	r.Post("/api/certifications", a.handleCertifyUser)

	// Demo dialogues endpoint (public)
	r.Get("/api/demo-dialogues", a.handleGetDemoDialogues)
}

// handleGetTutorialProgress returns the current user's tutorial progress
func (a *App) handleGetTutorialProgress(w http.ResponseWriter, r *http.Request) {
	sessionData := a.getCurrentUser(r)
	if sessionData == nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"has_progress": false,
			"message":      "Not logged in",
		})
		return
	}

	progress, err := a.tutorialStore.GetTutorialProgress(sessionData.UserID)
	if err != nil {
		log.Printf("Error getting tutorial progress: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if progress == nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"has_progress": false,
			"user_id":      sessionData.UserID,
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"has_progress":    true,
		"progress":        progress,
		"user_id":         sessionData.UserID,
	})
}

// handleStartTutorial starts or retrieves the tutorial for the current user
func (a *App) handleStartTutorial(w http.ResponseWriter, r *http.Request) {
	sessionData := a.requireAuth(w, r)
	if sessionData == nil {
		return
	}

	progress, err := a.tutorialStore.StartTutorial(sessionData.UserID)
	if err != nil {
		log.Printf("Error starting tutorial: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success":  true,
		"progress": progress,
	})
}

// handleUpdateTutorialProgress updates the user's tutorial progress
func (a *App) handleUpdateTutorialProgress(w http.ResponseWriter, r *http.Request) {
	sessionData := a.requireAuth(w, r)
	if sessionData == nil {
		return
	}

	var req struct {
		CurrentStep    int      `json:"current_step"`
		CompletedSteps []string `json:"completed_steps"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := a.tutorialStore.UpdateTutorialProgress(sessionData.UserID, req.CurrentStep, req.CompletedSteps)
	if err != nil {
		log.Printf("Error updating tutorial progress: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
	})
}

// handleCompleteTutorial marks the tutorial as completed and upgrades user to writer
func (a *App) handleCompleteTutorial(w http.ResponseWriter, r *http.Request) {
	sessionData := a.requireAuth(w, r)
	if sessionData == nil {
		return
	}

	var req struct {
		CharacterID string `json:"character_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Allow empty body for backward compatibility
		req.CharacterID = ""
	}

	// Mark tutorial as complete
	if err := a.tutorialStore.CompleteTutorial(sessionData.UserID); err != nil {
		log.Printf("Error completing tutorial: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Upgrade user to writer role
	user, err := a.authService.userStore.GetUserByID(sessionData.UserID)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if user.Role == RoleNewUser || user.Role == RoleViewer {
		user.Role = RoleWriter
		if err := a.authService.userStore.UpdateUser(user); err != nil {
			log.Printf("Error updating user role: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Update session with new role
		sessionData.Role = RoleWriter
		if err := a.authService.sessionManager.SetSession(w, r, user); err != nil {
			log.Printf("Error updating session: %v", err)
		}
	}

	// Certify user for the character if provided
	if req.CharacterID != "" {
		if err := a.tutorialStore.CertifyUser(sessionData.UserID, req.CharacterID, ""); err != nil {
			log.Printf("Error certifying user: %v", err)
			// Don't fail the whole request for certification error
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success":  true,
		"new_role": string(RoleWriter),
		"message":  "Tutorial completed! You are now a writer.",
	})
}

// handleGetCertifications returns all character certifications for the current user
func (a *App) handleGetCertifications(w http.ResponseWriter, r *http.Request) {
	sessionData := a.getCurrentUser(r)
	if sessionData == nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"certifications": []any{},
		})
		return
	}

	certs, err := a.tutorialStore.GetUserCertifications(sessionData.UserID)
	if err != nil {
		log.Printf("Error getting certifications: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"certifications": certs,
		"user_id":        sessionData.UserID,
	})
}

// handleCheckCertification checks if the current user is certified for a specific character
func (a *App) handleCheckCertification(w http.ResponseWriter, r *http.Request) {
	characterID := chi.URLParam(r, "characterId")
	if characterID == "" {
		http.Error(w, "Character ID required", http.StatusBadRequest)
		return
	}

	sessionData := a.getCurrentUser(r)
	if sessionData == nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"certified": false,
			"reason":    "not_logged_in",
		})
		return
	}

	certified, err := a.tutorialStore.IsCertified(sessionData.UserID, characterID)
	if err != nil {
		log.Printf("Error checking certification: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"certified":    certified,
		"character_id": characterID,
		"user_id":      sessionData.UserID,
	})
}

// handleCertifyUser allows editors to certify users for characters
func (a *App) handleCertifyUser(w http.ResponseWriter, r *http.Request) {
	sessionData := a.requireEditor(w, r)
	if sessionData == nil {
		return
	}

	var req struct {
		UserID      string `json:"user_id"`
		CharacterID string `json:"character_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserID == "" || req.CharacterID == "" {
		http.Error(w, "user_id and character_id are required", http.StatusBadRequest)
		return
	}

	err := a.tutorialStore.CertifyUser(req.UserID, req.CharacterID, sessionData.UserID)
	if err != nil {
		log.Printf("Error certifying user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success":      true,
		"certified_by": sessionData.Username,
	})
}

// handleGetDemoDialogues returns sample dialogues for the landing page
func (a *App) handleGetDemoDialogues(w http.ResponseWriter, r *http.Request) {
	// For now, return curated sample dialogues from presets
	// In the future, this could read from the demo_dialogues table
	
	demoDialogues := []map[string]any{}

	// Try to load a sample from each preset
	for presetName := range a.presetFolders {
		// Skip non-CN presets for demo (or adjust as needed)
		if presetName == "presets_lin_lu_CN" {
			sample := a.loadSampleDialogue(presetName, 3) // Load up to 3 turns
			if sample != nil {
				demoDialogues = append(demoDialogues, sample)
			}
			break // Just one demo for now
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"dialogues": demoDialogues,
	})
}

// loadSampleDialogue loads a sample dialogue from a preset for demo purposes
func (a *App) loadSampleDialogue(preset string, maxTurns int) map[string]any {
	folder, err := a.getPresetFolder(preset)
	if err != nil {
		return nil
	}

	entries, _, err := a.listPresetFiles(preset)
	if err != nil || len(entries) == 0 {
		return nil
	}

	// Find the first valid JSON file
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if _, ok := parseWriterFilename(name); !ok {
			continue
		}

		path := folder + "/" + name
		var data map[string]any
		if err := a.loadJSON(path, &data); err != nil {
			continue
		}

		// Extract dialogue
		dialogueRaw, ok := data["dialogue"]
		if !ok {
			continue
		}
		dialogue, ok := dialogueRaw.([]any)
		if !ok || len(dialogue) == 0 {
			continue
		}

		// Truncate to maxTurns
		if len(dialogue) > maxTurns {
			dialogue = dialogue[:maxTurns]
		}

		characterName := getCurrentCharacterDisplayName(preset)

		return map[string]any{
			"character":    characterName,
			"character_id": preset,
			"dialogue":     dialogue,
			"source":       "demo",
		}
	}

	return nil
}

