package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
)

// RegisterAPI wires all JSON endpoints with 1:1 parity.
func (a *App) RegisterAPI(r chi.Router) {
	// Create auth middleware for API endpoints (returns JSON errors, not redirects)
	authMiddleware := NewAuthMiddleware(a.authService)
	writerAuth := authMiddleware.RequireWriterOrEditorAPI
	editorAuth := authMiddleware.RequireEditorAPI

	// User info endpoint (public - returns role info for frontend)
	r.Get("/api/current_user", a.handleGetCurrentUser)

	// User characters endpoints (these already have auth in their handlers)
	r.Get("/api/user-characters", a.handleGetUserCharacters)
	r.Post("/api/user-characters", a.handleCreateUserCharacter)
	r.Get("/api/user-characters/{id}", a.handleGetUserCharacter)
	r.Put("/api/user-characters/{id}", a.handleUpdateUserCharacter)
	r.Delete("/api/user-characters/{id}", a.handleDeleteUserCharacter)

	// New character discussion messages
	r.Get("/api/newcharacter-messages", a.handleGetNewCharacterMessages)
	r.Post("/api/newcharacter-messages", a.handlePostNewCharacterMessage)

	// Public read-only endpoints
	r.Get("/load_structured_writer_files", a.handleLoadStructuredWriterFiles)
	r.Get("/get_categories", a.handleGetCategories)
	r.Post("/load_writer_file_content", a.handleLoadWriterFileContent)
	r.Post("/load_checklist_data", a.handleLoadChecklistData)
	r.Post("/get_deletable_characters", a.handleGetDeletableCharacters)

	// Debug logging endpoint (accepts client-side logs)
	r.Post("/api/debug_log", a.handleDebugLog)

	// Guest-accessible endpoint for tutorial temp characters
	r.Post("/add_temp_character", a.handleAddTempCharacter)

	// Writer/Editor protected endpoints (data modification)
	r.Post("/save_writer_file_content", writerAuth(a.handleSaveWriterFileContent))
	r.Post("/update_character_turn", writerAuth(a.handleUpdateCharacterTurn))
	r.Post("/update_kurisu_turn", writerAuth(a.handleUpdateCharacterTurn)) // backward compatible
	r.Post("/save_conversation_data", writerAuth(a.handleSaveConversationData))
	r.Post("/create_new_version", writerAuth(a.handleCreateNewVersion))
	r.Post("/delete_version", writerAuth(a.handleDeleteVersion))
	r.Post("/add_character", writerAuth(a.handleAddCharacter))
	r.Post("/delete_character", writerAuth(a.handleDeleteCharacter))
	r.Post("/update_real_name", writerAuth(a.handleUpdateRealName))
	r.Post("/add_day", writerAuth(a.handleAddDay))
	r.Post("/delete_day", writerAuth(a.handleDeleteDay))
	r.Post("/update_schedule", writerAuth(a.handleUpdateSchedule))
	r.Post("/move_to_legacy", writerAuth(a.handleMoveToLegacy))
	r.Post("/auto_archive_old_items", writerAuth(a.handleAutoArchiveOldItems))
	r.Post("/update_inner_thought_annotation", writerAuth(a.handleUpdateInnerThoughtAnnotation))
	r.Post("/delete_inner_thought_annotation", writerAuth(a.handleDeleteInnerThoughtAnnotation))
	r.Post("/save_checklist_data", writerAuth(a.handleSaveChecklistData))

	// Editor-only endpoints (QC, approval, category changes)
	r.Post("/update_qc_status", editorAuth(a.handleUpdateQCStatus))
	r.Post("/update_day_category", editorAuth(a.handleUpdateDayCategory))

	// Character defaults API - provides initial values, experiences, etc. for each character
	r.Get("/api/character-defaults/{preset}", a.handleGetCharacterDefaults)
	r.Post("/api/character-defaults/{preset}", writerAuth(a.handleUpdateCharacterDefaults))
	r.Post("/api/sync-character-defaults", writerAuth(a.handleSyncCharacterDefaults))

	// Character profiles API - provides profile data (values, experiences, judgements, abilities)
	// from the centralized character_profiles.json file (source of truth for /descriptions page)
	r.Get("/api/character-profiles", a.handleGetAllCharacterProfiles)
	r.Get("/api/character-profiles/{characterId}", a.handleGetCharacterProfile)
	r.Get("/api/character-profiles/{characterId}/schedule/{day}", a.handleGetCharacterScheduleForDay)
	r.Post("/api/character-profiles/{characterId}/schedule", writerAuth(a.handleUpdateCharacterSchedule))
	r.Delete("/api/character-profiles/{characterId}/schedule/{day}", editorAuth(a.handleDeleteCharacterSchedule))

	// Progress tracking API - returns total passed tokens count
	r.Get("/api/passed-tokens", a.handleGetPassedTokens)

	// Storyboard regeneration API - runs the Python export script
	r.Post("/api/storyboard/regenerate", editorAuth(a.handleRegenerateStoryboard))

	// Sticker API endpoints
	r.Get("/api/stickers", a.handleGetStickers)
	r.Post("/api/stickers/validate", writerAuth(a.handleValidateSticker))
	r.Post("/api/stickers/upload", writerAuth(a.handleUploadSticker))
}

// handleGetCurrentUser returns the current user's session info for frontend role-based UI
func (a *App) handleGetCurrentUser(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	log.Printf("[PERF] handleGetCurrentUser: START")

	sessionStart := time.Now()
	sessionData := a.getCurrentUser(r)
	log.Printf("[PERF] handleGetCurrentUser: getCurrentUser took %v", time.Since(sessionStart))

	if sessionData == nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"logged_in":      false,
			"role":           string(RoleNewUser),
			"username":       "",
			"can_write":      false,
			"tutorial_done":  false,
			"certifications": []string{},
		})
		log.Printf("[PERF] handleGetCurrentUser: TOTAL %v (not logged in)", time.Since(startTime))
		return
	}

	// Get tutorial progress and certifications
	tutorialStart := time.Now()
	tutorialDone := false
	if progress, err := a.tutorialStore.GetTutorialProgress(sessionData.UserID); err == nil && progress != nil {
		tutorialDone = progress.IsCompleted
	}
	log.Printf("[PERF] handleGetCurrentUser: GetTutorialProgress took %v", time.Since(tutorialStart))

	certStart := time.Now()
	certifications := []string{}
	if certs, err := a.tutorialStore.GetUserCertifications(sessionData.UserID); err == nil {
		for _, cert := range certs {
			certifications = append(certifications, cert.CharacterID)
		}
	}
	log.Printf("[PERF] handleGetCurrentUser: GetUserCertifications took %v", time.Since(certStart))

	canWrite := sessionData.Role == RoleWriter || sessionData.Role == RoleEditor

	writeJSON(w, http.StatusOK, map[string]any{
		"logged_in":      true,
		"role":           string(sessionData.Role),
		"username":       sessionData.Username,
		"user_id":        sessionData.UserID,
		"can_write":      canWrite,
		"tutorial_done":  tutorialDone,
		"certifications": certifications,
	})
	log.Printf("[PERF] handleGetCurrentUser: TOTAL %v", time.Since(startTime))
}

// ===== route handlers =====

func (a *App) handleLoadStructuredWriterFiles(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	log.Printf("[PERF] handleLoadStructuredWriterFiles: START")

	// Get current user session (allow new users, but filter content accordingly)
	sessionData := a.getCurrentUser(r)
	userRole := RoleNewUser
	if sessionData != nil {
		userRole = sessionData.Role
	}

	preset := r.URL.Query().Get("preset_set")
	categoryFilter := r.URL.Query().Get("category")
	log.Printf("[PERF] handleLoadStructuredWriterFiles: preset=%s, category=%s (session check: %v)", preset, categoryFilter, time.Since(startTime))

	// TODO: Implement file filtering based on userRole
	// For now, new users see all files. Future enhancement will filter based on starred status
	_ = userRole // Suppress unused variable warning - will be used in future enhancement
	folder, err := a.getPresetFolder(preset)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}

	hiddenStart := time.Now()
	hidden := a.loadHiddenUsers(preset)
	log.Printf("[PERF] handleLoadStructuredWriterFiles: loadHiddenUsers took %v", time.Since(hiddenStart))

	dirStart := time.Now()
	entries, err := os.ReadDir(folder)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}
	log.Printf("[PERF] handleLoadStructuredWriterFiles: ReadDir found %d entries in %v", len(entries), time.Since(dirStart))

	structured := map[string]map[string]map[string]map[string]any{}

	jsonStart := time.Now()
	filesLoaded := 0
	filesSkipped := 0
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		pf, ok := parseWriterFilename(name)
		if !ok {
			filesSkipped++
			continue
		}
		if hidden != nil && hidden[pf.UserID] {
			filesSkipped++
			continue
		}

		path := filepath.Join(folder, name)
		var data map[string]any
		if err := a.loadJSON(path, &data); err != nil {
			filesSkipped++
			continue
		}
		filesLoaded++

		fileCategory := stringField(data, "category", "pending")
		if categoryFilter != "" && fileCategory != categoryFilter {
			continue
		}

		dialogue, _ := data["dialogue"].([]any)
		info := map[string]any{
			"filename":                    name,
			"completed":                   boolField(data, "completed", false),
			"day_num":                     pf.DayNum,
			"dialogue_count":              len(dialogue),
			"dialogue_char_count":         countDialogueHanChars(dialogue), // reuse Han counter for char count
			"dialogue_chinese_char_count": countDialogueHanChars(dialogue),
			"is_excellent_case":           a.excellent[preset][pf.UserID],
			"category":                    fileCategory,
			"rejection_reason":            stringField(data, "rejection_reason", ""),
			"ready_for_qc":                boolField(data, "ready_for_qc", false),
		}
		for _, k := range []string{"name", "user_name", "real_name", "dialogue_trait"} {
			if v, ok := data[k]; ok {
				info[k] = v
			}
		}

		if _, ok := structured[pf.UserID]; !ok {
			structured[pf.UserID] = map[string]map[string]map[string]any{}
		}
		if _, ok := structured[pf.UserID][pf.DupID]; !ok {
			structured[pf.UserID][pf.DupID] = map[string]map[string]any{}
		}
		structured[pf.UserID][pf.DupID][pf.DayStr] = info
	}
	log.Printf("[PERF] handleLoadStructuredWriterFiles: JSON loading took %v (loaded: %d, skipped: %d)", time.Since(jsonStart), filesLoaded, filesSkipped)

	// sort days
	sortStart := time.Now()
	for userID := range structured {
		for dup := range structured[userID] {
			dayMap := structured[userID][dup]
			keys := make([]string, 0, len(dayMap))
			for k := range dayMap {
				keys = append(keys, k)
			}
			sort.Slice(keys, func(i, j int) bool {
				di, _ := strconv.Atoi(strings.TrimPrefix(keys[i], "Day"))
				dj, _ := strconv.Atoi(strings.TrimPrefix(keys[j], "Day"))
				return di < dj
			})
			sorted := map[string]map[string]any{}
			for _, k := range keys {
				sorted[k] = dayMap[k]
			}
			structured[userID][dup] = sorted
		}
	}
	log.Printf("[PERF] handleLoadStructuredWriterFiles: sorting took %v", time.Since(sortStart))

	writeJSON(w, http.StatusOK, structured)
	log.Printf("[PERF] handleLoadStructuredWriterFiles: TOTAL %v (users: %d)", time.Since(startTime), len(structured))
}

func (a *App) handleGetCategories(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	log.Printf("[PERF] handleGetCategories: START")

	preset := r.URL.Query().Get("preset_set")
	folder, err := a.getPresetFolder(preset)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}
	cats := map[string]bool{}
	entries, _ := os.ReadDir(folder)
	log.Printf("[PERF] handleGetCategories: ReadDir found %d entries", len(entries))

	jsonStart := time.Now()
	filesLoaded := 0
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		pf, ok := parseWriterFilename(e.Name())
		if !ok {
			continue
		}
		path := filepath.Join(folder, pf.Name)
		var data map[string]any
		if err := a.loadJSON(path, &data); err != nil {
			cats["pending"] = true
			continue
		}
		filesLoaded++
		cats[stringField(data, "category", "existing")] = true
	}
	log.Printf("[PERF] handleGetCategories: loaded %d files in %v", filesLoaded, time.Since(jsonStart))

	out := []string{}
	for k := range cats {
		out = append(out, k)
	}
	sort.Strings(out)
	writeJSON(w, http.StatusOK, map[string]any{"categories": out})
	log.Printf("[PERF] handleGetCategories: TOTAL %v", time.Since(startTime))
}

