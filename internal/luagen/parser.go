package luagen

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
)

type MethodComment struct {
	Skip        bool
	NilMap      bool
	Module      string
	Metamethod  string
	ForceMethod bool
	Method      bool
	Name        string
	Raises      bool
	RaisesType  string
	Raw         bool
}

type FuncComment struct {
	ModuleFn   string
	Metamethod string
	Raw        bool
}

type EmmyComment struct {
	Doc        string
	Deprecated string
	See        string
	Version    string
	Alias      string
	Enum       string
}

var validAnnotationKeys = map[string]bool{
	"skip-fields":   true,
	"nil-map":       true,
	"force-method":  true,
	"skip":          true,
	"module":        true,
	"skip-field":    true,
	"nil-map-field": true,
	"class":         true,
	"module-fn":     true,
	"metamethod":    true,
	"raw":           true,
	"method":        true,
	"name":          true,
	"raises":        true,
}

func ParseSource(sourceFile string, typeName string) (*TypeInfo, *GenConfigAnnotation, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, sourceFile, nil, parser.ParseComments)
	if err != nil {
		return nil, nil, err
	}

	info := &TypeInfo{
		Package: node.Name.Name,
		Name:    typeName,
	}

	annotations := extractTypeAnnotations(node, typeName)

	info.Class = annotations.Class

	if annotations.Module == "" {
		info.Module = extractPackageModule(node)
	} else {
		info.Module = annotations.Module
	}

	methodComments := extractMethodComments(node, typeName)
	emmyComments := extractEmmyComments(node, typeName)

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

		errorType := ""
		if returnKind == ReturnTuple {
			errorType = ptrType
		}

		pos := fset.Position(funcDecl.Pos())
		sourceFile := filepath.Base(sourceFile)

		emmy := emmyComments[methodName]

		info.Methods = append(info.Methods, MethodInfo{
			Name:           methodName,
			Params:         params,
			ReturnType:     returnType,
			ReturnKind:     returnKind,
			PtrType:        ptrType,
			IsNilMap:       comments.NilMap,
			IsForceMethod:  comments.ForceMethod,
			Method:         comments.Method,
			LuaName:        comments.Name,
			Raises:         comments.Raises,
			RaisesType:     comments.RaisesType,
			ErrorType:      errorType,
			SourceFile:     sourceFile,
			SourceLine:     pos.Line,
			EmmyDoc:        emmy.Doc,
			EmmyDeprecated: emmy.Deprecated,
			EmmySee:        emmy.See,
			EmmyVersion:    emmy.Version,
			EmmyAlias:      emmy.Alias,
			EmmyEnum:       emmy.Enum,
		})
	}

	// Parse module-level functions (//lua:module <name> and //lua: module-fn <name>)
	for _, decl := range node.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Recv != nil {
			continue
		}
		if funcDecl.Doc == nil {
			continue
		}
		text := funcDecl.Doc.Text()

		hasModuleFn := false
		hasMetamethod := false
		hasRaw := false
		luaName := ""
		metamethodName := ""

		for _, line := range strings.Split(text, "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "//lua:module ") {
				luaName = strings.TrimSpace(strings.TrimPrefix(line, "//lua:module "))
				hasModuleFn = true
			} else if strings.HasPrefix(line, "lua:module ") {
				luaName = strings.TrimSpace(strings.TrimPrefix(line, "lua:module "))
				hasModuleFn = true
			} else if strings.HasPrefix(line, "//lua: module-fn ") {
				luaName = strings.TrimSpace(strings.TrimPrefix(line, "//lua: module-fn "))
				hasModuleFn = true
			} else if strings.HasPrefix(line, "//lua:module-fn ") {
				luaName = strings.TrimSpace(strings.TrimPrefix(line, "//lua:module-fn "))
				hasModuleFn = true
			} else if strings.HasPrefix(line, "lua: module-fn ") {
				luaName = strings.TrimSpace(strings.TrimPrefix(line, "lua: module-fn "))
				hasModuleFn = true
			} else if strings.HasPrefix(line, "lua:module-fn ") {
				luaName = strings.TrimSpace(strings.TrimPrefix(line, "lua:module-fn "))
				hasModuleFn = true
			} else if strings.HasPrefix(line, "//lua: metamethod ") {
				metamethodName = strings.TrimSpace(strings.TrimPrefix(line, "//lua: metamethod "))
				hasMetamethod = true
			} else if strings.HasPrefix(line, "//lua:metamethod ") {
				metamethodName = strings.TrimSpace(strings.TrimPrefix(line, "//lua:metamethod "))
				hasMetamethod = true
			} else if strings.HasPrefix(line, "lua: metamethod ") {
				metamethodName = strings.TrimSpace(strings.TrimPrefix(line, "lua: metamethod "))
				hasMetamethod = true
			} else if strings.HasPrefix(line, "lua:metamethod ") {
				metamethodName = strings.TrimSpace(strings.TrimPrefix(line, "lua:metamethod "))
				hasMetamethod = true
			} else if line == "//lua: raw" || line == "//lua:raw" || line == "lua: raw" || line == "lua:raw" {
				hasRaw = true
			}
		}

		if hasModuleFn && luaName != "" {
			info.ModuleFuncs = append(info.ModuleFuncs, ModuleFuncInfo{
				Name:    funcDecl.Name.Name,
				LuaName: luaName,
				Raw:     hasRaw,
			})
		}

		if hasMetamethod && metamethodName != "" {
			info.Metamethods = append(info.Metamethods, MetamethodInfo{
				Name:    funcDecl.Name.Name,
				LuaName: metamethodName,
				Raw:     hasRaw,
			})
		}
	}

	// Parse struct fields
	info.Fields = extractStructFields(node, typeName)

	return info, annotations, nil
}

