package std_test

import (
	"strings"
	"testing"
	"unicode"

	lua "github.com/speedata/go-lua"
	"github.com/stretchr/testify/assert"
)

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

			l.CreateTable(0, 0)
			if err := l.ProtectedCall(1, 1, 0); err != nil {
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
