package std_test

import (
	"strings"
	"testing"
	"unicode"

	lua "github.com/speedata/go-lua"
	"github.com/stretchr/testify/assert"
)

func pushCtxTable(l *lua.State) {
	l.CreateTable(0, 0)
	l.PushString("skip")
	l.PushGoFunction(func(l *lua.State) int {
		argc := l.Top()
		var skip bool
		var note string

		switch argc {
		case 0:
			skip = true
		case 1:
			if l.IsString(1) {
				skip = true
				note, _ = l.ToString(1)
			} else if l.IsBoolean(1) {
				skip = l.ToBoolean(1)
			} else {
				lua.ArgumentError(l, 1, "string or boolean expected")
				panic("unreachable")
			}
		case 2:
			if !l.IsBoolean(1) {
				lua.ArgumentError(l, 1, "boolean expected")
				panic("unreachable")
			}
			skip = l.ToBoolean(1)
			if l.IsString(2) {
				note, _ = l.ToString(2)
			} else {
				lua.ArgumentError(l, 2, "string expected")
				panic("unreachable")
			}
		default:
			lua.ArgumentError(l, 3, "too many arguments")
			panic("unreachable")
		}

		if !skip {
			return 0
		}

		l.PushString("__SKIP__" + note)
		l.Error()
		return 0 // unreachable
	})
	l.RawSet(-3)
}

func runLuaSuite(t *testing.T, l *lua.State, path string) {
	err := lua.DoFile(l, path)
	if !assert.NoError(t, err) {
		t.Fatalf("failed to execute test file %q: %v", path, err)
	}

	if l.TypeOf(-1) != lua.TypeTable {
		t.Fatalf("expected test file %q to return a table, got %s", path, lua.TypeNameOf(l, -1))
	}

	l.PushString("cases")
	l.RawGet(-2)
	if l.TypeOf(-1) != lua.TypeTable {
		t.Fatalf("expected suite to have a 'cases' table, got %s", lua.TypeNameOf(l, -1))
	}

	l.PushNil()
	for l.Next(-2) {
		if !l.IsTable(-1) {
			l.Pop(1)
			continue
		}

		l.PushString("name")
		l.RawGet(-2)
		caseName, _ := l.ToString(-1)
		l.Pop(1)

		sanitized := sanitizeTestName(caseName)

		t.Run(sanitized, func(t *testing.T) {
			top := l.Top()
			defer l.SetTop(top)

			l.PushString("fn")
			l.RawGet(-2)
			if !l.IsFunction(-1) {
				t.Fatalf("case %q: 'fn' is not a function", caseName)
			}

			pushCtxTable(l)
			if err := l.ProtectedCall(1, 1, 0); err != nil {
				msg := err.Error()
				if after, ok := strings.CutPrefix(msg, "__SKIP__"); ok {
					t.Skip(after)
				}
				if l.IsString(-1) {
					s, _ := l.ToString(-1)
					t.Fatalf("%s: %s", caseName, s)
				}
				t.Fatalf("%s: %v", caseName, err)
			}
			if l.IsString(-1) {
				s, _ := l.ToString(-1)
				t.Fatalf("%s: %s", caseName, s)
			}
		})

		l.Pop(1) // pop case table
	}
	l.Pop(1) // pop cases table
	l.Pop(1) // pop suite table
}

func sanitizeTestName(name string) string {
	var b strings.Builder
	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			b.WriteRune(r)
		} else {
			b.WriteRune('_')
		}
	}
	s := strings.Trim(b.String(), "_")
	if s == "" {
		return "unnamed"
	}
	return s
}