func extractTypeAnnotations(node *ast.File, typeName string) *GenConfigAnnotation {
	ann := &GenConfigAnnotation{}

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

			text := extractAnnotationText(node, typeSpec, genDecl)
			if text == "" {
				return ann
			}

			hasLuaBindgenMarker := false
			hasBuildTag := false
			hasLuaAnnotation := false

			for _, line := range strings.Split(text, "\n") {
				line = strings.TrimSpace(line)
				if strings.Contains(line, "+lua-bindgen.sh") {
					hasLuaBindgenMarker = true
				}
				if strings.Contains(line, "//go:build") && strings.Contains(line, "lua_bindgen") {
					hasBuildTag = true
				}
				if strings.HasPrefix(line, "lua: module") || strings.HasPrefix(line, "lua:module") || strings.HasPrefix(line, "lua: class") || strings.HasPrefix(line, "lua:class") {
					hasLuaAnnotation = true
				}
			}

			if !hasLuaBindgenMarker && !hasBuildTag && !hasLuaAnnotation {
				return ann
			}

			for _, line := range strings.Split(text, "\n") {
				line = strings.TrimSpace(line)
				if !strings.HasPrefix(line, "+lua-bindgen.sh") && !strings.HasPrefix(line, "lua:") {
					continue
				}

				if strings.HasPrefix(line, "lua: ") || strings.HasPrefix(line, "lua:") {
					// Parse inline lua: annotations
					line = strings.TrimPrefix(line, "lua: ")
					line = strings.TrimPrefix(line, "lua:")
					line = strings.TrimSpace(line)

					// Handle key=value format
					if strings.Contains(line, "=") {
						pairs := parseAnnotationLine(line)
						for key, value := range pairs {
							if !validAnnotationKeys[key] {
								continue
							}
							switch key {
							case "module":
								ann.Module = value
							case "class":
								ann.Class = value
							case "skip-fields":
								ann.SkipFields = parseCommaList(value)
							case "nil-map":
								ann.NilMap = parseCommaList(value)
							case "force-method":
								ann.ForceMethod = parseCommaList(value)
							case "skip":
								ann.Skip = parseCommaList(value)
							}
						}
					} else {
						// Handle key value format (e.g., "lua: class Query")
						parts := strings.SplitN(line, " ", 2)
						if len(parts) == 2 {
							key := strings.TrimSpace(parts[0])
							value := strings.TrimSpace(parts[1])
							if validAnnotationKeys[key] {
								switch key {
								case "module":
									ann.Module = value
								case "class":
									ann.Class = value
								}
							}
						}
					}
					continue
				}

				line = strings.TrimPrefix(line, "+lua-bindgen.sh")
				line = strings.TrimSpace(line)

				pairs := parseAnnotationLine(line)
				for key, value := range pairs {
					if !validAnnotationKeys[key] {
						continue
					}
					switch key {
					case "module":
						ann.Module = value
					case "class":
						ann.Class = value
					case "skip-fields":
						ann.SkipFields = parseCommaList(value)
					case "nil-map":
						ann.NilMap = parseCommaList(value)
					case "force-method":
						ann.ForceMethod = parseCommaList(value)
					case "skip":
						ann.Skip = parseCommaList(value)
					}
				}
			}
		}
	}
	return ann
}

