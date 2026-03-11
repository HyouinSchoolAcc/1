# TTS Integration Guide

## Overview

The video chat renderer now supports **Text-to-Speech (TTS)** integration using the Index-TTS-VLLM server. The video will automatically sync with the audio duration - **each message stays on screen until the audio is fully finished playing**.

## ✨ Features

✅ **Automatic Audio Sync** - Video timing syncs with TTS audio duration
✅ **No Manual Timing** - `delay_after` is automatically calculated from audio length
✅ **Multiple Voices** - Assign different TTS voices to each character
✅ **Seamless Integration** - Audio composited into final video file
✅ **Quality Audio** - 24kHz sample rate, AAC encoding

## 🚀 Quick Start

### 1. Start the TTS Server

First, make sure the Index-TTS server is running:

```bash
cd /home/exx/Desktop/fine-tune/index-tts-vllm
python api_server.py
```

The server will start on `http://localhost:6006` by default.

### 2. Generate Video with TTS

```bash
cd video_chat_renderer/video_renderer

# Basic usage with TTS
python generator.py script.json output.mp4 --style novel --tts

# Custom voices
python generator.py script.json output.mp4 \
  --style novel \
  --tts \
  --voice-a kurisu \
  --voice-b jay_klee
```

## 📝 Command-Line Options

### TTS-Specific Options

| Option | Default | Description |
|--------|---------|-------------|
| `--tts` | False | Enable TTS generation |
| `--tts-server` | http://localhost:6006 | TTS API server URL |
| `--voice-a` | kurisu | Voice/character for person A |
| `--voice-b` | jay_klee | Voice/character for person B |

### Combined Example

```bash
python generator.py dialogue.json output_tts.mp4 \
  --style novel \
  --tts \
  --tts-server http://localhost:6006 \
  --voice-a kurisu \
  --voice-b jay_klee
```

## 🎭 Available Voices

Default voices configured in `index-tts-vllm/assets/speaker.json`:

| Voice Name | Description | Reference Audio |
|------------|-------------|-----------------|
| `kurisu` | Female voice | kurisu_13s.wav |
| `jay_klee` | Mixed voice | jay_promptvn.wav + klee audio |

To add custom voices, edit `assets/speaker.json` in the index-tts-vllm directory.

## 📖 How It Works

### Workflow

```
1. Load script.json
   ↓
2. [TTS Enabled] Generate audio for each message
   ├─ Call TTS API for each message
   ├─ Get audio duration
   └─ Update message.delay_after = audio_duration + 0.3s
   ↓
3. Generate video frames
   ├─ Timeline calculated with new durations
   └─ Each message visible for full audio duration
   ↓
4. Encode video (no audio)
   ↓
5. Create single audio track
   ├─ Place each audio clip at correct timestamp
   └─ Fill gaps with silence
   ↓
6. Composite audio + video
   └─ FFmpeg combines into final MP4
```

### Timing Synchronization

**Without TTS:**
```json
{
  "text": "Hello, how are you?",
  "delay_after": 2.0  ← Manual timing
}
```
Video shows message for exactly 2.0 seconds.

**With TTS:**
```json
{
  "text": "Hello, how are you?",
  "delay_after": 2.0  ← Ignored when TTS enabled
}
```
1. TTS generates audio: 1.8 seconds
2. System updates: `delay_after = 1.8 + 0.3 = 2.1`
3. Video shows message for 2.1 seconds
4. Audio plays for full 1.8 seconds

**Result:** Message doesn't advance until audio is complete! ✨

## 📋 Example Script

```json
{
  "person_a": "Alice",
  "person_b": "Bob",
  "messages": [
    {
      "sender": "b",
      "text": "Hey! How's your day going?",
      "delay_after": 2.0,
      "typing_duration": 1.0
    },
    {
      "sender": "a",
      "text": "Pretty good! Just finished my presentation.",
      "delay_after": 2.5,
      "typing_duration": 1.2
    },
    {
      "sender": "b",
      "text": "Nice! How did it go?",
      "delay_after": 1.5,
      "typing_duration": 0.8
    }
  ],
  "config": {
    "fps": 30,
    "width": 1080,
    "height": 1920
  }
}
```

Generate with TTS:
```bash
python generator.py example.json output_with_voice.mp4 \
  --style novel \
  --tts \
  --voice-a kurisu \
  --voice-b jay_klee
```

## 🎯 Best Practices

### 1. Message Length

```
✅ Good: "Hey! How's it going?"
✅ Good: "I just finished the project and sent it over."
⚠️ Long: "So I was thinking that maybe we could possibly..."
```

Shorter messages (5-15 words) work best for natural pacing.

### 2. Voice Assignment

```bash
# Assign appropriate voices to characters
--voice-a kurisu      # Female character → kurisu voice
--voice-b jay_klee    # Male character → jay_klee voice
```

### 3. Novel Style + TTS = Perfect Match

