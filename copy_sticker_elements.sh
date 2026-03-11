#!/bin/bash

# Script to copy all emote/sticker-related web elements
# This creates a clean copy of only the sticker functionality

DEST_DIR="sticker_elements_copy"
SOURCE_DIR="."

# Create destination directory structure
mkdir -p "$DEST_DIR"
mkdir -p "$DEST_DIR/templates/includes"
mkdir -p "$DEST_DIR/internal/web"
mkdir -p "$DEST_DIR/stickers/GIFs"

# Copy Go module files (needed for compilation)
echo "Copying Go module files..."
cp go.mod "$DEST_DIR/" 2>/dev/null
cp go.sum "$DEST_DIR/" 2>/dev/null

# Copy sticker-related Go files
echo "Copying Go files..."
# api.go contains sticker API handlers
cp internal/web/api.go "$DEST_DIR/internal/web/" 2>/dev/null
# router.go contains sticker route
cp internal/web/router.go "$DEST_DIR/internal/web/" 2>/dev/null

# Copy sticker-related HTML files
echo "Copying HTML files..."
# Main writing page with sticker JavaScript functionality
cp templates/writing_main.html "$DEST_DIR/templates/" 2>/dev/null
# Modal with sticker picker UI
cp templates/includes/modals.html "$DEST_DIR/templates/includes/" 2>/dev/null
# Main panel with sticker button
cp templates/includes/main_panel.html "$DEST_DIR/templates/includes/" 2>/dev/null

# Copy sticker data files
echo "Copying sticker data..."
cp stickers/stickers_map.json "$DEST_DIR/stickers/" 2>/dev/null
# Copy all sticker GIF files
if [ -d "stickers/GIFs" ]; then
    cp -r stickers/GIFs/* "$DEST_DIR/stickers/GIFs/" 2>/dev/null
    echo "  Copied sticker GIF files"
fi

# Copy necessary supporting files
echo "Copying supporting files..."
# Check if there are other files needed
if [ -f "internal/web/app.go" ]; then
    cp internal/web/app.go "$DEST_DIR/internal/web/" 2>/dev/null
    echo "  Copied app.go"
fi
if [ -f "internal/web/config.go" ]; then
    cp internal/web/config.go "$DEST_DIR/internal/web/" 2>/dev/null
    echo "  Copied config.go"
fi
if [ -f "internal/web/models.go" ]; then
    cp internal/web/models.go "$DEST_DIR/internal/web/" 2>/dev/null
    echo "  Copied models.go"
fi
if [ -f "internal/web/middleware.go" ]; then
    cp internal/web/middleware.go "$DEST_DIR/internal/web/" 2>/dev/null
    echo "  Copied middleware.go"
fi
if [ -f "internal/web/auth.go" ]; then
    cp internal/web/auth.go "$DEST_DIR/internal/web/" 2>/dev/null
    echo "  Copied auth.go"
fi
if [ -f "internal/web/handlers.go" ]; then
    cp internal/web/handlers.go "$DEST_DIR/internal/web/" 2>/dev/null
    echo "  Copied handlers.go"
fi
if [ -f "internal/web/templates.go" ]; then
    cp internal/web/templates.go "$DEST_DIR/internal/web/" 2>/dev/null
    echo "  Copied templates.go"
fi
if [ -f "internal/web/language_middleware.go" ]; then
    cp internal/web/language_middleware.go "$DEST_DIR/internal/web/" 2>/dev/null
    echo "  Copied language_middleware.go"
fi
if [ -f "cmd/server/main.go" ]; then
    mkdir -p "$DEST_DIR/cmd/server"
    cp cmd/server/main.go "$DEST_DIR/cmd/server/" 2>/dev/null
    echo "  Copied main.go"
fi

echo ""
echo "Done! All sticker/emote elements copied to: $DEST_DIR"
echo ""
echo "Files copied:"
echo "  ✓ Go files:"
echo "    - internal/web/api.go (contains handleGetStickers, handleValidateSticker)"
echo "    - internal/web/router.go (contains sticker route)"
echo "    - Supporting Go files (app.go, config.go, models.go, etc.)"
echo ""
echo "  ✓ HTML files:"
echo "    - templates/writing_main.html (sticker JavaScript functions)"
echo "    - templates/includes/modals.html (sticker picker modal)"
echo "    - templates/includes/main_panel.html (sticker button)"
echo ""
echo "  ✓ Sticker data:"
echo "    - stickers/stickers_map.json"
echo "    - stickers/GIFs/* (all sticker images)"
echo ""
echo "  ✓ Go module files:"
echo "    - go.mod"
echo "    - go.sum"
echo ""
echo "NOTE: The copied api.go and router.go files contain other functionality too."
echo "      See STICKER_ELEMENTS_GUIDE.txt for details on which parts are sticker-related."
