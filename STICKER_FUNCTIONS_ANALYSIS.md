# Sticker Functionality - Complete Function & Folder Analysis

## Overview
**Note:** There is no function explicitly named `add_sticker`. The main function that adds stickers to conversations is **`selectSticker()`** in JavaScript.

---

## 📁 FOLDERS

### 1. `/stickers/` - Sticker Data Directory
- **Location:** `data_labler_UI_production/stickers/`
- **Contents:**
  - `stickers_map.json` - Maps sticker filenames to metadata (emotion, description)
  - `GIFs/` - Directory containing all sticker GIF image files

### 2. `/templates/` - Frontend Templates
- **Location:** `data_labler_UI_production/templates/`
- **Relevant Files:**
  - `writing_main.html` - Main writing page with sticker JavaScript functions
  - `includes/modals.html` - Sticker picker modal UI
  - `includes/main_panel.html` - Sticker button in controls bar

### 3. `/internal/web/` - Backend Go Files
- **Location:** `data_labler_UI_production/internal/web/`
- **Relevant Files:**
  - `api.go` - Contains sticker API handlers
  - `router.go` - Contains sticker file serving route

---

## 🔧 FUNCTIONS

### Frontend JavaScript Functions (in `templates/writing_main.html`)

#### 1. **`selectSticker(sticker)`** - Main Function (Lines 8576-8602)
   - **Purpose:** Adds a sticker to the conversation
   - **Location:** `templates/writing_main.html:8576`
   - **What it does:**
     - Checks if user is new (prevents new users from adding stickers)
     - Creates a new turn object with `content_type: "sticker"`
     - Adds sticker metadata (filename, emotion, description)
     - Pushes to `currentWriterFileData.dialogue` array
     - Refreshes the UI and closes the picker modal
   - **Called by:** `renderStickerGrid()` when user clicks a sticker item

#### 2. **`showStickerPicker()`** - Opens Sticker Modal (Lines 8526-8531)
   - **Purpose:** Displays the sticker picker modal
   - **Location:** `templates/writing_main.html:8526`
   - **What it does:**
     - Loads stickers if not already loaded
     - Shows the sticker picker modal
   - **Called by:** Button click in `templates/includes/main_panel.html:18`

#### 3. **`hideStickerPicker()`** - Closes Sticker Modal (Lines 8533-8535)
   - **Purpose:** Hides the sticker picker modal
   - **Location:** `templates/writing_main.html:8533`
   - **Called by:** Modal close button and `selectSticker()`

#### 4. **`loadStickers()`** - Fetches Stickers from API (Lines 8511-8524)
   - **Purpose:** Loads available stickers from backend
   - **Location:** `templates/writing_main.html:8511`
   - **What it does:**
     - Fetches from `/api/stickers` endpoint
     - Populates `availableStickers` array
     - Calls `renderStickerGrid()` to display stickers
   - **Called by:** `showStickerPicker()` if stickers not loaded

#### 5. **`renderStickerGrid()`** - Displays Sticker Grid (Lines 8548-8574)
   - **Purpose:** Renders the sticker grid in the modal
   - **Location:** `templates/writing_main.html:8548`
   - **What it does:**
     - Creates DOM elements for each sticker
     - Sets up click handlers that call `selectSticker()`
     - Displays sticker images and emotions
   - **Called by:** `loadStickers()` after fetching data

#### 6. **`selectStickerRole(role)`** - Sets User/Character Role (Lines 8537-8546)
   - **Purpose:** Sets which role (User or Character) will send the sticker
   - **Location:** `templates/writing_main.html:8537`
   - **What it does:**
     - Updates `selectedStickerRole` variable
     - Updates UI button states
   - **Called by:** Role selection buttons in modal (`templates/includes/modals.html:741,744`)

#### 7. **Sticker Display Rendering** (Lines 2413-2453)
   - **Purpose:** Renders stickers in the conversation view
   - **Location:** `templates/writing_main.html:2413`
   - **What it does:**
     - Checks if `turnData.content_type === 'sticker'`
     - Creates sticker container with image and metadata
     - Displays sticker GIF and emotion/description info
   - **Called by:** Conversation rendering logic when displaying turns

### Backend Go Functions (in `internal/web/api.go`)

