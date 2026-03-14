import json

# Read the CURRENT Day 6 file
with open(r'presets\presets_lin_lu_CN\user_1_Day6_dup_1_simplified.json', 'r', encoding='utf-8') as f:
    d6_current = json.load(f)

# Read the current Day 5 file
with open(r'presets\presets_lin_lu_CN\user_1_Day5_dup_1_simplified.json', 'r', encoding='utf-8') as f:
    d5 = json.load(f)

# Read the oldest backup of Day 6
with open(r'backups\presets\presets_lin_lu_CN\user_1_Day6_dup_1_simplified.json.20260312T050941000Z.bak', 'r', encoding='utf-8') as f:
    d6_oldest = json.load(f)

dlg5 = d5.get('dialogue', [])
dlg6_cur = d6_current.get('dialogue', [])
dlg6_old = d6_oldest.get('dialogue', [])

print('Day 5 current: {} turns'.format(len(dlg5)))
print('Day 6 current: {} turns'.format(len(dlg6_cur)))
print('Day 6 oldest backup: {} turns'.format(len(dlg6_old)))
print()

# Compare current Day 6 with current Day 5
print('=== Comparing CURRENT Day 5 vs CURRENT Day 6 (first 10 turns) ===')
for i in range(min(10, len(dlg5), len(dlg6_cur))):
    c5 = dlg5[i].get('content', '')
    r5 = dlg5[i].get('role', '?')
    c6 = dlg6_cur[i].get('content', '')
    r6 = dlg6_cur[i].get('role', '?')
    content_match = 'SAME' if c5 == c6 else 'DIFF'
    role_match = 'SAME' if r5 == r6 else 'SWAPPED'
    print('  Turn {}: content={}, role={}'.format(i, content_match, role_match))
    if content_match == 'SAME':
        print('    Content: {}...'.format(c5[:60]))
        print('    D5 role: {}, D6 role: {}'.format(r5, r6))

print()
# Compare current Day 6 with oldest Day 6 backup
print('=== Comparing oldest Day 6 backup vs CURRENT Day 6 (first 10 turns) ===')
for i in range(min(10, len(dlg6_old), len(dlg6_cur))):
    co = dlg6_old[i].get('content', '')
    ro = dlg6_old[i].get('role', '?')
    cc = dlg6_cur[i].get('content', '')
    rc = dlg6_cur[i].get('role', '?')
    content_match = 'SAME' if co == cc else 'DIFF'
    role_match = 'SAME' if ro == rc else 'SWAPPED'
    print('  Turn {}: content={}, role={}'.format(i, content_match, role_match))
    if content_match == 'DIFF':
        print('    Old: [{}] {}...'.format(ro, co[:60]))
        print('    Cur: [{}] {}...'.format(rc, cc[:60]))
