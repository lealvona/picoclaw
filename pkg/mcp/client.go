package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// MCPTool represents a tool exposed via MCP
type MCPTool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// MCPResource represents a resource exposed via MCP
type MCPResource struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MimeType    string `json:"mimeType,omitempty"`
}

// MCPServer represents an MCP server configuration
type MCPServer struct {
	Name    string `json:"name"`
	Endpoint string `json:"endpoint"`
	Enabled  bool   `json:"enabled"`
	Tools    []MCPTool    `json:"tools"`
	Resources []MCPResource `json:"resources"`
}

// MCPClient handles communication with MCP servers
type MCPClIENT struct {
	servers map[string]*MCPServer
	client  *http.Client
	mu      sync.RWMutex
	timeout time.Duration
}

// NewMCPClient creates a new MCP client
func NewMCPClient() *MCPClIENT {
	return &MCPClIENT{
		servers: make(map[string]*MCPServer),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		timeout: 30 * time.Second,
	}
}

// AddServer adds an MCP server to the client
func (c *MCPClIENT) AddServer(name, endpoint string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	server := &MCPServer{
		Name:    name,
		Endpoint: endpoint,
		Enabled:  true,
	}

	// Discover tools and resources from the server
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	tools, err := c.discoverTools(ctx, endpoint)
	if err != nil {
		return fmt.Errorf("failed to discover tools from MCP server %s: %w", name, err)
	}
	server.Tools = tools

	resources, err := c.discoverResources(ctx, endpoint)
	if err == nil {
		server.Resources = resources
	}

	c.servers[name] = server
	return nil
}

// GetServers returns all configured MCP servers
func (c *MCPClIENT) GetServers() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	names := make([]string, 0, len(c.servers))
	for name := range c.servers {
		names = append(names, name)
	}
	return names
}

// GetTool returns a specific tool from MCP servers
func (c *MCPClIENT) GetTool(serverName, toolName string) (*MCPTool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	server, ok := c.servers[serverName]
	if !ok {
		return nil, fmt.Errorf("MCP server %s not found", serverName)
	}

	for _, tool := range server.Tools {
		if tool.Name == toolName {
			return &tool, nil
		}
	}

	return nil, fmt.Errorf("tool %s not found in MCP server %s", toolName, serverName)
}

// ExecuteTool executes a tool on an MCP server
func (c *MCPClIENT) ExecuteTool(ctx context.Context, serverName, toolName string, arguments map[string]interface{}) (string, error) {
	c.mu.RLock()
	server, ok := c.servers[serverName]
	c.mu.RUnlock()

	if !ok {
		return "", fmt.Errorf("MCP server %s not found", serverName)
	}

	if !server.Enabled {
		return "", fmt.Errorf("MCP server %s is disabled", serverName)
	}

	// Build request
	req := map[string]interface{}{
		"method": "tools/call",
		"params": map[string]interface{}{
			"name":      toolName,
			"arguments": arguments,
		},
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Send request to MCP server
	httpReq, err := http.NewRequestWithContext(ctx, "POST", server.Endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Body = io.NopCloser(jsonData)

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("MCP server returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var result struct {
		Result interface{} `json:"result"`
		Error  *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Error != nil {
		return "", fmt.Errorf("MCP error %d: %s", result.Error.Code, result.Error.Message)
	}

	resultJSON, _ := json.Marshal(result.Result)
	return string(resultJSON), nil
}

// discoverTools discovers available tools from an MCP server
func (c *MCPClIENT) discoverTools(ctx context.Context, endpoint string) ([]MCPTool, error) {
	req := map[string]interface{}{
		"method": "tools/list",
		"params": map[string]interface{}{},
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Body = io.NopCloser(jsonData)

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result struct {
		Result struct {
			Tools []MCPTool `json:"tools"`
		} `json:"result"`
		Error *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Error != nil {
		return nil, fmt.Errorf("MCP error %d: %s", result.Error.Code, result.Error.Message)
	}

	return result.Result.Tools, nil
}

// discoverResources discovers available resources from an MCP server
func (c *MCPClIENT) discoverResources(ctx context.Context, endpoint string) ([]MCPResource, error) {
	req := map[string]interface{}{
		"method": "resources/list",
		"params": map[string]interface{}{},
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Body = io.NopCloser(jsonData)

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result struct {
		Result struct {
			Resources []MCPResource `json:"resources"`
		} `json:"result"`
		Error *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Error != nil {
		return nil, fmt.Errorf("MCP error %d: %s", result.Error.Code, result.Error.Message)
	}

	return result.Result.Resources, nil
}

// RemoveServer removes an MCP server from the client
func (c *MCPClIENT) RemoveServer(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.servers, name)
}

// EnableServer enables an MCP server
func (c *MCPClIENT) EnableServer(name string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	server, ok := c.servers[name]
	if !ok {
		return fmt.Errorf("MCP server %s not found", name)
	}

	server.Enabled = true
	return nil
}

// DisableServer disables an MCP server
func (c *MCPClIENT) DisableServer(name string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	server, ok := c.servers[name]
	if !ok {
		return fmt.Errorf("MCP server %s not found", name)
	}

	server.Enabled = false
	return nil
}