func (a *App) handleLoadWriterFileContent(w http.ResponseWriter, r *http.Request) {
	// Allow all users (including new users) to load file content for viewing
	// File content filtering based on user role happens at the UI level
	loadStartTime := time.Now()

	body, err := decodeBody(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}
	filename := stringField(body, "filename", "")
	preset := stringField(body, "preset_set", "presets_kurisu")

	log.Printf("[LOAD] Loading file: %s (preset: %s)", filename, preset)

	folder, err := a.getPresetFolder(preset)
	if err != nil || filename == "" {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("Filename or preset set not provided")))
		return
	}
	pf, ok := parseWriterFilename(filename)
	if !ok {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("Invalid filename format")))
		return
	}
	path := filepath.Join(folder, filename)
	userData, err := a.loadUserDataWithSchedule(path, preset)
	if err != nil {
		log.Printf("[LOAD] ERROR loading %s: %v", filename, err)
		writeJSON(w, http.StatusInternalServerError, errJSON(err))
		return
	}

	// Log dialogue stats for debugging
	loadedTurns := 0
	loadedWordCount := 0
	if dlg, ok := userData["dialogue"].([]any); ok {
		loadedTurns = len(dlg)
		for _, turn := range dlg {
			if turnMap, ok := turn.(map[string]any); ok {
				if content, ok := turnMap["content"].(string); ok {
					words := len(strings.Fields(content))
					for _, r := range content {
						if r >= 0x4e00 && r <= 0x9fff {
							loadedWordCount++
						}
					}
					loadedWordCount += words
				}
			}
		}
	}
	log.Printf("[LOAD] Loaded %s: %d turns, ~%d words (took %v)", filename, loadedTurns, loadedWordCount, time.Since(loadStartTime))

	// recalc starting intimacy from previous day
	if pf.DayNum > 1 {
		prev := a.getPreviousDayIntimacy(pf, preset)
		old := intVal(userData["starting_intimacy_level"])
		userData["starting_intimacy_level"] = prev
		if val, ok := userData["intimacy_level"]; ok {
			userData["intimacy_level"] = intVal(val) - old + prev
		} else {
			userData["intimacy_level"] = prev
		}
		if old != prev {
			_ = a.atomicWriteJSON(path, userData)
		}
	} else {
		if _, ok := userData["starting_intimacy_level"]; !ok {
			userData["starting_intimacy_level"] = 0
		}
		if _, ok := userData["intimacy_level"]; !ok {
			userData["intimacy_level"] = userData["starting_intimacy_level"]
		}
	}

	// format schedules as strings for display
	if sched, ok := userData["user_schedule"].(map[string]any); ok {
		userData["user_schedule_raw"] = sched
		userData["user_schedule"] = formatScheduleToStringAny(sched, stringField(userData, "user_name", "Unknown User"), preset)
	}
	if sched, ok := userData["character_schedule"].(map[string]any); ok {
		userData["character_schedule_raw"] = sched
		userData["character_schedule"] = formatScheduleToStringAny(sched, stringField(userData, "character_name", "Unknown Character"), preset)
	}

	writeJSON(w, http.StatusOK, userData)
}

func (a *App) handleSaveWriterFileContent(w http.ResponseWriter, r *http.Request) {
	saveStartTime := time.Now()
	saveID := fmt.Sprintf("%d", saveStartTime.UnixNano()%1000000) // Short ID for log correlation

	// Only writers and editors can save content
	sessionData := a.requireWriterOrEditor(w, r)
	if sessionData == nil {
		return
	}

	body, err := decodeBody(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}
	filename := stringField(body, "filename", "")
	preset := stringField(body, "preset_set", "presets_kurisu")

	// Calculate incoming dialogue stats for logging
	incomingTurns := 0
	incomingWordCount := 0
	if dlg, ok := body["dialogue_list"].([]any); ok {
		incomingTurns = len(dlg)
		for _, turn := range dlg {
			if turnMap, ok := turn.(map[string]any); ok {
				if content, ok := turnMap["content"].(string); ok {
					// Count words (English) + characters (Chinese)
					words := len(strings.Fields(content))
					for _, r := range content {
						if r >= 0x4e00 && r <= 0x9fff {
							incomingWordCount++
						}
					}
					incomingWordCount += words
				}
			}
		}
	}

	log.Printf("[SAVE %s] === SAVE REQUEST RECEIVED ===", saveID)
	log.Printf("[SAVE %s] File: %s", saveID, filename)
	log.Printf("[SAVE %s] Preset: %s", saveID, preset)
	log.Printf("[SAVE %s] Incoming: %d turns, ~%d words", saveID, incomingTurns, incomingWordCount)
	log.Printf("[SAVE %s] Client IP: %s", saveID, r.RemoteAddr)

	folder, err := a.getPresetFolder(preset)
	if err != nil || filename == "" {
		log.Printf("[SAVE %s] ERROR: Invalid filename or preset", saveID)
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("Filename or preset set not provided")))
		return
	}
	path := filepath.Join(folder, filename)

	// Log existing file state before modification
	var fileData map[string]any
	if err := a.loadJSON(path, &fileData); err != nil {
		log.Printf("[SAVE %s] ERROR: File not found at %s", saveID, path)
		writeJSON(w, http.StatusNotFound, errJSON(errors.New("File not found, cannot save")))
		return
	}

	// Check ownership: users can only edit their own characters
	if !a.checkFileOwnershipFromData(w, sessionData, fileData) {
		log.Printf("[SAVE %s] DENIED: User %s not authorized to edit %s", saveID, sessionData.Username, filename)
		return
	}

	// Calculate existing dialogue stats
	existingTurns := 0
	existingWordCount := 0
	if existingDlg, ok := fileData["dialogue"].([]any); ok {
		existingTurns = len(existingDlg)
		for _, turn := range existingDlg {
			if turnMap, ok := turn.(map[string]any); ok {
				if content, ok := turnMap["content"].(string); ok {
					words := len(strings.Fields(content))
					for _, r := range content {
						if r >= 0x4e00 && r <= 0x9fff {
							existingWordCount++
						}
					}
					existingWordCount += words
				}
			}
		}
	}
	log.Printf("[SAVE %s] Existing file: %d turns, ~%d words", saveID, existingTurns, existingWordCount)

	// WARN if incoming has less content than existing (potential data loss)
	if incomingWordCount < existingWordCount && existingWordCount > 0 {
		log.Printf("[SAVE %s] ⚠️ WARNING: Incoming data (%d words) is SMALLER than existing (%d words)!",
			saveID, incomingWordCount, existingWordCount)
		log.Printf("[SAVE %s] ⚠️ This may indicate data loss - check client logs", saveID)
	}

	if dlg, ok := body["dialogue_list"]; ok {
		fileData["dialogue"] = dlg
	}
	fileData["completed"] = boolField(body, "is_complete", false)
	if v, ok := body["history"]; ok {
		fileData["history"] = v
	}
	// New profile fields
	if v, ok := body["values"]; ok {
		fileData["values"] = v
	}
	if v, ok := body["experiences"]; ok {
		fileData["experiences"] = v
	}
	if v, ok := body["judgements"]; ok {
		fileData["judgements"] = v
	}
	if v, ok := body["abilities"]; ok {
		fileData["abilities"] = v
	}
	if v, ok := body["intimacy_level"]; ok {
		fileData["intimacy_level"] = v
	}
	if v, ok := body["starting_intimacy_level"]; ok {
		fileData["starting_intimacy_level"] = v
	}
	if v, ok := body["user_schedule"]; ok {
		fileData["user_schedule"] = v
	}
	if v, ok := body["character_schedule"]; ok {
		day := intVal(getNested(v.(map[string]any), "day"))
		if day == 0 {
			if pf, ok := parseWriterFilename(filename); ok {
				day = pf.DayNum
			} else {
				day = 1
			}
		}
		_ = a.updateUniversalCharacterSchedule(day, v.(map[string]any), preset)
	}

	for _, k := range []string{"character_name", "user_name", "real_name", "relationship", "character_description", "character_motivation", "comments"} {
		if v, ok := body[k]; ok {
			fileData[k] = v
		}
	}

	if err := a.atomicWriteJSON(path, fileData); err != nil {
		log.Printf("[SAVE %s] ERROR: Write failed: %v", saveID, err)
		writeJSON(w, http.StatusInternalServerError, errJSON(err))
		return
	}

	elapsed := time.Since(saveStartTime)
	log.Printf("[SAVE %s] ✅ SUCCESS - Saved %d turns, ~%d words in %v", saveID, incomingTurns, incomingWordCount, elapsed)

	// update user_info if names/description changed
	a.updateUserInfoFromFile(body, fileData, preset)

	writeJSON(w, http.StatusOK, map[string]any{"success": true, "message": "File saved."})
}

func (a *App) handleUpdateQCStatus(w http.ResponseWriter, r *http.Request) {
	// Only writers and editors can update QC status
	sessionData := a.requireWriterOrEditor(w, r)
	if sessionData == nil {
		return
	}
	body, err := decodeBody(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}
	filename := stringField(body, "filename", "")
	preset := stringField(body, "preset_set", "presets_kurisu")
	ready := boolField(body, "ready_for_qc", false)
	folder, err := a.getPresetFolder(preset)
	if err != nil || filename == "" {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("Filename or preset set not provided")))
		return
	}
	path := filepath.Join(folder, filename)
	var data map[string]any
	if err := a.loadJSON(path, &data); err != nil {
		writeJSON(w, http.StatusNotFound, errJSON(errors.New("File not found")))
		return
	}

	// Check ownership: users can only update QC status for their own characters
	if !a.checkFileOwnershipFromData(w, sessionData, data) {
		return
	}

	data["ready_for_qc"] = ready
	if err := a.atomicWriteJSON(path, data); err != nil {
		writeJSON(w, http.StatusInternalServerError, errJSON(err))
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"success": true, "message": "QC status updated successfully.", "ready_for_qc": ready})
}

// conversation buffers
var savingChangedList []map[string]any
var savingUnchangedList []map[string]any
var convMu sync.Mutex

func (a *App) handleUpdateCharacterTurn(w http.ResponseWriter, r *http.Request) {
	// Only writers and editors can update character turns
	if sessionData := a.requireWriterOrEditor(w, r); sessionData == nil {
		return
	}
	body, err := decodeBody(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}
	idx := numberField(body, "scl_index", -1)
	newContent := stringField(body, "content", "")
	newChoices, _ := body["choices"].(map[string]any)
	preset := stringField(body, "preset_set", "presets_kurisu")
	role := getCurrentCharacterRole(preset)

	convMu.Lock()
	defer convMu.Unlock()
	if idx < 0 || idx >= len(savingChangedList) {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("Invalid index")))
		return
	}
	target := savingChangedList[idx]
	if stringField(target, "role", "") != role {
		writeJSON(w, http.StatusBadRequest, errJSON(fmt.Errorf("Can only update %s turns", role)))
		return
	}
	changed := false
	if newContent != "" && stringField(target, "content", "") != newContent {
		target["content"] = newContent
		changed = true
	}
	if newChoices != nil {
		if targetChoices, ok := target["choices"].(map[string]any); ok {
			target["choices"] = targetChoices
		} else {
			target["choices"] = map[string]any{}
		}
		tc := target["choices"].(map[string]any)
		if v, ok := newChoices["inner_thoughts_all"]; ok {
			tc["inner_thoughts_all"] = v
		}
		if v, ok := newChoices["inner_thoughts"]; ok {
			if fmt.Sprint(tc["inner_thoughts"]) != fmt.Sprint(v) {
				tc["inner_thoughts"] = v
				changed = true
			}
		}
	}
	if changed {
		target["dpo"] = true
	}
	writeJSON(w, http.StatusOK, map[string]any{"success": true, "updated_turn": target, "scl_index": idx})
}

func (a *App) handleSaveConversationData(w http.ResponseWriter, r *http.Request) {
	// Only writers and editors can save conversation data
	if sessionData := a.requireWriterOrEditor(w, r); sessionData == nil {
		return
	}
	body, err := decodeBody(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}
	prefix := strings.TrimSpace(stringField(body, "prefix", "chatlog"))
	if prefix == "" {
		prefix = "chatlog"
	}
	convMu.Lock()
	defer convMu.Unlock()
	files := map[string][]byte{}
	if b, err := json.MarshalIndent(savingChangedList, "", "  "); err == nil {
		files[prefix+"_changed.json"] = b
	}
	if b, err := json.MarshalIndent(savingUnchangedList, "", "  "); err == nil {
		files[prefix+"_unchanged.json"] = b
	}
	zipBytes, err := zipBuffers(files)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errJSON(err))
		return
	}
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.zip\"", prefix))
	w.Write(zipBytes)
}

func (a *App) handleCreateNewVersion(w http.ResponseWriter, r *http.Request) {
	sessionData := a.getCurrentUser(r)

	body, err := decodeBody(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}
	source := stringField(body, "source_filename", "")
	preset := stringField(body, "preset_set", "presets_kurisu")
	folder, err := a.getPresetFolder(preset)
	if err != nil || source == "" {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("Source filename or preset set not provided")))
		return
	}
	pf, ok := parseWriterFilename(source)
	if !ok {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("Invalid source filename format")))
		return
	}

	// Check ownership: users can only create versions for their own characters
	if !a.checkPresetCharacterOwnership(w, sessionData, pf.UserID, preset) {
		return
	}
	entries, _ := os.ReadDir(folder)
	var existingDup []int
	var templateFiles []struct {
		name string
		pf   *ParsedFilename
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		p, ok := parseWriterFilename(e.Name())
		if !ok || p.UserID != pf.UserID {
			continue
		}
		if strings.HasPrefix(p.DupID, "dup_") {
			if n, err := strconv.Atoi(strings.TrimPrefix(p.DupID, "dup_")); err == nil {
				existingDup = append(existingDup, n)
			}
		}
		if p.DupID == "dup_0" {
			templateFiles = append(templateFiles, struct {
				name string
				pf   *ParsedFilename
			}{e.Name(), p})
		}
	}
	nextDup := 0
	if len(existingDup) > 0 {
		sort.Ints(existingDup)
		nextDup = existingDup[len(existingDup)-1] + 1
	}
	if len(templateFiles) == 0 && len(existingDup) > 0 {
		lowest := existingDup[0]
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			p, ok := parseWriterFilename(e.Name())
			if ok && p.UserID == pf.UserID && p.DupID == fmt.Sprintf("dup_%d", lowest) {
				templateFiles = append(templateFiles, struct {
					name string
					pf   *ParsedFilename
				}{e.Name(), p})
			}
		}
	}
	if len(templateFiles) == 0 {
		writeJSON(w, http.StatusNotFound, errJSON(fmt.Errorf("No template files found for %s", pf.UserID)))
		return
	}
	created := []string{}
	for _, t := range templateFiles {
		srcPath := filepath.Join(folder, t.name)
		var data map[string]any
		if err := a.loadJSON(srcPath, &data); err != nil {
			continue
		}
		day := t.pf.DayNum
		newName := fmt.Sprintf("%s_Day%d_dup_%d_simplified.json", pf.UserID, day, nextDup)
		// Get default values for the character from defaults file
		defaults := a.getCharacterDefaultsFromFile(preset)
		// Copy owner_user_id from source file's metadata
		newCreationMeta := map[string]any{
			"created_at":            time.Now().Unix(),
			"created_by_ip":         r.RemoteAddr,
			"created_by_user_agent": r.UserAgent(),
			"version":               "1.0",
		}
		if srcMeta, ok := data["_creation_metadata"].(map[string]any); ok {
			if ownerID, ok := srcMeta["owner_user_id"].(string); ok && ownerID != "" {
				newCreationMeta["owner_user_id"] = ownerID
			}
			if ownerUsername, ok := srcMeta["owner_username"].(string); ok && ownerUsername != "" {
				newCreationMeta["owner_username"] = ownerUsername
			}
		}
		newData := map[string]any{
			"dialogue":                []any{},
			"completed":               false,
			"history":                 "",
			"user_schedule":           data["user_schedule"],
			"relationship":            data["relationship"],
			"user_name":               data["user_name"],
			"character_name":          data["character_name"],
			"user_id":                 data["user_id"],
			"real_name":               data["real_name"],
			"intimacy_level":          0,
			"starting_intimacy_level": 0,
			"values":                  defaults.Values,
			"experiences":             defaults.Experiences,
			"judgements":              defaults.Judgements,
			"abilities":               defaults.Abilities,
			"_creation_metadata":      newCreationMeta,
		}
		if err := a.atomicWriteJSON(filepath.Join(folder, newName), newData); err == nil {
			created = append(created, newName)
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"success":       true,
		"created_files": created,
		"message":       fmt.Sprintf("New version %d created with %d days copied", nextDup, len(created)),
	})
}

