package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"ACGTF/internal/evaluator"
	"ACGTF/internal/generator"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	// Get API key from environment variable
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	// Create generator config
	config := generator.GeneratorConfig{
		Provider:  "openai",
		Model:     "gpt-4o-mini",
		APIKey:    apiKey,
		MaxTokens: 2000,
	}

	// Initialize generator
	gen, err := generator.NewGenerator(config)
	if err != nil {
		log.Fatalf("Failed to create generator: %v", err)
	}

	// Create problem
	p := generator.Problem{
		Title:      "Two Sum",
		Difficulty: "Easy",
		Statement:  "Given an array of integers nums and an integer target, return indices of the two numbers such that they add up to target.",
		Parameters: []generator.ProblemParameter{
			{Name: "nums", Type: "[]int", LowerBound: 2, UpperBound: 10000},
			{Name: "target", Type: "int", LowerBound: -1000000000, UpperBound: 1000000000},
		},
	}

	// Create generated directory if it doesn't exist
	if err := os.MkdirAll("generated", 0755); err != nil {
		log.Fatalf("Failed to create generated directory: %v", err)
	}

	// Generate code
	code, err := gen.GenerateCode(&p)
	if err != nil {
		log.Fatalf("Failed to generate code: %v", err)
	}

	// Save the generated code
	filepath, err := generator.SaveGeneratedCode(code, &p, config.Model)
	if err != nil {
		log.Fatalf("Failed to save generated code: %v", err)
	}

	// Evaluate the code quality
	// fmt.Printf("Successfully generated and saved code for %s\n", p.Title)
	evaluationResult, err := evaluator.EvaluateCode(filepath)
	if err != nil {
		log.Fatalf("Failed to evaluate code: %v", err)
	}

	if data, err := json.Marshal(*evaluationResult); err == nil {
		fmt.Printf("Evaluation result: %v\n", string(data))
	}
}
