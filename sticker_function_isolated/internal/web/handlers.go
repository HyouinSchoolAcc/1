package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// TemplateData holds data for template rendering
type TemplateData struct {
	LoggedIn    bool
	CurrentUser string
	UserRole    string
	Characters  []Character
	Language    string // "zh" or "en"
}

// Character represents a character for the writing page
type Character struct {
	ID          string
	Route       string // URL route for this character (may differ from ID)
	Name        string
	Image       string
	Tagline     string
	Description string
	Tags        []string
	Source      string
}

// loadCharacters loads available characters from preset folders
func (a *App) loadCharacters(lang string) []Character {
	isEn := lang == "en"
	// Use /e suffix for English, no suffix for Chinese
	langSuffix := ""
	if isEn {
		langSuffix = "/e"
	}
	
	characters := []Character{
		{
			ID:    "kurisu",
			Route: "/kurisu" + langSuffix,
			Name: func() string {
				if isEn {
					return "Makise Kurisu"
				}
				return "牧濑红莉栖"
			}(),
			Image: "kurisu/kurisu_avatar.png",
			Tagline: func() string {
				if isEn {
					return "Researcher · Data from lab survivors"
				}
				return "研究员 · 数据出自住过实验室的人"
			}(),
			Description: func() string {
				if isEn {
					return "Lab-to-dorm grind, staring at equipment, drawing MATLAB plots, wiring circuits — physics/chem/bio PhDs who've wasted a year on dead ends."
				}
				return "两点一线，干瞪器械，画MATLAB，接线，浪费过一年在无用功的物理/化学/生物PhD们"
			}(),
			Tags: func() []string {
				if isEn {
					return []string{"Physics", "Bio Labs", "PhD Life"}
				}
				return []string{"物理", "生物实验室", "科研"}
			}(),
			Source: func() string {
				if isEn {
					return "Steins;Gate"
				}
				return "命运石之门"
			}(),
		},
		{
			ID:    "kafka",
			Route: "/kafka" + langSuffix,
			Name: func() string {
				if isEn {
					return "Kafka"
				}
				return "卡夫卡"
			}(),
			Image: "kafka/kafka_avatar.png",
			Tagline: func() string {
				if isEn {
					return "Stellaron Hunter · Data from thrill seekers"
				}
				return "星际危险分子 · 数据出自情感老手"
			}(),
			Description: func() string {
				if isEn {
					return "Veterans of multiple relationships, emotionally numb yet driven to seek thrills and excitement."
				}
				return "谈过多段恋爱，情感麻木后自身能动力爆表追寻刺激的人"
			}(),
			Tags: func() []string {
				if isEn {
					return []string{"Romance", "Star Rail", "Thrill"}
				}
				return []string{"恋爱", "星铁", "刺激"}
			}(),
			Source: func() string {
				if isEn {
					return "Honkai: Star Rail"
				}
				return "崩坏：星穹铁道"
			}(),
		},
		{
			ID:    "lin_lu", // Preset folder ID
			Route: "/linlu" + langSuffix,
			Name: func() string {
				if isEn {
					return "Lin Lu"
				}
				return "林路"
			}(),
			Image: "lin_lu/word.png",
			Tagline: func() string {
				if isEn {
					return "Literature Professor · Data from literary minds"
				}
				return "文学教授 · 数据出自文学家们"
			}(),
			Description: func() string {
				if isEn {
					return "Coming soon — looking for writers with humanities background."
				}
				return "招募中 — 寻找人文背景的写手"
			}(),
			Tags: func() []string {
				if isEn {
					return []string{"Literature", "Humanities", "Poetry"}
				}
				return []string{"文学", "人文", "诗词"}
			}(),
			Source: func() string {
				if isEn {
					return "Original"
				}
				return "原创"
			}(),
		},
		{
			ID:    "newcharacter_1",
			Route: "/newcharacter_1" + langSuffix,
			Name: func() string {
				if isEn {
					return "New Character 1"
				}
				return "新角色 1"
			}(),
			Image: "user/simple-user-default-icon-free-png.png",
			Tagline: func() string {
				if isEn {
					return "Coming Soon"
				}
				return "即将推出的角色"
			}(),
			Description: func() string {
				if isEn {
					return "Participate in deciding the next character to be created."
				}
				return "参与决定下一个要创建的角色"
			}(),
			Tags: func() []string {
				if isEn {
					return []string{"Story", "Interaction", "New"}
				}
				return []string{"故事", "互动", "新"}
			}(),
			Source: func() string {
				if isEn {
					return "Original"
				}
				return "原创"
			}(),
		},
	}

	// Filter characters based on available presets
	available := []Character{}
	for _, char := range characters {
		// Check if preset folder exists - try both with and without _CN suffix
		presetNames := []string{"presets_" + char.ID, "presets_" + char.ID + "_CN"}
		found := false
		for _, presetName := range presetNames {
			if folder, err := a.getPresetFolder(presetName); err == nil {
				if _, err := os.Stat(folder); err == nil {
					found = true
					break
				}
			}
		}
		if found {
			available = append(available, char)
		}
	}

	return available
}

// getTemplateData creates template data with current user session info
func (a *App) getTemplateData(r *http.Request) TemplateData {
	sessionData := a.getCurrentUser(r)

	// Get language from context (set by LanguageMiddleware)
	language := GetLanguage(r)

	data := TemplateData{
		LoggedIn:    false,
		CurrentUser: "",
		UserRole:    string(RoleNewUser),
		Characters:  []Character{},
		Language:    language,
	}

	if sessionData != nil {
		data.LoggedIn = true
		data.CurrentUser = sessionData.Username
		data.UserRole = string(sessionData.Role)
	}

	return data
}