func (a *App) handleDeleteVersion(w http.ResponseWriter, r *http.Request) {
	sessionData := a.getCurrentUser(r)

	body, err := decodeBody(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}
	userID := stringField(body, "user_id", "")
	dupID := stringField(body, "dup_id", "")
	preset := stringField(body, "preset_set", "presets_kurisu")
	folder, err := a.getPresetFolder(preset)
	if err != nil || userID == "" || dupID == "" {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("User ID, dup ID, or preset set not provided")))
		return
	}

	// Check ownership: users can only delete versions of their own characters
	if !a.checkPresetCharacterOwnership(w, sessionData, userID, preset) {
		return
	}
	entries, _ := os.ReadDir(folder)
	var toDelete []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		pf, ok := parseWriterFilename(e.Name())
		if ok && pf.UserID == userID && pf.DupID == dupID {
			toDelete = append(toDelete, e.Name())
		}
	}
	if len(toDelete) == 0 {
		writeJSON(w, http.StatusNotFound, errJSON(fmt.Errorf("No files found for %s with %s", userID, dupID)))
		return
	}
	deleted := []string{}
	for _, name := range toDelete {
		if err := os.Remove(filepath.Join(folder, name)); err == nil {
			deleted = append(deleted, name)
		}
	}
	if len(deleted) == 0 {
		writeJSON(w, http.StatusInternalServerError, errJSON(errors.New("Failed to delete any files")))
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"success":       true,
		"deleted_files": deleted,
		"message":       fmt.Sprintf("Successfully deleted %d file(s) for %s %s", len(deleted), userID, dupID),
	})
}

func (a *App) handleAddCharacter(w http.ResponseWriter, r *http.Request) {
	// Get session data for ownership tracking
	sessionData := a.getCurrentUser(r)

	body, err := decodeBody(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}
	required := []string{"character_name", "character_description", "character_motivation", "character_age", "character_profession", "character_real_name"}
	for _, k := range required {
		if stringField(body, k, "") == "" {
			writeJSON(w, http.StatusBadRequest, errJSON(errors.New("角色姓名、真实姓名、描述、动机、年龄和职业为必填项。")))
			return
		}
	}
	day1Schedule, ok := body["day1_schedule"].(map[string]any)
	if !ok || len(day1Schedule) == 0 {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("第1天的日程至少需要包含一个时间段。")))
		return
	}
	preset := stringField(body, "preset_set", "presets_kurisu")
	folder, err := a.getPresetFolder(preset)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}

	// find next user num
	entries, _ := os.ReadDir(folder)
	userNums := map[int]bool{}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if pf, ok := parseWriterFilename(e.Name()); ok {
			n, _ := strconv.Atoi(strings.TrimPrefix(pf.UserID, "user_"))
			userNums[n] = true
		}
	}
	nextUser := 0
	for userNums[nextUser] {
		nextUser++
	}
	userID := fmt.Sprintf("user_%d", nextUser)

	characterName := stringField(body, "character_name", "")
	characterReal := stringField(body, "character_real_name", "")
	characterDesc := stringField(body, "character_description", "")
	characterMot := stringField(body, "character_motivation", "")
	characterAge := stringField(body, "character_age", "")
	characterProf := stringField(body, "character_profession", "")

	creationMeta := map[string]any{
		"created_at":            time.Now().Unix(),
		"created_by_ip":         r.RemoteAddr,
		"created_by_user_agent": r.UserAgent(),
		"version":               "1.0",
	}
	// Store owner_user_id to link preset character to the logged-in user account
	if sessionData != nil {
		creationMeta["owner_user_id"] = sessionData.UserID
		creationMeta["owner_username"] = sessionData.Username
	}

	// save user info file
	a.addUserInfoEntry(folder, preset, userID, nextUser, characterName, characterReal, characterAge, characterProf, characterMot, characterDesc, creationMeta)

	// create day1 schedule dict (english/chinese keys)
	userSchedule := map[string]any{"day": 1}
	timePeriods := []string{"morning", "noon", "afternoon", "evening", "night"}
	keyMapping := map[string]string{}
	if strings.Contains(preset, "_CN") {
		keyMapping = map[string]string{
			"morning":   "早晨",
			"noon":      "中午",
			"afternoon": "下午",
			"evening":   "晚上",
			"night":     "夜晚",
		}
	} else {
		for _, p := range timePeriods {
			keyMapping[p] = p
		}
	}
	for _, p := range timePeriods {
		if v, ok := day1Schedule["Day1_"+p]; ok {
			if s, ok := v.(string); ok && strings.TrimSpace(s) != "" {
				userSchedule[keyMapping[p]] = strings.TrimSpace(s)
			}
		} else if v, ok := day1Schedule[p]; ok {
			if s, ok := v.(string); ok && strings.TrimSpace(s) != "" {
				userSchedule[keyMapping[p]] = strings.TrimSpace(s)
			}
		}
	}

	relationshipDesc := "Two people have just met and are still unfamiliar with each other. Their conversations are mostly polite and superficial, limited to everyday topics like work, studies, or the weather. At this stage, both sides are cautious and have yet to reveal their true selves."
	aiCharacterName := getCurrentCharacterDisplayName(preset)

	// Get default character values (values, experiences, judgements, abilities) from defaults file
	defaults := a.getCharacterDefaultsFromFile(preset)

	fileData := map[string]any{
		"user_schedule":           userSchedule,
		"history":                 "",
		"relationship":            relationshipDesc,
		"dialogue":                []any{},
		"intimacy_level":          0,
		"starting_intimacy_level": 0,
		"completed":               false,
		"user_id":                 userID,
		"user_name":               characterName,
		"real_name":               characterReal,
		"character_name":          aiCharacterName,
		"values":                  defaults.Values,
		"experiences":             defaults.Experiences,
		"judgements":              defaults.Judgements,
		"abilities":               defaults.Abilities,
		"_creation_metadata":      creationMeta,
	}
	filename := fmt.Sprintf("%s_Day1_dup_1_simplified.json", userID)
	if err := a.atomicWriteJSON(filepath.Join(folder, filename), fileData); err != nil {
		writeJSON(w, http.StatusInternalServerError, errJSON(err))
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success":        true,
		"user_id":        userID,
		"filename":       filename,
		"character_name": characterName,
		"message":        fmt.Sprintf("Character '%s' created successfully as %s", characterName, userID),
	})
}

// handleAddTempCharacter creates a temporary character for tutorial/guest users
// These characters are prefixed with "temp_" and can be created without authentication
func (a *App) handleAddTempCharacter(w http.ResponseWriter, r *http.Request) {
	body, err := decodeBody(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}
	
	// Less strict validation for temp characters - only require name
	characterName := stringField(body, "character_name", "")
	if characterName == "" {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("角色姓名为必填项。")))
		return
	}
	
	// Get other fields with defaults
	characterReal := stringField(body, "character_real_name", "游客")
	characterDesc := stringField(body, "character_description", "临时角色")
	characterMot := stringField(body, "character_motivation", "体验写作")
	characterAge := stringField(body, "character_age", "25")
	characterProf := stringField(body, "character_profession", "学生")
	
	// Day1 schedule - use defaults if not provided
	day1Schedule, ok := body["day1_schedule"].(map[string]any)
	if !ok || len(day1Schedule) == 0 {
		day1Schedule = map[string]any{
			"morning": "起床，准备开始新的一天",
		}
	}
	
	preset := stringField(body, "preset_set", "presets_kurisu")
	folder, err := a.getPresetFolder(preset)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}

	// Find next temp user number (use temp_ prefix)
	entries, _ := os.ReadDir(folder)
	tempNums := map[int]bool{}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if pf, ok := parseWriterFilename(e.Name()); ok && strings.HasPrefix(pf.UserID, "temp_") {
			n, _ := strconv.Atoi(strings.TrimPrefix(pf.UserID, "temp_"))
			tempNums[n] = true
		}
	}
	nextTemp := 0
	for tempNums[nextTemp] {
		nextTemp++
	}
	userID := fmt.Sprintf("temp_%d", nextTemp)

	creationMeta := map[string]any{
		"created_at":            time.Now().Unix(),
		"created_by_ip":         r.RemoteAddr,
		"created_by_user_agent": r.UserAgent(),
		"version":               "1.0",
		"is_temporary":          true,
	}
	// Store owner_user_id if the user is logged in
	if sessionData := a.getCurrentUser(r); sessionData != nil {
		creationMeta["owner_user_id"] = sessionData.UserID
		creationMeta["owner_username"] = sessionData.Username
	}

	// save user info file
	a.addUserInfoEntry(folder, preset, userID, nextTemp, characterName, characterReal, characterAge, characterProf, characterMot, characterDesc, creationMeta)

	// create day1 schedule dict
	userSchedule := map[string]any{"day": 1}
	timePeriods := []string{"morning", "noon", "afternoon", "evening", "night"}
	keyMapping := map[string]string{}
	if strings.Contains(preset, "_CN") {
		keyMapping = map[string]string{
			"morning":   "早晨",
			"noon":      "中午",
			"afternoon": "下午",
			"evening":   "晚上",
			"night":     "夜晚",
		}
	} else {
		for _, p := range timePeriods {
			keyMapping[p] = p
		}
	}
	for _, p := range timePeriods {
		if v, ok := day1Schedule["Day1_"+p]; ok {
			if s, ok := v.(string); ok && strings.TrimSpace(s) != "" {
				userSchedule[keyMapping[p]] = strings.TrimSpace(s)
			}
		} else if v, ok := day1Schedule[p]; ok {
			if s, ok := v.(string); ok && strings.TrimSpace(s) != "" {
				userSchedule[keyMapping[p]] = strings.TrimSpace(s)
			}
		}
	}

	relationshipDesc := "Two people have just met and are still unfamiliar with each other. Their conversations are mostly polite and superficial, limited to everyday topics like work, studies, or the weather. At this stage, both sides are cautious and have yet to reveal their true selves."
	aiCharacterName := getCurrentCharacterDisplayName(preset)

	// Get default character values
	defaults := a.getCharacterDefaultsFromFile(preset)

	fileData := map[string]any{
		"user_schedule":           userSchedule,
		"history":                 "",
		"relationship":            relationshipDesc,
		"dialogue":                []any{},
		"intimacy_level":          0,
		"starting_intimacy_level": 0,
		"completed":               false,
		"user_id":                 userID,
		"user_name":               characterName,
		"real_name":               characterReal,
		"character_name":          aiCharacterName,
		"values":                  defaults.Values,
		"experiences":             defaults.Experiences,
		"judgements":              defaults.Judgements,
		"abilities":               defaults.Abilities,
		"_creation_metadata":      creationMeta,
	}
	filename := fmt.Sprintf("%s_Day1_dup_1_simplified.json", userID)
	if err := a.atomicWriteJSON(filepath.Join(folder, filename), fileData); err != nil {
		writeJSON(w, http.StatusInternalServerError, errJSON(err))
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success":        true,
		"user_id":        userID,
		"filename":       filename,
		"character_name": characterName,
		"is_temporary":   true,
		"message":        fmt.Sprintf("Temporary character '%s' created successfully as %s", characterName, userID),
	})
}

