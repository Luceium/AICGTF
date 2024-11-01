import re

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

class GenerationError(Exception):
    pass

async def prompt(prompt: str) -> str:
    return "def solution(x, y, z):\n\treturn 'hi'"

async def generateCode(prompt) -> bool:
    # uses ai and writes to file
    generatedCode = prompt(prompt)

    with open("tmp.py", "w") as f:
        f.write(generatedCode)
    
    return validateCode()

# check that ai code is in the correct format
def validateCode() -> bool:
    with open("tmp.py", "r") as f:
        code = f.read()
        return len(re.findall(r"^def solution\(", code)) == 1