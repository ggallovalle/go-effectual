package std

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	lua "github.com/speedata/go-lua"

	"github.com/ggallovalle/go-effectual"
)

const (
	ModTestingName      = "std.testing"
	slugTestCtxHandle   = "go/std/testing/TestCtx*"
	slugExpectHandle    = "go/std/testing/Expect*"
)

type ModTesting struct {
	name string
}

type ModTestingApi struct {
	mod *ModTesting
	lua *lua.State
}

func MakeModTesting() effectual.LuaMod[ModTestingApi] {
	return &ModTesting{name: ModTestingName}
}

func (lib *ModTesting) Name() string {
	return lib.name
}

func (lib *ModTesting) Annotations() string {
	return ""
}

func (lib *ModTesting) Api(l *lua.State) ModTestingApi {
	return ModTestingApi{mod: lib, lua: l}
}

func (lib *ModTesting) Open(l *lua.State) int {
	lua.NewLibrary(l, testingLibrary)

	lua.NewMetaTable(l, slugTestCtxHandle)
	l.PushValue(-1)
	l.SetField(-2, "__index")
	lua.SetFunctions(l, testCtxMetatable, 0)
	l.Pop(1)

	lua.NewMetaTable(l, slugExpectHandle)
	l.PushValue(-1)
	l.SetField(-2, "__index")
	lua.SetFunctions(l, expectMetatable, 0)
	l.Pop(1)

	return 1
}

func (lib *ModTesting) OpenLib(l *lua.State) {
	lua.Require(l, lib.name, lib.Open, false)
	l.Pop(1)
}

func (lib *ModTesting) Require(l *lua.State) {
	l.Global("require")
	l.PushString(lib.Name())
	l.Call(1, 1)
}

var testingLibrary = []lua.RegistryFunction{
	{Name: "ctx", Function: func(l *lua.State) int {
		if l.TypeOf(1) != lua.TypeUserData {
			lua.ArgumentError(l, 1, "*testing.T expected")
			panic("unreachable")
		}
		t, ok := l.ToUserData(1).(*testing.T)
		if !ok {
			lua.ArgumentError(l, 1, "*testing.T expected")
			panic("unreachable")
		}

		tc := &testCtx{t: t}

		if l.Top() >= 2 && l.TypeOf(2) == lua.TypeString {
			tc.suiteName, _ = l.ToString(2)
		}
		if l.Top() >= 3 && l.TypeOf(3) == lua.TypeString {
			tc.caseName, _ = l.ToString(3)
		}

		l.PushUserData(tc)
		lua.SetMetaTableNamed(l, slugTestCtxHandle)
		return 1
	}},
}

type testCtx struct {
	t         *testing.T
	suiteName string
	caseName  string
}

type CtxExtEntry struct {
	Ext    LuaTestCtxExtension
	Params map[string]any
}

var ctxExtRegistry = make(map[*testing.T][]CtxExtEntry)

func SetCtxExt(t *testing.T, entries []CtxExtEntry) {
	ctxExtRegistry[t] = entries
}

type LuaTestCtxExtension interface {
	Name() string
	Build(l *lua.State, params map[string]any)
}

func toTestCtx(l *lua.State) *testCtx {
	return lua.CheckUserData(l, 1, slugTestCtxHandle).(*testCtx)
}

var testCtxGetters = map[string]func(*lua.State){
	"name": func(l *lua.State) {
		tc := toTestCtx(l)
		if tc.caseName != "" {
			if tc.suiteName != "" {
				l.PushString(tc.suiteName + "/" + tc.caseName)
			} else {
				l.PushString(tc.caseName)
			}
		} else {
			l.PushString(tc.t.Name())
		}
	},
	"ext": func(l *lua.State) {
		tc := toTestCtx(l)
		l.NewTable()
		if entries, ok := ctxExtRegistry[tc.t]; ok {
			for _, entry := range entries {
				l.PushString(entry.Ext.Name())
				entry.Ext.Build(l, entry.Params)
				l.RawSet(-3)
			}
		}
	},
}