func (a *App) handleDeleteCharacter(w http.ResponseWriter, r *http.Request) {
	sessionData := a.getCurrentUser(r)

	body, err := decodeBody(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}
	userID := stringField(body, "user_id", "")
	preset := stringField(body, "preset_set", "presets_kurisu")
	characterName := stringField(body, "character_name", "")
	if userID == "" {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("User ID is required")))
		return
	}

	// Check ownership: users can only delete their own characters
	if !a.checkPresetCharacterOwnership(w, sessionData, userID, preset) {
		return
	}

	folder, err := a.getPresetFolder(preset)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}
	entries, _ := os.ReadDir(folder)
	var files []string
	var creationMeta map[string]any
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		pf, ok := parseWriterFilename(e.Name())
		if ok && pf.UserID == userID {
			files = append(files, e.Name())
			if creationMeta == nil {
				var data map[string]any
				_ = a.loadJSON(filepath.Join(folder, e.Name()), &data)
				if meta, ok := data["_creation_metadata"].(map[string]any); ok {
					creationMeta = meta
				}
			}
		}
	}
	if len(files) == 0 {
		writeJSON(w, http.StatusNotFound, errJSON(fmt.Errorf("No files found for user %s", userID)))
		return
	}
	if creationMeta == nil {
		writeJSON(w, http.StatusForbidden, errJSON(errors.New("Cannot delete character: Creation metadata not found.")))
		return
	}
	createdAt := int64(intVal(creationMeta["created_at"]))
	createdIP := stringField(creationMeta, "created_by_ip", "")
	if time.Since(time.Unix(createdAt, 0)) > 24*time.Hour {
		writeJSON(w, http.StatusForbidden, errJSON(fmt.Errorf("Cannot delete character: This character was created %.1f hours ago. Characters can only be deleted within 24 hours of creation.", time.Since(time.Unix(createdAt, 0)).Hours())))
		return
	}
	if createdIP != "" && !strings.HasPrefix(createdIP, r.RemoteAddr) {
		writeJSON(w, http.StatusForbidden, errJSON(fmt.Errorf("Cannot delete character: This character was created by a different user (IP: %s...). You can only delete characters you created yourself.", createdIP[:min(8, len(createdIP))])))
		return
	}
	deleted := []string{}
	for _, name := range files {
		_ = os.Remove(filepath.Join(folder, name))
		deleted = append(deleted, name)
	}
	// update user_info file
	a.removeUserInfoEntry(folder, preset, userID)
	// remove from user_info_CN.js (best-effort)
	if characterName != "" {
		a.removeUserInfoCNLine(preset, characterName)
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"success":           true,
		"deleted_files":     len(deleted),
		"deleted_file_list": deleted,
		"message":           fmt.Sprintf("Successfully deleted %d files for user %s", len(deleted), userID),
	})
}

func (a *App) handleGetDeletableCharacters(w http.ResponseWriter, r *http.Request) {
	body, err := decodeBody(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}
	preset := stringField(body, "preset_set", "presets_kurisu")
	folder, err := a.getPresetFolder(preset)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}
	currentIP := r.RemoteAddr
	entries, _ := os.ReadDir(folder)
	userFiles := map[string][]string{}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if pf, ok := parseWriterFilename(e.Name()); ok {
			userFiles[pf.UserID] = append(userFiles[pf.UserID], e.Name())
		}
	}
	deletable := []map[string]any{}
	for user, files := range userFiles {
		var data map[string]any
		if err := a.loadJSON(filepath.Join(folder, files[0]), &data); err != nil {
			continue
		}
		meta, _ := data["_creation_metadata"].(map[string]any)
		if meta == nil {
			continue
		}
		createdAt := intVal(meta["created_at"])
		createdIP := stringField(meta, "created_by_ip", "")
		if createdAt == 0 || createdIP == "" {
			continue
		}
		age := time.Since(time.Unix(int64(createdAt), 0))
		if age <= 24*time.Hour && strings.HasPrefix(currentIP, createdIP) {
			deletable = append(deletable, map[string]any{
				"user_id":                 user,
				"character_name":          stringField(data, "user_name", "Unknown Character"),
				"file_count":              len(files),
				"hours_since_creation":    round1(age.Hours()),
				"can_delete":              true,
				"created_by_current_user": true,
			})
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"success":              true,
		"deletable_characters": deletable,
		"message":              fmt.Sprintf("Found %d deletable characters", len(deletable)),
	})
}

func round1(f float64) float64 {
	return math.Round(f*10) / 10
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (a *App) handleUpdateRealName(w http.ResponseWriter, r *http.Request) {
	sessionData := a.getCurrentUser(r)

	body, err := decodeBody(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}
	userID := stringField(body, "user_id", "")
	realName := stringField(body, "real_name", "")
	preset := stringField(body, "preset_set", "presets_kurisu")
	if userID == "" || realName == "" {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("User ID and real name are required")))
		return
	}

	// Check ownership: users can only update real name for their own characters
	if !a.checkPresetCharacterOwnership(w, sessionData, userID, preset) {
		return
	}

	folder, err := a.getPresetFolder(preset)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}
	entries, _ := os.ReadDir(folder)
	updated := 0
	var characterName string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		pf, ok := parseWriterFilename(e.Name())
		if ok && pf.UserID == userID {
			var data map[string]any
			if err := a.loadJSON(filepath.Join(folder, e.Name()), &data); err != nil {
				continue
			}
			data["real_name"] = realName
			if characterName == "" {
				characterName = stringField(data, "user_name", "")
			}
			if err := a.atomicWriteJSON(filepath.Join(folder, e.Name()), data); err == nil {
				updated++
			}
		}
	}
	a.updateRealNameInUserInfo(folder, preset, userID, realName, characterName)
	writeJSON(w, http.StatusOK, map[string]any{
		"success":       true,
		"message":       fmt.Sprintf("Real name updated successfully for %d files", updated),
		"updated_files": updated,
	})
}

func (a *App) handleAddDay(w http.ResponseWriter, r *http.Request) {
	sessionData := a.getCurrentUser(r)

	body, err := decodeBody(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}
	userID := stringField(body, "user_id", "")
	dayNumber := numberField(body, "day_number", 0)
	daySchedule, _ := body["day_schedule"].(map[string]any)
	preset := stringField(body, "preset_set", "presets_kurisu")
	if userID == "" || dayNumber < 1 || len(daySchedule) == 0 {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("User ID, valid day number and schedule required")))
		return
	}

	// Check ownership: users can only add days to their own characters
	if !a.checkPresetCharacterOwnership(w, sessionData, userID, preset) {
		return
	}

	folder, err := a.getPresetFolder(preset)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}
	entries, _ := os.ReadDir(folder)
	var existing []struct {
		name string
		pf   *ParsedFilename
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if pf, ok := parseWriterFilename(e.Name()); ok && pf.UserID == userID {
			existing = append(existing, struct {
				name string
				pf   *ParsedFilename
			}{e.Name(), pf})
		}
	}
	if len(existing) == 0 {
		writeJSON(w, http.StatusNotFound, errJSON(fmt.Errorf("Character %s not found. Create the character first.", userID)))
		return
	}
	for _, ex := range existing {
		if ex.pf.DayNum == dayNumber {
			writeJSON(w, http.StatusConflict, errJSON(fmt.Errorf("Day %d already exists for %s. Files: %s", dayNumber, userID, ex.name)))
			return
		}
	}
	firstFile := filepath.Join(folder, existing[0].name)
	var firstData map[string]any
	_ = a.loadJSON(firstFile, &firstData)
	characterName := stringField(firstData, "user_name", "Unknown Character")
	aiName := getCurrentCharacterDisplayName(preset)

	// starting intimacy from previous day
	startIntimacy := 0
	var values, experiences, judgements, abilities interface{}

	for _, ex := range existing {
		if ex.pf.DayNum == dayNumber-1 {
			var prev map[string]any
			if err := a.loadJSON(filepath.Join(folder, ex.name), &prev); err == nil {
				startIntimacy = intVal(prev["intimacy_level"])
				values = prev["values"]
				experiences = prev["experiences"]
				judgements = prev["judgements"]
				abilities = prev["abilities"]
			}
			break
		}
	}

	// If no previous day values found, use character defaults from defaults file
	if values == nil || values == "" {
		defaults := a.getCharacterDefaultsFromFile(preset)
		values = defaults.Values
		experiences = defaults.Experiences
		judgements = defaults.Judgements
		abilities = defaults.Abilities
	}

	dupID := existing[0].pf.DupID
	filename := fmt.Sprintf("%s_Day%d_%s_simplified.json", userID, dayNumber, dupID)

	// Fetch character schedule from character_profiles.json for this day
	characterSchedule := a.getCharacterScheduleForDay(preset, dayNumber)

	// Copy owner_user_id from the first existing file's metadata
	newCreationMeta := map[string]any{
		"created_at":            time.Now().Unix(),
		"created_by_ip":         r.RemoteAddr,
		"created_by_user_agent": r.UserAgent(),
		"version":               "1.0",
	}
	if srcMeta, ok := firstData["_creation_metadata"].(map[string]any); ok {
		if ownerID, ok := srcMeta["owner_user_id"].(string); ok && ownerID != "" {
			newCreationMeta["owner_user_id"] = ownerID
		}
		if ownerUsername, ok := srcMeta["owner_username"].(string); ok && ownerUsername != "" {
			newCreationMeta["owner_username"] = ownerUsername
		}
	}

	fileData := map[string]any{
		"user_schedule":           map[string]any{"day": dayNumber},
		"character_schedule":      characterSchedule,
		"history":                 "",
		"values":                  values,
		"experiences":             experiences,
		"judgements":              judgements,
		"abilities":               abilities,
		"relationship":            "Continuing relationship from previous interactions.",
		"dialogue":                []any{},
		"intimacy_level":          startIntimacy,
		"starting_intimacy_level": startIntimacy,
		"completed":               false,
		"user_name":               characterName,
		"character_name":          aiName,
		"_creation_metadata":      newCreationMeta,
	}
	for k, v := range daySchedule {
		if s, ok := v.(string); ok && strings.TrimSpace(s) != "" {
			fileData["user_schedule"].(map[string]any)[k] = strings.TrimSpace(s)
		}
	}
	if err := a.atomicWriteJSON(filepath.Join(folder, filename), fileData); err != nil {
		writeJSON(w, http.StatusInternalServerError, errJSON(err))
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"success":                 true,
		"filename":                filename,
		"character_name":          characterName,
		"day_number":              dayNumber,
		"starting_intimacy_level": startIntimacy,
		"message":                 fmt.Sprintf("Day %d created successfully for %s (%s)", dayNumber, characterName, userID),
	})
}

func (a *App) handleDeleteDay(w http.ResponseWriter, r *http.Request) {
	sessionData := a.getCurrentUser(r)

	body, err := decodeBody(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}
	userID := stringField(body, "user_id", "")
	dayNumber := numberField(body, "day_number", 0)
	dupID := stringField(body, "dup_id", "dup_1")
	preset := stringField(body, "preset_set", "presets_kurisu")
	if userID == "" || dayNumber < 1 {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("User ID and valid day number are required")))
		return
	}

	// Check ownership: users can only delete days from their own characters
	if !a.checkPresetCharacterOwnership(w, sessionData, userID, preset) {
		return
	}

	folder, err := a.getPresetFolder(preset)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}
	filename := fmt.Sprintf("%s_Day%d_%s_simplified.json", userID, dayNumber, dupID)
	path := filepath.Join(folder, filename)
	if _, err := os.Stat(path); err != nil {
		writeJSON(w, http.StatusNotFound, errJSON(fmt.Errorf("Day %d file not found for %s", dayNumber, userID)))
		return
	}
	// check later days
	var later []int
	entries, _ := os.ReadDir(folder)
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		pf, ok := parseWriterFilename(e.Name())
		if ok && pf.UserID == userID && pf.DupID == dupID && pf.DayNum > dayNumber {
			later = append(later, pf.DayNum)
		}
	}
	var characterName string
	if b, err := os.ReadFile(path); err == nil {
		var data map[string]any
		if err := json.Unmarshal(b, &data); err == nil {
			characterName = stringField(data, "user_name", "Unknown Character")
		}
	}
	_ = os.Remove(path)
	warn := ""
	if len(later) > 0 {
		sort.Ints(later)
		warn = fmt.Sprintf(" Warning: Later days (%v) still exist and may reference the deleted day.", later)
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"success":        true,
		"filename":       filename,
		"character_name": characterName,
		"day_number":     dayNumber,
		"message":        fmt.Sprintf("Day %d deleted successfully for %s (%s).%s", dayNumber, characterName, userID, warn),
	})
}

func (a *App) handleUpdateSchedule(w http.ResponseWriter, r *http.Request) {
	sessionData := a.getCurrentUser(r)

	body, err := decodeBody(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}
	filename := stringField(body, "filename", "")
	scheduleType := stringField(body, "schedule_type", "")
	scheduleData, ok := body["schedule_data"].(map[string]any)
	preset := stringField(body, "preset_set", "presets_kurisu")
	if filename == "" || scheduleType == "" || !ok {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("Filename, schedule_type, and schedule_data are required")))
		return
	}
	if scheduleType != "user_schedule" && scheduleType != "character_schedule" {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("schedule_type must be 'user_schedule' or 'character_schedule'")))
		return
	}
	folder, err := a.getPresetFolder(preset)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}
	path := filepath.Join(folder, filename)
	var data map[string]any
	if err := a.loadJSON(path, &data); err != nil {
		writeJSON(w, http.StatusNotFound, errJSON(errors.New("File not found")))
		return
	}

	// Check ownership: users can only update schedules for their own characters
	if !a.checkFileOwnershipFromData(w, sessionData, data) {
		return
	}

	data[scheduleType] = scheduleData
	if err := a.atomicWriteJSON(path, data); err != nil {
		writeJSON(w, http.StatusInternalServerError, errJSON(err))
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"success": true, "message": fmt.Sprintf("Schedule updated successfully for %s", filename)})
}

