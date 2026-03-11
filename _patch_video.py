#!/usr/bin/env python3
"""Patch clean backup files to add video dropdown functionality."""
import os

os.chdir(r'C:\Users\user\Desktop\data_labeler')

# ============================================================
# 1. Patch main_panel.html - add video dropdown HTML + CSS
# ============================================================
with open('additional_contents/templates/includes/main_panel.html', 'rb') as f:
    panel = f.read()

# Video CSS
video_css = b'<style>\n  .video-dropdown-item:hover:not(:disabled) {\n    background: #f0f0f0 !important;\n    color: #e05585 !important;\n  }\n</style>\n'

# Video dropdown HTML - using raw bytes for emoji/CJK characters
# \xf0\x9f\x8e\xac = movie clapper emoji
# \xe8\xa7\x86\xe9\xa2\x91 = video (Chinese)
# \xe2\x96\xb2 = up triangle
# \xf0\x9f\x94\x84 = counterclockwise arrows emoji
# \xe7\x94\x9f\xe6\x88\x90\xe8\xa7\x86\xe9\xa2\x91 = generate video (Chinese)
# \xe2\x96\xb6\xef\xb8\x8f = play button emoji
# \xe9\xa2\x84\xe8\xa7\x88\xe8\xa7\x86\xe9\xa2\x91 = preview video (Chinese)
# \xe2\xac\x87\xef\xb8\x8f = down arrow emoji
# \xe4\xb8\x8b\xe8\xbd\xbd\xe8\xa7\x86\xe9\xa2\x91 = download video (Chinese)

lines = []
lines.append(b'    <!-- Video dropdown -->')
lines.append(b'    <div class="video-dropdown" style="position: relative; display: inline-block;">')
lines.append(b'      <button id="btnVideoDropdown"')
lines.append(b'              class="btn pink-text-button"')
lines.append(b'              onclick="toggleVideoDropdown()"')
lines.append(b'              style="display: inline-flex; align-items: center; gap: 0.4rem; padding: 0.5rem 1rem; font-size: 0.875rem;">')
lines.append(b'        <span>\xf0\x9f\x8e\xac</span><span>{{if eq .Language "en"}}Video{{else}}\xe8\xa7\x86\xe9\xa2\x91{{end}}</span><span style="font-size:0.6rem;margin-left:0.15rem;">\xe2\x96\xb2</span>')
lines.append(b'      </button>')
lines.append(b'      <div id="videoDropdownMenu" style="display:none; position:absolute; bottom:100%; left:0; z-index:1050; min-width:180px; background:#fff; border:1px solid #dee2e6; border-radius:0.5rem; box-shadow:0 4px 16px rgba(0,0,0,0.13); padding:0.35rem 0; margin-bottom:0.25rem;">')
lines.append(b'        <button class="video-dropdown-item" onclick="generateVideo()" style="display:flex;align-items:center;gap:0.5rem;width:100%;padding:0.5rem 1rem;border:none;background:none;cursor:pointer;font-size:0.85rem;color:#333;text-align:left;">')
lines.append(b'          <span>\xf0\x9f\x94\x84</span><span>{{if eq .Language "en"}}Generate Video{{else}}\xe7\x94\x9f\xe6\x88\x90\xe8\xa7\x86\xe9\xa2\x91{{end}}</span>')
lines.append(b'        </button>')
lines.append(b'        <button class="video-dropdown-item" onclick="previewVideo()" style="display:flex;align-items:center;gap:0.5rem;width:100%;padding:0.5rem 1rem;border:none;background:none;cursor:pointer;font-size:0.85rem;color:#333;text-align:left;">')
lines.append(b'          <span>\xe2\x96\xb6\xef\xb8\x8f</span><span>{{if eq .Language "en"}}Preview Video{{else}}\xe9\xa2\x84\xe8\xa7\x88\xe8\xa7\x86\xe9\xa2\x91{{end}}</span>')
lines.append(b'        </button>')
lines.append(b'        <button id="btnDownloadVideo" class="video-dropdown-item" onclick="downloadVideo()" disabled style="display:flex;align-items:center;gap:0.5rem;width:100%;padding:0.5rem 1rem;border:none;background:none;cursor:pointer;font-size:0.85rem;color:#999;text-align:left;opacity:0.5;">')
lines.append(b'          <span>\xe2\xac\x87\xef\xb8\x8f</span><span>{{if eq .Language "en"}}Download Video{{else}}\xe4\xb8\x8b\xe8\xbd\xbd\xe8\xa7\x86\xe9\xa2\x91{{end}}</span>')
lines.append(b'        </button>')
lines.append(b'      </div>')
lines.append(b'    </div>')
video_html = b'\n'.join(lines) + b'\n'

# Insert CSS at the beginning
panel = video_css + panel

# Insert video dropdown after the "Generate Image" button
marker = b'Generate Image{{else}}'
idx = panel.find(marker)
if idx < 0:
    print("ERROR: Could not find Generate Image marker")
    exit(1)

end_btn = panel.find(b'</button>\n', idx)
if end_btn < 0:
    print("ERROR: Could not find button end")
    exit(1)

insert_pos = end_btn + len(b'</button>\n')
panel = panel[:insert_pos] + video_html + panel[insert_pos:]
print("OK: main_panel.html patched with video dropdown")

with open('templates/includes/main_panel.html', 'wb') as f:
    f.write(panel)
print("  Wrote %d bytes" % len(panel))


# ============================================================
# 2. Patch writing_main.html - add video JS functions
# ============================================================
with open('writing_main.html', 'rb') as f:
    html = f.read()

with open('_video_js_block.js', 'r', encoding='utf-8') as f:
    js_text = f.read()
video_js = js_text.encode('utf-8')

dl_func_pos = html.find(b'function downloadDialogueImage()')
if dl_func_pos < 0:
    print("ERROR: Could not find downloadDialogueImage")
    exit(1)

script_end = html.find(b'\n\n  </script>', dl_func_pos)
if script_end < 0:
    script_end = html.find(b'\n  </script>', dl_func_pos)
if script_end < 0:
    print("ERROR: Could not find </script>")
    exit(1)

html = html[:script_end] + video_js + html[script_end:]
print("OK: writing_main.html patched with video JS")

with open('templates/writing_main.html', 'wb') as f:
    f.write(html)
print("  Wrote %d bytes" % len(html))


# ============================================================
# 3. Verify no corruption
# ============================================================
print("\nVerification:")
for p in ['templates/includes/main_panel.html', 'templates/writing_main.html']:
    with open(p, 'rb') as f:
        d = f.read()
    c = sum(1 for i in range(len(d)-2) if d[i] >= 0xc0 and d[i+1] >= 0x80 and d[i+1] <= 0xbf and d[i+2] == 0x3f)
    print("  %s: %d bytes, %d corruptions" % (p, len(d), c))

print("\nDone!")
