package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
)

type MCPClient struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	reader *bufio.Reader
}

type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id,omitempty"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

type JSONRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *JSONRPCError   `json:"error,omitempty"`
}

type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewMCPClient(serverPath string) (*MCPClient, error) {
	cmd := exec.Command(serverPath)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	// Capture stderr for debugging
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start MCP server: %w", err)
	}

	// Log stderr in background
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Printf("MCP Server stderr: %s\n", scanner.Text())
		}
	}()

	return &MCPClient{
		cmd:    cmd,
		stdin:  stdin,
		stdout: stdout,
		reader: bufio.NewReader(stdout),
	}, nil
}

func (c *MCPClient) SendRequest(req JSONRPCRequest) (*JSONRPCResponse, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Log the request being sent
	fmt.Printf("MCP Request: %s\n", string(data))

	// Write request with newline
	if _, err := c.stdin.Write(append(data, '\n')); err != nil {
		return nil, fmt.Errorf("failed to write request: %w", err)
	}

	// Read responses until we get one matching our request ID
	// (skip notifications and other responses)
	expectedID := req.ID
	for {
		line, err := c.reader.ReadBytes('\n')
		if err != nil {
			fmt.Printf("Error reading response: %v\n", err)
			return nil, fmt.Errorf("failed to read response: %w", err)
		}

		// Log the response received
		fmt.Printf("MCP Raw Response: %s\n", string(line))

		// Check if response is empty
		if len(line) == 0 || (len(line) == 1 && line[0] == '\n') {
			continue // Skip empty lines
		}

		var resp JSONRPCResponse
		if err := json.Unmarshal(line, &resp); err != nil {
			fmt.Printf("Failed to unmarshal: %v, data: %s\n", err, string(line))
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}

		// If this is a notification (no ID), skip it
		if resp.ID == 0 {
			fmt.Printf("Skipping notification\n")
			continue
		}

		// If this response matches our request, return it
		if resp.ID == expectedID {
			return &resp, nil
		}

		// Otherwise, skip this response (might be from a different request)
		fmt.Printf("Skipping response with ID %d (expected %d)\n", resp.ID, expectedID)
	}
}

func (c *MCPClient) Initialize() error {
	initReq := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
		Params: map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]interface{}{},
			"clientInfo": map[string]string{
				"name":    "kong-mcp-plugin",
				"version": "1.0.0",
			},
		},
	}

	if _, err := c.SendRequest(initReq); err != nil {
		return err
	}

	// Send initialized notification (no response expected)
	initNotif := JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "notifications/initialized",
	}

	data, err := json.Marshal(initNotif)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	if _, err := c.stdin.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("failed to write notification: %w", err)
	}

	return nil
}

func (c *MCPClient) CallTool(toolName string, arguments map[string]interface{}) (*JSONRPCResponse, error) {
	req := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      2,
		Method:  "tools/call",
		Params: map[string]interface{}{
			"name":      toolName,
			"arguments": arguments,
		},
	}

	resp, err := c.SendRequest(req)
	if err != nil {
		return nil, err
	}

	// Log response details for debugging
	fmt.Printf("CallTool response - ID: %d, Error: %v, Result length: %d\n",
		resp.ID, resp.Error, len(resp.Result))

	return resp, nil
}

func (c *MCPClient) Close() error {
	c.stdin.Close()
	c.stdout.Close()
	return c.cmd.Wait()
}