func (a *App) handleUpdateDayCategory(w http.ResponseWriter, r *http.Request) {
	body, err := decodeBody(r.Body)
	if err != nil || body == nil {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("No JSON received")))
		return
	}
	filename := stringField(body, "filename", "")
	category := stringField(body, "category", "")
	assessmentTS := numberField(body, "assessment_timestamp", int(time.Now().Unix()))
	problems := numberField(body, "problems_count", 0)
	preset := stringField(body, "preset_set", "")
	dialogueTrait := stringField(body, "dialogue_trait", "")
	password := stringField(body, "approval_password", "")

	if filename == "" || category == "" {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("Filename and category are required")))
		return
	}

	// For approval/rejection operations, require editor role
	if category == "passed" || category == "rejected" {
		if sessionData := a.requireEditor(w, r); sessionData == nil {
			return // requireEditor already sent the error response
		}
	}

	if category == "passed" {
		if password == "" {
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "Password is required for approval", "message": "审批需要密码！\n\nPassword is required for approval!"})
			return
		}
		if password != "Yanqing" {
			writeJSON(w, http.StatusForbidden, map[string]any{"error": "Incorrect password", "message": "密码错误！请重试。\n\nIncorrect password! Please try again."})
			return
		}
	}
	path, actualPreset, err := a.findFileAcrossPresets(filename, preset)
	if err != nil {
		writeJSON(w, http.StatusNotFound, errJSON(errors.New("File not found")))
		return
	}
	var data map[string]any
	if err := a.loadJSON(path, &data); err != nil {
		writeJSON(w, http.StatusInternalServerError, errJSON(err))
		return
	}
	data["category"] = category
	data["assessment_timestamp"] = assessmentTS
	data["problems_count"] = problems
	if dialogueTrait != "" {
		data["dialogue_trait"] = dialogueTrait
	}
	if category == "rejected" {
		data["rejection_timestamp"] = assessmentTS
	}
	if strings.HasPrefix(category, "legacy_") {
		if _, ok := data["legacy_timestamp"]; !ok {
			data["legacy_timestamp"] = assessmentTS
		}
	}
	if err := a.atomicWriteJSON(path, data); err != nil {
		writeJSON(w, http.StatusInternalServerError, errJSON(err))
		return
	}
	resp := map[string]any{
		"success":  true,
		"message":  fmt.Sprintf("Category updated to %s", category),
		"filename": filename,
		"category": category,
	}
	if dialogueTrait != "" {
		resp["dialogue_trait"] = dialogueTrait
		resp["message"] = fmt.Sprintf("Category updated to %s with trait: %s", category, dialogueTrait)
	}
	resp["preset_set"] = actualPreset
	writeJSON(w, http.StatusOK, resp)
}

func (a *App) handleMoveToLegacy(w http.ResponseWriter, r *http.Request) {
	// Only editors can move items to legacy
	if sessionData := a.requireEditor(w, r); sessionData == nil {
		return
	}

	body, err := decodeBody(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}
	filename := stringField(body, "filename", "")
	preset := stringField(body, "preset_set", "")
	if filename == "" {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("Filename is required")))
		return
	}
	path, _, err := a.findFileAcrossPresets(filename, preset)
	if err != nil {
		writeJSON(w, http.StatusNotFound, errJSON(errors.New("File not found")))
		return
	}
	var data map[string]any
	if err := a.loadJSON(path, &data); err != nil {
		writeJSON(w, http.StatusInternalServerError, errJSON(err))
		return
	}
	category := stringField(data, "category", "pending")
	switch category {
	case "passed":
		data["category"] = "legacy_passed"
	case "rejected":
		data["category"] = "legacy_rejected"
	default:
		writeJSON(w, http.StatusBadRequest, errJSON(fmt.Errorf("Cannot move %s items to legacy. Only 'passed' or 'rejected' items can be archived.", category)))
		return
	}
	data["legacy_timestamp"] = time.Now().Unix()
	if err := a.atomicWriteJSON(path, data); err != nil {
		writeJSON(w, http.StatusInternalServerError, errJSON(err))
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"success":  true,
		"message":  fmt.Sprintf("Moved to %s", data["category"]),
		"filename": filename,
		"category": data["category"],
	})
}

func (a *App) handleAutoArchiveOldItems(w http.ResponseWriter, r *http.Request) {
	// Only editors can auto-archive items
	if sessionData := a.requireEditor(w, r); sessionData == nil {
		return
	}

	body, err := decodeBody(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}
	preset := stringField(body, "preset_set", "")
	days := numberField(body, "days_threshold", 30)
	if days <= 0 {
		days = 30
	}
	current := time.Now().Unix()
	threshold := int64(days) * 24 * 3600

	folders := map[string]string{}
	if preset != "" {
		if folder, err := a.getPresetFolder(preset); err == nil {
			folders[preset] = folder
		}
	} else {
		folders = a.presetFolders
	}
	archived := []map[string]any{}
	for ps, folder := range folders {
		entries, _ := os.ReadDir(folder)
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
				continue
			}
			path := filepath.Join(folder, e.Name())
			var data map[string]any
			if err := a.loadJSON(path, &data); err != nil {
				continue
			}
			category := stringField(data, "category", "pending")
			ts := int64(intVal(data["assessment_timestamp"]))
			if category != "passed" && category != "rejected" {
				continue
			}
			if ts == 0 {
				continue
			}
			if current-ts > threshold {
				newCat := "legacy_" + category
				data["category"] = newCat
				data["legacy_timestamp"] = current
				if err := a.atomicWriteJSON(path, data); err == nil {
					archived = append(archived, map[string]any{
						"filename":      e.Name(),
						"preset":        ps,
						"from_category": category,
						"to_category":   newCat,
						"days_old":      int((current - ts) / 86400),
					})
				}
			}
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"success":        true,
		"message":        fmt.Sprintf("Archived %d items older than %d days", len(archived), days),
		"archived_count": len(archived),
		"archived_files": archived,
	})
}

func (a *App) handleUpdateInnerThoughtAnnotation(w http.ResponseWriter, r *http.Request) {
	// Only editors can update inner thought annotations
	if sessionData := a.requireEditor(w, r); sessionData == nil {
		return
	}

	body, err := decodeBody(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}
	filename := stringField(body, "filename", "")
	preset := stringField(body, "preset_set", "")
	utterIdx := numberField(body, "utterance_index", -1)
	actualThought := stringField(body, "actual_thought", "")
	correctThought := stringField(body, "correct_thought", "")
	reviewerNote := stringField(body, "reviewer_note", "")
	if filename == "" || utterIdx < 0 {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("Filename and utterance_index are required")))
		return
	}
	path, preset, err := a.findFileAcrossPresets(filename, preset)
	if err != nil {
		writeJSON(w, http.StatusNotFound, errJSON(errors.New("File not found")))
		return
	}
	var data map[string]any
	if err := a.loadJSON(path, &data); err != nil {
		writeJSON(w, http.StatusInternalServerError, errJSON(err))
		return
	}
	ann, _ := data["inner_thought_annotations"].(map[string]any)
	if ann == nil {
		ann = map[string]any{}
	}
	ann[strconv.Itoa(utterIdx)] = map[string]any{
		"actual_thought":  actualThought,
		"correct_thought": correctThought,
		"reviewer_note":   reviewerNote,
		"timestamp":       time.Now().Unix(),
		"reviewer_ip":     r.RemoteAddr,
	}
	data["inner_thought_annotations"] = ann
	if err := a.atomicWriteJSON(path, data); err != nil {
		writeJSON(w, http.StatusInternalServerError, errJSON(err))
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"success": true, "message": fmt.Sprintf("Inner thought annotation saved for utterance %d", utterIdx), "preset_set": preset})
}

func (a *App) handleDeleteInnerThoughtAnnotation(w http.ResponseWriter, r *http.Request) {
	// Only editors can delete inner thought annotations
	if sessionData := a.requireEditor(w, r); sessionData == nil {
		return
	}

	body, err := decodeBody(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}
	filename := stringField(body, "filename", "")
	preset := stringField(body, "preset_set", "")
	utterIdx := numberField(body, "utterance_index", -1)
	if filename == "" || utterIdx < 0 {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("Filename and utterance_index are required")))
		return
	}
	path, _, err := a.findFileAcrossPresets(filename, preset)
	if err != nil {
		writeJSON(w, http.StatusNotFound, errJSON(errors.New("File not found")))
		return
	}
	var data map[string]any
	if err := a.loadJSON(path, &data); err != nil {
		writeJSON(w, http.StatusInternalServerError, errJSON(err))
		return
	}
	if ann, ok := data["inner_thought_annotations"].(map[string]any); ok {
		delete(ann, strconv.Itoa(utterIdx))
		data["inner_thought_annotations"] = ann
		_ = a.atomicWriteJSON(path, data)
	}
	writeJSON(w, http.StatusOK, map[string]any{"success": true, "message": fmt.Sprintf("Inner thought annotation deleted for utterance %d", utterIdx)})
}

func (a *App) handleSaveChecklistData(w http.ResponseWriter, r *http.Request) {
	// Only editors can save checklist data (part of QC process)
	if sessionData := a.requireEditor(w, r); sessionData == nil {
		return
	}

	body, err := decodeBody(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}
	filename := stringField(body, "filename", "")
	preset := stringField(body, "preset_set", "presets_kurisu")
	checklist := body["checklist"]
	reasons := body["reasons"]
	selectedTrait := stringField(body, "selected_trait", "")
	if filename == "" {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("Filename or preset set not provided")))
		return
	}
	folder, err := a.getPresetFolder(preset)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}
	path := filepath.Join(folder, filename)
	var data map[string]any
	if err := a.loadJSON(path, &data); err != nil {
		writeJSON(w, http.StatusNotFound, errJSON(errors.New("File not found")))
		return
	}
	data["quality_checklist"] = map[string]any{
		"checklist":      checklist,
		"reasons":        reasons,
		"selected_trait": selectedTrait,
		"last_updated":   time.Now().Format(time.RFC3339),
		"updated_by":     r.RemoteAddr,
	}
	if err := a.atomicWriteJSON(path, data); err != nil {
		writeJSON(w, http.StatusInternalServerError, errJSON(err))
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"success": true, "message": "Checklist data saved successfully"})
}

func (a *App) handleLoadChecklistData(w http.ResponseWriter, r *http.Request) {
	body, err := decodeBody(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}
	filename := stringField(body, "filename", "")
	preset := stringField(body, "preset_set", "presets_kurisu")
	if filename == "" {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("Filename is required")))
		return
	}
	folder, err := a.getPresetFolder(preset)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}
	path := filepath.Join(folder, filename)
	var data map[string]any
	if err := a.loadJSON(path, &data); err != nil {
		writeJSON(w, http.StatusNotFound, errJSON(errors.New("File not found")))
		return
	}
	if qc, ok := data["quality_checklist"]; ok {
		writeJSON(w, http.StatusOK, map[string]any{"success": true, "quality_checklist": qc})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"success": true, "quality_checklist": map[string]any{}})
}

// CharacterDefaults contains the default profile values for a character
type CharacterDefaults struct {
	Values      string `json:"values"`
	Experiences string `json:"experiences"`
	Judgements  string `json:"judgements"`
	Abilities   string `json:"abilities"`
}

// getCharacterDefaultsPath returns the path to the _character_defaults.json file for a preset
func (a *App) getCharacterDefaultsPath(preset string) (string, error) {
	folder, err := a.getPresetFolder(preset)
	if err != nil {
		return "", err
	}
	return filepath.Join(folder, "_character_defaults.json"), nil
}

// getCharacterDefaults reads default values from _character_defaults.json file in the preset folder
func (a *App) getCharacterDefaultsFromFile(preset string) *CharacterDefaults {
	defaultsPath, err := a.getCharacterDefaultsPath(preset)
	if err != nil {
		return &CharacterDefaults{}
	}

	data, err := os.ReadFile(defaultsPath)
	if err != nil {
		// File doesn't exist, return empty defaults
		return &CharacterDefaults{}
	}

	var defaults CharacterDefaults
	if err := json.Unmarshal(data, &defaults); err != nil {
		return &CharacterDefaults{}
	}

	return &defaults
}

// getCharacterScheduleForDay fetches the character schedule for a specific day
// from the centralized character_profiles.json file
func (a *App) getCharacterScheduleForDay(preset string, dayNumber int) map[string]any {
	// Map preset to character ID
	var characterId string
	switch {
	case strings.Contains(preset, "kurisu"):
		characterId = "kurisu"
	case strings.Contains(preset, "linlu"):
		characterId = "linlu"
	default:
		return map[string]any{"day": dayNumber}
	}

	// Determine if we should use English based on preset name
	// English presets don't have "_CN" suffix
	useEnglish := !strings.Contains(preset, "_CN")

	// Load character profiles
	profiles, err := a.loadCharacterProfiles()
	if err != nil {
		return map[string]any{"day": dayNumber}
	}

	profile, exists := profiles[characterId]
	if !exists || profile.Schedules == nil {
		return map[string]any{"day": dayNumber}
	}

	// Find the schedule for the requested day
	for _, schedule := range profile.Schedules {
		if schedule.Day == dayNumber {
			// Use English fields if available and preset is English, otherwise use Chinese
			morning := schedule.Morning
			noon := schedule.Noon
			afternoon := schedule.Afternoon
			evening := schedule.Evening
			night := schedule.Night

			if useEnglish {
				if schedule.MorningEn != "" {
					morning = schedule.MorningEn
				}
				if schedule.NoonEn != "" {
					noon = schedule.NoonEn
				}
				if schedule.AfternoonEn != "" {
					afternoon = schedule.AfternoonEn
				}
				if schedule.EveningEn != "" {
					evening = schedule.EveningEn
				}
				if schedule.NightEn != "" {
					night = schedule.NightEn
				}
			}

			return map[string]any{
				"day":       dayNumber,
				"morning":   morning,
				"noon":      noon,
				"afternoon": afternoon,
				"evening":   evening,
				"night":     night,
			}
		}
	}

	// No schedule found for this day, return empty with day number
	return map[string]any{"day": dayNumber}
}

