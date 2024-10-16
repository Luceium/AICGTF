# Example Leet Code Question
### Description:
You are given two non-empty linked lists representing two non-negative integers. The digits are stored in reverse order, and each of their nodes contains a single digit. Add the two numbers and return the sum as a linked list.

You may assume the two numbers do not contain any leading zero, except the number 0 itself.

 

### Example 1:



Input: l1 = [2,4,3], l2 = [5,6,4]
Output: [7,0,8]
Explanation: 342 + 465 = 807.
Example 2:

Input: l1 = [0], l2 = [0]
Output: [0]
Example 3:

Input: l1 = [9,9,9,9,9,9,9], l2 = [9,9,9,9]
Output: [8,9,9,9,0,0,0,1]
 

### Constraints:

The number of nodes in each linked list is in the range [1, 100].
0 <= Node.val <= 9
It is guaranteed that the list represents a number that does not have leading zeros.



# Automated QA

import solution from tmpSolution

class Problem:
  statement = “”
  class Params:
    paramType
    lowerBound
    upperBound
  params = []

  
  class Test:
    args = {}
    expectedOutput

problems = [] # get from pre made DB (scraped from leetcode)


for problem in problems:
  code = ai.prompt(problem.statement)
  with open(“tmpSolution.py”, “w”) as runner:
    runner.write(code)
  test(solution)

def test(solution):
  # defined in following sections   




Compilation & Code quality - Does the code compile and is it idiomatic

def compilation(solution, args) -> bool:
  try:
    solution(**args)
    return True
  catch e:
    return False


	
def codeQuality(code) -> bool:
  prompt = f“””
  You are a code quality checker.
  Evaluate the program on best practices, readability, and idiomatic code.
  For example, Python code should use SnakeCase...
  ...
  Return a JSON response in the form {codeQuality: score}
  where score is a number between 0 and 1.

  Evaluate the following program:
  ${code}
  “””
  rsp = ai.prompt(prompt, temperature=0)
  return rsp.codeQuality



Objective & test cases - Does the code meet the requirements and pass test cases

def objective(code, statement) -> bool:
  prompt = f“””
  You are a code objective checker.
  Evaluate the program on how much it matches the problem objective
  ...
  Return a JSON response in the form {objectiveScore: score}
  where score is a number between 0 and 1.

  Given the problem statement:
  ${statement}

  Evaluate the following program:
  ${code}
  “””
  rsp = ai.prompt(prompt, temperature=0)
  return rsp.codeQuality


def testCases(solution) -> bool:
  map(lambda test ->
    try:
      rsp = solution(**args)
      return rsp == test.expectedOutput
    catch e:
      return False
  , tests)



Edge cases & Stress testing - Does it handle edge cases well and what are the maximum inputs it can handle within time constraints

Max input and min input per param in x 5 seconds (using binary search approach starting at top option.



Performance & Benchmarking - How fast can it complete the standard test case & how much memory does it take.
Compare oLama vs GPT-4o mini

Clock measure how much time it takes to run the first test


Check the max memory for the running program