#!/usr/bin/env python3
"""
Generate WeChat-style chat videos from user data files.
Converts user dialogue JSON files to video renderer format and generates videos.
"""

import json
import os
import sys
import subprocess
import math
import random


# Paths
BASE_DIR = os.path.dirname(os.path.abspath(__file__))
VIDEO_RENDERER_DIR = os.path.join(BASE_DIR, "video_renderer")
PRESETS_DIR = "/home/exx/Desktop/fine-tune/data_labler_UI_production/presets/presets_kurisu_CN"

# Avatar paths
USER_AVATAR = "/home/exx/Desktop/fine-tune/data_labler_UI_production/static/user/user_avatar.png"
KURISU_AVATAR = "/home/exx/Desktop/fine-tune/data_labler_UI_production/static/kurisu/kurisu_avatar.png"

# Stickers directory
STICKERS_DIR = "/home/exx/Desktop/fine-tune/data_labler_UI_production/stickers"

# Files to process: (source_file, output_video_name, header_text)
FILES_TO_PROCESS = [
    ("user_32_Day1_dup_13_simplified.json", "user32_v13_day1.mp4", "stranger"),
    ("user_32_Day2_dup_13_simplified.json", "user32_v13_day2.mp4", "屑爆王18"),
    ("user_32_Day1_dup_2_simplified.json",  "user32_v2_day1.mp4",  "stranger"),
    ("user_32_Day2_dup_2_simplified.json",  "user32_v2_day2.mp4",  "屑爆王18"),
]


def convert_dialogue_to_script(data: dict, header_text: str = "chat") -> dict:
    """
    Convert a user data JSON file to the video renderer script format.
    
    Args:
        data: Parsed user data dictionary with 'dialogue', 'user_name', 'character_name', etc.
        header_text: Text to display in the header bar
    
    Returns:
        Video renderer script dictionary
    """
    user_name = data.get("user_name", "User")
    character_name = data.get("character_name", "Kurisu")
    
    # Inner thought annotations keyed by dialogue index (e.g. {"0": {...}, "2": {...}})
    inner_thoughts = data.get("inner_thought_annotations", {})
    
    messages = []
    
    # --- Simulated clock for timestamps ---
    # Start at 20:00 (8 PM) and advance realistically
    sim_minutes = 20 * 60  # 20:00 in minutes from midnight
    last_timestamp_min = -999  # force first timestamp
    TIMESTAMP_INTERVAL = 5  # show a new timestamp every 5+ simulated minutes
    
    random.seed(42)  # deterministic but varied gaps
    
    sticker_count = 0
    thought_count = 0
    
    for i, msg in enumerate(data.get("dialogue", [])):
        content = msg.get("content", "").strip()
        role = msg.get("role", "")
        
        # Skip empty messages
        if not content:
            continue
        
        # Determine sender
        if role == character_name or role == "Kurisu":
            sender = "b"
        else:
            sender = "a"
        
        # --- Advance simulated clock ---
        sim_gap = random.uniform(0.5, 2.0)
        if random.random() < 0.08:
            sim_gap += random.uniform(3, 10)
        sim_minutes += sim_gap
        
        # Handle sticker messages
        if content.startswith("[[sticker:") and content.endswith("]]"):
            sticker_rel = content[len("[[sticker:"):-2]  # e.g. "GIFs/file.gif"
            sticker_path = os.path.join(STICKERS_DIR, sticker_rel)
            if not os.path.exists(sticker_path):
                print(f"    ⚠ Sticker not found, skipping: {sticker_path}")
                continue
            
            msg_dict = {
                "sender": sender,
                "type": "sticker",
                "sticker_path": sticker_path,
                "text": "",
                "typing_duration": 0.5,
                "delay_after": 1.5
            }
            sticker_count += 1
        else:
            # Strip trailing newlines
            content = content.rstrip("\n").strip()
            if not content:
                continue
            
            # Calculate typing duration based on text length
            # Shorter messages: -0.2s, longer messages: +0.2s compared to before
            char_count = len(content)
            typing_duration = max(0.2, min(2.0, char_count * 0.05))
            
            # Calculate delay_after based on text length (reading time)
            delay_after = max(1.0, min(3.0, char_count * 0.06 + 0.8))
            
            msg_dict = {
                "sender": sender,
                "text": content,
                "typing_duration": round(typing_duration, 2),
                "delay_after": round(delay_after, 2)
            }
        
        # Add timestamp if enough simulated time has passed
        if sim_minutes - last_timestamp_min >= TIMESTAMP_INTERVAL:
            total_mins = int(sim_minutes) % (24 * 60)
            hours = total_mins // 60
            mins = total_mins % 60
            msg_dict["timestamp"] = f"{hours:02d}:{mins:02d}"
            last_timestamp_min = sim_minutes
        
        # --- Inject inner thought BEFORE this message ---
        # The thought appears first so the viewer sees the character's
        # reasoning before the spoken response.
        # Source 1: per-message "reasoning_chain" field (on character messages)
        # Source 2: top-level "inner_thought_annotations" dict (keyed by index)
        thought_text = ""
        
        # Check reasoning_chain first (directly on the message, character only)
        if sender == "b":
            reasoning = msg.get("reasoning_chain", "").strip()
            if reasoning:
                thought_text = reasoning
        
        # Override with inner_thought_annotations if present (higher priority)
        thought_key = str(i)
        if thought_key in inner_thoughts:
            thought_data = inner_thoughts[thought_key]
            annotated = (thought_data.get("correct_thought") or
                         thought_data.get("actual_thought", "")).strip()
            if annotated:
                thought_text = annotated
        
        if thought_text:
            thought_msg = {
                "sender": "b",       # always the character's inner voice
                "type": "thought",
                "text": thought_text,
                "typing_duration": 0, # thoughts appear instantly, no typing dots
                "delay_after": 1.2    # brief pause so the viewer can read it
            }
            messages.append(thought_msg)
            thought_count += 1
        
        messages.append(msg_dict)
    
    # --- Append character schedule card as the final message ---
    character_schedule = data.get("character_schedule")
    if character_schedule and isinstance(character_schedule, dict):
        # Check if there are any non-empty time period entries
        periods = ["morning", "noon", "afternoon", "evening", "night"]
        has_content = any((character_schedule.get(p) or "").strip() for p in periods)
        if has_content:
            schedule_msg = {
                "sender": "b",
                "type": "schedule",
                "text": "",
                "schedule": character_schedule,
                "character_name": character_name,
                "typing_duration": 0,
                "delay_after": 3.0  # hold the schedule card on screen
            }
            messages.append(schedule_msg)
    
    # Build the video renderer script
    script = {
        "person_a": user_name,
        "person_b": character_name,
        "messages": messages,
        "config": {
            "fps": 30,
            "width": 1080,
            "height": 1920,
            "hold_at_end": 2.0,
            "person_a_image": USER_AVATAR,
            "person_b_image": KURISU_AVATAR,
            "header_text": header_text
        }
    }
    
    return script


