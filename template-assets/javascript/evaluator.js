const Ajv = require("ajv");
const fs = require("fs");
const path = require('path');
const converter = require('./converter');

function loadSchema() {
  try {
    const schemaPath = path.join(__dirname, "../feedback_schema.json");
    const schema = JSON.parse(fs.readFileSync(schemaPath, "utf8"));
    return schema;
  } catch (e) {
    throw new Error(`Error loading schema: ${e.message}`);
  }
}

function initializeAjv(schema) {
  try {
    const ajv = new Ajv();
    return ajv.compile(schema);
  } catch (e) {
    throw new Error(`Error initializing Ajv: ${e.message}`);
  }
}

function parseInput(inputString) {
  try {
    return JSON.parse(inputString); // Safely parse JSON input
  } catch (e) {
    return `Error parsing input: ${e.message}`;
  }
}

function validateResponse(response, validate) {
  try {
    const valid = validate(response);
    if (!valid) {
      // Prepare validation error details
      const errors = validate.errors.map(
        (err) => `${err.instancePath} ${err.message}`
      );
      return {
        status: "fail",
        results: [],
        error: "internal server error",
        details: `Schema validation failed: ${errors.join(", ")}`,
      };
    }
    return null; // Indicates valid response
  } catch (e) {
    return {
      status: "fail",
      results: [],
      error: "internal server error",
      details: `Error during validation: ${e.message}`,
    };
  }
}

function runTestCases(userFunction, testCases, validate,functionConfig) {
  const results = [];
  let allPassed = true;

  for (const testCase of testCases) {
    try {
      const inputs = testCase.parameters.map((param, index) => converter.listyToType(param, functionConfig.parameters[index].param_type));
      const expectedOutput = converter.listyToType(testCase.expected_output, functionConfig.return_type);

      const actualOutput = userFunction(...inputs);

      if (JSON.stringify(actualOutput) === JSON.stringify(expectedOutput)) {
        results.push({
          status: "pass",
          parameters: testCase.parameters,
          expected_output: String(expectedOutput),
          actual_output: String(actualOutput),
        });
      } else {
        allPassed = false;
        results.push({
          status: "fail",
          parameters: testCase.parameters,
          expected_output: String(expectedOutput),
          actual_output: String(actualOutput),
        });
      }
    } catch (e) {
      allPassed = false;
      results.push({
        status: "fail",
        parameters: testCase.parameters,
        expected_output: testCase.expected_output,
        actual_output: `Error: ${e.message}`,
      });
    }
  }

  const response = {
    status: allPassed ? "success" : "fail",
    results,
    error: allPassed ? null : "fail tests",
    details: allPassed ? null : "Some test cases failed.",
  };

  // Validate the response against the schema
  const validationError = validateResponse(response, validate);
  if (validationError) {
    return validationError;
  }

  return response;
}

function evaluateUserCode(userCode, testCases, functionName, functionConfig) {
  let userFunction;
  let validate;
  try {
    const schema = loadSchema();
    validate = initializeAjv(schema);
  } catch (e) {
    return {
      status: "fail",
      results: [],
      error: "internal server error",
      details: e.message,
    };
  }

  try {
    const wrappedCode = `
    const utils = require('./ds_utils.js');
    ${userCode}
    return ${functionName};
    `;

    userFunction = new Function("require", wrappedCode)(require);
  } catch (e) {
    const response = {
      status: "fail",
      results: [],
      error: "compilation",
      details: e.message,
    };

    // Validate the response against the schema
    const validationError = validateResponse(response, validate);
    if (validationError) {
      return validationError;
    }
    return response;
  }

  if (typeof userFunction !== "function") {
    const response = {
      status: "fail",
      results: [],
      error: "compilation",
      details: `${functionName} is not defined or not a function`,
    };

    const validationError = validateResponse(response, validate);
    if (validationError) {
      return validationError;
    }
    return response;
  }

  return runTestCases(userFunction, testCases, validate,functionConfig);
}

module.exports = { evaluateUserCode };
