#!/usr/bin/env python3
"""
Export Storyboard
Exports all character storylines, schedules, and daily interactions into a comprehensive
JSON document for the storyboard editor.
"""

import json
import os
import glob
from pathlib import Path
from datetime import datetime
from typing import Dict, Any, List, Optional
import argparse
import re


class StoryboardExporter:
    def __init__(self, base_dir: Optional[str] = None):
        if base_dir is None:
            script_dir = Path(__file__).parent
            self.base_dir = script_dir
        else:
            self.base_dir = Path(base_dir)
        
        self.presets_dir = self.base_dir / "presets"
        self.data_dir = self.base_dir / "data"
        self.output_file = self.base_dir / "static" / "storyboard.json"
    
    def load_character_profiles(self) -> Dict[str, Any]:
        """Load character profiles from data/character_profiles.json"""
        profiles_file = self.data_dir / "character_profiles.json"
        if not profiles_file.exists():
            return {}
        with open(profiles_file, 'r', encoding='utf-8') as f:
            return json.load(f)
    
    def get_preset_folders(self, include_cn: bool = True) -> List[Path]:
        """Get all preset folders"""
        if not self.presets_dir.exists():
            return []
        folders = []
        for d in sorted(self.presets_dir.iterdir()):
            if d.is_dir() and d.name.startswith("presets_"):
                if include_cn or not d.name.endswith("_CN"):
                    folders.append(d)
        return folders
    
    def load_day_files(self, preset_dir: Path) -> Dict[int, Dict[int, List[Dict[str, Any]]]]:
        """
        Load all day files for a preset, organized by user and day.
        Returns: {user_id: {day: [duplicate_files_data...]}}
        """
        pattern = str(preset_dir / "user_*_Day*_dup_*_simplified.json")
        files = glob.glob(pattern)
        
        result = {}
        for file_path in files:
            filename = os.path.basename(file_path)
            match = re.match(r'user_(\d+)_Day(\d+)_dup_(\d+)_simplified\.json', filename)
            if not match:
                continue
            
            user_id = int(match.group(1))
            day_num = int(match.group(2))
            dup_num = int(match.group(3))
            
            try:
                with open(file_path, 'r', encoding='utf-8') as f:
                    data = json.load(f)
                    data['_filename'] = filename
                    data['_filepath'] = file_path
                    data['_dup'] = dup_num
                    
                    if user_id not in result:
                        result[user_id] = {}
                    if day_num not in result[user_id]:
                        result[user_id][day_num] = []
                    result[user_id][day_num].append(data)
            except Exception as e:
                print(f"Error reading {file_path}: {e}")
        
        return result
    
    def load_user_info(self, preset_dir: Path) -> List[Dict[str, Any]]:
        """Load user info file for the preset"""
        preset_name = preset_dir.name.replace("presets_", "").replace("_CN", "")
        patterns = [
            preset_dir / f"new_user_info_{preset_name}.json",
            preset_dir / f"new_user_info_{preset_name.replace('_', '-')}.json",
            preset_dir / f"new_user_info_{preset_name}_cn.json",
        ]
        
        for pattern in patterns:
            if pattern.exists():
                try:
                    with open(pattern, 'r', encoding='utf-8') as f:
                        return json.load(f)
                except:
                    pass
        return []
    
    def extract_schedule(self, schedule: Dict[str, Any]) -> Dict[str, Any]:
        """Extract schedule with both CN and EN versions"""
        if not schedule:
            return {}
        
        return {
            "day": schedule.get("day", 0),
            "title": schedule.get("title", ""),
            "title_en": schedule.get("title_en", schedule.get("title", "")),
            "morning": schedule.get("morning", schedule.get("早晨", "")),
            "morning_en": schedule.get("morning_en", schedule.get("morning", "")),
            "noon": schedule.get("noon", schedule.get("中午", "")),
            "noon_en": schedule.get("noon_en", schedule.get("noon", "")),
            "afternoon": schedule.get("afternoon", schedule.get("下午", "")),
            "afternoon_en": schedule.get("afternoon_en", schedule.get("afternoon", "")),
            "evening": schedule.get("evening", schedule.get("晚上", "")),
            "evening_en": schedule.get("evening_en", schedule.get("evening", "")),
            "night": schedule.get("night", schedule.get("夜晚", "")),
            "night_en": schedule.get("night_en", schedule.get("night", "")),
        }
    
    def build_character_data(self, char_id: str, profile: Dict[str, Any], 
                              en_preset_dir: Path, cn_preset_dir: Optional[Path]) -> Dict[str, Any]:
        """Build complete character data structure"""
        
        # Basic profile info
        char_data = {
            "id": char_id,
            "name": profile.get("name", char_id),
            "name_en": profile.get("english_name", profile.get("name", char_id)),
            "tagline": profile.get("tagline", ""),
            "tagline_en": profile.get("tagline_en", profile.get("tagline", "")),
            "description": profile.get("description", ""),
            "description_en": profile.get("description_en", profile.get("description", "")),
            "values": profile.get("values", ""),
            "values_en": profile.get("values_en", profile.get("values", "")),
            "experiences": profile.get("experiences", ""),
            "experiences_en": profile.get("experiences_en", profile.get("experiences", "")),
            "judgements": profile.get("judgements", ""),
            "judgements_en": profile.get("judgements_en", profile.get("judgements", "")),
            "abilities": profile.get("abilities", ""),
            "abilities_en": profile.get("abilities_en", profile.get("abilities", "")),
            "relationships": [],
            "official_schedules": [],
            "story": profile.get("story"),
            "users": []
        }
        
        # Relationships
        for rel in profile.get("relationships", []):
            char_data["relationships"].append({
                "name": rel.get("name", ""),
                "name_en": rel.get("name_en", rel.get("name", "")),
                "style": rel.get("style", ""),
                "style_en": rel.get("style_en", rel.get("style", "")),
                "description": rel.get("description", ""),
                "description_en": rel.get("description_en", rel.get("description", "")),
            })
        
        # Official schedules from profile
        for sched in profile.get("schedules", []):
            char_data["official_schedules"].append(self.extract_schedule(sched))
        
        # Load user scenarios from EN preset
        user_info = self.load_user_info(en_preset_dir)
        user_names = {u.get('id', i): {
            "name": u.get('name', f"User {u.get('id', i)}"),
            "name_en": u.get('english_name', u.get('name', f"User {u.get('id', i)}"))
        } for i, u in enumerate(user_info)}
        
        days_by_user = self.load_day_files(en_preset_dir)
        
        # Also load CN data if available
        cn_days_by_user = {}
        if cn_preset_dir and cn_preset_dir.exists():
            cn_days_by_user = self.load_day_files(cn_preset_dir)
        
        # Build user scenarios
        for user_id in sorted(days_by_user.keys()):
            user_data = {
                "id": user_id,
                "name": user_names.get(user_id, {}).get("name", f"User {user_id}"),
                "name_en": user_names.get(user_id, {}).get("name_en", f"User {user_id}"),
                "days": []
            }
            
            for day_num in sorted(days_by_user[user_id].keys()):
                day_files = days_by_user[user_id][day_num]
                data = day_files[0] if day_files else {}
                
                # Try to get CN version
                cn_data = {}
                if user_id in cn_days_by_user and day_num in cn_days_by_user[user_id]:
                    cn_files = cn_days_by_user[user_id][day_num]
                    cn_data = cn_files[0] if cn_files else {}
                
                char_sched = data.get('character_schedule', {})
                user_sched = data.get('user_schedule', {})
                cn_char_sched = cn_data.get('character_schedule', {})
                cn_user_sched = cn_data.get('user_schedule', {})
                
                day_data = {
                    "day": day_num,
                    "filename": data.get('_filename', ''),
                    "filepath": data.get('_filepath', ''),
                    "duplicates": len(day_files),
                    "category": data.get('category', 'pending'),
                    "relationship": data.get('relationship', ''),
                    "history": data.get('history', ''),
                    "starting_intimacy": data.get('starting_intimacy_level', 0),
                    "intimacy": data.get('intimacy_level', 0),
                    "dialogue_count": len(data.get('dialogue', [])),
                    "character_schedule": {
                        "day": day_num,
                        "morning": char_sched.get('morning', '') if isinstance(char_sched, dict) else '',
                        "morning_cn": cn_char_sched.get('morning', cn_char_sched.get('早晨', '')) if isinstance(cn_char_sched, dict) else '',
                        "noon": char_sched.get('noon', '') if isinstance(char_sched, dict) else '',
                        "noon_cn": cn_char_sched.get('noon', cn_char_sched.get('中午', '')) if isinstance(cn_char_sched, dict) else '',
                        "afternoon": char_sched.get('afternoon', '') if isinstance(char_sched, dict) else '',
                        "afternoon_cn": cn_char_sched.get('afternoon', cn_char_sched.get('下午', '')) if isinstance(cn_char_sched, dict) else '',
                        "evening": char_sched.get('evening', '') if isinstance(char_sched, dict) else '',
                        "evening_cn": cn_char_sched.get('evening', cn_char_sched.get('晚上', '')) if isinstance(cn_char_sched, dict) else '',
                        "night": char_sched.get('night', '') if isinstance(char_sched, dict) else '',
                        "night_cn": cn_char_sched.get('night', cn_char_sched.get('夜晚', '')) if isinstance(cn_char_sched, dict) else '',
                    },
                    "user_schedule": {
                        "day": day_num,
                        "morning": user_sched.get('morning', '') if isinstance(user_sched, dict) else '',
                        "morning_cn": cn_user_sched.get('morning', cn_user_sched.get('早晨', '')) if isinstance(cn_user_sched, dict) else '',
                        "noon": user_sched.get('noon', '') if isinstance(user_sched, dict) else '',
                        "noon_cn": cn_user_sched.get('noon', cn_user_sched.get('中午', '')) if isinstance(cn_user_sched, dict) else '',
                        "afternoon": user_sched.get('afternoon', '') if isinstance(user_sched, dict) else '',
                        "afternoon_cn": cn_user_sched.get('afternoon', cn_user_sched.get('下午', '')) if isinstance(cn_user_sched, dict) else '',
                        "evening": user_sched.get('evening', '') if isinstance(user_sched, dict) else '',
                        "evening_cn": cn_user_sched.get('evening', cn_user_sched.get('晚上', '')) if isinstance(cn_user_sched, dict) else '',
                        "night": user_sched.get('night', '') if isinstance(user_sched, dict) else '',
                        "night_cn": cn_user_sched.get('night', cn_user_sched.get('夜晚', '')) if isinstance(cn_user_sched, dict) else '',
                    },
                    "dialogue_preview": []
                }
                
                # Add dialogue preview (first and last few messages)
                dialogue = data.get('dialogue', [])
                if dialogue:
                    preview = []
                    for msg in dialogue[:3]:
                        preview.append({
                            "role": msg.get('role', '?'),
                            "content": msg.get('content', '')[:100],
                            "timestamp": msg.get('timestamp', '')
                        })
                    if len(dialogue) > 6:
                        preview.append({"role": "...", "content": f"({len(dialogue) - 6} more messages)", "timestamp": ""})
                    for msg in dialogue[-3:]:
                        preview.append({
                            "role": msg.get('role', '?'),
                            "content": msg.get('content', '')[:100],
                            "timestamp": msg.get('timestamp', '')
                        })
                    day_data["dialogue_preview"] = preview
                
                user_data["days"].append(day_data)
            
            char_data["users"].append(user_data)
        
        return char_data
    
    def export(self) -> Dict[str, Any]:
        """Generate the complete storyboard data"""
        profiles = self.load_character_profiles()
        
        # Map character IDs to preset folder names
        char_to_preset = {
            "kurisu": "presets_kurisu",

            "linlu": "presets_lin_lu",
            "lin_lu": "presets_lin_lu",
            "newcharacter_1": "presets_newcharacter_1",
        }
        
        storyboard = {
            "generated": datetime.now().isoformat(),
            "characters": []
        }
        
        # Process each character
        processed = set()
        for char_id, profile in profiles.items():
            if char_id in processed:
                continue
            processed.add(char_id)
            
            preset_name = char_to_preset.get(char_id, f"presets_{char_id}")
            en_preset_dir = self.presets_dir / preset_name
            cn_preset_dir = self.presets_dir / f"{preset_name}_CN"
            
            if not en_preset_dir.exists():
                # Try alternate naming
                alt_names = [f"presets_{char_id.replace('_', '')}", f"presets_{char_id}"]
                for alt in alt_names:
                    if (self.presets_dir / alt).exists():
                        en_preset_dir = self.presets_dir / alt
                        cn_preset_dir = self.presets_dir / f"{alt}_CN"
                        break
            
            if en_preset_dir.exists():
                char_data = self.build_character_data(
                    char_id, profile, en_preset_dir, 
                    cn_preset_dir if cn_preset_dir.exists() else None
                )
                storyboard["characters"].append(char_data)
                print(f"Processed {char_id}: {len(char_data['users'])} users, {sum(len(u['days']) for u in char_data['users'])} days")
        
        return storyboard
    
    def save(self):
        """Export and save to file"""
        data = self.export()
        
        # Ensure output directory exists
        self.output_file.parent.mkdir(parents=True, exist_ok=True)
        
        with open(self.output_file, 'w', encoding='utf-8') as f:
            json.dump(data, f, ensure_ascii=False, indent=2)
        
        print(f"Storyboard exported to: {self.output_file}")
        return self.output_file


def main():
    parser = argparse.ArgumentParser(description="Export storyboard to JSON")
    parser.add_argument("--base-dir", default=None, help="Base directory of the project")
    parser.add_argument("--output", default=None, help="Output file path")
    
    args = parser.parse_args()
    
    exporter = StoryboardExporter(args.base_dir)
    if args.output:
        exporter.output_file = Path(args.output)
    
    exporter.save()


if __name__ == "__main__":
    main()
