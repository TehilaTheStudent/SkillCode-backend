package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestCase represents a test case for a coding question
type TestCase struct {
	Input          string `bson:"input" json:"input" validate:"required"`
	ExpectedOutput string `bson:"expected_output" json:"expected_output" validate:"required"`
}

// LanguageConfig holds the function signature for a specific language
type LanguageConfig struct {
	Language          string `bson:"language" json:"language" validate:"required"`                     // Programming language (e.g., python, golang)
	FunctionSignature string `bson:"function_signature" json:"function_signature" validate:"required"` // Function signature for the language
}

// Solution represents a user-provided solution for a coding question
type Solution struct {
	Language string `json:"language"`
	Function string `json:"user_function"`
}

// Question represents a coding question with test cases and language configurations
type Question struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string             `bson:"title" json:"title" validate:"required"`
	Description string             `bson:"description" json:"description" validate:"required"`
	TestCases   []TestCase         `bson:"test_cases" json:"test_cases" validate:"omitempty,dive"` // Optional but validates nested TestCases
	Languages   []LanguageConfig   `bson:"languages" json:"languages" validate:"required,dive"`    // Language-specific configurations
	Visibility  string             `bson:"visibility" json:"visibility" validate:"required,oneof=public private"`
	CreatedBy   string             `bson:"created_by" json:"created_by" validate:"required"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}