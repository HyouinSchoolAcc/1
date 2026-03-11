package web

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const LanguageKey contextKey = "language"

// LanguageMiddleware handles language detection using /e suffix pattern
// - /page = Chinese
// - /page/e = English
func LanguageMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lang := "zh" // default language
		originalPath := r.URL.Path
		modifiedPath := originalPath

		// Check if path ends with /e (English suffix)
		if strings.HasSuffix(originalPath, "/e") && len(originalPath) > 2 {
			lang = "en"
			// Strip /e suffix for internal routing
			modifiedPath = strings.TrimSuffix(originalPath, "/e")
		} else if originalPath == "/e" {
			// Just /e goes to English landing
			lang = "en"
			modifiedPath = "/"
		}

		// Set language cookie for future requests
		http.SetCookie(w, &http.Cookie{
			Name:     "lang",
			Value:    lang,
			Path:     "/",
			MaxAge:   365 * 24 * 60 * 60, // 1 year
			HttpOnly: false,               // Allow JS access for dynamic content
			SameSite: http.SameSiteLaxMode,
		})

		// Store language in request context
		ctx := context.WithValue(r.Context(), LanguageKey, lang)

		// Update the request path for routing (strip /e suffix)
		if modifiedPath != originalPath {
			r.URL.Path = modifiedPath
		}

		// Serve the request with updated context and path
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetLanguage extracts language from request context
func GetLanguage(r *http.Request) string {
	if lang, ok := r.Context().Value(LanguageKey).(string); ok {
		return lang
	}
	return "zh" // default fallback
}

// RedirectLegacyLanguageURLs is now a no-op pass-through (kept for compatibility)
// The /e suffix pattern is now the primary pattern, not legacy
func RedirectLegacyLanguageURLs(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
