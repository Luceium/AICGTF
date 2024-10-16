"""
Ai compatability layer
"""
SYSTEM_PROPMT = """
    SYSTEM: Only return the code, do not print anything else.
    Code must be in the form:
    def solution(x, y, z, ...):
        # code here
    Note: that there should be one function named solution and the parameters come from the problem description
"""

def prompt(prompt: str) -> str:
    return "def solution(x, y, z):\n\treturn 'hi'"
