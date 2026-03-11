package web

import (
	"bytes"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

// TemplateResolver is a minimal include-aware renderer for the existing Jinja-style templates.
// It only supports `{% include "path.html" %}` and otherwise leaves the file content intact.
// This keeps the HTML identical while avoiding a full templating rewrite.
type TemplateResolver struct {
	Root     string
	BaseDir  string
	DevMode  bool
	staticFS fs.FS

	mu    sync.RWMutex
	cache map[string]string
}

// ServeTemplate returns an http.HandlerFunc that writes the resolved HTML.
func (tr *TemplateResolver) ServeTemplate(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		html, err := tr.Resolve(name)
		if err != nil {
			http.Error(w, fmt.Sprintf("template error: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(html))
	}
}

var includeRE = regexp.MustCompile(`\{\%\s*include\s+"([^"]+)"\s*\%\}`)
var staticRE = regexp.MustCompile(`\{\{\s*url_for\(['"]static['"],\s*filename=['"]([^'"]+)['"]\)\s*\}\}`)

// Jinja2 block patterns - these will be stripped out
var jinjaWithRE = regexp.MustCompile(`\{\%\s*with\s+[^%]+\%\}`)
var jinjaIfRE = regexp.MustCompile(`\{\%\s*if\s+[^%]+\%\}`)
var jinjaForRE = regexp.MustCompile(`\{\%\s*for\s+[^%]+\%\}`)
var jinjaEndRE = regexp.MustCompile(`\{\%\s*(endif|endfor|endwith)\s*\%\}`)
var jinjaVarRE = regexp.MustCompile(`\{\{\s*[^}]+\s*\}\}`)

// Resolve loads a template and expands include directives recursively.
func (tr *TemplateResolver) Resolve(name string) (string, error) {
	if !tr.DevMode {
		tr.mu.RLock()
		if cached, ok := tr.cache[name]; ok {
			tr.mu.RUnlock()
			return cached, nil
		}
		tr.mu.RUnlock()
	}

	fullPath := filepath.Join(tr.BaseDir, name)
	content, err := fs.ReadFile(tr.staticFS, name)
	if err != nil {
		// fallback to absolute path if needed
		content, err = os.ReadFile(fullPath)
		if err != nil {
			return "", fmt.Errorf("read %s: %w", fullPath, err)
		}
	}

	expanded, err := tr.expandIncludes(content, filepath.Dir(fullPath))
	if err != nil {
		return "", err
	}

	if !tr.DevMode {
		tr.mu.Lock()
		tr.cache[name] = expanded
		tr.mu.Unlock()
	}
	return expanded, nil
}

func (tr *TemplateResolver) expandIncludes(content []byte, currentDir string) (string, error) {
	matches := includeRE.FindAllSubmatchIndex(content, -1)
	if matches == nil {
		return string(content), nil
	}

	var buf bytes.Buffer
	last := 0
	for _, m := range matches {
		start, end := m[0], m[1]
		pathStart, pathEnd := m[2], m[3]
		buf.Write(content[last:start])

		includeName := string(content[pathStart:pathEnd])
		includePath := filepath.Join(currentDir, includeName)
		includeBytes, err := fs.ReadFile(os.DirFS(filepath.Dir(includePath)), filepath.Base(includePath))
		if err != nil {
			return "", fmt.Errorf("include %s: %w", includePath, err)
		}
		expanded, err := tr.expandIncludes(includeBytes, filepath.Dir(includePath))
		if err != nil {
			return "", err
		}
		buf.WriteString(expanded)
		last = end
	}
	buf.Write(content[last:])
	out := buf.String()
	// Replace Flask static helpers with direct /static paths
	out = staticRE.ReplaceAllString(out, `/static/$1`)
	
	// Strip out Jinja2 blocks that we can't process
	out = stripJinja2Blocks(out)
	
	return out, nil
}

// stripJinja2Blocks removes Jinja2 conditionals, loops, and variable interpolation
// that we can't process, leaving the static HTML intact
func stripJinja2Blocks(content string) string {
	// Remove {% with ... %} blocks (including nested if/for/end blocks)
	content = removeJinjaBlock(content, `\{\%\s*with\s+[^%]+\%\}`, `\{\%\s*endwith\s*\%\}`, `\{\%\s*endwith\s*\%\}`)
	
	// Remove {% if ... %} blocks (handles {% else %} as well)
	content = removeJinjaBlock(content, `\{\%\s*if\s+[^%]+\%\}`, `\{\%\s*endif\s*\%\}`, `\{\%\s*(else|endif)\s*\%\}`)
	
	// Remove {% for ... %} blocks
	content = removeJinjaBlock(content, `\{\%\s*for\s+[^%]+\%\}`, `\{\%\s*endfor\s*\%\}`, `\{\%\s*endfor\s*\%\}`)
	
	// Remove standalone {% endif %}, {% endfor %}, {% endwith %}, {% else %}
	content = jinjaEndRE.ReplaceAllString(content, "")
	content = regexp.MustCompile(`\{\%\s*else\s*\%\}`).ReplaceAllString(content, "")
	
	// Remove variable interpolation {{ ... }}
	content = jinjaVarRE.ReplaceAllString(content, "")
	
	return content
}

// removeJinjaBlock removes a Jinja2 block from start pattern to end pattern
// Handles nested blocks by counting depth
// middlePattern is used to detect intermediate markers like {% else %} for if blocks
func removeJinjaBlock(content string, startPattern, endPattern, middlePattern string) string {
	startRE := regexp.MustCompile(startPattern)
	endRE := regexp.MustCompile(endPattern)
	middleRE := regexp.MustCompile(middlePattern)
	
	var result strings.Builder
	lastIndex := 0
	
	for {
		// Find next start
		startMatch := startRE.FindStringIndex(content[lastIndex:])
		if startMatch == nil {
			// No more blocks, append rest of content
			result.WriteString(content[lastIndex:])
			break
		}
		
		startPos := lastIndex + startMatch[0]
		endPos := lastIndex + startMatch[1]
		
		// Append content before the block
		result.WriteString(content[lastIndex:startPos])
		
		// Find matching end, handling nesting
		depth := 1
		searchPos := endPos
		
		for depth > 0 && searchPos < len(content) {
			// Check for nested starts
			if nestedStart := startRE.FindStringIndex(content[searchPos:]); nestedStart != nil {
				nestedStartPos := searchPos + nestedStart[0]
				// Check for matching end or middle before nested start
				if nestedEnd := endRE.FindStringIndex(content[searchPos:nestedStartPos]); nestedEnd != nil {
					searchPos = searchPos + nestedEnd[1]
					depth--
					continue
				}
				if nestedMiddle := middleRE.FindStringIndex(content[searchPos:nestedStartPos]); nestedMiddle != nil && depth == 1 {
					// Found else/middle marker at top level, skip to it
					searchPos = searchPos + nestedMiddle[1]
					continue
				}
				searchPos = nestedStartPos + (nestedStart[1] - nestedStart[0])
				depth++
				continue
			}
			
			// Check for middle markers (like {% else %}) at current depth
			if middleMatch := middleRE.FindStringIndex(content[searchPos:]); middleMatch != nil && depth == 1 {
				// Found else/middle at top level, skip to it
				searchPos = searchPos + middleMatch[1]
				continue
			}
			
			// Check for end
			if endMatch := endRE.FindStringIndex(content[searchPos:]); endMatch != nil {
				searchPos = searchPos + endMatch[1]
				depth--
			} else {
				// No matching end found, break
				break
			}
		}
		
		lastIndex = searchPos
	}
	
	return result.String()
}
