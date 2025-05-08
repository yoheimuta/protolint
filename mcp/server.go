// Package mcp implements the Model Context Protocol (MCP) server for protolint.
package mcp

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/yoheimuta/protolint/internal/osutil"
)

// Server represents an MCP server
type Server struct {
	tools  []Tool
	stdout io.Writer
	stderr io.Writer
}

// NewServer creates a new MCP server
func NewServer(stdout, stderr io.Writer) *Server {
	return &Server{
		tools: []Tool{
			NewLintFilesTool(),
			// Other tools can be added in the future
		},
		stdout: stdout,
		stderr: stderr,
	}
}

// Run starts the MCP server
func (s *Server) Run() osutil.ExitCode {
	_, _ = fmt.Fprintf(s.stderr, "protolint MCP server is running. cwd: %s\n", getCurrentDir())

	decoder := json.NewDecoder(os.Stdin)
	encoder := json.NewEncoder(s.stdout)

	for {
		var request Request
		if err := decoder.Decode(&request); err != nil {
			if err == io.EOF {
				// Normal termination
				return osutil.ExitSuccess
			}
			_, _ = fmt.Fprintf(s.stderr, "Error decoding request: %v\n", err)
			return osutil.ExitInternalFailure
		}

		// Ensure JSONRPC version is set
		if request.JSONRPC == "" {
			request.JSONRPC = "2.0"
		}

		// Process the request
		response := s.handleRequest(&request)

		// Only encode and send a response if there is one
		// Notifications don't require a response
		if response != nil {
			if err := encoder.Encode(response); err != nil {
				_, _ = fmt.Fprintf(s.stderr, "Error encoding response: %v\n", err)
				return osutil.ExitInternalFailure
			}
		}
	}
}

// handleRequest handles a single JSON-RPC request
func (s *Server) handleRequest(req *Request) *Response {
	switch req.Method {
	case "initialize":
		return s.handleInitialize(req)
	case "notifications/initialized":
		// This is a notification, so we don't need to send a response
		_, _ = fmt.Fprintf(s.stderr, "Received initialized notification, client is ready\n")
		return nil
	case "tools/list":
		return s.handleToolsList(req)
	case "tools/call":
		return s.handleToolsCall(req)
	default:
		return &Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &Error{
				Code:    -32601, // Method not found
				Message: fmt.Sprintf("Method not found: %s", req.Method),
			},
		}
	}
}

// handleInitialize handles initialization requests according to the MCP protocol
func (s *Server) handleInitialize(req *Request) *Response {
	// Parse the initialize parameters
	var params InitializeParams
	if req.Params != nil && len(req.Params) > 0 {
		if err := json.Unmarshal(req.Params, &params); err != nil {
			_, _ = fmt.Fprintf(s.stderr, "Warning: failed to parse initialize params: %v\n", err)
			return &Response{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error: &Error{
					Code:    -32602, // Invalid params
					Message: fmt.Sprintf("Failed to parse initialize params: %v", err),
				},
			}
		}
	}

	// Log information about the client
	if params.ClientInfo.Name != "" {
		_, _ = fmt.Fprintf(s.stderr, "Client info: %s %s\n", params.ClientInfo.Name, params.ClientInfo.Version)
	}

	// Log the client's protocol version
	if params.ProtocolVersion != "" {
		_, _ = fmt.Fprintf(s.stderr, "Client protocol version: %s\n", params.ProtocolVersion)
	}

	// Check if we support the client's protocol version
	// We only support 2024-11-05
	supportedVersion := "2024-11-05"
	if params.ProtocolVersion != "" && params.ProtocolVersion != supportedVersion {
		_, _ = fmt.Fprintf(s.stderr, "Warning: Client requested protocol version %s, but we're responding with %s\n",
			params.ProtocolVersion, supportedVersion)
		// We still continue with our supported version - the client will judge compatibility
	}

	// Create initialize result with proper MCP protocol capabilities
	result := InitializeResult{
		ProtocolVersion: supportedVersion,
		ServerInfo: ServerInfo{
			Name:    "protolint-mcp",
			Version: "1.0.0",
		},
		Capabilities: ServerCapabilities{
			// We only support tools with listChanged capability
			Tools: map[string]interface{}{
				"listChanged": true,
			},
		},
		Instructions: "Protocol Buffer linter for enforcing proto style guide rules",
	}

	return &Response{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}
}

// handleToolsList handles tools/list request
func (s *Server) handleToolsList(req *Request) *Response {
	toolInfos := make([]ToolInfo, 0, len(s.tools))
	for _, tool := range s.tools {
		toolInfos = append(toolInfos, tool.GetInfo())
	}

	return &Response{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: &ListToolsResponse{
			Tools: toolInfos,
		},
	}
}

// handleToolsCall handles tools/call request
func (s *Server) handleToolsCall(req *Request) *Response {
	var payload CallToolPayload
	if err := json.Unmarshal(req.Params, &payload); err != nil {
		return &Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &Error{
				Code:    -32700, // Parse error
				Message: fmt.Sprintf("Invalid payload: %v", err),
			},
		}
	}

	// Find the tool
	var tool Tool
	for _, t := range s.tools {
		if t.GetInfo().Name == payload.Name {
			tool = t
			break
		}
	}

	if tool == nil {
		return &Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &Error{
				Code:    -32601, // Method not found
				Message: fmt.Sprintf("Tool not found: %s", payload.Name),
			},
		}
	}

	// Execute the tool
	result, err := tool.Execute(payload.Arguments)
	if err != nil {
		return &Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &Error{
				Code:    -32000, // Server error
				Message: fmt.Sprintf("Tool execution failed: %v", err),
			},
		}
	}

	// Convert the result to a text content item
	resultText, err := json.Marshal(result)
	if err != nil {
		return &Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &Error{
				Code:    -32000, // Server error
				Message: fmt.Sprintf("Failed to marshal result: %v", err),
			},
		}
	}

	return &Response{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: &CallToolResponse{
			Content: []ContentItem{
				{
					Type: "text",
					Text: string(resultText),
				},
			},
			IsError: false,
		},
	}
}

// getCurrentDir returns the current working directory
func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "<unknown>"
	}
	return dir
}
