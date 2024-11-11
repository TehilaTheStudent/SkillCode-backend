package tester

import (
	"fmt"
	"os"
	"strings"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
)

// generatePythonScript dynamically generates a Python script for testing
func generatePythonScript(filePath string, userFunction string, functionSignature string, testCases []model.TestCase) error {
	// Start building the Python script with the function signature
	script := fmt.Sprintf(`%s`, functionSignature)

	// Add user function body (properly indented)
	for _, line := range splitLines(userFunction) {
		script += fmt.Sprintf("        %s\n", line) // Proper indentation for Python code
	}

	// Add test cases
	script += "test_cases = [\n"
	for _, testCase := range testCases {
		script += fmt.Sprintf("    %s,\n", formatTestCase(testCase))
	}
	script += "]\n\n"

	// Add the test harness logic inside a main function
	script += `
# Initialize the Solution class
solution = Solution()

def main():
    # Results container
    results = []

    # Execute all test cases
    for i, test_case in enumerate(test_cases):
        try:
            # Extract inputs dynamically
            inputs = test_case['input']
            expected_output = test_case['expected_output']

            # Dynamically call the user's function
            actual_output = solution.merge(**inputs)  # Adjust function name if needed

            # Compare actual and expected outputs
            if actual_output == expected_output:
                results.append({
                    'status': 'passed',
                    'input': inputs,
                    'expected': expected_output,
                    'actual': actual_output
                })
            else:
                results.append({
                    'status': 'failed',
                    'input': inputs,
                    'expected': expected_output,
                    'actual': actual_output
                })
        except Exception as e:
            # Handle runtime errors gracefully
            results.append({
                'status': 'runtime_error',
                'input': inputs,
                'expected': expected_output,
                'error': str(e)
            })

    # Print results in JSON-like format
    import json
    print(json.dumps(results, indent=4))

if __name__ == "__main__":
    main()
`

	// Write the script to the file
	err := os.WriteFile(filePath, []byte(script), 0644)
	if err != nil {
		return fmt.Errorf("failed to write script to file: %w", err)
	}

	return nil
}

// Helper to split lines of code
func splitLines(code string) []string {
	return strings.Split(code, "\n")
}

// Helper to format test cases as Python dictionaries
func formatTestCase(tc model.TestCase) string {
	return fmt.Sprintf(`{"input": %v, "expected_output": %v}`, tc.Input, tc.ExpectedOutput)
}


