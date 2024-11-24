package coding

//generate function signature
import (
	"errors"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
)

var languageGenerators = map[model.PredefinedSupportedLanguage]func(model.Question) (string, error){
	model.Python:     GeneratePythonSignature,
	model.Java:     GenerateJavaSignature,
	model.JavaScript: GenerateJavaScriptSignature,

	// Add more languages here
}

func GenerateByQuestionAndLanguage(question model.Question, language model.PredefinedSupportedLanguage) (string, error) {
	generator, exists := languageGenerators[language]
	if !exists {
		return "", errors.New("unsupported language")
	}
	return generator(question)
}
