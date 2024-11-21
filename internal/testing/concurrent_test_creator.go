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

	// Step 2: Generate test runner file
	testRunnerPath := filepath.Join(uniqueDir, fmt.Sprintf("main.%s", model.GetFileExtension(submission.Language)))
	err := CreateTestRunner(submission.Language, testRunnerPath, question, submission.Code)
	if err != nil {
		// Cleanup the directory if generation fails
		_ = os.RemoveAll(uniqueDir)
		return "", fmt.Errorf("failed to generate test runner: %v", err)
	}

	return uniqueDir, nil
}


func (t *UniqueTester) ExecuteUniqueTest(uniqueDir string,string code) (string, error) {
	params := map[string]string{
		"JOB_NAME":       t.jobName,
		"IMAGE_NAME":     t.imageName,
		"RUNTIME_COMMAND": t.runtimeCommand,
		"FILE_EXTENSION": t.fileExtension,
		"REQUEST_ID":     t.requestID,
	}

	// Use job template and uniqueDir
	return t.ExecuteWithJobTemplate(params, filepath.Join(config.GlobalConfigAPI.JobTemplatePath), )
}


func CleanupUniqueAssets(path string) error {
	return os.RemoveAll(path)
}
