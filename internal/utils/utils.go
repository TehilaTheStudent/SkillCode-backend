package utils

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
)

// EnsureWorkingDirectory ensures the working directory is set to the project root specified by an environment variable
func EnsureWorkingDirectory() {
	projectRoot := os.Getenv("PROJECT_ROOT")
	if projectRoot == "" {
		log.Println("PROJECT_ROOT environment variable is not set")
		os.Exit(1)
	}
	if err := changeDirectory(projectRoot); err != nil {
		log.Fatalf("Failed to change directory to PROJECT_ROOT: %v", err)
	}
}

// changeDirectory changes the current working directory to the specified path
func changeDirectory(path string) error {
	if err := os.Chdir(path); err != nil {
		return fmt.Errorf("failed to change directory to %s: %w", path, err)
	}
	return nil
}

func LowerToEnum(language string) (model.PredefinedSupportedLanguage, error) {
	// Convert language to lowercase
	language = strings.ToLower(language)
	var langEnum model.PredefinedSupportedLanguage
	switch language {
	case "python":
		langEnum = model.Python
	case "javascript":
		langEnum = model.JavaScript
	case "java":
		langEnum = model.Java
	case "go":
		langEnum = model.Go
	case "csharp":
		langEnum = model.CSharp
	case "cpp":
		langEnum = model.Cpp
	default:
		return "", fmt.Errorf("unsupported language: %s", language)
	}
	return langEnum, nil
}


