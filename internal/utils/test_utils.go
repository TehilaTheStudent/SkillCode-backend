package utils

import (
	"encoding/json"
	"time"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GenerateQuestion generates a `Question` struct with default values.
// You can override specific fields by passing a map with the key as the struct field name.
func GenerateQuestion(overrides map[string]interface{}) model.Question {
	// Default values for the Question struct
	question := model.Question{
		ID:          primitive.NewObjectID(),
		Title:       "Default Title",
		Description: "Default Description",
		TestCases: []model.TestCase{
			{Input: "[1, 2, 3, 4, 5], 3", ExpectedOutput: "2"},
			{Input: "[10, 20, 30, 40], 30", ExpectedOutput: "2"},
		},
		Languages: []model.LanguageConfig{
			{
				Language:          "golang",
				FunctionSignature: "func binarySearch(nums []int, target int) int",
			},
			{
				Language:          "python",
				FunctionSignature: "def binary_search(nums: List[int], target: int) -> int:",
			},
		},
		Visibility: "public",
		CreatedBy:  "test_user",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Apply overrides to customize specific fields
	for key, value := range overrides {
		switch key {
		case "ID":
			question.ID = value.(primitive.ObjectID)
		case "Title":
			question.Title = value.(string)
		case "Description":
			question.Description = value.(string)
		case "TestCases":
			question.TestCases = value.([]model.TestCase)
		case "Languages":
			question.Languages = value.([]model.LanguageConfig)
		case "Visibility":
			question.Visibility = value.(string)
		case "CreatedBy":
			question.CreatedBy = value.(string)
		case "CreatedAt":
			question.CreatedAt = value.(time.Time)
		case "UpdatedAt":
			question.UpdatedAt = value.(time.Time)
		}
	}

	return question
}

// GenerateCreateQuestionPayload generates a JSON payload for creating a question
func GenerateCreateQuestionPayload(overrides map[string]interface{}) string {
	payload := map[string]interface{}{
		"title":       "Default Title",
		"description": "Default Description",
		"test_cases": []map[string]string{
			{"input": "[1,2,3]", "expected_output": "6"},
		},
		"languages": []map[string]string{
			{"language": "golang", "function_signature": "func example() string"},
		},
		"visibility": "public",
		"created_by": "test_user",
	}

	// Apply overrides
	for key, value := range overrides {
		payload[key] = value
	}

	jsonData, _ := json.Marshal(payload)
	return string(jsonData)
}

// GenerateUserFunction generates a user function for the given language.
// Supports optional overrides for customization.
func GenerateUserFunction(language string, overrides map[string]interface{}) model.Solution {
	// Get the default function for the language
	solution := model.Solution{
		Function: generateUserFunction(language),
		Language: language,
	}

	// Apply overrides to customize specific fields
	for key, value := range overrides {
		switch key {
		case "Function":
			solution.Function = value.(string)
		case "Language":
			solution.Language = value.(string)
		}
	}

	return solution
}

// generateUserFunction generates a default user function for the specified language.
func generateUserFunction(language string) string {
	switch language {
	case "golang":
		return `
func binarySearch(nums []int, target int) int {
	low, high := 0, len(nums)-1
	for low <= high {
		mid := (low + high) / 2
		if nums[mid] == target {
			return mid
		} else if nums[mid] < target {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return -1
}`
	case "python":
		return `
def binary_search(nums, target):
	low, high = 0, len(nums) - 1
	while low <= high:
		mid = (low + high) // 2
		if nums[mid] == target:
			return mid
		elif nums[mid] < target:
			low = mid + 1
		else:
			high = mid - 1
	return -1`
	default:
		return "" // Return an empty string for unsupported languages
	}
}

// GenerateUserFunctionPayload generates a JSON payload for the user function, including language and code.
func GenerateUserFunctionPayload(language string) string {
	// Generate the function for the specified language
	function := generateUserFunction(language)

	// Build the payload
	payload := map[string]interface{}{
		"language":      language,
		"user_function": function,
	}

	// Convert payload to JSON
	jsonData, _ := json.Marshal(payload)
	return string(jsonData)
}
