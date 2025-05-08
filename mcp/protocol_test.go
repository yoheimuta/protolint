package mcp

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestRequest_Unmarshal(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    Request
		wantErr bool
	}{
		{
			name: "list_tools request",
			json: `{"jsonrpc":"2.0","method":"list_tools","id":"request-1234"}`,
			want: Request{
				JSONRPC: "2.0",
				Method:  "list_tools",
				ID:      "request-1234",
			},
			wantErr: false,
		},
		{
			name: "call_tool request",
			json: `{"jsonrpc":"2.0","method":"call_tool","id":"request-5678","params":{"name":"lint-files","arguments":{"files":["/path/to/file.proto"]}}}`,
			want: Request{
				JSONRPC: "2.0",
				Method:  "call_tool",
				ID:      "request-5678",
				Params:  json.RawMessage(`{"name":"lint-files","arguments":{"files":["/path/to/file.proto"]}}`),
			},
			wantErr: false,
		},
		{
			name: "request with numeric ID",
			json: `{"jsonrpc":"2.0","method":"list_tools","id":42}`,
			want: Request{
				JSONRPC: "2.0",
				Method:  "list_tools",
				ID:      float64(42), // JSON numbers are unmarshaled as float64
			},
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			json:    `{"jsonrpc":"2.0","method":"list_tools","id":`,
			want:    Request{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Request
			err := json.Unmarshal([]byte(tt.json), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if got.JSONRPC != tt.want.JSONRPC || got.Method != tt.want.Method || got.ID != tt.want.ID {
				t.Errorf("Unmarshal() = %v, want %v", got, tt.want)
			}

			// For params, we need to compare the string representation
			if !reflect.DeepEqual(string(got.Params), string(tt.want.Params)) {
				t.Errorf("Params = %v, want %v", string(got.Params), string(tt.want.Params))
			}
		})
	}
}

