package main

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/Kong/go-pdk"
)

func handleHTTPProxy(kong *pdk.PDK, conf Config) {
	proxyURL := conf.ProxyURL
	if proxyURL == "" {
		proxyURL = "http://localhost:9000"
	}

	method, _ := kong.Request.GetMethod()
	path, _ := kong.Request.GetPath()
	body, _ := kong.Request.GetRawBody()
	headers, _ := kong.Request.GetHeaders(-1) // -1 means get all headers

	// Create HTTP client
	client := &http.Client{
		Timeout: time.Duration(conf.Timeout) * time.Millisecond,
	}

	// Create request to proxy
	req, err := http.NewRequest(method, proxyURL+path, nil)
	if err != nil {
		kong.Log.Err("Failed to create proxy request: ", err.Error())
		kong.Response.Exit(500, []byte(`{"error":"Internal server error"}`), nil)
		return
	}

	// Set body for POST/PUT requests
	if len(body) > 0 {
		req.Body = io.NopCloser(bytes.NewReader(body))
		req.ContentLength = int64(len(body))
	}

	// Forward headers
	for key, values := range headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Make request to proxy
	resp, err := client.Do(req)
	if err != nil {
		kong.Log.Err("Proxy request failed: ", err.Error())
		kong.Response.Exit(502, []byte(`{"error":"Bad gateway"}`), nil)
		return
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		kong.Log.Err("Failed to read proxy response: ", err.Error())
		kong.Response.Exit(500, []byte(`{"error":"Internal server error"}`), nil)
		return
	}

	// Forward response headers
	respHeaders := make(map[string][]string)
	for key, values := range resp.Header {
		respHeaders[key] = values
	}

	// Return proxy response
	kong.Response.Exit(resp.StatusCode, respBody, respHeaders)
}
