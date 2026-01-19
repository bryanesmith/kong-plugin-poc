package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	mcpServerPath := os.Getenv("MCP_SERVER_PATH")
	if mcpServerPath == "" {
		mcpServerPath = "/usr/local/bin/mcp_server"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	log.Printf("Starting MCP HTTP Proxy on port %s", port)
	log.Printf("MCP Server: %s", mcpServerPath)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create MCP client to connect to the stdio server
	client := mcp.NewClient(&mcp.Implementation{
		Name:    "kong-mcp-proxy",
		Version: "1.0.0",
	}, nil)

	// Create transport to stdio MCP server
	cmdTransport := &mcp.CommandTransport{
		Command: exec.Command(mcpServerPath),
	}

	// Connect to MCP server
	log.Println("Connecting to MCP server via stdio...")
	clientSession, err := client.Connect(ctx, cmdTransport, nil)
	if err != nil {
		log.Fatalf("Failed to connect to MCP server: %v", err)
	}
	defer clientSession.Close()

	log.Println("Successfully connected to MCP server")

	// Create a proxy server that forwards requests to the stdio MCP server
	proxyServer := mcp.NewServer(&mcp.Implementation{
		Name:    "kong-mcp-http-proxy",
		Version: "1.0.0",
	}, nil)

	// Add a dynamic tool handler that proxies to the stdio server
	// First, get the list of tools from the stdio server
	toolsList, err := clientSession.ListTools(ctx, &mcp.ListToolsParams{})
	if err != nil {
		log.Fatalf("Failed to list tools: %v", err)
	}

	// Register each tool with a proxy handler
	for _, tool := range toolsList.Tools {
		toolName := tool.Name
		log.Printf("Registering tool: %s", toolName)

		// Use the generic AddTool with map types for flexibility
		mcp.AddTool(proxyServer, &mcp.Tool{
			Name:        tool.Name,
			Description: tool.Description,
			InputSchema: tool.InputSchema,
		}, func(ctx context.Context, req *mcp.CallToolRequest, args map[string]any) (*mcp.CallToolResult, map[string]any, error) {
			// Forward the call to the stdio MCP server
			result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
				Name:      toolName,
				Arguments: args,
			})
			if err != nil {
				return nil, nil, err
			}
			return result, nil, nil
		})
	}

	// Create Streamable HTTP transport
	httpTransport := &mcp.StreamableServerTransport{
		Stateless: false,
	}

	// Create HTTP handler with CORS
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, MCP-Session-Id, MCP-Protocol-Version, Last-Event-ID")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Use SDK's transport to handle the request
		httpTransport.ServeHTTP(w, r)
	})

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down...")
		cancel()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
	}()

	// Connect the proxy server to the HTTP transport in a goroutine
	go func() {
		log.Println("Starting MCP server with HTTP transport...")
		_, err := proxyServer.Connect(ctx, httpTransport, nil)
		if err != nil {
			log.Printf("Server connection error: %v", err)
		}
	}()

	// Give the server a moment to connect
	time.Sleep(100 * time.Millisecond)

	// Start HTTP server
	log.Printf("MCP HTTP Proxy listening on http://localhost:%s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}

	log.Println("Server stopped")
}
