package tester_test

import (
	"os"
	"testing"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/testing"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/utils"
)

func TestTester(t *testing.T) {

	// Step 1: Define the Question with Function Signatures for Multiple Languages
	question := utils.GenerateQuestion(nil)

	// Step 2: Define the User Function (Body Only)
	// For Python, the user provides only the function body (no class or method signature)
	userFunctionPython := utils.GenerateUserFunction("python", nil)

	// Step 5: Call TestUserSolution (Generate and Run Test Harness)
	results, err := tester.TestUserSolution(&question, userFunctionPython.Function, userFunctionPython.Language)
	if err != nil {
		t.Fatalf("TestUserSolution failed: %v", err)
	}

	// Step 6: Log the Results
	// fmt.Println("Test results:\n" + results)
	t.Logf("Test results:\n%s", results)
}

func TestMain(m *testing.M) {
	utils.EnsureWorkingDirectory()
	os.Exit(m.Run())
}
