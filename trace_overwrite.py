import json, os, glob

backup_dir = r'backups\presets\presets_lin_lu_CN'
pattern = os.path.join(backup_dir, 'user_1_Day6_dup_1_simplified.json.*.bak')
backups = sorted(glob.glob(pattern))

# Read current Day 5 first turn content for comparison
with open(r'presets\presets_lin_lu_CN\user_1_Day5_dup_1_simplified.json', 'r', encoding='utf-8') as f:
    d5 = json.load(f)
d5_first_content = d5['dialogue'][0]['content'] if d5.get('dialogue') else ''

print('Day 5 first turn content: {}...'.format(d5_first_content[:40]))
print()
print('=== Day 6 backup history ===')
for bak_path in backups:
    fname = os.path.basename(bak_path)
    try:
        with open(bak_path, 'r', encoding='utf-8') as f:
            data = json.load(f)
        dlg = data.get('dialogue', [])
        turns = len(dlg)
        first_content = dlg[0]['content'][:40] if dlg else ''
        first_role = dlg[0]['role'] if dlg else ''
        matches_d5 = (first_content == d5_first_content[:40]) if first_content else False
        flag = ' *** MATCHES DAY 5 ***' if matches_d5 else ''
        print('{}: {} turns, first=[{}] {}...{}'.format(fname, turns, first_role, first_content, flag))
    except Exception as e:
        print('{}: ERROR: {}'.format(fname, e))
