#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
写手日记系统 - 简单的思考记录工具
"""

import os
import json
import time
from datetime import datetime
from typing import List, Dict, Any, Optional

class WriterJournal:
    """写手日记管理器"""
    
    def __init__(self, journal_file: str = "writer_journal.txt"):
        self.journal_file = journal_file
        self.entries_file = journal_file.replace('.txt', '_entries.json')
        self.ensure_files_exist()
    
    def ensure_files_exist(self):
        """确保日记文件存在"""
        if not os.path.exists(self.journal_file):
            with open(self.journal_file, 'w', encoding='utf-8') as f:
                f.write("# 写手日记\n\n")
        
        if not os.path.exists(self.entries_file):
            with open(self.entries_file, 'w', encoding='utf-8') as f:
                json.dump([], f, ensure_ascii=False, indent=2)
    
    def add_entry(self, title: str, content: str, tags: Optional[List[str]] = None, mood: Optional[str] = None) -> Dict[str, Any]:
        """添加日记条目"""
        timestamp = time.time()
        datetime_str = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        
        entry = {
            "id": int(timestamp * 1000),  # 使用毫秒作为ID
            "timestamp": timestamp,
            "datetime": datetime_str,
            "title": title.strip(),
            "content": content.strip(),
            "tags": tags or [],
            "mood": mood,
            "word_count": len(content.strip().split()) if content.strip() else 0
        }
        
        # 保存到JSON文件（用于前端显示）
        try:
            with open(self.entries_file, 'r', encoding='utf-8') as f:
                entries = json.load(f)
        except:
            entries = []
        
        entries.append(entry)
        entries.sort(key=lambda x: x['timestamp'], reverse=True)  # 最新的在前
        
        with open(self.entries_file, 'w', encoding='utf-8') as f:
            json.dump(entries, f, ensure_ascii=False, indent=2)
        
        # 同时追加到文本文件（可读性强）
        self._append_to_text_file(entry)
        
        return entry
    
    def _append_to_text_file(self, entry: Dict[str, Any]):
        """追加条目到文本文件"""
        with open(self.journal_file, 'a', encoding='utf-8') as f:
            f.write(f"\n---\n")
            f.write(f"**{entry['datetime']}** - {entry['title']}\n\n")
            
            if entry['tags']:
                f.write(f"🏷️ 标签: {', '.join(entry['tags'])}\n")
            
            if entry['mood']:
                f.write(f"😊 心情: {entry['mood']}\n")
            
            f.write(f"📝 字数: {entry['word_count']}\n\n")
            f.write(f"{entry['content']}\n")
    
    def get_entries(self, limit: int = 50, tag: Optional[str] = None, search: Optional[str] = None) -> List[Dict[str, Any]]:
        """获取日记条目"""
        try:
            with open(self.entries_file, 'r', encoding='utf-8') as f:
                entries = json.load(f)
        except:
            return []
        
        # 应用过滤器
        filtered_entries = []
        for entry in entries:
            # 标签过滤
            if tag and tag not in entry.get('tags', []):
                continue
            
            # 搜索过滤
            if search:
                search_text = f"{entry['title']} {entry['content']}".lower()
                if search.lower() not in search_text:
                    continue
            
            filtered_entries.append(entry)
        
        return filtered_entries[:limit]
    
    def update_entry(self, entry_id: int, title: Optional[str] = None, content: Optional[str] = None, 
                    tags: Optional[List[str]] = None, mood: Optional[str] = None) -> Optional[Dict[str, Any]]:
        """更新日记条目"""
        try:
            with open(self.entries_file, 'r', encoding='utf-8') as f:
                entries = json.load(f)
        except:
            return None
        
        # 查找条目
        for entry in entries:
            if entry['id'] == entry_id:
                if title is not None:
                    entry['title'] = title.strip()
                if content is not None:
                    entry['content'] = content.strip()
                    entry['word_count'] = len(content.strip().split()) if content.strip() else 0
                if tags is not None:
                    entry['tags'] = tags
                if mood is not None:
                    entry['mood'] = mood
                
                entry['updated_at'] = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
                
                # 保存更新
                with open(self.entries_file, 'w', encoding='utf-8') as f:
                    json.dump(entries, f, ensure_ascii=False, indent=2)
                
                # 重新生成文本文件
                self._regenerate_text_file(entries)
                
                return entry
        
        return None
    
    def delete_entry(self, entry_id: int) -> bool:
        """删除日记条目"""
        try:
            with open(self.entries_file, 'r', encoding='utf-8') as f:
                entries = json.load(f)
        except:
            return False
        
        # 删除条目
        original_count = len(entries)
        entries = [entry for entry in entries if entry['id'] != entry_id]
        
        if len(entries) < original_count:
            # 保存更新
            with open(self.entries_file, 'w', encoding='utf-8') as f:
                json.dump(entries, f, ensure_ascii=False, indent=2)
            
            # 重新生成文本文件
            self._regenerate_text_file(entries)
            return True
        
        return False
    
    def _regenerate_text_file(self, entries: List[Dict[str, Any]]):
        """重新生成文本文件"""
        with open(self.journal_file, 'w', encoding='utf-8') as f:
            f.write("# 写手日记\n\n")
            f.write(f"📅 最后更新: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n")
            f.write(f"📊 总条目数: {len(entries)}\n\n")
            
            for entry in sorted(entries, key=lambda x: x['timestamp'], reverse=True):
                f.write(f"\n---\n")
                f.write(f"**{entry['datetime']}** - {entry['title']}\n\n")
                
                if entry.get('tags'):
                    f.write(f"🏷️ 标签: {', '.join(entry['tags'])}\n")
                
                if entry.get('mood'):
                    f.write(f"😊 心情: {entry['mood']}\n")
                
                f.write(f"📝 字数: {entry['word_count']}\n\n")
                f.write(f"{entry['content']}\n")
    
    def get_statistics(self) -> Dict[str, Any]:
        """获取统计信息"""
        try:
            with open(self.entries_file, 'r', encoding='utf-8') as f:
                entries = json.load(f)
        except:
            entries = []
        
        # 统计
        stats = {
            "total_entries": len(entries),
            "total_words": sum(entry.get('word_count', 0) for entry in entries),
            "tags": {},
            "moods": {},
            "recent_activity": []
        }
        
        # 统计标签和心情
        for entry in entries:
            for tag in entry.get('tags', []):
                stats['tags'][tag] = stats['tags'].get(tag, 0) + 1
            
            mood = entry.get('mood')
            if mood:
                stats['moods'][mood] = stats['moods'].get(mood, 0) + 1
        
        # 最近活动（按日期分组）
        daily_counts = {}
        for entry in entries:
            date = entry['datetime'].split(' ')[0]
            daily_counts[date] = daily_counts.get(date, 0) + 1
        
        stats['recent_activity'] = sorted(daily_counts.items(), reverse=True)[:7]
        
        return stats

# 全局实例
writer_journal = WriterJournal() 