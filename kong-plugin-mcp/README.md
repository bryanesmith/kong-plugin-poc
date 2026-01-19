# Kong MCP Bridge Plugin

A Kong Gateway plugin written in Go that bridges HTTP API requests to MCP (Model Context Protocol) servers via stdio.

## Overview

This plugin intercepts HTTP requests, spawns an MCP server process, communicates via JSON-RPC over stdin/stdout, and returns the MCP response as an HTTP response.

## Architecture

```
HTTP Request → Kong → Go Plugin → MCP Server (stdio)
                                       ↓
HTTP Response ← Kong ← Go Plugin ← JSON-RPC Response
```

## Configuration

The plugin accepts the following configuration parameters:

- `mcp_server_path` (string, required): Absolute path to the MCP server executable
- `tool_name` (string, required): Name of the MCP tool to call
- `timeout` (int, default: 5000): Timeout in milliseconds for MCP operations

## Example Configuration

```yaml
plugins:
  - name: mcp-bridge
    config:
      mcp_server_path: /usr/local/bin/mcp_server
      tool_name: get_wordle_suggestions
      timeout: 5000
```

## Building

```bash
make build
```

This creates `mcp-bridge.so` which can be loaded by Kong.

## Files

- `main.go`: Kong plugin handler and entry point
- `config.go`: Configuration schema
- `mcp_client.go`: MCP server communication logic
- `Makefile`: Build automation

## How It Works

1. Plugin intercepts HTTP request in the `Access` phase
2. Reads and parses JSON request body
3. Spawns MCP server process using `exec.Cmd`
4. Sends MCP initialization handshake:
   - `initialize` request
   - `notifications/initialized` notification
5. Sends `tools/call` request with tool name and arguments
6. Reads JSON-RPC response from MCP server stdout
7. Translates MCP response to HTTP response
8. Closes MCP server process

## Error Handling

The plugin handles various error scenarios:

- Invalid request body (400)
- Invalid JSON (400)
- Failed to start MCP server (500)
- Failed to initialize MCP (500)
- MCP tool execution errors (500)
- Request timeout (504)

All errors are returned as JSON with an `error` field.
