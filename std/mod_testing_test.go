package std_test

import (
	"testing"

	lua "github.com/speedata/go-lua"
	"github.com/stretchr/testify/assert"

	"github.com/ggallovalle/go-effectual"
	sut "github.com/ggallovalle/go-effectual/std"
)

func Test_LibGoTesting(t *testing.T) {
	t.Run("ctx.name returns test name", func(t *testing.T) {
		l := setupTestingCtx(t)

		err := lua.DoString(l, `
			local n = ctx.name
			assert(type(n) == "string", "expected string, got " .. type(n))
			assert(n:match("ctx%.name_returns_test_name$"), "got: " .. tostring(n))
		`)
		assert.NoError(t, err)
	})

	t.Run("ctx:log does not error", func(t *testing.T) {
		l := setupTestingCtx(t)

		err := lua.DoString(l, `
			ctx:log("hello world")
			ctx:log("with attrs", {key = "value", num = 42})
		`)
		assert.NoError(t, err)
	})

	t.Run("ctx:skip unconditional", func(t *testing.T) {
		l := setupTestingCtx(t)

		err := lua.DoString(l, `ctx:skip()`)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "__SKIP__")
	})

	t.Run("ctx:skip with note", func(t *testing.T) {
		l := setupTestingCtx(t)

		err := lua.DoString(l, `ctx:skip("not yet")`)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "__SKIP__not yet")
	})

	t.Run("ctx:skip conditional no-op", func(t *testing.T) {
		l := setupTestingCtx(t)

		err := lua.DoString(l, `
			ctx:skip(false)
			ctx:skip(false, "should not skip")
		`)
		assert.NoError(t, err)
	})

	t.Run("ctx:skip conditional skip", func(t *testing.T) {
		l := setupTestingCtx(t)

		err := lua.DoString(l, `ctx:skip(true, "conditional")`)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "__SKIP__conditional")
	})

	t.Run("ctx:skip invalid arg errors", func(t *testing.T) {
		l := setupTestingCtx(t)

		err := lua.DoString(l, `ctx:skip(123)`)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "string or boolean expected")
	})
}

func setupTestingCtx(t *testing.T) *lua.State {
	l := lua.NewState()
	lua.OpenLibraries(l)
	effectual.LuaModOpen(l, sut.MakeModTesting())

	l.Global("require")
	l.PushString("std.testing")
	l.Call(1, 1)
	l.PushString("ctx")
	l.RawGet(-2)
	l.PushUserData(t)
	l.Call(1, 1)
	l.SetGlobal("ctx")
	l.Pop(1) // pop std.testing module

	return l
}
