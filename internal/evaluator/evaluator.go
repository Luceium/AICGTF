package evaluator

type EvaluatorInterface interface {
	// Evaluates the code and returns an evaluation result
	EvaluateCode() (ResultInterface, error)
	// Calculates the score of the evaluation result
	// This sets the Score field of the evaluation result and returns it
	CalculateScore(extensionResults ResultInterface) int
}

type Evaluator struct {
	filepath string
	Name     string
}

type ResultInterface interface {
	Score() int
}

type ComprehensiveEvaluationResult struct {
	EvaluationResults []*ResultInterface
	Score             int
}

func EvaluateCode(filepath string) (*ComprehensiveEvaluationResult, error) {
	evaluators := [1]EvaluatorInterface{
		GetQualityEvaluator(filepath),
		// GetTestCaseEvaluator(filepath),
	}

	// build ComprehensiveEvaluationResult by looping through
	comprehensiveResult := &ComprehensiveEvaluationResult{[]*ResultInterface{}, 0}
	for _, evaluator := range evaluators {
		result, err := evaluator.EvaluateCode()
		if err != nil {
			return nil, err
		}
		comprehensiveResult.EvaluationResults = append(comprehensiveResult.EvaluationResults, &result)
	}

	// average the scores of the results
	score := 0
	for _, result := range comprehensiveResult.EvaluationResults {
		score += (*result).Score()
	}
	score /= len(comprehensiveResult.EvaluationResults)

	return &ComprehensiveEvaluationResult{
		EvaluationResults: comprehensiveResult.EvaluationResults,
		Score:             score,
	}, nil
}
