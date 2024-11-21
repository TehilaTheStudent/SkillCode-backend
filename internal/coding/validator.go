package coding

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
)

// Regex for valid names (must start with a letter or underscore, followed by letters, digits, or underscores, spaces too)
var validNameRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_ ]*$`)

// ValidateCharacters checks the validity of function and parameter names
func ValidateCharacters(question *model.Question) error {
	if question == nil {
		return errors.New("question cannot be nil")
	}

	// Validate the function name
	funcName := question.FunctionConfig.Name
	if !validNameRegex.MatchString(funcName) {
		return fmt.Errorf("invalid function name: '%s'. Must start with a letter or underscore and contain only letters, digits, or underscores", funcName)
	}

	// Validate parameters
	paramNames := make(map[string]struct{})
	if question.FunctionConfig.Parameters != nil {
		for _, param := range *question.FunctionConfig.Parameters {
			// Check for duplicate parameter names
			if _, exists := paramNames[param.Name]; exists {
				return fmt.Errorf("duplicate parameter name: '%s'", param.Name)
			}
			paramNames[param.Name] = struct{}{}

			// Validate parameter name
			if !validNameRegex.MatchString(param.Name) {
				return fmt.Errorf("invalid parameter name: '%s'. Must start with a letter or underscore and contain only letters, digits, or underscores", param.Name)
			}
		}
	}

	return nil
}
