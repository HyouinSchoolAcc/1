package web

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"data_labler_ui_go/internal/database"
	"data_labler_ui_go/internal/email"
)

// App holds server state and helpers.
type App struct {
	cfg               Config
	presetFolders     map[string]string
	publicUsers       map[string]map[string]bool
	publicUsersMu     sync.RWMutex
	fileLocks         sync.Map // path -> *sync.Mutex
	authService        *AuthService
	editorRateLimiter  *EditorRateLimiter
	loungeStore        *LoungeStore
	characterStore    *CharacterStore
	registrationStore *database.RegistrationStore
	emailService      *email.Service
	db                *database.DB
	llmService        *LLMService
	tutorialStore     *database.TutorialStore
	chatStore         *database.ChatStore
	profileStore      *database.ProfileStore
	
	// Payment summary cache (expensive operation)
	paymentSummaryCache   *PaymentSummary
	paymentSummaryCacheAt time.Time
	paymentSummaryMutex   sync.RWMutex

	// Admin-editable payment overrides (only admin can set; persisted to data/payment_overrides.json)
	paymentOverridesMu      sync.RWMutex
	pendingPaymentsOverride *float64
	paidAmountOverride      *float64

	// Video generation job tracking (async)
	videoJobs   map[string]*VideoJob
	videoJobsMu sync.RWMutex

	// Editor nudge/rush tracking (persisted to data/editor_nudges.json)
	editorNudgesMu sync.RWMutex
	editorNudges   []EditorNudge
}

// VideoJob tracks the status of an async video generation.
type VideoJob struct {
	Status    string    `json:"status"`
	Message   string    `json:"message,omitempty"`
	Error     string    `json:"error,omitempty"`
	SizeMB    string    `json:"size_mb,omitempty"`
	Filename  string    `json:"filename,omitempty"`
	StartedAt time.Time `json:"started_at"`
}

func NewApp(cfg Config) *App {
	presets := []string{
		"presets_kurisu",
		"presets_kurisu_CN",
		"presets_lin_lu_CN",
		"presets_lin_lu",
		"presets_newcharacter_1",
	}
	presetFolders := make(map[string]string, len(presets))
	for _, p := range presets {
		presetFolders[p] = filepath.Join(cfg.PresetBaseDir, p)
		_ = os.MkdirAll(presetFolders[p], 0o755)
	}

	publicUsers := make(map[string]map[string]bool, len(presets))
	defaults := map[string]map[string]bool{
		"presets_kurisu_CN": setOf("user_32", "user_33"),
		"presets_lin_lu_CN": setOf("user_1", "user_4"),
	}
	for _, p := range presets {
		publicUsers[p] = loadPublicUsersFile(filepath.Join(presetFolders[p], "public_users.json"), defaults[p])
	}

	// Initialize database
	db, err := database.New(filepath.Join(cfg.RootDir, "data"))
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize user store with database
	userStore := database.NewUserStore(db)

	// Initialize auth service with SQL-based user store
	authService := NewAuthService(userStore)

	// Initialize lounge store with database
	loungeStore := database.NewLoungeStore(db)

	// Initialize character store with database
	characterStore := database.NewCharacterStore(db)

	// Initialize registration store
	registrationStore := database.NewRegistrationStore(db)

	// Initialize tutorial store
	tutorialStore := database.NewTutorialStore(db)

	// Initialize chat store for Discord-like features
	chatStore := database.NewChatStore(db)
	// Ensure default #general channel exists
	if err := chatStore.EnsureDefaultChannel(); err != nil {
		log.Printf("Warning: Failed to ensure default channel: %v", err)
	}

	// Initialize email service
	emailConfig := email.Config{
		SMTPHost:     "smtp.gmail.com",
		SMTPPort:     "587",
		SMTPUsername: "divtwop@gmail.com",
		SMTPPassword: "jimmy123!",
		FromEmail:    "divtwop@gmail.com",
		FromName:     "Divergence 2% Writer Portal",
	}
	emailService := email.NewService(emailConfig)

	app := &App{
		cfg:                cfg,
		presetFolders:      presetFolders,
		publicUsers:        publicUsers,
		authService:        authService,
		editorRateLimiter:  NewEditorRateLimiter(),
		loungeStore:        loungeStore,
		characterStore:    characterStore,
		registrationStore: registrationStore,
		emailService:      emailService,
		db:                db,
		tutorialStore:     tutorialStore,
		chatStore:         chatStore,
		profileStore:      database.NewProfileStore(db),
		videoJobs:          make(map[string]*VideoJob),
	}
	app.loadPaymentOverrides(cfg.RootDir)
	app.loadEditorNudges(cfg.RootDir)
	return app
}

// paymentOverridesFile holds admin-editable payment values for persistence
type paymentOverridesFile struct {
	PendingPayments *float64 `json:"pending_payments,omitempty"`
	PaidAmount      *float64 `json:"paid_amount,omitempty"`
}

func (a *App) loadPaymentOverrides(rootDir string) {
	path := filepath.Join(rootDir, "data", "payment_overrides.json")
	b, err := os.ReadFile(path)
	if err != nil {
		return
	}
	var f paymentOverridesFile
	if err := json.Unmarshal(b, &f); err != nil {
		return
	}
	a.paymentOverridesMu.Lock()
	a.pendingPaymentsOverride = f.PendingPayments
	a.paidAmountOverride = f.PaidAmount
	a.paymentOverridesMu.Unlock()
}

