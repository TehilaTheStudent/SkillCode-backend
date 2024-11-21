package coding

//generate function signature
import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
	"github.com/ettle/strcase"
)

// Type mappings for JavaScript
var javascriptTypeMappings = map[string]string{
	string(model.Integer):  "number",
	string(model.Double):   "number",
	string(model.String):   "string",
	string(model.Boolean):  "boolean",
	string(model.Array):    "Array",
	string(model.Matrix):   "Array<Array>",
	string(model.Graph):    "Graph",
	string(model.TreeNode): "TreeNode",
	string(model.ListNode): "ListNode",
}

var pythonTypeMappings = map[string]string{
	string(model.Integer):  "int",
	string(model.Double):   "float",
	string(model.String):   "str",
	string(model.Boolean):  "bool",
	string(model.Array):    "list",
	string(model.Matrix):   "list[list]",
	string(model.Graph):    "Graph",
	string(model.TreeNode): "TreeNode",
	string(model.ListNode): "ListNode",
}

// mapToJSType maps abstract types to JavaScript types
func mapToJSType(paramType model.AbstractType) string {
	baseType, exists := javascriptTypeMappings[paramType.Type]
	if !exists {
		return "any"
	}

	// Handle nested structures
	if paramType.TypeChildren != nil {
		switch paramType.Type {
		case string(model.Array):
			return fmt.Sprintf("Array<%s>", mapToJSType(*paramType.TypeChildren))
		case string(model.Matrix):
			return fmt.Sprintf("Array<Array<%s>>", mapToJSType(*paramType.TypeChildren))
		case string(model.Graph), string(model.TreeNode), string(model.ListNode):
			return fmt.Sprintf("%s<%s>", baseType, mapToJSType(*paramType.TypeChildren))
		default:
			return baseType
		}
	}

	return baseType
}

// mapToPythonType maps abstract types to Python types
func mapToPythonType(paramType model.AbstractType) string {
	baseType, exists := pythonTypeMappings[paramType.Type]
	if !exists {
		return "any"
	}

	// Handle nested structures
	if paramType.TypeChildren != nil {
		switch paramType.Type {
		case string(model.Array):
			return fmt.Sprintf("list[%s]", mapToPythonType(*paramType.TypeChildren))
		case string(model.Matrix):
			return fmt.Sprintf("list[list[%s]]", mapToPythonType(*paramType.TypeChildren))
		case string(model.Graph), string(model.TreeNode), string(model.ListNode):
			return fmt.Sprintf("%s[%s]", baseType, mapToPythonType(*paramType.TypeChildren))
		default:
			return baseType
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

const jsFunctionTemplate = `/**
 * {{- range .ParamsDocs }}
 * @param {{ .Type }} {{ .Name }}
 * {{- end }}
 * @returns {{ .ReturnType }}
 */
function {{.FunctionName}}({{.Params}}) {
    // TODO: Implement this function
}`

const tsFunctionTemplate = `function {{.FunctionName}}({{.Params}}) {
    // TODO: Implement this function
}`

func GenerateJavaScriptSignature(question model.Question) (string, error) {
	// Prepare data for JSDoc and function signature
	paramList := []string{}
	paramDocs := []map[string]string{}

	for _, param := range *question.FunctionConfig.Parameters {
		paramList = append(paramList, param.Name)

		// Add JSDoc details for each parameter
		paramDocs = append(paramDocs, map[string]string{
			"Name": param.Name,
			"Type": mapToJSType(param.ParamType),
		})
	}

	data := map[string]interface{}{
		"FunctionName": ToJSStyle(question.FunctionConfig.Name),
		"Params":       strings.Join(paramList, ", "),
		"ParamsDocs":   paramDocs,
		"ReturnType":   mapToJSType(*question.FunctionConfig.ReturnType),
	}

	// Render the template
	tmpl, err := template.New("jsFunc").Parse(jsFunctionTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// ToJSStyle converts a string to JavaScript-style camelCase
func ToJSStyle(functionName string) string {
	return strcase.ToCamel(functionName)
}


var languageGenerators = map[model.PredefinedSupportedLanguage]func(model.Question) (string, error){
	model.Python:     GeneratePythonSignature,
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
