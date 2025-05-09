package mcp

import (
	"io"
	"reflect"
	"testing"
)

func TestServer_handleInitialize_Success(t *testing.T) {
	server := NewServer(io.Discard, io.Discard)
	req := &Request{
		JSONRPC: "2.0",
		Method:  "initialize",
		Params: []byte(`{
			"protocolVersion": "2024-11-05",
			"capabilities": {
				"roots": {},
				"sampling": {}
			},
			"clientInfo": {
				"name": "TestClient",
				"version": "1.0.0"
			}
		}`),
		ID: float64(0),
	}

	resp := server.handleInitialize(req)

	if resp.JSONRPC != "2.0" {
		t.Errorf("Expected JSONRPC version to be '2.0', got '%s'", resp.JSONRPC)
	}

	if resp.ID != float64(0) {
		t.Errorf("Expected response ID 0, got '%v'", resp.ID)
	}

	if resp.Error != nil {
		t.Errorf("Expected Error to be nil, got %v", resp.Error)
	}

	result, ok := resp.Result.(InitializeResult)
	if !ok {
		t.Fatalf("Expected Result to be InitializeResult, got %T", resp.Result)
	}

	// Check protocol version
	if result.ProtocolVersion != "2024-11-05" {
		t.Errorf("Expected protocolVersion to be '2024-11-05', got '%v'", result.ProtocolVersion)
	}

	// Check server info
	if result.ServerInfo.Name == "" {
		t.Errorf("Expected serverInfo.name to be non-empty")
	}
	if result.ServerInfo.Version == "" {
		t.Errorf("Expected serverInfo.version to be non-empty")
	}

	// Check capabilities
	if result.Capabilities.Tools != nil {
		t.Errorf("Expected capabilities.tools to be nil")
	}
}

func TestServer_handleInitialize_DifferentVersion(t *testing.T) {
	server := NewServer(io.Discard, io.Discard)
	req := &Request{
		JSONRPC: "2.0",
		Method:  "initialize",
		Params: []byte(`{
			"protocolVersion": "2025-03-26",
			"capabilities": {
				"roots": {},
				"sampling": {}
			},
			"clientInfo": {
				"name": "TestClient",
				"version": "1.0.0"
			}
		}`),
		ID: float64(1),
	}

	resp := server.handleInitialize(req)

	if resp.JSONRPC != "2.0" {
		t.Errorf("Expected JSONRPC version to be '2.0', got '%s'", resp.JSONRPC)
	}

	if resp.ID != float64(1) {
		t.Errorf("Expected response ID 1, got '%v'", resp.ID)
	}

	// We should NOT get an error for unsupported version
	if resp.Error != nil {
		t.Fatalf("Expected Error to be nil for different version, got %v", resp.Error)
	}

	// Instead, we should get a result with our supported version
	result, ok := resp.Result.(InitializeResult)
	if !ok {
		t.Fatalf("Expected Result to be InitializeResult, got %T", resp.Result)
	}

	// Check protocol version is our supported version
	if result.ProtocolVersion != "2024-11-05" {
		t.Errorf("Expected protocolVersion to be '2024-11-05', got '%v'", result.ProtocolVersion)
	}
}

func TestServer_handleListTools(t *testing.T) {
	server := NewServer(io.Discard, io.Discard)
	req := &Request{
		JSONRPC: "2.0",
		Method:  "list_tools",
		ID:      "test-1",
	}

	resp := server.handleToolsList(req)

	if resp.JSONRPC != "2.0" {
		t.Errorf("Expected JSONRPC version to be '2.0', got '%s'", resp.JSONRPC)
	}

	if resp.ID != "test-1" {
		t.Errorf("Expected response ID 'test-1', got '%v'", resp.ID)
	}

	if resp.Error != nil {
		t.Errorf("Expected Error to be nil, got %v", resp.Error)
	}

	result, ok := resp.Result.(*ListToolsResponse)
	if !ok {
		t.Fatalf("Expected Result to be *ListToolsResponse")
	}

	if len(result.Tools) == 0 {
		t.Error("Expected at least one tool")
	}

	// Check if lint-files tool is present
	foundLintFiles := false
	for _, tool := range result.Tools {
		if tool.Name == "lint-files" {
			foundLintFiles = true
			break
		}
	}

	if !foundLintFiles {
		t.Error("Expected to find lint-files tool")
	}
}

func TestServer_handleInitializedNotification(t *testing.T) {
	server := NewServer(io.Discard, io.Discard)
	req := &Request{
		JSONRPC: "2.0",
		Method:  "notifications/initialized",
		Params:  nil,
		ID:      nil, // Notifications don't have IDs
	}

	// For notifications, we expect nil response
	resp := server.handleRequest(req)
	if resp != nil {
		t.Errorf("Expected nil response for notification, got %v", resp)
	}
}

