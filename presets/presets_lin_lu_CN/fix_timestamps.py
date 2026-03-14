import json
import os

script_dir = os.path.dirname(os.path.abspath(__file__))

with open(os.path.join(script_dir, 'user_7_Day2_dup_1_simplified.json'), 'r', encoding='utf-8') as f:
    day2 = json.load(f)

msgs = day2['dialogue']

# The first 59 messages are the ones moved from Day 1
# Let's verify them and apply correct timestamps

# Correct timestamps for all 59 moved messages
correct_timestamps = [
    # 0-4: Age discussion (5 msgs)
    "00:00", "00:00", "00:01", "00:02", "00:02",
    # 5-8: Young teacher discussion (4 msgs)
    "00:03", "00:03", "00:04", "00:04",
    # 9-15: 恋与深空 discussion (7 msgs)
    "00:05", "00:05", "00:05", "00:06", "00:06", "00:06", "00:07",
    # 16-21: Students & game discussion (6 msgs)
    "00:07", "00:08", "00:08", "00:09", "00:09", "00:10",
    # 22-25: Teacher-student relationship comments (4 msgs)
    "00:10", "00:11", "00:11", "00:12",
    # 26-29: 尽心尽责 + sticker (4 msgs)
    "00:12", "00:12", "00:12", "00:12",
    # 30-35: Height discussion (6 msgs)
    "00:13", "00:16", "00:16", "00:18", "00:18", "00:18",
    # 36-39: Northern/Southern + Zhejiang (4 msgs)
    "00:19", "00:19", "00:19", "00:19",
    # 40-42: Hometown exchange (3 msgs)
    "00:20", "00:20", "00:20",
    # 43-48: Nanjing discussion (6 msgs)
    "00:21", "00:21", "00:21", "00:22", "00:22", "00:22",
    # 49: 好吧 (1 msg)
    "00:23",
    # 50-56: Going to sleep exchange (7 msgs)
    "00:25", "00:25", "00:25", "00:26", "00:26", "00:27", "00:28",
    # 57-58: Unanswered messages (2 msgs)
    "00:35", "00:50",
]

assert len(correct_timestamps) == 59, f"Expected 59 timestamps, got {len(correct_timestamps)}"

# Print what's being fixed
print("Fixing timestamps for first 59 messages in Day 2:")
changes = 0
for i in range(59):
    old_ts = msgs[i].get('timestamp', '??')
    new_ts = correct_timestamps[i]
    content = msgs[i]['content'][:50].replace('\n', ' ')
    if old_ts != new_ts:
        print(f"  [{i}] {old_ts} -> {new_ts}: {msgs[i]['role']}: {content}")
        changes += 1
    msgs[i]['timestamp'] = new_ts
    msgs[i]['timestampManuallyEdited'] = True

print(f"\nFixed {changes} timestamps")

# Verify transition
print(f"\nTransition point:")
print(f"  [{msgs[58].get('timestamp')}] {msgs[58]['role']}: {msgs[58]['content'][:50]}")
print(f"  [{msgs[59].get('timestamp')}] {msgs[59]['role']}: {msgs[59]['content'][:50]}")

with open(os.path.join(script_dir, 'user_7_Day2_dup_1_simplified.json'), 'w', encoding='utf-8') as f:
    json.dump(day2, f, ensure_ascii=False, indent=2)

print("\nDone! File written.")
