package wordlemcpserver

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GetWordleSuggestionsInput struct {
	Guesses []string `json:"guesses" jsonschema:"array of previous guesses"`
}

type GetWordleSuggestionsOutput struct {
	Suggestions []string `json:"suggestions" jsonschema:"suggested words to try next"`
}

func GetWordleSuggestions(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetWordleSuggestionsInput,
) (*mcp.CallToolResult, GetWordleSuggestionsOutput, error) {
	return nil, GetWordleSuggestionsOutput{
		Suggestions: []string{"apple"},
	}, nil
}
