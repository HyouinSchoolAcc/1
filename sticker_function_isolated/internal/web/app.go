package web

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
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

	"data_labler_ui_go/internal/database"
	"data_labler_ui_go/internal/email"
)

// App holds server state and helpers.
type App struct {
	cfg               Config
	presetFolders     map[string]string
	excellent         map[string]map[string]bool
	fileLocks         sync.Map // path -> *sync.Mutex
	authService       *AuthService
	loungeStore       *LoungeStore
	characterStore    *CharacterStore
	registrationStore *database.RegistrationStore
	emailService      *email.Service
	db                *database.DB
	llmService        *LLMService
	tutorialStore     *database.TutorialStore
	chatStore         *database.ChatStore
}

func NewApp(cfg Config) *App {
	presets := []string{
		"presets_kurisu",
		"presets_kurisu_CN",
		"presets_kafka",
		"presets_kafka_CN",
		"presets_lin_lu_CN",
		"presets_lin_lu",
		"presets_newcharacter_1",
	}
	presetFolders := make(map[string]string, len(presets))
	for _, p := range presets {
		presetFolders[p] = filepath.Join(cfg.PresetBaseDir, p)
		_ = os.MkdirAll(presetFolders[p], 0o755)
	}

	excellent := map[string]map[string]bool{
		"presets_kurisu_CN":     setOf("user_32", "user_33"),
		"presets_kafka_CN":      setOf("user_0", "user_6"),
		"presets_lin_lu_CN": setOf("user_1", "user_4"),
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

	return &App{
		cfg:               cfg,
		presetFolders:     presetFolders,
		excellent:         excellent,
		authService:       authService,
		loungeStore:       loungeStore,
		characterStore:    characterStore,
		registrationStore: registrationStore,
		emailService:      emailService,
		db:                db,
		tutorialStore:     tutorialStore,
		chatStore:         chatStore,
	}
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
	log.Println("✓ Temp character cleanup routine started (runs every 1 hour)")
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
			_ = copyFile(path, filepath.Join(backupDir, backupName))
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

var filenameRe = regexp.MustCompile(`((?:user|temp)_\d+)_Day(\d+)_dup_(\d+)_simplified\.json`)

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

// === character & schedule helpers ===

func getCurrentCharacterRole(preset string) string {
	l := strings.ToLower(preset)
	switch {
	case strings.Contains(l, "kafka"):
		return "Kafka"
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
	case strings.Contains(l, "kafka"):
		return "Kafka"
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

// default schedule fallback matching python
func fallbackSchedule(day int, preset string) map[string]any {
	role := getCurrentCharacterRole(preset)
	isCN := strings.Contains(preset, "_CN")
	switch role {
	case "Kafka":
		if isCN {
			return map[string]any{"day": day, "早晨": "查阅最新情报报告，制定战略行动计划。", "下午": "与其他星核猎手协调，评估当前任务参数。", "晚上": "反思一天的发展，享受片刻宁静。"}
		}
		return map[string]any{"day": day, "morning": "Review recent intelligence reports and plan strategic operations.", "afternoon": "Coordinate with fellow Stellaron Hunters and assess current mission parameters.", "evening": "Reflect on the day's developments while enjoying a quiet moment."}
	case "Kurisu":
		if isCN {
			return map[string]any{"day": day, "早晨": "启动加速器，校准设备参数，准备实验。", "下午": "进行时间旅行理论研究，分析实验数据。", "晚上": "通过阅读科学论文或观看轻松娱乐来放松身心。"}
		}
		return map[string]any{"day": day, "morning": "Start the accelerator, calibrate the equipment, and prepare for experiments.", "afternoon": "Conduct research on time travel theories and analyze experimental data.", "evening": "Unwind by reviewing scientific papers or watching some light entertainment."}
	case "lin_lu", "林路":
		if isCN {
			return map[string]any{"day": day, "早晨": "备课，准备当天的文学史课程资料。", "下午": "给学生上课，讲解唐诗宋词鉴赏。", "晚上": "写自己的历史小说，或者在学校大银杏树下看云放松。"}
		}
		return map[string]any{"day": day, "morning": "Prepare for classes and organize course materials for literature history.", "afternoon": "Teach students about Tang poetry and Song lyrics appreciation.", "evening": "Write my historical novel or relax under the big ginkgo tree at school."}
	default:
		return map[string]any{"day": day}
	}
}

func (a *App) getUniversalCharacterSchedule(day int, preset string) map[string]any {
	folder, err := a.getPresetFolder(preset)
	if err != nil {
		return fallbackSchedule(day, preset)
	}
	path := filepath.Join(folder, "character_schedules", "default_schedule.json")
	b, err := os.ReadFile(path)
	if err != nil {
		return fallbackSchedule(day, preset)
	}
	var data map[string]any
	if err := json.Unmarshal(b, &data); err != nil {
		return fallbackSchedule(day, preset)
	}
	arr, ok := data["default_schedule"].([]any)
	if !ok {
		return fallbackSchedule(day, preset)
	}
	for _, v := range arr {
		if m, ok := v.(map[string]any); ok {
			if intVal(m["day"]) == day {
				return m
			}
		}
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

func (a *App) updateUniversalCharacterSchedule(day int, schedule map[string]any, preset string) error {
	folder, err := a.getPresetFolder(preset)
	if err != nil {
		return err
	}
	path := filepath.Join(folder, "character_schedules", "default_schedule.json")
	_ = os.MkdirAll(filepath.Dir(path), 0o755)
	var data map[string]any
	if b, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(b, &data)
	}
	if data == nil {
		data = map[string]any{
			"sche_id":          0,
			"chinese_trans":    "通用角色日程",
			"default_schedule": []any{},
		}
	}
	schedule["day"] = day
	arr, _ := data["default_schedule"].([]any)
	updated := false
	for i, v := range arr {
		if m, ok := v.(map[string]any); ok && intVal(m["day"]) == day {
			arr[i] = schedule
			updated = true
			break
		}
	}
	if !updated {
		arr = append(arr, schedule)
	}
	// sort
	sort.Slice(arr, func(i, j int) bool {
		di := 0
		if m, ok := arr[i].(map[string]any); ok {
			di = intVal(m["day"])
		}
		dj := 0
		if m, ok := arr[j].(map[string]any); ok {
			dj = intVal(m["day"])
		}
		return di < dj
	})
	data["default_schedule"] = arr
	return a.atomicWriteJSON(path, data)
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
	case "presets_kafka_CN":
		return "new_user_info_kafka_cn.json"
	case "presets_kurisu_CN":
		return "new_user_info_kurisu_cn.json"
	case "presets_kafka":
		return "new_user_info_kafka.json"
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

		converted := map[string]any{
			"ai_name":       firstNonEmpty(raw["character_name"], raw["ai_name"], raw["character"]),
			"ai_schedule":   scheduleToString(raw["character_schedule"]),
			"user_name":     firstNonEmpty(raw["user_name"]),
			"user_schedule": scheduleToString(raw["user_schedule"]),
			"dialogues":     dialogues,
		}
		if converted["ai_name"] == "" {
			converted["ai_name"] = "AI"
		}
		if converted["user_name"] == "" {
			converted["user_name"] = "用户"
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

		cmd := exec.CommandContext(r.Context(), "python3", scriptPath, "--input", tmpJSON.Name(), "--output", tmpPNG.Name())
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
		w.Header().Set("Content-Disposition", "attachment; filename=\"dialogue.png\"")
		_, _ = w.Write(img)
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
