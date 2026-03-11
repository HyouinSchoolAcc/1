# Divergence 2% Writer Portal

A web application for managing conversational data labeling with multiple AI character presets. Writers create dialogue data, editors review quality, and the system provides tools for collaboration, certification, and payment tracking.

## Quick Start

```bash
# Windows
run_server.bat
```

This starts the server on port 5002 and creates an ngrok tunnel at `https://wl2.studio`.

**Local URL:** http://localhost:5002  
**Public URL:** https://wl2.studio

## Test Credentials

| Role | Username | Password |
|------|----------|----------|
| Editor | `testeditor` | `password123` |
| Writer | `testwriter` | `password123` |

## Available Characters

| Character | Route | Source | Description |
|-----------|-------|--------|-------------|
| Kurisu (牧濑红莉栖) | `/kurisu` | Steins;Gate | Time travel researcher |
| Lin Lu (林路) | `/linlu` | Original | Literature professor |
| New Character | `/newcharacter_1` | Community | Vote for next character |

### Character Data

All character stories, schedules, and profile information live in **`data/character_profiles.json`**. This is the single source of truth and is the same data displayed at https://wl2.studio/descriptions.

The `presets/presets_*/` directories hold **per-user dialogue files and user profile data only** — they no longer contain character stories or schedules.

## Features

### Core Functionality
- **Dialogue Editor** — Create and edit conversation data with AI characters
- **Quality Control** — Editors review, approve, or request revisions
- **Character Schedules** — Daily activity templates sourced from `data/character_profiles.json`
- **Bilingual Support** — All pages available in Chinese (`/page`) and English (`/page/e`)

### Media Generation
- **Dialogue Image Export** — Generate PNG images from dialogue files
- **Video Generation** — Render dialogue as chat-style videos (WeChat/iMessage-like)

### Community Features
- **Writers' Lounge** (`/lounge`) — Forum for writers to share tips and chat
- **Discord-style Chat** — Real-time channels with reactions and DMs
- **Character Voting** — Community decides which character to develop next

### Writer Management
- **Tutorial System** — Onboarding for new writers
- **Certification** — Writers certify for specific characters
- **Payment Tracking** — Earnings and payment history

### AI Integration
- **LLM Chat** (`/llm`) — Chat with fine-tuned Qwen3-14B model

## Project Structure

```
data_labeler/
├── cmd/server/main.go          # Application entry point
├── internal/
│   ├── database/               # SQLite database layer
│   │   ├── schema.sql          # Table definitions
│   │   ├── user_store.go       # Authentication
│   │   ├── chat_store.go       # Discord-style chat
│   │   ├── lounge_store.go     # Writers' Lounge
│   │   └── tutorial_store.go   # Onboarding
│   ├── email/                  # Email confirmations
│   └── web/                    # HTTP handlers and routing
│       ├── app.go              # Application state, helpers
│       ├── router.go           # Route registration
│       ├── api.go              # Core REST endpoints
│       ├── chat_api.go         # Discord-style features
│       ├── lounge_api.go       # Writers' Lounge
│       ├── payment_api.go      # Earnings tracking
│       └── llm_api.go          # vLLM integration
├── templates/                  # HTML templates
├── static/                     # CSS, JS, images
├── presets/                    # User dialogue files & per-character user profiles
├── data/character_profiles.json # Source of truth for character stories & schedules
├── stickers/                   # Sticker assets and metadata
├── video_chat_renderer/        # Python video renderer
├── convert_dialogue_to_image.py # Dialogue-to-PNG generator
├── data/                       # SQLite database (app.db)
├── docs/                       # API, development, deployment docs
├── go.mod / go.sum             # Go module files
├── server_sql.exe              # Compiled server binary (Windows)
├── run_server.bat              # Start server + ngrok (Windows)
└── stop_services.bat           # Stop all services (Windows)
```

## Database

SQLite database at `data/app.db`. Key tables:

| Table | Purpose |
|-------|---------|
| `users` | User accounts and authentication |
| `lounge_posts`, `lounge_replies` | Writers' Lounge content |
| `chat_channels`, `chat_messages` | Discord-style chat |
| `tutorial_progress` | Writer onboarding progress |
| `character_certifications` | Writer-character permissions |

## User Roles

| Role | Permissions |
|------|-------------|
| `new_user` | View public pages only |
| `writer` | Create dialogues, post in lounge, edit own work |
| `editor` | All writer permissions + QC approval, archive, admin functions |

## Critical Rules

> **NEVER modify production (`data_labeler/`) directly.** All changes must be made in `data_labeler_testing_env/` and deployed to production via `deploy_to_production.bat`. This applies to code, templates, static assets, and configuration. The deploy script is the only sanctioned path to production.

## Development

### Prerequisites
- Go 1.21+
- SQLite3
- ngrok (with auth token configured)
- Python 3 (for dialogue image/video generation)
- FFmpeg (for video rendering)

### Build and Run

```bash
go mod download
go build -o server_sql ./cmd/server
./server_sql
```

### Stopping Services

```bash
# Windows
stop_services.bat
```

---

**Version:** 2.1.0  
**Last Updated:** March 2026
