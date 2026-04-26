//go:build lua_bindgen

package testdata

// +lua-bindgen.sh module=std.test
type Foo struct {
	Name string
}