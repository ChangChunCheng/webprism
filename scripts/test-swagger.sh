#!/bin/bash

# Test Swagger UI integration
# This script starts the server briefly to test Swagger endpoints

set -e

echo "🧪 Testing Swagger UI Integration"
echo "=================================="
echo ""

# Export required environment variables
export WEBPRISM_DATABASE_PASSWORD=webprism
export WEBPRISM_SECURITY_ENCRYPTION_KEY=dev-key-12345678901234567890abcd

# Start server in background
echo "📌 Starting WEBPRISM server in background..."
./bin/webprism-server &
SERVER_PID=$!

# Wait for server to start
echo "⏳ Waiting for server to start..."
sleep 3

# Test endpoints
echo ""
echo "🔍 Testing Swagger endpoints..."
echo ""

# Test swagger.json
echo "1. Testing /swagger.json"
if curl -s -f http://localhost:8080/swagger.json > /dev/null; then
    echo "   ✅ /swagger.json is accessible"
    echo "   📊 Swagger spec preview:"
    curl -s http://localhost:8080/swagger.json | head -20
else
    echo "   ❌ /swagger.json is NOT accessible"
fi

echo ""
echo "2. Testing /swagger-ui/"
if curl -s -f http://localhost:8080/swagger-ui/ | grep -q "WEBPRISM"; then
    echo "   ✅ /swagger-ui/ is accessible"
else
    echo "   ❌ /swagger-ui/ is NOT accessible"
fi

echo ""
echo "3. Testing /swagger-ui (without trailing slash)"
if curl -s -f http://localhost:8080/swagger-ui | grep -q "WEBPRISM"; then
    echo "   ✅ /swagger-ui is accessible"
else
    echo "   ❌ /swagger-ui is NOT accessible"
fi

# Cleanup
echo ""
echo "🧹 Cleaning up..."
kill $SERVER_PID 2>/dev/null || true
sleep 1

echo ""
echo "✨ Test completed!"
echo ""
echo "📚 Access Swagger UI at: http://localhost:8080/swagger-ui/"
echo "📄 Access OpenAPI spec at: http://localhost:8080/swagger.json"
