package utils

import (
	"encoding/json"

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
		Difficulty:  string(model.Easy),
		Category:    string(model.Array),
		Stats:       0,
		Examples: []model.InputOutput{
			{Parameters: []string{"[1, 2, 3, 4, 5]", "3"}, ExpectedOutput: "2"},
		},
		TestCases: []model.InputOutput{
			{Parameters: []string{"[1, 2, 3, 4, 5]", "3"}, ExpectedOutput: "2"},
			{Parameters: []string{"[10, 20, 30, 40], 30"}, ExpectedOutput: "2"},
		},
		FunctionConfig: model.FunctionConfig{
			Name: "binarySearch",
			Parameters: &[]model.Parameter{{Name: "nums",
				ParamType: model.AbstractType{Type: string(model.Array),
					TypeChildren: &model.AbstractType{Type: string(model.Integer)}}},
				{Name: "target",
					ParamType: model.AbstractType{Type: string(model.Integer)}}},
			ReturnType: &model.AbstractType{Type: string(model.Integer)},
		},
		Languages: []string{string(model.Python), string(model.Java)},
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
			question.TestCases = value.([]model.InputOutput)
		case "Languages":
			question.Languages = value.([]string)
		case "Difficulty":
			question.Difficulty = value.(string)
		case "Category":
			question.Category = value.(string)
		case "Stats":
			question.Stats = value.(int)
		case "Examples":
			question.Examples = value.([]model.InputOutput)
		case "FunctionConfig":
			question.FunctionConfig = value.(model.FunctionConfig)
		}
	}

	return question
}

// GenerateCreateQuestionPayload generates a JSON payload for creating a question
func GenerateCreateQuestionPayload(overrides map[string]interface{}) string {
	payload := map[string]interface{}{
		"title":       "Default Title",
		"description": "Default Description",
		"difficulty":  "Easy",
		"category":    "Array",
		"stats":       0,
		"examples": []map[string]interface{}{
			{"parameters": []string{"[1, 2, 3, 4, 5]", "3"}, "expected_output": "2"},
		},
		"test_cases": []map[string]string{
			{"input": "[1,2,3]", "expected_output": "6"},
		},
		"function_config": map[string]interface{}{
			"name":        "binarySearch",
			"parameters":  []map[string]string{{"name": "nums", "param_type": "[]int"}, {"name": "target", "param_type": "int"}},
			"return_type": "int",
		},
		"languages":  []string{"golang", "python"},
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
func GenerateUserFunction(language string, overrides map[string]interface{}) model.Submission {
	// Get the default function for the language
	solution := model.Submission{
		Code:     generateUserFunction(language),
		Language: model.Python,
	}

	// Apply overrides to customize specific fields
	for key, value := range overrides {
		switch key {
		case "Function":
			solution.Code = value.(string)
		case "Language":
			solution.Language = value.(model.PredefinedSupportedLanguage)
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
