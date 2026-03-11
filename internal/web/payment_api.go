package web

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

// WriterStats represents statistics for a single writer
type WriterStats struct {
	UserID             string   `json:"user_id"`
	Username           string   `json:"username"`
	RealName           string   `json:"real_name"`
	TotalDialogues     int      `json:"total_dialogues"`
	ApprovedDialogues  int      `json:"approved_dialogues"`
	PendingDialogues   int      `json:"pending_dialogues"`
	RejectedDialogues  int      `json:"rejected_dialogues"`
	TotalCharacters    int      `json:"total_characters"`
	ApprovedChars      int      `json:"approved_chars"`
	ExcellentCount     int      `json:"excellent_count"`
	Earnings           float64  `json:"earnings"`
	ExpectedEarnings   float64  `json:"expected_earnings"`
	ApprovalRate       float64  `json:"approval_rate"`
	Characters         []string `json:"characters"`
}

// ReviewAlert represents a notification about a reviewed dialogue
type ReviewAlert struct {
	Filename        string   `json:"filename"`
	CharacterName   string   `json:"character_name"`
	Category        string   `json:"category"`
	Timestamp       int64    `json:"timestamp"`
	RejectionReason string   `json:"rejection_reason,omitempty"`
	DialogueTrait   string   `json:"dialogue_trait,omitempty"`
	ProblemsCount   int      `json:"problems_count,omitempty"`
	Comments        []any    `json:"comments,omitempty"`
	CharCount       int      `json:"char_count"`
	Preset          string   `json:"preset"`
}

// PersonalDashboard contains the logged-in writer's personal stats and alerts
type PersonalDashboard struct {
	Username         string        `json:"username"`
	RealName         string        `json:"real_name"`
	WordsPassed      int           `json:"words_passed"`
	IncomeReceived   float64       `json:"income_received"`
	IncomeExpected   float64       `json:"income_expected"`
	Upvotes          int           `json:"upvotes"`
	TotalDialogues   int           `json:"total_dialogues"`
	ApprovedDialogues int          `json:"approved_dialogues"`
	PendingDialogues  int          `json:"pending_dialogues"`
	RejectedDialogues int          `json:"rejected_dialogues"`
	TotalCharacters   int          `json:"total_characters"`
	CharsMissing      int          `json:"chars_missing"`
	ApprovalRate      float64      `json:"approval_rate"`
	Alerts            []ReviewAlert `json:"alerts"`
	LeaderboardRank   int          `json:"leaderboard_rank"`
}

// BestWork represents an excellent dialogue sample
type BestWork struct {
	Filename      string   `json:"filename"`
	UserID        string   `json:"user_id"`
	CharacterName string   `json:"character_name"`
	RealName      string   `json:"real_name"`
	DialogueTrait string   `json:"dialogue_trait"`
	CharCount     int      `json:"char_count"`
	DialogueCount int      `json:"dialogue_count"`
	Preview       []string `json:"preview"`
	Category      string   `json:"category"`
	Preset        string   `json:"preset"`
}

// PaymentSummary contains the overall payment statistics
type PaymentSummary struct {
	TotalWriters       int          `json:"total_writers"`
	TotalApprovedChars int          `json:"total_approved_chars"`
	TotalEarnings      float64      `json:"total_earnings"`
	Leaderboard        []WriterStats `json:"leaderboard"`
	BestWorks          []BestWork    `json:"best_works"`
	PendingPayments    float64      `json:"pending_payments"`
	PaidAmount         float64      `json:"paid_amount"`
}

// RegisterPaymentAPI wires payment-related API endpoints
func (a *App) RegisterPaymentAPI(r chi.Router) {
	r.Get("/api/payment/summary", a.handleGetPaymentSummary)
	r.Get("/api/payment/leaderboard", a.handleGetLeaderboard)
	r.Get("/api/payment/best-works", a.handleGetBestWorks)
	r.Get("/api/payment/writer/{user_id}", a.handleGetWriterEarnings)
	r.Get("/api/payment/my-dashboard", a.handleGetMyDashboard)
	r.Put("/api/payment/settings", a.handleUpdatePaymentSettings)
	r.Post("/api/payment/rush-editor", a.handleRushEditor)
	r.Get("/api/payment/rush-status", a.handleGetRushStatus)
	r.Get("/api/payment/editor-nudges", a.handleGetEditorNudges)
}

