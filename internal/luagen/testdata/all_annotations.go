//go:build lua_bindgen

package testdata

// +lua-bindgen.sh module=std.all skip-fields=A,B nil-map=N force-method=F1,F2 skip=S1,S2
type All struct{}