Novel style (one message at a time) works perfectly with TTS:
- Clear focus on current narration
- No visual clutter
- Natural pacing
- Professional result

```bash
# Recommended combination
python generator.py script.json output.mp4 \
  --style novel \  # One message at a time
  --tts            # Auto-sync with audio
```

## 🔧 Troubleshooting

### TTS Server Not Responding

```
❌ Error: TTS server is not responding!
```

**Solution:**
```bash
# Start the TTS server
cd /home/exx/Desktop/fine-tune/index-tts-vllm
python api_server.py

# Wait for: ">> TTS model initialized successfully"
```

### Voice Not Found

```
❌ TTS request failed with status 400
```

**Solution:**
Check available voices in `index-tts-vllm/assets/speaker.json` and use exact names.

### Audio Out of Sync

If audio doesn't match video timing:

1. Check FPS consistency (should be 30)
2. Verify all audio files generated successfully
3. Look for warnings in console output

### Dependencies Missing

```
❌ ModuleNotFoundError: No module named 'soundfile'
```

**Solution:**
```bash
cd video_renderer
pip install -r requirements.txt
```

## 🎬 Complete Example

### 1. Prepare Script

`conversation.json`:
```json
{
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
      "text": "I need help with my project.",
      "delay_after": 1.5,
      "typing_duration": 0.7
    },
    {
      "sender": "b",
      "text": "I'd be happy to assist! What kind of project are you working on?",
      "delay_after": 2.5,
      "typing_duration": 1.0
    }
  ]
}
```

### 2. Start TTS Server

Terminal 1:
```bash
cd index-tts-vllm
python api_server.py
```

Wait for initialization message.

### 3. Generate Video

Terminal 2:
```bash
cd video_chat_renderer/video_renderer

python generator.py conversation.json output_final.mp4 \
  --style novel \
  --tts \
  --voice-a kurisu \
  --voice-b jay_klee
```

### 4. Output

```
Loading script: conversation.json
Video config: 1080x1920 @ 30 FPS
Messages: 3
Style: novel

============================================================
🎤 TTS Generation Enabled
============================================================
TTS Server: http://localhost:6006
Voice A (User): kurisu
Voice B (Assistant): jay_klee

Checking TTS server health...
✓ TTS server is healthy

Generating TTS audio for all messages...
------------------------------------------------------------
Generating audio for message 1/3: Hello! How can I help you today?...
  ✓ Generated audio: 1.92s (original delay: 2.00s)
Generating audio for message 2/3: I need help with my project....
  ✓ Generated audio: 1.45s (original delay: 1.50s)
Generating audio for message 3/3: I'd be happy to assist! What kind of...
  ✓ Generated audio: 2.58s (original delay: 2.50s)
------------------------------------------------------------
✓ Generated 3 audio clips

Starting video generation: output_final.mp4

Using novel style renderer (showing ONE message at a time)
Generating 213 frames (7.10 seconds)...
[Progress bar...]
Encoding with ffmpeg (may take 10-60 seconds)...

============================================================
🎵 Compositing Audio with Video
============================================================
Creating audio track...
✓ Audio track created: /tmp/audio_track_123.wav
Adding audio to video...
✓ Audio added successfully

✓ Video with TTS audio created successfully: output_final.mp4
```

## 🔄 Without TTS (Comparison)

```bash
# Generate without TTS (manual timing)
python generator.py conversation.json output_no_tts.mp4 --style novel

# Messages use delay_after from JSON
# No audio in final video
```

## 📊 Performance

| Stage | Typical Duration |
|-------|-----------------|
| TTS generation (per message) | 0.5-2s |
| Video rendering | 10-30s |
| Audio composition | 5-10s |
| **Total** | **~1-2 minutes for 10 messages** |

## 🚀 Advanced Usage

### Custom TTS Server

```bash
# Use different TTS server
python generator.py script.json output.mp4 \
  --tts \
  --tts-server http://192.168.1.100:6006
```

### Batch Processing

```bash
#!/bin/bash
# Generate multiple videos with TTS

for script in scripts/*.json; do
  output="output/$(basename $script .json).mp4"
  python generator.py "$script" "$output" \
    --style novel \
    --tts \
    --voice-a kurisu \
    --voice-b jay_klee
done
```

## 🎓 Tips

1. **Start TTS server first** - Always ensure server is running before generating
2. **Use novel style** - Best combination with TTS for one-message-at-a-time narration
3. **Short messages** - Work best with TTS (5-15 words)
4. **Test voices** - Try different voice combinations for your characters
5. **Check console** - Watch for TTS generation progress and any errors

## 📚 Related Documentation

- [USAGE_GUIDE.md](USAGE_GUIDE.md) - General usage without TTS
- [NOVEL_STYLE_UPDATE.md](NOVEL_STYLE_UPDATE.md) - Novel style details
- [DEADPAN_UPDATE.md](DEADPAN_UPDATE.md) - Deadpan effect details

---

**Now you can create professional narrated chat videos with perfect audio sync! 🎙️🎬**
