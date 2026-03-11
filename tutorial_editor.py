#!/usr/bin/env python3
"""
Tutorial Dialogue Script Editor
A Flask-based visual editor for the Divergence 2% tutorial dialogue script.
"""

import os
import re
import json
from flask import Flask, render_template_string, request, jsonify
from datetime import datetime

app = Flask(__name__)

# Paths to the markdown files (Chinese and English)
SCRIPT_PATH_CN = os.path.join(os.path.dirname(__file__), 'docs', 'tutorial_dialogue_script.md')
SCRIPT_PATH_EN = os.path.join(os.path.dirname(__file__), 'docs', 'tutorial_dialogue_script_EN.md')

# Default to Chinese
SCRIPT_PATH = SCRIPT_PATH_CN

# HTML template for the editor
EDITOR_HTML = '''
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>📜 Tutorial Dialogue Editor</title>
    <link href="https://fonts.googleapis.com/css2?family=Noto+Sans+SC:wght@300;400;500;700&family=JetBrains+Mono:wght@400;500&display=swap" rel="stylesheet">
    <style>
        :root {
            --bg-dark: #0d1117;
            --bg-card: #161b22;
            --bg-hover: #21262d;
            --border: #30363d;
            --text-primary: #e6edf3;
            --text-secondary: #8b949e;
            --text-muted: #6e7681;
            --accent-blue: #58a6ff;
            --accent-purple: #bc8cff;
            --accent-green: #3fb950;
            --accent-yellow: #d29922;
            --accent-red: #f85149;
            --accent-orange: #db6d28;
            --cao-cao: #ffa657;
            --user: #79c0ff;
            --choice: #a5d6ff;
            --highlight: #7ee787;
            --action: #d2a8ff;
        }

        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: 'Noto Sans SC', sans-serif;
            background: var(--bg-dark);
            color: var(--text-primary);
            min-height: 100vh;
            line-height: 1.6;
        }

        /* Header */
        .header {
            background: linear-gradient(135deg, #1a1f35 0%, #0d1117 100%);
            border-bottom: 1px solid var(--border);
            padding: 1.5rem 2rem;
            position: sticky;
            top: 0;
            z-index: 100;
            backdrop-filter: blur(10px);
        }

        .header-content {
            max-width: 1600px;
            margin: 0 auto;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }

        .title {
            font-size: 1.75rem;
            font-weight: 700;
            background: linear-gradient(135deg, var(--cao-cao), var(--accent-purple));
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            display: flex;
            align-items: center;
            gap: 0.75rem;
        }

        .title-emoji {
            font-size: 2rem;
            -webkit-text-fill-color: initial;
        }

        .subtitle {
            color: var(--text-secondary);
            font-size: 0.9rem;
            margin-top: 0.25rem;
        }

        .header-actions {
            display: flex;
            gap: 1rem;
            align-items: center;
        }

        .btn {
            padding: 0.65rem 1.25rem;
            border-radius: 8px;
            border: none;
            font-family: inherit;
            font-size: 0.9rem;
            font-weight: 500;
            cursor: pointer;
            transition: all 0.2s ease;
            display: flex;
            align-items: center;
            gap: 0.5rem;
        }

        .btn-primary {
            background: linear-gradient(135deg, var(--accent-green), #238636);
            color: white;
        }

        .btn-primary:hover {
            transform: translateY(-1px);
            box-shadow: 0 4px 12px rgba(63, 185, 80, 0.3);
        }

        .btn-secondary {
            background: var(--bg-hover);
            color: var(--text-primary);
            border: 1px solid var(--border);
        }

        .btn-secondary:hover {
            background: var(--border);
        }

        .save-status {
            font-size: 0.85rem;
            color: var(--text-muted);
            display: flex;
            align-items: center;
            gap: 0.5rem;
        }

        .save-status.saved {
            color: var(--accent-green);
        }

        .save-status.unsaved {
            color: var(--accent-yellow);
        }

        /* Main layout */
        .main-container {
            max-width: 1600px;
            margin: 0 auto;
            padding: 2rem;
            display: grid;
            grid-template-columns: 280px 1fr;
            gap: 2rem;
        }

        /* Sidebar */
        .sidebar {
            position: sticky;
            top: 100px;
            height: fit-content;
            max-height: calc(100vh - 120px);
            overflow-y: auto;
        }

        .sidebar-card {
            background: var(--bg-card);
            border: 1px solid var(--border);
            border-radius: 12px;
            padding: 1.25rem;
        }

        .sidebar-title {
            font-size: 0.8rem;
            text-transform: uppercase;
            letter-spacing: 0.5px;
            color: var(--text-muted);
            margin-bottom: 1rem;
            font-weight: 500;
        }

        .section-nav {
            list-style: none;
        }

        .section-nav-item {
            padding: 0.6rem 0.75rem;
            border-radius: 6px;
            cursor: pointer;
            transition: all 0.15s ease;
            font-size: 0.9rem;
            color: var(--text-secondary);
            display: flex;
            align-items: center;
            gap: 0.5rem;
            margin-bottom: 0.25rem;
        }

        .section-nav-item:hover {
            background: var(--bg-hover);
            color: var(--text-primary);
        }

        .section-nav-item.active {
            background: rgba(88, 166, 255, 0.15);
            color: var(--accent-blue);
        }

        .section-num {
            font-family: 'JetBrains Mono', monospace;
            font-size: 0.75rem;
            background: var(--bg-hover);
            padding: 0.15rem 0.4rem;
            border-radius: 4px;
            min-width: 1.5rem;
            text-align: center;
        }

        /* Editor area */
        .editor-area {
            display: flex;
            flex-direction: column;
            gap: 1.5rem;
        }

        .section-card {
            background: var(--bg-card);
            border: 1px solid var(--border);
            border-radius: 12px;
            overflow: hidden;
            transition: border-color 0.2s ease;
        }

        .section-card:hover {
            border-color: var(--border);
        }

        .section-card.editing {
            border-color: var(--accent-blue);
            box-shadow: 0 0 0 1px var(--accent-blue);
        }

        .section-header {
            background: linear-gradient(135deg, rgba(88, 166, 255, 0.1), rgba(188, 140, 255, 0.05));
            padding: 1rem 1.25rem;
            display: flex;
            justify-content: space-between;
            align-items: center;
            border-bottom: 1px solid var(--border);
        }

        .section-title {
            font-size: 1.1rem;
            font-weight: 600;
            display: flex;
            align-items: center;
            gap: 0.75rem;
        }

        .section-badge {
            font-family: 'JetBrains Mono', monospace;
            font-size: 0.75rem;
            background: var(--accent-purple);
            color: white;
            padding: 0.2rem 0.5rem;
            border-radius: 4px;
        }

        .section-desc {
            font-size: 0.85rem;
            color: var(--text-secondary);
            font-style: italic;
        }

        .section-actions {
            display: flex;
            gap: 0.5rem;
        }

        .btn-icon {
            padding: 0.5rem;
            background: transparent;
            border: 1px solid transparent;
            border-radius: 6px;
            color: var(--text-secondary);
            cursor: pointer;
            transition: all 0.15s ease;
        }

        .btn-icon:hover {
            background: var(--bg-hover);
            color: var(--text-primary);
            border-color: var(--border);
        }

        .section-content {
            padding: 1.25rem;
        }

        /* Dialogue display */
        .dialogue-preview {
            font-family: 'Noto Sans SC', sans-serif;
            line-height: 1.8;
        }

        .dialogue-line {
            padding: 0.5rem 0.75rem;
            margin: 0.25rem 0;
            border-radius: 8px;
            background: var(--bg-dark);
        }

        .dialogue-line.cao-cao {
            border-left: 3px solid var(--cao-cao);
        }

        .dialogue-line.user {
            border-left: 3px solid var(--user);
        }

        .dialogue-line.choice {
            border-left: 3px solid var(--choice);
            background: rgba(165, 214, 255, 0.05);
        }

        .dialogue-line.highlight {
            border-left: 3px solid var(--highlight);
            background: rgba(126, 231, 135, 0.05);
        }

        .dialogue-line.action {
            border-left: 3px solid var(--action);
            background: rgba(210, 168, 255, 0.05);
        }

        .speaker {
            font-weight: 600;
            margin-right: 0.5rem;
        }

        .speaker.cao-cao {
            color: var(--cao-cao);
        }

        .speaker.user {
            color: var(--user);
        }

        /* Editor textarea */
        .editor-textarea {
            width: 100%;
            min-height: 300px;
            background: var(--bg-dark);
            border: 1px solid var(--border);
            border-radius: 8px;
            color: var(--text-primary);
            font-family: 'JetBrains Mono', monospace;
            font-size: 0.9rem;
            line-height: 1.7;
            padding: 1rem;
            resize: vertical;
        }

        .editor-textarea:focus {
            outline: none;
            border-color: var(--accent-blue);
        }

        /* Mode toggle */
        .mode-toggle {
            display: flex;
            background: var(--bg-dark);
            border-radius: 8px;
            padding: 0.25rem;
            border: 1px solid var(--border);
        }

        .mode-btn {
            padding: 0.4rem 0.75rem;
            border: none;
            background: transparent;
            color: var(--text-secondary);
            font-size: 0.85rem;
            cursor: pointer;
            border-radius: 6px;
            transition: all 0.15s ease;
        }

        .mode-btn.active {
            background: var(--bg-hover);
            color: var(--text-primary);
        }

        /* Raw markdown view */
        .raw-editor {
            padding: 1.5rem;
        }

        .raw-textarea {
            width: 100%;
            min-height: calc(100vh - 200px);
            background: var(--bg-dark);
            border: 1px solid var(--border);
            border-radius: 8px;
            color: var(--text-primary);
            font-family: 'JetBrains Mono', monospace;
            font-size: 0.9rem;
            line-height: 1.7;
            padding: 1.25rem;
            resize: vertical;
        }

        .raw-textarea:focus {
            outline: none;
            border-color: var(--accent-blue);
        }

        /* Toast notification */
        .toast {
            position: fixed;
            bottom: 2rem;
            right: 2rem;
            padding: 1rem 1.5rem;
            background: var(--bg-card);
            border: 1px solid var(--border);
            border-radius: 10px;
            box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
            display: flex;
            align-items: center;
            gap: 0.75rem;
            transform: translateY(100px);
            opacity: 0;
            transition: all 0.3s ease;
            z-index: 1000;
        }

        .toast.show {
            transform: translateY(0);
            opacity: 1;
        }

        .toast.success {
            border-color: var(--accent-green);
        }

        .toast.error {
            border-color: var(--accent-red);
        }

        /* Legend */
        .legend {
            margin-top: 1.5rem;
            padding-top: 1rem;
            border-top: 1px solid var(--border);
        }

        .legend-item {
            display: flex;
            align-items: center;
            gap: 0.5rem;
            font-size: 0.8rem;
            color: var(--text-secondary);
            margin-bottom: 0.5rem;
        }

        .legend-color {
            width: 12px;
            height: 12px;
            border-radius: 3px;
        }

        /* Scrollbar */
        ::-webkit-scrollbar {
            width: 8px;
            height: 8px;
        }

        ::-webkit-scrollbar-track {
            background: var(--bg-dark);
        }

        ::-webkit-scrollbar-thumb {
            background: var(--border);
            border-radius: 4px;
        }

        ::-webkit-scrollbar-thumb:hover {
            background: var(--text-muted);
        }

        /* Loading state */
        .loading {
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 200px;
            color: var(--text-muted);
        }

        .spinner {
            width: 24px;
            height: 24px;
            border: 2px solid var(--border);
            border-top-color: var(--accent-blue);
            border-radius: 50%;
            animation: spin 0.8s linear infinite;
        }

        @keyframes spin {
            to { transform: rotate(360deg); }
        }

        /* Metadata section */
        .metadata-section {
            background: var(--bg-card);
            border: 1px solid var(--border);
            border-radius: 12px;
            padding: 1.25rem;
            margin-bottom: 1.5rem;
        }

        .metadata-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 1rem;
        }

        .metadata-item {
            display: flex;
            flex-direction: column;
            gap: 0.25rem;
        }

        .metadata-label {
            font-size: 0.75rem;
            text-transform: uppercase;
            letter-spacing: 0.5px;
            color: var(--text-muted);
        }

        .metadata-value {
            font-size: 0.95rem;
            color: var(--text-primary);
        }
    </style>
</head>
<body>
    <header class="header">
        <div class="header-content">
            <div>
                <h1 class="title">
                    <span class="title-emoji">📜</span>
                    Tutorial Dialogue Editor
                </h1>
                <p class="subtitle">Divergence 2% — Cao Cao's Guide to Writing Characters</p>
            </div>
            <div class="header-actions">
                <span class="save-status" id="saveStatus">
                    <span>●</span> Saved
                </span>
                <div class="mode-toggle">
                    <button class="mode-btn" data-lang="cn" onclick="setLanguage('cn')">中文</button>
                    <button class="mode-btn" data-lang="en" onclick="setLanguage('en')">English</button>
                </div>
                <div class="mode-toggle" style="margin-left: 0.5rem;">
                    <button class="mode-btn active" data-mode="visual" onclick="setMode('visual')">Visual</button>
                    <button class="mode-btn" data-mode="raw" onclick="setMode('raw')">Raw Markdown</button>
                </div>
                <button class="btn btn-secondary" onclick="reloadContent()">
                    ↻ Reload
                </button>
                <button class="btn btn-primary" onclick="saveContent()">
                    💾 Save Changes
                </button>
            </div>
        </div>
    </header>

    <main class="main-container" id="visualView">
        <aside class="sidebar">
            <div class="sidebar-card">
                <h3 class="sidebar-title">Sections</h3>
                <ul class="section-nav" id="sectionNav">
                    <!-- Populated by JS -->
                </ul>
                <div class="legend">
                    <h3 class="sidebar-title">Legend</h3>
                    <div class="legend-item">
                        <div class="legend-color" style="background: var(--cao-cao)"></div>
                        <span id="legend-caocao">曹操 (Cao Cao)</span>
                    </div>
                    <div class="legend-item">
                        <div class="legend-color" style="background: var(--user)"></div>
                        <span id="legend-user">用户 (User)</span>
                    </div>
                    <div class="legend-item">
                        <div class="legend-color" style="background: var(--choice)"></div>
                        <span>→ CHOICE</span>
                    </div>
                    <div class="legend-item">
                        <div class="legend-color" style="background: var(--highlight)"></div>
                        <span>[ HIGHLIGHT ]</span>
                    </div>
                    <div class="legend-item">
                        <div class="legend-color" style="background: var(--action)"></div>
                        <span>→ [Action]</span>
                    </div>
                </div>
            </div>
        </aside>
        
        <div class="editor-area" id="editorArea">
            <div class="loading">
                <div class="spinner"></div>
            </div>
        </div>
    </main>

    <div class="raw-editor" id="rawView" style="display: none;">
        <textarea class="raw-textarea" id="rawEditor" spellcheck="false"></textarea>
    </div>

    <div class="toast" id="toast">
        <span id="toastIcon">✓</span>
        <span id="toastMessage">Saved successfully!</span>
    </div>

    <script>
        let sections = [];
        let rawContent = '';
        let hasChanges = false;
        let currentMode = 'visual';
        let currentLanguage = 'cn';

        // Initialize
        document.addEventListener('DOMContentLoaded', () => {
            // Check URL for language param
            const urlParams = new URLSearchParams(window.location.search);
            currentLanguage = urlParams.get('lang') || 'cn';
            updateLanguageButtons();
            loadContent();
        });

        function updateLanguageButtons() {
            document.querySelectorAll('[data-lang]').forEach(btn => {
                btn.classList.toggle('active', btn.dataset.lang === currentLanguage);
            });
        }

        function setLanguage(lang) {
            if (hasChanges && !confirm('You have unsaved changes. Switch language anyway?')) {
                return;
            }
            currentLanguage = lang;
            updateLanguageButtons();
            loadContent();
            // Update URL without reload
            const url = new URL(window.location);
            url.searchParams.set('lang', lang);
            window.history.pushState({}, '', url);
        }

        async function loadContent() {
            try {
                const response = await fetch(`/api/script?lang=${currentLanguage}`);
                const data = await response.json();
                rawContent = data.content;
                parseAndRender(rawContent);
                document.getElementById('rawEditor').value = rawContent;
                updateSaveStatus(false);
                hasChanges = false;
            } catch (error) {
                showToast('Failed to load content', 'error');
                console.error(error);
            }
        }

        async function reloadContent() {
            if (hasChanges && !confirm('You have unsaved changes. Reload anyway?')) {
                return;
            }
            await loadContent();
            showToast('Content reloaded', 'success');
        }

        function parseAndRender(content) {
            sections = parseSections(content);
            renderSections(sections);
            renderNav(sections);
        }

        function parseSections(content) {
            const sectionRegex = /## SECTION (\d+(?:\.\d+)?): ([^\\n]+)/g;
            const parts = [];
            let lastIndex = 0;
            let match;

            // Extract header/metadata first
            const headerMatch = content.match(/^([\s\S]*?)(?=## SECTION)/);
            if (headerMatch) {
                parts.push({
                    type: 'header',
                    title: 'Document Header',
                    content: headerMatch[1],
                    startIndex: 0,
                    endIndex: headerMatch[1].length
                });
                lastIndex = headerMatch[1].length;
            }

            while ((match = sectionRegex.exec(content)) !== null) {
                if (match.index > lastIndex) {
                    // There's content between sections (shouldn't happen normally)
                }
                
                const nextMatch = sectionRegex.exec(content);
                const endIndex = nextMatch ? nextMatch.index : content.length;
                sectionRegex.lastIndex = match.index + match[0].length; // Reset for next iteration
                
                // Find the actual end by looking for next section header or end of content
                const sectionEnd = content.indexOf('\\n## SECTION', match.index + 1);
                const actualEnd = sectionEnd === -1 ? 
                    (content.indexOf('\\n## 📝', match.index + 1) !== -1 ? 
                        content.indexOf('\\n## 📝', match.index + 1) : content.length) 
                    : sectionEnd;
                
                parts.push({
                    type: 'section',
                    number: match[1],
                    title: match[2].trim(),
                    content: content.substring(match.index, actualEnd),
                    startIndex: match.index,
                    endIndex: actualEnd
                });
                
                lastIndex = actualEnd;
            }

            // Check for implementation notes at the end
            const notesMatch = content.match(/(## 📝 Implementation Notes[\s\S]*)/);
            if (notesMatch) {
                parts.push({
                    type: 'notes',
                    title: 'Implementation Notes',
                    content: notesMatch[1],
                    startIndex: content.indexOf(notesMatch[1]),
                    endIndex: content.length
                });
            }

            return parts;
        }

        function renderSections(sections) {
            const container = document.getElementById('editorArea');
            container.innerHTML = '';

            sections.forEach((section, index) => {
                const card = document.createElement('div');
                card.className = 'section-card';
                card.id = `section-${index}`;

                const badge = section.type === 'header' ? 'META' : 
                              section.type === 'notes' ? 'IMPL' : 
                              `S${section.number}`;

                card.innerHTML = `
                    <div class="section-header">
                        <div class="section-title">
                            <span class="section-badge">${badge}</span>
                            ${section.title}
                        </div>
                        <div class="section-actions">
                            <button class="btn-icon" onclick="toggleEdit(${index})" title="Edit">✏️</button>
                            <button class="btn-icon" onclick="expandSection(${index})" title="Expand">⬜</button>
                        </div>
                    </div>
                    <div class="section-content">
                        <div class="dialogue-preview" id="preview-${index}">
                            ${renderDialoguePreview(section.content)}
                        </div>
                        <textarea 
                            class="editor-textarea" 
                            id="editor-${index}" 
                            style="display: none;"
                            oninput="onSectionEdit(${index})"
                        >${escapeHtml(section.content)}</textarea>
                    </div>
                `;

                container.appendChild(card);
            });
        }

        function renderDialoguePreview(content) {
            // Extract dialogue blocks
            const lines = content.split('\\n');
            let html = '';
            let inCodeBlock = false;

            for (const line of lines) {
                if (line.trim().startsWith('```')) {
                    inCodeBlock = !inCodeBlock;
                    continue;
                }

                if (!inCodeBlock && !line.startsWith('##') && line.trim()) {
                    if (line.startsWith('曹操:') || line.startsWith('曹操：') || line.startsWith('Cao Cao:')) {
                        const speakerLabel = currentLanguage === 'en' ? 'Cao Cao:' : '曹操:';
                        html += `<div class="dialogue-line cao-cao">
                            <span class="speaker cao-cao">${speakerLabel}</span>
                            ${escapeHtml(line.replace(/^(曹操[:：]|Cao Cao:)/, '').trim())}
                        </div>`;
                    } else if (line.startsWith('用户:') || line.startsWith('用户：') || line.startsWith('User:')) {
                        const speakerLabel = currentLanguage === 'en' ? 'User:' : '用户:';
                        html += `<div class="dialogue-line user">
                            <span class="speaker user">${speakerLabel}</span>
                            ${escapeHtml(line.replace(/^(用户[:：]|User:)/, '').trim())}
                        </div>`;
                    } else if (line.trim().startsWith('→ CHOICE')) {
                        html += `<div class="dialogue-line choice">${escapeHtml(line)}</div>`;
                    } else if (line.includes('[ HIGHLIGHT')) {
                        html += `<div class="dialogue-line highlight">${escapeHtml(line)}</div>`;
                    } else if (line.trim().startsWith('→ [')) {
                        html += `<div class="dialogue-line action">${escapeHtml(line)}</div>`;
                    } else if (line.trim().startsWith('-')) {
                        html += `<div class="dialogue-line choice" style="margin-left: 1rem;">${escapeHtml(line)}</div>`;
                    } else if (line.trim().startsWith('*')) {
                        html += `<div style="color: var(--text-muted); font-style: italic; padding: 0.5rem 0;">${escapeHtml(line)}</div>`;
                    } else {
                        html += `<div style="padding: 0.25rem 0; color: var(--text-secondary);">${escapeHtml(line)}</div>`;
                    }
                } else if (inCodeBlock) {
                    if (line.startsWith('曹操:') || line.startsWith('曹操：') || line.startsWith('Cao Cao:')) {
                        const speakerLabel = currentLanguage === 'en' ? 'Cao Cao:' : '曹操:';
                        html += `<div class="dialogue-line cao-cao">
                            <span class="speaker cao-cao">${speakerLabel}</span>
                            ${escapeHtml(line.replace(/^(曹操[:：]|Cao Cao:)/, '').trim())}
                        </div>`;
                    } else if (line.startsWith('用户:') || line.startsWith('用户：') || line.startsWith('User:')) {
                        const speakerLabel = currentLanguage === 'en' ? 'User:' : '用户:';
                        html += `<div class="dialogue-line user">
                            <span class="speaker user">${speakerLabel}</span>
                            ${escapeHtml(line.replace(/^(用户[:：]|User:)/, '').trim())}
                        </div>`;
                    } else if (line.trim().startsWith('→ CHOICE')) {
                        html += `<div class="dialogue-line choice">${escapeHtml(line)}</div>`;
                    } else if (line.includes('[ HIGHLIGHT')) {
                        html += `<div class="dialogue-line highlight">${escapeHtml(line)}</div>`;
                    } else if (line.trim().startsWith('→ [')) {
                        html += `<div class="dialogue-line action">${escapeHtml(line)}</div>`;
                    } else if (line.trim().startsWith('[Click')) {
                        html += `<div class="dialogue-line action">${escapeHtml(line)}</div>`;
                    } else if (line.trim().startsWith('-')) {
                        html += `<div class="dialogue-line choice" style="margin-left: 1rem;">${escapeHtml(line)}</div>`;
                    } else if (line.trim()) {
                        html += `<div style="padding: 0.25rem 0; color: var(--text-secondary);">${escapeHtml(line)}</div>`;
                    }
                } else if (line.startsWith('##')) {
                    // Skip headers in preview, they're in the card header
                } else if (line.trim()) {
                    html += `<div style="padding: 0.25rem 0;">${escapeHtml(line)}</div>`;
                }
            }

            return html || '<div style="color: var(--text-muted);">No dialogue content</div>';
        }

        function renderNav(sections) {
            const nav = document.getElementById('sectionNav');
            nav.innerHTML = sections.map((section, index) => `
                <li class="section-nav-item" onclick="scrollToSection(${index})">
                    <span class="section-num">${section.type === 'header' ? 'H' : section.type === 'notes' ? 'N' : section.number}</span>
                    ${truncate(section.title, 20)}
                </li>
            `).join('');
        }

        function scrollToSection(index) {
            const section = document.getElementById(`section-${index}`);
            section.scrollIntoView({ behavior: 'smooth', block: 'start' });
            
            // Update nav active state
            document.querySelectorAll('.section-nav-item').forEach((item, i) => {
                item.classList.toggle('active', i === index);
            });
        }

        function toggleEdit(index) {
            const preview = document.getElementById(`preview-${index}`);
            const editor = document.getElementById(`editor-${index}`);
            const card = document.getElementById(`section-${index}`);
            
            const isEditing = preview.style.display === 'none';
            
            if (isEditing) {
                // Switch to preview mode
                preview.style.display = 'block';
                editor.style.display = 'none';
                card.classList.remove('editing');
                
                // Update preview with edited content
                sections[index].content = editor.value;
                preview.innerHTML = renderDialoguePreview(editor.value);
            } else {
                // Switch to edit mode
                preview.style.display = 'none';
                editor.style.display = 'block';
                card.classList.add('editing');
                editor.focus();
            }
        }

        function expandSection(index) {
            const card = document.getElementById(`section-${index}`);
            card.scrollIntoView({ behavior: 'smooth', block: 'start' });
        }

        function onSectionEdit(index) {
            const editor = document.getElementById(`editor-${index}`);
            sections[index].content = editor.value;
            hasChanges = true;
            updateSaveStatus(true);
        }

        async function saveContent() {
            let contentToSave;
            
            if (currentMode === 'raw') {
                contentToSave = document.getElementById('rawEditor').value;
            } else {
                // Reconstruct full content from sections
                contentToSave = sections.map(s => s.content).join('\\n');
            }

            try {
                const response = await fetch('/api/script', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ content: contentToSave, language: currentLanguage })
                });

                if (response.ok) {
                    rawContent = contentToSave;
                    hasChanges = false;
                    updateSaveStatus(false);
                    showToast('Saved successfully!', 'success');
                    
                    // If we were in raw mode, re-parse for visual mode
                    if (currentMode === 'raw') {
                        parseAndRender(contentToSave);
                    }
                } else {
                    throw new Error('Save failed');
                }
            } catch (error) {
                showToast('Failed to save', 'error');
                console.error(error);
            }
        }

        function setMode(mode) {
            currentMode = mode;
            
            document.querySelectorAll('.mode-btn').forEach(btn => {
                btn.classList.toggle('active', btn.dataset.mode === mode);
            });

            if (mode === 'visual') {
                document.getElementById('visualView').style.display = 'grid';
                document.getElementById('rawView').style.display = 'none';
                
                // Sync raw content to sections if there were changes
                if (hasChanges) {
                    const rawEditor = document.getElementById('rawEditor');
                    parseAndRender(rawEditor.value);
                }
            } else {
                document.getElementById('visualView').style.display = 'none';
                document.getElementById('rawView').style.display = 'block';
                
                // Sync sections to raw editor if there were changes
                if (hasChanges) {
                    document.getElementById('rawEditor').value = sections.map(s => s.content).join('\\n');
                }
            }
        }

        function updateSaveStatus(unsaved) {
            const status = document.getElementById('saveStatus');
            if (unsaved) {
                status.className = 'save-status unsaved';
                status.innerHTML = '<span>●</span> Unsaved changes';
            } else {
                status.className = 'save-status saved';
                status.innerHTML = '<span>●</span> Saved';
            }
        }

        function showToast(message, type = 'success') {
            const toast = document.getElementById('toast');
            const icon = document.getElementById('toastIcon');
            const msg = document.getElementById('toastMessage');
            
            toast.className = `toast ${type}`;
            icon.textContent = type === 'success' ? '✓' : '✕';
            msg.textContent = message;
            
            toast.classList.add('show');
            setTimeout(() => toast.classList.remove('show'), 3000);
        }

        function escapeHtml(text) {
            const div = document.createElement('div');
            div.textContent = text;
            return div.innerHTML;
        }

        function truncate(str, len) {
            return str.length > len ? str.substring(0, len) + '...' : str;
        }

        // Handle raw editor changes
        document.getElementById('rawEditor').addEventListener('input', () => {
            hasChanges = true;
            updateSaveStatus(true);
        });

        // Warn before leaving with unsaved changes
        window.addEventListener('beforeunload', (e) => {
            if (hasChanges) {
                e.preventDefault();
                e.returnValue = '';
            }
        });

        // Keyboard shortcuts
        document.addEventListener('keydown', (e) => {
            if ((e.ctrlKey || e.metaKey) && e.key === 's') {
                e.preventDefault();
                saveContent();
            }
        });
    </script>
</body>
</html>
'''


