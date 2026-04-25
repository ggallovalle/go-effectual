package luagen

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type MethodComment struct {
	Skip       bool
	NilMap     bool
	Module     string
	Metamethod string
}

func ParseSource(sourceFile string, typeName string) (*TypeInfo, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, sourceFile, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	info := &TypeInfo{
		Package: node.Name.Name,
		Name:    typeName,
	}

	methodComments := extractMethodComments(node, typeName)

	for _, decl := range node.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if !isMethodOf(funcDecl, typeName) {
			continue
		}

		methodName := funcDecl.Name.Name
		if !ast.IsExported(methodName) {
			continue
		}

		comments := methodComments[methodName]
		if comments.Skip {
			info.Methods = append(info.Methods, MethodInfo{
				Name:      methodName,
				IsSkipped: true,
			})
			continue
		}

		if comments.Metamethod != "" {
			info.Metamethods = append(info.Metamethods, MetamethodInfo{
				Name:    methodName,
				LuaName: comments.Metamethod,
			})
		}

		params := extractParams(funcDecl.Type.Params)
		returnType, returnKind, ptrType := extractReturn(funcDecl.Type.Results)

		info.Methods = append(info.Methods, MethodInfo{
			Name:       methodName,
			Params:     params,
			ReturnType: returnType,
			ReturnKind: returnKind,
			PtrType:    ptrType,
			IsNilMap:   comments.NilMap,
		})
	}

	// Parse module-level functions (//lua:module <name>)
	for _, decl := range node.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Recv != nil {
			continue
		}
		if funcDecl.Doc == nil {
			continue
		}
		text := funcDecl.Doc.Text()
		for _, line := range strings.Split(text, "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "//lua:module ") {
				luaName := strings.TrimSpace(strings.TrimPrefix(line, "//lua:module "))
				if luaName != "" {
					info.ModuleFuncs = append(info.ModuleFuncs, ModuleFuncInfo{
						Name:    funcDecl.Name.Name,
						LuaName: luaName,
					})
				}
			} else if strings.HasPrefix(line, "lua:module ") {
				luaName := strings.TrimSpace(strings.TrimPrefix(line, "lua:module "))
				if luaName != "" {
					info.ModuleFuncs = append(info.ModuleFuncs, ModuleFuncInfo{
						Name:    funcDecl.Name.Name,
						LuaName: luaName,
					})
				}
			}
		}
	}

	// Parse struct fields
	info.Fields = extractStructFields(node, typeName)

	return info, nil
}

func isMethodOf(fd *ast.FuncDecl, typeName string) bool {
	_ = typeName
	if fd.Recv == nil || len(fd.Recv.List) == 0 {
		return false
	}
	recv := fd.Recv.List[0]
	if star, ok := recv.Type.(*ast.StarExpr); ok {
		if ident, ok := star.X.(*ast.Ident); ok {
			return ident.Name == typeName
		}
	}
	if ident, ok := recv.Type.(*ast.Ident); ok {
		return ident.Name == typeName
	}
	return false
}

func extractParams(fl *ast.FieldList) []ParamInfo {
	if fl == nil {
		return nil
	}
	var params []ParamInfo
	for _, field := range fl.List {
		typeStr := exprString(field.Type)
		for _, name := range field.Names {
			params = append(params, ParamInfo{
				Name: name.Name,
				Type: typeStr,
			})
		}
	}
	return params
}

func extractReturn(fl *ast.FieldList) (string, ReturnKind, string) {
	if fl == nil || len(fl.List) == 0 {
		return "", ReturnVoid, ""
	}
	if len(fl.List) > 1 {
		return "", ReturnComplex, ""
	}
	field := fl.List[0]
	if len(field.Names) > 1 {
		return "", ReturnComplex, ""
	}
	typeStr := exprString(field.Type)
	kind, ptrType := classifyReturn(typeStr)
	return typeStr, kind, ptrType
}

func classifyReturn(typeStr string) (ReturnKind, string) {
	switch typeStr {
	case "bool":
		return ReturnBool, ""
	case "int":
		return ReturnInt, ""
	case "int64":
		return ReturnInt64, ""
	case "string":
		return ReturnString, ""
	case "[]string":
		return ReturnStringSlice, ""
	case "[][2]string":
		return ReturnTupleSlice, ""
	}
	if strings.HasPrefix(typeStr, "[]*") {
		return ReturnPointerSlice, typeStr[3:]
	}
	// Check if it's a pointer to a known type (not built-in)
	if strings.HasPrefix(typeStr, "*") {
		return ReturnPointer, typeStr[1:]
	}
	return ReturnComplex, ""
}

func extractMethodComments(node *ast.File, typeName string) map[string]MethodComment {
	comments := make(map[string]MethodComment)

	for _, decl := range node.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || !isMethodOf(funcDecl, typeName) {
			continue
		}
		if funcDecl.Doc == nil {
			continue
		}
		text := funcDecl.Doc.Text()
		mc := MethodComment{}
		for _, line := range strings.Split(text, "\n") {
			line = strings.TrimSpace(line)
			if line == "//lua:skip" {
				mc.Skip = true
			}
			_ = line
		}
		if strings.Contains(text, "//lua:skip") {
			mc.Skip = true
		}
		if strings.Contains(text, "//lua:nil-map") {
			mc.NilMap = true
		}
		// Check for //lua:metamethod <name> or lua:metamethod <name>
		for _, line := range strings.Split(text, "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "//lua:metamethod ") {
				mc.Metamethod = strings.TrimSpace(strings.TrimPrefix(line, "//lua:metamethod "))
			} else if strings.HasPrefix(line, "lua:metamethod ") {
				mc.Metamethod = strings.TrimSpace(strings.TrimPrefix(line, "lua:metamethod "))
			}
		}
		if mc.Skip || mc.NilMap || mc.Metamethod != "" {
			comments[funcDecl.Name.Name] = mc
		}
	}
	return comments
}

func extractStructFields(node *ast.File, typeName string) []FieldInfo {
	var fields []FieldInfo

	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok || typeSpec.Name.Name != typeName {
				continue
			}
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}
			if structType.Fields == nil {
				continue
			}
			for _, field := range structType.Fields.List {
				if len(field.Names) == 0 {
					continue // embedded field, skip
				}
				for _, name := range field.Names {
					fields = append(fields, FieldInfo{
						Name: name.Name,
						Type: exprString(field.Type),
					})
				}
			}
		}
	}
	return fields
}

func exprString(e ast.Expr) string {
	switch t := e.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.BasicLit:
		return t.Value
	case *ast.StarExpr:
		return "*" + exprString(t.X)
	case *ast.ArrayType:
		if t.Len != nil {
			return "[" + exprString(t.Len) + "]" + exprString(t.Elt)
		}
		return "[]" + exprString(t.Elt)
	case *ast.SelectorExpr:
		return exprString(t.X) + "." + t.Sel.Name
	case *ast.MapType:
		return "map[" + exprString(t.Key) + "]" + exprString(t.Value)
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.ChanType:
		return "chan " + exprString(t.Value)
	default:
		return ""
	}
}
