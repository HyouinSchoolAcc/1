package web

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"data_labler_ui_go/internal/database"
	"github.com/go-chi/chi/v5"
)

// XPTier defines what a given XP level unlocks
type XPTier struct {
	Level        int    `json:"level"`
	Name         string `json:"name"`
	XPRequired   int    `json:"xp_required"`
	Title        string `json:"title"`
	BorderClass  string `json:"border_class"`
	ChatColor    string `json:"chat_color"`
	ChatColorHex string `json:"chat_color_hex"`
}

var xpTiers = []XPTier{
	{1, "Newcomer",    0,    "Newcomer",    "border-none",     "Default", "#6b7280"},
	{2, "Regular",     100,  "Regular",     "border-silver",   "Sky Blue", "#38bdf8"},
	{3, "Contributor", 500,  "Contributor", "border-gold",     "Emerald", "#34d399"},
	{4, "Veteran",     1500, "Veteran",     "border-gradient", "Purple", "#a78bfa"},
	{5, "Elite",       3000, "Elite",       "border-rainbow",  "Gold", "#fbbf24"},
	{6, "Legend",      6000, "Legend",      "border-legend",   "Crimson", "#f87171"},
}

// CharacterFlair is a character-based contributor badge
type CharacterFlair struct {
	CharacterID string `json:"character_id"`
	Name        string `json:"name"`
	Level       string `json:"level"`
	Threshold   int    `json:"threshold"`
}

// StoreItem represents an item in the XP store catalog
type StoreItem struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	Cost     int    `json:"cost"`
	Icon     string `json:"icon"`
	Category string `json:"category"`
	Preview  string `json:"preview"`
}

var storeItems = []StoreItem{
	{"silver_frame",   "Silver Frame",       "A sleek silver border for your avatar.",          100,  "🖼️", "border", "border-silver"},
	{"gold_glow",      "Gold Glow",          "Radiant gold avatar glow.",                       300,  "✨", "border", "border-gold"},
	{"gradient_frame", "Gradient Frame",     "Animated gradient border.",                       600,  "🌈", "border", "border-gradient"},
	{"sky_name",       "Sky Name Color",     "Display your name in sky blue in chat.",          150,  "💬", "chat",   "#38bdf8"},
	{"emerald_name",   "Emerald Name",       "Emerald green chat name color.",                  250,  "💬", "chat",   "#34d399"},
	{"purple_aura",    "Purple Aura",        "Purple chat name — radiates mystery.",            400,  "💜", "chat",   "#a78bfa"},
	{"star_flair",     "⭐ Star Flair",       "Add a star flair next to your username.",        200,  "⭐", "flair",  ""},
	{"ink_quill",      "🪶 Ink Quill Flair", "The quill of a seasoned writer.",                500,  "🪶", "flair",  ""},
}

