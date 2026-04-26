//go:build lua_bindgen

package testdata

// +lua-bindgen.sh nil-map=Get force-method=ToString,Keys
type Query struct{}