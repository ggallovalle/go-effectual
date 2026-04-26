package std_test

import (
	"path/filepath"
	"testing"

	lua "github.com/speedata/go-lua"

	sut "github.com/ggallovalle/go-effectual/std"
	"github.com/ggallovalle/go-effectual"
)

func TestLuaSuite(t *testing.T) {
	l := lua.NewState()
	lua.OpenLibraries(l)
	effectual.LuaModOpen(l, sut.MakeModSemver())

	testFile := filepath.Join("..", "luahome", "std-test", "semver_test.lua")
	runLuaSuite(t, l, testFile)
}