func TestResponse_Marshal(t *testing.T) {
	tests := []struct {
		name     string
		response Response
		want     string
		wantErr  bool
	}{
		{
			name: "list_tools_response",
			response: Response{
				JSONRPC: "2.0",
				ID:      "request-1234",
				Result: &ListToolsResponse{
					Tools: []ToolInfo{
						{
							Name:        "lint-files",
							Description: "Lint Protocol Buffer files using protolint",
							InputSchema: map[string]interface{}{
								"type": "object",
								"properties": map[string]interface{}{
									"files": map[string]interface{}{
										"type": "array",
										"items": map[string]interface{}{
											"type": "string",
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "call_tool_response",
			response: Response{
				JSONRPC: "2.0",
				ID:      "request-5678",
				Result: &CallToolResponse{
					Content: []ContentItem{
						{
							Type: "text",
							Text: `{"exit_code":0,"results":[{"file_path":"/path/to/file.proto","failures":[]}]}`,
						},
					},
					IsError: false,
				},
			},
			wantErr: false,
		},
		{
			name: "response with numeric ID",
			response: Response{
				JSONRPC: "2.0",
				ID:      float64(42),
				Result: &ListToolsResponse{
					Tools: []ToolInfo{},
				},
			},
			wantErr: false,
		},
		{
			name: "error response",
			response: Response{
				JSONRPC: "2.0",
				ID:      "request-9012",
				Error: &Error{
					Code:    -32000,
					Message: "Tool not found: invalid-tool",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.response)

			if (err != nil) != tt.wantErr {
				t.Errorf("json.Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// Deserialize and compare
			var gotResponse, wantResponse map[string]interface{}
			if err := json.Unmarshal(got, &gotResponse); err != nil {
				t.Errorf("json.Unmarshal() error = %v", err)
				return
			}

			wantJSON, _ := json.Marshal(tt.response)
			if err := json.Unmarshal(wantJSON, &wantResponse); err != nil {
				t.Errorf("json.Unmarshal() error = %v", err)
				return
			}

			// Check jsonrpc and id
			if gotResponse["jsonrpc"] != wantResponse["jsonrpc"] || gotResponse["id"] != wantResponse["id"] {
				t.Errorf("Marshal() = %v, want %v", string(got), string(wantJSON))
			}

			// Check result or error
			if tt.response.Result != nil {
				if !reflect.DeepEqual(gotResponse["result"], wantResponse["result"]) {
					t.Errorf("Result mismatch:\ngot:  %v\nwant: %v", gotResponse["result"], wantResponse["result"])
				}
			} else if tt.response.Error != nil {
				if !reflect.DeepEqual(gotResponse["error"], wantResponse["error"]) {
					t.Errorf("Error mismatch:\ngot:  %v\nwant: %v", gotResponse["error"], wantResponse["error"])
				}
			}
		})
	}
}

func TestRequest_ParseParams(t *testing.T) {
	tests := []struct {
		name    string
		request Request
		target  interface{}
		wantErr bool
	}{
		{
			name: "parse call_tool params",
			request: Request{
				JSONRPC: "2.0",
				Method:  "call_tool",
				ID:      "request-1234",
				Params:  json.RawMessage(`{"name":"lint-files","arguments":{"files":["/path/to/file.proto"]}}`),
			},
			target:  &CallToolPayload{},
			wantErr: false,
		},
		{
			name: "parse invalid params",
			request: Request{
				JSONRPC: "2.0",
				Method:  "call_tool",
				ID:      "request-1234",
				Params:  json.RawMessage(`{"name":"lint-files","arguments":invalid}`),
			},
			target:  &CallToolPayload{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := json.Unmarshal(tt.request.Params, tt.target)

			if (err != nil) != tt.wantErr {
				t.Errorf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// Check if the target was populated
			payload, ok := tt.target.(*CallToolPayload)
			if !ok {
				t.Fatalf("Expected target to be *CallToolPayload")
			}

			if payload.Name == "" {
				t.Errorf("Expected Name to be populated, got empty string")
			}

			if len(payload.Arguments) == 0 {
				t.Errorf("Expected Arguments to be populated, got empty")
			}
		})
	}
}

func TestInitializeParams_Unmarshal(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    InitializeParams
		wantErr bool
	}{
		{
			name: "basic params",
			json: `{"protocolVersion":"2025-03-26","capabilities":{"roots":{},"sampling":{}}}`,
			want: InitializeParams{
				ProtocolVersion: "2025-03-26",
				Capabilities: ClientCapabilities{
					Roots:    map[string]interface{}{},
					Sampling: map[string]interface{}{},
				},
			},
			wantErr: false,
		},
		{
			name: "with client info",
			json: `{"protocolVersion":"2025-03-26","capabilities":{"roots":{}},"clientInfo":{"name":"TestClient","version":"1.0.0"}}`,
			want: InitializeParams{
				ProtocolVersion: "2025-03-26",
				Capabilities: ClientCapabilities{
					Roots: map[string]interface{}{},
				},
				ClientInfo: ClientInfo{
					Name:    "TestClient",
					Version: "1.0.0",
				},
			},
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			json:    `{"protocolVersion":"2025-03-26","capabilities":`,
			want:    InitializeParams{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got InitializeParams
			err := json.Unmarshal([]byte(tt.json), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if got.ProtocolVersion != tt.want.ProtocolVersion {
				t.Errorf("ProtocolVersion = %v, want %v", got.ProtocolVersion, tt.want.ProtocolVersion)
			}

			// Check ClientInfo if provided
			if tt.want.ClientInfo.Name != "" {
				if got.ClientInfo.Name != tt.want.ClientInfo.Name {
					t.Errorf("ClientInfo.Name = %v, want %v", got.ClientInfo.Name, tt.want.ClientInfo.Name)
				}
				if got.ClientInfo.Version != tt.want.ClientInfo.Version {
					t.Errorf("ClientInfo.Version = %v, want %v", got.ClientInfo.Version, tt.want.ClientInfo.Version)
				}
			}
		})
	}
}

func TestError_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		error   Error
		wantErr bool
	}{
		{
			name: "basic error",
			error: Error{
				Code:    -32000,
				Message: "Server error",
			},
			wantErr: false,
		},
		{
			name: "error with data",
			error: Error{
				Code:    -32602,
				Message: "Invalid params",
				Data:    "Missing required field",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.error)

			if (err != nil) != tt.wantErr {
				t.Errorf("json.Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// Deserialize and compare
			var gotError map[string]interface{}
			if err := json.Unmarshal(got, &gotError); err != nil {
				t.Errorf("json.Unmarshal() error = %v", err)
				return
			}

			// Check code and message
			if int(gotError["code"].(float64)) != tt.error.Code {
				t.Errorf("Expected code to be %d, got %v", tt.error.Code, gotError["code"])
			}

			if gotError["message"] != tt.error.Message {
				t.Errorf("Expected message to be %s, got %v", tt.error.Message, gotError["message"])
			}

			// Check data if present
			if tt.error.Data != nil {
				if gotError["data"] == nil {
					t.Errorf("Expected data to be non-nil")
				}
			}
		})
	}
}