func (a *App) savePaymentOverrides(rootDir string) error {
	a.paymentOverridesMu.RLock()
	f := paymentOverridesFile{
		PendingPayments: a.pendingPaymentsOverride,
		PaidAmount:      a.paidAmountOverride,
	}
	a.paymentOverridesMu.RUnlock()
	path := filepath.Join(rootDir, "data", "payment_overrides.json")
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

// EditorNudge represents a single "rush the editor" nudge from a writer
type EditorNudge struct {
	WriterID     string `json:"writer_id"`
	WriterName   string `json:"writer_name"`
	RealName     string `json:"real_name,omitempty"`
	Timestamp    int64  `json:"timestamp"`
	PendingCount int    `json:"pending_count"`
	Message      string `json:"message,omitempty"`
}

type editorNudgesFile struct {
	Nudges []EditorNudge `json:"nudges"`
}

func (a *App) loadEditorNudges(rootDir string) {
	path := filepath.Join(rootDir, "data", "editor_nudges.json")
	b, err := os.ReadFile(path)
	if err != nil {
		return
	}
	var f editorNudgesFile
	if err := json.Unmarshal(b, &f); err != nil {
		return
	}
	a.editorNudgesMu.Lock()
	a.editorNudges = f.Nudges
	a.editorNudgesMu.Unlock()
}

func (a *App) saveEditorNudges(rootDir string) error {
	a.editorNudgesMu.RLock()
	f := editorNudgesFile{Nudges: a.editorNudges}
	a.editorNudgesMu.RUnlock()
	path := filepath.Join(rootDir, "data", "editor_nudges.json")
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

// Close closes the database connection
func (a *App) Close() error {
	if a.db != nil {
		return a.db.Close()
	}
	return nil
}

// StartTempCleanup starts a background goroutine that cleans up temp characters every hour
func (a *App) StartTempCleanup() {
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		
		// Run cleanup immediately on start, then every hour
		a.cleanupTempCharacters()
		
		for range ticker.C {
			a.cleanupTempCharacters()
		}
	}()
	log.Println("✅Temp character cleanup routine started (runs every 1 hour)")
}

// cleanupTempCharacters removes temp character files older than 1 hour
func (a *App) cleanupTempCharacters() {
	log.Println("Running temp character cleanup...")
	totalDeleted := 0
	cutoffTime := time.Now().Add(-1 * time.Hour).Unix()
	
	for presetName, folder := range a.presetFolders {
		entries, err := os.ReadDir(folder)
		if err != nil {
			continue
		}
		
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			name := e.Name()
			// Only process temp character files
			if !strings.HasPrefix(name, "temp_") || !strings.HasSuffix(name, ".json") {
				continue
			}
			
			path := filepath.Join(folder, name)
			
			// Load file to check creation time
			var data map[string]any
			if err := a.loadJSON(path, &data); err != nil {
				continue
			}
			
			// Check creation metadata
			meta, ok := data["_creation_metadata"].(map[string]any)
			if !ok {
				// No metadata, delete if file is old based on file mod time
				info, err := e.Info()
				if err != nil {
					continue
				}
				if info.ModTime().Unix() < cutoffTime {
					os.Remove(path)
					totalDeleted++
					log.Printf("  Deleted old temp file (no metadata): %s", name)
				}
				continue
			}
			
			// Check created_at timestamp
			createdAt, ok := meta["created_at"].(float64)
			if !ok {
				continue
			}
			
			if int64(createdAt) < cutoffTime {
				os.Remove(path)
				totalDeleted++
				log.Printf("  Deleted temp file: %s (preset: %s)", name, presetName)
			}
		}
	}
	
	if totalDeleted > 0 {
		log.Printf("Temp cleanup complete: deleted %d files", totalDeleted)
	} else {
		log.Println("Temp cleanup complete: no files to delete")
	}
}

func setOf(items ...string) map[string]bool {
	m := make(map[string]bool, len(items))
	for _, v := range items {
		m[v] = true
	}
	return m
}

// === ownership checks ===

// checkCharacterOwnership checks if the current user can edit a character based on the file data.
// Returns true if:
//   - user is an editor (editors can edit everything)
//   - data has no personnel_id (unassigned characters are accessible to all writers)
//   - data's personnel_id matches the session user ID
func (a *App) checkCharacterOwnership(sessionData *SessionData, data map[string]any) bool {
	if sessionData == nil {
		return false
	}
	if sessionData.Role == RoleEditor {
		return true
	}

	personnelID := stringField(data, "personnel_id", "")
	if personnelID == "" {
		return true // No personnel assigned - allow access (do nothing)
	}

	return personnelID == sessionData.UserID
}

// canEditFile checks if the current user can edit a specific file.
// Loads the file and checks personnel_id ownership.
func (a *App) canEditFile(sessionData *SessionData, filePath string) bool {
	if sessionData == nil {
		return false
	}
	if sessionData.Role == RoleEditor {
		return true
	}

	var data map[string]any
	if err := a.loadJSON(filePath, &data); err != nil {
		return true // If can't load, let the handler deal with the error
	}

	return a.checkCharacterOwnership(sessionData, data)
}

// checkUserOwnership checks if the current user owns a character (by user_id) in a preset.
// Checks the first file found for the user_id.
func (a *App) checkUserOwnership(sessionData *SessionData, userID string, preset string) bool {
	if sessionData == nil {
		return false
	}
	if sessionData.Role == RoleEditor {
		return true
	}

	folder, err := a.getPresetFolder(preset)
	if err != nil {
		return false
	}

	entries, _ := os.ReadDir(folder)
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		pf, ok := parseWriterFilename(e.Name())
		if !ok || pf.UserID != userID {
			continue
		}

		var data map[string]any
		if err := a.loadJSON(filepath.Join(folder, e.Name()), &data); err != nil {
			continue
		}

		personnelID := stringField(data, "personnel_id", "")
		if personnelID == "" {
			return true // Unassigned - allow access
		}
		return personnelID == sessionData.UserID
	}

	return true // No files found - allow (handler will deal with it)
}

