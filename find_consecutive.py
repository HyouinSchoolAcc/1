import json
def check(f):
    print(f'\n{f}')
    with open(f, encoding='utf-8') as file:
        data = json.load(file)
    d = data['dialogue']
    
    last_role = None
    count = 0
    start_idx = 0
    
    for i in range(len(d)):
        role = d[i]['role']
        if role == last_role:
            count += 1
        else:
            if count > 2: # More than 2 consecutive messages by the same person could be suspect
                pass # We will just print them all
            last_role = role
            count = 1
            start_idx = i
            
    # Actually, better to print where the speaker starts and how many messages they sent, to identify large blocks.
    blocks = []
    current_block = {'role': d[0]['role'], 'start': 0, 'count': 1, 'text': d[0]['content'][:20]}
    
    for i in range(1, len(d)):
        if d[i]['role'] == current_block['role']:
            current_block['count'] += 1
        else:
            blocks.append(current_block)
            current_block = {'role': d[i]['role'], 'start': i, 'count': 1, 'text': d[i]['content'][:20]}
    blocks.append(current_block)
    
    for b in blocks:
        if b['count'] >= 4:
            print(f"Large block: {b['start']} to {b['start']+b['count']-1} | {b['role']} ({b['count']} msgs) | starts with: {b['text']}")

check('data_labeler/presets/presets_lin_lu_CN/user_7_Day6_dup_1_simplified.json')
check('data_labeler/presets/presets_lin_lu_CN/user_7_Day8_dup_1_simplified.json')
