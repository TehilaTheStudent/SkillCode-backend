function parseInput(inputString) {
  try {
    return JSON.parse(inputString); // Safely parse JSON input
  } catch (e) {
    return `Error parsing input: ${e.message}`;
  }
}

function runTestCases(userFunction, testCases) {
  const results = [];

  for (const testCase of testCases) {
    try {
      // Parse each parameter string individually
      const inputs = testCase.parameters.map(parseInput);
      const expectedOutput = parseInput(testCase.expected_output);

      // Invoke the user's function
      const actualOutput = userFunction(...inputs);

      // Compare outputs
      if (JSON.stringify(actualOutput) === JSON.stringify(expectedOutput)) {
        results.push({
          status: "pass",
          expected_output: expectedOutput,
          actual_output: actualOutput,
        });
      } else {
        results.push({
          status: "fail",
          expected_output: expectedOutput,
          actual_output: actualOutput,
        });
      }
    } catch (e) {
      // Handle errors in execution
      results.push({
        status: "fail",
        expected_output: testCase.expected_output,
        actual_output: `Error: ${e.message}`,
      });
    }
  }

  return {
    status: "success",
    results,
    error: null,
    details: null,
  };
}

function evaluateUserCode(userCode, testCases, functionName) {
  let userFunction;
  try {
    const wrappedCode = `
    const utils = require('./ds_utils.js');
    ${userCode}
    return ${functionName};
    `;

    // Safely wrap and execute user code
    userFunction = new Function("require", wrappedCode)(require);
  } catch (e) {
    return {
      status: "fail",
      results: [],
      error: "Compilation failed",
      details: e.message,
    };
  }

  if (typeof userFunction !== "function") {
    return {
      status: "fail",
      results: [],
      error: `${functionName} is not defined or not a function`,
      details: null,
    };
  }

  // Run test cases
  return runTestCases(userFunction, testCases);
}

// Example usage
const userCode = `function binarySearch(arr, target) {
    let left = 0, right = arr.length - 1;
    w
        }
    }`;

const testCases = [
  { parameters: ["[1, 2, 3, 4, 5]", "3"], expected_output: "2" },
  { parameters: ["[1, 2, 3, 4, 5]", "6"], expected_output: "-1" },
];

const functionName = "binarySearch";

const results = evaluateUserCode(userCode, testCases, functionName);
// console.log(results);

module.exports = { evaluateUserCode };
