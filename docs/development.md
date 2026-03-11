# Development Guide

## Architecture Overview

The application is a Go web server using the Chi router. It serves HTML templates with embedded JavaScript for interactivity.

### Core Components

```
internal/
├── database/           # Data access layer
│   ├── db.go          # Database connection
│   ├── schema.sql     # Table definitions
│   ├── user_store.go  # User CRUD operations
│   ├── lounge_store.go
│   ├── chat_store.go
│   └── tutorial_store.go
├── email/             # Email service for confirmations
└── web/               # HTTP layer
    ├── app.go         # Application state
    ├── router.go      # Route registration
    ├── handlers.go    # Page handlers
    ├── auth.go        # Authentication
    ├── api.go         # Core REST endpoints
    ├── chat_api.go    # Discord-style chat
    ├── lounge_api.go  # Writers' Lounge
    ├── payment_api.go # Payment tracking
    └── llm_api.go     # LLM integration
```

## Templates

Go templates with Jinja-style include syntax. The `expandJinjaIncludes` function in `handlers.go` processes `{% include "path" %}` directives.

### Template Files

| File | Purpose |
|------|---------|
| `landing.html`, `landing_en.html` | Home page |
| `writing.html` | Character selection |
| `writing_main.html` | Dialogue editor workspace |
| `lounge.html` | Writers' Lounge |
| `faq.html`, `guide.html` | Information pages |

### Shared Components

Templates in `templates/includes/`:
- `header.html` — Top navigation
- `sidebar.html` — Left navigation and user status
- `footer.html` — Page footer

## Adding a New Character

1. **Create preset folder:**
   ```bash
   mkdir -p presets/presets_newchar
   mkdir -p presets/presets_newchar_CN
   ```

2. **Add to app.go presets list:**
   ```go
   presets := []string{
       // ... existing presets
       "presets_newchar",
       "presets_newchar_CN",
   }
   ```

3. **Add character definition in handlers.go:**
   ```go
   {
       ID:    "newchar",
       Route: "/newchar",
       Name: "Character Name",
       // ... other fields
   }
   ```

4. **Register route in router.go:**
   ```go
   r.Get("/newchar", app.handleTemplate(writingMainTpl))
   ```

5. **Add character helper functions in app.go:**
   - Update `getCurrentCharacterRole()`
   - Update `getCurrentCharacterDisplayName()`
   - Update `fallbackSchedule()`

## API Conventions

### Response Format

```json
{
    "success": true,
    "data": { ... },
    "error": ""
}
```

### Authentication Check

Use middleware for protected routes:
```go
r.Get("/protected", authMiddleware.RequireWriterOrEditor(handler))
```

Or check in handler:
```go
user := a.getCurrentUser(r)
if user == nil || user.Role == RoleNewUser {
    http.Error(w, "Unauthorized", http.StatusUnauthorized)
    return
}
```

## Database Operations

### Adding a New Table

1. Add schema in `internal/database/schema.sql`
2. Create store file (e.g., `feature_store.go`)
3. Initialize in `app.go`

### Store Pattern

```go
type FeatureStore struct {
    db *DB
}

func NewFeatureStore(db *DB) *FeatureStore {
    return &FeatureStore{db: db}
}

func (s *FeatureStore) Create(item *Feature) error {
    _, err := s.db.Exec(`INSERT INTO features ...`, ...)
    return err
}
```

## Character Presets

Each character has a preset folder structure:

```
presets/presets_character_CN/
├── hidden_users.json           # Users to hide from listings
├── new_user_info_*.json        # Default user template
└── user_*_Day*_dup_*_simplified.json  # Dialogue files
```

> **Note:** Character stories and schedules are stored in `data/character_profiles.json` (the same source as https://wl2.studio/descriptions), not in the preset folders.

### Dialogue File Format

```json
{
    "user_name": "用户名",
    "character_name": "角色名",
    "user_schedule": { "day": 1, "早晨": "..." },
    "character_schedule": { "day": 1, "早晨": "..." },
    "dialogues": [
        { "speaker": "user", "text": "..." },
        { "speaker": "ai", "text": "..." }
    ],
    "qc_status": "pending"
}
```

## Testing

### Manual Testing

1. Start server: `./run_go3.sh`
2. Open browser to http://localhost:5002
3. Test routes and features

### Database Testing

```bash
# Check database state
sqlite3 data/app.db "SELECT * FROM users;"

# Reset database (caution!)
rm data/app.db
# Server will recreate on next start
```

## Debugging

### Enable Debug Mode

```bash
DEBUG=true ./server_sql
```

### View Logs

```bash
tail -f server.log
```

### Common Issues

**Port in use:**
```bash
lsof -ti:5002 | xargs kill -9
```

**Database locked:**
```bash
# Kill all server instances
pkill -f server_sql
```

**Template not updating:**
Hard refresh browser (Ctrl+Shift+R)
