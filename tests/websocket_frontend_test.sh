#!/bin/bash

# WebSocket Frontend UX Test Script
# This script tests the WebSocket functionality from a user perspective

echo "==================================="
echo "WebSocket Frontend UX Test"
echo "==================================="
echo ""

# Check if server is running
echo "1. Checking if server is running on port 8080..."
if curl -s http://localhost:8080/api/health > /dev/null 2>&1; then
    echo "   ✅ Server is running"
else
    echo "   ❌ Server is not running on port 8080"
    echo "   Please start the server with: ./forge-orchestrator"
    exit 1
fi

echo ""
echo "2. Building frontend application..."
cd frontend
if npm run build > /dev/null 2>&1; then
    echo "   ✅ Frontend build successful"
else
    echo "   ❌ Frontend build failed"
    exit 1
fi

echo ""
echo "3. Testing WebSocket connection from browser simulation..."
cd ..

# Create a simple Node.js WebSocket test
cat > /tmp/ws_frontend_test.js << 'EOF'
const WebSocket = require('ws');

const ws = new WebSocket('ws://localhost:8080/ws');

ws.on('open', function open() {
  console.log('   ✅ WebSocket connected from frontend client');
  
  // Send a test message
  const message = JSON.stringify({
    type: 'FRONTEND_TEST',
    payload: {
      message: 'Test from simulated frontend',
      timestamp: new Date().toISOString()
    }
  });
  
  ws.send(message);
  console.log('   ✅ Test message sent');
});

ws.on('message', function message(data) {
  console.log('   ✅ Received message:', data.toString());
  
  // Parse and validate
  try {
    const parsed = JSON.parse(data);
    if (parsed.type) {
      console.log('   ✅ Message type:', parsed.type);
    }
  } catch (e) {
    console.log('   ⚠️  Message is not JSON');
  }
  
  // Close after receiving response
  setTimeout(() => {
    ws.close();
    console.log('\n   ✅ Connection closed gracefully');
    console.log('\n=================================');
    console.log('🎉 All frontend UX tests PASSED!');
    console.log('=================================');
    process.exit(0);
  }, 1000);
});

ws.on('error', function error(err) {
  console.error('   ❌ WebSocket error:', err.message);
  process.exit(1);
});

// Timeout after 10 seconds
setTimeout(() => {
  console.error('   ❌ Test timeout - no response received');
  process.exit(1);
}, 10000);
EOF

# Run the test
node /tmp/ws_frontend_test.js

# Clean up
rm /tmp/ws_frontend_test.js
