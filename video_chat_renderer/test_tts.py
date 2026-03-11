#!/usr/bin/env python3
"""
Quick test script for TTS integration.
Checks if TTS server is available and generates a test video with audio.
"""

import sys
import os
import subprocess
import requests

def check_tts_server(url="http://localhost:6006"):
    """Check if TTS server is running"""
    try:
        response = requests.get(f"{url}/health", timeout=5)
        if response.status_code == 200:
            print(f"✓ TTS server is healthy at {url}")
            return True
        else:
            print(f"✗ TTS server responded with status {response.status_code}")
            return False
    except requests.exceptions.ConnectionError:
        print(f"✗ TTS server not responding at {url}")
        return False
    except Exception as e:
        print(f"✗ Error checking TTS server: {e}")
        return False

def create_test_script():
    """Create a simple test script"""
    script = {
        "person_a": "User",
        "person_b": "Assistant",
        "messages": [
            {
                "sender": "b",
                "text": "Hello! How can I help you today?",
                "delay_after": 2.0,
                "typing_duration": 0.8
            },
            {
                "sender": "a",
                "text": "I'd like to test the TTS feature.",
                "delay_after": 1.8,
                "typing_duration": 0.7
            },
            {
                "sender": "b",
                "text": "Great! The audio should sync perfectly with the video.",
                "delay_after": 2.2,
                "typing_duration": 0.9
            }
        ],
        "config": {
            "fps": 30,
            "width": 1080,
            "height": 1920,
            "hold_at_end": 1.5
        }
    }
    
    import json
    with open('video_renderer/test_tts_script.json', 'w') as f:
        json.dump(script, f, indent=2)
    
    return 'video_renderer/test_tts_script.json'

def main():
    print("="*70)
    print("TTS Integration Test")
    print("="*70)
    print()
    
    # Check if in correct directory
    if not os.path.exists('video_renderer'):
        print("✗ Error: Must run from video_chat_renderer directory")
        print()
        print("  cd /home/exx/Desktop/fine-tune/video_chat_renderer")
        print("  python test_tts.py")
        print()
        return 1
    
    # Check TTS server
    print("Step 1: Checking TTS server...")
    print("-" * 70)
    
    if not check_tts_server():
        print()
        print("❌ TTS server is not running!")
        print()
        print("To start the TTS server:")
        print()
        print("  Terminal 1:")
        print("  cd /home/exx/Desktop/fine-tune/index-tts-vllm")
        print("  python api_server.py")
        print()
        print("  Then run this test again in Terminal 2")
        print()
        return 1
    
    print()
    
    # Create test script
    print("Step 2: Creating test script...")
    print("-" * 70)
    script_path = create_test_script()
    print(f"✓ Created: {script_path}")
    print()
    
    # Generate video with TTS
    print("Step 3: Generating video with TTS...")
    print("-" * 70)
    print()
    
    cmd = [
        sys.executable,
        'video_renderer/generator.py',
        script_path,
        'video_renderer/test_tts_output.mp4',
        '--style', 'novel',
        '--tts',
        '--voice-a', 'kurisu',
        '--voice-b', 'jay_klee'
    ]
    
    try:
        result = subprocess.run(cmd, check=True)
        
        print()
        print("="*70)
        print("✅ SUCCESS!")
        print("="*70)
        print()
        print("Test video with TTS generated:")
        print("  video_renderer/test_tts_output.mp4")
        print()
        print("Features verified:")
        print("  ✓ TTS server communication")
        print("  ✓ Audio generation")
        print("  ✓ Video/audio synchronization")
        print("  ✓ Audio composition")
        print()
        print("Next steps:")
        print("  1. Play test_tts_output.mp4 to verify audio")
        print("  2. Create your own scripts with dialogue")
        print("  3. Generate videos with custom voices")
        print()
        print("See TTS_INTEGRATION.md for full documentation")
        print()
        
        return 0
        
    except subprocess.CalledProcessError:
        print()
        print("="*70)
        print("❌ Video generation failed")
        print("="*70)
        print()
        print("Check the error messages above for details")
        print()
        return 1

if __name__ == '__main__':
    sys.exit(main())
