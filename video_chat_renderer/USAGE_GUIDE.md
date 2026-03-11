# Video Chat Renderer - Usage Guide

A Python-based tool to generate vertical chat-style videos from JSON scripts. Now supports two rendering styles!

## Features

- ✅ **Two Rendering Styles**
  - **Chat Style** (Traditional): Scrolling chat with all messages visible
  - **Novel Style** (New): Prominent full-screen messages that flip-flop between speakers
- ✅ Profile pictures with rounded corners
- ✅ Typing indicators
- ✅ Smooth animations
- ✅ Customizable colors and layout
- ✅ Support for Chinese characters (CJK fonts)
- ✅ Prepared for TTS integration (coming soon)

## Installation

```bash
cd video_renderer
pip install -r requirements.txt
```

## Quick Start

### Traditional Chat Style (Default)

```bash
python generator.py example_script.json output_chat.mp4
```

### Video Novel Style

```bash
python generator.py example_script.json output_novel.mp4 --style novel
```

### Novel Style with More Messages

```bash
# Show last 5 messages from each side
python generator.py example_script.json output_novel.mp4 --style novel --max-messages 5
```

## Command Line Options

```
python generator.py <input.json> <output.mp4> [options]

Required Arguments:
  input.json          Path to input JSON script file
  output.mp4          Path to output video file

Optional Arguments:
  --style {chat,novel}
                      Rendering style (default: chat)
                      - 'chat': Traditional scrolling chat with all messages
                      - 'novel': Video novel style with prominent messages
  
  --max-messages N    For novel style only (default: 3)
                      Maximum number of recent messages to show from each side
```

## Rendering Styles Explained

### Chat Style (Traditional)

The traditional chat interface style:
- All messages scroll up from the bottom
- Messages stay visible throughout the video
- Similar to WeChat/iMessage appearance
- Good for: Full conversations, context-heavy dialogues

**Best for:** When you want viewers to see the full conversation history

### Novel Style (Video Novel)

A prominent, cinematic presentation style:
- Each message takes up significant screen space
- Only shows the last N messages from each person
- Messages alternate/flip-flop between left and right alignment
- Larger fonts and more prominent avatars
- Smooth slide-in animations from sides

**Best for:** 
- Dramatic dialogues
- Visual novels
- TikTok/Instagram Reels style content
- When you want focus on current conversation, not history

## JSON Script Format

```json
{
  "person_a": "Alice",
  "person_b": "Bob",
  "messages": [
    {
      "sender": "a",
      "text": "Hello, how are you?",
      "delay_after": 2.0,
      "typing_duration": 1.0
    },
    {
      "sender": "b",
      "text": "I'm doing great, thanks!",
      "delay_after": 1.5,
      "typing_duration": 0.8
    }
  ],
  "config": {
    "fps": 30,
    "width": 1080,
    "height": 1920,
    "hold_at_end": 2.0,
    "person_a_image": "/path/to/avatar_a.png",
    "person_b_image": "/path/to/avatar_b.png"
  }
}
```

### Message Fields

- `sender`: Either `"a"` or `"b"` (person_a or person_b)
- `text`: The message content (supports Chinese/emoji)
- `delay_after`: Seconds to wait after message appears before next event
- `typing_duration`: (Optional) Seconds to show typing indicator before message

### Config Fields (Optional)

- `fps`: Video frame rate (default: 30)
- `width`: Video width in pixels (default: 1080)
- `height`: Video height in pixels (default: 1920)
- `hold_at_end`: Seconds to hold final frame (default: 2.0)
- `person_a_image`: Path to person A's profile picture
- `person_b_image`: Path to person B's profile picture

## Examples

### Example 1: Basic Chat Style

```bash
python generator.py example_script.json output_chat.mp4
```

Creates a traditional scrolling chat video.

### Example 2: Novel Style with Default Settings

```bash
python generator.py example_novel.json output_novel.mp4 --style novel
```

Creates a video novel with last 3 messages from each side.

### Example 3: Novel Style with More Context

```bash
python generator.py example_novel.json output_novel_full.mp4 --style novel --max-messages 6
```

Shows more message history in novel style (last 6 from each side).

### Example 4: Short-Form Content (Novel Style)

Perfect for TikTok/Reels:

```bash
# Show only the most recent message from each side
python generator.py short_script.json tiktok_video.mp4 --style novel --max-messages 1
```

## Style Comparison

| Feature | Chat Style | Novel Style |
|---------|-----------|-------------|
| Message visibility | All messages | Last N from each side |
| Screen usage | Efficient, compact | Prominent, full-screen |
| Font size | Standard (36px) | Large (60px) |
| Avatar size | 90px | 150px |
| Best for | Full conversations | Dramatic moments |
| Animation | Pop-in from position | Slide-in from sides |
| Viewing history | Complete | Recent only |

## TTS Integration (Coming Soon)

The novel style is designed with TTS (Text-to-Speech) in mind:

- Prominent text display works well with voice narration
- Per-message timing allows audio sync
- `delay_after` will eventually sync with audio duration
- Each speaker can have a distinct voice

**Planned features:**
- `--tts` flag to enable text-to-speech
- Voice selection for each character
- Automatic timing based on audio duration
- Lip-sync support (stretch goal)

## Tips & Best Practices

### For Chat Style

- Use natural conversation pacing (1-2 seconds per message)
- Add typing indicators for realism
- Keep messages concise for better readability

### For Novel Style

- **Less is more**: Use `--max-messages 2` or `3` for best results
- Longer `delay_after` times (2-3 seconds) work better
- Great for emotional or dramatic exchanges
- Perfect for short-form vertical video platforms
- Works excellently with longer messages that need emphasis

### General Tips

- Use high-quality avatar images (at least 200x200px)
- Test with different `hold_at_end` values for your platform
- Chinese text works great with both styles
- Consider your target platform's video specs

## Troubleshooting

### Fonts not displaying correctly

The tool tries multiple font paths. If Chinese characters don't display:
```bash
# Install Noto CJK fonts (Ubuntu/Debian)
sudo apt-get install fonts-noto-cjk

# Or specify a custom font by modifying config.py
```

### Video encoding fails

Make sure ffmpeg is installed:
```bash
# Ubuntu/Debian
sudo apt-get install ffmpeg

# macOS
brew install ffmpeg
```

### Avatar images not loading

- Check file paths in your JSON config
- Ensure images are readable (PNG, JPG supported)
- Images will be automatically resized and cropped

## File Structure

```
video_renderer/
├── generator.py          # Main entry point
├── parser.py            # JSON validation
├── timeline.py          # Frame timing calculation
├── renderer.py          # Chat style renderer
├── novel_renderer.py    # Novel style renderer (NEW!)
├── encoder.py           # Video encoding
├── config.py            # Style configuration
├── example_script.json  # Chat style example
├── example_novel.json   # Novel style example (NEW!)
└── requirements.txt     # Dependencies
```

## Next Steps

1. **Try both styles** with your own scripts
2. **Experiment** with `--max-messages` values for novel style
3. **Prepare for TTS** by timing your messages appropriately
4. **Create content** for different platforms (TikTok, Instagram, YouTube Shorts)

## Contributing

Future enhancements planned:
- [ ] TTS integration with multiple voice engines
- [ ] Background music support
- [ ] Custom color schemes
- [ ] Gradient backgrounds
- [ ] Emoji reaction overlays
- [ ] Multi-language support

---

**Questions? Issues?** Open an issue on GitHub or reach out!