// setPersonnelForUser sets the personnel_id and personnel_username on all files for a user_id in a preset.
func (a *App) setPersonnelForUser(userID, personnelID, personnelUsername, preset string) (int, error) {
	folder, err := a.getPresetFolder(preset)
	if err != nil {
		return 0, err
	}

	entries, _ := os.ReadDir(folder)
	updated := 0
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		pf, ok := parseWriterFilename(e.Name())
		if !ok || pf.UserID != userID {
			continue
		}

		path := filepath.Join(folder, e.Name())
		var data map[string]any
		if err := a.loadJSON(path, &data); err != nil {
			continue
		}

		data["personnel_id"] = personnelID
		data["personnel_username"] = personnelUsername
		if err := a.atomicWriteJSON(path, data); err == nil {
			updated++
		}
	}

	return updated, nil
}

// === locking and atomic writes ===

func (a *App) getLock(path string) *sync.Mutex {
	if val, ok := a.fileLocks.Load(path); ok {
		return val.(*sync.Mutex)
	}
	mu := &sync.Mutex{}
	actual, _ := a.fileLocks.LoadOrStore(path, mu)
	return actual.(*sync.Mutex)
}

func (a *App) atomicWriteJSON(path string, data any) error {
	lock := a.getLock(path)
	lock.Lock()
	defer lock.Unlock()

	lockFile := path + ".lock"
	start := time.Now()
	for {
		fd, err := os.OpenFile(lockFile, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0o644)
		if err == nil {
			fd.Close()
			break
		}
		if time.Since(start) > 10*time.Second {
			return fmt.Errorf("timeout acquiring lock for %s", path)
		}
		time.Sleep(50 * time.Millisecond)
	}
	defer os.Remove(lockFile)

	// backup existing
	if _, err := os.Stat(path); err == nil {
		backupDir, backupName := a.backupRelPaths(path)
		if err := os.MkdirAll(backupDir, 0o755); err == nil {
			backupPath := filepath.Join(backupDir, backupName)
			if err := copyFile(path, backupPath); err == nil {
				log.Printf("[BACKUP] Created backup: %s", backupPath)
			} else {
				log.Printf("[BACKUP] WARNING: Failed to create backup %s: %v", backupPath, err)
			}
		}
	}

	tmp := path + ".tmp"
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	if err := enc.Encode(data); err != nil {
		f.Close()
		return err
	}
	if err := f.Sync(); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

// backupRelPaths mirrors python logic
func (a *App) backupRelPaths(original string) (string, string) {
	baseDir := a.cfg.RootDir
	rel, err := filepath.Rel(baseDir, original)
	if err != nil {
		rel = filepath.Base(original)
	}
	relDir := filepath.Dir(rel)
	base := filepath.Base(original)
	ts := time.Now().UTC().Format("20060102T150405000Z")
	backupDir := filepath.Join(a.cfg.BackupRootDir, relDir)
	backupName := fmt.Sprintf("%s.%s.bak", base, ts)
	return backupDir, backupName
}

// === helpers ===

var filenameRe = regexp.MustCompile(`((?:user|temp)_\d+)_Day(\d+)_dup_(\d+)_simplified\.json$`)

type ParsedFilename struct {
	UserID string
	DayStr string
	DayNum int
	DupID  string
	Name   string
}

func parseWriterFilename(name string) (*ParsedFilename, bool) {
	m := filenameRe.FindStringSubmatch(name)
	if len(m) != 4 {
		return nil, false
	}
	dayNum, _ := strconv.Atoi(m[2])
	return &ParsedFilename{
		UserID: m[1],
		DayStr: "Day" + m[2],
		DayNum: dayNum,
		DupID:  "dup_" + m[3],
		Name:   name,
	}, true
}

func (a *App) getPresetFolder(preset string) (string, error) {
	if p, ok := a.presetFolders[preset]; ok {
		return p, nil
	}
	return "", fmt.Errorf("preset folder '%s' not found", preset)
}

func (a *App) loadJSON(path string, v any) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

func (a *App) listPresetFiles(preset string) ([]fs.DirEntry, string, error) {
	folder, err := a.getPresetFolder(preset)
	if err != nil {
		return nil, "", err
	}
	entries, err := os.ReadDir(folder)
	return entries, folder, err
}

// hidden users file can be array or {hidden_users:[]}
func (a *App) loadHiddenUsers(preset string) map[string]bool {
	folder, err := a.getPresetFolder(preset)
	if err != nil {
		return nil
	}
	path := filepath.Join(folder, "hidden_users.json")
	b, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var raw any
	if err := json.Unmarshal(b, &raw); err != nil {
		return nil
	}
	out := map[string]bool{}
	switch val := raw.(type) {
	case []any:
		for _, v := range val {
			if s, ok := v.(string); ok {
				out[s] = true
			}
		}
	case map[string]any:
		if arr, ok := val["hidden_users"].([]any); ok {
			for _, v := range arr {
				if s, ok := v.(string); ok {
					out[s] = true
				}
			}
		}
	}
	return out
}

// loadPublicUsersFile reads public_users.json; falls back to defaults if missing.
func loadPublicUsersFile(path string, defaults map[string]bool) map[string]bool {
	b, err := os.ReadFile(path)
	if err != nil {
		if defaults != nil {
			return defaults
		}
		return map[string]bool{}
	}
	var arr []string
	if err := json.Unmarshal(b, &arr); err != nil {
		if defaults != nil {
			return defaults
		}
		return map[string]bool{}
	}
	out := make(map[string]bool, len(arr))
	for _, v := range arr {
		out[v] = true
	}
	return out
}

func (a *App) savePublicUsers(preset string) error {
	folder, err := a.getPresetFolder(preset)
	if err != nil {
		return err
	}
	a.publicUsersMu.RLock()
	users := a.publicUsers[preset]
	a.publicUsersMu.RUnlock()

	arr := make([]string, 0, len(users))
	for u := range users {
		if users[u] {
			arr = append(arr, u)
		}
	}
	sort.Strings(arr)
	b, err := json.MarshalIndent(arr, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(folder, "public_users.json"), b, 0o644)
}

func (a *App) isPublicUser(preset, userID string) bool {
	a.publicUsersMu.RLock()
	defer a.publicUsersMu.RUnlock()
	return a.publicUsers[preset][userID]
}

func (a *App) togglePublicUser(preset, userID string) (bool, error) {
	a.publicUsersMu.Lock()
	if a.publicUsers[preset] == nil {
		a.publicUsers[preset] = map[string]bool{}
	}
	newVal := !a.publicUsers[preset][userID]
	a.publicUsers[preset][userID] = newVal
	a.publicUsersMu.Unlock()
	return newVal, a.savePublicUsers(preset)
}

// === character & schedule helpers ===

func getCurrentCharacterRole(preset string) string {
	l := strings.ToLower(preset)
	switch {
	case strings.Contains(l, "kurisu"):
		return "Kurisu"
	case strings.Contains(l, "lin_lu"):
		return "lin_lu"
	default:
		return ""
	}
}

func getCurrentCharacterDisplayName(preset string) string {
	l := strings.ToLower(preset)
	switch {
	case strings.Contains(l, "kurisu"):
		return "Kurisu"
	case strings.Contains(l, "lin_lu") && strings.Contains(preset, "_CN"):
		return "林路"
	case strings.Contains(l, "lin_lu"):
		return "Lin Lu"
	case strings.Contains(l, "林路"):
		return "林路"
	default:
		return ""
	}
}

// default schedule fallback -- only used when character_profiles.json has no entry for the day.
// Uses English keys consistently (matching the rest of the codebase and the frontend).
// Chinese/English content is selected based on the preset name.
func fallbackSchedule(day int, preset string) map[string]any {
	role := getCurrentCharacterRole(preset)
	isCN := strings.Contains(preset, "_CN")
	switch role {
	case "Kurisu":
		if isCN {
			return map[string]any{"day": day, "morning": "启动加速器，校准设备参数，准备实验。", "afternoon": "进行时间旅行理论研究，分析实验数据。", "evening": "通过阅读科学论文或观看轻松娱乐来放松身心。"}
		}
		return map[string]any{"day": day, "morning": "Start the accelerator, calibrate the equipment, and prepare for experiments.", "afternoon": "Conduct research on time travel theories and analyze experimental data.", "evening": "Unwind by reviewing scientific papers or watching some light entertainment."}
	case "lin_lu", "林路":
		if isCN {
			return map[string]any{"day": day, "morning": "备课，准备当天的文学史课程资料。", "afternoon": "给学生上课，讲解唐诗宋词鉴赏。", "evening": "写自己的历史小说，或者在学校大銀杏树下看云放松。"}
		}
		return map[string]any{"day": day, "morning": "Prepare for classes and organize course materials for literature history.", "afternoon": "Teach students about Tang poetry and Song lyrics appreciation.", "evening": "Write my historical novel or relax under the big ginkgo tree at school."}
	default:
		return map[string]any{"day": day}
	}
}

func (a *App) getUniversalCharacterSchedule(day int, preset string) map[string]any {
	// Use the centralized character_profiles.json as the single source of truth.
	// This is the same data source used by https://wl2.studio/descriptions.
	// The old preset-local default_schedule.json files are deprecated.
	schedule := a.getCharacterScheduleForDay(preset, day)
	if len(schedule) > 1 {
		return schedule
	}
	return fallbackSchedule(day, preset)
}

func intVal(v any) int {
	switch t := v.(type) {
	case float64:
		return int(t)
	case int:
		return t
	case json.Number:
		i, _ := t.Int64()
		return int(i)
	}
	return 0
}


// load user data and merge universal schedule
func (a *App) loadUserDataWithSchedule(path, preset string) (map[string]any, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var user map[string]any
	if err := json.Unmarshal(b, &user); err != nil {
		return nil, err
	}
	day := intVal(getNested(user, "character_schedule", "day"))
	if day == 0 {
		day = intVal(getNested(user, "user_schedule", "day"))
	}
	if day == 0 {
		if pf, ok := parseWriterFilename(filepath.Base(path)); ok {
			day = pf.DayNum
		} else {
			day = 1
		}
	}
	user["character_schedule"] = a.getUniversalCharacterSchedule(day, preset)
	return user, nil
}

func getNested(m map[string]any, keys ...string) any {
	cur := any(m)
	for _, k := range keys {
		if mm, ok := cur.(map[string]any); ok {
			cur = mm[k]
		} else {
			return nil
		}
	}
	return cur
}

// === char counters ===

func isHan(ch rune) bool {
	code := ch
	return (code >= 0x4E00 && code <= 0x9FFF) ||
		(code >= 0x3400 && code <= 0x4DBF) ||
		(code >= 0xF900 && code <= 0xFAFF) ||
		(code >= 0x20000 && code <= 0x2A6DF) ||
		(code >= 0x2A700 && code <= 0x2B73F) ||
		(code >= 0x2B740 && code <= 0x2B81F) ||
		(code >= 0x2B820 && code <= 0x2CEAF)
}

func countHanChars(text string) int {
	n := 0
	for _, r := range text {
		if isHan(r) {
			n++
		}
	}
	return n
}

func countDialogueHanChars(dialogue any) int {
	arr, ok := dialogue.([]any)
	if !ok {
		return 0
	}
	total := 0
	for _, item := range arr {
		switch v := item.(type) {
		case map[string]any:
			s := firstString(v, "content", "text", "message", "utterance")
			total += countHanChars(s)
		case string:
			total += countHanChars(v)
		default:
			total += countHanChars(fmt.Sprint(v))
		}
	}
	return total
}

func firstString(m map[string]any, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k]; ok {
			if s, ok := v.(string); ok {
				return s
			}
		}
	}
	return ""
}

