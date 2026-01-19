FROM kong:latest

USER root

# Copy MCP server binary
COPY ./mcp_server/mcp_server /usr/local/bin/mcp_server
RUN chmod +x /usr/local/bin/mcp_server

# Copy Go plugin binary (embedded server)
COPY ./kong-plugin-mcp/mcp-bridge /usr/local/bin/mcp-bridge
RUN chmod +x /usr/local/bin/mcp-bridge

# Enable Go plugin support (embedded server style)
ENV KONG_PLUGINS=bundled,mcp-bridge
ENV KONG_PLUGINSERVER_NAMES=mcp-bridge
ENV KONG_PLUGINSERVER_MCP_BRIDGE_START_CMD=/usr/local/bin/mcp-bridge
ENV KONG_PLUGINSERVER_MCP_BRIDGE_QUERY_CMD="/usr/local/bin/mcp-bridge -dump"

USER kong

EXPOSE 8000 8443 8001 8444
