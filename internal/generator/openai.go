package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const openaiEndpoint = "https://api.openai.com/v1/chat/completions"

// OpenAIGenerator implements the Generator interface for OpenAI
type OpenAIGenerator struct {
	apiKey string
	model  string
	client *http.Client
}

// singleton instance
var (
	instance *OpenAIGenerator
	once     sync.Once
)

// GetOpenAIGenerator returns the singleton instance of OpenAIGenerator
func GetOpenAIGenerator(apiKey, model string) *OpenAIGenerator {
	once.Do(func() {
		instance = &OpenAIGenerator{
			apiKey: apiKey,
			model:  model,
			client: &http.Client{},
		}
	})
	return instance
}

// request represents a simplified chat completion request
type request struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
}

// response represents a simplified chat completion response
type response struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// GenerateCode generates Go code for the given problem using OpenAI
func (g *OpenAIGenerator) GenerateCode(problem *Problem) (string, error) {
	prompt := g.createPrompt(problem)

	// Create request body
	reqBody := request{
		Model: g.model,
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	// Marshal request to JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", openaiEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.apiKey)

	fmt.Printf("Request: %s\n", string(jsonData))

	// Make request
	resp, err := g.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var respBody response
	if err := json.Unmarshal(body, &respBody); err != nil {
		return "", fmt.Errorf("error parsing response: %v", err)
	}

	if len(respBody.Choices) == 0 {
		return "", fmt.Errorf("no response from API")
	}

	return g.cleanGeneratedCode(respBody.Choices[0].Message.Content), nil
}

// SaveGeneratedCode saves the generated code to a file
func (g *OpenAIGenerator) SaveGeneratedCode(code string, problem *Problem) error {
	filename := fmt.Sprintf("%s_%s.go", strings.ToLower(g.model), problem.Title)
	filepath := filepath.Join("generated", filename)

	if err := os.WriteFile(filepath, []byte(code), 0644); err != nil {
		return fmt.Errorf("error saving generated code: %v", err)
	}

	return nil
}

func (g *OpenAIGenerator) createPrompt(problem *Problem) string {
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

func (g *OpenAIGenerator) cleanGeneratedCode(code string) string {
	code = strings.TrimPrefix(code, "```go")
	code = strings.TrimPrefix(code, "```")
	code = strings.TrimSuffix(code, "```")
	return strings.TrimSpace(code)
}
