#!/bin/bash

echo "Testing Kong MCP Bridge Plugin"
echo "================================"
echo ""

echo "Test 1: Valid Wordle request"
echo "----------------------------"
curl -X POST http://localhost:8000/mcp/wordle \
  -H "Content-Type: application/json" \
  -d '{"guesses": ["slate", "crane"]}'
echo ""
echo ""

echo "Test 2: Empty guesses array"
echo "---------------------------"
curl -X POST http://localhost:8000/mcp/wordle \
  -H "Content-Type: application/json" \
  -d '{"guesses": []}'
echo ""
echo ""

echo "Test 3: Invalid JSON"
echo "--------------------"
curl -X POST http://localhost:8000/mcp/wordle \
  -H "Content-Type: application/json" \
  -d 'invalid json'
echo ""
echo ""

echo "Test 4: Empty request body"
echo "--------------------------"
curl -X POST http://localhost:8000/mcp/wordle \
  -H "Content-Type: application/json" \
  -d '{}'
echo ""
echo ""

echo "Test 5: Check Kong health"
echo "-------------------------"
curl -s http://localhost:8001/status
echo ""
echo ""
