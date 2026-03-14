import json, os, glob

backup_dir = r'backups\presets\presets_lin_lu_CN'
pattern = os.path.join(backup_dir, 'user_1_Day6_dup_1_simplified.json.*.bak')
backups = sorted(glob.glob(pattern))

for bak_path in backups[-4:]:
    fname = os.path.basename(bak_path)
    with open(bak_path, 'r', encoding='utf-8') as f:
        data = json.load(f)
    dlg = data.get('dialogue', [])
    vals = str(data.get('values', ''))[:60]
    judg = str(data.get('judgements', ''))[:60]
    exp = str(data.get('experiences', ''))[:60]
    abil = str(data.get('abilities', ''))[:60]
    first_content = dlg[0].get('content', '')[:40] if dlg else 'EMPTY'
    print('{}:'.format(fname))
    print('  turns: {}'.format(len(dlg)))
    print('  first: {}...'.format(first_content))
    print('  values: {}...'.format(vals))
    print('  judgements: {}...'.format(judg))
    print('  experiences: {}...'.format(exp))
    print('  abilities: {}...'.format(abil))
    print()
