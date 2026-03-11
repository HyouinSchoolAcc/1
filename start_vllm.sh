#!/bin/bash

# vLLM Server Startup Script for 林路 (Lin Lu) Model
# This script starts the vLLM server with the fine-tuned Qwen3-14B model

set -e

# Configuration
MODEL_ID="Qwen/Qwen3-14B"
LORA_PATH="/home/exx/Desktop/fine-tune/12_29_train/outputs/qwen3_masked_full_20251229_014726"
PORT=8000
HOST="0.0.0.0"

# GPU Configuration
GPU_MEMORY_UTILIZATION=0.8

echo "=========================================="
echo "vLLM Server - 林路 (Lin Lu)"
echo "=========================================="
echo "Model: $MODEL_ID"
echo "LoRA: $LORA_PATH"
echo "Port: $PORT"
echo "=========================================="

# Check if LoRA path exists
if [ ! -d "$LORA_PATH" ]; then
    echo "ERROR: LoRA path not found: $LORA_PATH"
    echo "Please ensure the fine-tuned model exists."
    exit 1
fi

# Check if vLLM is installed
if ! python3 -c "import vllm" 2>/dev/null; then
    echo "ERROR: vLLM is not installed."
    echo "Install it with: pip install vllm"
    exit 1
fi

# Start vLLM server with LoRA support
echo "Starting vLLM server..."
python3 -m vllm.entrypoints.openai.api_server \
    --model "$MODEL_ID" \
    --host "$HOST" \
    --port "$PORT" \
    --enable-lora \
    --lora-modules "linlu=$LORA_PATH" \
    --max-lora-rank 64 \
    --trust-remote-code \
    --dtype bfloat16 \
    --gpu-memory-utilization "$GPU_MEMORY_UTILIZATION" \
    --max-model-len 4096

echo "=========================================="
echo "vLLM server stopped."
echo "=========================================="

