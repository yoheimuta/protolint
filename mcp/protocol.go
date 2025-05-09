// Package mcp implements the Model Context Protocol (MCP) server for protolint.
package mcp

import (
	"encoding/json"
)

// Request represents a JSON-RPC 2.0 request
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	ID      interface{}     `json:"id"`
}

// Response represents a JSON-RPC 2.0 response
type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *Error      `json:"error,omitempty"`
	ID      interface{} `json:"id"`
}

// Error represents a JSON-RPC 2.0 error
type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ClientInfo represents information about the client
type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// ClientCapabilities represents client capabilities
type ClientCapabilities struct {
	Roots        map[string]interface{} `json:"roots,omitempty"`
	Sampling     map[string]interface{} `json:"sampling,omitempty"`
	Experimental map[string]interface{} `json:"experimental,omitempty"`
}

// InitializeParams represents the parameters for initialize request
type InitializeParams struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ClientCapabilities `json:"capabilities"`
	ClientInfo      ClientInfo         `json:"clientInfo,omitempty"`
}

// ServerInfo represents information about the server
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// ServerCapabilities represents the server's capabilities
type ServerCapabilities struct {
	Prompts      map[string]interface{} `json:"prompts,omitempty"`
	Resources    map[string]interface{} `json:"resources,omitempty"`
	Tools        map[string]interface{} `json:"tools,omitempty"`
	Logging      map[string]interface{} `json:"logging,omitempty"`
	Experimental map[string]interface{} `json:"experimental,omitempty"`
}

// InitializeResult represents the response for initialize request
type InitializeResult struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ServerCapabilities `json:"capabilities"`
	ServerInfo      ServerInfo         `json:"serverInfo"`
	Instructions    string             `json:"instructions,omitempty"`
}

// ListToolsResponse represents the response for list_tools request
type ListToolsResponse struct {
	Tools []ToolInfo `json:"tools"`
}

// ToolInfo represents information about a tool
type ToolInfo struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema,omitempty"`
}

// CallToolPayload represents the payload for call_tool request
type CallToolPayload struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

// ContentItem represents a content item in a tool result
type ContentItem struct {
	Type     string      `json:"type"`
	Text     string      `json:"text,omitempty"`
	Data     string      `json:"data,omitempty"`
	MimeType string      `json:"mimeType,omitempty"`
	Resource interface{} `json:"resource,omitempty"`
}

// CallToolResponse represents the response for call_tool request
type CallToolResponse struct {
	Content []ContentItem `json:"content"`
	IsError bool          `json:"isError"`
}
