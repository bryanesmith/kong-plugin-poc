# kong-plugin-poc
POC of a Kong plugin that routes HTTP API calls to an MCP server.

## Prerequisites

- Maker sure Docker is installed and running
- Install [decK](https://developer.konghq.com/deck/) CLI version 1.43 or later

## Get Started

```bash
make run

# validate admin API is running
deck gateway ping 

# validate the test endpoint
export KONNECT_PROXY_URL=http://localhost:8000 
curl "$KONNECT_PROXY_URL/mcp/anything" \
     --no-progress-meter --fail-with-body
```

## Commands

```bash
make all # Build
make run # Build and run 
make logs # View Kong logs
make stop # Stop the container
make clean # Clean up
```

## More Info

- [Get started with Kong Gateway](https://developer.konghq.com/gateway/get-started/)