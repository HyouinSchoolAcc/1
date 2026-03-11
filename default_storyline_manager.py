#!/usr/bin/env python3
"""
Default Storyline Manager
This script manages default storylines for all characters by:
1. Extracting storylines from user_0 files 
2. Storing them in a central default storylines file
3. Applying default storylines to all character files
"""

import json
import os
import glob
import shutil
from pathlib import Path
from typing import Dict, Any, List, Optional
import argparse
from datetime import datetime

class DefaultStorylineManager:
    def __init__(self, base_dir: Optional[str] = None):
        if base_dir is None:
            # Use absolute path relative to this script's directory
            script_dir = Path(__file__).parent
            self.base_dir = script_dir / "presets"
        else:
            self.base_dir = Path(base_dir)
        self.default_storylines_file = self.base_dir / "default_storylines.json"
        self.backup_dir = None
        
    def extract_default_storylines(self) -> Dict[str, Dict[int, Dict[str, Any]]]:
        """
        Extract default storylines from user_0_Day*_dup_1 files across all character presets.
        Returns a structure like:
        {
            "presets_kurisu": {
                1: {"day": 1, "morning": "...", "noon": "...", ...},
                2: {"day": 2, "morning": "...", "noon": "...", ...},
                ...
            },
            ...
        }
        """
        default_storylines = {}
        
        # Find all preset directories
        preset_dirs = [d for d in self.base_dir.iterdir() if d.is_dir() and d.name.startswith("presets_")]
        
        for preset_dir in preset_dirs:
            preset_name = preset_dir.name
            default_storylines[preset_name] = {}
            
            # Find all user_0_Day*_dup_1_simplified.json files first, then try dup_0 if none found
            pattern_1 = str(preset_dir / "user_0_Day*_dup_1_simplified.json")
            user_0_files = glob.glob(pattern_1)
            
            # If no dup_1 files found, try dup_0 files (some presets use dup_0 as canonical)
            if not user_0_files:
                pattern_0 = str(preset_dir / "user_0_Day*_dup_0_simplified.json")
                user_0_files = glob.glob(pattern_0)
            
            print(f"Found {len(user_0_files)} user_0 files in {preset_name}")
            
            for file_path in user_0_files:
                try:
                    with open(file_path, 'r', encoding='utf-8') as f:
                        data = json.load(f)
                    
                    if 'character_schedule' in data:
                        schedule = data['character_schedule']
                        if 'day' in schedule:
                            day_num = str(schedule['day'])  # Convert to string for consistent keys
                            default_storylines[preset_name][day_num] = schedule
                            print(f"  Extracted Day {day_num} from {os.path.basename(file_path)}")
                        else:
                            print(f"  Warning: No 'day' field in {file_path}")
                    else:
                        print(f"  Warning: No 'character_schedule' in {file_path}")
                        
                except Exception as e:
                    print(f"  Error reading {file_path}: {e}")
        
        return default_storylines
    
    def save_default_storylines(self, storylines: Dict[str, Dict[int, Dict[str, Any]]]):
        """Save default storylines to the central file."""
        with open(self.default_storylines_file, 'w', encoding='utf-8') as f:
            json.dump(storylines, f, indent=2, ensure_ascii=False)
        print(f"Default storylines saved to {self.default_storylines_file}")
    
    def load_default_storylines(self) -> Dict[str, Dict[int, Dict[str, Any]]]:
        """Load default storylines from the central file."""
        if not self.default_storylines_file.exists():
            print(f"Default storylines file not found: {self.default_storylines_file}")
            return {}
        
        with open(self.default_storylines_file, 'r', encoding='utf-8') as f:
            return json.load(f)
    
    def create_backup_directory(self) -> Path:
        """Create a timestamped backup directory."""
        timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
        backup_dir = self.base_dir / f"backup_{timestamp}"
        backup_dir.mkdir(exist_ok=True)
        self.backup_dir = backup_dir
        print(f"Created backup directory: {backup_dir}")
        return backup_dir
    
    def backup_file(self, file_path: str) -> bool:
        """Create a backup of a file, maintaining directory structure."""
        if not self.backup_dir:
            return False
        
        try:
            # Get relative path from base_dir
            rel_path = Path(file_path).relative_to(self.base_dir)
            backup_file_path = self.backup_dir / rel_path
            
            # Create parent directories if they don't exist
            backup_file_path.parent.mkdir(parents=True, exist_ok=True)
            
            # Copy the file
            shutil.copy2(file_path, backup_file_path)
            return True
        except Exception as e:
            print(f"  Warning: Failed to backup {file_path}: {e}")
            return False
    
    def apply_default_storylines_to_all(self, dry_run: bool = True, create_backup: bool = True):
        """
        Apply default storylines to ALL character files (except user_0 files).
        
        Args:
            dry_run: If True, only show what would be changed without making changes
            create_backup: If True and not dry_run, create backups before making changes
        """
        default_storylines = self.load_default_storylines()
        
        if not default_storylines:
            print("No default storylines found. Run extract_and_save_defaults() first.")
            return
        
        # Create backup directory if we're making real changes
        if not dry_run and create_backup:
            self.create_backup_directory()
        
        changes_made = 0
        files_processed = 0
        
        # Process all preset directories
        for preset_dir in self.base_dir.iterdir():
            if not preset_dir.is_dir() or not preset_dir.name.startswith("presets_"):
                continue
                
            preset_name = preset_dir.name
            if preset_name not in default_storylines:
                print(f"Warning: No default storylines found for {preset_name}")
                continue
            
            print(f"\nProcessing {preset_name}...")
            
            # Find all character files except user_0 files
            pattern = str(preset_dir / "*_simplified.json")
            all_files = glob.glob(pattern)
            
            # Filter out user_0 files
            character_files = [f for f in all_files if not os.path.basename(f).startswith("user_0_")]
            
            for file_path in character_files:
                files_processed += 1
                try:
                    # Extract day number from filename
                    filename = os.path.basename(file_path)
                    # Parse filename like: user_20_Day3_dup_1_simplified.json
                    parts = filename.split('_')
                    day_part = None
                    for part in parts:
                        if part.startswith('Day'):
                            day_part = part
                            break
                    
                    if not day_part:
                        print(f"  Warning: Could not extract day number from {filename}")
                        continue
                    
                    day_num = day_part[3:]  # Remove 'Day' prefix, keep as string
                    
                    if day_num not in default_storylines[preset_name]:
                        print(f"  Warning: No default storyline for Day {day_num} in {preset_name}")
                        continue
                    
                    # Load the character file
                    with open(file_path, 'r', encoding='utf-8') as f:
                        character_data = json.load(f)
                    
                    # Check if character_schedule needs updating
                    default_schedule = default_storylines[preset_name][day_num]
                    current_schedule = character_data.get('character_schedule', {})
                    
                    if current_schedule != default_schedule:
                        if not dry_run:
                            # Create backup first if backup directory exists
                            if create_backup and self.backup_dir:
                                self.backup_file(file_path)
                            
                            character_data['character_schedule'] = default_schedule.copy()
                            with open(file_path, 'w', encoding='utf-8') as f:
                                json.dump(character_data, f, indent=2, ensure_ascii=False)
                        
                        changes_made += 1
                        action = "Would update" if dry_run else "Updated"
                        print(f"  {action}: {filename} (Day {day_num})")
                    
                except Exception as e:
                    print(f"  Error processing {file_path}: {e}")
        
        print(f"\nSummary:")
        print(f"Files processed: {files_processed}")
        print(f"Changes {'planned' if dry_run else 'made'}: {changes_made}")
        
        if not dry_run and changes_made > 0 and create_backup and self.backup_dir:
            print(f"Backups created in: {self.backup_dir}")
        
        if dry_run and changes_made > 0:
            print(f"\nTo apply these changes, run without --dry-run flag")
    
    def extract_and_save_defaults(self):
        """Extract default storylines from user_0 files and save them."""
        print("Extracting default storylines from user_0 files...")
        storylines = self.extract_default_storylines()
        
        if storylines:
            self.save_default_storylines(storylines)
            
            # Print summary
            print(f"\nExtracted storylines summary:")
            for preset_name, days in storylines.items():
                print(f"  {preset_name}: Days {sorted(days.keys())}")
        else:
            print("No storylines extracted.")
    
    def show_default_storylines(self):
        """Display the current default storylines."""
        storylines = self.load_default_storylines()
        
        if not storylines:
            print("No default storylines found.")
            return
        
        print("Current Default Storylines:")
        print("=" * 50)
        
        for preset_name, days in storylines.items():
            print(f"\n{preset_name.upper()}:")
            for day_num in sorted(days.keys(), key=int):
                schedule = days[day_num]
                print(f"  Day {day_num}:")
                # Handle both English and Chinese keys
                period_mappings = [
                    ('morning', '早晨'),
                    ('noon', '中午'), 
                    ('afternoon', '下午'),
                    ('evening', '晚上'),
                    ('night', '夜晚')
                ]
                
                for eng_period, cn_period in period_mappings:
                    # Try English key first, then Chinese key
                    content = schedule.get(eng_period) or schedule.get(cn_period, "")
                    display_period = eng_period
                    if cn_period in schedule:
                        display_period = f"{eng_period} ({cn_period})"
                    
                    if content:
                        print(f"    {display_period}: {content}")
                    else:
                        print(f"    {display_period}: (empty)")

