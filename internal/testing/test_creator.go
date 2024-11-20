package tester

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
)


func CreateTestRunner(language model.PredefinedSupportedLanguage, path string, question model.Question, userCode string) error {
	switch language {
	case model.Python:
		return CreatePythonTestRunner(path, question, userCode)
	case model.JavaScript:
		return CreateJavaScriptTestRunner(path, question, userCode)
	// Add other languages here (e.g., Java, Go, etc.)
	default:
		return fmt.Errorf("unsupported language: %s", language)
	}
}

func CreatePythonTestRunner(path string, question model.Question, userCode string) error {
	testCasesJSON, _ := json.Marshal(question.TestCases)
	content := fmt.Sprintf(`from evaluator import evaluate_user_code
import json
user_code = """%s"""
test_cases = %s
function_name = "%s"

results = evaluate_user_code(user_code, test_cases, function_name)
print(json.dumps(results, indent=2))`,
		userCode, string(testCasesJSON), question.FunctionConfig.Name,
	)

	return os.WriteFile(path, []byte(content), 0644)
}

func CreateJavaScriptTestRunner(path string, question model.Question, userCode string) error {
	// Convert test cases to JSON format
	testCasesJSON, err := json.Marshal(question.TestCases)
	if err != nil {
		return fmt.Errorf("failed to marshal test cases: %v", err)
	}

	// Escape backticks in userCode to prevent issues in JavaScript template literals
	escapedUserCode := strings.ReplaceAll(userCode, "`", "\\`")
	backtick:="`"
	// Generate the test runner content with proper escaping
	content := fmt.Sprintf(`const { evaluateUserCode } = require('./evaluator.js');
const userCode = %s%s%s;

const testCases = %s;
const functionName = "%s";

const results = evaluateUserCode(userCode, testCases, functionName);
console.log(JSON.stringify(results, null, 2));
`,backtick, escapedUserCode,backtick, string(testCasesJSON), question.FunctionConfig.Name)

	// Write the generated JavaScript test runner to the specified path
	return os.WriteFile(path, []byte(content), 0644)
}
