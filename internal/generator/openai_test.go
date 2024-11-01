package generator

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

// At the top of the file, add this interface
type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// mockHTTPClient implements a mock http.Client for testing
type mockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestGetOpenAIGenerator(t *testing.T) {
	// Reset singleton for testing
	instance = nil

	// Get first instance
	gen1 := GetOpenAIGenerator("test-key", "gpt-4")
	if gen1 == nil {
		t.Fatal("Expected non-nil generator")
	}

	// Get second instance
	gen2 := GetOpenAIGenerator("different-key", "different-model")
	if gen2 != gen1 {
		t.Error("Expected same instance")
	}

	// Verify first configuration was kept
	if gen2.apiKey != "test-key" || gen2.model != "gpt-4" {
		t.Error("Expected original configuration to be maintained")
	}
}

func TestOpenAIGenerator_GenerateCode(t *testing.T) {
	// Create mock response
	mockResp := response{
		Choices: []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		}{
			{
				Message: struct {
					Content string `json:"content"`
				}{
					Content: "func twoSum(nums []int, target int) []int {\n\treturn nil\n}",
				},
			},
		},
	}

	// Create mock client
	mockClient := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// Verify request
			if req.Header.Get("Authorization") != "Bearer test-key" {
				t.Error("Missing or incorrect Authorization header")
			}

			// Return mock response
			respBody, _ := json.Marshal(mockResp)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBuffer(respBody)),
			}, nil
		},
	}

	// Create generator instance
	gen := &OpenAIGenerator{
		apiKey: "test-key",
		model:  "gpt-4",
		client: mockClient,
	}

	// Test problem
	problem := &Problem{
		Title:      "Two Sum",
		Difficulty: "Easy",
		Statement:  "Test problem statement",
		Parameters: []ProblemParameter{
			{Name: "nums", Type: "[]int"},
			{Name: "target", Type: "int"},
		},
	}

	// Generate code
	code, err := gen.GenerateCode(problem)
	if err != nil {
		t.Fatalf("GenerateCode() error = %v", err)
	}

	expectedCode := "func twoSum(nums []int, target int) []int {\n\treturn nil\n}"
	if code != expectedCode {
		t.Errorf("GenerateCode() = %v, want %v", code, expectedCode)
	}
}

func TestOpenAIGenerator_SaveGeneratedCode(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "generator_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to temp directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	// Create generator
	gen := &OpenAIGenerator{
		model: "gpt-4",
	}

	// Test saving code
	problem := &Problem{
		Title:      "Two Sum",
		Difficulty: "Easy",
		Statement:  "Given an array of integers nums and an integer target, return indices of the two numbers such that they add up to target.",
		Parameters: []ProblemParameter{
			{Name: "nums", Type: "[]int", LowerBound: 2, UpperBound: 10000},
			{Name: "target", Type: "int", LowerBound: -1000000000, UpperBound: 1000000000},
		},
	}
	testCode := "func twoSum(nums []int, target int) []int {\n\treturn nil\n}"

	if err := gen.SaveGeneratedCode(testCode, problem); err != nil {
		t.Fatalf("SaveGeneratedCode() error = %v", err)
	}

	// Verify file
	expectedPath := filepath.Join("generated", "gpt-4_test-problem.go")
	content, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	if string(content) != testCode {
		t.Errorf("File contents = %v, want %v", string(content), testCode)
	}
}

func TestOpenAIGenerator_CleanGeneratedCode(t *testing.T) {
	gen := &OpenAIGenerator{}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "clean code without markdown",
			input:    "func test() {}",
			expected: "func test() {}",
		},
		{
			name:     "code with markdown",
			input:    "```go\nfunc test() {}\n```",
			expected: "func test() {}",
		},
		{
			name:     "code with whitespace",
			input:    "\n\n  func test() {}\n  \n",
			expected: "func test() {}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := gen.cleanGeneratedCode(tt.input)
			if result != tt.expected {
				t.Errorf("cleanGeneratedCode() = %q, want %q", result, tt.expected)
			}
		})
	}
}
