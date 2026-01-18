package main

import (
	"context"
	"log"

	"github.com/bryanesmith/kong-plugin-poc/mcp_server/wordlemcpserver"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "wordle-mcp-server",
			Version: "1.0.0",
		},
		nil,
	)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "get_wordle_suggestions",
			Description: "Get word suggestions for Wordle based on previous guesses",
		},
		wordlemcpserver.GetWordleSuggestions,
	)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