// handleGetPaymentSummary returns the complete payment summary
func (a *App) handleGetPaymentSummary(w http.ResponseWriter, r *http.Request) {
	// Get current user session
	sessionData := a.getCurrentUser(r)
	if sessionData == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "Not logged in"})
		return
	}

	summary := a.getPaymentSummaryWithOverrides()
	writeJSON(w, http.StatusOK, summary)
}

// handleGetLeaderboard returns the writer leaderboard
func (a *App) handleGetLeaderboard(w http.ResponseWriter, r *http.Request) {
	sessionData := a.getCurrentUser(r)
	if sessionData == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "Not logged in"})
		return
	}

	summary := a.getPaymentSummary()
	writeJSON(w, http.StatusOK, map[string]any{
		"leaderboard": summary.Leaderboard,
		"total":       len(summary.Leaderboard),
	})
}

// handleGetBestWorks returns the best/excellent works
func (a *App) handleGetBestWorks(w http.ResponseWriter, r *http.Request) {
	sessionData := a.getCurrentUser(r)
	if sessionData == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "Not logged in"})
		return
	}

	summary := a.getPaymentSummary()
	writeJSON(w, http.StatusOK, map[string]any{
		"best_works": summary.BestWorks,
		"total":      len(summary.BestWorks),
	})
}

// handleGetWriterEarnings returns earnings for a specific writer
func (a *App) handleGetWriterEarnings(w http.ResponseWriter, r *http.Request) {
	sessionData := a.getCurrentUser(r)
	if sessionData == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "Not logged in"})
		return
	}

	userID := chi.URLParam(r, "user_id")
	summary := a.getPaymentSummary()

	for _, writer := range summary.Leaderboard {
		if writer.UserID == userID || writer.Username == userID {
			writeJSON(w, http.StatusOK, writer)
			return
		}
	}

	writeJSON(w, http.StatusNotFound, map[string]any{"error": "Writer not found"})
}

// handleUpdatePaymentSettings allows admin to set PendingPayments and PaidAmount overrides (PUT /api/payment/settings).
// Body: {"pending_payments": 123.45, "paid_amount": 678.90}. Send a key with value null to clear that override; omit a key to leave it unchanged.
func (a *App) handleUpdatePaymentSettings(w http.ResponseWriter, r *http.Request) {
	sessionData := a.getCurrentUser(r)
	if sessionData == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "Not logged in"})
		return
	}
	if sessionData.Role != "editor" {
		writeJSON(w, http.StatusForbidden, map[string]any{"error": "Admin only"})
		return
	}
	if r.Method != http.MethodPut {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"error": "Method not allowed"})
		return
	}
	var body map[string]*float64
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "Invalid JSON"})
		return
	}
	a.paymentOverridesMu.Lock()
	if v, ok := body["pending_payments"]; ok {
		a.pendingPaymentsOverride = v
	}
	if v, ok := body["paid_amount"]; ok {
		a.paidAmountOverride = v
	}
	a.paymentOverridesMu.Unlock()
	if err := a.savePaymentOverrides(a.cfg.RootDir); err != nil {
		log.Printf("savePaymentOverrides: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "Failed to save settings"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"success": true})
}
// handleGetMyDashboard returns the personal dashboard for the logged-in user
func (a *App) handleGetMyDashboard(w http.ResponseWriter, r *http.Request) {
	sessionData := a.getCurrentUser(r)
	if sessionData == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "Not logged in"})
		return
	}

	dashboard := a.calculatePersonalDashboard(sessionData)
	writeJSON(w, http.StatusOK, dashboard)
}

