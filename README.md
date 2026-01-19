# kong-plugin-poc

Kong Gateway plugin that provides **Streamable HTTP (SSE)** transport for MCP servers, enabling MCP clients like Windsurf to connect via HTTP.

## Architecture

This project demonstrates a complete integration between Kong Gateway, an MCP HTTP Proxy, and an MCP server:

- **Kong Gateway**: Routes HTTP requests and provides API gateway features (auth, rate limiting, etc.)
- **Kong MCP Bridge Plugin** (Go): Simple HTTP proxy that forwards requests to MCP HTTP Proxy
- **MCP HTTP Proxy** (Go + MCP SDK): Implements Streamable HTTP transport (SSE), handles MCP protocol
- **MCP Server** (Go): Provides Wordle suggestions via the `get_wordle_suggestions` tool

## Request Flow

```
┌─────────────────┐
│ MCP Client      │ (Windsurf, curl, etc.)
│ (Streamable     │
│  HTTP/SSE)      │
└────────┬────────┘
         │ HTTP POST/GET (port 8000)
         ↓
┌─────────────────┐
│ Kong Gateway    │ (API Gateway - routing, auth, rate limiting)
└────────┬────────┘
         │ HTTP forward
         ↓
┌─────────────────┐
│ Kong MCP Bridge │ (Simple HTTP proxy plugin)
│ Plugin (Go)     │
└────────┬────────┘
         │ HTTP (port 9000, internal)
         ↓
┌─────────────────┐
│ MCP HTTP Proxy  │ (MCP Go SDK - handles SSE, sessions, protocol)
│ (Go + SDK)      │
└────────┬────────┘
         │ stdio JSON-RPC
         ↓
┌─────────────────┐
│ MCP Server      │ (Wordle tool implementation)
│ (Go)            │
└─────────────────┘
```

**Key Features:**
- ✅ **Streamable HTTP (SSE)** transport for MCP protocol
- ✅ **Compatible with Windsurf** and other MCP clients
- ✅ **Spec-compliant** using official MCP Go SDK
- ✅ **Simple architecture** - Kong plugin is just HTTP forwarding (~100 lines)
- ✅ **Scalable** - Easy to add auth, rate limiting, multiple MCP servers

## Prerequisites

- Maker sure Docker is installed and running
- Install [decK](https://developer.konghq.com/deck/) CLI version 1.43 or later

## Get Started

### Quick Start

```bash
# Build and start Kong with the MCP plugin
make run-kong

# Test the integration
make test-kong
```

### Step-by-Step

**1. Test the MCP server standalone:**
```bash
make test-mcp
```

**2. Build and run Kong:**
```bash
make run-kong
```

**3. Verify Kong is running:**
```bash
curl http://localhost:8001/status
```

**4. Test the MCP integration:**
```bash
# Initialize MCP session
curl -X POST http://localhost:8000/mcp/wordle \
  -H "Content-Type: application/json" \
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

# Call the Wordle tool
curl -X POST http://localhost:8000/mcp/wordle \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 2,
    "method": "tools/call",
    "params": {
      "name": "get_wordle_suggestions",
      "arguments": {"guesses": ["slate", "crane"]}
    }
  }'
```

Expected response (SSE format):
```
event: message
id: <session_id>_0
data: {"jsonrpc":"2.0","id":2,"result":{"content":[{"type":"text","text":"{\"suggestions\":[\"apple\"]}"}]}}
```

**5. Run test script:**
```bash
./test_streamable_http.sh
```

## Project Structure

```
.
├── kong-plugin-mcp/          # Kong plugin (HTTP proxy)
│   ├── main.go               # Plugin entry point
│   ├── config.go             # Configuration schema
│   ├── http_proxy.go         # HTTP forwarding handler
│   └── Makefile              # Plugin build
├── mcp-http-proxy/           # MCP HTTP Proxy service
│   ├── main.go               # Proxy using MCP Go SDK
│   └── go.mod                # Go module with SDK dependency
├── mcp_server/               # MCP server
│   ├── main.go               # Server entry point
│   ├── wordlemcpserver/      # Wordle tool implementation
│   └── test_request.sh       # Standalone test
├── kong/
│   ├── Dockerfile            # Kong image with all components
│   ├── kong.yml              # Kong declarative config
│   └── start.sh              # Startup script (runs proxy + Kong)
├── test_streamable_http.sh   # Streamable HTTP tests
├── WINDSURF_SETUP.md         # Windsurf integration guide
├── IMPLEMENTATION_SUMMARY.md # Technical details
└── Makefile                  # Main build automation
```

## Reference

### MCP Server

```bash
make build-mcp  # Build MCP server binary
make test-mcp   # Build and run standalone test
make clean-mcp  # Remove built files
```

### Kong Plugin

```bash
make build-plugin  # Build Go plugin (.so file)
make clean-plugin  # Remove built files
```

### Kong Gateway

```bash
make build-kong  # Build Docker image with plugin and MCP server
make run-kong    # Build and run Kong in background
make logs-kong   # View Kong logs
make test-kong   # Run integration tests
make stop-kong   # Stop the container
make clean-kong  # Remove Docker image
```

### All

```bash
make all    # Build everything
make clean  # Clean everything
```

## Windsurf Integration

This Kong MCP proxy can be used with Windsurf and other MCP clients that support Streamable HTTP transport.

**Add to `~/.codeium/windsurf/mcp_config.json`:**

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

After restarting Windsurf, you'll have access to the `get_wordle_suggestions` tool.

See **[WINDSURF_SETUP.md](WINDSURF_SETUP.md)** for detailed setup instructions.

## More Info

- [Get started with Kong Gateway](https://developer.konghq.com/gateway/get-started/)
- [MCP Specification](https://spec.modelcontextprotocol.io/)
- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk)
- [WINDSURF_SETUP.md](WINDSURF_SETUP.md) - Windsurf integration guide
- [IMPLEMENTATION_SUMMARY.md](IMPLEMENTATION_SUMMARY.md) - Technical details