#### 1. **`handleGetStickers(w, r)`** - API Handler (Lines 2974-2999)
   - **Purpose:** Returns list of all available stickers with metadata
   - **Location:** `internal/web/api.go:2974`
   - **Endpoint:** `GET /api/stickers`
   - **What it does:**
     - Reads `stickers/stickers_map.json`
     - Converts map to array format
     - Adds URL paths for each sticker
     - Returns JSON response with stickers array
   - **Called by:** Frontend `loadStickers()` function

#### 2. **`handleValidateSticker(w, r)`** - Validation Handler (Lines 3001-3043)
   - **Purpose:** Validates sticker exists before saving
   - **Location:** `internal/web/api.go:3001`
   - **Endpoint:** `POST /api/stickers/validate`
   - **Authentication:** Requires writer/editor auth (`writerAuth` middleware)
   - **What it does:**
     - Checks if sticker exists in `stickers_map.json`
     - Verifies physical file exists in `stickers/GIFs/`
     - Returns sticker metadata if valid
   - **Called by:** Backend validation (likely during save operations)

### Router Configuration (in `internal/web/router.go`)

#### 1. **Sticker File Serving Route** (Lines 186-187)
   - **Purpose:** Serves sticker GIF files
   - **Location:** `internal/web/router.go:186`
   - **Route:** `/static/stickers/*`
   - **What it does:**
     - Maps URL path to `stickers/` directory
     - Serves static GIF files
   - **Used by:** Frontend to display sticker images

---

## 🔗 FUNCTION CALL FLOW

### Adding a Sticker:
1. User clicks "Add Sticker" button → calls `showStickerPicker()`
2. `showStickerPicker()` → calls `loadStickers()` if needed
3. `loadStickers()` → fetches `/api/stickers` → calls `renderStickerGrid()`
4. `renderStickerGrid()` → creates clickable sticker items
5. User clicks sticker → calls `selectSticker(sticker)`
6. `selectSticker()` → adds turn to dialogue → calls `hideStickerPicker()`
7. UI refreshes → sticker display logic renders the sticker

### Displaying Stickers:
1. Conversation rendering checks `content_type === 'sticker'`
2. Creates sticker container with image and metadata
3. Displays sticker GIF from `/static/stickers/GIFs/{filename}`

---

## 📊 DATA STRUCTURES

### Sticker Object (in dialogue):
```javascript
{
  role: "User" | "Character",
  content: "",
  content_type: "sticker",
  sticker: {
    filename: "example.gif",
    emotion: "happy",
    description: "description text"
  },
  intimacy_delta: 0.01
}
```

### Sticker Metadata (from API):
```javascript
{
  filename: "example.gif",
  emotion: "happy",
  description: "description text",
  url: "/static/stickers/GIFs/example.gif"
}
```

### Stickers Map JSON:
```json
{
  "filename.gif": {
    "emotion": "happy",
    "description": "description text"
  }
}
```

---

## 🎯 KEY FILES SUMMARY

| File | Purpose | Key Functions/Lines |
|------|---------|---------------------|
| `templates/writing_main.html` | Main sticker JavaScript logic | `selectSticker()`, `loadStickers()`, `showStickerPicker()`, `renderStickerGrid()`, `selectStickerRole()`, sticker display (2413-2453) |
| `templates/includes/modals.html` | Sticker picker modal UI | Modal HTML (732-752), CSS styles (675-704) |
| `templates/includes/main_panel.html` | Sticker button | Button with `onclick="showStickerPicker()"` (16-20) |
| `internal/web/api.go` | Backend API handlers | `handleGetStickers()` (2974-2999), `handleValidateSticker()` (3001-3043) |
| `internal/web/router.go` | File serving route | Sticker route (186-187) |
| `stickers/stickers_map.json` | Sticker metadata | Maps filenames to emotions/descriptions |
| `stickers/GIFs/` | Sticker images | All sticker GIF files |

---

## 📝 NOTES

- The main function that **adds** stickers is `selectSticker()` (not `add_sticker`)
- Stickers are added directly to the dialogue array in the frontend
- Backend validation is available via `handleValidateSticker()` but may not be called by `selectSticker()`
- Sticker files are served statically via the router
- New users are prevented from adding stickers (checked in `selectSticker()`)
