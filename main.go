package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"your-module/leetcode"
)

func main() {
	// Create problems directory if it doesn't exist
	if err := os.MkdirAll("problems", 0755); err != nil {
		log.Fatalf("Error creating problems directory: %v", err)
	}

	// Example problems to fetch
	problems := []string{
		"two-sum",
		"add-two-numbers",
		"longest-substring-without-repeating-characters",
	}

	for _, titleSlug := range problems {
		problem, err := leetcode.FetchProblem(titleSlug)
		if err != nil {
			log.Printf("Error fetching problem %s: %v", titleSlug, err)
			continue
		}

		// Save the problem to the problems directory
		problemPath := filepath.Join("problems", fmt.Sprintf("%s_%s.json", problem.ID, problem.TitleSlug))
		if err := leetcode.SaveProblem(problem, problemPath); err != nil {
			log.Printf("Error saving problem %s: %v", titleSlug, err)
			continue
		}

		fmt.Printf("Successfully fetched and saved problem: %s\n", problem.Title)
	}
} 