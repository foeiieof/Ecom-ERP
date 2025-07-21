#!/bin/bash

# Load environment variables from .env file
export $(grep -v '^#' .env | xargs)

# Start the server in background
go run ./cmd/server/main.go &
SERVER_PID=$!

# Wait a bit for server to start
sleep 3

# Test health endpoint
echo "Testing health endpoint..."
curl -s http://localhost:8080/health || echo "Health check failed"

# Test root endpoint
echo -e "\nTesting root endpoint..."
curl -s http://localhost:8080/ || echo "Root endpoint failed"

# Kill the server
kill $SERVER_PID 2>/dev/null || true
wait $SERVER_PID 2>/dev/null || true

echo -e "\nServer test completed."
