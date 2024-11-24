package coding

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
	"github.com/ettle/strcase"
)

var javaTypeMappings = map[string]string{
	string(model.Integer):  "Integer",
	string(model.Double):   "Double",
	string(model.String):   "String",
	string(model.Boolean):  "Boolean",
	string(model.Array):    "List",
	string(model.Matrix):   "List<List>",
	string(model.Graph):    "Graph",
	string(model.TreeNode): "TreeNode",
	string(model.ListNode): "ListNode",
}

// mapToJavaType maps abstract types to Java types
func mapToJavaType(paramType model.AbstractType) string {
	baseType, exists := javaTypeMappings[paramType.Type]
	if !exists {
		return "Object"
	}

	// Handle nested structures
	if paramType.TypeChildren != nil {
		switch paramType.Type {
		case string(model.Array):
			return fmt.Sprintf("List<%s>", mapToJavaType(*paramType.TypeChildren))
		case string(model.Matrix):
			return fmt.Sprintf("List<List<%s>>", mapToJavaType(*paramType.TypeChildren))
		case string(model.Graph), string(model.TreeNode), string(model.ListNode):
			return baseType // These types are already defined in Java
		default:
			return baseType
		}
	}

	return baseType
}

const javaFunctionTemplate = `public class {{.ClassName}} {
    public  {{.ReturnType}} {{.FunctionName}}({{.Params}}) {
       //TODO: implement this function
    }
}`

// question -> java signature
func GenerateJavaSignature(question model.Question) (string, error) {
	// Prepare data for template
	paramList := []string{}
	for _, param := range *question.FunctionConfig.Parameters {
		paramList = append(paramList, fmt.Sprintf("%s %s", mapToJavaType(param.ParamType), ToJavaStyle(param.Name)))
	}
	data := map[string]string{
		"ClassName":    "UserSolution",
		"FunctionName": ToJavaStyle(question.FunctionConfig.Name),
		"Params":       strings.Join(paramList, ", "),
		"ReturnType":   mapToJavaType(*question.FunctionConfig.ReturnType),
	}

	// Render the template
	tmpl, err := template.New("javaFunc").Parse(javaFunctionTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// ToJavaStyle converts a string to Java-style camelCase
func ToJavaStyle(functionName string) string {
	return strcase.ToCamel(functionName)
}