// handleRushEditor lets a writer "催人" — nudge editors to review their pending work.
// Writers can nudge once every 2 hours. The nudge is stored and visible to editors.
func (a *App) handleRushEditor(w http.ResponseWriter, r *http.Request) {
	sessionData := a.getCurrentUser(r)
	if sessionData == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "Not logged in"})
		return
	}

	// Parse optional message from body
	var body struct {
		Message string `json:"message"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)

	const cooldownSeconds = 7200 // 2 hours

	now := time.Now().Unix()

	// Check cooldown
	a.editorNudgesMu.RLock()
	for _, nudge := range a.editorNudges {
		if nudge.WriterID == sessionData.UserID || nudge.WriterName == sessionData.Username {
			elapsed := now - nudge.Timestamp
			if elapsed < cooldownSeconds {
				remaining := cooldownSeconds - elapsed
				a.editorNudgesMu.RUnlock()
				writeJSON(w, http.StatusTooManyRequests, map[string]any{
					"error":             "Cooldown active",
					"cooldown_remaining": remaining,
					"last_nudge":         nudge.Timestamp,
				})
				return
			}
		}
	}
	a.editorNudgesMu.RUnlock()

	// Count this writer's pending files for context
	pendingCount := 0
	realName := ""
	for _, folder := range a.presetFolders {
		entries, err := os.ReadDir(folder)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") || e.Name() == "_character_defaults.json" {
				continue
			}
			if _, ok := parseWriterFilename(e.Name()); !ok {
				continue
			}
			path := filepath.Join(folder, e.Name())
			var data map[string]any
			if err := a.loadJSON(path, &data); err != nil {
				continue
			}
			personnelID := stringField(data, "personnel_id", "")
			userName := stringField(data, "user_name", "")
			isMyFile := (personnelID != "" && personnelID == sessionData.UserID) ||
				(userName != "" && userName == sessionData.Username)
			if !isMyFile {
				continue
			}
			if realName == "" {
				realName = stringField(data, "real_name", "")
			}
			cat := stringField(data, "category", "pending")
			if cat == "pending" || cat == "existing" {
				pendingCount++
			}
		}
	}

	// Build the nudge
	nudge := EditorNudge{
		WriterID:     sessionData.UserID,
		WriterName:   sessionData.Username,
		RealName:     realName,
		Timestamp:    now,
		PendingCount: pendingCount,
		Message:      body.Message,
	}

	// Replace existing nudge for this writer or append
	a.editorNudgesMu.Lock()
	found := false
	for i, existing := range a.editorNudges {
		if existing.WriterID == sessionData.UserID || existing.WriterName == sessionData.Username {
			a.editorNudges[i] = nudge
			found = true
			break
		}
	}
	if !found {
		a.editorNudges = append(a.editorNudges, nudge)
	}
	a.editorNudgesMu.Unlock()

	// Persist
	if err := a.saveEditorNudges(a.cfg.RootDir); err != nil {
		log.Printf("saveEditorNudges: %v", err)
	}

	log.Printf("[NUDGE] Writer %s (%s) is rushing editors! Pending files: %d", sessionData.Username, realName, pendingCount)

	writeJSON(w, http.StatusOK, map[string]any{
		"success":       true,
		"pending_count": pendingCount,
		"message":       "催审成功！编辑已收到通知。/ Editors have been nudged!",
		"cooldown":      cooldownSeconds,
	})
}

// handleGetRushStatus returns whether the current writer has an active nudge and cooldown info
func (a *App) handleGetRushStatus(w http.ResponseWriter, r *http.Request) {
	sessionData := a.getCurrentUser(r)
	if sessionData == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "Not logged in"})
		return
	}

	now := time.Now().Unix()
	const cooldownSeconds int64 = 7200

	a.editorNudgesMu.RLock()
	defer a.editorNudgesMu.RUnlock()

	for _, nudge := range a.editorNudges {
		if nudge.WriterID == sessionData.UserID || nudge.WriterName == sessionData.Username {
			elapsed := now - nudge.Timestamp
			if elapsed < cooldownSeconds {
				writeJSON(w, http.StatusOK, map[string]any{
					"has_nudged":         true,
					"cooldown_remaining": cooldownSeconds - elapsed,
					"last_nudge":         nudge.Timestamp,
				})
				return
			}
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"has_nudged":         false,
		"cooldown_remaining": 0,
	})
}

// handleGetEditorNudges returns all active nudges (for editors to see who's waiting)
func (a *App) handleGetEditorNudges(w http.ResponseWriter, r *http.Request) {
	sessionData := a.getCurrentUser(r)
	if sessionData == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "Not logged in"})
		return
	}

	// Only editors can see all nudges
	if sessionData.Role != "editor" {
		writeJSON(w, http.StatusForbidden, map[string]any{"error": "Editor access required"})
		return
	}

	now := time.Now().Unix()
	const maxAge int64 = 86400 // Show nudges from the last 24 hours

	a.editorNudgesMu.RLock()
	var active []EditorNudge
	for _, nudge := range a.editorNudges {
		if now-nudge.Timestamp < maxAge {
			active = append(active, nudge)
		}
	}
	a.editorNudgesMu.RUnlock()

	// Sort by most recent first
	sort.Slice(active, func(i, j int) bool {
		return active[i].Timestamp > active[j].Timestamp
	})

	writeJSON(w, http.StatusOK, map[string]any{
		"nudges": active,
		"total":  len(active),
	})
}

// calculatePersonalDashboard scans all files to build the current user's personal dashboard
func (a *App) calculatePersonalDashboard(sessionData *SessionData) *PersonalDashboard {
	ratePerChar := 0.08

	dashboard := &PersonalDashboard{
		Username: sessionData.Username,
		Alerts:   []ReviewAlert{},
	}

	// Scan all preset folders for files belonging to this user
	for preset, folder := range a.presetFolders {
		entries, err := os.ReadDir(folder)
		if err != nil {
			continue
		}

		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
				continue
			}
			if e.Name() == "_character_defaults.json" {
				continue
			}

			pf, ok := parseWriterFilename(e.Name())
			if !ok {
				continue
			}

			path := filepath.Join(folder, e.Name())
			var data map[string]any
			if err := a.loadJSON(path, &data); err != nil {
				continue
			}

			// Match by personnel_id (primary) or user_name (fallback)
			personnelID := stringField(data, "personnel_id", "")
			userName := stringField(data, "user_name", "")
			realName := stringField(data, "real_name", "")

			isMyFile := false
			if personnelID != "" && personnelID == sessionData.UserID {
				isMyFile = true
			} else if userName != "" && userName == sessionData.Username {
				isMyFile = true
			}
			if !isMyFile {
				continue
			}

			// Set display name from first matched file
			if dashboard.RealName == "" && realName != "" {
				dashboard.RealName = realName
			}

			category := stringField(data, "category", "pending")
			dialogue, _ := data["dialogue"].([]any)
			charCount := countDialogueHanChars(dialogue)
			dialogueTrait := stringField(data, "dialogue_trait", "")
			charName := stringField(data, "character_name", "")
			rejectionReason := stringField(data, "rejection_reason", "")
			problemsCount := numberField(data, "problems_count", 0)
			assessmentTS := int64(numberField(data, "assessment_timestamp", 0))
			comments, _ := data["comments"].([]any)

			dashboard.TotalDialogues++
			dashboard.TotalCharacters += charCount

			switch category {
			case "passed":
				dashboard.ApprovedDialogues++
				dashboard.WordsPassed += charCount
			case "rejected":
				dashboard.RejectedDialogues++
			case "pending", "existing":
				dashboard.PendingDialogues++
			}

			if a.isPublicUser(preset, pf.UserID) {
				dashboard.Upvotes++
			}

			// Collect alerts for reviewed files (passed or rejected)
			if category == "passed" || category == "rejected" {
				alert := ReviewAlert{
					Filename:        e.Name(),
					CharacterName:   charName,
					Category:        category,
					Timestamp:       assessmentTS,
					RejectionReason: rejectionReason,
					DialogueTrait:   dialogueTrait,
					ProblemsCount:   problemsCount,
					CharCount:       charCount,
					Preset:          preset,
				}
				if len(comments) > 0 {
					alert.Comments = comments
				}
				dashboard.Alerts = append(dashboard.Alerts, alert)
			}
		}
	}

	// Calculate derived fields
	dashboard.IncomeReceived = float64(dashboard.WordsPassed) * ratePerChar
	dashboard.IncomeExpected = float64(dashboard.TotalCharacters) * ratePerChar
	dashboard.CharsMissing = dashboard.TotalCharacters - dashboard.WordsPassed
	if dashboard.TotalDialogues > 0 {
		dashboard.ApprovalRate = float64(dashboard.ApprovedDialogues) / float64(dashboard.TotalDialogues) * 100
	}

	// Sort alerts by timestamp descending (most recent first)
	sort.Slice(dashboard.Alerts, func(i, j int) bool {
		return dashboard.Alerts[i].Timestamp > dashboard.Alerts[j].Timestamp
	})

	// Limit alerts to 20 most recent
	if len(dashboard.Alerts) > 20 {
		dashboard.Alerts = dashboard.Alerts[:20]
	}

	// Determine leaderboard rank
	summary := a.getPaymentSummary()
	for i, writer := range summary.Leaderboard {
		if writer.Username == sessionData.Username || writer.RealName == sessionData.Username {
			dashboard.LeaderboardRank = i + 1
			break
		}
	}

	return dashboard
}

// Cache TTL is 5 minutes to balance freshness with performance
func (a *App) getPaymentSummary() *PaymentSummary {
	const cacheTTL = 5 * time.Minute
	
	// Try to read from cache first
	a.paymentSummaryMutex.RLock()
	if a.paymentSummaryCache != nil && time.Since(a.paymentSummaryCacheAt) < cacheTTL {
		result := a.paymentSummaryCache
		a.paymentSummaryMutex.RUnlock()
		log.Printf("[PERF] getPaymentSummary: returned cached result (age: %v)", time.Since(a.paymentSummaryCacheAt))
		return result
	}
	a.paymentSummaryMutex.RUnlock()
	
	// Cache miss or stale - calculate fresh
	result := a.calculatePaymentSummary()
	
	// Update cache
	a.paymentSummaryMutex.Lock()
	a.paymentSummaryCache = result
	a.paymentSummaryCacheAt = time.Now()
	a.paymentSummaryMutex.Unlock()
	
	return result
}

// getPaymentSummaryWithOverrides returns a copy of the payment summary with admin overrides applied for PendingPayments and PaidAmount.
func (a *App) getPaymentSummaryWithOverrides() *PaymentSummary {
	summary := a.getPaymentSummary()
	a.paymentOverridesMu.RLock()
	pending := a.pendingPaymentsOverride
	paid := a.paidAmountOverride
	a.paymentOverridesMu.RUnlock()
	if pending == nil && paid == nil {
		return summary
	}
	// Copy so we don't mutate the cache
	out := *summary
	if pending != nil {
		out.PendingPayments = *pending
	}
	if paid != nil {
		out.PaidAmount = *paid
	}
	return &out
}

// calculatePaymentSummary aggregates data across all presets
func (a *App) calculatePaymentSummary() *PaymentSummary {
	startTime := time.Now()
	log.Printf("[PERF] calculatePaymentSummary: START (fresh calculation)")
	writerMap := make(map[string]*WriterStats)
	var bestWorks []BestWork
	totalApprovedChars := 0
	totalEarnings := 0.0

	// Rate per character (in CNY)
	ratePerChar := 0.08

	// Iterate through all preset folders
	for preset, folder := range a.presetFolders {
		entries, err := os.ReadDir(folder)
		if err != nil {
			continue
		}

		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
				continue
			}
			if e.Name() == "_character_defaults.json" {
				continue
			}

			pf, ok := parseWriterFilename(e.Name())
			if !ok {
				continue
			}

			path := filepath.Join(folder, e.Name())
			var data map[string]any
			if err := a.loadJSON(path, &data); err != nil {
				continue
			}

			// Extract data
			category := stringField(data, "category", "pending")
			dialogue, _ := data["dialogue"].([]any)
			charCount := countDialogueHanChars(dialogue)
			dialogueTrait := stringField(data, "dialogue_trait", "")
			userName := stringField(data, "user_name", "Unknown")
			realName := stringField(data, "real_name", "")

			// Get or create writer stats
			writerKey := pf.UserID
			if _, exists := writerMap[writerKey]; !exists {
				writerMap[writerKey] = &WriterStats{
					UserID:     pf.UserID,
					Username:   userName,
					RealName:   realName,
					Characters: []string{},
				}
			}
			writer := writerMap[writerKey]

			// Update totals
			writer.TotalDialogues++
			writer.TotalCharacters += charCount

			// Track which characters this writer has worked on
			charName := stringField(data, "character_name", "")
			if charName != "" {
				found := false
				for _, c := range writer.Characters {
					if c == charName {
						found = true
						break
					}
				}
				if !found {
					writer.Characters = append(writer.Characters, charName)
				}
			}

			// Count by category
			switch category {
			case "passed":
				writer.ApprovedDialogues++
				writer.ApprovedChars += charCount
				totalApprovedChars += charCount
			case "rejected":
				writer.RejectedDialogues++
			case "pending", "existing":
				writer.PendingDialogues++
			}

			if a.isPublicUser(preset, pf.UserID) {
				writer.ExcellentCount++
			}

			// Collect best works (passed with dialogue traits)
			if (category == "passed" ) && dialogueTrait != "" {
				preview := getDialoguePreview(dialogue, 3)
				bestWorks = append(bestWorks, BestWork{
					Filename:      e.Name(),
					UserID:        pf.UserID,
					CharacterName: userName,
					RealName:      realName,
					DialogueTrait: dialogueTrait,
					CharCount:     charCount,
					DialogueCount: len(dialogue),
					Preview:       preview,
					Category:      category,
					Preset:        preset,
				})
			}
		}
	}

	// Calculate earnings and approval rates
	var leaderboard []WriterStats
	for _, writer := range writerMap {
		writer.Earnings = float64(writer.ApprovedChars) * ratePerChar
		writer.ExpectedEarnings = float64(writer.TotalCharacters) * ratePerChar
		totalEarnings += writer.Earnings
		if writer.TotalDialogues > 0 {
			writer.ApprovalRate = float64(writer.ApprovedDialogues) / float64(writer.TotalDialogues) * 100
		}
		if writer.ApprovedChars > 0 {
			leaderboard = append(leaderboard, *writer)
		}
	}

	// Sort leaderboard by approved characters (descending)
	sort.Slice(leaderboard, func(i, j int) bool {
		return leaderboard[i].ApprovedChars > leaderboard[j].ApprovedChars
	})

	// Sort best works by char count (descending) and limit to top 10
	sort.Slice(bestWorks, func(i, j int) bool {
		return bestWorks[i].CharCount > bestWorks[j].CharCount
	})
	if len(bestWorks) > 10 {
		bestWorks = bestWorks[:10]
	}

	log.Printf("[PERF] calculatePaymentSummary: TOTAL %v (writers: %d)", time.Since(startTime), len(leaderboard))

	return &PaymentSummary{
		TotalWriters:       len(leaderboard),
		TotalApprovedChars: totalApprovedChars,
		TotalEarnings:      totalEarnings,
		Leaderboard:        leaderboard,
		BestWorks:          bestWorks,
		PendingPayments:    totalEarnings * 0.3, // Example: 30% pending
		PaidAmount:         totalEarnings * 0.7, // Example: 70% paid
	}
}

// getDialoguePreview extracts the first n dialogue entries as preview
func getDialoguePreview(dialogue []any, n int) []string {
	var preview []string
	for i, d := range dialogue {
		if i >= n {
			break
		}
		if dm, ok := d.(map[string]any); ok {
			role := stringField(dm, "role", "")
			content := stringField(dm, "content", "")
			if content != "" {
				if len(content) > 80 {
					content = content[:80] + "..."
				}
				preview = append(preview, role+": "+content)
			}
		}
	}
	return preview
}

// PaymentTemplateData extends TemplateData with payment-specific fields
type PaymentTemplateData struct {
	TemplateData
	Summary *PaymentSummary
}

// handlePaymentWithData renders the payment page with summary data
func (a *App) handlePaymentWithData(w http.ResponseWriter, r *http.Request) {
	data := a.getTemplateData(r)
	summary := a.getPaymentSummaryWithOverrides()

	// Convert to JSON for embedding in template
	summaryJSON, err := json.Marshal(summary)
	if err != nil {
		summaryJSON = []byte("{}")
	}

	paymentData := map[string]any{
		"LoggedIn":    data.LoggedIn,
		"CurrentUser": data.CurrentUser,
		"UserRole":    data.UserRole,
		"Characters":  data.Characters,
		"SummaryJSON": string(summaryJSON),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	
	// For now, we'll inject the data as a script variable
	// The template will use JavaScript to render the dynamic content
	_ = paymentData
}