@app.route('/')
def index():
    """Serve the main editor page."""
    return render_template_string(EDITOR_HTML)


@app.route('/api/script', methods=['GET'])
def get_script():
    """Get the tutorial dialogue script content."""
    lang = request.args.get('lang', 'cn')
    script_path = SCRIPT_PATH_EN if lang == 'en' else SCRIPT_PATH_CN
    try:
        with open(script_path, 'r', encoding='utf-8') as f:
            content = f.read()
        return jsonify({'content': content, 'language': lang})
    except FileNotFoundError:
        return jsonify({'error': 'Script file not found'}), 404
    except Exception as e:
        return jsonify({'error': str(e)}), 500


@app.route('/api/script', methods=['POST'])
def save_script():
    """Save the tutorial dialogue script content."""
    try:
        data = request.get_json()
        content = data.get('content', '')
        lang = data.get('language', 'cn')
        script_path = SCRIPT_PATH_EN if lang == 'en' else SCRIPT_PATH_CN
        
        # Create a backup before saving
        backup_dir = os.path.join(os.path.dirname(script_path), 'backups')
        os.makedirs(backup_dir, exist_ok=True)
        
        if os.path.exists(script_path):
            timestamp = datetime.now().strftime('%Y%m%d_%H%M%S')
            suffix = '_EN' if lang == 'en' else ''
            backup_path = os.path.join(backup_dir, f'tutorial_dialogue_script{suffix}_{timestamp}.md')
            with open(script_path, 'r', encoding='utf-8') as f:
                with open(backup_path, 'w', encoding='utf-8') as bf:
                    bf.write(f.read())
        
        # Save the new content
        with open(script_path, 'w', encoding='utf-8') as f:
            f.write(content)
        
        return jsonify({'success': True, 'message': 'Script saved successfully', 'language': lang})
    except Exception as e:
        return jsonify({'error': str(e)}), 500


@app.route('/api/script/backup', methods=['GET'])
def list_backups():
    """List available backups."""
    try:
        backup_dir = os.path.join(os.path.dirname(SCRIPT_PATH), 'backups')
        if not os.path.exists(backup_dir):
            return jsonify({'backups': []})
        
        backups = []
        for filename in os.listdir(backup_dir):
            if filename.startswith('tutorial_dialogue_script_') and filename.endswith('.md'):
                filepath = os.path.join(backup_dir, filename)
                backups.append({
                    'filename': filename,
                    'modified': os.path.getmtime(filepath)
                })
        
        backups.sort(key=lambda x: x['modified'], reverse=True)
        return jsonify({'backups': backups})
    except Exception as e:
        return jsonify({'error': str(e)}), 500


if __name__ == '__main__':
    print(f"🚀 Tutorial Dialogue Editor starting...")
    print(f"📁 Script path: {SCRIPT_PATH}")
    print(f"🌐 Open http://localhost:5050 in your browser")
    app.run(host='0.0.0.0', port=5050, debug=True)
