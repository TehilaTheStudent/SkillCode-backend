package tester

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/config"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
)

func GenerateUniqueAssets(requestID string, question model.Question, submission model.Submission) (string, error) {
	// Step 1: Create unique directory
	uniqueDir := filepath.Join(config.GlobalConfigAPI.UniqueAssetsDir, requestID)
	if err := os.MkdirAll(uniqueDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create unique directory: %v", err)
	}
	var testRunnerPath string
	if config.GlobalConfigAPI.ModeEnv != "production" {
		testRunnerPath = filepath.Join(config.GlobalLanguageConfigs[submission.Language].AssetsDir, fmt.Sprintf("dev.%s", model.GetFileExtension(submission.Language)))
	} else {
		testRunnerPath = filepath.Join(uniqueDir, fmt.Sprintf("main.%s", model.GetFileExtension(submission.Language)))
	}
	// Step 2: Generate test runner file
	err := CreateTestRunner(submission.Language, testRunnerPath, question, submission.Code)
	if err != nil {
		// Cleanup the directory if generation fails
		_ = os.RemoveAll(uniqueDir)
		return "", fmt.Errorf("failed to generate test runner: %v", err)
	}

	return uniqueDir, nil
}

func (t *UniqueTester) ExecuteUniqueTest(uniqueDir string, scriptContent string) (string, error) {
	params := map[string]string{
		"JOB_NAME":        t.jobName,
		"IMAGE_NAME":      t.imageName,
		"RUNTIME_COMMAND": t.runtimeCommand,
		"FILE_EXTENSION":  t.fileExtension,
		"REQUEST_ID":      t.requestID,
	}

	// Construct the path to `main.<fileExtension>`
	filePath := filepath.Join(uniqueDir, "main."+t.fileExtension)

	// Read the contents of `main.<fileExtension>`
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %v", filePath, err)
	}

	// Convert fileContent to string
	fileContentStr := string(fileContent)

	// Use fileContentStr and params
	return t.ExecuteWithJobTemplate(params, config.GlobalConfigAPI.JobTemplatePath, fileContentStr)
}

func CleanupUniqueAssets(path string) error {
	return os.RemoveAll(path)
}
