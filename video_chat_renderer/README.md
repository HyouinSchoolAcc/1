# Video Chat Renderer

Generate vertical chat-style videos from JSON scripts with **two rendering styles**: traditional chat and cinematic video novel!

![Python](https://img.shields.io/badge/python-3.8+-blue.svg)
![PIL](https://img.shields.io/badge/PIL-Pillow-green.svg)
![FFmpeg](https://img.shields.io/badge/FFmpeg-required-red.svg)

## ✨ Features

- 🎭 **Two Rendering Styles**
  - **Chat Style**: Traditional scrolling chat (WeChat/iMessage-like)
  - **Novel Style**: Cinematic full-screen messages with deadpan effect
- 🎤 **TTS Integration** (NEW!)
  - Text-to-Speech narration with perfect audio sync
  - Video doesn't advance until audio finishes
  - Multiple voice support
- 📱 Optimized for vertical video (1080x1920)
- 🎨 Customizable colors and layout
- 🖼️ Profile pictures with rounded corners
- ⌨️ Typing indicators
- 🎬 Deadpan instant appearance (no animations)
- 🌏 Full CJK (Chinese/Japanese/Korean) support
- 🖤 Screenshot-style black bars

## 🚀 Quick Start

### Installation

```bash
cd video_chat_renderer/video_renderer
pip install -r requirements.txt

# Make sure ffmpeg is installed
# Ubuntu/Debian: sudo apt-get install ffmpeg
# macOS: brew install ffmpeg
```

### Generate Your First Video

```bash
# Traditional chat style
python generator.py example_script.json output.mp4

# Video novel style (deadpan, one message at a time)
python generator.py example_novel.json output_novel.mp4 --style novel

# NEW: Video with TTS (Text-to-Speech)
python generator.py example_novel.json output_tts.mp4 --style novel --tts
```

**Note:** TTS requires Index-TTS server running. See [TTS_INTEGRATION.md](TTS_INTEGRATION.md) for setup.

### Quick Test

```bash
# Run quick test (generates test video automatically)
cd video_chat_renderer
python quick_test.py

# Or run comprehensive tests
./test_novel_style.sh
```

## 📖 Rendering Styles

### Chat Style (Traditional)

<table>
<tr>
<td width="50%">

**Features:**
- All messages visible
- Scrolling timeline
- Compact layout
- WeChat-like appearance

**Best for:**
- Full conversations
- Context-heavy dialogues
- Tutorial videos
- Documentation

</td>
<td width="50%">

```bash
python generator.py script.json \
  output.mp4 \
  --style chat
```

</td>
</tr>
</table>

### Novel Style (NEW!) 🎬

<table>
<tr>
<td width="50%">

**Features:**
- Prominent full-screen messages
- Only last N messages visible
- Flip-flop between speakers
- Larger fonts & avatars
- Cinematic animations

**Best for:**
- Dramatic dialogues
- TikTok/Reels/Shorts
- Visual novels
- Emotional exchanges
- TTS narration

</td>
<td width="50%">

```bash
python generator.py script.json \
  output.mp4 \
  --style novel \
  --max-messages 3
```

</td>
</tr>
</table>

## 📝 Usage

### Basic Commands

```bash
# Chat style (default)
python generator.py input.json output.mp4

# Novel style with default settings (3 messages per side)
python generator.py input.json output.mp4 --style novel

# Novel style with custom message limit
python generator.py input.json output.mp4 --style novel --max-messages 5

# TikTok/Reels style (minimal - 1 message per side)
python generator.py input.json output.mp4 --style novel --max-messages 1
```

### Command-Line Options

```
python generator.py <input.json> <output.mp4> [options]

Options:
  --style {chat,novel}   Rendering style (default: chat)
  --max-messages N       For novel: max messages per side (default: 3)
  
  --tts                  Enable TTS (Text-to-Speech)
  --tts-server URL       TTS server URL (default: http://localhost:6006)
  --voice-a NAME         Voice for person A (default: kurisu)
  --voice-b NAME         Voice for person B (default: jay_klee)
```

## 📄 JSON Script Format

Create a JSON file with your dialogue:

```json
{
  "person_a": "Alice",
  "person_b": "Bob",
  "messages": [
    {
      "sender": "a",
      "text": "Hey! How are you doing?",
      "delay_after": 2.0,
      "typing_duration": 1.0
    },
    {
      "sender": "b",
      "text": "I'm great, thanks for asking!",
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

### Field Descriptions

| Field | Required | Description |
|-------|----------|-------------|
| `person_a` | ✅ | Name of person A |
| `person_b` | ✅ | Name of person B |
| `messages` | ✅ | Array of message objects |
| `sender` | ✅ | `"a"` or `"b"` |
| `text` | ✅ | Message content |
| `delay_after` | ✅ | Seconds after message before next event |
| `typing_duration` | ❌ | Seconds to show typing indicator |
| `config` | ❌ | Video configuration |
| `fps` | ❌ | Frame rate (default: 30) |
| `width` | ❌ | Video width (default: 1080) |
| `height` | ❌ | Video height (default: 1920) |
| `hold_at_end` | ❌ | Hold final frame seconds (default: 2.0) |
| `person_a_image` | ❌ | Path to avatar A |
| `person_b_image` | ❌ | Path to avatar B |

## 🎯 Use Cases

### TikTok/Instagram Reels

```bash
# Minimal style for short-form content
python generator.py dialogue.json reels.mp4 --style novel --max-messages 1
```

### YouTube Shorts

```bash
# Bit more context
python generator.py dialogue.json shorts.mp4 --style novel --max-messages 2
```

### Visual Novel / Story Content

```bash
# Full novel experience
python generator.py story.json novel.mp4 --style novel --max-messages 3
```

### Tutorial / Documentation

```bash
# Full conversation history
python generator.py tutorial.json docs.mp4 --style chat
```

## 📊 Style Comparison

| Feature | Chat Style | Novel Style |
|---------|-----------|-------------|
| Message visibility | ✅ All | 🎯 Last N per side |
| Screen usage | Compact | Full-screen |
| Font size | 36px | 60px |
| Avatar size | 90px | 150px |
| Animation | Pop-in | Slide-in from sides |
| Best for | Full context | Dramatic moments |
| Typing indicator | ✅ Yes | 🚧 Coming soon |

## 🎨 Customization

### Colors & Layout

Edit `video_renderer/config.py`:

```python
# Colors
COLOR_PERSON_A_BUBBLE = (149, 236, 105)  # Green
COLOR_PERSON_B_BUBBLE = (255, 255, 255)  # White

# Sizes
MESSAGE_FONT_SIZE = 36
BUBBLE_MAX_WIDTH = 850
```

### Novel Style Settings

Edit `video_renderer/novel_renderer.py`:

```python
# Larger fonts
self.message_font_size = 80  # Default: 60

# Bigger avatars  
self.avatar_size = 200  # Default: 150
```

## 🔮 Coming Soon: TTS Integration

The novel style is designed with Text-to-Speech in mind:

```bash
# Planned usage
python generator.py script.json output.mp4 \
  --style novel \
  --tts \
  --voice-a "en-US-Neural2-A" \
  --voice-b "en-US-Neural2-C"
```

**Why novel style + TTS is perfect:**
- Prominent text matches audio narration
- Natural pacing for voice synthesis
- Less visual clutter = better focus
- Per-message timing for audio sync

## 📁 Project Structure

```
video_chat_renderer/
├── README.md                    # This file
├── USAGE_GUIDE.md              # Detailed usage guide
├── NOVEL_STYLE_README.md       # Novel style feature docs
├── quick_test.py               # Quick test script
├── test_novel_style.sh         # Comprehensive test suite
│
└── video_renderer/
    ├── generator.py            # Main entry point
    ├── parser.py              # JSON validation
    ├── timeline.py            # Frame timing
    ├── renderer.py            # Chat style renderer
    ├── novel_renderer.py      # Novel style renderer (NEW!)
    ├── encoder.py             # Video encoding
    ├── config.py              # Configuration
    ├── requirements.txt       # Dependencies
    ├── example_script.json    # Chat example
    └── example_novel.json     # Novel example (NEW!)
```

## 🛠️ Requirements

- Python 3.8+
- Pillow (PIL)
- FFmpeg (system installation)
- Noto CJK fonts (for Chinese/Japanese/Korean)

Install Python dependencies:
```bash
pip install -r video_renderer/requirements.txt
```

Install FFmpeg:
```bash
# Ubuntu/Debian
sudo apt-get install ffmpeg

# macOS
brew install ffmpeg

# Windows
# Download from https://ffmpeg.org/download.html
```

## 🐛 Troubleshooting

### "FFmpeg not found"

```bash
# Install ffmpeg
sudo apt-get install ffmpeg  # Linux
brew install ffmpeg          # macOS
```

### Chinese characters not displaying

```bash
# Install Noto CJK fonts
sudo apt-get install fonts-noto-cjk
```

### Avatar images not loading

- Check file paths in JSON
- Ensure images are readable (PNG/JPG)
- Images will auto-resize and crop

## 💡 Tips & Best Practices

### For Chat Style
- Use 1-2 second delays for natural pacing
- Add typing indicators for realism
- Keep messages concise

### For Novel Style
- Use 2-3 second delays (longer for emphasis)
- `--max-messages 2-3` works best for most content
- Perfect for emotional/dramatic moments
- Great for TikTok/Reels (use `--max-messages 1`)
- Longer messages get more emphasis

## 📚 Documentation

- **[USAGE_GUIDE.md](USAGE_GUIDE.md)** - Complete usage documentation
- **[NOVEL_STYLE_README.md](NOVEL_STYLE_README.md)** - Novel style feature deep-dive
- **[video_renderer/example_script.json](video_renderer/example_script.json)** - Chat style example
- **[video_renderer/example_novel.json](video_renderer/example_novel.json)** - Novel style example

## 🎬 Examples

See the `video_renderer/` directory for example scripts:

1. **example_script.json** - Basic chat conversation
2. **example_novel.json** - Novel style dialogue
3. **example_with_features.json** - All features demo

## 🚀 Next Steps

1. **Try it out:**
   ```bash
   cd video_chat_renderer
   python quick_test.py
   ```

2. **Read the guides:**
   - [USAGE_GUIDE.md](USAGE_GUIDE.md) for detailed usage
   - [NOVEL_STYLE_README.md](NOVEL_STYLE_README.md) for novel style specifics

3. **Create your content:**
   - Copy `example_novel.json`
   - Customize with your dialogue
   - Generate video!

4. **Experiment:**
   - Try different `--max-messages` values
   - Compare chat vs novel styles
   - Prepare for TTS integration

## 🤝 Contributing

Future enhancements:
- [ ] TTS integration (multiple engines)
- [ ] Background music support
- [ ] Custom color schemes per message
- [ ] Gradient backgrounds
- [ ] Emoji reactions
- [ ] Multi-speaker support (3+ people)
- [ ] Typing indicator for novel style
- [ ] Custom fonts per character
- [ ] Video backgrounds
- [ ] Subtitle/translation support

## 📝 License

[Add your license here]

## 🙏 Acknowledgments

- Built with Pillow (PIL) for image generation
- FFmpeg for video encoding
- Noto fonts for CJK support

---

**Ready to create cinematic chat videos? Start now! 🎬**

```bash
cd video_chat_renderer
python quick_test.py
```