func TestServer_handleToolsList(t *testing.T) {
	server := NewServer(io.Discard, io.Discard)
	req := &Request{
		JSONRPC: "2.0",
		Method:  "tools/list",
		ID:      "test-alt",
	}

	resp := server.handleRequest(req)

	if resp.JSONRPC != "2.0" {
		t.Errorf("Expected JSONRPC version to be '2.0', got '%s'", resp.JSONRPC)
	}

	if resp.ID != "test-alt" {
		t.Errorf("Expected response ID 'test-alt', got '%v'", resp.ID)
	}

	if resp.Error != nil {
		t.Errorf("Expected Error to be nil, got %v", resp.Error)
	}

	result, ok := resp.Result.(*ListToolsResponse)
	if !ok {
		t.Fatalf("Expected Result to be *ListToolsResponse, got %T", resp.Result)
	}

	if len(result.Tools) == 0 {
		t.Error("Expected at least one tool")
	}

	// Check if lint-files tool is present
	foundLintFiles := false
	for _, tool := range result.Tools {
		if tool.Name == "lint-files" {
			foundLintFiles = true
			break
		}
	}

	if !foundLintFiles {
		t.Error("Expected to find lint-files tool")
	}
}

func TestServer_handleCallTool(t *testing.T) {
	tests := []struct {
		name      string
		request   *Request
		wantError bool
	}{
		{
			name: "tool not found",
			request: &Request{
				JSONRPC: "2.0",
				Method:  "call_tool",
				ID:      "test-1",
				Params: []byte(`{
					"name": "non-existent-tool",
					"arguments": {}
				}`),
			},
			wantError: true,
		},
		{
			name: "invalid payload",
			request: &Request{
				JSONRPC: "2.0",
				Method:  "call_tool",
				ID:      "test-2",
				Params: []byte(`{
					"invalid": "json"
				}`),
			},
			wantError: true,
		},
		// Note: We can't easily test a successful call_tool because it would require
		// setting up mock files and overriding the lib.Lint function.
		// That would be better suited for an integration test.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewServer(io.Discard, io.Discard)
			resp := server.handleToolsCall(tt.request)

			if resp.JSONRPC != "2.0" {
				t.Errorf("Expected JSONRPC version to be '2.0', got '%s'", resp.JSONRPC)
			}

			if resp.ID != tt.request.ID {
				t.Errorf("Expected response ID '%v', got '%v'", tt.request.ID, resp.ID)
			}

			// Check if error response as expected
			if tt.wantError {
				if resp.Error == nil {
					t.Fatalf("Expected Error to be non-nil")
				}

				if resp.Error.Message == "" {
					t.Error("Expected non-empty error message")
				}

				if resp.Result != nil {
					t.Errorf("Expected Result to be nil for error response")
				}
			} else {
				if resp.Error != nil {
					t.Errorf("Expected Error to be nil, got %v", resp.Error)
				}

				if resp.Result == nil {
					t.Errorf("Expected Result to be non-nil")
				}
			}
		})
	}
}

func TestServer_handleRequest(t *testing.T) {
	tests := []struct {
		name           string
		request        *Request
		wantError      bool
		wantResult     bool
		wantPayload    reflect.Type
		isNotification bool
	}{
		{
			name: "list_tools request",
			request: &Request{
				JSONRPC: "2.0",
				Method:  "tools/list",
				ID:      "test-1",
			},
			wantError:      false,
			wantResult:     true,
			wantPayload:    reflect.TypeOf(&ListToolsResponse{}),
			isNotification: false,
		},
		{
			name: "initialized notification",
			request: &Request{
				JSONRPC: "2.0",
				Method:  "notifications/initialized",
				ID:      nil,
			},
			wantError:      false,
			wantResult:     false,
			wantPayload:    nil,
			isNotification: true,
		},
		{
			name: "unknown method",
			request: &Request{
				JSONRPC: "2.0",
				Method:  "unknown_method",
				ID:      "test-2",
			},
			wantError:      true,
			wantResult:     false,
			wantPayload:    nil,
			isNotification: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewServer(io.Discard, io.Discard)
			resp := server.handleRequest(tt.request)

			// Notifications should return nil responses
			if tt.isNotification {
				if resp != nil {
					t.Errorf("Expected nil response for notification, got %v", resp)
				}
				return
			}

			// Regular requests should have responses
			if resp == nil {
				t.Fatalf("Expected non-nil response for regular request")
			}

			if resp.JSONRPC != "2.0" {
				t.Errorf("Expected JSONRPC version to be '2.0', got '%s'", resp.JSONRPC)
			}

			if resp.ID != tt.request.ID {
				t.Errorf("Expected response ID '%v', got '%v'", tt.request.ID, resp.ID)
			}

			// Check error
			if tt.wantError {
				if resp.Error == nil {
					t.Errorf("Expected Error to be non-nil")
				}
			} else {
				if resp.Error != nil {
					t.Errorf("Expected Error to be nil, got %v", resp.Error)
				}
			}

			// Check result
			if tt.wantResult {
				if resp.Result == nil {
					t.Errorf("Expected Result to be non-nil")
				} else {
					resultType := reflect.TypeOf(resp.Result)
					if resultType != tt.wantPayload {
						t.Errorf("Expected Result type '%v', got '%v'", tt.wantPayload, resultType)
					}
				}
			} else {
				if resp.Result != nil {
					t.Errorf("Expected Result to be nil, got %v", resp.Result)
				}
			}
		})
	}
}
