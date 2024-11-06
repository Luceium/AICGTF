package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

const openaiEndpoint = "https://api.openai.com/v1/chat/completions"

// OpenAIGenerator implements the Generator interface for OpenAI
type OpenAIGenerator struct {
	apiKey string
	model  string
	client *http.Client
	Generator
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
	prompt := createPrompt(problem)

	// Create request body
	reqBody := request{
		Model: g.model,
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{
				Role:    "system",
				Content: "You are a Go code generator for LeetCode problems. Act as software and only output the code. You will loose points for including any other text.",
			},
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

	// fmt.Printf("Request: %s\n", string(jsonData))

	// // Make request
	// resp, err := g.client.Do(req)
	// if err != nil {
	// 	return "", fmt.Errorf("error making request: %v", err)
	// }
	// defer resp.Body.Close()

	// // Read response body
	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return "", fmt.Errorf("error reading response: %v", err)
	// }

	// // Check status code
	// if resp.StatusCode != http.StatusOK {
	// 	return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	// }

	// // Parse response
	// var respBody response
	// if err := json.Unmarshal(body, &respBody); err != nil {
	// 	return "", fmt.Errorf("error parsing response: %v", err)
	// }

	// if len(respBody.Choices) == 0 {
	// 	return "", fmt.Errorf("no response from API")
	// }

	// var code string = respBody.Choices[0].Message.Content

	// TODO: Revert temporary change to use pregenerated code to save on LLM costs
	code := `package main

import "fmt"

func twoSum(nums []int, target int) []int {
    numMap := make(map[int]int)
    for i, num := range nums {
        complement := target - num
        if idx, found := numMap[complement]; found {
            return []int{idx, i}
        }
        numMap[num] = i
    }
    return nil
}

func main() {
    nums := []int{2, 7, 11, 15}
    target := 9
    result := twoSum(nums, target)
    fmt.Println(result)
}`

	return cleanGeneratedCode(code), nil
}
