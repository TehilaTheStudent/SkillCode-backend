package coding

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
	"github.com/ettle/strcase"
)

var pythonTypeMappings = map[string]string{
	string(model.Integer):  "int",
	string(model.Double):   "float",
	string(model.String):   "str",
	string(model.Boolean):  "bool",
	string(model.Array):    "list",
	string(model.Matrix):   "list[list]",
	string(model.Graph):    "utils.Graph",
	string(model.TreeNode): "utils.TreeNode",
	string(model.ListNode): "utils.ListNode",
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
		paramList = append(paramList, fmt.Sprintf("%s: %s",  ToPythonStyle(param.Name), mapToPythonType(param.ParamType)))
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
