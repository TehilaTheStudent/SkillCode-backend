package parser_validator_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/parser_validator"
)

// GenerateAbstractType generates a random AbstractType.
func GenerateAbstractType(depth int) *model.AbstractType {
	if depth == 0 {
		selectedType := model.AtomicTypes[rand.Intn(len(model.AtomicTypes))]
		abstractType := &model.AbstractType{Type: string(selectedType), TypeChildren: nil}
		return abstractType
	}

	selectedType := model.CompositeTypes[rand.Intn(len(model.CompositeTypes))]

	childType := GenerateAbstractType(depth - 1)
	return &model.AbstractType{
		Type:         string(selectedType),
		TypeChildren: childType,
	}
}

//i dont use it
func ReformatStringOfType(input string) (string, error) {
	var parsed interface{}

	// Parse the input string as JSON
	if err := json.Unmarshal([]byte(input), &parsed); err != nil {
		return "", fmt.Errorf("failed to parse input: %w", err)
	}

	// Reformat into a clean JSON string
	formatted, err := json.Marshal(parsed)
	if err != nil {
		return "", fmt.Errorf("failed to format input: %w", err)
	}

	// Iterate over the formatted string and add a space after each comma
	result := ""
	for i := 0; i < len(formatted); i++ {
		result += string(formatted[i])
		if formatted[i] == ',' {
			result += " "
		}
	}

	return result, nil
}

func TestAbstractTypeValidation(t *testing.T) {
	for i := 0; i < 1; i++ {
		abstractType := GenerateAbstractType(2)
		fmt.Println(abstractType.ToPrint())
		validInput := parser_validator.GenerateValidString(abstractType)
		fmt.Printf("Valid input: %s\n", validInput)
		err := parser_validator.ValidateAbstractType(validInput, abstractType)
		if err != nil {
			t.Errorf("Valid input failed validation:\nAbstractType: %+v\nInput: %s\nError: %v", abstractType, validInput, err)
		}
	}
}
