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
	UserID            string   `json:"user_id"`
	Username          string   `json:"username"`
	RealName          string   `json:"real_name"`
	TotalDialogues    int      `json:"total_dialogues"`
	ApprovedDialogues int      `json:"approved_dialogues"`
	TotalCharacters   int      `json:"total_characters"`
	ApprovedChars     int      `json:"approved_chars"`
	ExcellentCount    int      `json:"excellent_count"`
	Earnings          float64  `json:"earnings"`
	ApprovalRate      float64  `json:"approval_rate"`
	Characters        []string `json:"characters"`
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
}

// handleGetPaymentSummary returns the complete payment summary
func (a *App) handleGetPaymentSummary(w http.ResponseWriter, r *http.Request) {
	// Get current user session
	sessionData := a.getCurrentUser(r)
	if sessionData == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "Not logged in"})
		return
	}

	summary := a.getPaymentSummary()
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

// getPaymentSummary returns cached payment summary or calculates it if stale
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
	totalFilesLoaded := 0
	for preset, folder := range a.presetFolders {
		presetStart := time.Now()
		entries, err := os.ReadDir(folder)
		if err != nil {
			continue
		}

		filesLoaded := 0
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
			filesLoaded++

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

			// Count approved work
			if category == "passed" || category == "legacy_passed" {
				writer.ApprovedDialogues++
				writer.ApprovedChars += charCount
				totalApprovedChars += charCount
			}

			// Check for excellent cases
			if a.excellent[preset] != nil && a.excellent[preset][pf.UserID] {
				writer.ExcellentCount++
			}

			// Collect best works (passed with dialogue traits)
			if (category == "passed" || category == "legacy_passed") && dialogueTrait != "" {
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
		totalFilesLoaded += filesLoaded
		log.Printf("[PERF] calculatePaymentSummary: preset %s loaded %d files in %v", preset, filesLoaded, time.Since(presetStart))
	}

	log.Printf("[PERF] calculatePaymentSummary: total files loaded: %d", totalFilesLoaded)

	// Calculate earnings and approval rates
	calcStart := time.Now()
	var leaderboard []WriterStats
	for _, writer := range writerMap {
		writer.Earnings = float64(writer.ApprovedChars) * ratePerChar
		totalEarnings += writer.Earnings
		if writer.TotalDialogues > 0 {
			writer.ApprovalRate = float64(writer.ApprovedDialogues) / float64(writer.TotalDialogues) * 100
		}
		if writer.ApprovedChars > 0 || writer.TotalDialogues > 0 {
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

	log.Printf("[PERF] calculatePaymentSummary: calculations and sorting took %v", time.Since(calcStart))
	log.Printf("[PERF] calculatePaymentSummary: TOTAL %v (writers: %d, files: %d)", time.Since(startTime), len(leaderboard), totalFilesLoaded)

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
	summary := a.getPaymentSummary()

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

