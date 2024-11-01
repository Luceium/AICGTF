// Package generator provides functionality for generating code with AI.
package generator

import (
	"fmt"
)

// Generator defines the interface for code generation
type Generator interface {
	GenerateCode(problem *Problem) (string, error)
	SaveGeneratedCode(code string, problem *Problem) error
}

// Problem represents the input for code generation
type Problem struct {
	Title      string             `json:"title"`
	Difficulty string             `json:"difficulty"`
	Statement  string             `json:"statement"`
	Parameters []ProblemParameter `json:"parameters"`
}

// ProblemParameter represents an input parameter for the problem
type ProblemParameter struct {
	Name       string      `json:"name"`
	Type       string      `json:"type"`
	LowerBound interface{} `json:"lowerBound"`
	UpperBound interface{} `json:"upperBound"`
}

// GeneratorConfig holds configuration for code generation
type GeneratorConfig struct {
	Provider  string            // The provider (e.g., "openai", "anthropic")
	Model     string            // The model to use (e.g., "gpt-4", "claude-2")
	APIKey    string            // API key for the service
	MaxTokens int               // Maximum tokens for generation
	Options   map[string]string // Additional provider-specific options
}

// NewGenerator creates a new code generator based on the provider
func NewGenerator(config GeneratorConfig) (Generator, error) {
	switch config.Provider {
	case "openai":
		return GetOpenAIGenerator(config.APIKey, config.Model), nil
	// Add other providers here as needed
	default:
		return nil, fmt.Errorf("unsupported provider: %s", config.Provider)
	}
}
