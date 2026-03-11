package web

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
)

const (
	SessionName = "writer-portal-session"
	SessionKey  = "user-session-key-change-in-production" // In production, use environment variable
)

// SessionManager handles user sessions
type SessionManager struct {
	store sessions.Store
}

// NewSessionManager creates a new session manager
func NewSessionManager() *SessionManager {
	store := sessions.NewCookieStore([]byte(SessionKey))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400, // 24 hours
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	}
	return &SessionManager{store: store}
}

// SetSession creates a new session for the user
func (sm *SessionManager) SetSession(w http.ResponseWriter, r *http.Request, user *User) error {
	session, err := sm.store.Get(r, SessionName)
	if err != nil {
		return err
	}

	session.Values["user_id"] = user.ID
	session.Values["username"] = user.Username
	session.Values["role"] = string(user.Role)

	return session.Save(r, w)
}

// GetSession retrieves the current session data
func (sm *SessionManager) GetSession(r *http.Request) (*SessionData, error) {
	session, err := sm.store.Get(r, SessionName)
	// Note: Gorilla sessions.Get() can return an error for non-fatal issues
	// (e.g., cookie decode problems, checksum mismatches from old cookies)
	// but still returns a usable session object. We should check for session
	// data presence rather than treating all errors as fatal.
	if err != nil {
		// Log the error for debugging but don't fail immediately
		log.Printf("Session store warning (may be non-fatal): %v", err)
		// If we got no session object at all, that's a real error
		if session == nil {
			return nil, err
		}
		// Otherwise, try to read values anyway - the session might still be valid
	}

	userID, ok := session.Values["user_id"].(string)
	if !ok || userID == "" {
		return nil, nil
	}

	username, ok := session.Values["username"].(string)
	if !ok || username == "" {
		return nil, nil
	}

	role, ok := session.Values["role"].(string)
	if !ok || role == "" {
		return nil, nil
	}

	return &SessionData{
		UserID:   userID,
		Username: username,
		Role:     UserRole(role),
	}, nil
}

// ClearSession removes the user's session
func (sm *SessionManager) ClearSession(w http.ResponseWriter, r *http.Request) error {
	session, err := sm.store.Get(r, SessionName)
	if err != nil {
		return err
	}

	session.Options.MaxAge = -1
	return session.Save(r, w)
}

// AuthService handles authentication operations
type AuthService struct {
	userStore      *UserStore
	sessionManager *SessionManager
}

// NewAuthService creates a new authentication service with SQL-based user store
func NewAuthService(userStore *UserStore) *AuthService {
	return &AuthService{
		userStore:      userStore,
		sessionManager: NewSessionManager(),
	}
}

// handleLogin renders the login page or processes login
func (a *App) handleLogin(tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// Get redirect parameter to pass to template
			redirect := r.URL.Query().Get("redirect")

			// Get language from middleware context
			language := GetLanguage(r)

			// Always use English for login to ensure consistent English experience
			language = "en"

			// Render login page
			data := TemplateData{
				LoggedIn:    false,
				CurrentUser: "",
				Characters:  []Character{},
				Language:    language,
			}

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			if err := tpl.Execute(w, data); err != nil {
				log.Printf("Error rendering login template: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}

			// Store redirect in session for POST
			if redirect != "" {
				session, _ := a.authService.sessionManager.store.Get(r, "login-redirect")
				session.Values["redirect"] = redirect
				session.Save(r, w)
			}

		case http.MethodPost:
			// Process login
			a.processLogin(w, r)
		}
	}
}

// processLogin handles login form submission
func (a *App) processLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "无效的表单数据", http.StatusBadRequest)
		return
	}

	username := strings.TrimSpace(r.Form.Get("username"))
	password := r.Form.Get("password")

	if username == "" || password == "" {
		http.Error(w, "需要用户名和密码", http.StatusBadRequest)
		return
	}

	user, err := a.authService.userStore.ValidateUser(username, password)
	if err != nil {
		http.Error(w, "凭据无效", http.StatusUnauthorized)
		return
	}

	if err := a.authService.sessionManager.SetSession(w, r, user); err != nil {
		log.Printf("Error setting session: %v", err)
		http.Error(w, "内部服务器错误", http.StatusInternalServerError)
		return
	}

	// Determine if this is an English endpoint
	isEnglish := strings.HasSuffix(r.URL.Path, "/e") || r.URL.Path == "/e"
	langSuffix := ""
	if isEnglish {
		langSuffix = "/e"
	}

	// Check for redirect URL - first from form, then from session
	redirectPath := "/landing" + langSuffix // default with language

	// Try to get redirect from form data first
	formRedirect := r.Form.Get("redirect")
	if formRedirect != "" && strings.HasPrefix(formRedirect, "/") && !strings.HasPrefix(formRedirect, "//") {
		redirectPath = formRedirect
		// Ensure the redirect has the correct language suffix if needed
		if isEnglish && !strings.HasSuffix(redirectPath, "/e") && !strings.Contains(redirectPath, "/e/") {
			redirectPath = redirectPath + "/e"
		}
	} else {
		// Fallback to session
		session, err := a.authService.sessionManager.store.Get(r, "login-redirect")
		if err == nil {
			if redirect, ok := session.Values["redirect"].(string); ok && redirect != "" {
				// Validate redirect path (must be internal)
				if strings.HasPrefix(redirect, "/") && !strings.HasPrefix(redirect, "//") {
					redirectPath = redirect
					// Ensure the redirect has the correct language suffix if needed
					if isEnglish && !strings.HasSuffix(redirectPath, "/e") && !strings.Contains(redirectPath, "/e/") {
						redirectPath = redirectPath + "/e"
					}
				}
				// Clear the redirect from session
				delete(session.Values, "redirect")
				session.Save(r, w)
			}
		}
	}

	// Redirect after successful login
	http.Redirect(w, r, redirectPath, http.StatusFound)
}