def main():
    parser = argparse.ArgumentParser(description="Manage default storylines for characters")
    parser.add_argument("--base-dir", default="presets", help="Base directory containing preset folders")
    
    subparsers = parser.add_subparsers(dest='command', help='Available commands')
    
    # Extract command
    extract_parser = subparsers.add_parser('extract', help='Extract default storylines from user_0 files')
    
    # Apply command
    apply_parser = subparsers.add_parser('apply', help='Apply default storylines to all character files')
    apply_parser.add_argument('--dry-run', action='store_true', help='Show what would be changed without making changes')
    apply_parser.add_argument('--no-backup', action='store_true', help='Skip creating backups (not recommended)')
    
    # Show command
    show_parser = subparsers.add_parser('show', help='Show current default storylines')
    
    args = parser.parse_args()
    
    manager = DefaultStorylineManager(args.base_dir)
    
    if args.command == 'extract':
        manager.extract_and_save_defaults()
    elif args.command == 'apply':
        dry_run = args.dry_run if hasattr(args, 'dry_run') else True
        create_backup = not (args.no_backup if hasattr(args, 'no_backup') else False)
        manager.apply_default_storylines_to_all(dry_run=dry_run, create_backup=create_backup)
    elif args.command == 'show':
        manager.show_default_storylines()
    else:
        parser.print_help()

if __name__ == "__main__":
    main() 