// Deprecated: getCharacterDefaults is now a wrapper that calls the App method
// Keep for backward compatibility with existing code
func getCharacterDefaults(preset string) *CharacterDefaults {
	// Fallback hardcoded defaults when App instance is not available
	switch {
	case strings.Contains(preset, "kurisu"):
		return &CharacterDefaults{
			Values:      "1. 理性\n2. 创造力\n3. 成就\n4. 尊重\n5. 自由\n6. 和谐\n7. 秩序\n8. 正直\n9. 健康\n10. 传统、信仰\n11. 权力\n12. 利他、仁慈",
			Experiences: "研究项目：使用本地的粒子加速器，将特定粒子束射入实验对象的大脑，来研究其内部映射图。\n\n学业经历：二年级时开始读物理，闲的没事，但挺好玩",
			Judgements:  "• 人再笨四年级学不会微积分么？\n• 哲学……我对自己想要什么虽然不是很明白，但是一个莫名的人跟我说起话来还是有点奇怪吧",
			Abilities:   "模型训练、算法设计、物理半导体芯片、蛋白药物靶向、嵌入式设计",
		}
	default:
		return &CharacterDefaults{}
	}
}

// handleGetCharacterDefaults returns the default profile values for a character preset
func (a *App) handleGetCharacterDefaults(w http.ResponseWriter, r *http.Request) {
	preset := chi.URLParam(r, "preset")
	if preset == "" {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("preset parameter required")))
		return
	}

	defaults := a.getCharacterDefaultsFromFile(preset)
	writeJSON(w, http.StatusOK, map[string]any{
		"success":  true,
		"preset":   preset,
		"defaults": defaults,
	})
}

// handleUpdateCharacterDefaults updates the _character_defaults.json file
func (a *App) handleUpdateCharacterDefaults(w http.ResponseWriter, r *http.Request) {
	preset := chi.URLParam(r, "preset")
	if preset == "" {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("preset parameter required")))
		return
	}

	body, err := decodeBody(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}

	defaults := CharacterDefaults{
		Values:      stringField(body, "values", ""),
		Experiences: stringField(body, "experiences", ""),
		Judgements:  stringField(body, "judgements", ""),
		Abilities:   stringField(body, "abilities", ""),
	}

	defaultsPath, err := a.getCharacterDefaultsPath(preset)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}

	if err := a.atomicWriteJSON(defaultsPath, defaults); err != nil {
		writeJSON(w, http.StatusInternalServerError, errJSON(err))
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"message": fmt.Sprintf("Character defaults updated for %s", preset),
	})
}

// handleSyncCharacterDefaults propagates character defaults to all Day 1 files in a preset
func (a *App) handleSyncCharacterDefaults(w http.ResponseWriter, r *http.Request) {
	body, err := decodeBody(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}

	preset := stringField(body, "preset_set", "")
	syncAllDays := boolField(body, "sync_all_days", false)

	if preset == "" {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("preset_set is required")))
		return
	}

	folder, err := a.getPresetFolder(preset)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}

	// Read defaults from file
	defaults := a.getCharacterDefaultsFromFile(preset)

	entries, err := os.ReadDir(folder)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errJSON(err))
		return
	}

	updated := []string{}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		// Skip the defaults file itself
		if e.Name() == "_character_defaults.json" {
			continue
		}

		pf, ok := parseWriterFilename(e.Name())
		if !ok {
			continue
		}

		// Only update Day 1 files unless syncAllDays is true
		if !syncAllDays && pf.DayNum != 1 {
			continue
		}

		path := filepath.Join(folder, e.Name())
		var data map[string]any
		if err := a.loadJSON(path, &data); err != nil {
			continue
		}

		// Update the values, experiences, judgements, abilities
		changed := false
		if defaults.Values != "" {
			data["values"] = defaults.Values
			changed = true
		}
		if defaults.Experiences != "" {
			data["experiences"] = defaults.Experiences
			changed = true
		}
		if defaults.Judgements != "" {
			data["judgements"] = defaults.Judgements
			changed = true
		}
		if defaults.Abilities != "" {
			data["abilities"] = defaults.Abilities
			changed = true
		}

		if changed {
			if err := a.atomicWriteJSON(path, data); err == nil {
				updated = append(updated, e.Name())
			}
		}
	}

	scope := "Day 1 files"
	if syncAllDays {
		scope = "all files"
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success":       true,
		"updated_count": len(updated),
		"updated_files": updated,
		"message":       fmt.Sprintf("Synced character defaults to %d %s in %s", len(updated), scope, preset),
	})
}

// ===== helper funcs used by handlers =====

func (a *App) getPreviousDayIntimacy(pf *ParsedFilename, preset string) int {
	folder, err := a.getPresetFolder(preset)
	if err != nil {
		return 0
	}
	prevDay := pf.DayNum - 1
	if prevDay < 1 {
		return 0
	}
	filename := fmt.Sprintf("%s_Day%d_%s_simplified.json", pf.UserID, prevDay, pf.DupID)
	path := filepath.Join(folder, filename)
	var data map[string]any
	if err := a.loadJSON(path, &data); err != nil {
		return 0
	}
	return intVal(data["intimacy_level"])
}

func formatScheduleToStringAny(schedule map[string]any, name string, preset string) string {
	isCN := strings.Contains(preset, "_CN")
	timeOrder := []string{"morning", "noon", "afternoon", "evening", "night"}
	timeMap := map[string]string{
		"morning":   ternary(isCN, "早晨", "morning"),
		"noon":      ternary(isCN, "中午", "noon"),
		"afternoon": ternary(isCN, "下午", "afternoon"),
		"evening":   ternary(isCN, "晚上", "evening"),
		"night":     ternary(isCN, "夜晚", "night"),
	}
	labels := map[string]string{
		"morning":   ternary(isCN, "早晨活动", "Morning Event"),
		"noon":      ternary(isCN, "中午活动", "Noon Event"),
		"afternoon": ternary(isCN, "下午活动", "Afternoon Event"),
		"evening":   ternary(isCN, "晚上活动", "Evening Event"),
		"night":     ternary(isCN, "夜晚活动", "Night Event"),
	}
	var parts []string
	for _, p := range timeOrder {
		key := timeMap[p]
		if v, ok := schedule[p]; ok {
			if s, ok := v.(string); ok && s != "" {
				parts = append(parts, fmt.Sprintf("%s: %s", labels[p], s))
				continue
			}
		}
		if v, ok := schedule[key]; ok {
			if s, ok := v.(string); ok && s != "" {
				parts = append(parts, fmt.Sprintf("%s: %s", labels[p], s))
			}
		}
	}
	return strings.Join(parts, "\n")
}

func ternary[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}

func (a *App) updateUserInfoFromFile(body map[string]any, fileData map[string]any, preset string) {
	filename := stringField(body, "filename", "")
	pf, ok := parseWriterFilename(filename)
	if !ok {
		return
	}
	userID := pf.UserID
	userNum, _ := strconv.Atoi(strings.TrimPrefix(userID, "user_"))
	folder, err := a.getPresetFolder(preset)
	if err != nil {
		return
	}
	userInfoFile := filepath.Join(folder, getUserInfoFilename(preset))
	var list []map[string]any
	if b, err := os.ReadFile(userInfoFile); err == nil {
		_ = json.Unmarshal(b, &list)
	}
	updated := false
	for i := range list {
		if intVal(list[i]["id"]) == userNum {
			for _, k := range []string{"real_name", "name", "description", "motivation"} {
				switch k {
				case "real_name":
					if v, ok := body["real_name"]; ok {
						list[i]["real_name"] = v
						updated = true
					}
				case "name":
					if v, ok := body["user_name"]; ok {
						list[i]["name"] = v
						updated = true
					}
				case "description":
					if v, ok := body["character_description"]; ok {
						list[i]["description"] = v
						updated = true
					}
				case "motivation":
					if v, ok := body["character_motivation"]; ok {
						list[i]["motivation"] = v
						updated = true
					}
				}
			}
			break
		}
	}
	if updated {
		_ = a.atomicWriteJSON(userInfoFile, list)
	}
	// update user_info_CN.js best effort
	if strings.Contains(preset, "_CN") {
		currentCharacter := stringField(body, "user_name", stringField(fileData, "user_name", ""))
		currentReal := stringField(body, "real_name", stringField(fileData, "real_name", ""))
		a.updateUserInfoCNLine(preset, currentCharacter, currentReal)
	}
}

func (a *App) addUserInfoEntry(folder, preset, userID string, userNum int, name, realName, age, profession, motivation, description string, meta map[string]any) {
	userInfoFile := filepath.Join(folder, getUserInfoFilename(preset))
	var list []map[string]any
	if b, err := os.ReadFile(userInfoFile); err == nil {
		_ = json.Unmarshal(b, &list)
	}
	englishName := name
	if containsChinese(name) {
		englishName = toPinyin(name)
	}
	list = append(list, map[string]any{
		"name":               name,
		"english_name":       englishName,
		"real_name":          realName,
		"id":                 userNum,
		"created":            true,
		"age":                age,
		"profession":         profession,
		"motivation":         motivation,
		"description":        description,
		"_creation_metadata": meta,
	})
	_ = a.atomicWriteJSON(userInfoFile, list)
	if strings.Contains(preset, "_CN") {
		a.insertUserInfoCNLine(preset, name, realName)
	}
}

func (a *App) removeUserInfoEntry(folder, preset, userID string) {
	userInfoFile := filepath.Join(folder, getUserInfoFilename(preset))
	var list []map[string]any
	if b, err := os.ReadFile(userInfoFile); err == nil {
		_ = json.Unmarshal(b, &list)
	}
	userNum, _ := strconv.Atoi(strings.TrimPrefix(userID, "user_"))
	newList := []map[string]any{}
	for _, u := range list {
		if intVal(u["id"]) != userNum {
			newList = append(newList, u)
		}
	}
	_ = a.atomicWriteJSON(userInfoFile, newList)
}

// Placeholder pinyin conversion: strip spaces; keep ASCII
func toPinyin(s string) string {
	return strings.Map(func(r rune) rune {
		if r == ' ' {
			return -1
		}
		return r
	}, s)
}

func containsChinese(s string) bool {
	for _, r := range s {
		if isHan(r) {
			return true
		}
	}
	return false
}

func (a *App) updateRealNameInUserInfo(folder, preset, userID, realName, characterName string) {
	userInfoFile := filepath.Join(folder, getUserInfoFilename(preset))
	var list []map[string]any
	if b, err := os.ReadFile(userInfoFile); err == nil {
		_ = json.Unmarshal(b, &list)
	}
	userNum, _ := strconv.Atoi(strings.TrimPrefix(userID, "user_"))
	for i := range list {
		if intVal(list[i]["id"]) == userNum {
			list[i]["real_name"] = realName
			if characterName == "" {
				characterName = stringField(list[i], "name", characterName)
			}
			break
		}
	}
	_ = a.atomicWriteJSON(userInfoFile, list)
	if strings.Contains(preset, "_CN") && characterName != "" {
		a.updateUserInfoCNLine(preset, characterName, realName)
	}
}

func (a *App) updateUserInfoCNLine(preset, characterName, realName string) {
	jsPath := a.userInfoCNPath(preset)
	if jsPath == "" {
		return
	}
	b, err := os.ReadFile(jsPath)
	if err != nil {
		return
	}
	lines := strings.Split(string(b), "\n")
	pat := regexp.MustCompile(fmt.Sprintf("^\\s*\"%s\":`[^`]*`,\\s*$", regexp.QuoteMeta(characterName)))
	newLine := fmt.Sprintf(`  "%s":`+"`作者 ——%s`,", characterName, realName)
	updated := false
	for i, line := range lines {
		if pat.MatchString(line) {
			lines[i] = newLine
			updated = true
			break
		}
	}
	if updated {
		_ = os.WriteFile(jsPath, []byte(strings.Join(lines, "\n")), 0o644)
	}
}

func (a *App) insertUserInfoCNLine(preset, characterName, realName string) {
	jsPath := a.userInfoCNPath(preset)
	if jsPath == "" {
		return
	}
	b, err := os.ReadFile(jsPath)
	if err != nil {
		return
	}
	lines := strings.Split(string(b), "\n")
	insertAt := len(lines)
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.TrimSpace(lines[i]) == "};" {
			insertAt = i
			break
		}
	}
	newLine := fmt.Sprintf(`  "%s":`+"`作者 ——%s`,", characterName, realName)
	lines = append(lines[:insertAt], append([]string{newLine}, lines[insertAt:]...)...)
	_ = os.WriteFile(jsPath, []byte(strings.Join(lines, "\n")), 0o644)
}

func (a *App) removeUserInfoCNLine(preset, characterName string) {
	jsPath := a.userInfoCNPath(preset)
	if jsPath == "" {
		return
	}
	b, err := os.ReadFile(jsPath)
	if err != nil {
		return
	}
	lines := strings.Split(string(b), "\n")
	pat := regexp.MustCompile(fmt.Sprintf("^\\s*\"%s\":`[^`]*`,\\s*$", regexp.QuoteMeta(characterName)))
	newLines := []string{}
	for _, line := range lines {
		if pat.MatchString(line) {
			continue
		}
		newLines = append(newLines, line)
	}
	_ = os.WriteFile(jsPath, []byte(strings.Join(newLines, "\n")), 0o644)
}

func (a *App) userInfoCNPath(preset string) string {
	if !strings.Contains(preset, "_CN") {
		return ""
	}
	subdir := strings.TrimPrefix(preset, "presets_")
	if strings.HasSuffix(subdir, "_CN") {
		subdir = strings.TrimSuffix(subdir, "_CN")
	}
	subdir = strings.ReplaceAll(subdir, "-", " ")
	return filepath.Join(a.cfg.StaticDir, subdir, "user_info_CN.js")
}

