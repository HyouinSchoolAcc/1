import json

# Read the earliest backup of Day 6
with open(r'backups\presets\presets_lin_lu_CN\user_1_Day6_dup_1_simplified.json.20260312T050941000Z.bak', 'r', encoding='utf-8') as f:
    d6_oldest = json.load(f)
dlg = d6_oldest.get('dialogue', [])
us = d6_oldest.get('user_schedule', {})
print('Oldest backup of Day 6: {} turns'.format(len(dlg)))
print('user_schedule.day: {}'.format(us.get('day', 'N/A') if isinstance(us, dict) else us))
for i, t in enumerate(dlg[:5]):
    content = t.get('content', '')[:80]
    print('  [{:03d}] {} | {}'.format(i, t.get('role','?'), content))
print()

# Compare with Day 5 first turns 
with open(r'presets\presets_lin_lu_CN\user_1_Day5_dup_1_simplified.json', 'r', encoding='utf-8') as f:
    d5 = json.load(f)
dlg5 = d5.get('dialogue', [])
print('Current Day 5: {} turns'.format(len(dlg5)))
for i, t in enumerate(dlg5[:5]):
    content = t.get('content', '')[:80]
    print('  [{:03d}] {} | {}'.format(i, t.get('role','?'), content))

# Check if first dialogue content matches but roles differ
print()
print('=== Comparing first 10 turns content ===')
for i in range(min(10, len(dlg), len(dlg5))):
    c5 = dlg5[i].get('content','')[:50]
    r5 = dlg5[i].get('role','?')
    c6 = dlg[i].get('content','')[:50]
    r6 = dlg[i].get('role','?')
    match = 'SAME' if c5 == c6 else 'DIFF'
    role_match = 'SAME' if r5 == r6 else 'SWAPPED'
    print('  Turn {}: content={}, role={} (D5={}, D6={})'.format(i, match, role_match, r5, r6))
