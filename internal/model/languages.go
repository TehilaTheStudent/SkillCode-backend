package model

import (
	"fmt"
	"strings"
)

// GetFileExtension returns the file extension for the given language.
func GetFileExtension(language PredefinedSupportedLanguage) string {
	switch language {
	case Python:
		return "py"
	case JavaScript:
		return "js"
	default:
		return "" // Unsupported language
	}
}


func GetRuntime(language PredefinedSupportedLanguage) string{
	switch language {
    case Python:
        return "python3"
    case JavaScript:
        return "node"
    default:
        return "" // Unsupported language
    }
}


var PredefinedSupportedLanguages = []PredefinedSupportedLanguage{
	Python,
	JavaScript,
}

func LowerToEnum(language string) (PredefinedSupportedLanguage, error) {
	// Convert language to lowercase
	language = strings.ToLower(language)
	var langEnum PredefinedSupportedLanguage
	switch language {
	case "python":
		langEnum = Python
	case "javascript":
		langEnum = JavaScript
	default:
		return "", NewCustomError(400, fmt.Sprintf("unsupported language: %s", language))
	}
	return langEnum, nil
}
