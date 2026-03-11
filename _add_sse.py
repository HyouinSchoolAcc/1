import json

with open('internal/web/app.go', 'rb') as f:
    data = f.read()

# Find the handleVideoStatus function and add SSE handler after it
old = b'// handleVideoStatus reports the current status of a video generation job.\nfunc (a *App) handleVideoStatus(w http.ResponseWriter, r *http.Request) {\n'

if old not in data:
    print('ERROR: handleVideoStatus not found')
    exit(1)

# Find the end of handleVideoStatus (closing brace on its own line after the function)
start_idx = data.index(old)
# Find the function body - look for the closing } followed by \nfunc or \n\n
search_from = start_idx + len(old)
# Find "writeJSON(w, http.StatusOK, job)\n}\n"
end_marker = b'writeJSON(w, http.StatusOK, job)\n}\n'
end_idx = data.index(end_marker, search_from)
end_of_func = end_idx + len(end_marker)

sse_handler = b'''
// handleVideoStatusSSE streams video job status via Server-Sent Events.
// The client opens ONE connection; the server pushes updates when status changes.
func (a *App) handleVideoStatusSSE(w http.ResponseWriter, r *http.Request) {
\tjobKey := r.URL.Query().Get("job_key")
\tif jobKey == "" {
\t\thttp.Error(w, "missing job_key", http.StatusBadRequest)
\t\treturn
\t}

\tflusher, ok := w.(http.Flusher)
\tif !ok {
\t\thttp.Error(w, "streaming not supported", http.StatusInternalServerError)
\t\treturn
\t}

\tw.Header().Set("Content-Type", "text/event-stream")
\tw.Header().Set("Cache-Control", "no-cache")
\tw.Header().Set("Connection", "keep-alive")
\tw.Header().Set("X-Accel-Buffering", "no")

\tctx := r.Context()
\tticker := time.NewTicker(1 * time.Second)
\tdefer ticker.Stop()

\tlastStatus := ""
\tfor {
\t\tselect {
\t\tcase <-ctx.Done():
\t\t\treturn
\t\tcase <-ticker.C:
\t\t\ta.videoJobsMu.RLock()
\t\t\tjob, exists := a.videoJobs[jobKey]
\t\t\tvar jobCopy VideoJob
\t\t\tif exists {
\t\t\t\tjobCopy = *job
\t\t\t}
\t\t\ta.videoJobsMu.RUnlock()

\t\t\tif !exists {
\t\t\t\tfmt.Fprintf(w, "data: {\\"status\\":\\"unknown\\"}\\n\\n")
\t\t\t\tflusher.Flush()
\t\t\t\treturn
\t\t\t}

\t\t\tif jobCopy.Status != lastStatus {
\t\t\t\tlastStatus = jobCopy.Status
\t\t\t\tjsonBytes, _ := json.Marshal(jobCopy)
\t\t\t\tfmt.Fprintf(w, "data: %s\\n\\n", jsonBytes)
\t\t\t\tflusher.Flush()

\t\t\t\tif jobCopy.Status == "done" || jobCopy.Status == "error" {
\t\t\t\t\treturn
\t\t\t\t}
\t\t\t}
\t\t}
\t}
}
'''

data = data[:end_of_func] + sse_handler + data[end_of_func:]

with open('internal/web/app.go', 'wb') as f:
    f.write(data)

print('SUCCESS: Added handleVideoStatusSSE after handleVideoStatus')
