#!/bin/bash

# Script to isolate sticker functionality into a transfer-ready folder
# This creates a clean folder with ONLY sticker-related code and assets

SOURCE_DIR="/home/exx/Desktop/fine-tune/data_labler_UI_production"
DEST_DIR="/home/exx/Desktop/fine-tune/data_labler_UI_production/sticker_function_isolated"

echo "=========================================="
echo "Isolating Sticker Function for Transfer"
echo "=========================================="
echo ""

# Remove existing folder if it exists
if [ -d "$DEST_DIR" ]; then
    echo "Removing existing folder..."
    rm -rf "$DEST_DIR"
fi

# Create directory structure
echo "Creating directory structure..."
mkdir -p "$DEST_DIR"
mkdir -p "$DEST_DIR/templates/includes"
mkdir -p "$DEST_DIR/internal/web"
mkdir -p "$DEST_DIR/stickers/GIFs"
mkdir -p "$DEST_DIR/cmd/server"

# Copy Go module files (needed for compilation)
echo ""
echo "📦 Copying Go module files..."
cp "$SOURCE_DIR/go.mod" "$DEST_DIR/" 2>/dev/null && echo "  ✓ go.mod" || echo "  ✗ go.mod (not found)"
cp "$SOURCE_DIR/go.sum" "$DEST_DIR/" 2>/dev/null && echo "  ✓ go.sum" || echo "  ✗ go.sum (not found)"

# Copy sticker-related Go files
echo ""
echo "🔧 Copying Go backend files..."
cp "$SOURCE_DIR/internal/web/api.go" "$DEST_DIR/internal/web/" 2>/dev/null && echo "  ✓ internal/web/api.go (contains handleGetStickers, handleValidateSticker)" || echo "  ✗ api.go (not found)"
cp "$SOURCE_DIR/internal/web/router.go" "$DEST_DIR/internal/web/" 2>/dev/null && echo "  ✓ internal/web/router.go (contains sticker route)" || echo "  ✗ router.go (not found)"

# Copy supporting Go files (needed for compilation)
echo ""
echo "📚 Copying supporting Go files..."
for file in app.go config.go models.go middleware.go auth.go handlers.go templates.go language_middleware.go; do
    if [ -f "$SOURCE_DIR/internal/web/$file" ]; then
        cp "$SOURCE_DIR/internal/web/$file" "$DEST_DIR/internal/web/" 2>/dev/null && echo "  ✓ internal/web/$file" || echo "  ✗ $file (copy failed)"
    else
        echo "  ⚠ $file (not found - may be needed)"
    fi
done

# Copy main.go if it exists
if [ -f "$SOURCE_DIR/cmd/server/main.go" ]; then
    cp "$SOURCE_DIR/cmd/server/main.go" "$DEST_DIR/cmd/server/" 2>/dev/null && echo "  ✓ cmd/server/main.go" || echo "  ✗ main.go (copy failed)"
fi

# Copy sticker-related HTML files
echo ""
echo "🎨 Copying frontend template files..."
cp "$SOURCE_DIR/templates/writing_main.html" "$DEST_DIR/templates/" 2>/dev/null && echo "  ✓ templates/writing_main.html (contains sticker JavaScript)" || echo "  ✗ writing_main.html (not found)"
cp "$SOURCE_DIR/templates/includes/modals.html" "$DEST_DIR/templates/includes/" 2>/dev/null && echo "  ✓ templates/includes/modals.html (sticker picker modal)" || echo "  ✗ modals.html (not found)"
cp "$SOURCE_DIR/templates/includes/main_panel.html" "$DEST_DIR/templates/includes/" 2>/dev/null && echo "  ✓ templates/includes/main_panel.html (sticker button)" || echo "  ✗ main_panel.html (not found)"

# Copy sticker data files
echo ""
echo "🖼️  Copying sticker data files..."
cp "$SOURCE_DIR/stickers/stickers_map.json" "$DEST_DIR/stickers/" 2>/dev/null && echo "  ✓ stickers/stickers_map.json" || echo "  ✗ stickers_map.json (not found)"

