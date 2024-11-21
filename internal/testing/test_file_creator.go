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
)

// CreateTestRunner generates a test runner script using templates for the specified language
func CreateTestRunner(language model.PredefinedSupportedLanguage, path string, question model.Question, userCode string) error {
	// Map language to its template file path
	templatePath := filepath.Join(config.GlobalLanguageConfigs[language].AssetsDir, "main.tmpl")

	// Load test cases as JSON
	testCasesJSON, err := json.Marshal(question.TestCases)
	if err != nil {
		return fmt.Errorf("failed to marshal test cases: %v", err)
	}

	// Prepare template data
	var functionName string
	if language == model.Python {
		functionName = coding.ToPythonStyle(question.FunctionConfig.Name)
	} else if language == model.JavaScript {
		functionName = coding.ToJSStyle(question.FunctionConfig.Name)
	} else {
		return fmt.Errorf("unsupported language: %v", language)
	}

	data := map[string]string{
		"UserCode":     userCode,
		"TestCases":    string(testCasesJSON),
		"FunctionName": functionName,
	}

	// Generate the test runner
	return generateFromTemplate(templatePath, path, data)
}

// generateFromTemplate processes the template with given data and writes to the output path
func generateFromTemplate(templatePath, outputPath string, data map[string]string) error {
	// Parse the template
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed to parse template file %s: %v", templatePath, err)
	}

	// Execute the template with the provided data
	var output bytes.Buffer
	if err := tmpl.Execute(&output, data); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	// Write the generated content to the output file
	if err := os.WriteFile(outputPath, output.Bytes(), 0644); err != nil {
		// Ensure the directory exists
		if err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directories for output file: %v", err)
		}
		return fmt.Errorf("failed to write output file: %v", err)
	}

	return nil
}
