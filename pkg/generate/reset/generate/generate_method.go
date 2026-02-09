package generate

import (
	"fmt"
	"go/ast"
	"io"
	"strings"
)

const receiverName = "rcvr"

var numberTypes = []string{"int", "float", "complex"}

const funcDeclarationTemplate = `
func (rcvr *%s) Reset() {
	%s
}
`

func GenerateResetMethod(w io.Writer, structName string, structType *ast.StructType) {
	funcBody := generateResetMethodBody(structType)
	funcDeclaration := fmt.Sprintf(funcDeclarationTemplate, structName, funcBody)
	fmt.Fprintln(w, funcDeclaration)
}

func generateResetMethodBody(structType *ast.StructType) string {
	funcBodySlice := []string{}
	for _, field := range structType.Fields.List {
		for _, fieldName := range field.Names {
			fieldSetString := generateFieldReset(receiverName, fieldName.Name, field.Type)
			funcBodySlice = append(funcBodySlice, fieldSetString)
		}
	}
	return strings.Join(funcBodySlice, "\n")
}

func generateFieldReset(rcvr, fieldName string, fieldType ast.Expr) string {
	switch t := fieldType.(type) {
	case *ast.Ident:
		defaultValue, isBasic := getDefautValue(t.Name)
		if isBasic {
			return fmt.Sprintf("%s.%s = %s", rcvr, fieldName, defaultValue)
		} else {
			return fmt.Sprintf("// %s TODO", fieldName)
			// typ := types.Info{}[]
			// fmt.Println(hasResetMethod(fieldType))
		}
	case *ast.StarExpr:
		newRcvr := "*" + rcvr
		fieldReset := generateFieldReset(newRcvr, fieldName, t.X)
		return fmt.Sprintf(`if %s.%s != nil {
			%s
			}`, rcvr, fieldName, fieldReset)

	case *ast.ArrayType:
		return fmt.Sprintf("%s.%s = (%s.%s)[:0]", rcvr, fieldName, rcvr, fieldName)
	case *ast.MapType:
		return fmt.Sprintf("clear(%s.%s)", rcvr, fieldName)
	default:
		return fmt.Sprintf("// type of %s doesn't support ", fieldName)
	}
}

func getDefautValue(typeName string) (string, bool) {
	for _, numberType := range numberTypes {
		if strings.Contains(typeName, numberType) || strings.Contains(typeName, "int") {
			return "0", true
		}
	}
	if typeName == "bool" {
		return "false", true
	}
	if typeName == "string" {
		return `""`, true
	}
	return `""`, false
}

// type notResetable struct{}
// type resetable struct {
// }

// func (r *resetable) Reset() {
// }

// // generate:reset
// type strct struct {
// 	iint, iint2   int
// 	bl            bool
// 	i64           int64
// 	str           string
// 	intP          *int
// 	intPP         **int
// 	strP          *string
// 	slc           []*string
// 	slcP          *[]*string
// 	mp            map[string]string
// 	mpP           *map[string]string
// 	resetable     resetable
// 	notResetable  notResetable
// 	resetableP    *resetable
// 	notResetableP *notResetable
// }
