#!/usr/bin/env python3
"""
Quick test script to verify novel style renderer works.
Generates a simple test video without needing external files.
"""

import json
import os
import sys

def create_test_script():
    """Create a simple test JSON script"""
    script = {
        "person_a": "Alice",
        "person_b": "Bob",
        "messages": [
            {"sender": "b", "text": "Hey! How's the project going?", "delay_after": 2.0, "typing_duration": 0.8},
            {"sender": "a", "text": "Pretty good! Almost done with the backend.", "delay_after": 1.8, "typing_duration": 1.0},
            {"sender": "b", "text": "Nice! When can I test it?", "delay_after": 1.5, "typing_duration": 0.7},
            {"sender": "a", "text": "Should be ready by tomorrow afternoon.", "delay_after": 2.0, "typing_duration": 0.9},
            {"sender": "b", "text": "Perfect! I'll prepare some test cases.", "delay_after": 2.0, "typing_duration": 1.0},
            {"sender": "a", "text": "Great! Let's sync up at 2pm?", "delay_after": 1.5, "typing_duration": 0.8},
            {"sender": "b", "text": "Sounds good. See you then!", "delay_after": 2.0, "typing_duration": 0.7}
        ],
        "config": {
            "fps": 30,
            "width": 1080,
            "height": 1920,
            "hold_at_end": 2.5
        }
    }
    
    # Write test script
    with open('test_quick.json', 'w') as f:
        json.dump(script, f, indent=2)
    
    return 'test_quick.json'


def main():
    print("🎬 Quick Test - Video Novel Style")
    print("=" * 50)
    print()
    
    # Change to video_renderer directory
    if os.path.exists('video_renderer'):
        os.chdir('video_renderer')
    
    # Create test script
    print("📝 Creating test script...")
    script_path = create_test_script()
    print(f"   Created: {script_path}")
    print()
    
    # Test novel style
    print("🎨 Generating novel style video...")
    print("   (This may take 30-60 seconds)")
    print()
    
    import subprocess
    
    cmd = [
        sys.executable, 
        'generator.py', 
        script_path, 
        'test_novel_output.mp4',
        '--style', 'novel',
        '--max-messages', '3'
    ]
    
    try:
        result = subprocess.run(cmd, check=True, capture_output=False, text=True)
        print()
        print("✅ Success!")
        print()
        print("Generated video: video_renderer/test_novel_output.mp4")
        print()
        print("To test different configurations:")
        print()
        print("  # Traditional chat style")
        print(f"  python generator.py {script_path} output_chat.mp4")
        print()
        print("  # Novel style with 1 message (TikTok style)")
        print(f"  python generator.py {script_path} output_minimal.mp4 --style novel --max-messages 1")
        print()
        print("  # Novel style with 5 messages (more context)")
        print(f"  python generator.py {script_path} output_full.mp4 --style novel --max-messages 5")
        print()
        
        # Cleanup
        if os.path.exists(script_path):
            os.remove(script_path)
            
        return 0
        
    except subprocess.CalledProcessError as e:
        print()
        print("❌ Error generating video")
        print()
        print("Make sure you're in the video_chat_renderer directory")
        print("and have installed requirements:")
        print()
        print("  cd video_chat_renderer/video_renderer")
        print("  pip install -r requirements.txt")
        print()
        return 1
    
    except FileNotFoundError:
        print()
        print("❌ Error: generator.py not found")
        print()
        print("Make sure you're running this from the video_chat_renderer directory:")
        print()
        print("  cd video_chat_renderer")
        print("  python quick_test.py")
        print()
        return 1


if __name__ == '__main__':
    sys.exit(main())