type FullProfile struct {
	UserID           string     `json:"user_id"`
	Username         string     `json:"username"`
	RealName         string     `json:"real_name"`
	DisplayName      string     `json:"display_name"`
	Role             string     `json:"role"`
	Bio              string     `json:"bio"`
	AvatarColor      string     `json:"avatar_color"`
	IsPublic         bool       `json:"is_public"`
	JoinedAt         time.Time  `json:"joined_at"`

	XP               int        `json:"xp"`
	XPLevel          int        `json:"xp_level"`
	XPLevelName      string     `json:"xp_level_name"`
	XPToNext         int        `json:"xp_to_next"`
	XPProgressPct    int        `json:"xp_progress_pct"`
	XPTitle          string     `json:"xp_title"`
	XPBorderClass    string     `json:"xp_border_class"`
	XPChatColor      string     `json:"xp_chat_color"`
	XPChatColorHex   string     `json:"xp_chat_color_hex"`
	DaysLoggedIn     int        `json:"days_logged_in"`

	InkTokens        int        `json:"ink_tokens"`
	BonusInk         int        `json:"bonus_ink"`
	InkLevel         int        `json:"ink_level"`
	InkLevelName     string     `json:"ink_level_name"`
	InkToNext        int        `json:"ink_to_next"`
	InkProgressPct   int        `json:"ink_progress_pct"`

	ApprovedDialogues int        `json:"approved_dialogues"`
	TotalDialogues   int        `json:"total_dialogues"`
	WordsPassed      int        `json:"words_passed"`
	ApprovalRate     float64    `json:"approval_rate"`
	LeaderboardRank  int        `json:"leaderboard_rank"`
	ExcellentCount   int        `json:"excellent_count"`

	Certifications    []database.CharacterCertification `json:"certifications"`
	TutorialCompleted bool                              `json:"tutorial_completed"`
	CharacterFlairs   []CharacterFlair                  `json:"character_flairs"`
	UnlockedTiers     []XPTier                          `json:"unlocked_tiers"`

	ActiveBanner     string   `json:"active_banner"`
	ActiveBackground string   `json:"active_background"`
	OwnedBanners     []string `json:"owned_banners"`
	OwnedBackgrounds []string `json:"owned_backgrounds"`

	InkHistory []database.PointTransaction `json:"ink_history,omitempty"`
	IsOwner    bool                         `json:"is_owner"`
}

// RegisterProfileAPI wires profile and store endpoints
func (a *App) RegisterProfileAPI(r chi.Router) {
	r.Get("/api/profile/{username}", a.handleGetProfile)
	r.Get("/api/profile/me/points", a.handleGetMyPoints)
	r.Put("/api/profile/me", a.handleUpdateMyProfile)
	r.Post("/api/profile/{username}/award-points", a.handleAwardPoints)
	r.Get("/api/store/items", a.handleGetStoreItems)
	r.Post("/api/store/purchase", a.handleStorePurchase)
	r.Post("/api/cosmetics/grant", a.handleGrantCosmetic)
}

func (a *App) handleProfile(tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")
		sessionData := a.getCurrentUser(r)
		profile, err := a.buildFullProfile(username, sessionData)
		if err != nil {
			http.Error(w, "Profile not found", http.StatusNotFound)
			return
		}
		profileJSON, _ := json.Marshal(profile)
		storeJSON, _ := json.Marshal(storeItems)
		cu, ur := "", "new_user"
		if sessionData != nil {
			cu = sessionData.Username
			ur = string(sessionData.Role)
		}
		data := map[string]any{
			"LoggedIn":    sessionData != nil,
			"CurrentUser": cu,
			"UserRole":    ur,
			"Profile":     profile,
			"ProfileJSON": string(profileJSON),
			"StoreItems":  storeItems,
			"StoreJSON":   string(storeJSON),
			"Language":    GetLanguage(r),
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tpl.Execute(w, data); err != nil {
			log.Printf("Error rendering profile template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func (a *App) handleGetProfile(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	sessionData := a.getCurrentUser(r)
	profile, err := a.buildFullProfile(username, sessionData)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "Profile not found"})
		return
	}
	writeJSON(w, http.StatusOK, profile)
}

func (a *App) handleGetMyPoints(w http.ResponseWriter, r *http.Request) {
	sessionData := a.getCurrentUser(r)
	if sessionData == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "Not logged in"})
		return
	}
	history, err := a.profileStore.GetInkHistory(sessionData.UserID, 50)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "Failed"})
		return
	}
	if history == nil {
		history = []database.PointTransaction{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"history": history})
}

