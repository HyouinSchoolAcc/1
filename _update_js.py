with open('templates/writing_main.html', 'rb') as f:
    data = f.read()

# Find pollVideoStatus start
start_marker = b'function pollVideoStatus(jobKey) {'
start_idx = data.index(start_marker)
# Walk backward to include the leading whitespace/indentation
line_start = data.rfind(b'\n', 0, start_idx) + 1
func_start = line_start

# Find the closing brace by counting braces
brace_count = 0
i = data.index(b'{', start_idx)
while i < len(data):
    if data[i:i+1] == b'{':
        brace_count += 1
    elif data[i:i+1] == b'}':
        brace_count -= 1
        if brace_count == 0:
            func_end = i + 1
            break
    i += 1

old_func = data[func_start:func_end]
print(f'Found pollVideoStatus: bytes {func_start}-{func_end} ({func_end-func_start} bytes)')
print(f'First 80 chars: {old_func[:80]}')
print(f'Last 80 chars: {old_func[-80:]}')

# Build replacement
new_func = b"""    function streamVideoStatus(jobKey) {
      const btn = document.getElementById('btnVideoDropdown');
      const startTime = Date.now();

      const source = new EventSource('/api/video/status/stream?job_key=' + encodeURIComponent(jobKey));

      // Update button with elapsed time (local timer, no network traffic)
      const timeUpdater = setInterval(() => {
        const elapsed = Math.floor((Date.now() - startTime) / 1000);
        const mins = Math.floor(elapsed / 60);
        const secs = elapsed % 60;
        const timeStr = mins > 0 ? mins + 'm' + secs + 's' : secs + 's';
        if (btn) {
          btn.innerHTML = '<span>\xe2\x8f\xb3</span><span>' + (languageMode === 'en' ? 'Generating ' : '\xe7\x94\x9f\xe6\x88\x90\xe4\xb8\xad ') + timeStr + '</span>';
        }
      }, 1000);

      source.onmessage = function(event) {
        const job = JSON.parse(event.data);

        if (job.status === 'done') {
          source.close();
          clearInterval(timeUpdater);
          resetVideoButton();
          checkVideoAvailability();
          alert(languageMode === 'en'
            ? 'Video generated!\\n' + job.message + '\\nSize: ' + job.size_mb + ' MB'
            : '\xe8\xa7\x86\xe9\xa2\x91\xe7\x94\x9f\xe6\x88\x90\xe6\x88\x90\xe5\x8a\x9f\xef\xbc\x81\\n' + job.message + '\\n\xe5\xa4\xa7\xe5\xb0\x8f: ' + job.size_mb + ' MB');
        } else if (job.status === 'error') {
          source.close();
          clearInterval(timeUpdater);
          resetVideoButton();
          alert((languageMode === 'en' ? 'Video generation failed: ' : '\xe8\xa7\x86\xe9\xa2\x91\xe7\x94\x9f\xe6\x88\x90\xe5\xa4\xb1\xe8\xb4\xa5: ') + (job.error || 'Unknown error'));
        }
      };

      source.onerror = function() {
        source.close();
        clearInterval(timeUpdater);
        resetVideoButton();
      };
    }"""

data = data[:func_start] + new_func + data[func_end:]

# Replace call site: pollVideoStatus -> streamVideoStatus
data = data.replace(b'pollVideoStatus(data.job_key)', b'streamVideoStatus(data.job_key)')

# Remove unused _videoPollTimer variable
data = data.replace(b'    let _videoPollTimer = null;\n', b'')

with open('templates/writing_main.html', 'wb') as f:
    f.write(data)

print('SUCCESS: Replaced polling with SSE EventSource')
