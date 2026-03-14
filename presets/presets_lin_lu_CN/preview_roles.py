import json
import os

def preview_roles(filepath):
    print(f"--- Preview for {os.path.basename(filepath)} ---")
    with open(filepath, 'r', encoding='utf-8') as f:
        data = json.load(f)
    for i, msg in enumerate(data.get('dialogue', [])):
        role = msg.get('role', 'Unknown')
        content = msg.get('content', '')
        if isinstance(content, str):
            content = content.replace('\n', ' ')[:50]
        print(f"{i:03d} | {role:10s} | {content}")

if __name__ == "__main__":
    preview_roles('data_labeler/presets/presets_lin_lu_CN/user_7_Day6_dup_1_simplified.json')
    print("\n")
    preview_roles('data_labeler/presets/presets_lin_lu_CN/user_7_Day8_dup_1_simplified.json')
