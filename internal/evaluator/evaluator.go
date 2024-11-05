package evaluator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// EvaluationResult contains the results of code evaluation
type EvaluationResult struct {
	Compiles      bool
	CompileErrors string
	GoFmtClean    bool
	GoLintIssues  []string
	Score         int
}

var WORK_DIR = "out/generated"

// EvaluateCode performs compilation and quality checks on the generated code
func EvaluateCode(filename string) (*EvaluationResult, error) {
	result := &EvaluationResult{}

	// Get absolute path
	filepath := filepath.Join(WORK_DIR, filename)

	// Check if file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %s", filepath)
	}

	// Run go vet
	vetIssues := runGoVet(filepath)
	if len(vetIssues) > 0 {
		result.GoLintIssues = append(result.GoLintIssues, vetIssues...)
	}

	// Try to compile
	compileErr := compileCode(filepath)
	if compileErr != nil {
		result.Compiles = false
		result.CompileErrors = compileErr.Error()
	} else {
		result.Compiles = true
	}

	// Calculate score
	result.Score = calculateScore(result, filepath)

	return result, nil
}

func runGoVet(filepath string) []string {
	cmd := exec.Command("go", "vet", filepath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return []string{string(output)}
	}
	if len(output) > 0 {
		return strings.Split(string(output), "\n")
	}
	return nil
}

func compileCode(filepath string) error {
	cmd := exec.Command("go", "build", filepath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("compilation error: %s", string(output))
	}
	return nil
}

// calculateScore computes a quality score based on formatting and linting results
func calculateScore(result *EvaluationResult, filepath string) int {
	// Start with a perfect score of 100
	score := 100

	// Check if code compiles (critical requirement)
	if !result.Compiles {
		return 0
	}

	// Get the diff size by running gofmt
	diffSize := getDiffSize(filepath)
	// Subtract points based on diff size (1 point per 10 characters of diff)
	score -= diffSize / 10

	// Subtract points for linting issues (3 points per issue)
	score -= len(result.GoLintIssues) * 3

	// Ensure score doesn't go below 0
	if score < 0 {
		score = 0
	}

	return score
}

// getDiffSize calculates the size of the diff between current and formatted code
func getDiffSize(filepath string) int {
	// Run gofmt with -d flag to get diff
	cmd := exec.Command("gofmt", "-d", filepath)
	output, err := cmd.Output()
	if err != nil {
		return 500 // Return large diff size on error
	}

	// Parse the diff to get added and removed lines
	diffText := string(output)
	addedLines := make([]string, 0)
	removedLines := make([]string, 0)

	// Separate added and removed lines
	for _, line := range strings.Split(diffText, "\n") {
		if strings.HasPrefix(line, "+") {
			addedLines = append(addedLines, strings.TrimPrefix(line, "+"))
		} else if strings.HasPrefix(line, "-") {
			removedLines = append(removedLines, strings.TrimPrefix(line, "-"))
		}
	}

	// Calculate actual diff size by comparing line lengths
	diffSize := 0
	maxLen := len(addedLines)
	if len(removedLines) > maxLen {
		maxLen = len(removedLines)
	}

	// Compare corresponding lines and count actual differences
	for i := 0; i < maxLen; i++ {
		var added, removed string
		if i < len(addedLines) {
			added = addedLines[i]
		}
		if i < len(removedLines) {
			removed = removedLines[i]
		}
		// Only count the length of the line that changed more
		if len(added) > len(removed) {
			diffSize += len(added)
		} else {
			diffSize += len(removed)
		}
	}

	return diffSize
}
