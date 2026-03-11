package web

import (
	"net/http"
)

// AuthMiddleware provides authentication-related middleware
type AuthMiddleware struct {
	authService *AuthService
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(authService *AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// RequireAuth middleware ensures the user is authenticated
func (am *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionData, err := am.authService.sessionManager.GetSession(r)
		if err != nil || sessionData == nil {
			// Redirect to login page
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// RequireRole middleware ensures the user has a specific role
func (am *AuthMiddleware) RequireRole(role UserRole) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			sessionData, err := am.authService.sessionManager.GetSession(r)
			if err != nil || sessionData == nil {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			if sessionData.Role != role {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}

// RequireWriterOrEditor middleware ensures the user is at least a writer
func (am *AuthMiddleware) RequireWriterOrEditor(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionData, err := am.authService.sessionManager.GetSession(r)
		if err != nil || sessionData == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		if sessionData.Role != RoleWriter && sessionData.Role != RoleEditor {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// RequireEditor middleware ensures the user is an editor
func (am *AuthMiddleware) RequireEditor(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionData, err := am.authService.sessionManager.GetSession(r)
		if err != nil || sessionData == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		if sessionData.Role != RoleEditor {
			http.Error(w, "Forbidden: Editor access required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// InjectUserContext middleware adds user context to all requests
func (am *AuthMiddleware) InjectUserContext(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// This middleware doesn't block - it just makes session data available
		// Individual handlers can check for authentication as needed
		next.ServeHTTP(w, r)
	}
}

// RequireWriterOrEditorAPI middleware for API endpoints (returns JSON errors, not redirects)
func (am *AuthMiddleware) RequireWriterOrEditorAPI(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionData, err := am.authService.sessionManager.GetSession(r)
		if err != nil || sessionData == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"success":false,"error":"Authentication required. Please login to perform this action."}`))
			return
		}

		if sessionData.Role != RoleWriter && sessionData.Role != RoleEditor {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"success":false,"error":"Forbidden: Writer or Editor access required."}`))
			return
		}

		next.ServeHTTP(w, r)
	}
}

// RequireEditorAPI middleware for API endpoints (returns JSON errors, not redirects)
func (am *AuthMiddleware) RequireEditorAPI(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionData, err := am.authService.sessionManager.GetSession(r)
		if err != nil || sessionData == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"success":false,"error":"Authentication required. Please login to perform this action."}`))
			return
		}

		if sessionData.Role != RoleEditor {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"success":false,"error":"Forbidden: Editor access required."}`))
			return
		}

		next.ServeHTTP(w, r)
	}
}