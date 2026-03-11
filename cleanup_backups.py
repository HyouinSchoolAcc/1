#!/usr/bin/env python3
"""
Cleanup script for backup files.
Keeps only the most recent backup for each unique file and deletes all older versions.
"""

import os
import re
from pathlib import Path
from collections import defaultdict

def extract_base_filename(backup_path):
    """
    Extract the base filename from a backup file.
    Example: user_5_Day3_dup_1_simplified.json.20260117T092919000Z.bak
             -> user_5_Day3_dup_1_simplified.json
    """
    filename = os.path.basename(backup_path)
    # Remove the timestamp and .bak extension
    # Pattern: .YYYYMMDDTHHMMSS[microseconds]Z.bak
    match = re.match(r'(.+?)\.(\d{8}T\d{6}.*?Z)\.bak$', filename)
    if match:
        return match.group(1), match.group(2)  # base_name, timestamp
    return None, None

def main():
    backups_dir = Path('/home/exx/Desktop/fine-tune/data_labler_UI_production/backups')
    
    if not backups_dir.exists():
        print(f"Backups directory not found: {backups_dir}")
        return
    
    print(f"Scanning backup files in: {backups_dir}")
    
    # Group backups by (directory, base_filename)
    backup_groups = defaultdict(list)
    
    # Find all .bak files
    for backup_file in backups_dir.rglob('*.bak'):
        base_name, timestamp = extract_base_filename(str(backup_file))
        if base_name and timestamp:
            # Group by directory and base filename
            key = (backup_file.parent, base_name)
            backup_groups[key].append((timestamp, backup_file))
    
    print(f"Found {len(backup_groups)} unique files with backups")
    
    total_backups = 0
    files_to_delete = []
    
    # For each group, keep only the most recent
    for (directory, base_name), backups in backup_groups.items():
        total_backups += len(backups)
        
        if len(backups) > 1:
            # Sort by timestamp (descending) - most recent first
            backups.sort(reverse=True, key=lambda x: x[0])
            
            # Keep the first (most recent), delete the rest
            most_recent = backups[0][1]
            to_delete = [backup[1] for backup in backups[1:]]
            
            files_to_delete.extend(to_delete)
            
            if len(to_delete) > 0:
                print(f"\n{base_name} ({directory.relative_to(backups_dir)})")
                print(f"  Keeping: {most_recent.name}")
                print(f"  Deleting: {len(to_delete)} older backup(s)")
    
    print(f"\n{'='*60}")
    print(f"Summary:")
    print(f"  Total backup files found: {total_backups}")
    print(f"  Unique files: {len(backup_groups)}")
    print(f"  Files to keep: {len(backup_groups)}")
    print(f"  Files to delete: {len(files_to_delete)}")
    print(f"{'='*60}")
    
    if files_to_delete:
        response = input(f"\nDelete {len(files_to_delete)} old backup files? (yes/no): ")
        if response.lower() in ['yes', 'y']:
            deleted_count = 0
            for file_path in files_to_delete:
                try:
                    file_path.unlink()
                    deleted_count += 1
                except Exception as e:
                    print(f"Error deleting {file_path}: {e}")
            
            print(f"\n✓ Successfully deleted {deleted_count} backup files")
            print(f"✓ Kept {len(backup_groups)} most recent backups")
        else:
            print("Cleanup cancelled.")
    else:
        print("\nNo duplicate backups to delete. All files are unique.")

if __name__ == '__main__':
    main()

