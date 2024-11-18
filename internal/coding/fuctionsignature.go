package coding

//generate function signature
import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"strings"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
	"github.com/ettle/strcase"
)

var typeMappings = map[string]string{
	string(model.Integer):   "int",
	string(model.Double):    "float",
	string(model.String):    "str",
	string(model.Boolean):   "bool",
	string(model.Array):     "list",
	string(model.Matrix):    "list[list]",
	string(model.GraphNode): "GraphNode",
	string(model.TreeNode):  "TreeNode",
	string(model.ListNode):  "ListNode",
}

// abstract type-> string
func mapToPythonType(paramType model.AbstractType) string {
	baseType, exists := typeMappings[paramType.Type]
	if !exists {
		return "any"
	}

	// Handle children types recursively (for nested structures like list, TreeNode, etc.)
	if paramType.TypeChildren != nil {
		if paramType.Type == string(model.Array) {
			return fmt.Sprintf("list[%s]", mapToPythonType(*paramType.TypeChildren))
		} else if paramType.Type == string(model.Matrix) {
			return fmt.Sprintf("list[list[%s]]", mapToPythonType(*paramType.TypeChildren))
		} else {
			return fmt.Sprintf("%s[%s]", baseType, mapToPythonType(*paramType.TypeChildren))
		}
	}

	return baseType
}

const pythonFunctionTemplate = `def {{.FunctionName}}({{.Params}}) -> {{.ReturnType}}:`

// question -> python signature
func GeneratePythonSignature(question model.Question) (string, error) {
	// Prepare data for template
	paramList := []string{}
	for _, param := range *question.FunctionConfig.Parameters {
		paramList = append(paramList, fmt.Sprintf("%s: %s", param.Name, mapToPythonType(param.ParamType)))
	}
	data := map[string]string{
		"FunctionName": ToPythonStyle(question.FunctionConfig.Name),
		"Params":       strings.Join(paramList, ", "),
		"ReturnType":   mapToPythonType(*question.FunctionConfig.ReturnType),
	}

	// Render the template
	tmpl, err := template.New("pythonFunc").Parse(pythonFunctionTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// ToPythonStyle converts a string to Python-style snake_case
func ToPythonStyle(functionName string) string {
	return strcase.ToSnake(functionName)
}

// Placeholder for JavaScript signature generation
func GenerateJavaScriptSignature(question model.Question) (string, error) {
	// TODO: Implement JavaScript signature generation
	return fmt.Sprintf("function %s(...) { /* JavaScript implementation */ }", question.FunctionConfig.Name), nil
}

// Placeholder for Java signature generation
func GenerateJavaSignature(question model.Question) (string, error) {
	// TODO: Implement Java signature generation
	return fmt.Sprintf("public %s %s(...) { }", "void", question.FunctionConfig.Name), nil
}

var languageGenerators = map[model.PredefinedSupportedLanguage]func(model.Question) (string, error){
	model.Python: GeneratePythonSignature,
	// Add more languages here
}

func GenerateByQuestionAndLanguage(question model.Question, language model.PredefinedSupportedLanguage) (string, error) {
	generator, exists := languageGenerators[language]
	if !exists {
		return "", errors.New("unsupported language")
	}
	return generator(question)
}