func (a *App) handleUpdateMyProfile(w http.ResponseWriter, r *http.Request) {
	sessionData := a.getCurrentUser(r)
	if sessionData == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "Not logged in"})
		return
	}
	var body struct {
		Bio              string
		DisplayName      string
		AvatarColor      string
		IsPublic         *bool
		ActiveBanner     string
		ActiveBackground string
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "Invalid JSON"})
		return
	}
	profile, err := a.profileStore.GetProfile(sessionData.UserID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "Failed"})
		return
	}
	if len(body.Bio) > 500 {
		body.Bio = body.Bio[:500]
	}
	profile.Bio = body.Bio
	if body.DisplayName != "" {
		if len(body.DisplayName) > 50 {
			body.DisplayName = body.DisplayName[:50]
		}
		profile.DisplayName = body.DisplayName
	}
	if body.AvatarColor != "" && strings.HasPrefix(body.AvatarColor, "#") {
		profile.AvatarColor = body.AvatarColor
	}
	if body.IsPublic != nil {
		profile.IsPublic = *body.IsPublic
	}
	if body.ActiveBanner == "none" {
		profile.ActiveBanner = ""
	} else if body.ActiveBanner != "" {
		if owned, _ := a.profileStore.HasCosmetic(sessionData.UserID, "banner", body.ActiveBanner); owned {
			profile.ActiveBanner = body.ActiveBanner
		}
	}
	if body.ActiveBackground == "none" {
		profile.ActiveBackground = ""
	} else if body.ActiveBackground != "" {
		if owned, _ := a.profileStore.HasCosmetic(sessionData.UserID, "background", body.ActiveBackground); owned {
			profile.ActiveBackground = body.ActiveBackground
		}
	}
	if err := a.profileStore.UpdateProfile(profile); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "Save failed"})
		return
	}
	if body.Bio != "" {
		_ = a.profileStore.AddXP(sessionData.UserID, 5)
	}
	writeJSON(w, http.StatusOK, map[string]any{"success": true})
}

func (a *App) handleAwardPoints(w http.ResponseWriter, r *http.Request) {
	sessionData := a.getCurrentUser(r)
	if sessionData == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "Not logged in"})
		return
	}
	if string(sessionData.Role) != "editor" {
		writeJSON(w, http.StatusForbidden, map[string]any{"error": "Editor only"})
		return
	}
	username := chi.URLParam(r, "username")
	targetUser, found := a.authService.userStore.GetUser(username)
	if !found {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "User not found"})
		return
	}
	var body struct {
		Amount int
		Reason string
		Type   string
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "Invalid JSON"})
		return
	}
	if body.Amount == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "Amount required"})
		return
	}
	if body.Reason == "" {
		body.Reason = "Manual award by editor"
	}
	if body.Type == "xp" {
		if err := a.profileStore.AddXP(targetUser.ID, body.Amount); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "Failed"})
			return
		}
	} else {
		if err := a.profileStore.AwardInk(targetUser.ID, body.Amount, body.Reason, sessionData.Username); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "Failed"})
			return
		}
	}
	log.Printf("[AWARD] %s awarded %d %s to %s", sessionData.Username, body.Amount, body.Type, username)
	writeJSON(w, http.StatusOK, map[string]any{"success": true, "awarded": body.Amount})
}

func (a *App) handleGetStoreItems(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"items": storeItems})
}

func (a *App) handleStorePurchase(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"success": false,
		"message": "The store is coming soon!",
	})
}

func (a *App) handleGrantCosmetic(w http.ResponseWriter, r *http.Request) {
	sessionData := a.getCurrentUser(r)
	if sessionData == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "Not logged in"})
		return
	}
	if string(sessionData.Role) != "editor" {
		writeJSON(w, http.StatusForbidden, map[string]any{"error": "Editor only"})
		return
	}
	var body struct {
		Username     string
		CosmeticType string
		CosmeticID   string
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "Invalid JSON"})
		return
	}
	target, found := a.authService.userStore.GetUser(body.Username)
	if !found {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "User not found"})
		return
	}
	if err := a.profileStore.GrantCosmetic(target.ID, body.CosmeticType, body.CosmeticID, sessionData.Username); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "Failed"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"success": true})
}

