import json
import os

script_dir = os.path.dirname(os.path.abspath(__file__))

with open(os.path.join(script_dir, 'user_7_Day1_dup_1_simplified.json'), 'r', encoding='utf-8') as f:
    day1 = json.load(f)
with open(os.path.join(script_dir, 'user_7_Day2_dup_1_simplified.json'), 'r', encoding='utf-8') as f:
    day2 = json.load(f)

msgs1 = day1['dialogue']
msgs2 = day2['dialogue']

print(f"Day 1: {len(msgs1)} messages")
print("Day 1 - Last 5 messages:")
for m in msgs1[-5:]:
    ts = m.get('timestamp', '??')
    content = m['content'][:60].replace('\n', ' ')
    print(f"  [{ts}] {m['role']}: {content}")

print(f"\nDay 2: {len(msgs2)} messages")
print("Day 2 - First 10 messages:")
for m in msgs2[:10]:
    ts = m.get('timestamp', '??')
    content = m['content'][:60].replace('\n', ' ')
    print(f"  [{ts}] {m['role']}: {content}")

# Find transition point in Day 2 (where moved messages end and original begin)
print("\nDay 2 - Looking for transition (moved -> original):")
for i, m in enumerate(msgs2):
    ts = m.get('timestamp', '??')
    if ts == '09:15' or (i > 0 and msgs2[i-1].get('timestamp','') < '09:00' and ts >= '09:00'):
        print(f"  Transition at index {i}:")
        for j in range(max(0,i-2), min(len(msgs2), i+3)):
            mj = msgs2[j]
            tsj = mj.get('timestamp', '??')
            cj = mj['content'][:60].replace('\n', ' ')
            marker = " <-- transition" if j == i else ""
            print(f"    [{tsj}] {mj['role']}: {cj}{marker}")
        break