// handleSignup processes user registration
func (a *App) handleSignup(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "无效的表单数据", http.StatusBadRequest)
		return
	}

	username := strings.TrimSpace(r.Form.Get("username"))
	email := strings.TrimSpace(r.Form.Get("email"))
	password := r.Form.Get("password")
	role := r.Form.Get("role")

	if username == "" || email == "" || password == "" {
		http.Error(w, "所有字段都是必需的", http.StatusBadRequest)
		return
	}

	if len(password) < 8 {
		http.Error(w, "密码必须至少 8 个字符", http.StatusBadRequest)
		return
	}

	// Default role for new registrations is "writer"
	userRole := RoleWriter
	if role == "editor" {
		// Only allow creating editor accounts with the correct editor password
		editorPassword := r.Form.Get("editor_password")
		if editorPassword != "yanqing" {
			http.Error(w, "Incorrect editor password. Only authorized users can create editor accounts.", http.StatusForbidden)
			return
		}
		userRole = RoleEditor
	}

	if err := a.authService.userStore.CreateUser(username, email, password, userRole); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			http.Error(w, "用户名已存在", http.StatusConflict)
		} else {
			log.Printf("Error creating user: %v", err)
			http.Error(w, "内部服务器错误", http.StatusInternalServerError)
		}
		return
	}

	// Auto-login after successful signup
	user, err := a.authService.userStore.ValidateUser(username, password)
	if err != nil {
		log.Printf("Error validating new user: %v", err)
		http.Error(w, "内部服务器错误", http.StatusInternalServerError)
		return
	}

	if err := a.authService.sessionManager.SetSession(w, r, user); err != nil {
		log.Printf("Error setting session after signup: %v", err)
		http.Error(w, "内部服务器错误", http.StatusInternalServerError)
		return
	}

	// Return success response for AJAX
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// handleLogout processes user logout
func (a *App) handleLogout(w http.ResponseWriter, r *http.Request) {
	if err := a.authService.sessionManager.ClearSession(w, r); err != nil {
		log.Printf("Error clearing session: %v", err)
	}

	// Determine if this request came from an English page
	referer := r.Header.Get("Referer")
	isEnglish := strings.Contains(referer, "/e")
	langSuffix := ""
	if isEnglish {
		langSuffix = "/e"
	}

	// Check if there's a return_to parameter or use Referer
	returnTo := r.URL.Query().Get("return_to")
	if returnTo == "" {
		// Use Referer header if available
		if referer != "" {
			// Extract path from referer URL
			if strings.Contains(referer, r.Host) || strings.Contains(referer, "localhost") || strings.Contains(referer, "ngrok") {
				// Parse the referer to get just the path
				parts := strings.Split(referer, "//")
				if len(parts) > 1 {
					pathParts := strings.SplitN(parts[1], "/", 2)
					if len(pathParts) > 1 {
						returnTo = "/" + pathParts[1]
					}
				}
			}
		}
	}

	// If no return path or it's a protected page, use landing with language suffix
	if returnTo == "" || returnTo == "/login" || returnTo == "/signup" || returnTo == "/payment" ||
		returnTo == "/login/e" || returnTo == "/signup/e" || returnTo == "/payment/e" {
		returnTo = "/landing" + langSuffix
	}

	http.Redirect(w, r, returnTo, http.StatusFound)
}

// getCurrentUser gets the current user session data
func (a *App) getCurrentUser(r *http.Request) *SessionData {
	sessionData, err := a.authService.sessionManager.GetSession(r)
	if err != nil {
		log.Printf("Error getting session: %v", err)
		return nil
	}
	return sessionData
}

// requireAuth checks if user is authenticated for API endpoints
func (a *App) requireAuth(w http.ResponseWriter, r *http.Request) *SessionData {
	sessionData := a.getCurrentUser(r)
	if sessionData == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return nil
	}
	return sessionData
}

// requireWriterOrEditor checks if user has at least writer permissions
func (a *App) requireWriterOrEditor(w http.ResponseWriter, r *http.Request) *SessionData {
	sessionData := a.requireAuth(w, r)
	if sessionData == nil {
		return nil
	}

	if sessionData.Role != RoleWriter && sessionData.Role != RoleEditor {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return nil
	}

	return sessionData
}

// requireEditor checks if user is an editor (for approval/QC operations)
func (a *App) requireEditor(w http.ResponseWriter, r *http.Request) *SessionData {
	sessionData := a.requireAuth(w, r)
	if sessionData == nil {
		return nil
	}

	if sessionData.Role != RoleEditor {
		http.Error(w, "Forbidden: Editor access required", http.StatusForbidden)
		return nil
	}

	return sessionData
}

// handleRecoverPassword handles password recovery requests
// This is a simple development-only feature that resets password to a default
func (a *App) handleRecoverPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
		return
	}

	email := strings.TrimSpace(request.Email)
	if email == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Email is required"})
		return
	}

	username, newPassword, err := a.authService.userStore.ResetPasswordByEmail(email)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "No account found with that email"})
		return
	}

	log.Printf("Password reset for user %s (email: %s)", username, email)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"username": username,
		"password": newPassword,
		"message":  "Password has been reset successfully",
	})
}