func (a *App) buildFullProfile(username string, viewer *SessionData) (*FullProfile, error) {
	user, found := a.authService.userStore.GetUser(username)
	if !found {
		return nil, os.ErrNotExist
	}
	dbProfile, err := a.profileStore.GetProfile(user.ID)
	if err != nil {
		return nil, err
	}
	profile := &FullProfile{
		UserID:       user.ID,
		Username:     user.Username,
		Role:         string(user.Role),
		Bio:          dbProfile.Bio,
		AvatarColor:  dbProfile.AvatarColor,
		DisplayName:  dbProfile.DisplayName,
		IsPublic:     dbProfile.IsPublic,
		JoinedAt:     user.CreatedAt,
		XP:           dbProfile.XP,
		DaysLoggedIn: dbProfile.DaysLoggedIn,
	}
	if viewer != nil && viewer.Username == username {
		profile.IsOwner = true
	}
	a.fillWritingStats(profile)
	certs, _ := a.tutorialStore.GetUserCertifications(user.ID)
	if certs == nil {
		certs = []database.CharacterCertification{}
	}
	profile.Certifications = certs
	if progress, _ := a.tutorialStore.GetTutorialProgress(user.ID); progress != nil {
		profile.TutorialCompleted = progress.IsCompleted
	}
	profile.XPLevel, profile.XPLevelName, profile.XPToNext, profile.XPProgressPct = calcXPLevel(profile.XP)
	if profile.XPLevel <= len(xpTiers) {
		tier := xpTiers[profile.XPLevel-1]
		profile.XPTitle = tier.Title
		profile.XPBorderClass = tier.BorderClass
		profile.XPChatColor = tier.ChatColor
		profile.XPChatColorHex = tier.ChatColorHex
	}
	for _, t := range xpTiers {
		if t.XPRequired <= profile.XP {
			profile.UnlockedTiers = append(profile.UnlockedTiers, t)
		}
	}
	bonusInk, _ := a.profileStore.GetBonusInk(user.ID)
	profile.BonusInk = bonusInk
	tutorialBonus := 0
	if profile.TutorialCompleted {
		tutorialBonus = 100
	}
	profile.InkTokens = (profile.WordsPassed/10) + tutorialBonus + (len(certs) * 50) + (profile.ExcellentCount * 200) + bonusInk
	profile.InkLevel, profile.InkLevelName, profile.InkToNext, profile.InkProgressPct = calcInkLevel(profile.InkTokens)
	profile.CharacterFlairs = a.buildCharacterFlairs(profile)
	summary := a.getPaymentSummary()
	for i, w := range summary.Leaderboard {
		if w.Username == username {
			profile.LeaderboardRank = i + 1
			break
		}
	}
	profile.RealName = a.getRealNameFromFiles(user.ID, user.Username)
	if profile.IsOwner {
		if history, _ := a.profileStore.GetInkHistory(user.ID, 20); history != nil {
			profile.InkHistory = history
		}
	}
	profile.ActiveBanner = dbProfile.ActiveBanner
	profile.ActiveBackground = dbProfile.ActiveBackground
	if banners, _ := a.profileStore.GetUserCosmetics(user.ID, "banner"); banners != nil {
		profile.OwnedBanners = banners
	} else {
		profile.OwnedBanners = []string{}
	}
	if bgs, _ := a.profileStore.GetUserCosmetics(user.ID, "background"); bgs != nil {
		profile.OwnedBackgrounds = bgs
	} else {
		profile.OwnedBackgrounds = []string{}
	}
	return profile, nil
}

func (a *App) fillWritingStats(profile *FullProfile) {
	for preset, folder := range a.presetFolders {
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
			var data map[string]any
			if err := a.loadJSON(filepath.Join(folder, e.Name()), &data); err != nil {
				continue
			}
			pid := stringField(data, "personnel_id", "")
			uname := stringField(data, "user_name", "")
			if !((pid != "" && pid == profile.UserID) || (uname != "" && uname == profile.Username)) {
				continue
			}
			dialogue, _ := data["dialogue"].([]any)
			charCount := countDialogueHanChars(dialogue)
			profile.TotalDialogues++
			if stringField(data, "category", "pending") == "passed" {
				profile.ApprovedDialogues++
				profile.WordsPassed += charCount
			}
			pf, _ := parseWriterFilename(e.Name())
			if a.isPublicUser(preset, pf.UserID) {
				profile.ExcellentCount++
			}
		}
	}
	if profile.TotalDialogues > 0 {
		profile.ApprovalRate = float64(profile.ApprovedDialogues) / float64(profile.TotalDialogues) * 100
	}
}