// JSON response helper
func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// CharacterProfile represents the profile data for a character
// This is the centralized source of truth used by /descriptions page
type CharacterProfile struct {
	ID             string                   `json:"id"`
	Name           string                   `json:"name"`
	EnglishName    string                   `json:"english_name"`
	Tagline        string                   `json:"tagline"`
	TaglineEn      string                   `json:"tagline_en,omitempty"`
	Description    string                   `json:"description"`
	DescriptionEn  string                   `json:"description_en,omitempty"`
	Values         string                   `json:"values"`
	ValuesEn       string                   `json:"values_en,omitempty"`
	Experiences    string                   `json:"experiences"`
	ExperiencesEn  string                   `json:"experiences_en,omitempty"`
	Judgements     string                   `json:"judgements"`
	JudgementsEn   string                   `json:"judgements_en,omitempty"`
	Abilities      string                   `json:"abilities"`
	AbilitiesEn    string                   `json:"abilities_en,omitempty"`
	Relationships  []map[string]interface{} `json:"relationships,omitempty"`
	Schedules      []CharacterDaySchedule   `json:"schedules,omitempty"`
	Story          map[string]interface{}   `json:"story,omitempty"`
}

// CharacterDaySchedule represents a single day's schedule for a character
type CharacterDaySchedule struct {
	Day          int    `json:"day"`
	Title        string `json:"title"`
	TitleEn      string `json:"title_en,omitempty"`
	Morning      string `json:"morning"`
	MorningEn    string `json:"morning_en,omitempty"`
	Noon         string `json:"noon"`
	NoonEn       string `json:"noon_en,omitempty"`
	Afternoon    string `json:"afternoon"`
	AfternoonEn  string `json:"afternoon_en,omitempty"`
	Evening      string `json:"evening"`
	EveningEn    string `json:"evening_en,omitempty"`
	Night        string `json:"night"`
	NightEn      string `json:"night_en,omitempty"`
}

// getCharacterProfilesPath returns the path to the character_profiles.json file
func (a *App) getCharacterProfilesPath() string {
	return filepath.Join(a.cfg.RootDir, "data", "character_profiles.json")
}

// loadCharacterProfiles loads all character profiles from the JSON file
func (a *App) loadCharacterProfiles() (map[string]*CharacterProfile, error) {
	profilesPath := a.getCharacterProfilesPath()
	data, err := os.ReadFile(profilesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read character profiles: %w", err)
	}

	var profiles map[string]*CharacterProfile
	if err := json.Unmarshal(data, &profiles); err != nil {
		return nil, fmt.Errorf("failed to parse character profiles: %w", err)
	}

	return profiles, nil
}

// handleGetAllCharacterProfiles returns all character profiles
func (a *App) handleGetAllCharacterProfiles(w http.ResponseWriter, r *http.Request) {
	profiles, err := a.loadCharacterProfiles()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success":  true,
		"profiles": profiles,
	})
}

// handleGetCharacterProfile returns profile data for a specific character
// Character ID mapping:
//   - "kurisu" -> Kurisu (牧濑红莉栖)
//   - "linlu" -> Lin Lu (林路)
func (a *App) handleGetCharacterProfile(w http.ResponseWriter, r *http.Request) {
	characterId := chi.URLParam(r, "characterId")
	if characterId == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "character ID required",
		})
		return
	}

	// Check language preference from cookie or query param
	lang := "zh" // default to Chinese
	if cookie, err := r.Cookie("language"); err == nil && cookie.Value == "en" {
		lang = "en"
	}
	if q := r.URL.Query().Get("lang"); q == "en" {
		lang = "en"
	}

	// Normalize character ID - map route names to profile IDs
	normalizedId := characterId
	switch strings.ToLower(characterId) {
	case "kurisu", "makise-kurisu":
		normalizedId = "kurisu"
	case "linlu", "lin-lu", "lin_lu":
		normalizedId = "linlu"
	}

	profiles, err := a.loadCharacterProfiles()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	profile, exists := profiles[normalizedId]
	if !exists {
		writeJSON(w, http.StatusNotFound, map[string]any{
			"success": false,
			"error":   fmt.Sprintf("character profile not found: %s", characterId),
		})
		return
	}

	// If English is requested, swap in English fields where available
	if lang == "en" {
		if profile.ValuesEn != "" {
			profile.Values = profile.ValuesEn
		}
		if profile.ExperiencesEn != "" {
			profile.Experiences = profile.ExperiencesEn
		}
		if profile.JudgementsEn != "" {
			profile.Judgements = profile.JudgementsEn
		}
		if profile.AbilitiesEn != "" {
			profile.Abilities = profile.AbilitiesEn
		}
		if profile.DescriptionEn != "" {
			profile.Description = profile.DescriptionEn
		}
		if profile.TaglineEn != "" {
			profile.Tagline = profile.TaglineEn
		}
		// Use English name as primary name for English mode
		if profile.EnglishName != "" {
			profile.Name = profile.EnglishName
		}
		// Swap relationship fields too
		for i := range profile.Relationships {
			if descEn, ok := profile.Relationships[i]["description_en"].(string); ok && descEn != "" {
				profile.Relationships[i]["description"] = descEn
			}
			if nameEn, ok := profile.Relationships[i]["name_en"].(string); ok && nameEn != "" {
				profile.Relationships[i]["name"] = nameEn
			}
			if styleEn, ok := profile.Relationships[i]["style_en"].(string); ok && styleEn != "" {
				profile.Relationships[i]["style"] = styleEn
			}
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"profile": profile,
	})
}

// CharacterSchedule represents a day's schedule for a character
type CharacterSchedule struct {
	Day       int    `json:"day"`
	Title     string `json:"title"`
	Morning   string `json:"morning"`
	Noon      string `json:"noon"`
	Afternoon string `json:"afternoon"`
	Evening   string `json:"evening"`
	Night     string `json:"night"`
}

// handleGetCharacterScheduleForDay returns the schedule for a specific day for a character
// This is used to populate the character_schedule field when creating new days
func (a *App) handleGetCharacterScheduleForDay(w http.ResponseWriter, r *http.Request) {
	characterId := chi.URLParam(r, "characterId")
	dayStr := chi.URLParam(r, "day")

	if characterId == "" || dayStr == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "character ID and day number required",
		})
		return
	}

	day, err := strconv.Atoi(dayStr)
	if err != nil || day < 1 {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "invalid day number",
		})
		return
	}

	// Normalize character ID
	normalizedId := characterId
	switch strings.ToLower(characterId) {
	case "kurisu", "makise-kurisu":
		normalizedId = "kurisu"
	case "linlu", "lin-lu", "lin_lu":
		normalizedId = "linlu"
	}

	profiles, err := a.loadCharacterProfiles()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	profile, exists := profiles[normalizedId]
	if !exists {
		writeJSON(w, http.StatusNotFound, map[string]any{
			"success": false,
			"error":   fmt.Sprintf("character profile not found: %s", characterId),
		})
		return
	}

	// Check language preference from cookie or query param
	lang := "zh" // default to Chinese
	if cookie, err := r.Cookie("language"); err == nil && cookie.Value == "en" {
		lang = "en"
	}
	if q := r.URL.Query().Get("lang"); q == "en" {
		lang = "en"
	}

	// Find the schedule for the requested day
	var schedule *CharacterSchedule
	if profile.Schedules != nil {
		for _, s := range profile.Schedules {
			if s.Day == day {
				// Use English fields if language is English and English content exists
				title := s.Title
				morning := s.Morning
				noon := s.Noon
				afternoon := s.Afternoon
				evening := s.Evening
				night := s.Night
				
				if lang == "en" {
					if s.TitleEn != "" {
						title = s.TitleEn
					}
					if s.MorningEn != "" {
						morning = s.MorningEn
					}
					if s.NoonEn != "" {
						noon = s.NoonEn
					}
					if s.AfternoonEn != "" {
						afternoon = s.AfternoonEn
					}
					if s.EveningEn != "" {
						evening = s.EveningEn
					}
					if s.NightEn != "" {
						night = s.NightEn
					}
				}
				
				schedule = &CharacterSchedule{
					Day:       s.Day,
					Title:     title,
					Morning:   morning,
					Noon:      noon,
					Afternoon: afternoon,
					Evening:   evening,
					Night:     night,
				}
				break
			}
		}
	}

	if schedule == nil {
		writeJSON(w, http.StatusNotFound, map[string]any{
			"success":    false,
			"error":      fmt.Sprintf("no schedule found for day %d", day),
			"total_days": len(profile.Schedules),
			"character":  normalizedId,
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success":   true,
		"schedule":  schedule,
		"character": normalizedId,
		"day":       day,
	})
}

