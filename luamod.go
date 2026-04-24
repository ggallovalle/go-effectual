package effectual

import "github.com/Shopify/go-lua"

// LuaMod is the interface implemented by types that
// can exposed a lua module.
type LuaMod[T any] interface {
	Name() string
	Annotations() string
	Open(l *lua.State) int
	OpenLib(l *lua.State)
	Require(l *lua.State)
	Api(l *lua.State) T
}

func LuaModOpen[T any](l *lua.State, mod LuaMod[T]) T {
	mod.OpenLib(l)
	return mod.Api(l)
}

func LuaMetaIndex(getters map[string]func(*lua.State), methods map[string]lua.Function) lua.Function {
	return func(l *lua.State) int {
		key := lua.CheckString(l, 2)
		if l.MetaTable(1) {
			l.Field(-1, key)
			if !l.IsNil(-1) {
				return 1
			}
			l.Pop(1)
		}
		if getter, ok := getters[key]; ok {
			getter(l)
			return 1
		}
		if method, ok := methods[key]; ok {
			l.PushGoFunction(method)
			return 1
		}
		l.PushNil()
		return 1
	}
}
