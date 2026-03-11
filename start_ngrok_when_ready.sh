#!/bin/bash

# Script to monitor and start ngrok tunnel for wl2.studio when available

DOMAIN="wl2.studio"
PORT=5002
LOG_FILE="/home/exx/Desktop/fine-tune/data_labler_UI_production/ngrok.log"

echo "=============================================="
echo "Ngrok Tunnel Starter for $DOMAIN"
echo "=============================================="
echo ""
echo "📋 Instructions to stop the old tunnel:"
echo "   1. Open: https://dashboard.ngrok.com/cloud-edge/endpoints"
echo "   2. Find: $DOMAIN"
echo "   3. Click: Stop or Delete"
echo ""
echo "⏳ Waiting for endpoint to become available..."
echo ""

# Keep trying to start ngrok
attempt=1
while true; do
    echo "Attempt $attempt: Trying to start ngrok tunnel..."
    
    # Try to start ngrok
    ngrok http --domain=$DOMAIN $PORT > "$LOG_FILE" 2>&1 &
    NGROK_PID=$!
    
    # Wait a moment for ngrok to start
    sleep 3
    
    # Check if it succeeded
    if grep -q "started tunnel" "$LOG_FILE" 2>/dev/null || ! grep -q "ERR_NGROK_334" "$LOG_FILE" 2>/dev/null; then
        echo ""
        echo "✅ SUCCESS! Ngrok tunnel started!"
        echo ""
        echo "Tunnel Details:"
        echo "==============="
        echo "Domain:  https://$DOMAIN"
        echo "Port:    $PORT"
        echo "PID:     $NGROK_PID"
        echo "Log:     $LOG_FILE"
        echo ""
        echo "🌐 Your site is now live at: https://$DOMAIN/"
        echo ""
        break
    else
        # Kill the failed attempt
        kill $NGROK_PID 2>/dev/null
        
        # Show error
        echo "❌ Endpoint still in use. Waiting 5 seconds before retry..."
        echo "   (Stop the old endpoint at: https://dashboard.ngrok.com/cloud-edge/endpoints)"
        echo ""
        sleep 5
    fi
    
    attempt=$((attempt + 1))
done
