import json, sys

def show_day(path, label):
    with open(path, 'r', encoding='utf-8') as f:
        data = json.load(f)
    dlg = data.get('dialogue', [])
    print(f'=== {label} ===')
    print(f'  user_schedule.day: {data.get("user_schedule", {}).get("day", "N/A")}')
    print(f'  character_schedule.day: {data.get("character_schedule", {}).get("day", "N/A")}')
    print(f'  user_name: {data.get("user_name", "N/A")}')
    print(f'  character_name: {data.get("character_name", "N/A")}')
    print(f'  completed: {data.get("completed", "N/A")}')
    print(f'  intimacy_level: {data.get("intimacy_level", "N/A")}')
    print(f'  starting_intimacy_level: {data.get("starting_intimacy_level", "N/A")}')
    print(f'  judgements: {str(data.get("judgements", ""))[:120]}')
    print(f'  values: {str(data.get("values", ""))[:120]}')
    print(f'  Total dialogue turns: {len(dlg)}')
    print()
    # First 5 turns
    for i, t in enumerate(dlg[:5]):
        content = t.get('content', '')[:80]
        print(f'  [{i:03d}] {t.get("role","?"):12s} | {content}')
    if len(dlg) > 10:
        print(f'  ... ({len(dlg) - 10} turns omitted) ...')
    # Last 5 turns
    for i in range(max(5, len(dlg)-5), len(dlg)):
        t = dlg[i]
        content = t.get('content', '')[:80]
        print(f'  [{i:03d}] {t.get("role","?"):12s} | {content}')
    print()

base = r'C:\Users\user\Desktop\data_labeler\presets\presets_lin_lu_CN'
show_day(f'{base}\\user_1_Day5_dup_1_simplified.json', 'Day 5')
show_day(f'{base}\\user_1_Day6_dup_1_simplified.json', 'Day 6')

# Check if dialogue content is identical
with open(f'{base}\\user_1_Day5_dup_1_simplified.json', 'r', encoding='utf-8') as f:
    d5 = json.load(f)
with open(f'{base}\\user_1_Day6_dup_1_simplified.json', 'r', encoding='utf-8') as f:
    d6 = json.load(f)

dlg5 = d5.get('dialogue', [])
dlg6 = d6.get('dialogue', [])

if dlg5 == dlg6:
    print('!!! DIALOGUES ARE IDENTICAL !!!')
elif len(dlg5) == len(dlg6):
    print(f'Dialogues have same length ({len(dlg5)}) but differ')
    for i in range(len(dlg5)):
        if dlg5[i] != dlg6[i]:
            print(f'  First diff at turn {i}:')
            print(f'    Day5: role={dlg5[i].get("role")} content={dlg5[i].get("content","")[:60]}')
            print(f'    Day6: role={dlg6[i].get("role")} content={dlg6[i].get("content","")[:60]}')
            break
else:
    print(f'Dialogues differ in length: Day5={len(dlg5)}, Day6={len(dlg6)}')

# Check for role swaps
print('\n=== Role distribution ===')
for label, dlg in [('Day5', dlg5), ('Day6', dlg6)]:
    roles = {}
    for t in dlg:
        r = t.get('role', '?')
        roles[r] = roles.get(r, 0) + 1
    print(f'  {label}: {roles}')
