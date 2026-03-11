#!/bin/bash

# This script stops all running services (Go server and ngrok)

set -euo pipefail

echo "Stopping all services..."
echo ""

# Function to check if lsof is available
require_cmd() {
    if ! command -v "$1" >/dev/null 2>&1; then
        echo "WARNING: Command '$1' not found. Some cleanup may be incomplete."
        return 1
    fi
    return 0
}

STOPPED_SOMETHING=false

# Stop processes on port 5002
if require_cmd lsof; then
    echo "Checking for processes on port 5002..."
    if lsof -ti:5002 > /dev/null 2>&1; then
        echo "Stopping processes on port 5002..."
        lsof -ti:5002 | xargs kill -9 2>/dev/null
        echo "✓ Stopped processes on port 5002"
        STOPPED_SOMETHING=true
        sleep 1
    else
        echo "  No processes found on port 5002"
    fi
fi

# Stop processes on port 5003
if require_cmd lsof; then
    echo "Checking for processes on port 5003..."
    if lsof -ti:5003 > /dev/null 2>&1; then
        echo "Stopping processes on port 5003..."
        lsof -ti:5003 | xargs kill -9 2>/dev/null
        echo "✓ Stopped processes on port 5003"
        STOPPED_SOMETHING=true
        sleep 1
    else
        echo "  No processes found on port 5003"
    fi
fi

# Stop ngrok processes
echo "Checking for ngrok processes..."
if pgrep -x ngrok > /dev/null; then
    echo "Stopping ngrok..."
    pkill -9 ngrok
    echo "✓ Stopped ngrok"
    STOPPED_SOMETHING=true
    sleep 1
else
    echo "  No ngrok processes found"
fi

# Stop server_sql processes by name
echo "Checking for server_sql processes..."
if pgrep -f "server_sql" > /dev/null; then
    echo "Stopping server_sql..."
    pkill -9 -f "server_sql"
    echo "✓ Stopped server_sql"
    STOPPED_SOMETHING=true
    sleep 1
else
    echo "  No server_sql processes found"
fi

# Stop any Go processes running from this directory (backup check)
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
echo "Checking for Go processes in $SCRIPT_DIR..."
if pgrep -f "$SCRIPT_DIR.*go" > /dev/null; then
    echo "Stopping Go processes in this directory..."
    pkill -9 -f "$SCRIPT_DIR.*go"
    echo "✓ Stopped Go processes"
    STOPPED_SOMETHING=true
    sleep 1
else
    echo "  No Go processes found in this directory"
fi

echo ""
echo "=========================================="
if [ "$STOPPED_SOMETHING" = true ]; then
    echo "✓ All services stopped successfully"
else
    echo "✓ No services were running"
fi
echo "=========================================="
echo ""
echo "You can now safely start services again with:"
echo "  - Linux/Mac: ./run_go3.sh (from parent directory)"
echo "  - Or from this directory: ../run_go3.sh"
echo ""

