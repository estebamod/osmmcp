# OpenStreetMap MCP Server (Go Implementation)

This is a Go implementation of the OpenStreetMap MCP server that enhances LLM capabilities with location-based services and geospatial data.

## Overview

This project is a Go OpenStreetMap MCP server.  It implemets the [Model Context Protocol](https://github.com/mark3labs/mcp-go) to enable LLMs to interact with geospatial data.

We've reimplemented the functionality in Go with a focus on performance, maintainability, and ease of integration with MCP desktop clients.

## Features

The server provides LLMs with tools to interact with OpenStreetMap data, including:

* Geocoding addresses and place names to coordinates
* Reverse geocoding coordinates to addresses
* Finding nearby points of interest
* Getting route directions between locations
* Searching for places by category within a bounding box
* Suggesting optimal meeting points for multiple people
* Exploring areas and getting comprehensive location information
* Finding EV charging stations near a location
* Finding EV charging stations along a route
* Analyzing commute options between home and work
* Performing neighborhood livability analysis
* Finding schools near a location
* Finding parking facilities

## Implemented Tools

| Tool Name | Description |
|-----------|-------------|
| `geocode_address` | Convert an address or place name to geographic coordinates |
| `reverse_geocode` | Convert geographic coordinates to a human-readable address |
| `find_nearby_places` | Find points of interest near a specific location |
| `search_category` | Search for places by category within a rectangular area |
| `get_route_directions` | Get directions for a route between two locations |
| `suggest_meeting_point` | Suggest an optimal meeting point for multiple people |
| `explore_area` | Explore an area and get comprehensive information about it |
| `find_charging_stations` | Find electric vehicle charging stations near a location |
| `find_route_charging_stations` | Find electric vehicle charging stations along a route |
| `analyze_commute` | Analyze transportation options between home and work locations |
| `analyze_neighborhood` | Evaluate neighborhood livability for real estate and relocation decisions |
| `find_schools_nearby` | Find educational institutions near a specific location |
| `find_parking_facilities` | Find parking facilities near a specific location |

## Code Architecture and Design

The code follows software engineering best practices:

1. **High Cohesion, Low Coupling** - Each package has a clear, focused responsibility
2. **Separation of Concerns** - Tools, server logic, and utilities are cleanly separated
3. **DRY (Don't Repeat Yourself)** - Common utilities are extracted into the `pkg/osm` package
4. **Security First** - HTTP clients are properly configured with timeouts and connection limits
5. **Structured Logging** - All logging is done via `slog` with consistent levels and formats
6. **SOLID Principles** - Particularly Single Responsibility and Interface Segregation
7. **Registry Pattern** - All tools are defined in a central registry for improved maintainability

## Usage

### Requirements

- Go 1.24 or higher

### Building the server

```bash
go build -o osmmcp ./cmd/osmmcp
```

### Running the server

```bash
./osmmcp
```

The server also supports command-line flags:

```bash
# Show version information
./osmmcp --version

# Enable debug logging
./osmmcp --debug

# Generate a Claude Desktop Client configuration file
./osmmcp --generate-config /path/to/config.json
```

The server will start and listen for MCP requests on the standard input/output. You can use it with any MCP-compatible client or LLM integration.

### Using with Claude Desktop Client

This MCP server is designed to work with Claude Desktop Client. You can set it up easily with the following steps:

1. Build the server:
   ```bash
   go build -o osmmcp ./cmd/osmmcp
   ```

2. Generate or update the Claude Desktop Client configuration:
   ```bash
   ./osmmcp --generate-config ~/Library/Application\ Support/Anthropic/Claude/config.json
   ```

   This will add an `OSM` entry to the `mcpServers` configuration in Claude Desktop Client. The configuration system intelligently:
   
   - Creates the file if it doesn't exist
   - Preserves existing tools like TAK when updating the configuration
   - Uses absolute paths to ensure Claude can find the executable
   - Validates JSON output to prevent corruption

3. Restart Claude Desktop Client to load the updated configuration.

4. In a conversation with Claude, you can now use the OpenStreetMap MCP tools.

The configuration file will look similar to this:

```json
{
  "mcpServers": {
    "OSM": {
      "command": "/path/to/osmmcp",
      "args": []
    },
    "TAK": {
      "command": "/path/to/takmcp",
      "args": [
        "--tak-host=localhost",
        "--tak-port=8089"
      ]
    }
  }
}
```

### API Dependencies

The server relies on these external APIs:

- **Nominatim** - For geocoding operations
- **Overpass API** - For OpenStreetMap data queries
- **OSRM** - For routing calculations

No API keys are required as these are open public APIs, but the server follows usage policies including proper user agent identification and request rate limiting.

## Development

### Project Structure

- `cmd/osmmcp` - Main application entry point
- `pkg/server` - MCP server implementation
- `pkg/tools` - OpenStreetMap tool implementations and tool registry
- `pkg/osm` - Common OpenStreetMap utilities and helpers

### Adding New Tools

To add a new tool:

1. Implement the tool functions in a new or existing file in `pkg/tools`
2. Add the tool definition to the registry in `pkg/tools/registry.go`

The registry-based design makes it easy to add new tools without modifying multiple files. All tool definitions are centralized in one place, making the codebase more maintainable.

### Troubleshooting

If you encounter build errors about redeclared types or functions, you might have older files from previous implementations. Check for and remove any conflicting files:

```bash
# Check for specialized.go which might conflict with newer implementations
rm -f pkg/tools/specialized.go pkg/tools/specialized_*.go

# Check for mock.go which might contain test implementations
rm -f pkg/tools/mock.go
```

## Acknowledgments

This implementation is based on two excellent sources:
- [jagan-shanmugam/open-streetmap-mcp](https://github.com/jagan-shanmugam/open-streetmap-mcp) - The original Python implementation
- [MCPLink OSM MCP Server](https://www.mcplink.ai/mcp/jagan-shanmugam/osm-mcp-server) - The MCPLink version with additional features

## License

MIT License 