package parser_validator

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
)

// Validators for atomic types
var (
	integerRegex = regexp.MustCompile(`^-?\d+$`)
	doubleRegex  = regexp.MustCompile(`^-?\d+(\.\d+)?$`)
	stringRegex  = regexp.MustCompile(`^[^\[\],]+$`)
)

// ValidateAtomicType validates an atomic type
func ValidateAtomicType(input string, atomicType string) error {
	switch model.AtomicType(atomicType) {
	case model.Integer:
		if !integerRegex.MatchString(input) {
			return fmt.Errorf("invalid Integer: %s not all digits", input)
		}
	case model.Double:
		if !doubleRegex.MatchString(input) {
			return fmt.Errorf("invalid Double: %s not all digits or digits.digits", input)
		}
	case model.Boolean:
		if input != "true" && input != "false" {
			return fmt.Errorf("invalid Boolean: %s not true or false", input)
		}
	case model.String:
		if !stringRegex.MatchString(input) {
			return fmt.Errorf("invalid String: %s contains special characters: [ ] ,", input)
		}
		return nil
	default:
		return fmt.Errorf("unknown atomic type: %s", atomicType)
	}
	return nil
}

func ValidateCompositeType(input string, abstractType *model.AbstractType) error {
	var parsed interface{}

	// Parse the input as JSON
	if err := json.Unmarshal([]byte(input), &parsed); err != nil {
		return fmt.Errorf("invalid composite type: %s", err)
	}

	switch model.CompositeType(abstractType.Type) {
	case model.Array, model.ListNode, model.TreeNode:
		array, ok := parsed.([]interface{})
		if !ok {
			return fmt.Errorf("expected array, got: %T", parsed)
		}
		for _, elem := range array {
			elemJSON, err := json.Marshal(elem)
			if err != nil {
				return fmt.Errorf("failed to marshal element: %v", err)
			}
			if err := ValidateAbstractType(string(elemJSON), abstractType.TypeChildren); err != nil {
				return err
			}
		}
	case model.Matrix:
		matrix, ok := parsed.([]interface{})
		if !ok {
			return fmt.Errorf("expected 2D array for Matrix, got: %T", parsed)
		}
		// Check all rows are arrays of the same length
		var rowLength int
		for i, row := range matrix {
			rowArray, ok := row.([]interface{})
			if !ok {
				return fmt.Errorf("matrix row %d is not an array", i)
			}
			if i == 0 {
				rowLength = len(rowArray)
			} else if len(rowArray) != rowLength {
				return fmt.Errorf("matrix rows have inconsistent lengths")
			}
			for _, elem := range rowArray {
				elemJSON, err := json.Marshal(elem)
				if err != nil {
					return fmt.Errorf("failed to marshal element: %v", err)
				}
				if err := ValidateAbstractType(string(elemJSON), abstractType.TypeChildren); err != nil {
					return err
				}
			}
		}
	case model.Graph:
		graph, ok := parsed.([]interface{})
		if !ok {
			return fmt.Errorf("expected array of pairs for Graph, got: %T", parsed)
		}
		for _, edge := range graph {
			pair, ok := edge.([]interface{})
			if !ok || len(pair) != 2 {
				return fmt.Errorf("invalid edge: %v", edge)
			}
			for _, node := range pair {
				nodeJSON, err := json.Marshal(node)
				if err != nil {
					return fmt.Errorf("failed to marshal node: %v", err)
				}
				if err := ValidateAbstractType(string(nodeJSON), abstractType.TypeChildren); err != nil {
					return err
				}
			}
		}
	default:
		return fmt.Errorf("unknown composite type: %s", abstractType.Type)
	}
	return nil
}

// ValidateAbstractType validates an abstract type (atomic or composite)
func ValidateAbstractType(input string, abstractType *model.AbstractType) error {
	if abstractType.TypeChildren == nil {
		// Atomic type
		// fmt.Println("is " + input + " = " + abstractType.ToPrint())
		err := ValidateAtomicType(input, abstractType.Type)
		if err != nil {
			// fmt.Println("error " + err.Error())
			return err
		}
	} else {
		// Composite type
		// fmt.Println("is " + input + " = " + abstractType.ToPrint())
		err := ValidateCompositeType(input, abstractType)
		if err != nil {
			// fmt.Println("error " + err.Error())
			return err
		}
	}
	return nil
}