var testCtxMethods = map[string]lua.Function{
	"skip": func(l *lua.State) int {
		argc := l.Top()
		var note string

		switch argc {
		case 1: // only self
			l.PushString("__SKIP__" + note)
			l.Error()
		case 2: // self + one arg
			if l.TypeOf(2) == lua.TypeString {
				note, _ = l.ToString(2)
				l.PushString("__SKIP__" + note)
				l.Error()
			} else if l.IsBoolean(2) {
				if !l.ToBoolean(2) {
					return 0
				}
				l.PushString("__SKIP__" + note)
				l.Error()
			} else {
				lua.ArgumentError(l, 2, "string or boolean expected")
				panic("unreachable")
			}
		case 3: // self + two args
			if !l.IsBoolean(2) {
				lua.ArgumentError(l, 2, "boolean expected")
				panic("unreachable")
			}
			if !l.ToBoolean(2) {
				return 0
			}
			if l.IsString(3) {
				note, _ = l.ToString(3)
			} else {
				lua.ArgumentError(l, 3, "string expected")
				panic("unreachable")
			}
			l.PushString("__SKIP__" + note)
			l.Error()
		default:
			lua.ArgumentError(l, 4, "too many arguments")
			panic("unreachable")
		}

		return 0 // unreachable
	},
	"log": func(l *lua.State) int {
		tc := toTestCtx(l)
		msg, _ := l.ToString(2)

		lua.Where(l, 1)
		loc, _ := l.ToString(-1)
		l.Pop(1)

		if l.Top() >= 3 && l.IsTable(3) {
			var parts []string
			l.PushNil()
			for l.Next(3) {
				if key, ok := l.ToString(-2); ok {
					val := effectual.ConvertLuaToAny(l, -1)
					parts = append(parts, fmt.Sprintf("%s=%v", key, val))
				}
				l.Pop(1)
			}
			if len(parts) > 0 {
				tc.t.Log(loc + msg + " " + strings.Join(parts, " "))
				return 0
			}
		}

		tc.t.Log(loc + msg)
		return 0
	},
	"expect": func(l *lua.State) int {
		tc := toTestCtx(l)

		e := &expectObj{t: tc.t}

		if l.Top() >= 2 {
			if l.TypeOf(2) == lua.TypeUserData {
				e.value = l.ToUserData(2)
			} else {
				e.value = effectual.ConvertLuaToAny(l, 2)
			}
		} else {
			e.value = nil
		}

		if l.Top() >= 3 && l.TypeOf(3) == lua.TypeString {
			e.msg, _ = l.ToString(3)
		}

		lua.Where(l, 1)
		e.exprLocation, _ = l.ToString(-1)
		l.Pop(1)

		expr := extractExpectExpr(l, e.exprLocation)
		e.exprString = expr

		l.PushUserData(e)
		lua.SetMetaTableNamed(l, slugExpectHandle)
		return 1
	},
}

var testCtxMetatable = []lua.RegistryFunction{
	effectual.LuaMetaIndex(testCtxGetters, testCtxMethods),
}

type expectObj struct {
	t            *testing.T
	value        any
	msg          string
	exprLocation string
	exprString   string
}

func toExpect(l *lua.State) *expectObj {
	return lua.CheckUserData(l, 1, slugExpectHandle).(*expectObj)
}

func findLuaFrame(l *lua.State, startLevel int) (lua.Frame, int, bool) {
	for level := startLevel; ; level++ {
		ci, ok := lua.Stack(l, level)
		if !ok {
			return nil, 0, false
		}
		d, ok := lua.Info(l, "S", ci)
		if ok && d.What == "Lua" {
			return ci, level, true
		}
	}
}

func collectLocalVariables(l *lua.State, level int) map[string]string {
	result := make(map[string]string)

	for testLevel := level; testLevel <= level+10; testLevel++ {
		script := fmt.Sprintf(`
			local lvl = %d
			local results = {}
			for i = 1, 100 do
				local name, val = debug.getlocal(lvl, i)
				if name == nil then break end
				if name ~= "(vararg)" and name ~= "(temporary)" and name ~= "(C temporary)" then
					local ok, str = pcall(tostring, val)
					if ok then
						results[name] = str
					else
						results[name] = tostring(val)
					end
				end
			end
			return results
		`, testLevel)

		if err := lua.LoadString(l, script); err != nil {
			continue
		}
		if err := l.ProtectedCall(0, 1, 0); err != nil {
			l.Pop(1)
			continue
		}

		if l.IsTable(-1) {
			l.PushNil()
			for l.Next(-2) {
				if keyStr, ok := l.ToString(-2); ok {
					if valStr, ok := l.ToString(-1); ok {
						if keyStr != "(C temporary)" && keyStr != "(vararg)" && keyStr != "(temporary)" {
							result[keyStr] = valStr
						}
					}
				}
				l.Pop(1)
			}
		}
		l.Pop(1)

		if len(result) > 0 {
			break
		}
	}

	return result
}

