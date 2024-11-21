package parser_validator

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
)


const numElements = 2


func joinStrings(arr []string) string {
	// Join the array elements with commas
	return strings.Join(arr, ", ")
}
// GenerateValidString generates a valid string for a given AbstractType.
func GenerateValidString(abstractType *model.AbstractType) string {
	switch abstractType.Type {
	case string(model.Boolean):
		return "false"
	case string(model.String):
		return `"str"`
	case string(model.Integer):
		return "1"
	case string(model.Double):
		return "1.2"
	case string(model.Array), string(model.ListNode), string(model.TreeNode):
		numElements := numElements
		childValue := GenerateValidString(abstractType.TypeChildren)
		validValues := make([]string, numElements)
		for i := 0; i < numElements; i++ {
			validValues[i] = childValue
		}
		return fmt.Sprintf("[%s]", joinStrings(validValues))
	case string(model.Matrix):
		numRows := numElements
		numCols := numElements
		rowValues := make([]string, numRows)
		childValue := GenerateValidString(abstractType.TypeChildren)
		for i := 0; i < numRows; i++ {
			row := make([]string, numCols)
			for j := 0; j < numCols; j++ {
				row[j] = childValue
			}
			rowValues[i] = fmt.Sprintf("[%s]", joinStrings(row))
		}
		return fmt.Sprintf("[%s]", joinStrings(rowValues))
	case string(model.Graph):
		numEdges := numElements
		edgeValues := make([]string, numEdges)
		childValue := GenerateValidString(abstractType.TypeChildren)
		for i := 0; i < numEdges; i++ {
			edgeValues[i] = fmt.Sprintf("[%s, %s]", childValue, childValue)
		}
		return fmt.Sprintf("[%s]", joinStrings(edgeValues))
	}
	return ""
}


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
