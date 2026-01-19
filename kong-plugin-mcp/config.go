package main

type Config struct {
	MCPServerPath string `json:"mcp_server_path"`
	ToolName      string `json:"tool_name"`
	Timeout       int    `json:"timeout"`
	Mode          string `json:"mode"`      // "tool" or "http-proxy"
	ProxyURL      string `json:"proxy_url"` // URL of MCP HTTP proxy
}
