package luagen

import "strings"

type GenConfig struct {
	Package     string
	TypeName    string
	Module      string
	Skip        map[string]bool
	NilMap      map[string]bool
	ForceMethod map[string]bool
	SkipFields  map[string]bool
	Imports     map[string]string
	SourceDir   string
}

type GenConfigAnnotation struct {
	Module      string
	Class       string
	Skip        map[string]bool
	NilMap      map[string]bool
	ForceMethod map[string]bool
	SkipFields  map[string]bool
}

func (a *GenConfigAnnotation) IsEmpty() bool {
	return a.Module == "" && a.Class == "" && len(a.Skip) == 0 && len(a.NilMap) == 0 && len(a.ForceMethod) == 0 && len(a.SkipFields) == 0
}

type ParamInfo struct {
	Name string
	Type string
}

type ReturnKind int

const (
	ReturnVoid ReturnKind = iota
	ReturnBool
	ReturnInt
	ReturnInt64
	ReturnString
	ReturnStringSlice
	ReturnPointer
	ReturnPointerSlice
	ReturnTupleSlice
	ReturnTuple
	ReturnComplex
)

type MethodInfo struct {
	Name           string
	Params         []ParamInfo
	ReturnType     string
	ReturnKind     ReturnKind
	IsGetter       bool
	IsSkipped      bool
	IsNilMap       bool
	IsForceMethod  bool
	IsField        bool
	IsVerbatim     bool
	Method         bool
	LuaName        string
	PtrType        string
	ErrorType      string
	Raises         bool
	RaisesType     string
	SourceFile     string
	SourceLine     int
	EmmyDoc        string
	EmmyDeprecated string
	EmmySee        string
	EmmyVersion    string
	EmmyAlias      string
	EmmyEnum       string
}

type FieldInfo struct {
	Name      string
	Type      string
	IsSkipped bool
}

type ModuleFuncInfo struct {
	Name    string
	LuaName string
	Raw     bool
}

type MetamethodInfo struct {
	Name    string
	LuaName string
	Raw     bool
}

type TypeInfo struct {
	Package     string
	Name        string
	Class       string
	Module      string
	Methods     []MethodInfo
	Fields      []FieldInfo
	ModuleFuncs []ModuleFuncInfo
	Metamethods []MetamethodInfo
	Handle      string
	ImportPkg   string
}

func (c *GenConfig) IsSkipped(name string) bool {
	return c.Skip[name]
}

func (c *GenConfig) IsNilMapped(name string) bool {
	return c.NilMap[name]
}

func (c *GenConfig) IsForceMethod(name string) bool {
	return c.ForceMethod[name]
}

func (c *GenConfig) IsFieldSkipped(name string) bool {
	return c.SkipFields[name]
}

func ToSnake(s string) string {
	var b strings.Builder
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				b.WriteByte('_')
			}
			b.WriteByte(byte(r + 'a' - 'A'))
		} else {
			b.WriteByte(byte(r))
		}
	}
	return b.String()
}

func (t *TypeInfo) VarName() string {
	return strings.ToLower(t.Name[:1]) + t.Name[1:]
}

func (t *TypeInfo) OutputFileName() string {
	return ToSnake(t.Name) + "_bindings.go"
}