# Copy all sticker GIF files
if [ -d "$SOURCE_DIR/stickers/GIFs" ]; then
    gif_count=$(find "$SOURCE_DIR/stickers/GIFs" -type f -name "*.gif" | wc -l)
    cp -r "$SOURCE_DIR/stickers/GIFs"/* "$DEST_DIR/stickers/GIFs/" 2>/dev/null
    if [ $? -eq 0 ]; then
        echo "  ✓ stickers/GIFs/ ($gif_count GIF files copied)"
    else
        echo "  ✗ Failed to copy GIF files"
    fi
else
    echo "  ⚠ stickers/GIFs/ directory not found"
fi

# Copy documentation files
echo ""
echo "📄 Copying documentation..."
cp "$SOURCE_DIR/STICKER_ELEMENTS_GUIDE.txt" "$DEST_DIR/" 2>/dev/null && echo "  ✓ STICKER_ELEMENTS_GUIDE.txt" || echo "  ⚠ STICKER_ELEMENTS_GUIDE.txt (not found)"
cp "$SOURCE_DIR/STICKER_FUNCTIONS_ANALYSIS.md" "$DEST_DIR/" 2>/dev/null && echo "  ✓ STICKER_FUNCTIONS_ANALYSIS.md" || echo "  ⚠ STICKER_FUNCTIONS_ANALYSIS.md (not found)"

# Create README for transfer
echo ""
echo "📝 Creating transfer README..."
cat > "$DEST_DIR/README_TRANSFER.md" << 'EOF'
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
EOF

echo "  ✓ README_TRANSFER.md created"

# Create file manifest
echo ""
echo "📋 Creating file manifest..."
cat > "$DEST_DIR/MANIFEST.txt" << EOF
Sticker Function Transfer Package - File Manifest
Generated: $(date)

FILES INCLUDED:
================

Backend Go Files:
- cmd/server/main.go
- internal/web/api.go (sticker handlers at lines 2974-2999, 3001-3043)
- internal/web/router.go (sticker route at lines 186-187)
- internal/web/app.go
- internal/web/config.go
- internal/web/models.go
- internal/web/middleware.go
- internal/web/auth.go
- internal/web/handlers.go
- internal/web/templates.go
- internal/web/language_middleware.go

Frontend Template Files:
- templates/writing_main.html (sticker JS: lines 8507-8602, display: 2413-2453)
- templates/includes/modals.html (sticker modal: lines 732-752, CSS: 675-704)
- templates/includes/main_panel.html (sticker button: lines 16-20)

Sticker Data:
- stickers/stickers_map.json
- stickers/GIFs/* (all GIF files)

Configuration:
- go.mod
- go.sum

Documentation:
- STICKER_ELEMENTS_GUIDE.txt
- STICKER_FUNCTIONS_ANALYSIS.md
- README_TRANSFER.md
- MANIFEST.txt

TOTAL FILES: $(find "$DEST_DIR" -type f | wc -l)
TOTAL SIZE: $(du -sh "$DEST_DIR" | cut -f1)
EOF

echo "  ✓ MANIFEST.txt created"

echo ""
echo "=========================================="
echo "✅ Isolation Complete!"
echo "=========================================="
echo ""
echo "📂 Destination folder:"
echo "   $DEST_DIR"
echo ""
echo "📊 Summary:"
echo "   - Go files: $(find "$DEST_DIR/internal/web" -name "*.go" 2>/dev/null | wc -l) files"
echo "   - Template files: $(find "$DEST_DIR/templates" -name "*.html" 2>/dev/null | wc -l) files"
echo "   - Sticker GIFs: $(find "$DEST_DIR/stickers/GIFs" -name "*.gif" 2>/dev/null | wc -l) files"
echo "   - Total files: $(find "$DEST_DIR" -type f | wc -l) files"
echo "   - Total size: $(du -sh "$DEST_DIR" | cut -f1)"
echo ""
echo "📖 Next steps:"
echo "   1. Review README_TRANSFER.md for transfer instructions"
echo "   2. Check MANIFEST.txt for complete file list"
echo "   3. Transfer the entire folder to your production server"
echo ""
