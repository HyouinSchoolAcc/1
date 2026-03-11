#!/bin/bash
# Quick test script for the new novel style renderer

echo "🎬 Testing Video Novel Style Renderer"
echo "======================================"
echo ""

cd video_renderer

# Test 1: Novel style with default settings (3 messages per side)
echo "Test 1: Novel style with default settings (3 messages per side)"
python generator.py example_novel.json output_novel_default.mp4 --style novel
echo "✓ Created: output_novel_default.mp4"
echo ""

# Test 2: Novel style with only 1 message per side (TikTok style)
echo "Test 2: Novel style with 1 message per side (TikTok/Reels style)"
python generator.py example_novel.json output_novel_minimal.mp4 --style novel --max-messages 1
echo "✓ Created: output_novel_minimal.mp4"
echo ""

# Test 3: Novel style with 5 messages per side (more context)
echo "Test 3: Novel style with 5 messages per side (more context)"
python generator.py example_novel.json output_novel_full.mp4 --style novel --max-messages 5
echo "✓ Created: output_novel_full.mp4"
echo ""

# Test 4: Traditional chat style for comparison
echo "Test 4: Traditional chat style (for comparison)"
python generator.py example_novel.json output_chat_traditional.mp4 --style chat
echo "✓ Created: output_chat_traditional.mp4"
echo ""

echo "======================================"
echo "✅ All tests complete!"
echo ""
echo "Generated videos:"
echo "  - output_novel_default.mp4 (novel style, 3 messages)"
echo "  - output_novel_minimal.mp4 (novel style, 1 message - TikTok)"
echo "  - output_novel_full.mp4 (novel style, 5 messages)"
echo "  - output_chat_traditional.mp4 (traditional chat style)"
echo ""
echo "Compare the different styles to see which works best for your use case!"
