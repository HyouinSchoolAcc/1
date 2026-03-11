# Video Novel Style - Feature Overview

## What's New?

I've added a **Video Novel Style** rendering mode to the chat video renderer. This new style displays messages in a prominent, full-screen format that alternates between speakers - perfect for dramatic dialogues and short-form content!

## Quick Comparison

### Before (Chat Style)
```bash
python generator.py script.json output.mp4
```
- Traditional scrolling chat interface
- All messages visible
- Compact, WeChat-like appearance

### After (Novel Style) - NEW!
```bash
python generator.py script.json output.mp4 --style novel
```
- Prominent, full-screen message display
- Only shows last N messages from each side
- Flip-flops between speakers
- Larger fonts and avatars
- Cinematic slide-in animations

## Key Features

### 1. **Prominent Message Display**
- Each message takes up significant screen space
- Much larger fonts (60px vs 36px)
- Larger avatars (150px vs 90px)
- Perfect for mobile viewing

### 2. **Flip-Flop Layout**
- Person A messages: Right-aligned with slide-in from right
- Person B messages: Left-aligned with slide-in from left
- Creates a visual "back and forth" feel

### 3. **Configurable Message Limit**
```bash
# Show last 3 messages from each side (default)
--style novel --max-messages 3

# TikTok/Reels style (minimal)
--style novel --max-messages 1

# More context
--style novel --max-messages 5
```

### 4. **Smooth Animations**
- Slide-in from appropriate side
- Fade-in for text and avatars
- Staggered animation for name, avatar, and text
- Cubic easing for smooth motion

## File Changes

### New Files
- `video_renderer/novel_renderer.py` - Novel style renderer implementation
- `video_renderer/example_novel.json` - Example script for novel style
- `USAGE_GUIDE.md` - Comprehensive usage documentation
- `test_novel_style.sh` - Quick test script for all variants

### Modified Files
- `video_renderer/generator.py` - Added `--style` and `--max-messages` arguments

### Unchanged Files
- `renderer.py` - Original chat renderer (still works perfectly!)
- `parser.py`, `timeline.py`, `encoder.py`, `config.py` - No changes needed

## Usage Examples

### Example 1: Quick Test
```bash
cd video_chat_renderer
./test_novel_style.sh
```
This will generate 4 videos showing different configurations!

### Example 2: Create Novel Style Video
```bash
cd video_renderer
python generator.py example_novel.json my_novel_video.mp4 --style novel
```

### Example 3: TikTok/Reels Style (Minimal)
```bash
python generator.py script.json tiktok.mp4 --style novel --max-messages 1
```

### Example 4: Keep Traditional Style
```bash
# Everything still works exactly as before!
python generator.py example_script.json traditional.mp4
```

## When to Use Each Style?

### Use **Chat Style** for:
- ✅ Full conversation context needed
- ✅ Multiple speakers
- ✅ Long conversations
- ✅ Reference to previous messages
- ✅ Tutorial/documentation videos

### Use **Novel Style** for:
- ✅ Dramatic dialogues
- ✅ Emotional exchanges
- ✅ Short-form content (TikTok, Reels, Shorts)
- ✅ Focus on current conversation
- ✅ Cinematic presentation
- ✅ When preparing for TTS narration

## TTS Integration (Next Step)

The novel style is designed with TTS in mind:

```bash
# Future usage (coming soon!)
python generator.py script.json output.mp4 \
  --style novel \
  --max-messages 2 \
  --tts \
  --voice-a "en-US-Neural2-A" \
  --voice-b "en-US-Neural2-C"
```

**Why novel style works well with TTS:**
- Prominent text is easier to read while listening
- Per-message display allows natural pacing
- Slide-in animations can sync with audio
- Less visual clutter = better focus on narration

## Technical Details

### Novel Renderer Architecture

```python
class NovelRenderer:
    - Filters messages to show only last N from each side
    - Divides screen into horizontal sections
    - Each message gets a prominent panel
    - Alternates alignment based on sender
    - Smooth slide-in animations with cubic easing
    - Larger fonts and avatars for emphasis
```

### Animation Timeline

1. **Slide-in** (0.0 → 0.3): Panel slides from appropriate side
2. **Avatar fade-in** (0.3 → 1.0): Avatar fades in
3. **Name fade-in** (0.5 → 1.0): Speaker name appears
4. **Text fade-in** (0.6 → 1.0): Message text fades in

### Message Filtering Logic

```python
# Separate messages by sender
person_a_messages = [msg for msg in messages if msg['sender'] == 'a']
person_b_messages = [msg for msg in messages if msg['sender'] == 'b']

# Get last N from each
recent_a = person_a_messages[-max_messages:]
recent_b = person_b_messages[-max_messages:]

# Merge back chronologically
combined = sorted(recent_a + recent_b, by_index)
```

## Configuration Options

All existing configuration still works:

```json
{
  "config": {
    "fps": 30,
    "width": 1080,
    "height": 1920,
    "hold_at_end": 2.0,
    "person_a_image": "avatar_a.png",
    "person_b_image": "avatar_b.png"
  }
}
```

Novel-specific settings are command-line only:
- `--style novel` - Enable novel style
- `--max-messages N` - Messages per side to show

## Testing Checklist

- [x] Chat style still works (backward compatible)
- [x] Novel style with default settings (3 messages)
- [x] Novel style with 1 message (minimal/TikTok style)
- [x] Novel style with 5+ messages (more context)
- [x] Animations smooth and properly timed
- [x] Avatars display correctly
- [x] Chinese/CJK text renders properly
- [x] Long messages wrap correctly
- [x] Command-line arguments work
- [x] Error handling for invalid options

## Known Limitations

1. **Typing indicators** not yet implemented in novel style (coming soon)
2. **Fixed layout** - messages always split screen evenly
3. **No custom colors** yet - uses same colors as chat style

These will be addressed in future updates!

## Customization Ideas

Want to customize the novel style? Edit `novel_renderer.py`:

```python
# Larger fonts
self.message_font_size = 80  # Default: 60

# Bigger avatars
self.avatar_size = 200  # Default: 150

# More spacing
panel_margin = 100  # Default: 80

# Different animation
def _ease_out_bounce(self, t):  # Replace _ease_out_cubic
    # Your custom easing function
```

## Next Steps

1. **Test the new style** with your content
   ```bash
   cd video_chat_renderer
   ./test_novel_style.sh
   ```

2. **Create your own scripts** using `example_novel.json` as a template

3. **Experiment with `--max-messages`** to find what works best

4. **Prepare for TTS integration** by timing your messages appropriately

5. **Share feedback** on what features you'd like next!

## Questions?

- How do I switch between styles? → Use `--style chat` or `--style novel`
- Can I use both styles in one video? → Not yet, but planned!
- Will this break my existing scripts? → No! Default is still chat style
- When is TTS coming? → Next major feature update
- Can I customize colors? → Edit `config.py` (applies to both styles)

---

**Enjoy creating cinematic chat videos! 🎬**
