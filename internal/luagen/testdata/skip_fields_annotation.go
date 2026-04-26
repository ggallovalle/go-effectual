package testdata

// +lua-bindgen.sh skip-fields=FieldA,FieldB
type Bar struct {
	FieldA string
	FieldB int
}