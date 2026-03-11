
    // ========== Video Dropdown Logic ==========
    let _videoDropdownOpen = false;

    function toggleVideoDropdown() {
      const menu = document.getElementById('videoDropdownMenu');
      if (!menu) return;
      _videoDropdownOpen = !_videoDropdownOpen;
      menu.style.display = _videoDropdownOpen ? 'block' : 'none';
      if (_videoDropdownOpen) {
        checkVideoAvailability();
      }
    }

    // Close dropdown when clicking outside
    document.addEventListener('click', function(e) {
      const dropdown = document.querySelector('.video-dropdown');
      const menu = document.getElementById('videoDropdownMenu');
      if (dropdown && menu && !dropdown.contains(e.target)) {
        _videoDropdownOpen = false;
        menu.style.display = 'none';
      }
    });

    // Check if a video exists for the current file
    function checkVideoAvailability() {
      const dlBtn = document.getElementById('btnDownloadVideo');
      if (!dlBtn) return;
      if (!currentWriterFile || !currentWriterFile.filename || !presetSet) {
        dlBtn.disabled = true;
        dlBtn.style.opacity = '0.5';
        dlBtn.style.color = '#999';
        return;
      }
      fetch('/api/video/check?filename=' + encodeURIComponent(currentWriterFile.filename) + '&preset_set=' + encodeURIComponent(presetSet))
        .then(r => r.json())
        .then(data => {
          if (data.exists) {
            dlBtn.disabled = false;
            dlBtn.style.opacity = '1';
            dlBtn.style.color = '#333';
            dlBtn.title = (languageMode === 'en'
              ? 'Video ready (' + data.size_mb + ' MB, generated ' + data.generated + ')'
              : '\u89c6\u9891\u5df2\u5c31\u7eea (' + data.size_mb + ' MB, \u751f\u6210\u4e8e ' + data.generated + ')');
          } else {
            dlBtn.disabled = true;
            dlBtn.style.opacity = '0.5';
            dlBtn.style.color = '#999';
            dlBtn.title = languageMode === 'en' ? 'No video available - generate one first' : '\u5c1a\u65e0\u89c6\u9891 - \u8bf7\u5148\u751f\u6210';
          }
        })
        .catch(() => {
          dlBtn.disabled = true;
          dlBtn.style.opacity = '0.5';
        });
    }

    // Generate a video from the current dialogue (async with SSE)
    function generateVideo() {
      _videoDropdownOpen = false;
      document.getElementById('videoDropdownMenu').style.display = 'none';

      if (!currentWriterFile || !currentWriterFile.filename || !presetSet) {
        alert(languageMode === 'en' ? 'Please select character, user, version and day first' : '\u8bf7\u5148\u9009\u62e9AI\u89d2\u8272\u3001\u7528\u6237\u3001\u7248\u672c\u548c\u5929\u6570');
        return;
      }

      const btn = document.getElementById('btnVideoDropdown');
      if (btn) {
        btn.disabled = true;
        btn.innerHTML = '<span>\u231b</span><span>' + (languageMode === 'en' ? 'Generating...' : '\u751f\u6210\u4e2d...') + '</span>';
      }

      fetch('/api/video/generate', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          filename: currentWriterFile.filename,
          preset_set: presetSet
        })
      })
      .then(r => r.json())
      .then(data => {
        if (data.status === 'started' || data.status === 'running' || data.status === 'already_running') {
          streamVideoStatus(data.job_key);
        } else if (data.error) {
          resetVideoButton();
          alert((languageMode === 'en' ? 'Video generation failed: ' : '\u89c6\u9891\u751f\u6210\u5931\u8d25: ') + data.error);
        }
      })
      .catch(err => {
        console.error('Video generation error:', err);
        resetVideoButton();
        alert((languageMode === 'en' ? 'Video generation failed: ' : '\u89c6\u9891\u751f\u6210\u5931\u8d25: ') + err.message);
      });
    }

    function streamVideoStatus(jobKey) {
      const btn = document.getElementById('btnVideoDropdown');
      const startTime = Date.now();

      const source = new EventSource('/api/video/status/stream?job_key=' + encodeURIComponent(jobKey));

      const timeUpdater = setInterval(() => {
        const elapsed = Math.floor((Date.now() - startTime) / 1000);
        const mins = Math.floor(elapsed / 60);
        const secs = elapsed % 60;
        const timeStr = mins > 0 ? mins + 'm' + secs + 's' : secs + 's';
        if (btn) {
          btn.innerHTML = '<span>\u23f3</span><span>' + (languageMode === 'en' ? 'Generating ' : '\u751f\u6210\u4e2d ') + timeStr + '</span>';
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
            ? 'Video generated!\n' + job.message + '\nSize: ' + job.size_mb + ' MB'
            : '\u89c6\u9891\u751f\u6210\u6210\u529f\uff01\n' + job.message + '\n\u5927\u5c0f: ' + job.size_mb + ' MB');
        } else if (job.status === 'error') {
          source.close();
          clearInterval(timeUpdater);
          resetVideoButton();
          alert((languageMode === 'en' ? 'Video generation failed: ' : '\u89c6\u9891\u751f\u6210\u5931\u8d25: ') + (job.error || 'Unknown error'));
        }
      };

      source.onerror = function() {
        source.close();
        clearInterval(timeUpdater);
        resetVideoButton();
      };
    }

    function resetVideoButton() {
      const btn = document.getElementById('btnVideoDropdown');
      if (btn) {
        btn.disabled = false;
        btn.innerHTML = '<span>\ud83c\udfac</span><span>' + (languageMode === 'en' ? 'Video' : '\u89c6\u9891') + '</span><span style="font-size:0.6rem;margin-left:0.15rem;">\u25b2</span>';
      }
    }

    // Preview video in a modal overlay
    function previewVideo() {
      _videoDropdownOpen = false;
      document.getElementById('videoDropdownMenu').style.display = 'none';

      if (!currentWriterFile || !currentWriterFile.filename || !presetSet) {
        alert(languageMode === 'en' ? 'Please select character, user, version and day first' : '\u8bf7\u5148\u9009\u62e9AI\u89d2\u8272\u3001\u7528\u6237\u3001\u7248\u672c\u548c\u5929\u6570');
        return;
      }

      fetch('/api/video/check?filename=' + encodeURIComponent(currentWriterFile.filename) + '&preset_set=' + encodeURIComponent(presetSet))
        .then(r => r.json())
        .then(data => {
          if (!data.exists) {
            alert(languageMode === 'en' ? 'No video available yet. Please generate one first.' : '\u5c1a\u65e0\u89c6\u9891\uff0c\u8bf7\u5148\u751f\u6210\u3002');
            return;
          }
          const videoUrl = '/api/video/download?filename=' + encodeURIComponent(currentWriterFile.filename) + '&preset_set=' + encodeURIComponent(presetSet);
          showVideoPreviewModal(videoUrl, data.filename);
        })
        .catch(err => {
          alert((languageMode === 'en' ? 'Error checking video: ' : '\u68c0\u67e5\u89c6\u9891\u65f6\u51fa\u9519: ') + err.message);
        });
    }

    // Show a modal with embedded video player
    function showVideoPreviewModal(videoUrl, filename) {
      const existing = document.getElementById('videoPreviewModal');
      if (existing) existing.remove();

      const modal = document.createElement('div');
      modal.id = 'videoPreviewModal';
      modal.style.cssText = 'position:fixed;top:0;left:0;width:100%;height:100%;background:rgba(0,0,0,0.7);z-index:9999;display:flex;align-items:center;justify-content:center;';
      modal.innerHTML = `
        <div style="background:#fff;border-radius:1rem;padding:1.25rem;max-width:480px;width:95%;max-height:90vh;display:flex;flex-direction:column;box-shadow:0 8px 32px rgba(0,0,0,0.3);">
          <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:0.75rem;">
            <span style="font-weight:600;font-size:0.95rem;color:#333;">${filename}</span>
            <button onclick="document.getElementById('videoPreviewModal').remove()" style="background:none;border:none;font-size:1.4rem;cursor:pointer;color:#666;line-height:1;">\u2715</button>
          </div>
          <video controls autoplay style="width:100%;max-height:70vh;border-radius:0.5rem;background:#000;" src="${videoUrl}"></video>
        </div>
      `;
      modal.addEventListener('click', function(e) {
        if (e.target === modal) modal.remove();
      });
      document.body.appendChild(modal);
    }

    // Download the generated video
    function downloadVideo() {
      _videoDropdownOpen = false;
      document.getElementById('videoDropdownMenu').style.display = 'none';

      if (!currentWriterFile || !currentWriterFile.filename || !presetSet) {
        alert(languageMode === 'en' ? 'Please select character, user, version and day first' : '\u8bf7\u5148\u9009\u62e9AI\u89d2\u8272\u3001\u7528\u6237\u3001\u7248\u672c\u548c\u5929\u6570');
        return;
      }

      const downloadUrl = '/api/video/download?filename=' + encodeURIComponent(currentWriterFile.filename) + '&preset_set=' + encodeURIComponent(presetSet);
      const a = document.createElement('a');
      a.style.display = 'none';
      a.href = downloadUrl;
      a.download = '';
      document.body.appendChild(a);
      a.click();
      setTimeout(() => a.remove(), 3000);
    }