// handleLanding renders the landing page
func (a *App) handleLanding(tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := a.getTemplateData(r)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tpl.Execute(w, data); err != nil {
			log.Printf("Error rendering landing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

// handleLandingMultiLang renders the appropriate landing template based on language
func (a *App) handleLandingMultiLang(zhTpl, enTpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := a.getTemplateData(r)
		
		// Select template based on language
		tpl := zhTpl
		if data.Language == "en" {
			tpl = enTpl
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tpl.Execute(w, data); err != nil {
			log.Printf("Error rendering landing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

// handleWriting renders the navigation page with character selection
// (formerly called "writing page", now "navigation" - shows available characters)
func (a *App) handleWriting(tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := a.getTemplateData(r)
		data.Characters = a.loadCharacters(data.Language)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tpl.Execute(w, data); err != nil {
			log.Printf("Error rendering writing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

// handleFAQ renders the FAQ page
func (a *App) handleFAQ(tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := a.getTemplateData(r)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tpl.Execute(w, data); err != nil {
			log.Printf("Error rendering FAQ template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

// handleCharacter renders the Character Descriptions page (formerly "Character Background")
// Now accessible at /descriptions (with /character as backward compatibility)
func (a *App) handleCharacter(tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := a.getTemplateData(r)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tpl.Execute(w, data); err != nil {
			log.Printf("Error rendering character template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

// handlePayment renders the Payment page
func (a *App) handlePayment(tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := a.getTemplateData(r)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tpl.Execute(w, data); err != nil {
			log.Printf("Error rendering payment template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

// handleGuide renders the Writer's Guide page
func (a *App) handleGuide(tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := a.getTemplateData(r)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tpl.Execute(w, data); err != nil {
			log.Printf("Error rendering guide template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

// handleTemplate renders a simple template with common data
func (a *App) handleTemplate(tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := a.getTemplateData(r)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tpl.Execute(w, data); err != nil {
			log.Printf("Error rendering template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

// handleLounge renders the Writers' Lounge page
func (a *App) handleLounge(tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := a.getTemplateData(r)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tpl.Execute(w, data); err != nil {
			log.Printf("Error rendering lounge template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

// handleTutorial renders the interactive Cao Cao tutorial page
func (a *App) handleTutorial(tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := a.getTemplateData(r)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tpl.Execute(w, data); err != nil {
			log.Printf("Error rendering tutorial template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

// handleStoryboard renders the storyboard planning/viewing page
func (a *App) handleStoryboard(tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := a.getTemplateData(r)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tpl.Execute(w, data); err != nil {
			log.Printf("Error rendering storyboard template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

// loadGoTemplate loads and parses a Go template file, expanding Jinja-style includes first
func loadGoTemplate(templatesDir, name string) (*template.Template, error) {
	path := filepath.Join(templatesDir, name)
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Expand Jinja-style {% include %} directives before parsing as Go template
	expandedContent, err := expandJinjaIncludes(string(content), filepath.Dir(path))
	if err != nil {
		return nil, fmt.Errorf("expanding includes for %s: %w", name, err)
	}

	// Add custom template functions
	funcMap := template.FuncMap{
		"hasPrefix": func(s, prefix string) bool {
			return len(s) >= len(prefix) && s[:len(prefix)] == prefix
		},
		// langURL generates a language-aware URL using /e suffix pattern:
		// langURL "en" "/guide" -> "/guide/e"
		// langURL "zh" "/guide" -> "/guide"
		"langURL": func(lang, path string) string {
			if lang == "" {
				lang = "zh"
			}
			// Ensure path starts with /
			if !strings.HasPrefix(path, "/") {
				path = "/" + path
			}
			// For English, add /e suffix
			if lang == "en" {
				return path + "/e"
			}
			// For Chinese (default), no suffix
			return path
		},
	}

	tmpl, err := template.New(name).Funcs(funcMap).Parse(expandedContent)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

// expandJinjaIncludes expands {% include "path" %} directives recursively
func expandJinjaIncludes(content string, currentDir string) (string, error) {
	includeRE := regexp.MustCompile(`\{\%\s*include\s+"([^"]+)"\s*\%\}`)
	staticRE := regexp.MustCompile(`\{\{\s*url_for\(['"]static['"],\s*filename=['"]([^'"]+)['"]\)\s*\}\}`)

	// Keep expanding until no more includes found
	for {
		matches := includeRE.FindAllStringSubmatchIndex(content, -1)
		if matches == nil {
			break
		}

		var result strings.Builder
		last := 0
		for _, m := range matches {
			start, end := m[0], m[1]
			pathStart, pathEnd := m[2], m[3]
			result.WriteString(content[last:start])

			includeName := content[pathStart:pathEnd]
			includePath := filepath.Join(currentDir, includeName)

			// Read the include file
			includeBytes, err := os.ReadFile(includePath)
			if err != nil {
				return "", fmt.Errorf("include %s: %w", includePath, err)
			}

			// Recursively expand includes in the included file
			expanded, err := expandJinjaIncludes(string(includeBytes), filepath.Dir(includePath))
			if err != nil {
				return "", err
			}
			result.WriteString(expanded)
			last = end
		}
		result.WriteString(content[last:])
		content = result.String()
	}

	// Replace Flask static helpers with direct /static paths
	content = staticRE.ReplaceAllString(content, `/static/$1`)

	return content, nil
}
