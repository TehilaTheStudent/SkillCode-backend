package tester

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/coding"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/config"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/utils"
)


func (t *UniqueTester) ExecuteUniqueTestProducton(scriptContent string) (string, error) {
	params := map[string]string{
		"JOB_NAME":        t.jobName,
		"IMAGE_NAME":      t.imageName,
		"RUNTIME_COMMAND": t.runtimeCommand,
		"FILE_EXTENSION":  t.fileExtension,
		"REQUEST_ID":      t.requestID,
	}

	return t.ExecuteWithJobTemplate(params, config.GlobalConfigAPI.JobTemplatePath, scriptContent)
}

func (t *UniqueTester) ExecuteUniqueTestDevelopment(scriptContent string) (string, error) {
	uniqueTestRunnerPath := filepath.Join(config.GlobalLanguageConfigs[t.language].AssetsDir, fmt.Sprintf("%s.%s",t.requestID, model.GetFileExtension(t.language)))
	
	// Write the script content to the file
	err := os.WriteFile(uniqueTestRunnerPath, []byte(scriptContent), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write test runner to file: %v", err)
	}
	// defer os.RemoveAll(uniqueTestRunnerPath)
	// Execute the file using the runtime command
	rawLogs, err := utils.RunCommand(t.runtimeCommand,uniqueTestRunnerPath)
	if err != nil {
		return "", fmt.Errorf("failed to execute file %s: %v", uniqueTestRunnerPath, err)
	}

	// return raw logs
	return rawLogs, nil
}



// CreateTestRunner generates a test runner script using templates for the specified language
func CreateTestRunnerScript(language model.PredefinedSupportedLanguage,  question model.Question, userCode string) (string, error) {
	// Map language to its template file path
	templatePath := filepath.Join(config.GlobalLanguageConfigs[language].AssetsDir, "main.tmpl")

	// Load test cases as JSON
	testCasesJSON, err := json.Marshal(question.TestCases)
	if err != nil {
		return "", fmt.Errorf("failed to marshal test cases: %v", err)
	}

	// Prepare template data
	var functionName string
	if language == model.Python {
		functionName = coding.ToPythonStyle(question.FunctionConfig.Name)
	} else if language == model.JavaScript {
		functionName = coding.ToJSStyle(question.FunctionConfig.Name)
	} else {
		return "", fmt.Errorf("unsupported language: %v", language)
	}

	data := map[string]string{
		"UserCode":     userCode,
		"TestCases":    string(testCasesJSON),
		"FunctionName": functionName,
	}

	// Generate the test runner
	return generateFromTemplate(templatePath, data)
}

// generateFromTemplate processes the template with given data and returns the generated content as a string
func generateFromTemplate(templatePath string, data map[string]string) (string, error) {
	// Parse the template
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to parse template file %s: %v", templatePath, err)
	}

	// Execute the template with the provided data
	var output bytes.Buffer
	if err := tmpl.Execute(&output, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %v", err)
	}

	return output.String(), nil
}
