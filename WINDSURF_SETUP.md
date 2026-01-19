# Using Kong MCP Proxy with Windsurf

This guide explains how to configure Windsurf to use the Kong MCP HTTP proxy.

## Architecture

```
Windsurf (MCP Client)
    ↓ Streamable HTTP (SSE)
Kong Gateway (port 8000)
    ↓ HTTP Proxy Plugin
MCP HTTP Proxy (port 9000)
    ↓ MCP Go SDK
    ↓ stdio
MCP Server (Wordle)
```

## Setup

### 1. Start Kong with MCP HTTP Proxy

```bash
make run-kong
```

This will:
- Build the MCP server, MCP HTTP proxy, and Kong plugin
- Start Kong Gateway on port 8000
- Start MCP HTTP proxy on port 9000 (internal)
- Connect the proxy to the stdio MCP server

### 2. Configure Windsurf

Add the following to your Windsurf MCP configuration file (`~/.codeium/windsurf/mcp_config.json`):

```json
{
  "mcpServers": {
    "wordle-kong-proxy": {
      "disabled": false,
      "headers": {},
      "serverUrl": "http://localhost:8000/mcp/wordle"
    }
  }
}
```

### 3. Restart Windsurf

Restart Windsurf to load the new MCP server configuration.

### 4. Verify Connection

In Windsurf, you should now see the "wordle-kong-proxy" MCP server available. You can use it to:
- List available tools
- Call the `get_wordle_suggestions` tool with guesses

## Testing Without Windsurf

You can test the MCP HTTP proxy directly using curl:

```bash
# Run the test script
./test_streamable_http.sh
```

Or manually:

```bash
# Initialize session
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
      "clientInfo": {"name": "test", "version": "1.0.0"}
    }
  }'

# List tools
curl -X POST http://localhost:8000/mcp/wordle \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 2,
    "method": "tools/list"
  }'

# Call tool
curl -X POST http://localhost:8000/mcp/wordle \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 3,
    "method": "tools/call",
    "params": {
      "name": "get_wordle_suggestions",
      "arguments": {"guesses": ["slate", "crane"]}
    }
  }'
```

## How It Works

1. **Kong Gateway** receives HTTP requests on port 8000 at `/mcp/wordle`
2. **Kong Plugin** (in http-proxy mode) forwards requests to the MCP HTTP proxy on port 9000
3. **MCP HTTP Proxy** (using MCP Go SDK):
   - Implements Streamable HTTP transport (SSE)
   - Connects to the stdio MCP server
   - Proxies tool calls and responses
   - Handles session management
4. **MCP Server** processes tool calls and returns results

## Benefits of This Architecture

- ✅ **Spec-compliant**: Uses official MCP Go SDK for Streamable HTTP transport
- ✅ **Simple**: Kong plugin is just 70 lines of HTTP forwarding code
- ✅ **Maintainable**: SDK handles all SSE, session, and protocol complexity
- ✅ **Scalable**: Can add authentication, rate limiting, etc. at Kong layer
- ✅ **Flexible**: Easy to add more MCP servers behind different routes

## Troubleshooting

### Check if services are running

```bash
# Check Kong logs
docker logs kong-dev

# Look for these messages:
# - "Starting MCP HTTP Proxy..."
# - "Successfully connected to MCP server"
# - "MCP HTTP Proxy listening on http://localhost:9000"
```

### Test the MCP HTTP proxy directly

```bash
# From inside the container
docker exec -it kong-dev curl http://localhost:9000/mcp/wordle \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/list"}'
```

### Verify Kong routing

```bash
# Check Kong health
curl http://localhost:8001/status

# Check Kong routes
curl http://localhost:8001/routes
```

## Next Steps

- Add authentication to Kong routes
- Add rate limiting
- Add more MCP servers behind different routes
- Deploy to production with proper TLS/SSL