func extractAnnotationText(node *ast.File, typeSpec *ast.TypeSpec, genDecl *ast.GenDecl) string {
	var text string

	if genDecl != nil && genDecl.Doc != nil {
		text = genDecl.Doc.Text()
	}

	if text == "" && typeSpec.Doc != nil {
		text = typeSpec.Doc.Text()
	}

	if text == "" {
		for _, c := range node.Comments {
			if strings.Contains(c.Text(), "+lua-bindgen.sh") || strings.Contains(c.Text(), "lua: module") || strings.Contains(c.Text(), "lua: class") {
				text = c.Text()
				break
			}
		}
	}

	return text
}

func extractPackageModule(node *ast.File) string {
	for _, c := range node.Comments {
		text := c.Text()
		for _, line := range strings.Split(text, "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "//lua:module ") {
				return strings.TrimSpace(strings.TrimPrefix(line, "//lua:module "))
			}
			if strings.HasPrefix(line, "lua:module ") {
				return strings.TrimSpace(strings.TrimPrefix(line, "lua:module "))
			}
			if strings.HasPrefix(line, "//lua: module ") {
				return strings.TrimSpace(strings.TrimPrefix(line, "//lua: module "))
			}
			if strings.HasPrefix(line, "lua: module ") {
				return strings.TrimSpace(strings.TrimPrefix(line, "lua: module "))
			}
		}
	}
	return ""
}

func parseCommaList(s string) map[string]bool {
	result := make(map[string]bool)
	for _, v := range strings.Split(s, ",") {
		v = strings.TrimSpace(v)
		if v != "" {
			result[v] = true
		}
	}
	return result
}

func parseAnnotationLine(line string) map[string]string {
	result := make(map[string]string)
	pairs := strings.Fields(line)
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])
		result[key] = value
	}
	return result
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
	if len(fl.List) == 2 {
		valTypeStr := exprString(fl.List[0].Type)
		errTypeStr := exprString(fl.List[1].Type)
		if errTypeStr == "error" {
			return valTypeStr, ReturnTuple, "error"
		}
		return "", ReturnComplex, ""
	}
	if len(fl.List) > 2 {
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
			if line == "//lua:skip" || line == "//lua: skip" || line == "lua:skip" || line == "lua: skip" {
				mc.Skip = true
			}
			if line == "//lua:nil-map" || line == "//lua: nil-map" || line == "lua:nil-map" || line == "lua: nil-map" {
				mc.NilMap = true
			}
			if line == "//lua:force-method" || line == "//lua: force-method" || line == "lua:force-method" || line == "lua: force-method" {
				mc.ForceMethod = true
			}
			if line == "//lua: raw" || line == "//lua:raw" || line == "lua: raw" || line == "lua:raw" {
				mc.Raw = true
			}
			if line == "//lua:method" || line == "//lua: method" || line == "lua:method" || line == "lua: method" {
				mc.Method = true
			}
			if strings.HasPrefix(line, "//lua:name ") {
				mc.Name = strings.TrimSpace(strings.TrimPrefix(line, "//lua:name "))
			} else if strings.HasPrefix(line, "//lua: name ") {
				mc.Name = strings.TrimSpace(strings.TrimPrefix(line, "//lua: name "))
			} else if strings.HasPrefix(line, "lua:name ") {
				mc.Name = strings.TrimSpace(strings.TrimPrefix(line, "lua:name "))
			} else if strings.HasPrefix(line, "lua: name ") {
				mc.Name = strings.TrimSpace(strings.TrimPrefix(line, "lua: name "))
			}
			if strings.HasPrefix(line, "//lua:metamethod ") {
				mc.Metamethod = strings.TrimSpace(strings.TrimPrefix(line, "//lua:metamethod "))
			} else if strings.HasPrefix(line, "//lua: metamethod ") {
				mc.Metamethod = strings.TrimSpace(strings.TrimPrefix(line, "//lua: metamethod "))
			} else if strings.HasPrefix(line, "lua:metamethod ") {
				mc.Metamethod = strings.TrimSpace(strings.TrimPrefix(line, "lua:metamethod "))
			}
			if strings.HasPrefix(line, "//lua:raises ") {
				mc.Raises = true
				mc.RaisesType = strings.TrimSpace(strings.TrimPrefix(line, "//lua:raises "))
			} else if strings.HasPrefix(line, "//lua: raises ") {
				mc.Raises = true
				mc.RaisesType = strings.TrimSpace(strings.TrimPrefix(line, "//lua: raises "))
			} else if strings.HasPrefix(line, "lua:raises ") {
				mc.Raises = true
				mc.RaisesType = strings.TrimSpace(strings.TrimPrefix(line, "lua:raises "))
			} else if strings.HasPrefix(line, "lua: raises ") {
				mc.Raises = true
				mc.RaisesType = strings.TrimSpace(strings.TrimPrefix(line, "lua: raises "))
			}
		}
		if mc.Skip || mc.NilMap || mc.ForceMethod || mc.Method || mc.Name != "" || mc.Metamethod != "" || mc.Raises || mc.Raw {
			comments[funcDecl.Name.Name] = mc
		}
	}
	return comments
}