// === filename helpers ===

func getUserInfoFilename(preset string) string {
	switch preset {
	case "presets_kurisu_CN":
		return "new_user_info_kurisu_cn.json"
	case "presets_kurisu":
		return "new_user_info_kurisu.json"
	case "presets_lin_lu":
		return "new_user_info_lin_lu.json"
	case "presets_lin_lu_CN":
		return "new_user_info_lin_lu_cn.json"
	default:
		return "new_user_info.json"
	}
}

// zip helper
func zipBuffers(files map[string][]byte) ([]byte, error) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for name, content := range files {
		w, err := zw.Create(name)
		if err != nil {
			return nil, err
		}
		if _, err := w.Write(content); err != nil {
			return nil, err
		}
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// utility to pick preset by file search fallback
func (a *App) findFileAcrossPresets(filename string, preset string) (string, string, error) {
	if preset != "" {
		folder, err := a.getPresetFolder(preset)
		if err == nil {
			path := filepath.Join(folder, filename)
			if _, err := os.Stat(path); err == nil {
				return path, preset, nil
			}
		}
	}
	for ps, folder := range a.presetFolders {
		path := filepath.Join(folder, filename)
		if _, err := os.Stat(path); err == nil {
			return path, ps, nil
		}
	}
	return "", "", os.ErrNotExist
}

// convenience error wrapper for consistent messages
func errJSON(err error) map[string]string {
	return map[string]string{"error": err.Error()}
}

// helper: decode JSON body to map
func decodeBody(r io.Reader) (map[string]any, error) {
	var data map[string]any
	if err := json.NewDecoder(r).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func requireFields(m map[string]any, keys ...string) error {
	for _, k := range keys {
		if _, ok := m[k]; !ok {
			return fmt.Errorf("missing field: %s", k)
		}
	}
	return nil
}

func stringField(m map[string]any, key, def string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return def
}

func boolField(m map[string]any, key string, def bool) bool {
	if v, ok := m[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return def
}

func numberField(m map[string]any, key string, def int) int {
	if v, ok := m[key]; ok {
		return intVal(v)
	}
	return def
}

// --- Dialogue image generation helpers ---

// handleGenerateDialogueImage generates a PNG for a dialogue JSON file via the Python helper.
func (a *App) handleGenerateDialogueImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Filename string `json:"filename"`
			Preset   string `json:"preset_set"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("invalid request body: %v", err), http.StatusBadRequest)
			return
		}
		if req.Filename == "" {
			http.Error(w, "missing filename", http.StatusBadRequest)
			return
		}
		if req.Preset == "" {
			req.Preset = "presets_kurisu"
		}

		folder, err := a.getPresetFolder(req.Preset)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid preset_set: %v", err), http.StatusBadRequest)
			return
		}
		target := filepath.Join(folder, req.Filename)
		rawBytes, err := os.ReadFile(target)
		if err != nil {
			http.Error(w, fmt.Sprintf("cannot read file: %v", err), http.StatusNotFound)
			return
		}

		var raw map[string]any
		if err := json.Unmarshal(rawBytes, &raw); err != nil {
			http.Error(w, fmt.Sprintf("invalid JSON: %v", err), http.StatusBadRequest)
			return
		}

		dialogues := normalizeDialogues(raw)
		if len(dialogues) == 0 {
			http.Error(w, "no dialogues found to render", http.StatusBadRequest)
			return
		}

		// Extract ai_name: prefer top-level fields, then infer from dialogue speakers
		aiName := firstNonEmpty(raw["character_name"], raw["ai_name"], raw["character"])
		if aiName == "" {
			// Infer from the first non-user speaker in dialogues
			for _, d := range dialogues {
				sp := d["speaker"]
				if sp != "" && sp != "User" && sp != "user" && sp != "ai" {
					aiName = sp
					break
				}
			}
		}
		if aiName == "" {
			aiName = "AI"
		}

		// Determine language from preset
		isEnglish := !strings.Contains(req.Preset, "_CN")

		userName := firstNonEmpty(raw["user_name"])
		if userName == "" {
			if isEnglish {
				userName = "User"
			} else {
				userName = "用户"
			}
		}

		// Extract day number from filename (e.g. "Day3_dup_1_user_0.json" -> 3)
		dayNumber := 1
		dayRe := regexp.MustCompile(`Day(\d+)`)
		if m := dayRe.FindStringSubmatch(req.Filename); len(m) > 1 {
			if d, err := strconv.Atoi(m[1]); err == nil {
				dayNumber = d
			}
		}

		// Fetch the correct schedule from character_profiles.json
		// This ensures we get the right language version based on preset
		aiSchedule := a.getCharacterScheduleForDay(req.Preset, dayNumber)
		// If the profile schedule is empty (only has "day"), fall back to raw file data
		if len(aiSchedule) <= 1 {
			if rawSched, ok := raw["character_schedule"]; ok {
				aiSchedule, _ = rawSched.(map[string]any)
			}
		}

		converted := map[string]any{
			"ai_name":       aiName,
			"ai_schedule":   aiSchedule,
			"user_name":     userName,
			"user_schedule": raw["user_schedule"],
			"dialogues":     dialogues,
			"filename":      req.Filename,
			"preset_set":    req.Preset,
		}

		// Write converted JSON to temp file
		tmpJSON, err := os.CreateTemp("", "dialogue-*.json")
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to create temp json: %v", err), http.StatusInternalServerError)
			return
		}
		defer os.Remove(tmpJSON.Name())
		enc := json.NewEncoder(tmpJSON)
		enc.SetEscapeHTML(false)
		if err := enc.Encode(converted); err != nil {
			tmpJSON.Close()
			http.Error(w, fmt.Sprintf("failed to encode temp json: %v", err), http.StatusInternalServerError)
			return
		}
		tmpJSON.Close()

		// Prepare temp output image
		tmpPNG, err := os.CreateTemp("", "dialogue-*.png")
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to create temp image: %v", err), http.StatusInternalServerError)
			return
		}
		tmpPNG.Close()
		defer os.Remove(tmpPNG.Name())

		scriptPath := filepath.Join(a.cfg.RootDir, "convert_dialogue_to_image.py")
		if _, err := os.Stat(scriptPath); err != nil {
			http.Error(w, "image generation script missing", http.StatusInternalServerError)
			return
		}

		// Use "python" on Windows, "python3" on Unix
		pythonCmd := "python3"
		if runtime.GOOS == "windows" {
			pythonCmd = "python"
		}
		cmd := exec.CommandContext(r.Context(), pythonCmd, scriptPath, "--input", tmpJSON.Name(), "--output", tmpPNG.Name())
		cmd.Dir = a.cfg.RootDir
		out, err := cmd.CombinedOutput()
		if err != nil {
			http.Error(w, fmt.Sprintf("image generation failed: %v\n%s", err, string(out)), http.StatusInternalServerError)
			return
		}

		img, err := os.ReadFile(tmpPNG.Name())
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to read generated image: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Content-Length", strconv.Itoa(len(img)))
		w.Header().Set("Content-Disposition", "attachment; filename=\"dialogue.png\"")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(img)
	}
}

// videoOutputDir returns the directory where generated videos are stored, creating it if needed.
func (a *App) videoOutputDir() string {
	dir := filepath.Join(a.cfg.RootDir, "data", "generated_videos")
	_ = os.MkdirAll(dir, 0o755)
	return dir
}

// videoFilenameForDialogue returns a deterministic video filename based on the dialogue file.
func videoFilenameForDialogue(dialogueFilename, presetSet string) string {
	base := strings.TrimSuffix(dialogueFilename, filepath.Ext(dialogueFilename))
	return presetSet + "_" + base + ".mp4"
}

// handleCheckVideo checks whether a generated video exists for a given dialogue file.
func (a *App) handleCheckVideo(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")
	presetSet := r.URL.Query().Get("preset_set")
	if filename == "" || presetSet == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing filename or preset_set"})
		return
	}

	videoName := videoFilenameForDialogue(filename, presetSet)
	videoPath := filepath.Join(a.videoOutputDir(), videoName)

	if info, err := os.Stat(videoPath); err == nil && info.Size() > 0 {
		writeJSON(w, http.StatusOK, map[string]any{
			"exists":    true,
			"filename":  videoName,
			"size_mb":   fmt.Sprintf("%.1f", float64(info.Size())/(1024*1024)),
			"generated": info.ModTime().Format("2006-01-02 15:04:05"),
		})
	} else {
		writeJSON(w, http.StatusOK, map[string]any{"exists": false})
	}
}

// handleDownloadVideo serves a previously generated video file.
func (a *App) handleDownloadVideo(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")
	presetSet := r.URL.Query().Get("preset_set")
	if filename == "" || presetSet == "" {
		http.Error(w, "missing filename or preset_set", http.StatusBadRequest)
		return
	}

	videoName := videoFilenameForDialogue(filename, presetSet)
	videoPath := filepath.Join(a.videoOutputDir(), videoName)

	if _, err := os.Stat(videoPath); os.IsNotExist(err) {
		http.Error(w, "video not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "video/mp4")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", videoName))
	http.ServeFile(w, r, videoPath)
}

// handleGenerateVideo starts async video generation and returns immediately.
func (a *App) handleGenerateVideo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Filename string `json:"filename"`
			Preset   string `json:"preset_set"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("invalid request body: %v", err)})
			return
		}
		if req.Filename == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing filename"})
			return
		}
		if req.Preset == "" {
			req.Preset = "presets_kurisu"
		}

		// Job key = video filename
		videoName := videoFilenameForDialogue(req.Filename, req.Preset)

		// Check if already running
		a.videoJobsMu.RLock()
		if job, ok := a.videoJobs[videoName]; ok && job.Status == "running" {
			a.videoJobsMu.RUnlock()
			writeJSON(w, http.StatusOK, map[string]any{
				"status":  "running",
				"job_key": videoName,
				"message": "Video generation already in progress",
			})
			return
		}
		a.videoJobsMu.RUnlock()

		folder, err := a.getPresetFolder(req.Preset)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("invalid preset_set: %v", err)})
			return
		}

		target := filepath.Join(folder, req.Filename)
		rawBytes, err := os.ReadFile(target)
		if err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("cannot read file: %v", err)})
			return
		}

		var raw map[string]any
		if err := json.Unmarshal(rawBytes, &raw); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("invalid JSON: %v", err)})
			return
		}

		// Extract dialogue
		var dialogueItems []any
		if v, ok := raw["dialogue"]; ok {
			if arr, ok := v.([]any); ok {
				dialogueItems = arr
			}
		} else if v, ok := raw["dialogues"]; ok {
			if arr, ok := v.([]any); ok {
				dialogueItems = arr
			}
		}
		if len(dialogueItems) == 0 {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "no dialogue found in file"})
			return
		}

		aiName := firstNonEmpty(raw["character_name"], raw["ai_name"], raw["character"])
		if aiName == "" {
			aiName = "AI"
		}
		isEnglish := !strings.Contains(req.Preset, "_CN")
		userName := firstNonEmpty(raw["user_name"])
		if userName == "" {
			if isEnglish {
				userName = "User"
			} else {
				userName = "用户"
			}
		}

		// Extract inner_thought_annotations (top-level dict keyed by dialogue index)
		var innerThoughts map[string]any
		if v, ok := raw["inner_thought_annotations"]; ok {
			if obj, ok := v.(map[string]any); ok {
				innerThoughts = obj
			}
		}

		var messages []map[string]any
		lastShownTotalMins := -1 // track last displayed timestamp for 2-min gap filter
		for i, item := range dialogueItems {
			m, ok := item.(map[string]any)
			if !ok {
				continue
			}
			content := firstNonEmpty(m["content"], m["text"], m["message"], m["utterance"])
			if content == "" {
				continue
			}
			role := firstNonEmpty(m["role"], m["speaker"], m["identity"])
			sender := "a"
			if role == aiName || role == "Kurisu" || role == "ai" || role == "character" || role == "assistant" {
				sender = "b"
			}

			// Extract timestamp and apply 2-minute gap filter for periodic display
			showTimestamp := ""
			if ts, ok := m["timestamp"]; ok {
				if tsStr, ok := ts.(string); ok && tsStr != "" {
					tsParts := strings.SplitN(tsStr, ":", 2)
					if len(tsParts) == 2 {
						h, err1 := strconv.Atoi(strings.TrimSpace(tsParts[0]))
						mn, err2 := strconv.Atoi(strings.TrimSpace(tsParts[1]))
						if err1 == nil && err2 == nil && h >= 0 && h <= 23 && mn >= 0 && mn <= 59 {
							totalMins := h*60 + mn
							// Show timestamp only if first, >= 2 min gap, or day wraparound
							if lastShownTotalMins < 0 || totalMins-lastShownTotalMins >= 2 || totalMins < lastShownTotalMins {
								showTimestamp = tsStr
								lastShownTotalMins = totalMins
							}
						}
					}
				}
			}

			if strings.HasPrefix(content, "[[sticker:") && strings.HasSuffix(content, "]]") {
				stickerRel := content[len("[[sticker:"):len(content)-2]
				stickerPath := filepath.Join(a.cfg.RootDir, "stickers", stickerRel)
				if _, err := os.Stat(stickerPath); err == nil {
					stickerMsg := map[string]any{
						"sender":          sender,
						"type":            "sticker",
						"sticker_path":    stickerPath,
						"text":            "",
						"typing_duration": 0.5,
						"delay_after":     1.5,
					}
					if showTimestamp != "" {
						stickerMsg["timestamp"] = showTimestamp
					}
					messages = append(messages, stickerMsg)
				}
				continue
			}

			content = strings.TrimSpace(strings.TrimRight(content, "\n"))
			if content == "" {
				continue
			}

			charCount := len([]rune(content))
			typingDuration := math.Max(0.2, math.Min(2.0, float64(charCount)*0.05))
			delayAfter := math.Max(1.0, math.Min(3.0, float64(charCount)*0.06+0.8))

			msg := map[string]any{
				"sender":          sender,
				"text":            content,
				"typing_duration": math.Round(typingDuration*100) / 100,
				"delay_after":     math.Round(delayAfter*100) / 100,
			}
			if showTimestamp != "" {
				msg["timestamp"] = showTimestamp
			}
		// Inject inner thought BEFORE the spoken message so the viewer
		// sees the character's reasoning first, then the response.
		// Source 1: per-message "reasoning_chain" (character only)
		// Source 2: top-level "inner_thought_annotations" (higher priority)
		thoughtText := ""
		if sender == "b" {
			if rc, ok := m["reasoning_chain"]; ok {
				if rcStr, ok := rc.(string); ok {
					thoughtText = strings.TrimSpace(rcStr)
				}
			}
		}
		// Override with inner_thought_annotations if present
		idxKey := strconv.Itoa(i)
		if innerThoughts != nil {
			if ann, ok := innerThoughts[idxKey]; ok {
				if annMap, ok := ann.(map[string]any); ok {
					corrected := strings.TrimSpace(firstNonEmpty(annMap["correct_thought"]))
					actual := strings.TrimSpace(firstNonEmpty(annMap["actual_thought"]))
					if corrected != "" {
						thoughtText = corrected
					} else if actual != "" {
						thoughtText = actual
					}
				}
			}
		}
		if thoughtText != "" {
			messages = append(messages, map[string]any{
				"sender":          "b",
				"type":            "thought",
				"text":            thoughtText,
				"typing_duration": 0,
				"delay_after":     1.2,
			})
		}

		messages = append(messages, msg)
		}

		if len(messages) == 0 {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "no valid messages to render"})
			return
		}

		// Append character schedule card as the final message
		if cs, ok := raw["character_schedule"]; ok {
			if csMap, ok := cs.(map[string]any); ok {
				periods := []string{"morning", "noon", "afternoon", "evening", "night"}
				hasContent := false
				for _, p := range periods {
					if v, ok := csMap[p]; ok {
						if s, ok := v.(string); ok && strings.TrimSpace(s) != "" {
							hasContent = true
							break
						}
					}
				}
				if hasContent {
					messages = append(messages, map[string]any{
						"sender":          "b",
						"type":            "schedule",
						"text":            "",
						"schedule":        csMap,
						"character_name":  aiName,
						"typing_duration": 0,
						"delay_after":     3.0,
					})
				}
			}
		}

		headerText := "chat"
		dayRe := regexp.MustCompile(`Day(\d+)`)
		if m := dayRe.FindStringSubmatch(req.Filename); len(m) > 1 {
			if d, _ := strconv.Atoi(m[1]); d > 1 {
				headerText = userName
			}
		}

		charAvatar := filepath.Join(a.cfg.StaticDir, "kurisu", "kurisu_avatar.png")
		switch {
		case strings.Contains(req.Preset, "lin_lu"):
			charAvatar = filepath.Join(a.cfg.StaticDir, "lin_lu", "word.png")
		case strings.Contains(req.Preset, "kurisu"):
			charAvatar = filepath.Join(a.cfg.StaticDir, "kurisu", "kurisu_avatar.png")
		}

		script := map[string]any{
			"person_a": userName,
			"person_b": aiName,
			"messages": messages,
			"config": map[string]any{
				"fps":            30,
				"width":          1080,
				"height":         1920,
				"hold_at_end":    2.0,
				"person_a_image": filepath.Join(a.cfg.StaticDir, "user", "user_avatar.png"),
				"person_b_image": charAvatar,
				"header_text":    headerText,
			},
		}

		tmpJSON, err := os.CreateTemp("", "video-script-*.json")
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create temp file"})
			return
		}

		enc := json.NewEncoder(tmpJSON)
		enc.SetEscapeHTML(false)
		if err := enc.Encode(script); err != nil {
			tmpJSON.Close()
			os.Remove(tmpJSON.Name())
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to encode script"})
			return
		}
		tmpJSON.Close()

		videoPath := filepath.Join(a.videoOutputDir(), videoName)

		generatorScript := filepath.Join(a.cfg.RootDir, "video_chat_renderer", "video_renderer", "generator.py")
		if _, err := os.Stat(generatorScript); err != nil {
			os.Remove(tmpJSON.Name())
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "video generator script not found"})
			return
		}

		// Mark job as running
		a.videoJobsMu.Lock()
		a.videoJobs[videoName] = &VideoJob{
			Status:    "running",
			Filename:  videoName,
			StartedAt: time.Now(),
		}
		a.videoJobsMu.Unlock()

		msgCount := len(messages)
		tmpPath := tmpJSON.Name()
		reqFilename := req.Filename
		rootDir := a.cfg.RootDir

		// Launch generation in background goroutine
		go func() {
			defer os.Remove(tmpPath)

			pythonCmd := "python3"
			if runtime.GOOS == "windows" {
				pythonCmd = "python"
			}

			cmd := exec.Command(pythonCmd, generatorScript, tmpPath, videoPath, "--style", "chat", "--workers", "1")
			cmd.Dir = filepath.Join(rootDir, "video_chat_renderer", "video_renderer")

			log.Printf("[VIDEO] Generating video for %s -> %s (async)", reqFilename, videoPath)
			out, err := cmd.CombinedOutput()

			a.videoJobsMu.Lock()
			defer a.videoJobsMu.Unlock()

			if err != nil {
				log.Printf("[VIDEO] Generation failed: %v\n%s", err, string(out))
				a.videoJobs[videoName] = &VideoJob{
					Status:   "error",
					Filename: videoName,
					Error:    fmt.Sprintf("video generation failed: %v", err),
					Message:  string(out),
				}
				return
			}

			log.Printf("[VIDEO] Video generated successfully: %s", videoPath)
			sizeMB := "0"
			if info, statErr := os.Stat(videoPath); statErr == nil && info != nil {
				sizeMB = fmt.Sprintf("%.1f", float64(info.Size())/(1024*1024))
			}
			a.videoJobs[videoName] = &VideoJob{
				Status:   "done",
				Filename: videoName,
				SizeMB:   sizeMB,
				Message:  fmt.Sprintf("Video generated with %d messages", msgCount),
			}
		}()

		// Return immediately
		writeJSON(w, http.StatusOK, map[string]any{
			"status":  "started",
			"job_key": videoName,
			"message": fmt.Sprintf("Video generation started with %d messages", msgCount),
		})
	}
}

