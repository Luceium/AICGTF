import re
import asyncio
import ai

async def generateCode(prompt):
    # uses ai and writes to file
    generatedCode = ai.prompt(prompt)

    with open("tmp.py", "w") as f:
        f.write(generatedCode)
    return ""

# check that ai code is in the correct format
def validateCode() -> bool:
    with open("tmp.py", "r") as f:
        code = f.read()
        return len(re.findall(r"^def solution\(", code)) == 1

async def main():
    await generateCode("tmp")
    if validateCode():
        from tmp import solution
    else:
        raise ValueError


if __name__ == "__main__":
    asyncio.run(main())