var luaKeywords = map[string]bool{
	"and": true, "break": true, "do": true, "else": true, "elseif": true,
	"end": true, "false": true, "for": true, "function": true, "if": true,
	"in": true, "local": true, "nil": true, "not": true, "or": true,
	"repeat": true, "return": true, "then": true, "true": true, "until": true, "while": true,
}

func extractExpressionVariables(expr string) []string {
	var vars []string
	var current strings.Builder
	inIdentifier := false
	var identStart int

	for i := 0; i < len(expr); i++ {
		ch := expr[i]

		if isIdentChar(ch) {
			if !inIdentifier {
				identStart = i
			}
			current.WriteByte(ch)
			inIdentifier = true
		} else {
			if inIdentifier {
				ident := current.String()
				current.Reset()
				inIdentifier = false
				if !luaKeywords[ident] && !isMethodName(expr, identStart) {
					vars = append(vars, ident)
				}
			}
			if isOperatorChar(ch) {
				if ch == ':' || ch == '.' {
					continue
				}
				if inIdentifier {
					ident := current.String()
					current.Reset()
					inIdentifier = false
					if !luaKeywords[ident] && !isMethodName(expr, identStart) {
						vars = append(vars, ident)
					}
				}
			}
		}
	}

	if inIdentifier {
		ident := current.String()
		if !luaKeywords[ident] && !isMethodName(expr, identStart) {
			vars = append(vars, ident)
		}
	}

	return vars
}

func isIdentChar(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_'
}

func isOperatorChar(ch byte) bool {
	return strings.ContainsAny(string(ch), "+-*/%^#==~=<>[]{}();,")
}

func isMethodName(expr string, identStart int) bool {
	if identStart > 0 && (expr[identStart-1] == ':' || expr[identStart-1] == '.') {
		return true
	}
	return false
}

func expectFail(l *lua.State, expected, actual string) {
	e := toExpect(l)

	loc := e.exprLocation
	expr := e.exprString

	var msg string
	if e.msg != "" {
		msg = fmt.Sprintf("%s%s - expected %s, actual %s", loc, e.msg, expected, actual)
	} else if expr != "" {
		vars := extractExpressionVariables(expr)
		varValues := make(map[string]string)
		for _, v := range vars {
			varValues[v] = "?"
		}

		if len(vars) > 0 {
			_, level, found := findLuaFrame(l, 0)
			if found {
				locals := collectLocalVariables(l, level+1)
				for k, v := range locals {
					if _, exists := varValues[k]; exists {
						varValues[k] = v
					}
				}
			}
		}

		var varLines []string
		for _, v := range vars {
			varLines = append(varLines, fmt.Sprintf("- %s = %s", v, varValues[v]))
		}
		if len(varLines) > 0 {
			msg = fmt.Sprintf("%sexpected `%s` %s, actual %s\n%s", loc, expr, expected, actual, strings.Join(varLines, "\n"))
		} else {
			msg = fmt.Sprintf("%sexpected `%s` %s, actual %s", loc, expr, expected, actual)
		}
	} else {
		msg = fmt.Sprintf("%sexpected %s, actual %s", loc, expected, actual)
	}

	l.PushString(msg)
	l.Error()
}

func extractExpectExpr(l *lua.State, loc string) string {
	// loc is like "file.lua:10: " — format is path:line: suffix
	// Find last ": " to separate line number from trailing space
	suffixIdx := strings.LastIndex(loc, ": ")
	if suffixIdx == -1 {
		return ""
	}
	beforeSuffix := loc[:suffixIdx]
	// Find last ":" in "path:line" to extract line number
	colonIdx := strings.LastIndex(beforeSuffix, ":")
	if colonIdx == -1 {
		return ""
	}
	lineStr := beforeSuffix[colonIdx+1:]
	lineNum, err := strconv.Atoi(lineStr)
	if err != nil {
		return ""
	}

	filePath := beforeSuffix[:colonIdx]
	// Handle [string "..."] case - can't read source
	if strings.HasPrefix(filePath, "[string") {
		return ""
	}

	f, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for i := 1; scanner.Scan(); i++ {
		if i == lineNum {
			return parseExpectExpr(scanner.Text())
		}
	}
	return ""
}

