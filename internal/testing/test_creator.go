package tester

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
)

func CreatePythonTestRunner(path string, question model.Question, userCode string) error {
	testCasesJSON, _ := json.Marshal(question.TestCases)
	content := fmt.Sprintf(
		`from evaluator import evaluate_user_code
import ds_utils as utils

user_code = """%s"""
test_cases = %s
function_name = "%s"

results = evaluate_user_code(user_code, test_cases, function_name)
print(results)`,
		userCode, string(testCasesJSON), question.FunctionConfig.Name,
	)

	return os.WriteFile(path, []byte(content), 0644)
}
