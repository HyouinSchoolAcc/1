#!/bin/bash

# Backup script for JSON files before SQL migration

TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="data_backup_$TIMESTAMP"

echo "Creating backup at $BACKUP_DIR..."

# Create backup directory
mkdir -p "$BACKUP_DIR"

# Copy data directory
cp -r data/* "$BACKUP_DIR/"

echo "✓ Backup completed successfully!"
echo "Backup location: $(pwd)/$BACKUP_DIR"
echo ""
echo "If you need to restore, run:"
echo "  rm -rf data/*"
echo "  cp -r $BACKUP_DIR/* data/"

