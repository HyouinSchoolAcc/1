#!/bin/bash

# Run the data labeler UI server in debug mode on port 5004

# Get the directory where the script is located
APP_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

# Set port to 5004
export SERVE_PORT="5004"

# Change to the app directory
cd "$APP_DIR" || exit

echo "=========================================="
echo "Starting Data Labeler UI Server"
echo "=========================================="
echo "Port: $SERVE_PORT"
echo "Mode: DEBUG (Dev Reload Enabled)"
echo "Directory: $APP_DIR"
echo "=========================================="
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed or not in PATH"
    exit 1
fi

# Check if dependencies are installed
if [ ! -d "vendor" ] && [ ! -f "go.sum" ]; then
    echo "Installing Go dependencies..."
    go mod download
fi

# Run the server using go run (for development/debug mode)
# This allows hot-reloading of templates since EnableDevReload is true
echo "Starting server..."
go run cmd/server/main.go

