# protolint MCP Server

This document describes protolint's implementation of the [Model Context Protocol (MCP)](https://modelcontextprotocol.io), which allows AI models like Claude to interact with protolint directly.

## Overview

The Model Context Protocol (MCP) is a standardized protocol that facilitates communication between AI assistants and external tools. By implementing MCP, protolint can be used directly by AI assistants to lint Protocol Buffer files.

## Usage

To start protolint in MCP server mode:

```sh
protolint --mcp
```

This starts protolint as an MCP server, which listens for commands via stdin and writes responses to stdout.

## Integrating with Claude Desktop

To use protolint with Claude Desktop:

1. Edit the Claude Desktop configuration file:
   - macOS: `~/Library/Application Support/Claude/claude_desktop_config.json`
   - Windows: `%APPDATA%\Claude\claude_desktop_config.json`

2. Add the following configuration:
   ```json
   {
     "mcpServers": {
       "protolint": {
         "command": "/path/to/protolint",
         "args": ["--mcp"],
         "type": "stdio"
       }
     }
   }
   ```

3. Replace `/path/to/protolint` with the absolute path to your protolint executable.

4. Restart Claude Desktop.

## Protocol Implementation

The MCP server implementation follows the [Model Context Protocol specification](https://modelcontextprotocol.io) and uses JSON-RPC 2.0 for communication:

1. **Protocol Version**: The server supports version "2024-11-05" of the MCP protocol. If a client requests a different version, the server will respond with this supported version as specified in the protocol's version negotiation mechanism.

2. **Server Information**: The server identifies itself as "protolint-mcp" with version "1.0.0".

3. **Communication**: Uses stdio for communication between the client and server.

4. **Request Methods**:
   - `initialize`: Initializes the connection and negotiates protocol version
   - `notifications/initialized`: Notification that the client is ready
   - `tools/list`: Returns a list of available tools
   - `tools/call`: Calls a specific tool with arguments

5. **Response Format**: All responses follow the JSON-RPC 2.0 format with appropriate result or error fields.

## Available Tools

When running in MCP mode, protolint provides the following tools:

### lint-files

Lint Protocol Buffer files using protolint.

**Arguments:**
- `files`: (required) An array of file paths to lint
- `config_path`: (optional) Path to protolint config file
- `fix`: (optional) Fix lint errors if possible

**Example request:**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "id": "request-1234",
  "params": {
    "name": "lint-files",
    "arguments": {
      "files": ["/path/to/file1.proto", "/path/to/file2.proto"],
      "config_path": "/path/to/config.yaml",
      "fix": true
    }
  }
}
```

**Example response:**
```json
{
  "jsonrpc": "2.0",
  "id": "request-1234",
  "result": {
    "content": [
      {
        "type": "text",
        "text": "{\"exit_code\":0,\"results\":[{\"file_path\":\"/path/to/file1.proto\",\"failures\":[{\"rule_id\":\"ENUM_NAMES_UPPER_CAMEL_CASE\",\"message\":\"Enum name must be UpperCamelCase\",\"line\":5,\"column\":6,\"severity\":\"error\"}]}]}"
      }
    ],
    "isError": false
  }
}
```

## Example Usage in Claude Desktop

Once configured, you can ask Claude to lint your Protocol Buffer files:

```
Can you lint my protocol buffer files /path/to/file1.proto and /path/to/file2.proto?
```

Claude will use protolint to analyze the files and report any issues it finds.

## Exit Codes

The `lint-files` tool returns the following exit codes:

- `0`: No lint errors found
- `1`: Lint errors found
- `2`: Internal error occurred

## Development

The MCP implementation is located in the `mcp` directory and consists of:

- `protocol.go`: Protocol message definitions and JSON-RPC 2.0 structures
- `server.go`: The MCP server implementation with request handling
- `tools.go`: Tool implementations, currently only the lint-files tool

The server uses the MCP reporter for output formatting, which is configured when executing the lint command.

## Troubleshooting

If you encounter issues with the MCP server:

1. Check that the protolint executable is in your PATH or use an absolute path in the configuration.
2. Verify that the configuration file is correctly formatted.
3. Restart Claude Desktop after making changes to the configuration.
4. Check if there are any error messages in the Claude Desktop logs.
5. Ensure the protocol version in your client configuration matches "2024-11-05" or is omitted to allow version negotiation.
