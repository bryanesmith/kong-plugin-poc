#!/bin/bash

echo "Testing Kong MCP Streamable HTTP Proxy"
echo "========================================"
echo ""

echo "Test 1: Initialize MCP session"
echo "-------------------------------"
curl -X POST http://localhost:8000/mcp/wordle \
  -H "Content-Type: application/json" \
  -H "Accept: application/json, text/event-stream" \
  -H "MCP-Protocol-Version: 2024-11-05" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "initialize",
    "params": {
      "protocolVersion": "2024-11-05",
      "capabilities": {},
      "clientInfo": {
        "name": "test-client",
        "version": "1.0.0"
      }
    }
  }'
echo ""
echo ""

echo "Test 2: List available tools"
echo "-----------------------------"
curl -X POST http://localhost:8000/mcp/wordle \
  -H "Content-Type: application/json" \
  -H "Accept: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 2,
    "method": "tools/list",
    "params": {}
  }'
echo ""
echo ""

echo "Test 3: Call get_wordle_suggestions tool"
echo "-----------------------------------------"
curl -X POST http://localhost:8000/mcp/wordle \
  -H "Content-Type: application/json" \
  -H "Accept: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 3,
    "method": "tools/call",
    "params": {
      "name": "get_wordle_suggestions",
      "arguments": {
        "guesses": ["slate", "crane"]
      }
    }
  }'
echo ""
echo ""

echo "Test 4: Open SSE stream for server messages (GET)"
echo "--------------------------------------------------"
echo "This will keep the connection open. Press Ctrl+C to stop."
curl -N -X GET http://localhost:8000/mcp/wordle \
  -H "Accept: text/event-stream"
