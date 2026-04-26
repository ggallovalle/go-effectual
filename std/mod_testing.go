package std

import (
	"fmt"
	"strings"
	"testing"

	lua "github.com/speedata/go-lua"

	"github.com/ggallovalle/go-effectual"
)

const (
	ModTestingName     = "std.testing"
	slugTestCtxHandle = "go/std/testing/TestCtx*"
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
}

var testCtxMetatable = []lua.RegistryFunction{
	effectual.LuaMetaIndex(testCtxGetters, testCtxMethods),
}