def generate_video(script_path: str, output_path: str) -> bool:
    """
    Run the video generator for a given script.
    
    Args:
        script_path: Path to the video renderer JSON script
        output_path: Path for the output video
    
    Returns:
        True if successful, False otherwise
    """
    cmd = [
        sys.executable,
        os.path.join(VIDEO_RENDERER_DIR, "generator.py"),
        script_path,
        output_path,
        "--style", "chat",  # WeChat/chat style
        "--workers", "16",  # 16 parallel rendering workers
    ]
    
    try:
        result = subprocess.run(cmd, check=True, text=True, capture_output=False)
        return True
    except subprocess.CalledProcessError as e:
        print(f"  ❌ Error generating video (exit code {e.returncode})")
        if e.stderr:
            print(f"     stderr: {e.stderr[:500]}")
        return False


def main():
    print("=" * 60)
    print("🎬 User32 Video Generator")
    print("   Generating WeChat-style chat videos with typing indicators")
    print("=" * 60)
    print()
    
    # Output directory for generated videos
    output_dir = os.path.join(VIDEO_RENDERER_DIR, "user32_videos")
    os.makedirs(output_dir, exist_ok=True)
    
    # Temp directory for converted scripts
    temp_dir = os.path.join(output_dir, "scripts")
    os.makedirs(temp_dir, exist_ok=True)
    
    results = []
    
    for source_file, output_name, header_text in FILES_TO_PROCESS:
        source_path = os.path.join(PRESETS_DIR, source_file)
        output_path = os.path.join(output_dir, output_name)
        script_name = output_name.replace(".mp4", ".json")
        script_path = os.path.join(temp_dir, script_name)
        
        print("-" * 60)
        print(f"📄 Processing: {source_file}")
        print(f"   Output:     {output_name}")
        print(f"   Header:     {header_text}")
        print()
        
        # Check source file exists
        if not os.path.exists(source_path):
            print(f"  ❌ Source file not found: {source_path}")
            results.append((source_file, False, "Source file not found"))
            continue
        
        # Load source data
        with open(source_path, 'r', encoding='utf-8') as f:
            data = json.load(f)
        
        # Convert to video renderer format
        script = convert_dialogue_to_script(data, header_text=header_text)
        msg_count = len(script["messages"])
        sticker_count = sum(1 for m in script["messages"] if m.get("type") == "sticker")
        thought_count = sum(1 for m in script["messages"] if m.get("type") == "thought")
        schedule_count = sum(1 for m in script["messages"] if m.get("type") == "schedule")
        print(f"  ✓ Converted {msg_count} messages ({sticker_count} stickers, {thought_count} inner thoughts, {schedule_count} schedule card{'s' if schedule_count != 1 else ''})")
        
        if msg_count == 0:
            print(f"  ❌ No messages to render!")
            results.append((source_file, False, "No messages"))
            continue
        
        # Save converted script
        with open(script_path, 'w', encoding='utf-8') as f:
            json.dump(script, f, indent=2, ensure_ascii=False)
        print(f"  ✓ Script saved: {script_path}")
        
        # Generate video
        print(f"  🎥 Generating video...")
        success = generate_video(script_path, output_path)
        
        if success:
            # Get file size
            size_mb = os.path.getsize(output_path) / (1024 * 1024)
            print(f"  ✅ Video generated: {output_path} ({size_mb:.1f} MB)")
            results.append((source_file, True, f"{size_mb:.1f} MB"))
        else:
            results.append((source_file, False, "Generation failed"))
        
        print()
    
    # Summary
    print("=" * 60)
    print("📊 Summary")
    print("=" * 60)
    for source, success, info in results:
        status = "✅" if success else "❌"
        print(f"  {status} {source} → {info}")
    
    successful = sum(1 for _, s, _ in results if s)
    print()
    print(f"  {successful}/{len(results)} videos generated successfully")
    print(f"  Output directory: {output_dir}")
    print()
    
    return 0 if successful == len(results) else 1


if __name__ == "__main__":
    sys.exit(main())
