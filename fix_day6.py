import json, shutil

current_path = r'presets\presets_lin_lu_CN\user_1_Day6_dup_1_simplified.json'
good_backup = r'backups\presets\presets_lin_lu_CN\user_1_Day6_dup_1_simplified.json.20260313T041515000Z.bak'

# Read current (corrupted) file - has correct profile fields
with open(current_path, 'r', encoding='utf-8') as f:
    current = json.load(f)

# Read last good backup - has correct dialogue
with open(good_backup, 'r', encoding='utf-8') as f:
    good = json.load(f)

# Create a safety backup of the corrupted file
shutil.copy2(current_path, current_path + '.corrupted_backup')

# Restore dialogue from good backup
current['dialogue'] = good['dialogue']

# Write fixed file
with open(current_path, 'w', encoding='utf-8') as f:
    json.dump(current, f, ensure_ascii=False, indent=2)

print('Fixed! Restored {} dialogue turns from backup.'.format(len(good['dialogue'])))
print('Kept profile fields from current file:')
print('  values: {}...'.format(str(current.get('values', ''))[:80]))
print('  judgements: {}...'.format(str(current.get('judgements', ''))[:80]))
print('Corrupted version backed up to: {}'.format(current_path + '.corrupted_backup'))
