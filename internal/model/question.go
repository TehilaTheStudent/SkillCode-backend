package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type TestCase struct {
	Input          interface{} `bson:"input" json:"input" validate:"required"`
	ExpectedOutput interface{} `bson:"expected_output" json:"expected_output" validate:"required"`
}

type Question struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title             string             `bson:"title" json:"title" validate:"required"`
	Description       string             `bson:"description" json:"description" validate:"required"`
	FunctionSignature string             `bson:"function_signature" json:"function_signature" validate:"required"`
	TestCases         []TestCase         `bson:"test_cases" json:"test_cases" validate:"omitempty,dive"` // Optional but validates nested TestCases
	Visibility        string             `bson:"visibility" json:"visibility" validate:"required,oneof=public private"`
	CreatedBy         string             `bson:"created_by" json:"created_by" validate:"required"`
	CreatedAt         time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt         time.Time          `bson:"updated_at" json:"updated_at"`
}