// handleVideoStatus reports the current status of a video generation job.
func (a *App) handleVideoStatus(w http.ResponseWriter, r *http.Request) {
	jobKey := r.URL.Query().Get("job_key")
	if jobKey == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing job_key"})
		return
	}

	a.videoJobsMu.RLock()
	job, ok := a.videoJobs[jobKey]
	a.videoJobsMu.RUnlock()

	if !ok {
		writeJSON(w, http.StatusOK, map[string]any{"status": "unknown"})
		return
	}

	writeJSON(w, http.StatusOK, job)
}

// handleVideoStatusSSE streams video job status via Server-Sent Events.
// The client opens ONE connection; the server pushes updates when status changes.
func (a *App) handleVideoStatusSSE(w http.ResponseWriter, r *http.Request) {
	jobKey := r.URL.Query().Get("job_key")
	if jobKey == "" {
		http.Error(w, "missing job_key", http.StatusBadRequest)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	ctx := r.Context()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	lastStatus := ""
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			a.videoJobsMu.RLock()
			job, exists := a.videoJobs[jobKey]
			var jobCopy VideoJob
			if exists {
				jobCopy = *job
			}
			a.videoJobsMu.RUnlock()

			if !exists {
				fmt.Fprintf(w, "data: {\"status\":\"unknown\"}\n\n")
				flusher.Flush()
				return
			}

			if jobCopy.Status != lastStatus {
				lastStatus = jobCopy.Status
				jsonBytes, _ := json.Marshal(jobCopy)
				fmt.Fprintf(w, "data: %s\n\n", jsonBytes)
				flusher.Flush()

				if jobCopy.Status == "done" || jobCopy.Status == "error" {
					return
				}
			}
		}
	}
}
func normalizeDialogues(raw map[string]any) []map[string]string {
	var source any
	if v, ok := raw["dialogues"]; ok {
		source = v
	} else if v, ok := raw["dialogue"]; ok {
		source = v
	}

	var items []any
	switch t := source.(type) {
	case []any:
		items = t
	case map[string]any:
		for _, v := range t {
			items = append(items, v)
		}
	default:
		return nil
	}

	cycle := []string{"ai", "user"}
	var result []map[string]string
	for idx, item := range items {
		switch v := item.(type) {
		case map[string]any:
			speaker := firstNonEmpty(v["speaker"], v["role"], v["identity"])
			if speaker == "" {
				speaker = cycle[idx%2]
			}
			text := firstNonEmpty(v["text"], v["content"], v["message"], v["utterance"])
			if text == "" {
				text = fmt.Sprintf("%v", v)
			}
			result = append(result, map[string]string{"speaker": speaker, "text": text})
		case string:
			result = append(result, map[string]string{"speaker": cycle[idx%2], "text": v})
		default:
			result = append(result, map[string]string{"speaker": cycle[idx%2], "text": fmt.Sprintf("%v", v)})
		}
	}
	return result
}

func scheduleToString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case map[string]interface{}:
		var parts []string
		for key, val := range t {
			if key == "day" {
				continue
			}
			parts = append(parts, fmt.Sprintf("%s: %v", key, val))
		}
		sort.Strings(parts)
		return strings.Join(parts, "\n")
	default:
		return ""
	}
}

func firstNonEmpty(values ...any) string {
	for _, v := range values {
		if s := toString(v); s != "" {
			return s
		}
	}
	return ""
}

func toString(v any) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

