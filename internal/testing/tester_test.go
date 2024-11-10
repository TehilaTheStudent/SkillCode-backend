package tester_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/testing"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestTester(t *testing.T) {
	// Step 1: Create the Question
	question := model.Question{
		ID:          primitive.NewObjectID(),
		Title:       "Merge Two Sorted Arrays",
		Description: "Merge two sorted arrays in-place.",
		TestCases: []model.TestCase{
			{Input: "[[1,2,3,0,0,0],3,[2,5,6],3]", ExpectedOutput: "[1,2,2,3,5,6]"},
			{Input: "[[0],0,[1],1]", ExpectedOutput: "[1]"},
		},
		FunctionSignature: "func merge(nums1 []int, m int, nums2 []int, n int)",
	}

	// Step 2: Define User Function
	userFunction := `func merge(nums1 []int, m int, nums2 []int, n int) {
		copy(nums1[m:], nums2[:n])
		copy(nums1[:], nums1[:m+n])
	}`

	// Step 3: Specify Language
	language := "golang"

	// Step 4: Call TestUserSolution
	results, err := tester.TestUserSolution(&question, userFunction, language)
	if err != nil {
		t.Fatalf("TestUserSolution failed: %v", err)
	}

	t.Logf("Test results:\n%s", results)
}

func TestMain(m *testing.M) {
	// Change working directory to the project root
	if err := os.Chdir("../.."); err != nil {
		fmt.Printf("Failed to change working directory: %v\n", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}