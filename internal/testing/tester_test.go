package tester_test

import (
	"os"
	"testing"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/config"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/dependencies"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/testing"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/utils"
)

func TestTester(t *testing.T) {
	// Generate a question with function signatures for multiple languages
	question := utils.GenerateQuestion(nil)

	// Generate the user function body for Python
	userFunctionPython := utils.GenerateUserFunction("python", nil)
	sandboxConfig := config.GlobalConfigSandboxes[model.CSharp]
	
	// Test the user solution by generating and running the test harness
	results, err := tester.TestUserSolution(&question, userFunctionPython.Code,userFunctionPython.Language, *sandboxConfig)
	if err != nil {
		t.Fatalf("TestUserSolution failed: %v", err)
	}

	// Log the test results
	t.Logf("Test results:\n%s", results)
}

func TestMain(m *testing.M) {
	// Ensure the working directory is set correctly
	dependencies.EnsureWorkingDirectory()
	os.Exit(m.Run())
}
