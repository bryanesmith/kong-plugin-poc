package main

type Config struct {
	ProxyURL string `json:"proxy_url"` // URL of MCP HTTP proxy (default: http://localhost:9000)
	Timeout  int    `json:"timeout"`   // Request timeout in milliseconds (default: 30000)
}
