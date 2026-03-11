# Sticker Function - Transfer Package

This folder contains ONLY the sticker functionality code and assets, isolated for transfer to production server.

## 📁 Folder Structure

```
sticker_function_isolated/
├── cmd/server/
│   └── main.go                    # Server entry point
├── internal/web/
│   ├── api.go                     # Contains handleGetStickers() and handleValidateSticker()
│   ├── router.go                  # Contains sticker file serving route
│   └── [supporting Go files]      # Required for compilation
├── templates/
│   ├── writing_main.html          # Contains sticker JavaScript functions
│   └── includes/
│       ├── modals.html            # Sticker picker modal UI
│       └── main_panel.html        # Sticker button
├── stickers/
│   ├── stickers_map.json          # Sticker metadata mapping
│   └── GIFs/                      # All sticker GIF image files
├── go.mod                         # Go module file
├── go.sum                         # Go checksums
└── README_TRANSFER.md             # This file
```

## 🔧 Key Functions

### Backend (Go)
- `handleGetStickers()` - Returns list of available stickers (api.go:2974-2999)
- `handleValidateSticker()` - Validates sticker exists (api.go:3001-3043)
- Sticker file serving route (router.go:186-187)

### Frontend (JavaScript)
- `selectSticker()` - Adds sticker to conversation (writing_main.html:8576-8602)
- `showStickerPicker()` - Opens sticker picker modal (writing_main.html:8526-8531)
- `loadStickers()` - Fetches stickers from API (writing_main.html:8511-8524)
- `renderStickerGrid()` - Displays sticker grid (writing_main.html:8548-8574)
- `selectStickerRole()` - Sets User/Character role (writing_main.html:8537-8546)
- Sticker display rendering (writing_main.html:2413-2453)

## 📋 Transfer Instructions

1. **Copy the entire folder** to your production server
2. **Merge files** into your existing codebase:
   - Merge `internal/web/api.go` - add sticker handler functions
   - Merge `internal/web/router.go` - add sticker route
   - Merge `templates/writing_main.html` - add sticker JavaScript functions
   - Merge `templates/includes/modals.html` - add sticker modal HTML/CSS
   - Merge `templates/includes/main_panel.html` - add sticker button
   - Copy `stickers/` folder entirely to production
3. **Ensure dependencies** are met:
   - Go Chi router (github.com/go-chi/chi/v5)
   - All supporting Go files are present
4. **Test endpoints**:
   - GET /api/stickers
   - POST /api/stickers/validate
   - /static/stickers/* (file serving)

## 📍 File Locations in Source Code

See STICKER_FUNCTIONS_ANALYSIS.md for detailed line numbers and function locations.

## ⚠️ Important Notes

- The copied `api.go` and `router.go` files contain OTHER functionality too
- You'll need to extract ONLY the sticker-related parts when merging
- See STICKER_ELEMENTS_GUIDE.txt for exact line numbers of sticker code
- Ensure authentication middleware is properly configured for `/api/stickers/validate`
