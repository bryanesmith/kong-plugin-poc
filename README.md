# kong-plugin-poc
POC of a Kong plugin that routes HTTP API calls to an MCP server.

## Prerequisites

- Maker sure Docker is installed and running
- Install [decK](https://developer.konghq.com/deck/) CLI version 1.43 or later

## Get Started

### MCP Server

```bash
make run-mcp    # note: this is a foreground process
```

### Kong

```bash
make run-kong

# validate admin API is running
deck gateway ping 

# validate the test endpoint
export KONNECT_PROXY_URL=http://localhost:8000 
curl "$KONNECT_PROXY_URL/mcp/anything" \
     --no-progress-meter --fail-with-body
```

## Reference

### MCP Server 

```bash
make build-mcp # build 
make run-mcp # build and run 
make stop-mcp # stop
make clean-mcp # remove built files
```

### Kong

```bash
make build-kong # build 
make run-kong # build and run 
make logs-kong # view Kong logs
make stop-kong # stop the container
make clean-kong # remove built files
```

## More Info

- [Get started with Kong Gateway](https://developer.konghq.com/gateway/get-started/)