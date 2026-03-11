#!/bin/bash

# Script to copy all HTML and Go files to a new directory
# This creates a clean copy of only the web elements needed

DEST_DIR="web_elements_copy"
SOURCE_DIR="."

# Create destination directory
mkdir -p "$DEST_DIR"

# Copy Go module files
echo "Copying Go module files..."
cp go.mod "$DEST_DIR/" 2>/dev/null
cp go.sum "$DEST_DIR/" 2>/dev/null

# Copy all Go files preserving directory structure
echo "Copying Go files..."
find . -name "*.go" -type f ! -path "./.venv/*" ! -path "./.git/*" ! -path "./__pycache__/*" ! -path "./$DEST_DIR/*" | while read -r file; do
    # Get relative path
    rel_path="${file#./}"
    # Create directory structure in destination
    dest_file="$DEST_DIR/$rel_path"
    dest_dir=$(dirname "$dest_file")
    mkdir -p "$dest_dir"
    # Copy file
    cp "$file" "$dest_file"
    echo "  Copied: $rel_path"
done

# Copy all HTML files preserving directory structure
echo "Copying HTML files..."
find . -name "*.html" -type f ! -path "./.venv/*" ! -path "./.git/*" ! -path "./__pycache__/*" ! -path "./$DEST_DIR/*" | while read -r file; do
    # Get relative path
    rel_path="${file#./}"
    # Create directory structure in destination
    dest_file="$DEST_DIR/$rel_path"
    dest_dir=$(dirname "$dest_file")
    mkdir -p "$dest_dir"
    # Copy file
    cp "$file" "$dest_file"
    echo "  Copied: $rel_path"
done

# Also copy static directory if it exists (for CSS, JS, images)
if [ -d "static" ]; then
    echo "Copying static directory..."
    cp -r static "$DEST_DIR/" 2>/dev/null
fi

echo ""
echo "Done! All web elements copied to: $DEST_DIR"
echo ""
echo "Files copied:"
echo "  - All .go files (preserving directory structure)"
echo "  - All .html files (preserving directory structure)"
echo "  - go.mod and go.sum"
echo "  - static/ directory (if exists)"
echo ""
echo "You can now copy the entire '$DEST_DIR' directory to your new location."
