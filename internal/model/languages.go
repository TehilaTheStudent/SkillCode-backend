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
	case Go:
		return "go"
	case CSharp:
		return "cs"
	case Cpp:
		return "cpp"
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
    case Java:
        return "java"
    case Go:
        return "go"
    case CSharp:
        return "mono-runtime"
    case Cpp:
        return "g++"
    default:
        return "" // Unsupported language
    }
}


var PredefinedSupportedLanguages = []PredefinedSupportedLanguage{
	Python,
	JavaScript,
	Java,
	Go,
	CSharp,
	Cpp,
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
	case "go":
		langEnum = Go
	case "csharp":
		langEnum = CSharp
	case "cpp":
		langEnum = Cpp
	default:
		return "", NewCustomError(400, fmt.Sprintf("unsupported language: %s", language))
	}
	return langEnum, nil
}
