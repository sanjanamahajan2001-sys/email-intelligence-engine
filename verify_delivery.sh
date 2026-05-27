#!/bin/bash

echo "🔍 Starting Internal Verification Test..."

# 1. Build the latest binary
bash build_dist.sh

# 2. Start the API in the background on a custom port
cd api-dist
API_PORT=8081 ./email-api &
API_PID=$!

# Wait for API to start
sleep 3

echo "📡 Calling API for contact@stripe.com on port 8081..."
RESPONSE=$(curl -s -X POST "http://localhost:8080/v1/web-validate" \
     -H 'Content-Type: application/json' \
     -d '{"email": "contact@stripe.com"}')

echo "--------------------------------------------------------"
echo "RAW RESPONSE:"
echo $RESPONSE | jq . 2>/dev/null || echo $RESPONSE
echo "--------------------------------------------------------"

# Check for key rich fields
if [[ $RESPONSE == *"engagement"* ]] && [[ $RESPONSE == *"authenticity_status"* ]]; then
    echo "✅ SUCCESS: API returned rich intelligence data."
else
    echo "❌ FAILURE: API response is missing rich fields."
    kill $API_PID
    exit 1
fi

# Cleanup
echo "🧹 Cleaning up test process..."
kill $API_PID
echo "✅ Verification Complete!"
