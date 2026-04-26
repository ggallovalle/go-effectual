package std_test

import (
	"path/filepath"
	"testing"

	lua "github.com/speedata/go-lua"

	sut "github.com/ggallovalle/go-effectual/std"
	"github.com/ggallovalle/go-effectual"
	"github.com/stretchr/testify/assert"
)

func TestLuaSuite(t *testing.T) {
	l := lua.NewState()
	lua.OpenLibraries(l)
	effectual.LuaModOpen(l, sut.MakeModSemver())

	testFile := filepath.Join("..", "luahome", "std-test", "semver_test.lua")
	err := lua.DoFile(l, testFile)
	if !assert.NoError(t, err) {
		t.Fatalf("failed to execute test file: %v", err)
	}

	l.PushString("cases")
	l.RawGet(-2)

	l.PushNil()
	for l.Next(-2) {
		if l.IsTable(-1) {
			l.PushString("name")
			l.RawGet(-2)
			caseName, _ := l.ToString(-1)
			l.Pop(1)

			l.PushString("fn")
			l.RawGet(-2)
			if l.IsFunction(-1) {
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
				l.Pop(1)
			} else {
				l.Pop(1)
			}
		}
		l.Pop(1)
	}
	l.Pop(1)
}