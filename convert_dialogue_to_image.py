import json
from PIL import Image, ImageDraw, ImageFont, Image as PILImage
import os
import json as _json
import re

# 配置参数
IMG_WIDTH = 900
MARGIN = 40
PADDING = 20
BUBBLE_PADDING = 16
LINE_SPACING = 8
FONT_SIZE = 28
SMALL_FONT_SIZE = 22
BUBBLE_COLOR_AI = (220, 248, 198)
BUBBLE_COLOR_USER = (255, 255, 255)
TEXT_COLOR = (0, 0, 0)
BG_COLOR = (245, 245, 245)

# Auto-detect Chinese fonts (Windows + Linux)
import platform as _platform
FONT_PATHS = []
if _platform.system() == "Windows":
    _wf = os.path.join(os.environ.get("WINDIR", r"C:\Windows"), "Fonts")
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


def load_dialogue(json_path):
    """
    假设json结构：
    {
        "ai_name": "AI角色名",
        "ai_schedule": "AI日程",
        "user_name": "用户角色名",
        "user_schedule": "用户日程",
        "dialogues": [
            {"speaker": "ai", "text": "你好！"},
            {"speaker": "user", "text": "你好，AI。"},
            ...
        ]
    }
    """
    with open(json_path, 'r', encoding='utf-8') as f:
        return json.load(f)


def draw_text(draw, text, font, x, y, max_width):
    # 自动换行绘制文本，返回绘制后的高度
    lines = []
    words = text.split(' ')
    line = ''
    for word in words:
        test_line = line + word + ' '
        bbox = font.getbbox(test_line)
        w = bbox[2] - bbox[0]
        h = bbox[3] - bbox[1]
        if w > max_width and line:
            lines.append(line)
            line = word + ' '
        else:
            line = test_line
    if line:
        lines.append(line)
    height = 0
    for l in lines:
        draw.text((x, y + height), l, fill=TEXT_COLOR, font=font)
        bbox = font.getbbox(l)
        line_height = bbox[3] - bbox[1]
        height += line_height + LINE_SPACING
    return height


# 自动换行函数（提升到顶部，确保可用）
def draw_multiline_text(draw, text, x, y, font, max_width, fill):
    lines = []
    for raw_line in text.split('\n'):
        line = ''
        for char in raw_line:
            test_line = line + char
            bbox = font.getbbox(test_line)
            w = bbox[2] - bbox[0]
            if w > max_width and line:
                lines.append(line)
                line = char
            else:
                line = test_line
        if line:
            lines.append(line)
    for l in lines:
        draw.text((x, y), l, fill=fill, font=font)
        y += font.size + 6
    return y


