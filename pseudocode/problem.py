from typing import List

class Test:
    def __init__(self, args: dict, expectedOutput: any):
        self.args = args
        self.expectedOutput = expectedOutput

class Params:
    def __init__(self, paramName: str, paramType: str, lowerBound: int, upperBound: int):
        self.paramName = paramName
        self.paramType = paramType
        self.lowerBound = lowerBound
        self.upperBound = upperBound

class TestConfig:
    def __init__(self, timeout: int = 1000, memoryLimit: int = 1000, maxStyleDelta: int = 10):
        self.timeout = timeout
        self.memoryLimit = memoryLimit
        self.maxStyleDelta = maxStyleDelta

class Problem:
  def __init__(self, statement: str):
    self.statement = statement
    self.params : List[Params] = []
    self.testCases : List[Test] = []
    self.testConfig : TestConfig = TestConfig()
  