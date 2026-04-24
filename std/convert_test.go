package std_test

import (
	"testing"

	"github.com/Shopify/go-lua"
	sut "github.com/ggallovalle/go-effectual/std"
	"github.com/stretchr/testify/assert"
)

func Test_ConvertLuaToAny_Primitives(t *testing.T) {
	l := lua.NewState()

	t.Run("nil value", func(t *testing.T) {
		l.PushNil()
		assert.Nil(t, sut.ConvertLuaToAny(l, -1))
		l.Pop(1)
	})

	t.Run("string value", func(t *testing.T) {
		l.PushString("hello")
		result := sut.ConvertLuaToAny(l, -1)
		assert.Equal(t, "hello", result)
		l.Pop(1)
	})

	t.Run("number value", func(t *testing.T) {
		l.PushNumber(42.5)
		result := sut.ConvertLuaToAny(l, -1)
		assert.Equal(t, 42.5, result)
		l.Pop(1)
	})

	t.Run("boolean true", func(t *testing.T) {
		l.PushBoolean(true)
		result := sut.ConvertLuaToAny(l, -1)
		assert.Equal(t, true, result)
		l.Pop(1)
	})

	t.Run("boolean false", func(t *testing.T) {
		l.PushBoolean(false)
		result := sut.ConvertLuaToAny(l, -1)
		assert.Equal(t, false, result)
		l.Pop(1)
	})

	t.Run("function value returns nil", func(t *testing.T) {
		l.PushGoFunction(func(l *lua.State) int { return 0 })
		assert.Nil(t, sut.ConvertLuaToAny(l, -1))
		l.Pop(1)
	})

	t.Run("userdata returns nil", func(t *testing.T) {
		l.PushLightUserData(new(any))
		assert.Nil(t, sut.ConvertLuaToAny(l, -1))
		l.Pop(1)
	})

	t.Run("thread returns nil", func(t *testing.T) {
		l.PushThread()
		assert.Nil(t, sut.ConvertLuaToAny(l, -1))
		l.Pop(1)
	})
}

func Test_ConvertLuaToAny_Map(t *testing.T) {
	t.Run("empty table", func(t *testing.T) {
		l := lua.NewState()
		l.CreateTable(0, 0)
		result := sut.ConvertLuaToAny(l, -1)
		assert.Equal(t, map[string]any{}, result)
		l.Pop(1)
	})

	t.Run("string keys", func(t *testing.T) {
		l := lua.NewState()
		l.CreateTable(0, 2)
		l.PushString("value1")
		l.SetField(-2, "key1")
		l.PushString("value2")
		l.SetField(-2, "key2")
		result := sut.ConvertLuaToAny(l, -1)
		assert.Equal(t, map[string]any{"key1": "value1", "key2": "value2"}, result)
		l.Pop(1)
	})

	t.Run("nested tables", func(t *testing.T) {
		l := lua.NewState()
		l.CreateTable(0, 1)
		l.PushString("innerValue")
		l.SetField(-2, "nested")
		result := sut.ConvertLuaToAny(l, -1)
		assert.Equal(t, map[string]any{"nested": "innerValue"}, result)
		l.Pop(1)
	})

	t.Run("mixed types in table", func(t *testing.T) {
		l := lua.NewState()
		l.CreateTable(0, 3)
		l.PushString("string")
		l.SetField(-2, "strKey")
		l.PushNumber(123)
		l.SetField(-2, "numKey")
		l.PushBoolean(true)
		l.SetField(-2, "boolKey")
		result := sut.ConvertLuaToAny(l, -1)
		assert.Equal(t, map[string]any{"strKey": "string", "numKey": 123.0, "boolKey": true}, result)
		l.Pop(1)
	})
}

func Test_ConvertLuaToAny_Array(t *testing.T) {
	t.Run("sequential numeric keys becomes array", func(t *testing.T) {
		l := lua.NewState()
		l.CreateTable(3, 0)
		l.PushString("a")
		l.RawSetInt(-2, 1)
		l.PushString("b")
		l.RawSetInt(-2, 2)
		l.PushString("c")
		l.RawSetInt(-2, 3)
		result := sut.ConvertLuaToAny(l, -1)
		assert.Equal(t, []any{"a", "b", "c"}, result)
		l.Pop(1)
	})

	t.Run("sparse numeric keys remains array with nil", func(t *testing.T) {
		l := lua.NewState()
		l.CreateTable(3, 0)
		l.PushString("first")
		l.RawSetInt(-2, 1)
		l.PushNil()
		l.RawSetInt(-2, 2)
		l.PushString("third")
		l.RawSetInt(-2, 3)
		result := sut.ConvertLuaToAny(l, -1)
		assert.Equal(t, []any{"first", nil, "third"}, result)
		l.Pop(1)
	})
}

func Test_ConvertAnyToLua(t *testing.T) {
	l := lua.NewState()

	t.Run("nil value", func(t *testing.T) {
		sut.ConvertAnyToLua(l, nil)
		assert.Equal(t, lua.TypeNil, l.TypeOf(-1))
		l.Pop(1)
	})

	t.Run("string value", func(t *testing.T) {
		sut.ConvertAnyToLua(l, "hello")
		result, ok := l.ToString(-1)
		assert.True(t, ok)
		assert.Equal(t, "hello", result)
		l.Pop(1)
	})

	t.Run("number value float64", func(t *testing.T) {
		sut.ConvertAnyToLua(l, 42.5)
		result, ok := l.ToNumber(-1)
		assert.True(t, ok)
		assert.Equal(t, 42.5, result)
		l.Pop(1)
	})

	t.Run("number value int", func(t *testing.T) {
		sut.ConvertAnyToLua(l, 42)
		result, ok := l.ToInteger(-1)
		assert.True(t, ok)
		assert.Equal(t, 42, result)
		l.Pop(1)
	})

	t.Run("boolean true", func(t *testing.T) {
		sut.ConvertAnyToLua(l, true)
		assert.Equal(t, true, l.ToBoolean(-1))
		l.Pop(1)
	})

	t.Run("boolean false", func(t *testing.T) {
		sut.ConvertAnyToLua(l, false)
		assert.Equal(t, false, l.ToBoolean(-1))
		l.Pop(1)
	})

	t.Run("array value", func(t *testing.T) {
		sut.ConvertAnyToLua(l, []any{"a", "b", "c"})
		assert.Equal(t, lua.TypeTable, l.TypeOf(-1))
		l.Pop(1)
	})

	t.Run("map value", func(t *testing.T) {
		sut.ConvertAnyToLua(l, map[string]any{"key": "value"})
		assert.Equal(t, lua.TypeTable, l.TypeOf(-1))
		l.Pop(1)
	})
}