def create_dialogue_image(data, output_path, character_schedule_text=None, filename=None):
    order = ['morning', 'noon', 'afternoon', 'evening', 'night']
    font = ImageFont.truetype(FONT_PATH, FONT_SIZE)
    small_font = ImageFont.truetype(FONT_PATH, SMALL_FONT_SIZE)
    bold_font = small_font  # Pillow doesn't support bold; use same font
    # Estimate height
    dialogue_count = len(data.get('dialogues', []))
    est_height = 500 + dialogue_count * 150
    img = Image.new('RGB', (IMG_WIDTH, est_height), BG_COLOR)
    draw = ImageDraw.Draw(img)

    # Fill entire background to prevent black areas
    draw.rectangle([0, 0, IMG_WIDTH, est_height], fill=BG_COLOR)

    # Determine language from preset_set: if preset does NOT contain "_CN", use English
    preset_set = data.get('preset_set', '')
    use_english = '_CN' not in preset_set if preset_set else False

    # Use filename from data (passed by Go handler) if available, fallback to arg
    effective_filename = data.get('filename') or filename or ''

    # 1. Determine AI name
    ai_name_map_cn = {'Kurisu': '牧濑红莉栖'}
    ai_name_raw = data.get('ai_name') or data.get('character_name') or ''
    if not ai_name_raw:
        ai_name_raw = 'AI'
    # Only translate to Chinese name if in Chinese mode
    if use_english:
        ai_name = ai_name_raw
    else:
        ai_name = ai_name_map_cn.get(ai_name_raw, ai_name_raw)

    # Build a set of known AI speaker names for dynamic classification
    ai_speaker_names = {'ai', 'AI'}
    ai_speaker_names.add(ai_name_raw)
    ai_speaker_names.add(ai_name)
    for k, v in ai_name_map_cn.items():
        ai_speaker_names.add(k)
        ai_speaker_names.add(v)

    user_speaker_names = {'User', 'user', '用户'}

    # 2. Determine user name
    user_name = data.get('user_name') or ''
    if not user_name:
        for d in data.get('dialogues', []):
            if d.get('speaker', '') in user_speaker_names and d.get('name'):
                user_name = d['name']
                break
    if not user_name:
        profile = data.get('user_profile') or data.get('user_info') or data.get('user')
        if profile and isinstance(profile, dict):
            user_name = profile.get('name') or profile.get('full_name') or profile.get('nickname') or ''
    if not user_name:
        user_name = 'User' if use_english else '用户'
    user_speaker_names.add(user_name)

    # Extract version and day from fields, then filename
    version = data.get('version', None)
    day = data.get('day_simu', None)
    if version is None and effective_filename:
        m_day = re.search(r'Day(\d+)', effective_filename)
        if m_day:
            version = m_day.group(1)
    if day is None and effective_filename:
        m_dup = re.search(r'dup_(\d+)', effective_filename)
        if m_dup:
            day = m_dup.group(1)
    if version is None:
        version = data.get('preset_set', 'v1')
    if day is None:
        for d in data.get('dialogues', []):
            if 'day' in d:
                day = str(d['day'])
                break
        if day is None:
            sched = data.get('ai_schedule')
            if isinstance(sched, dict) and 'day' in sched:
                day = str(sched['day'])

    if use_english:
        top_text = f"Character: {ai_name}    User: {user_name}    Version: {version}    Day: {day}"
    else:
        top_text = f"角色：{ai_name}    用户：{user_name}    版本：{version}    天数：{day}"
    draw.text((MARGIN, MARGIN), top_text, fill=TEXT_COLOR, font=small_font)
    y = MARGIN + 40

    # Avatar image path map: character name -> (folder, image file)
    avatar_paths = {
        'Kurisu': ('kurisu', 'kurisu_image_gen.png'),
        '牧濑红莉栖': ('kurisu', 'kurisu_image_gen.png'),
    }
    # Try to find avatar for the current AI character
    # Check if there's a matching avatar by looking in static/<folder>/ directories
    base_dir = os.path.dirname(__file__)
    if ai_name_raw not in avatar_paths and ai_name not in avatar_paths:
        # Try to auto-discover avatar: look for <name>_image_gen.png or <name>_avatar.png
        for candidate_name in [ai_name_raw, ai_name]:
            if not candidate_name or candidate_name == 'AI':
                continue
            # Try common directory patterns
            for folder_name in [candidate_name.lower().replace(' ', '_'), candidate_name]:
                for img_name in [f'{folder_name}_image_gen.png', f'{folder_name}_avatar.png']:
                    candidate_path = os.path.join(base_dir, 'static', folder_name, img_name)
                    if os.path.exists(candidate_path):
                        avatar_paths[ai_name_raw] = (folder_name, img_name)
                        break
                if ai_name_raw in avatar_paths:
                    break

    # 2. 中间对话区（仿微信气泡）
    bubble_max_width = IMG_WIDTH - 2 * (MARGIN + PADDING + 60)
    avatar_radius = 32
    gap = 18
    for d in data['dialogues']:
        speaker = d.get('speaker', '')
        text = d['text']
        # Dynamic speaker classification using the name sets
        is_ai = speaker in ai_speaker_names or (speaker not in user_speaker_names and speaker != '')
        is_user = speaker in user_speaker_names
        # If speaker is unknown and not user, treat as AI (character)
        if not is_ai and not is_user:
            is_ai = True
        # 调试输出区分角色和AI
        if is_ai:
            print(f"[AI]：{text[:80]}")
        elif is_user:
            print(f"[USER]：{text[:80]}")
        # 计算文本尺寸
        lines = []
        words = list(text)
        line = ''
        for word in words:
            test_line = line + word
            bbox = font.getbbox(test_line)
            w = bbox[2] - bbox[0]
            if w > bubble_max_width and line:
                lines.append(line)
                line = word
            else:
                line = test_line
        if line:
            lines.append(line)
        text_h = 0
        text_w = 0
        for l in lines:
            bbox = font.getbbox(l)
            w = bbox[2] - bbox[0]
            h = bbox[3] - bbox[1]
            text_h += h + LINE_SPACING
            text_w = max(text_w, w)
        bubble_w = min(bubble_max_width, text_w + 2 * BUBBLE_PADDING)
        bubble_h = text_h + 2 * BUBBLE_PADDING
        # 气泡和头像位置
        if is_ai:
            avatar_x = MARGIN + PADDING
            bubble_x = avatar_x + avatar_radius * 2 + 10
            bubble_color = (200, 230, 255)  # AI浅蓝色
        else:
            bubble_x = IMG_WIDTH - MARGIN - PADDING - bubble_w
            avatar_x = IMG_WIDTH - MARGIN - PADDING - avatar_radius * 2 - 10
            bubble_color = (200, 255, 200)  # 用户浅绿色
        # 绘制头像（AI）
        if is_ai:
            avatar_y = y + bubble_h // 2 - avatar_radius
            # Try to load character avatar from avatar_paths
            avatar_loaded = False
            for name_key in [speaker, ai_name_raw, ai_name]:
                if name_key in avatar_paths:
                    folder_name, img_name = avatar_paths[name_key]
                    try:
                        avatar_img = PILImage.open(os.path.join(base_dir, 'static', folder_name, img_name)).convert('RGBA')
                        avatar_img = avatar_img.resize((avatar_radius * 2, avatar_radius * 2))
                        mask = PILImage.new('L', (avatar_radius * 2, avatar_radius * 2), 0)
                        mask_draw = ImageDraw.Draw(mask)
                        mask_draw.ellipse((0, 0, avatar_radius * 2, avatar_radius * 2), fill=255)
                        avatar_img.putalpha(mask)
                        img.paste(avatar_img, (avatar_x, avatar_y), avatar_img)
                        avatar_loaded = True
                        break
                    except Exception as e:
                        print(f"[WARN] {name_key}头像加载失败: {e}, 使用默认色块")
            if not avatar_loaded:
                # Default avatar: blue circle with first character of name
                draw.ellipse([avatar_x, avatar_y, avatar_x + avatar_radius * 2, avatar_y + avatar_radius * 2], fill=(120,180,255))
                label = ai_name[0] if ai_name else "AI"
                draw.text((avatar_x + avatar_radius - 12, avatar_y + avatar_radius - 16), label, fill=(255,255,255), font=small_font)
        # 绘制气泡
        draw.rounded_rectangle([bubble_x, y, bubble_x + bubble_w, y + bubble_h], radius=18, fill=bubble_color)
        # 绘制文本
        text_y = y + BUBBLE_PADDING
        for l in lines:
            draw.text((bubble_x + BUBBLE_PADDING, text_y), l, fill=TEXT_COLOR, font=font)
            bbox = font.getbbox(l)
            line_height = bbox[3] - bbox[1]
            text_y += line_height + LINE_SPACING
        y += bubble_h + gap

    # 3. Schedule section at the bottom
    y += 20
    ai_schedule = data.get('ai_schedule', '')
    user_schedule = data.get('user_schedule', '')
    ai_name_display = ai_name if ai_name else ('AI' if use_english else 'AI')
    user_name_display = user_name if user_name else ('User' if use_english else '用户')

    # Language-aware time labels and prefixes
    if use_english:
        time_map = {'morning': 'Morning', 'noon': 'Noon', 'afternoon': 'Afternoon', 'evening': 'Evening', 'night': 'Night'}
        time_prefix = ''       # no prefix word in English
        time_suffix = ': '     # "Morning: ..."
        schedule_label_ai = f"{ai_name_display} Schedule:"
        schedule_label_user = f"{user_name_display} Schedule:"
    else:
        time_map = {'morning': '早上', 'noon': '中午', 'afternoon': '下午', 'evening': '晚上', 'night': '夜晚'}
        time_prefix = '在'
        time_suffix = '，'
        schedule_label_ai = f"{ai_name_display}日程："
        schedule_label_user = f"{user_name_display}日程："

    # AI schedule
    if ai_schedule:
        ai_texts = []
        if isinstance(ai_schedule, dict):
            for key in order:
                val = str(ai_schedule.get(key, '')).strip()
                if val:
                    ai_texts.append((key, val))
        else:
            if str(ai_schedule).strip():
                ai_texts.append((None, str(ai_schedule).strip()))
        if ai_texts:
            bar_height = 36
            bar_color = (255, 100, 100)
            draw.rectangle([0, y, IMG_WIDTH, y + bar_height], fill=bar_color)
            draw.text((MARGIN, y), schedule_label_ai, fill=(255,255,255), font=bold_font)
            y += bar_height + 4
            for key, val in ai_texts:
                if key:
                    time_desc = time_map.get(key, key)
                    label = f"{time_prefix}{time_desc}{time_suffix}"
                    draw.text((MARGIN, y), label, fill=(120,60,60), font=small_font)
                    y = draw_multiline_text(draw, val, MARGIN + 100, y, small_font, IMG_WIDTH - 2*MARGIN - 100, (80,40,40))
                else:
                    y = draw_multiline_text(draw, val, MARGIN, y, small_font, IMG_WIDTH - 2*MARGIN, (80,40,40))
                y += 10
            y += 10
    # User schedule
    if user_schedule:
        user_texts = []
        user_schedule_dict = None
        if isinstance(user_schedule, dict):
            user_schedule_dict = user_schedule
        elif isinstance(user_schedule, str):
            try:
                user_schedule_dict = _json.loads(user_schedule)
                if not isinstance(user_schedule_dict, dict):
                    user_schedule_dict = None
            except Exception:
                user_schedule_dict = None
        if user_schedule_dict:
            for key in order:
                val = str(user_schedule_dict.get(key, '')).strip()
                if val:
                    user_texts.append((key, val))
        else:
            if str(user_schedule).strip():
                user_texts.append((None, str(user_schedule).strip()))
        if user_texts:
            bar_height = 36
            bar_color = (100, 200, 100)
            draw.rectangle([0, y, IMG_WIDTH, y + bar_height], fill=bar_color)
            draw.text((MARGIN, y), schedule_label_user, fill=(255,255,255), font=bold_font)
            y += bar_height + 4
            for key, val in user_texts:
                if key:
                    time_desc = time_map.get(key, key)
                    label = f"{time_prefix}{time_desc}{time_suffix}"
                    draw.text((MARGIN, y), label, fill=(80,80,80), font=small_font)
                    y = draw_multiline_text(draw, val, MARGIN + 100, y, small_font, IMG_WIDTH - 2*MARGIN - 100, (80,80,80))
                else:
                    y = draw_multiline_text(draw, val, MARGIN, y, small_font, IMG_WIDTH - 2*MARGIN, (80,80,80))
                y += 10
            y += 10

    # 4. Bottom info (repeated for reference)
    bottom_y = y + 30
    if use_english:
        bottom_text = f"Character: {ai_name}    User: {user_name}    Version: {version}    Day: {day}"
    else:
        bottom_text = f"角色：{ai_name}    用户：{user_name}    版本：{version}    天数：{day}"
    draw.text((MARGIN, bottom_y), bottom_text, fill=TEXT_COLOR, font=small_font)
    # 5. Studio logo/tagline
    if use_english:
        logo_text = "2025 Worldline 2 Studio · Building the Future of Emotional AI"
    else:
        logo_text = "2025 世界线二工作室 · 构建情感 AI 的未来"
    logo_y = bottom_y + 40
    draw.text((MARGIN, logo_y), logo_text, fill=(120,120,120), font=small_font)
    
    # 计算最终高度并创建新的精确尺寸图片
    final_height = logo_y + 50
    
    # 创建一个精确尺寸的新图片，确保背景色正确
    final_img = Image.new('RGB', (IMG_WIDTH, final_height), BG_COLOR)
    # 将内容复制到新图片
    final_img.paste(img.crop((0, 0, IMG_WIDTH, min(final_height, est_height))), (0, 0))
    
    final_img.save(output_path, 'PNG', optimize=True)
    print(f"图片已保存到: {output_path}")


def main():
    import argparse
    parser = argparse.ArgumentParser(description='将对话转为长图片')
    parser.add_argument('--input', required=True, help='输入json文件路径')
    parser.add_argument('--output', default='dialogue.png', help='输出图片路径')
    args = parser.parse_args()
    data = load_dialogue(args.input)
    # 传递文件名给create_dialogue_image
    create_dialogue_image(data, args.output, filename=os.path.basename(args.input))


if __name__ == '__main__':
    main()

