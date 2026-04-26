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

func Test_LibGoTesting_Expect_Pass(t *testing.T) {
	t.Run("is_nil passes on nil", func(t *testing.T) {
		l := setupTestingCtx(t)
		err := lua.DoString(l, `ctx:expect(nil):is_nil()`)
		assert.NoError(t, err)
	})

	t.Run("not_nil passes on non-nil", func(t *testing.T) {
		l := setupTestingCtx(t)
		err := lua.DoString(l, `ctx:expect(42):not_nil()`)
		assert.NoError(t, err)
	})

	t.Run("is_true passes on true", func(t *testing.T) {
		l := setupTestingCtx(t)
		err := lua.DoString(l, `ctx:expect(true):is_true()`)
		assert.NoError(t, err)
	})

	t.Run("is_false passes on false", func(t *testing.T) {
		l := setupTestingCtx(t)
		err := lua.DoString(l, `ctx:expect(false):is_false()`)
		assert.NoError(t, err)
	})

	t.Run("equals passes on equal values", func(t *testing.T) {
		l := setupTestingCtx(t)
		err := lua.DoString(l, `ctx:expect(1):equals(1)`)
		assert.NoError(t, err)
	})

	t.Run("equals passes on equal strings", func(t *testing.T) {
		l := setupTestingCtx(t)
		err := lua.DoString(l, `ctx:expect("hello"):equals("hello")`)
		assert.NoError(t, err)
	})

	t.Run("not_equals passes on unequal values", func(t *testing.T) {
		l := setupTestingCtx(t)
		err := lua.DoString(l, `ctx:expect(1):not_equals(2)`)
		assert.NoError(t, err)
	})

	t.Run("is_lt passes when less", func(t *testing.T) {
		l := setupTestingCtx(t)
		err := lua.DoString(l, `ctx:expect(1):is_lt(2)`)
		assert.NoError(t, err)
	})

	t.Run("not_lt passes when not less", func(t *testing.T) {
		l := setupTestingCtx(t)
		err := lua.DoString(l, `ctx:expect(2):not_lt(1)`)
		assert.NoError(t, err)
	})

	t.Run("is_le passes when equal", func(t *testing.T) {
		l := setupTestingCtx(t)
		err := lua.DoString(l, `ctx:expect(1):is_le(1)`)
		assert.NoError(t, err)
	})

	t.Run("not_le passes when greater", func(t *testing.T) {
		l := setupTestingCtx(t)
		err := lua.DoString(l, `ctx:expect(2):not_le(1)`)
		assert.NoError(t, err)
	})
}

func Test_LibGoTesting_Expect_Fail(t *testing.T) {
	t.Run("is_nil fails on non-nil", func(t *testing.T) {
		l := setupTestingCtx(t)
		err := lua.DoString(l, `ctx:expect(42):is_nil()`)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expected nil")
		assert.Contains(t, err.Error(), "actual 42")
	})

	t.Run("not_nil fails on nil", func(t *testing.T) {
		l := setupTestingCtx(t)
		err := lua.DoString(l, `ctx:expect(nil):not_nil()`)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expected non-nil")
		assert.Contains(t, err.Error(), "actual nil")
	})

	t.Run("is_true fails on false", func(t *testing.T) {
		l := setupTestingCtx(t)
		err := lua.DoString(l, `ctx:expect(false):is_true()`)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expected true")
		assert.Contains(t, err.Error(), "actual false")
	})

	t.Run("is_false fails on true", func(t *testing.T) {
		l := setupTestingCtx(t)
		err := lua.DoString(l, `ctx:expect(true):is_false()`)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expected false")
		assert.Contains(t, err.Error(), "actual true")
	})

	t.Run("equals fails on unequal", func(t *testing.T) {
		l := setupTestingCtx(t)
		err := lua.DoString(l, `ctx:expect(1):equals(2)`)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expected 2")
		assert.Contains(t, err.Error(), "actual 1")
	})

	t.Run("not_equals fails on equal", func(t *testing.T) {
		l := setupTestingCtx(t)
		err := lua.DoString(l, `ctx:expect(1):not_equals(1)`)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expected not 1")
		assert.Contains(t, err.Error(), "actual 1")
	})

	t.Run("is_lt fails when not less", func(t *testing.T) {
		l := setupTestingCtx(t)
		err := lua.DoString(l, `ctx:expect(2):is_lt(1)`)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expected < 1")
		assert.Contains(t, err.Error(), "actual 2")
	})

	t.Run("not_lt fails when less", func(t *testing.T) {
		l := setupTestingCtx(t)
		err := lua.DoString(l, `ctx:expect(1):not_lt(2)`)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expected not < 2")
		assert.Contains(t, err.Error(), "actual 1")
	})

	t.Run("is_le fails when greater", func(t *testing.T) {
		l := setupTestingCtx(t)
		err := lua.DoString(l, `ctx:expect(2):is_le(1)`)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expected <= 1")
		assert.Contains(t, err.Error(), "actual 2")
	})

	t.Run("not_le fails when less or equal", func(t *testing.T) {
		l := setupTestingCtx(t)
		err := lua.DoString(l, `ctx:expect(1):not_le(1)`)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expected not <= 1")
		assert.Contains(t, err.Error(), "actual 1")
	})

	t.Run("custom msg overrides expression", func(t *testing.T) {
		l := setupTestingCtx(t)
		err := lua.DoString(l, `ctx:expect(1, "check this"):equals(2)`)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "check this")
		assert.Contains(t, err.Error(), "expected 2")
		assert.Contains(t, err.Error(), "actual 1")
		assert.NotContains(t, err.Error(), "expr:")
	})
}
