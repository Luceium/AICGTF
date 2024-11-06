// Package generator provides functionality for generating code with AI.
package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

// SaveGeneratedCode saves the generated code to a file
func SaveGeneratedCode(code string, problem *Problem, model string) (string, error) {
	filename := fmt.Sprintf("%s_%s.go", strings.ToLower(model), problem.Title)

	//check that the parent path exists
	if _, err := os.Stat("out/generated"); os.IsNotExist(err) {
		if err := os.MkdirAll("out/generated", 0755); err != nil {
			return "", fmt.Errorf("error creating output directory: %v", err)
		}
	}

	filepath := filepath.Join("out/generated", filename)

	if err := os.WriteFile(filepath, []byte(code), 0644); err != nil {
		return "", fmt.Errorf("error saving generated code: %v", err)
	}

	return filepath, nil
}

// createPrompt creates a prompt for the AI model based on the problem
func createPrompt(problem *Problem) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Generate a Go solution for LeetCode problem: %s\n\n",
		problem.Title))
	sb.WriteString(fmt.Sprintf("Difficulty: %s\n\n", problem.Difficulty))
	sb.WriteString("Problem Statement:\n")
	sb.WriteString(problem.Statement)
	sb.WriteString("\n\nParameters:\n")

	for _, param := range problem.Parameters {
		sb.WriteString(fmt.Sprintf("- %s (%s)", param.Name, param.Type))
		if param.LowerBound != nil {
			sb.WriteString(fmt.Sprintf(", Lower bound: %v", param.LowerBound))
		}
		if param.UpperBound != nil {
			sb.WriteString(fmt.Sprintf(", Upper bound: %v", param.UpperBound))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// cleanGeneratedCode removes markdown code block syntax and trims whitespace
func cleanGeneratedCode(code string) string {
	code = strings.TrimPrefix(code, "```go")
	code = strings.TrimPrefix(code, "```")
	code = strings.TrimSuffix(code, "```")
	return strings.TrimSpace(code)
}