func parseExpectExpr(line string) string {
	// Find expect( and extract first argument
	idx := strings.Index(line, "expect(")
	if idx == -1 {
		return ""
	}
	start := idx + len("expect(")
	depth := 1
	for i := start; i < len(line); i++ {
		switch line[i] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return strings.TrimSpace(line[start:i])
			}
		}
	}
	return ""
}

func expectValueString(e *expectObj) string {
	if e.value == nil {
		return "nil"
	}
	switch v := e.value.(type) {
	case bool:
		if v {
			return "true"
		}
		return "false"
	case float64:
		if v == float64(int64(v)) {
			return fmt.Sprintf("%d", int64(v))
		}
		return fmt.Sprintf("%g", v)
	case string:
		return fmt.Sprintf("%q", v)
	case []any:
		return "table"
	case map[string]any:
		return "table"
	default:
		return fmt.Sprintf("%T", v)
	}
}

func expectPushValue(l *lua.State, e *expectObj) {
	if e.value == nil {
		l.PushNil()
		return
	}
	switch v := e.value.(type) {
	case string:
		l.PushString(v)
	case bool:
		l.PushBoolean(v)
	case float64:
		l.PushNumber(v)
	case int:
		l.PushInteger(v)
	case int64:
		l.PushInteger(int(v))
	default:
		l.PushUserData(v)
	}
}

func luaValueToString(l *lua.State, idx int) string {
	switch l.TypeOf(idx) {
	case lua.TypeNil:
		return "nil"
	case lua.TypeBoolean:
		if l.ToBoolean(idx) {
			return "true"
		}
		return "false"
	case lua.TypeNumber:
		n, _ := l.ToNumber(idx)
		if n == float64(int64(n)) {
			return fmt.Sprintf("%d", int64(n))
		}
		return fmt.Sprintf("%g", n)
	case lua.TypeString:
		s, _ := l.ToString(idx)
		return fmt.Sprintf("%q", s)
	case lua.TypeTable:
		return "table"
	case lua.TypeFunction:
		return "function"
	case lua.TypeUserData:
		return fmt.Sprintf("userdata(%T)", l.ToUserData(idx))
	case lua.TypeThread:
		return "thread"
	default:
		return l.TypeOf(idx).String()
	}
}

var expectMethods = map[string]lua.Function{
	"is_nil": func(l *lua.State) int {
		e := toExpect(l)
		if e.value != nil {
			expectFail(l, "nil", expectValueString(e))
		}
		return 0
	},
	"not_nil": func(l *lua.State) int {
		e := toExpect(l)
		if e.value == nil {
			expectFail(l, "non-nil", "nil")
		}
		return 0
	},
	"is_true": func(l *lua.State) int {
		e := toExpect(l)
		if b, ok := e.value.(bool); !ok || !b {
			expectFail(l, "true", expectValueString(e))
		}
		return 0
	},
	"is_false": func(l *lua.State) int {
		e := toExpect(l)
		if b, ok := e.value.(bool); !ok || b {
			expectFail(l, "false", expectValueString(e))
		}
		return 0
	},
	"equals": func(l *lua.State) int {
		e := toExpect(l)
		expectPushValue(l, e)
		l.PushValue(2)
		if !l.Compare(-2, -1, lua.OpEq) {
			actual := expectValueString(e)
			expected := luaValueToString(l, -1)
			l.Pop(2)
			expectFail(l, expected, actual)
		}
		l.Pop(2)
		return 0
	},
	"not_equals": func(l *lua.State) int {
		e := toExpect(l)
		expectPushValue(l, e)
		l.PushValue(2)
		if l.Compare(-2, -1, lua.OpEq) {
			actual := expectValueString(e)
			expected := luaValueToString(l, -1)
			l.Pop(2)
			expectFail(l, "not "+expected, actual)
		}
		l.Pop(2)
		return 0
	},
}

var expectMetatable = []lua.RegistryFunction{
	effectual.LuaMetaIndex(nil, expectMethods),
}