// handleUpdateCharacterSchedule creates or updates a schedule for a specific day
// This is used by editors to add/edit daily schedules from the /descriptions page
func (a *App) handleUpdateCharacterSchedule(w http.ResponseWriter, r *http.Request) {
	characterId := chi.URLParam(r, "characterId")
	if characterId == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "character ID required",
		})
		return
	}

	var body struct {
		Day         int    `json:"day"`
		Title       string `json:"title"`
		TitleEn     string `json:"title_en"`
		Morning     string `json:"morning"`
		MorningEn   string `json:"morning_en"`
		Noon        string `json:"noon"`
		NoonEn      string `json:"noon_en"`
		Afternoon   string `json:"afternoon"`
		AfternoonEn string `json:"afternoon_en"`
		Evening     string `json:"evening"`
		EveningEn   string `json:"evening_en"`
		Night       string `json:"night"`
		NightEn     string `json:"night_en"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "invalid JSON: " + err.Error(),
		})
		return
	}

	if body.Day < 1 {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "day number must be >= 1",
		})
		return
	}

	// Normalize character ID
	normalizedId := characterId
	switch strings.ToLower(characterId) {
	case "kurisu", "makise-kurisu":
		normalizedId = "kurisu"
	case "linlu", "lin-lu", "lin_lu":
		normalizedId = "linlu"
	}

	// Load current profiles
	profilesPath := a.getCharacterProfilesPath()
	data, err := os.ReadFile(profilesPath)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   "failed to read profiles: " + err.Error(),
		})
		return
	}

	var profiles map[string]json.RawMessage
	if err := json.Unmarshal(data, &profiles); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   "failed to parse profiles: " + err.Error(),
		})
		return
	}

	profileData, exists := profiles[normalizedId]
	if !exists {
		writeJSON(w, http.StatusNotFound, map[string]any{
			"success": false,
			"error":   fmt.Sprintf("character profile not found: %s", characterId),
		})
		return
	}

	// Parse the specific profile
	var profile map[string]interface{}
	if err := json.Unmarshal(profileData, &profile); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   "failed to parse profile: " + err.Error(),
		})
		return
	}

	// Create the new schedule entry
	newSchedule := map[string]interface{}{
		"day":          body.Day,
		"title":        body.Title,
		"title_en":     body.TitleEn,
		"morning":      body.Morning,
		"morning_en":   body.MorningEn,
		"noon":         body.Noon,
		"noon_en":      body.NoonEn,
		"afternoon":    body.Afternoon,
		"afternoon_en": body.AfternoonEn,
		"evening":      body.Evening,
		"evening_en":   body.EveningEn,
		"night":        body.Night,
		"night_en":     body.NightEn,
	}

	// Get or create schedules array
	schedules, ok := profile["schedules"].([]interface{})
	if !ok {
		schedules = []interface{}{}
	}

	// Check if schedule for this day exists - update or append
	found := false
	for i, s := range schedules {
		if sched, ok := s.(map[string]interface{}); ok {
			if dayNum, ok := sched["day"].(float64); ok && int(dayNum) == body.Day {
				schedules[i] = newSchedule
				found = true
				break
			}
		}
	}
	if !found {
		schedules = append(schedules, newSchedule)
	}

	// Sort schedules by day number
	sort.Slice(schedules, func(i, j int) bool {
		di, _ := schedules[i].(map[string]interface{})["day"].(float64)
		dj, _ := schedules[j].(map[string]interface{})["day"].(float64)
		return di < dj
	})

	profile["schedules"] = schedules

	// Marshal profile back
	newProfileData, err := json.Marshal(profile)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   "failed to marshal profile: " + err.Error(),
		})
		return
	}
	profiles[normalizedId] = newProfileData

	// Write back the entire profiles file with pretty formatting
	output, err := json.MarshalIndent(profiles, "", "  ")
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   "failed to marshal profiles: " + err.Error(),
		})
		return
	}

	if err := os.WriteFile(profilesPath, output, 0644); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   "failed to write profiles: " + err.Error(),
		})
		return
	}

	action := "updated"
	if !found {
		action = "created"
	}

	log.Printf("Schedule %s for %s day %d", action, normalizedId, body.Day)

	writeJSON(w, http.StatusOK, map[string]any{
		"success":   true,
		"action":    action,
		"character": normalizedId,
		"day":       body.Day,
	})
}

// handleDeleteCharacterSchedule deletes a schedule for a specific day (editors only)
func (a *App) handleDeleteCharacterSchedule(w http.ResponseWriter, r *http.Request) {
	characterId := chi.URLParam(r, "characterId")
	dayStr := chi.URLParam(r, "day")

	if characterId == "" || dayStr == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "character ID and day number required",
		})
		return
	}

	day, err := strconv.Atoi(dayStr)
	if err != nil || day < 1 {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "invalid day number",
		})
		return
	}

	// Normalize character ID
	normalizedId := characterId
	switch strings.ToLower(characterId) {
	case "kurisu", "makise-kurisu":
		normalizedId = "kurisu"
	case "linlu", "lin-lu", "lin_lu":
		normalizedId = "linlu"
	}

	// Load current profiles
	profilesPath := a.getCharacterProfilesPath()
	data, err := os.ReadFile(profilesPath)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   "failed to read profiles: " + err.Error(),
		})
		return
	}

	var profiles map[string]json.RawMessage
	if err := json.Unmarshal(data, &profiles); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   "failed to parse profiles: " + err.Error(),
		})
		return
	}

	profileData, exists := profiles[normalizedId]
	if !exists {
		writeJSON(w, http.StatusNotFound, map[string]any{
			"success": false,
			"error":   fmt.Sprintf("character profile not found: %s", characterId),
		})
		return
	}

	// Parse the specific profile
	var profile map[string]interface{}
	if err := json.Unmarshal(profileData, &profile); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   "failed to parse profile: " + err.Error(),
		})
		return
	}

	// Get schedules array
	schedules, ok := profile["schedules"].([]interface{})
	if !ok || len(schedules) == 0 {
		writeJSON(w, http.StatusNotFound, map[string]any{
			"success": false,
			"error":   "no schedules found",
		})
		return
	}

	// Find and remove the schedule for this day
	found := false
	newSchedules := []interface{}{}
	for _, s := range schedules {
		if sched, ok := s.(map[string]interface{}); ok {
			if dayNum, ok := sched["day"].(float64); ok && int(dayNum) == day {
				found = true
				continue // Skip this schedule (delete it)
			}
		}
		newSchedules = append(newSchedules, s)
	}

	if !found {
		writeJSON(w, http.StatusNotFound, map[string]any{
			"success": false,
			"error":   fmt.Sprintf("no schedule found for day %d", day),
		})
		return
	}

	profile["schedules"] = newSchedules

	// Marshal profile back
	newProfileData, err := json.Marshal(profile)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   "failed to marshal profile: " + err.Error(),
		})
		return
	}
	profiles[normalizedId] = newProfileData

	// Write back the entire profiles file with pretty formatting
	output, err := json.MarshalIndent(profiles, "", "  ")
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   "failed to marshal profiles: " + err.Error(),
		})
		return
	}

	if err := os.WriteFile(profilesPath, output, 0644); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   "failed to write profiles: " + err.Error(),
		})
		return
	}

	log.Printf("Schedule deleted for %s day %d", normalizedId, day)

	writeJSON(w, http.StatusOK, map[string]any{
		"success":   true,
		"action":    "deleted",
		"character": normalizedId,
		"day":       day,
	})
}

// handleGetPassedTokens returns the count of passed tokens (Chinese characters)
// per character from all approved dialogues
func (a *App) handleGetPassedTokens(w http.ResponseWriter, r *http.Request) {
	targetTokensPerChar := 500000

	// Map preset names to character IDs (must match IDs in handlers.go loadCharacters)
	presetToCharacter := map[string]string{
		"presets_kurisu":         "kurisu",
		"presets_kurisu_CN":      "kurisu",
		"presets_lin_lu":         "lin_lu",
		"presets_lin_lu_CN":      "lin_lu",
		"presets_newcharacter_1": "newcharacter_1",
	}

	// Track tokens per character
	characterTokens := map[string]int{
		"kurisu":         0,
		"lin_lu":         0,
		"newcharacter_1": 0,
	}

	// Iterate through all preset folders
	for preset, folder := range a.presetFolders {
		characterID := presetToCharacter[preset]
		if characterID == "" {
			continue
		}

		entries, err := os.ReadDir(folder)
		if err != nil {
			continue
		}

		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			name := e.Name()

			// Skip non-data files
			if name == "_character_defaults.json" || strings.HasPrefix(name, "new_user_info_") {
				continue
			}

			// Use parseWriterFilename to check if it's a valid writer file
			// This handles both .json and .bak files
			pf, ok := parseWriterFilename(name)
			if !ok {
				continue
			}
			_ = pf // suppress unused warning

			path := filepath.Join(folder, name)
			var data map[string]any
			if err := a.loadJSON(path, &data); err != nil {
				continue
			}

			// Only count passed dialogues
			category := stringField(data, "category", "pending")
			if category == "passed" || category == "legacy_passed" {
				dialogue, _ := data["dialogue"].([]any)
				charCount := countDialogueHanChars(dialogue)
				characterTokens[characterID] += charCount
			}
		}
	}

	// Build response with per-character stats
	characters := map[string]map[string]any{}
	totalPassedTokens := 0

	for charID, tokens := range characterTokens {
		totalPassedTokens += tokens
		percentage := float64(tokens) / float64(targetTokensPerChar) * 100
		if percentage > 100 {
			percentage = 100
		}
		characters[charID] = map[string]any{
			"passed_tokens": tokens,
			"target_tokens": targetTokensPerChar,
			"percentage":    math.Round(percentage*100) / 100,
		}
	}

	// Calculate total percentage
	totalTarget := targetTokensPerChar * 3 // 3 main characters
	totalPercentage := float64(totalPassedTokens) / float64(totalTarget) * 100
	if totalPercentage > 100 {
		totalPercentage = 100
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success":             true,
		"characters":          characters,
		"total_passed_tokens": totalPassedTokens,
		"total_target_tokens": totalTarget,
		"total_percentage":    math.Round(totalPercentage*100) / 100,
	})
}

// handleRegenerateStoryboard runs the Python export script to regenerate the storyboard markdown
func (a *App) handleRegenerateStoryboard(w http.ResponseWriter, r *http.Request) {
	// Get the base directory (where export_storyboard.py should be)
	baseDir := filepath.Dir(a.cfg.TemplatesDir)

	// Run the Python script
	cmd := exec.Command("python3", "export_storyboard.py")
	cmd.Dir = baseDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed to regenerate storyboard: %v\nOutput: %s", err, string(output))
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   fmt.Sprintf("failed to run export script: %v", err),
			"output":  string(output),
		})
		return
	}

	log.Printf("Storyboard regenerated successfully:\n%s", string(output))

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"message": "Storyboard regenerated successfully",
		"output":  string(output),
	})
}

// handleGetStickers returns the list of available stickers with their metadata
func (a *App) handleGetStickers(w http.ResponseWriter, r *http.Request) {
	stickersMapPath := filepath.Join(a.cfg.RootDir, "stickers", "stickers_map.json")

	var stickersMap map[string]map[string]string
	if err := a.loadJSON(stickersMapPath, &stickersMap); err != nil {
		writeJSON(w, http.StatusInternalServerError, errJSON(fmt.Errorf("failed to load stickers map: %w", err)))
		return
	}

	// Convert to array format for frontend
	stickers := make([]map[string]any, 0, len(stickersMap))
	for filename, metadata := range stickersMap {
		// Build URL based on folder field - if folder is empty, sticker is in root stickers dir
		folder := metadata["folder"]
		var url string
		if folder == "" {
			url = fmt.Sprintf("/static/stickers/%s", filename)
		} else {
			url = fmt.Sprintf("/static/stickers/%s/%s", folder, filename)
		}
		stickers = append(stickers, map[string]any{
			"filename":    filename,
			"description": metadata["description"],
			"folder":      folder,
			"url":         url,
		})
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success":  true,
		"stickers": stickers,
	})
}

// handleValidateSticker validates if a sticker exists before saving
func (a *App) handleValidateSticker(w http.ResponseWriter, r *http.Request) {
	body, err := decodeBody(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}

	filename := stringField(body, "filename", "")
	if filename == "" {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("filename is required")))
		return
	}

	// Check if sticker exists in map
	stickersMapPath := filepath.Join(a.cfg.RootDir, "stickers", "stickers_map.json")
	var stickersMap map[string]map[string]string
	if err := a.loadJSON(stickersMapPath, &stickersMap); err != nil {
		writeJSON(w, http.StatusInternalServerError, errJSON(fmt.Errorf("failed to load stickers map: %w", err)))
		return
	}

	metadata, exists := stickersMap[filename]
	if !exists {
		writeJSON(w, http.StatusNotFound, errJSON(fmt.Errorf("sticker not found: %s", filename)))
		return
	}

	// Check if file physically exists - use folder from metadata
	folder := metadata["folder"]
	var stickerPath string
	if folder == "" {
		stickerPath = filepath.Join(a.cfg.RootDir, "stickers", filename)
	} else {
		stickerPath = filepath.Join(a.cfg.RootDir, "stickers", folder, filename)
	}
	if _, err := os.Stat(stickerPath); os.IsNotExist(err) {
		writeJSON(w, http.StatusNotFound, errJSON(fmt.Errorf("sticker file not found: %s", filename)))
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success":     true,
		"filename":    filename,
		"description": metadata["description"],
		"folder":      folder,
		"message":     "Sticker is valid",
	})
}

// handleUploadSticker allows writers to upload new stickers
func (a *App) handleUploadSticker(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form with 10MB max
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(fmt.Errorf("failed to parse form: %w", err)))
		return
	}

	// Get the file
	file, header, err := r.FormFile("sticker")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(fmt.Errorf("failed to get file: %w", err)))
		return
	}
	defer file.Close()

	// Get metadata
	description := r.FormValue("description")
	if description == "" {
		description = "Uploaded Sticker / 上传贴纸"
	}

	// Validate file type
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".png" && ext != ".gif" && ext != ".jpg" && ext != ".jpeg" && ext != ".webp" {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("invalid file type: only png, gif, jpg, jpeg, webp allowed")))
		return
	}

	// Generate unique filename
	timestamp := time.Now().UnixNano()
	safeFilename := fmt.Sprintf("upload_%d%s", timestamp, ext)

	// Save to uploads subfolder (will be created if not exists)
	uploadsDir := filepath.Join(a.cfg.RootDir, "stickers", "uploads")
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		writeJSON(w, http.StatusInternalServerError, errJSON(fmt.Errorf("failed to create uploads directory: %w", err)))
		return
	}

	destPath := filepath.Join(uploadsDir, safeFilename)
	destFile, err := os.Create(destPath)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errJSON(fmt.Errorf("failed to create file: %w", err)))
		return
	}
	defer destFile.Close()

	// Copy file content
	buf := make([]byte, 1024*1024) // 1MB buffer
	for {
		n, readErr := file.Read(buf)
		if n > 0 {
			if _, writeErr := destFile.Write(buf[:n]); writeErr != nil {
				writeJSON(w, http.StatusInternalServerError, errJSON(fmt.Errorf("failed to write file: %w", writeErr)))
				return
			}
		}
		if readErr != nil {
			break
		}
	}

	// Update stickers_map.json
	stickersMapPath := filepath.Join(a.cfg.RootDir, "stickers", "stickers_map.json")
	var stickersMap map[string]map[string]string
	if err := a.loadJSON(stickersMapPath, &stickersMap); err != nil {
		// Create new map if file doesn't exist
		stickersMap = make(map[string]map[string]string)
	}

	stickersMap[safeFilename] = map[string]string{
		"description": description,
		"folder":      "uploads",
	}

	// Save updated map
	mapFile, err := os.Create(stickersMapPath)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errJSON(fmt.Errorf("failed to update stickers map: %w", err)))
		return
	}
	defer mapFile.Close()

	encoder := json.NewEncoder(mapFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(stickersMap); err != nil {
		writeJSON(w, http.StatusInternalServerError, errJSON(fmt.Errorf("failed to write stickers map: %w", err)))
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success":     true,
		"filename":    safeFilename,
		"description": description,
		"folder":      "uploads",
		"url":         fmt.Sprintf("/static/stickers/uploads/%s", safeFilename),
		"message":     "Sticker uploaded successfully",
	})
}

// handleDebugLog receives client-side debug logs and writes them to save_debug.log
func (a *App) handleDebugLog(w http.ResponseWriter, r *http.Request) {
	body, err := decodeBody(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}

	logType := stringField(body, "type", "DEBUG")
	message := stringField(body, "message", "")
	filename := stringField(body, "filename", "unknown")
	userAgent := r.UserAgent()
	clientIP := r.RemoteAddr

	// Extract additional context
	turnCount := numberField(body, "turn_count", 0)
	wordCount := numberField(body, "word_count", 0)
	queuedTurnCount := numberField(body, "queued_turn_count", -1)
	queuedWordCount := numberField(body, "queued_word_count", -1)

	// Build log entry
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	logEntry := fmt.Sprintf("[%s] [CLIENT %s] [%s] File: %s | Turns: %d | Words: %d",
		timestamp, logType, clientIP, filename, int(turnCount), int(wordCount))

	if queuedTurnCount >= 0 {
		logEntry += fmt.Sprintf(" | Queued: %d turns, %d words", int(queuedTurnCount), int(queuedWordCount))
	}

	if message != "" {
		logEntry += fmt.Sprintf(" | %s", message)
	}

	// Write to dedicated save_debug.log file
	debugLogPath := filepath.Join(a.cfg.RootDir, "save_debug.log")
	f, err := os.OpenFile(debugLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("[DEBUG LOG] Failed to open save_debug.log: %v", err)
		writeJSON(w, http.StatusInternalServerError, errJSON(err))
		return
	}
	defer f.Close()

	if _, err := f.WriteString(logEntry + "\n"); err != nil {
		log.Printf("[DEBUG LOG] Failed to write to save_debug.log: %v", err)
	}

	// Also log to main server.log for correlation
	log.Printf("[CLIENT DEBUG] %s | %s | %s", clientIP, filename, message)

	// Log user agent on first occurrence (helps identify browser issues)
	if logType == "SAVE_START" {
		if _, err := f.WriteString(fmt.Sprintf("[%s] [CLIENT INFO] UA: %s\n", timestamp, userAgent)); err != nil {
			log.Printf("[DEBUG LOG] Failed to write UA: %v", err)
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{"success": true})
}
