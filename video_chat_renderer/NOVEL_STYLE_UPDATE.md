# Novel Style Update - ONE Message at a Time

## What Changed?

I've updated the novel style renderer to match your vision: **ONE message displayed at a time, large and centered on screen**, like a visual novel or chat screenshot.

## New Behavior

### Before (Old Novel Style)
- Showed last N messages from each side
- Multiple messages on screen at once
- Messages stacked vertically

### Now (Fixed Novel Style) ✅
- Shows **ONLY ONE message** at a time
- Message takes up ~85% of screen width
- Large bubble centered vertically on screen
- Each message gets its own "scene"
- Flip-flops between left (person B) and right (person A)

## Visual Style

Matches your images exactly:

**Person B (left-aligned):**
```
[Avatar]  [Large white bubble with text taking up most of screen]
```

**Person A (right-aligned):**
```
           [Large green bubble with text taking up most of screen]  [Avatar]
```

## How It Works

### Timeline
```
Message 1: "Hey, did you finish the research paper?"
├─ Frames 0-60: Only this message visible, centered on screen
│
Message 2: "Yeah, I sent it to you this morning!"
├─ Frames 61-120: Message 1 gone, only message 2 visible
│
Message 3: "Oh, I must have missed it..."
└─ Frames 121-180: Message 2 gone, only message 3 visible
```

Each message has its own display duration (controlled by `delay_after` in JSON).

## Usage

```bash
# Generate novel style video (ONE message at a time)
python generator.py script.json output.mp4 --style novel

# The --max-messages flag is ignored in novel mode
# (it always shows 1 message)
python generator.py script.json output.mp4 --style novel --max-messages 3
# ↑ Still shows 1 message at a time
```

## JSON Script Format

Same as before - the `delay_after` now controls how long each individual message stays on screen:

```json
{
  "person_a": "Alice",
  "person_b": "Bob",
  "messages": [
    {
      "sender": "b",
      "text": "Hey, did you finish the research paper?",
      "delay_after": 3.0,  ← Message shows for 3 seconds
      "typing_duration": 1.0
    },
    {
      "sender": "a",
      "text": "Yeah, I sent it to you this morning!",
      "delay_after": 2.5,  ← Message shows for 2.5 seconds
      "typing_duration": 0.8
    }
  ]
}
```

## Key Features

✅ **ONE message at a time** - No clutter, full focus
✅ **Large display** - Bubble takes up 85% of screen width
✅ **Centered** - Message centered vertically on screen
✅ **Flip-flop layout** - Left for person B, right for person A
✅ **Smooth animations** - Slides in from appropriate side
✅ **Avatar display** - Shows avatar next to bubble
✅ **Perfect for TTS** - Each message has dedicated screen time

## Size Comparison

### Chat Style
- Bubble: Max 850px (80% width)
- Font: 36px
- Avatar: 90px
- Multiple messages visible

### Novel Style (New!)
- Bubble: ~917px (85% width) on 1080px screen
- Font: 60px (large and readable)
- Avatar: 150px (prominent)
- **ONE message visible**

## Test It

```bash
cd video_renderer
python generator.py example_novel.json test_novel.mp4 --style novel
```

Expected result:
- Each message appears individually
- Takes up most of the screen
- Looks like zoomed-in chat screenshots
- Messages flip-flop between sides

## Perfect For

✅ **Visual novels** - Story-driven content
✅ **TikTok/Reels** - Short-form vertical video
✅ **TTS narration** - One message = one voice clip
✅ **Dramatic dialogues** - Focus on each line
✅ **Clean presentation** - No visual clutter

## Technical Details

### Animation Timeline (per message)
- **0.0-1.0**: Slide in from appropriate side
- **0.3-1.0**: Avatar fades in
- **Fully visible**: Message displayed for `delay_after` seconds
- **Next message**: Previous message disappears, new one slides in

### Layout Calculations
```python
# Bubble width: 85% of screen
bubble_width = screen_width * 0.85

# Vertical centering
bubble_y = (screen_height - bubble_height) / 2

# Horizontal positioning
if person_a:
    bubble_x = screen_width - bubble_width - 100 - avatar_space
else:
    bubble_x = 100 + avatar_space
```

## Migration from Old Novel Style

If you used the old novel style with multiple messages:

**Old command:**
```bash
python generator.py script.json output.mp4 --style novel --max-messages 3
```

**New behavior:**
- Still accepts `--max-messages` flag (for backward compatibility)
- But **always shows 1 message** in novel mode
- Each message gets full screen focus

**If you want the old multi-message behavior:**
- Use `--style chat` (traditional scrolling chat)

## Examples

### Example 1: Simple Dialogue
```bash
python generator.py example_novel.json output.mp4 --style novel
```

Result:
1. First message fills screen for 3 seconds
2. Fades out, second message slides in
3. Each message gets dedicated display time

### Example 2: With Typing Indicators
```json
{
  "sender": "b",
  "text": "Let me think about that...",
  "delay_after": 2.5,
  "typing_duration": 1.5  ← Shows typing bubble for 1.5s first
}
```

Timeline:
- Typing indicator: 1.5 seconds
- Message appears: Displayed for 2.5 seconds
- Next message starts

## Comparison with Your Images

Your images show exactly what novel style now does:

**Your Image 1:**
`[Avatar] [Hey, did you finish the research paper?        ]`
- ✅ Left-aligned
- ✅ Avatar on left
- ✅ Large bubble
- ✅ Takes up most of width

**Your Image 2:**
`[        Yeah, I sent it to you this morning!] [Avatar]`
- ✅ Right-aligned
- ✅ Avatar on right
- ✅ Large bubble
- ✅ Takes up most of width

**Novel style = Full screen version of these images!**

## Next: TTS Integration

Now that each message has dedicated screen time, TTS integration is straightforward:

```python
# Future implementation
for message in messages:
    audio_clip = generate_tts(message['text'], voice)
    message['delay_after'] = audio_clip.duration  # Auto-sync
```

Each message:
1. Displays on screen
2. Has audio narrated
3. Stays visible for audio duration
4. Next message begins

Perfect sync! 🎯

---

**This is exactly what you wanted! Each message now gets the spotlight. 🎬**
