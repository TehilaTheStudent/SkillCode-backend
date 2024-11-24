package coding

import (
	"bytes"
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
	string(model.Graph):    "utils.Graph",
	string(model.TreeNode): "utils.TreeNode",
	string(model.ListNode): "utils.ListNode",
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
		paramList = append(paramList,ToJSStyle( param.Name))

		// Add JSDoc details for each parameter
		paramDocs = append(paramDocs, map[string]string{
			"Name": ToJSStyle(param.Name),
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
