with open('internal/web/router.go', 'rb') as f:
    data = f.read()
old = b'r.Get("/api/video/status", app.handleVideoStatus)'
new = old + b'\n\tr.Get("/api/video/status/stream", app.handleVideoStatusSSE)'
if old in data:
    data = data.replace(old, new, 1)
    with open('internal/web/router.go', 'wb') as f:
        f.write(data)
    print('SUCCESS: Added SSE route')
else:
    print('ERROR: route not found')