func (a *App) getRealNameFromFiles(userID, username string) string {
	for _, folder := range a.presetFolders {
		entries, _ := os.ReadDir(folder)
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") || e.Name() == "_character_defaults.json" {
				continue
			}
			var data map[string]any
			if err := a.loadJSON(filepath.Join(folder, e.Name()), &data); err != nil {
				continue
			}
			pid := stringField(data, "personnel_id", "")
			uname := stringField(data, "user_name", "")
			if (pid != "" && pid == userID) || (uname != "" && uname == username) {
				if rn := stringField(data, "real_name", ""); rn != "" {
					return rn
				}
			}
		}
	}
	return ""
}

func (a *App) buildCharacterFlairs(profile *FullProfile) []CharacterFlair {
	charCounts := make(map[string]int)
	for _, folder := range a.presetFolders {
		entries, _ := os.ReadDir(folder)
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") || e.Name() == "_character_defaults.json" {
				continue
			}
			var data map[string]any
			if err := a.loadJSON(filepath.Join(folder, e.Name()), &data); err != nil {
				continue
			}
			pid := stringField(data, "personnel_id", "")
			uname := stringField(data, "user_name", "")
			if !((pid != "" && pid == profile.UserID) || (uname != "" && uname == profile.Username)) {
				continue
			}
			if stringField(data, "category", "") != "passed" {
				continue
			}
			if charName := stringField(data, "character_name", ""); charName != "" {
				charCounts[charName]++
			}
		}
	}
	var flairs []CharacterFlair
	for charID, count := range charCounts {
		level, threshold := "", 0
		switch {
		case count >= 20:
			level, threshold = "Master", 20
		case count >= 10:
			level, threshold = "Expert", 10
		case count >= 3:
			level, threshold = "Specialist", 3
		}
		if level != "" {
			flairs = append(flairs, CharacterFlair{charID, charID, level, threshold})
		}
	}
	return flairs
}

func calcXPLevel(xp int) (int, string, int, int) {
	for i := len(xpTiers) - 1; i >= 0; i-- {
		t := xpTiers[i]
		if xp >= t.XPRequired {
			level := i + 1
			if level >= len(xpTiers) {
				return level, t.Name, 0, 100
			}
			next := xpTiers[level]
			span := next.XPRequired - t.XPRequired
			pct := 0
			if span > 0 {
				pct = (xp - t.XPRequired) * 100 / span
			}
			return level, t.Name, next.XPRequired - xp, pct
		}
	}
	return 1, xpTiers[0].Name, xpTiers[1].XPRequired - xp, xp
}

func calcInkLevel(ink int) (int, string, int, int) {
	thresholds := []struct {
		min  int
		name string
	}{
		{0, "Ink Drop"}, {100, "Apprentice Scribe"}, {500, "Skilled Writer"},
		{2000, "Senior Author"}, {5000, "Master Storyteller"}, {10000, "Legendary Chronicler"},
	}
	for i := len(thresholds) - 1; i >= 0; i-- {
		t := thresholds[i]
		if ink >= t.min {
			level := i + 1
			if level >= len(thresholds) {
				return level, t.name, 0, 100
			}
			next := thresholds[level]
			span := next.min - t.min
			pct := 0
			if span > 0 {
				pct = (ink - t.min) * 100 / span
			}
			return level, t.name, next.min - ink, pct
		}
	}
	return 1, thresholds[0].name, thresholds[1].min - ink, ink
}
