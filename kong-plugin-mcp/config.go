package main

type Config struct {
	MCPServerPath string `json:"mcp_server_path"`
	ToolName      string `json:"tool_name"`
	Timeout       int    `json:"timeout"`
}
