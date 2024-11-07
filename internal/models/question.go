package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type TestCase struct {
	Input          interface{} `bson:"input" json:"input"`
	ExpectedOutput interface{} `bson:"expected_output" json:"expected_output"`
}

type Question struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title             string             `bson:"title" json:"title"`
	Description       string             `bson:"description" json:"description"`
	FunctionSignature string             `bson:"function_signature" json:"function_signature"`
	TestCases         []TestCase         `bson:"test_cases" json:"test_cases"`
	Visibility        string             `bson:"visibility" json:"visibility"`
	CreatedBy         string             `bson:"created_by" json:"created_by"`
	CreatedAt         time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt         time.Time          `bson:"updated_at" json:"updated_at"`
}
