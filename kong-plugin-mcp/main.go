package main

import (
	"encoding/json"
	"time"

	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
)

const Version = "1.0.0"
const Priority = 1000

type MCPToolResult struct {
	Content           []ContentItem          `json:"content"`
	StructuredContent map[string]interface{} `json:"structuredContent"`
}

type ContentItem struct {
	Type string      `json:"type"`
	Text string      `json:"text,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

func main() {
	server.StartServer(New, Version, Priority)
}

func New() interface{} {
	return &Config{}
}

func (conf Config) Access(kong *pdk.PDK) {
	// Check if we're in HTTP proxy mode
	if conf.Mode == "http-proxy" {
		handleHTTPProxy(kong, conf)
		return
	}

	// Default to tool mode
	toolName := conf.ToolName
	if toolName == "" {
		toolName = "get_wordle_suggestions"
	}

	body, err := kong.Request.GetRawBody()
	if err != nil {
		kong.Log.Err("Failed to read request body: ", err.Error())
		errorBody, _ := json.Marshal(map[string]string{"error": "Invalid request body"})
		kong.Response.Exit(400, errorBody, nil)
		return
	}

	var arguments map[string]interface{}
	if err := json.Unmarshal([]byte(body), &arguments); err != nil {
		kong.Log.Err("Failed to parse JSON: ", err.Error())
		errorBody, _ := json.Marshal(map[string]string{"error": "Invalid JSON"})
		kong.Response.Exit(400, errorBody, nil)
		return
	}

	client, err := NewMCPClient(conf.MCPServerPath)
	if err != nil {
		kong.Log.Err("Failed to create MCP client: ", err.Error())
		errorBody, _ := json.Marshal(map[string]string{"error": "Failed to start MCP server"})
		kong.Response.Exit(500, errorBody, nil)
		return
	}
	defer client.Close()

	if err := client.Initialize(); err != nil {
		kong.Log.Err("Failed to initialize MCP: ", err.Error())
		errorBody, _ := json.Marshal(map[string]string{"error": "Failed to initialize MCP server"})
		kong.Response.Exit(500, errorBody, nil)
		return
	}

	resultChan := make(chan *JSONRPCResponse, 1)
	errChan := make(chan error, 1)

	go func() {
		resp, err := client.CallTool(conf.ToolName, arguments)
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- resp
	}()

	timeout := time.Duration(conf.Timeout) * time.Millisecond
	select {
	case resp := <-resultChan:
		if resp.Error != nil {
			errorBody, _ := json.Marshal(map[string]string{"error": resp.Error.Message})
			kong.Response.Exit(500, errorBody, nil)
			return
		}

		// Check if we got a response
		if len(resp.Result) == 0 {
			kong.Log.Err("Empty MCP response")
			errorBody, _ := json.Marshal(map[string]string{"error": "Empty MCP response"})
			kong.Response.Exit(500, errorBody, nil)
			return
		}

		kong.Log.Info("MCP Response length: ", len(resp.Result))

		// Parse the MCP tool result
		var result MCPToolResult
		if err := json.Unmarshal(resp.Result, &result); err != nil {
			kong.Log.Err("Failed to parse result: ", err.Error())
			errorBody, _ := json.Marshal(map[string]string{
				"error":   "Invalid MCP response",
				"details": err.Error(),
			})
			kong.Response.Exit(500, errorBody, nil)
			return
		}

		// Use structuredContent if available
		var responseBody map[string]interface{}
		if len(result.StructuredContent) > 0 {
			responseBody = result.StructuredContent
		} else if len(result.Content) > 0 {
			// Fallback: extract data from content items
			responseBody = make(map[string]interface{})
			for _, item := range result.Content {
				if item.Type == "text" && item.Text != "" {
					var textData map[string]interface{}
					if err := json.Unmarshal([]byte(item.Text), &textData); err == nil {
						for k, v := range textData {
							responseBody[k] = v
						}
					} else {
						responseBody["text"] = item.Text
					}
				}
			}
		} else {
			// Return the whole result if nothing else works
			kong.Log.Warn("No structured content or content array found")
			responseBody = map[string]interface{}{"result": result}
		}

		responseBytes, _ := json.Marshal(responseBody)
		kong.Response.Exit(200, responseBytes, nil)

	case err := <-errChan:
		kong.Log.Err("Tool call failed: ", err.Error())
		errorBody, _ := json.Marshal(map[string]string{"error": "Tool execution failed"})
		kong.Response.Exit(500, errorBody, nil)

	case <-time.After(timeout):
		kong.Log.Err("Tool call timed out")
		errorBody, _ := json.Marshal(map[string]string{"error": "Request timeout"})
		kong.Response.Exit(504, errorBody, nil)
	}
}
