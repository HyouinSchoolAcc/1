import json
import os

script_dir = os.path.dirname(os.path.abspath(__file__))

# Read Day 1
with open(os.path.join(script_dir, 'user_7_Day1_dup_1_simplified.json'), 'r', encoding='utf-8') as f:
    day1 = json.load(f)

# Read Day 2
with open(os.path.join(script_dir, 'user_7_Day2_dup_1_simplified.json'), 'r', encoding='utf-8') as f:
    day2 = json.load(f)

dialogue = day1['dialogue']

# Find the cutoff: keep everything up to and including "全国那么多大学，总要有人去当大学老师啊"
cutoff_idx = None
for i, msg in enumerate(dialogue):
    if msg['content'] == '全国那么多大学，总要有人去当大学老师啊':
        cutoff_idx = i
        break

if cutoff_idx is None:
    print("ERROR: Could not find cutoff message!")
    exit(1)

print(f"Cutoff index: {cutoff_idx}")
print(f"Total Day 1 messages: {len(dialogue)}")
print(f"Messages to keep in Day 1: {cutoff_idx + 1}")
print(f"Messages to move to Day 2: {len(dialogue) - cutoff_idx - 1}")

# Messages to keep in Day 1
keep = dialogue[:cutoff_idx + 1]

# Fix any entries in 'keep' that are missing timestamps
for msg in keep:
    if 'timestamp' not in msg:
        msg['timestamp'] = '23:59'
        msg['timestampManuallyEdited'] = True

# Messages to move to Day 2
move = dialogue[cutoff_idx + 1:]

# Print moved messages for verification
print(f"\n--- {len(move)} Messages being moved to Day 2 ---")
for i, msg in enumerate(move):
    old_ts = msg.get('timestamp', 'NO_TS')
    content_preview = msg['content'][:40].replace('\n', ' ')
    print(f"  {i}: [{old_ts}] {msg['role']}: {content_preview}")

# Assign realistic post-midnight timestamps
new_timestamps = [
    # 0-4: Age discussion (5 msgs)
    "00:00", "00:00", "00:01", "00:02", "00:02",
    # 5-8: Young teacher discussion (4 msgs)
    "00:03", "00:03", "00:04", "00:04",
    # 9-15: 恋与深空 discussion (7 msgs)
    "00:05", "00:05", "00:05", "00:06", "00:06", "00:06", "00:07",
    # 16-21: Students & game discussion (6 msgs)
    "00:07", "00:08", "00:08", "00:09", "00:09", "00:10",
    # 22-25: Teacher-student relationship (4 msgs)
    "00:10", "00:11", "00:11", "00:12",
    # 26-29: 尽心尽责 + sticker (4 msgs)
    "00:12", "00:12", "00:12", "00:12",
    # 30-35: Height discussion (6 msgs)
    "00:13", "00:16", "00:16", "00:18", "00:18", "00:18",
    # 36-39: Northern/Southern (4 msgs)
    "00:19", "00:19", "00:19", "00:19",
    # 40-42: Hometown (3 msgs)
    "00:20", "00:20", "00:20",
    # 43-48: Nanjing (6 msgs)
    "00:21", "00:21", "00:21", "00:22", "00:22", "00:22",
    # 49: 好吧 (1 msg)
    "00:23",
    # 50-56: Going to sleep (7 msgs)
    "00:25", "00:25", "00:25", "00:26", "00:26", "00:27", "00:28",
    # 57-58: Unanswered messages (2 msgs)
    "00:35", "00:50",
]

if len(new_timestamps) != len(move):
    print(f"\nWARNING: Timestamp count ({len(new_timestamps)}) != message count ({len(move)})")
    # Auto-adjust
    while len(new_timestamps) < len(move):
        new_timestamps.append(new_timestamps[-1])
    new_timestamps = new_timestamps[:len(move)]

# Apply new timestamps
for i, msg in enumerate(move):
    msg['timestamp'] = new_timestamps[i]
    msg['timestampManuallyEdited'] = True

# Update Day 1
day1['dialogue'] = keep

# Prepend moved messages to Day 2 dialogue
day2['dialogue'] = move + day2['dialogue']

# Write modified Day 1
with open(os.path.join(script_dir, 'user_7_Day1_dup_1_simplified.json'), 'w', encoding='utf-8') as f:
    json.dump(day1, f, ensure_ascii=False, indent=2)

# Write modified Day 2
with open(os.path.join(script_dir, 'user_7_Day2_dup_1_simplified.json'), 'w', encoding='utf-8') as f:
    json.dump(day2, f, ensure_ascii=False, indent=2)

print(f"\nDone! Day 1 modified: kept {len(keep)} messages (up to index {cutoff_idx})")
print(f"Done! Day 2 modified: prepended {len(move)} messages before existing dialogue")
print("Files written successfully!")
