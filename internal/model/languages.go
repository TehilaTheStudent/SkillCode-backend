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
	case Java:
		return "java"
	default:
		return "" // Unsupported language
	}
}



var PredefinedSupportedLanguages = []PredefinedSupportedLanguage{
	Python,
	JavaScript,
	Java,
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
	case "java":
		langEnum = Java
	default:
		return "", NewCustomError(400, fmt.Sprintf("unsupported language: %s", language))
	}
	return langEnum, nil
}
