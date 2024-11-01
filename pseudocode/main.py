import asyncio
import ai
import problem

def getProblems() -> List[problem.Problem]:
    problems = []
    # Problem 1
    problem1 = problem.Problem("Add two numbers")
    problem1.params.append(problem.Params("x", "int", 0, 100))
    problem1.params.append(problem.Params("y", "int", 0, 100))
    problem1.tests.append(problem.Test({"x": 1, "y": 2}, 3))
    problem1.tests.append(problem.Test({"x": 0, "y": 0}, 0))
    problems.append(problem1)
    # Problem 2
    problem2 = problem.Problem("Multiply two numbers")
    problem2.params.append(problem.Params("x", "int", 0, 100))
    problem2.params.append(problem.Params("y", "int", 0, 100))
    problem2.tests.append(problem.Test({"x": 1, "y": 2}, 2))
    problem2.tests.append(problem.Test({"x": 0, "y": 0}, 0))
    problems.append(problem2)
    return problems

def testProblem(problem: problem.Problem):
    # Generate code
    if ai.generateCode(problem.statement):
        from tmp import solution
    else:
        raise ai.GenerationError
    
    # Test: Compilation & Code Quality
    # Does the code compile?
    # Is the code idiomatic? And does it follow the style guidelines?

    # Test: Objective & Test Cases
    # Does meet the requirements?
    # Does the code pass all the test cases?

    # Test: Boundary Testing
    # Does the code handle edge cases?
    # What are the boundaries of the parameters within a fixed run time?

    # Test: Performance Testing & Benchmarking
    # How much time does it take for test cases to run?
    # What is the memory usage?

async def main():
    for problem in problem.getProblems():
        testProblem(problem)

if __name__ == "__main__":
    asyncio.run(main())