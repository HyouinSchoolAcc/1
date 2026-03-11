# CLAUDE.md

Guidance for AI assistants working with this codebase.

## Project Overview

**Divergence 2% Writer Portal** — A Go web application for managing conversational data labeling. Writers create dialogue data for AI characters, editors review quality, and the system provides collaboration tools.

## Quick Reference

**Start server (Linux/Mac):** `./run_go3.sh` (from parent directory)  
**Start server (Windows):** `run_go3.bat` (from this directory)  
**Stop services (Linux/Mac):** `./stop_services.sh` (from this directory)  
**Stop services (Windows):** `stop_services.bat` (from this directory)  
**Port:** 5002  
**Public URL:** https://wl2.studio  
**Database:** SQLite at `data/app.db`

## Architecture

```
cmd/server/main.go          # Entry point
internal/
├── database/               # SQLite data layer
│   ├── schema.sql          # Table definitions
│   ├── user_store.go       # Authentication
│   ├── chat_store.go       # Discord-style chat
│   ├── lounge_store.go     # Writers' Lounge
│   └── tutorial_store.go   # Onboarding
├── email/                  # Email confirmations
└── web/
    ├── app.go              # Application state, helpers
    ├── router.go           # Route registration
    ├── handlers.go         # Page handlers, template loading
    ├── auth.go             # Session management
    ├── api.go              # Core REST endpoints (~1700 lines)
    ├── chat_api.go         # Discord-style features
    ├── lounge_api.go       # Writers' Lounge
    ├── payment_api.go      # Earnings tracking
    └── llm_api.go          # vLLM integration
```

## Characters

| ID | Route | Display Name |
|----|-------|--------------|
| kurisu | /kurisu | Kurisu / 牧濑红莉栖 |
| kafka | /kafka | Kafka / 卡夫卡 |
| lin_lu | /linlu | Lin Lu / 林路 |
| newcharacter_1 | /newcharacter_1 | Community voting |

**Preset folders:** `presets/presets_{id}` and `presets/presets_{id}_CN`

## Key Patterns

### Template Loading
Templates use Jinja-style `{% include %}` expanded by Go:
```go
tmpl, err := loadGoTemplate(cfg.TemplatesDir, "page.html")
```

### Route Protection
```go
r.Get("/protected", authMiddleware.RequireWriterOrEditor(handler))
```

### User Roles
- `new_user` — View only
- `writer` — Create/edit dialogues
- `editor` — Full access + QC approval

### API Response Format
```go
// Success
json.NewEncoder(w).Encode(map[string]any{"success": true, "data": result})

// Error
json.NewEncoder(w).Encode(map[string]any{"error": "message"})
```

## Common Tasks

### Add New Character
1. Create preset folder in `presets/`
2. Add to `presets` slice in `app.go` NewApp()
3. Add to `loadCharacters()` in `handlers.go`
4. Add route in `router.go`
5. Update helper functions in `app.go`

### Add API Endpoint
1. Add handler in appropriate `*_api.go`
2. Register in corresponding `Register*API()` function
3. Follow existing patterns for auth checks and response format

### Database Changes
1. Add table in `internal/database/schema.sql`
2. Create `*_store.go` with CRUD operations
3. Initialize store in `app.go` NewApp()

## File Locations

| What | Where |
|------|-------|
| User templates | `templates/` |
| Static assets | `static/` |
| Character data | `presets/presets_*/` |
| Database | `data/app.db` |
| Server logs | `server.log` |

## Startup Scripts

### Start Services

#### Linux/Mac: `run_go3.sh`
Run from parent directory: `./run_go3.sh`
- Automatically handles all dependencies and cleanup
- Starts server on port 5002 with DEBUG mode
- Launches ngrok tunnel to https://wl2.studio
- Creates/migrates SQLite database if needed
- Does NOT start vLLM (LLM features disabled)

#### Windows: `run_go3.bat`
Run from this directory: `run_go3.bat`
- Windows equivalent of run_go3.sh
- Same functionality: server + ngrok, no vLLM
- Automatically checks for Go, ngrok, and sqlite3
- Kills existing processes on ports 5002/5003
- Keeps running in foreground (Ctrl+C to stop)
- Creates logs: `server.log` and `ngrok.log`

### Stop Services

#### Linux/Mac: `stop_services.sh`
Run from this directory: `./stop_services.sh`
- Stops all Go server processes
- Stops ngrok tunnel
- Kills processes on ports 5002/5003
- Safe cleanup of all background services

#### Windows: `stop_services.bat`
Run from this directory: `stop_services.bat`
- Windows equivalent of stop_services.sh
- Stops server_sql.exe and ngrok.exe
- Cleans up ports 5002/5003
- Kills background processes
- Press any key to close after completion

**Manual stop options:**
- Linux/Mac: `kill [PID]` (PIDs shown in startup output)
- Windows: Ctrl+C in the command window, or `taskkill /F /IM server_sql.exe`

**Requirements:**
- Go 1.21+ (https://go.dev/dl/)
- ngrok (https://ngrok.com/download)
- sqlite3 CLI (optional, for database inspection)

## Don't Forget

- Bilingual routes: most pages have `/page` and `/page/e` variants
- Character presets come in pairs: `presets_name` and `presets_name_CN`
- `bankclerk` is deprecated — use `lin_lu` / `linlu`
- Use `run_go3.sh` / `run_go3.bat` to start services
- Use `stop_services.sh` / `stop_services.bat` to stop all services
- Hard refresh browser after template changes
