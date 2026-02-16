package tools

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/sipeed/picoclaw/pkg/mcp"
)

// MCPTool integrates Model Context Protocol servers into PicoClaw
type MCPTool struct {
	client   *mcp.MCPClIENT
	enabled  bool
	mu       sync.RWMutex
}

// NewMCPTool creates a new MCP tool
func NewMCPTool(client *mcp.MCPClIENT) *MCPTool {
	return &MCPTool{
		client:  client,
		enabled: true,
	}
}

// Name returns the tool name
func (t *MCPTool) Name() string {
	return "mcp"
}

// Description returns the tool description
func (t *MCPTool) Description() string {
	return "Execute tools from configured MCP (Model Context Protocol) servers. MCP allows PicoClaw to connect to external services like Google Drive, Slack, GitHub, databases, and more without writing custom code."
}

// Parameters returns the tool parameters schema
func (t *MCPTool) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"server": map[string]interface{}{
				"type":        "string",
				"description": "The MCP server name to use (e.g., 'github', 'slack', 'database')",
			},
			"tool": map[string]interface{}{
				"type":        "string",
				"description": "The tool name to execute on the MCP server",
			},
			"arguments": map[string]interface{}{
				"type":        "object",
				"description": "Arguments to pass to the MCP tool (as JSON object)",
			},
			"list_servers": map[string]interface{}{
				"type":        "boolean",
				"description": "If true, list all available MCP servers and their tools instead of executing a tool",
			},
		},
	}
}

// Execute executes the MCP tool
func (t *MCPTool) Execute(ctx context.Context, args map[string]interface{}) *ToolResult {
	t.mu.RLock()
	if !t.enabled {
		t.mu.RUnlock()
		return ErrorResult("MCP tool is disabled")
	}
	t.mu.RUnlock()

	// Check if user wants to list servers
	if listServers, ok := args["list_servers"].(bool); ok && listServers {
		return t.listServers()
	}

	// Get server name
	server, ok := args["server"].(string)
	if !ok || server == "" {
		return ErrorResult("server parameter is required")
	}

	// Get tool name
	toolName, ok := args["tool"].(string)
	if !ok || toolName == "" {
		return ErrorResult("tool parameter is required")
	}

	// Get arguments
	arguments, ok := args["arguments"].(map[string]interface{})
	if !ok {
		arguments = make(map[string]interface{})
	}

	// Execute tool on MCP server
	result, err := t.client.ExecuteTool(ctx, server, toolName, arguments)
	if err != nil {
		return &ToolResult{
			ForLLM:  fmt.Sprintf("Failed to execute MCP tool: %v", err),
			ForUser: fmt.Sprintf("❌ Failed to execute tool '%s' on server '%s': %v", toolName, server, err),
			IsError: true,
		}
	}

	return &ToolResult{
		ForLLM:  result,
		ForUser: fmt.Sprintf("✅ Executed tool '%s' on server '%s'\n\n%s", toolName, server, result),
		IsError: false,
	}
}

// listServers lists all configured MCP servers and their tools
func (t *MCPTool) listServers() *ToolResult {
	servers := t.client.GetServers()

	if len(servers) == 0 {
		return &ToolResult{
			ForLLM:  "No MCP servers configured",
			ForUser: "📋 No MCP servers are currently configured.\n\nTo add an MCP server, add it to your config.json under providers.mcp section.",
			IsError: false,
		}
	}

	var sb strings.Builder
	sb.WriteString("📋 Available MCP Servers:\n\n")

	for _, serverName := range servers {
		sb.WriteString(fmt.Sprintf("🔹 **%s**\n", serverName))

		// Get tools from this server
		// Note: We'll need to enhance the MCP client to expose this
		sb.WriteString("   Tools: (fetching...)\n")
		sb.WriteString("\n")
	}

	result := sb.String()
	return &ToolResult{
		ForLLM:  result,
		ForUser: result,
		IsError: false,
	}
}

// SetEnabled enables or disables the MCP tool
func (t *MCPTool) SetEnabled(enabled bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.enabled = enabled
}

// IsEnabled returns whether the MCP tool is enabled
func (t *MCPTool) IsEnabled() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.enabled
}
