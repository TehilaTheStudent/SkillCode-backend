package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AbstractType represents a data type that can be atomic or composite.
type AbstractType struct {
	Type         string        `json:"type" bson:"type" validate:"required"`                   // AtomicType or CompositeType
	TypeChildren *AbstractType `json:"type_children,omitempty" bson:"type_children,omitempty"` // Recursive reference
}

// Parameter represents a function parameter.
type Parameter struct {
	Name      string       `json:"name" bson:"name" validate:"required"`             // Parameter name
	ParamType AbstractType `json:"param_type" bson:"param_type" validate:"required"` // Parameter type
}

type FunctionConfig struct {
	Name       string        `json:"name" bson:"name" validate:"required"`               // Function name
	Parameters *[]Parameter  `json:"parameters,omitempty" bson:"parameters,omitempty"`   // Pointer to slice for nil vs empty distinction
	ReturnType *AbstractType `json:"return_type,omitempty" bson:"return_type,omitempty"` // Nil means VoidType
}

// InputOutput represents example inputs and expected outputs for a function
type InputOutput struct {
	Parameters     []string `bson:"parameters" json:"parameters" validate:"required"`           // Input parameters
	ExpectedOutput string   `bson:"expected_output" json:"expected_output" validate:"required"` // Expected output
}

// Question represents a coding question with metadata and configurations
type Question struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`                                       // MongoDB ID
	Title          string             `bson:"title" json:"title" validate:"required"`                                  // Question title
	Description    string             `bson:"description" json:"description" validate:"required"`                      // Question description
	Difficulty     string             `bson:"difficulty" json:"difficulty" validate:"required,oneof=Easy Medium Hard"` // Difficulty level
	Category       string             `bson:"category" json:"category" validate:"required"`                            // Question category (e.g., Tree, Array)
	Stats          int                `bson:"stats" json:"stats"`                                                      // Submission stats
	Examples       []InputOutput      `bson:"examples" json:"examples" validate:"dive"`                                // Examples of input/output
	TestCases      []InputOutput      `bson:"test_cases" json:"test_cases" validate:"dive"`                            // Test cases
	FunctionConfig FunctionConfig     `bson:"function_config" json:"function_config"`                                  // Function signature configuration
	Languages      []string           `bson:"languages" json:"languages" validate:"dive"`                              // Supported programming languages
}

// Solution represents a user-provided solution for a coding question
type Submission struct {
	Language PredefinedSupportedLanguage `json:"language"`
	Code     string                      `json:"code"`
}

type QuestionQueryParams struct {
	Search       string   `json:"search"`
	Categories   []string `json:"categories"`
	Difficulties []string `json:"difficulties"`
	SortBy       string   `json:"sort_by"`
	SortOrder    string   `json:"order"`
}

