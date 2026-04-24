package effectual

import "github.com/Shopify/go-lua"

// ConvertLuaToAny converts a Lua value at the given stack index to a Go any value.
// Handles string, number, boolean, array (table with sequential integer keys),
// and dictionary (table with string keys) types. Returns nil for unhandled types.
func ConvertLuaToAny(l *lua.State, index int) any {
	switch l.TypeOf(index) {
	case lua.TypeString:
		if v, ok := l.ToString(index); ok {
			return v
		}
	case lua.TypeNumber:
		if v, ok := l.ToNumber(index); ok {
			return v
		}
	case lua.TypeBoolean:
		return l.ToBoolean(index)
	case lua.TypeTable:
		l.PushValue(index)
		defer l.Pop(1)

		if !l.IsTable(-1) {
			return nil
		}

		length := l.RawLength(-1)
		if length > 0 {
			isArray := true
			l.PushNil()
			for l.Next(-2) {
				if !l.IsNumber(-2) {
					isArray = false
					l.Pop(2)
					break
				}
				l.Pop(1)
			}
			if isArray {
				result := make([]any, length)
				for i := 1; i <= length; i++ {
					l.RawGetInt(-1, i)
					result[i-1] = ConvertLuaToAny(l, -1)
					l.Pop(1)
				}
				return result
			}
		}

		result := make(map[string]any)
		l.PushNil()
		for l.Next(-2) {
			if key, ok := l.ToString(-2); ok {
				result[key] = ConvertLuaToAny(l, -1)
			}
			l.Pop(1)
		}
		return result
	}
	return nil
}

// ConvertAnyToLua converts a Go any value to a Lua value pushed onto the stack.
// Handles string, bool, float64, float32, int, int64, uint, uint64, []any (as array),
// and map[string]any (as table). Pushes nil for nil values and unhandled types.
func ConvertAnyToLua(l *lua.State, value any) {
	if value == nil {
		l.PushNil()
		return
	}
	switch v := value.(type) {
	case string:
		l.PushString(v)
	case bool:
		l.PushBoolean(v)
	case float64:
		l.PushNumber(v)
	case float32:
		l.PushNumber(float64(v))
	case int:
		l.PushInteger(v)
	case int64:
		l.PushInteger(int(v))
	case uint:
		l.PushInteger(int(v))
	case uint64:
		l.PushInteger(int(v))
	case []any:
		l.CreateTable(len(v), 0)
		for i, item := range v {
			l.PushInteger(i + 1)
			ConvertAnyToLua(l, item)
			l.SetTable(-3)
		}
	case map[string]any:
		l.CreateTable(0, len(v))
		for key, val := range v {
			l.PushString(key)
			ConvertAnyToLua(l, val)
			l.SetTable(-3)
		}
	default:
		l.PushNil()
	}
}
