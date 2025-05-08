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
  "type": "call_tool",
  "id": "request-1234",
  "payload": {
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
  "type": "call_tool_response",
  "id": "request-1234",
  "payload": {
    "result": {
      "exit_code": 0,
      "results": [
        {
          "file_path": "/path/to/file1.proto",
          "failures": [
            {
              "rule_id": "ENUM_NAMES_UPPER_CAMEL_CASE",
              "message": "Enum name must be UpperCamelCase",
              "line": 5,
              "column": 6,
              "severity": "error"
            }
          ]
        }
      ]
    }
  }
}
```

## Example Usage in Claude Desktop

Once configured, you can ask Claude to lint your Protocol Buffer files:

```
Can you lint my protocol buffer files /path/to/file1.proto and /path/to/file2.proto?
```

Claude will use protolint to analyze the files and report any issues it finds.

## Protocol Implementation

The MCP server implementation follows the [Model Context Protocol specification](https://modelcontextprotocol.io) version 2024-11-05:

1. **Communication**: Uses stdio for communication between the client and server
2. **Request Types**:
   - `tools/list`: Returns a list of available tools
   - `tools/call`: Calls a specific tool with arguments
3. **Response Types**:
   - `list_tools_response`: Contains a list of available tools
   - `call_tool_response`: Contains the result of a tool call
   - `error`: Contains an error message

The server implements the protocol's version negotiation mechanism, responding with its supported version (2024-11-05) even if the client requests a different version. This follows the specification, which states that if the server doesn't support the requested version, it should respond with another version it does support.

## Exit Codes

The `lint-files` tool returns the following exit codes:

- `0`: No lint errors found
- `1`: Lint errors found
- `2`: Internal error occurred

## Development

The MCP implementation is located in the `mcp` directory and consists of:

- `server.go`: The MCP server implementation
- `protocol.go`: Protocol message definitions
- `tools.go`: Tool implementations

The reporter for MCP output format is in:
- `internal/linter/report/reporters/mcpReporter.go`

## Troubleshooting

If you encounter issues with the MCP server:

1. Check that the protolint executable is in your PATH or use an absolute path in the configuration.
2. Verify that the configuration file is correctly formatted.
3. Restart Claude Desktop after making changes to the configuration.
4. Check if there are any error messages in the Claude Desktop logs.