func extractEmmyComments(node *ast.File, typeName string) map[string]EmmyComment {
	comments := make(map[string]EmmyComment)

	for _, decl := range node.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || !isMethodOf(funcDecl, typeName) {
			continue
		}
		if funcDecl.Doc == nil {
			continue
		}
		text := funcDecl.Doc.Text()
		ec := EmmyComment{}
		for _, line := range strings.Split(text, "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "//emmylua:doc ") {
				ec.Doc = strings.TrimSpace(strings.TrimPrefix(line, "//emmylua:doc "))
			} else if strings.HasPrefix(line, "//emmylua: doc ") {
				ec.Doc = strings.TrimSpace(strings.TrimPrefix(line, "//emmylua: doc "))
			} else if strings.HasPrefix(line, "emmylua:doc ") {
				ec.Doc = strings.TrimSpace(strings.TrimPrefix(line, "emmylua:doc "))
			} else if strings.HasPrefix(line, "emmylua: doc ") {
				ec.Doc = strings.TrimSpace(strings.TrimPrefix(line, "emmylua: doc "))
			}
			if strings.HasPrefix(line, "//emmylua:deprecated ") {
				ec.Deprecated = strings.TrimSpace(strings.TrimPrefix(line, "//emmylua:deprecated "))
			} else if strings.HasPrefix(line, "//emmylua: deprecated ") {
				ec.Deprecated = strings.TrimSpace(strings.TrimPrefix(line, "//emmylua: deprecated "))
			}
			if strings.HasPrefix(line, "//emmylua:see ") {
				ec.See = strings.TrimSpace(strings.TrimPrefix(line, "//emmylua:see "))
			} else if strings.HasPrefix(line, "//emmylua: see ") {
				ec.See = strings.TrimSpace(strings.TrimPrefix(line, "//emmylua: see "))
			}
			if strings.HasPrefix(line, "//emmylua:version ") {
				ec.Version = strings.TrimSpace(strings.TrimPrefix(line, "//emmylua:version "))
			} else if strings.HasPrefix(line, "//emmylua: version ") {
				ec.Version = strings.TrimSpace(strings.TrimPrefix(line, "//emmylua: version "))
			}
			if strings.HasPrefix(line, "//emmylua:alias ") {
				ec.Alias = strings.TrimSpace(strings.TrimPrefix(line, "//emmylua:alias "))
			} else if strings.HasPrefix(line, "//emmylua: alias ") {
				ec.Alias = strings.TrimSpace(strings.TrimPrefix(line, "//emmylua: alias "))
			}
			if strings.HasPrefix(line, "//emmylua:enum ") {
				ec.Enum = strings.TrimSpace(strings.TrimPrefix(line, "//emmylua:enum "))
			} else if strings.HasPrefix(line, "//emmylua: enum ") {
				ec.Enum = strings.TrimSpace(strings.TrimPrefix(line, "//emmylua: enum "))
			}
		}
		if ec.Doc != "" || ec.Deprecated != "" || ec.See != "" || ec.Version != "" || ec.Alias != "" || ec.Enum != "" {
			comments[funcDecl.Name.Name] = ec
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
				isSkipped := false
				if field.Comment != nil {
					text := field.Comment.Text()
					for _, line := range strings.Split(text, "\n") {
						line = strings.TrimSpace(line)
						if line == "//lua:skip-field" || line == "//lua: skip-field" || line == "lua:skip-field" || line == "lua: skip-field" {
							isSkipped = true
							break
						}
					}
				}
				for _, name := range field.Names {
					fields = append(fields, FieldInfo{
						Name:      name.Name,
						Type:      exprString(field.Type),
						IsSkipped: isSkipped,
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
