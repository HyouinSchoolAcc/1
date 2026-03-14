"""Patch convert_dialogue_to_image.py to add Windows font paths."""
import os

filepath = 'convert_dialogue_to_image.py'
with open(filepath, 'r', encoding='utf-8') as f:
    lines = f.readlines()

# Find the FONT_PATHS block
start = None
end = None
for i, line in enumerate(lines):
    if 'FONT_PATHS = [' in line:
        start = i - 1 if i > 0 and lines[i-1].strip().startswith('#') else i
    if start is not None and 'raise RuntimeError' in line:
        end = i + 1
        break

if start is None or end is None:
    print("ERROR: Could not find FONT_PATHS block")
    exit(1)

print(f"Found block from line {start+1} to {end}")

new_block = """# Auto-detect Chinese fonts (Windows + Linux)
import platform as _platform
FONT_PATHS = []
if _platform.system() == "Windows":
    _wf = os.path.join(os.environ.get("WINDIR", r"C:\\Windows"), "Fonts")
    FONT_PATHS += [
        os.path.join(_wf, "msyh.ttc"),
        os.path.join(_wf, "msyhbd.ttc"),
        os.path.join(_wf, "simhei.ttf"),
        os.path.join(_wf, "simsun.ttc"),
        os.path.join(_wf, "meiryo.ttc"),
        os.path.join(_wf, "segoeui.ttf"),
    ]
FONT_PATHS += [
    "/usr/share/fonts/opentype/noto/NotoSansCJK-Regular.ttc",
    "/usr/share/fonts/truetype/noto/NotoSansCJK-Regular.ttc",
    "/usr/share/fonts/truetype/wqy/wqy-zenhei.ttc",
    "/usr/share/fonts/truetype/arphic/ukai.ttc",
    "/usr/share/fonts/truetype/arphic/uming.ttc",
    "/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
]
FONT_PATH = None
for path in FONT_PATHS:
    if os.path.exists(path):
        FONT_PATH = path
        break
if FONT_PATH is None:
    raise RuntimeError("No suitable font found for Chinese text!")
"""

new_lines = lines[:start] + [new_block] + lines[end:]

with open(filepath, 'w', encoding='utf-8') as f:
    f.writelines(new_lines)

print("SUCCESS: Patched font paths in convert_dialogue_to_image.py")
