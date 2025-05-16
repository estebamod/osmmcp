# AI Prompts for Geocoding Tools

This document contains prompts to help AI assistants properly use the geocoding tools. These prompts can be integrated into your MCP configuration to improve geocoding success rates.

## System Prompts

### General Geocoding System Prompt

```
You have access to geocoding tools that convert between addresses and coordinates. When using these:

1. Format addresses clearly without parentheses, e.g., "Blue Temple Chiang Rai Thailand" instead of "Blue Temple (Wat Rong Suea Ten)"
2. Always include city and country for international locations 
3. If geocoding fails, check the error message for suggestions
4. Try progressive simplification when address lookups fail
5. For reverse geocoding, ensure coordinates are in decimal form within valid ranges
```

### Geocode Address Tool Prompt

```
When using the geocode_address tool, follow these guidelines:

1. SIMPLIFY COMPLEX QUERIES
   - Remove parentheses and special characters 
   - Example: "Blue Temple Chiang Rai Thailand" NOT "Blue Temple (Wat Rong Suea Ten)"

2. ADD GEOGRAPHIC CONTEXT
   - Always include city, region, country for international locations
   - Example: "Eiffel Tower, Paris, France" NOT just "Eiffel Tower"

3. ERROR HANDLING
   - If you receive a NO_RESULTS error, follow the suggestions in the response
   - Try removing parenthetical information or simplifying the query
   - Try adding geographic context (city/country names)

4. PROGRESSIVE REFINEMENT
   - Start specific, then try broader forms if needed
   - If detailed address fails, try just the landmark name with location
```

### Reverse Geocode Tool Prompt

```
When using the reverse_geocode tool, follow these guidelines:

1. FORMAT COORDINATES PROPERLY
   - Use decimal degrees (e.g., 37.7749, -122.4194)
   - Do NOT use degrees/minutes/seconds format

2. VALIDATE COORDINATE RANGES
   - Latitude must be between -90 and 90
   - Longitude must be between -180 and 180

3. PRECISION MATTERS
   - Use at least 4 decimal places when available
   - More precise coordinates yield better results

4. HANDLE EMPTY RESULTS
   - If no meaningful address is returned, try coordinates slightly offset from original
   - Shift by 0.0001 degrees in any direction and try again
```

## Example Integration with MCP-Go

Here's how to integrate these prompts into your MCP server configuration:

```go
package main

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/NERVsystems/osmmcp/pkg/tools"
)

func main() {
	// Create an MCP server
	s := server.NewMCPServer(
		"Geocoding Example",
		"1.0",
		server.WithToolCapabilities(true),
	)

	// Register geocoding tools with improved descriptions
	geocodeAddressTool := tools.GeocodeAddressTool()
	reverseGeocodeTool := tools.ReverseGeocodeTool()

	// Add tools to server
	s.AddTool(geocodeAddressTool, tools.HandleGeocodeAddress)
	s.AddTool(reverseGeocodeTool, tools.HandleReverseGeocode)

	// Create a prompt with instructions for the AI
	systemPrompt := `You have access to geocoding tools that convert between addresses and coordinates. 
When using these tools:

1. Format addresses clearly without parentheses, e.g., "Blue Temple Chiang Rai Thailand" instead of "Blue Temple (Wat Rong Suea Ten)"
2. Always include city and country for international locations 
3. If geocoding fails, check the error message for suggestions
4. Try progressive simplification when address lookups fail
5. For reverse geocoding, ensure coordinates are in decimal form within valid ranges`

	// Set up an example prompt template
	promptTemplate := mcp.NewPromptTemplate("geocoding",
		[]mcp.PromptMessage{
			mcp.NewPromptMessage(
				mcp.RoleSystem,
				mcp.NewTextContent(systemPrompt),
			),
		},
	)

	// Add the prompt template to the server
	s.AddPromptTemplate(promptTemplate)

	// Start the server
	// ...
}
```

## Example Client Interaction

When integrating with an LLM client, ensure the client is instructed to follow these patterns:

```json
{
  "prompt": {
    "messages": [
      {
        "role": "system",
        "content": {
          "type": "text",
          "text": "You have access to geocoding tools. Format addresses clearly without parentheses. Always include city and country for international locations."
        }
      },
      {
        "role": "user",
        "content": {
          "type": "text",
          "text": "Can you find the coordinates of the Blue Temple in Thailand?"
        }
      }
    ]
  }
}
```

The AI should then use the `geocode_address` tool with a properly formatted query:

```json
{
  "name": "geocode_address",
  "arguments": {
    "address": "Blue Temple Chiang Rai Thailand"
  }
}
```

Rather than:

```json
{
  "name": "geocode_address",
  "arguments": {
    "address": "Blue Temple (Wat Rong Suea Ten)"
  }
}
```

Which would likely fail to return results. 