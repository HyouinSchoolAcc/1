package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,30}$`)
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

// handleRegisterPage renders the registration page
func (a *App) handleRegisterPage(tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// If already logged in, redirect to landing
		sessionData := a.getCurrentUser(r)
		if sessionData != nil {
			http.Redirect(w, r, "/landing", http.StatusSeeOther)
			return
		}

		if r.Method == "GET" {
			// Get language from middleware context
			language := GetLanguage(r)
			
			data := TemplateData{
				LoggedIn:    false,
				CurrentUser: "",
				Language:    language,
			}
			
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			if err := tpl.Execute(w, data); err != nil {
				log.Printf("Error rendering register template: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}

		// POST request handled by handleRegister
		a.handleRegister(w, r)
	}
}

// handleRegister processes registration requests
func (a *App) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "Invalid request format",
		})
		return
	}

	// Validate inputs
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if !usernameRegex.MatchString(req.Username) {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "Username must be 3-30 characters and contain only letters, numbers, and underscores",
		})
		return
	}

	if !emailRegex.MatchString(req.Email) {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "Invalid email address",
		})
		return
	}

	if len(req.Password) < 8 {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "Password must be at least 8 characters long",
		})
		return
	}

	// Check if username already exists
	usernameExists, err := a.registrationStore.CheckUsernameExists(req.Username)
	if err != nil {
		log.Printf("Error checking username existence: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   "Internal server error",
		})
		return
	}
	if usernameExists {
		writeJSON(w, http.StatusConflict, map[string]any{
			"success": false,
			"error":   "Username is already taken",
		})
		return
	}

	// Check if email already exists
	emailExists, err := a.registrationStore.CheckEmailExists(req.Email)
	if err != nil {
		log.Printf("Error checking email existence: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   "Internal server error",
		})
		return
	}
	if emailExists {
		writeJSON(w, http.StatusConflict, map[string]any{
			"success": false,
			"error":   "Email is already registered",
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   "Internal server error",
		})
		return
	}

	// Create user account immediately (skip email confirmation)
	// New registrations get "writer" role so they can start writing immediately
	userID := generateID()
	query := `INSERT INTO users (id, username, email, password, role, created_at) VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`
	_, err = a.db.Exec(query, userID, req.Username, req.Email, string(hashedPassword), "writer")
	if err != nil {
		log.Printf("Error creating user: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   "Failed to create account. Please try again.",
		})
		return
	}

	log.Printf("User registered successfully: %s (email: %s, ID: %s)", req.Username, req.Email, userID)
	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"message": "Registration successful! You can now log in.",
	})
}

// handleConfirmEmail processes email confirmation
func (a *App) handleConfirmEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Render the confirmation page
		http.ServeFile(w, r, a.cfg.TemplatesDir+"/confirm_email.html")
		return
	}

	// POST request - process confirmation
	token := r.URL.Query().Get("token")
	if token == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "Confirmation token is required",
		})
		return
	}

	// Confirm email and create user
	user, err := a.registrationStore.ConfirmEmail(token)
	if err != nil {
		log.Printf("Error confirming email: %v", err)
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	log.Printf("Email confirmed and user created: %s (ID: %s)", user.Username, user.ID)
	writeJSON(w, http.StatusOK, map[string]any{
		"success":  true,
		"message":  "Email confirmed successfully! You can now log in.",
		"username": user.Username,
	})
}

// getBaseURL gets the base URL for the application
func (a *App) getBaseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	// Check for X-Forwarded-Proto header (for proxies/load balancers)
	if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
		scheme = proto
	}

	host := r.Host
	if host == "" {
		host = "localhost:5002"
	}

	return fmt.Sprintf("%s://%s", scheme, host)
}
