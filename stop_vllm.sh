#!/bin/bash

# Stop vLLM Server Script

echo "Stopping vLLM server..."

# Find and kill vLLM processes
pkill -f "vllm.entrypoints.openai.api_server" 2>/dev/null && echo "vLLM server stopped." || echo "No vLLM server running."

