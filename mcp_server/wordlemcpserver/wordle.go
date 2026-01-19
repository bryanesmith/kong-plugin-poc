package wordlemcpserver

import (
	"context"

	"github.com/bryanesmith/wordle-help/go-sdk/recommendations"
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
	// Get recommendations from wordle-help library
	rated, err := recommendations.NextGuessRecommendations(input.Guesses, recommendations.DefaultDictionaryPath)
	if err != nil {
		return nil, GetWordleSuggestionsOutput{}, err
	}

	// Extract up to first 25 suggestions
	maxSuggestions := 25
	if len(rated) < maxSuggestions {
		maxSuggestions = len(rated)
	}

	suggestions := make([]string, maxSuggestions)
	for i := 0; i < maxSuggestions; i++ {
		suggestions[i] = rated[i].Guess
	}

	return nil, GetWordleSuggestionsOutput{
		Suggestions: suggestions,
	}, nil
}
