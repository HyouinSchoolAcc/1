package web

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	maxEditorAttempts  = 5
	editorLockoutTime  = 15 * time.Minute
	cleanupInterval    = 30 * time.Minute
)

type ipRecord struct {
	failures  int
	blockedAt time.Time
}

// EditorRateLimiter tracks failed editor-password attempts per IP
// and blocks IPs that exceed the threshold.
type EditorRateLimiter struct {
	mu      sync.Mutex
	records map[string]*ipRecord
}

func NewEditorRateLimiter() *EditorRateLimiter {
	rl := &EditorRateLimiter{records: make(map[string]*ipRecord)}
	go rl.cleanupLoop()
	return rl
}

func (rl *EditorRateLimiter) IsBlocked(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rec, ok := rl.records[ip]
	if !ok {
		return false
	}
	if rec.failures >= maxEditorAttempts {
		if time.Since(rec.blockedAt) < editorLockoutTime {
			return true
		}
		// Lockout expired — reset
		delete(rl.records, ip)
		return false
	}
	return false
}

func (rl *EditorRateLimiter) RecordFailure(ip string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rec, ok := rl.records[ip]
	if !ok {
		rec = &ipRecord{}
		rl.records[ip] = rec
	}
	rec.failures++
	if rec.failures >= maxEditorAttempts {
		rec.blockedAt = time.Now()
	}
}

func (rl *EditorRateLimiter) AttemptsRemaining(ip string) int {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rec, ok := rl.records[ip]
	if !ok {
		return maxEditorAttempts
	}
	remaining := maxEditorAttempts - rec.failures
	if remaining < 0 {
		return 0
	}
	return remaining
}

func (rl *EditorRateLimiter) ClearFailures(ip string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	delete(rl.records, ip)
}

func (rl *EditorRateLimiter) cleanupLoop() {
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		for ip, rec := range rl.records {
			if rec.failures >= maxEditorAttempts && time.Since(rec.blockedAt) >= editorLockoutTime {
				delete(rl.records, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// getRealIP extracts the client IP, respecting X-Forwarded-For / X-Real-IP.
func getRealIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.SplitN(xff, ",", 2)
		return strings.TrimSpace(parts[0])
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
