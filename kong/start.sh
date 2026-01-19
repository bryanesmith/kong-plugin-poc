#!/bin/bash

# Start MCP HTTP proxy in the background
echo "Starting MCP HTTP Proxy..."
/usr/local/bin/mcp-http-proxy &
PROXY_PID=$!

# Give the proxy a moment to start
sleep 2

# Start Kong
echo "Starting Kong..."
exec /docker-entrypoint.sh kong docker-start
