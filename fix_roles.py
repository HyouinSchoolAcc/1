import json
import os

def fix_roles(filepath, fixes):
    with open(filepath, 'r', encoding='utf-8') as f:
        data = json.load(f)
        
    changes_made = 0
    for msg in data.get('dialogue', []):
        content = msg.get('content', '')
        for fix_content, new_role in fixes.items():
            if fix_content in content:
                if msg['role'] != new_role:
                    print(f"Changing role for '{content}' to {new_role}")
                    msg['role'] = new_role
                    changes_made += 1
                    
    if changes_made > 0:
        with open(filepath, 'w', encoding='utf-8') as f:
            json.dump(data, f, ensure_ascii=False, indent=2)
        print(f"Saved {changes_made} changes to {os.path.basename(filepath)}")
    else:
        print(f"No changes needed for {os.path.basename(filepath)}")

if __name__ == "__main__":
    day6_fixes = {
        "姐就是女王，自信放光芒": "User",
        "能钓到吗？": "User"
    }
    fix_roles('data_labeler/presets/presets_lin_lu_CN/user_7_Day6_dup_1_simplified.json', day6_fixes)
    
    day8_fixes = {
        "大家都微死": "User",
        "不过有时我想聊一点更深的话题": "User",
        "更深的话题是指什么？": "林路"
    }
    fix_roles('data_labeler/presets/presets_lin_lu_CN/user_7_Day8_dup_1_simplified.json', day8_fixes)
