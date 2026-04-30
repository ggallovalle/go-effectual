package luagen

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractTypeAnnotations(t *testing.T) {
	tests := []struct {
		name     string
		typeName string
		want     *GenConfigAnnotation
	}{
		{
			name:     "basic_module_annotation",
			typeName: "Foo",
			want:     &GenConfigAnnotation{Module: "std.test"},
		},
		{
			name:     "skip_fields_annotation",
			typeName: "Bar",
			want:     &GenConfigAnnotation{SkipFields: map[string]bool{"FieldA": true, "FieldB": true}},
		},
		{
			name:     "nil_map_and_force_method",
			typeName: "Query",
			want: &GenConfigAnnotation{
				NilMap:      map[string]bool{"Get": true},
				ForceMethod: map[string]bool{"ToString": true, "Keys": true},
			},
		},
		{
			name:     "skip_methods_annotation",
			typeName: "MyType",
			want:     &GenConfigAnnotation{Skip: map[string]bool{"SkipMe": true, "AlsoSkip": true}},
		},
		{
			name:     "all_annotations",
			typeName: "All",
			want: &GenConfigAnnotation{
				Module:      "std.all",
				SkipFields:  map[string]bool{"A": true, "B": true},
				NilMap:      map[string]bool{"N": true},
				ForceMethod: map[string]bool{"F1": true, "F2": true},
				Skip:        map[string]bool{"S1": true, "S2": true},
			},
		},
		{
			name:     "no_annotation",
			typeName: "NoAnnotation",
			want:     &GenConfigAnnotation{},
		},
		{
			name:     "build_tag_without_marker",
			typeName: "NoMarker",
			want:     &GenConfigAnnotation{},
		},
		{
			name:     "annotation_marker_only",
			typeName: "HasMarker",
			want:     &GenConfigAnnotation{Module: "has.marker"},
		},
		{
			name:     "unknown_keys_ignored",
			typeName: "Unknown",
			want:     &GenConfigAnnotation{Module: "test", SkipFields: map[string]bool{"x": true}},
		},
		{
			name:     "whitespace_handling",
			typeName: "Whitespace",
			want:     &GenConfigAnnotation{Module: "test", SkipFields: map[string]bool{"a": true, "b": true}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sourceFile := filepath.Join("testdata", tt.name+".go")
			info, annotations, err := ParseSource(sourceFile, tt.typeName)
			if err != nil {
				t.Fatalf("ParseSource() error = %v", err)
			}
			if info == nil {
				t.Fatal("info is nil")
			}

			if tt.want.Module != "" && annotations.Module != tt.want.Module {
				t.Errorf("Module = %v, want %v", annotations.Module, tt.want.Module)
			}

			if len(tt.want.SkipFields) > 0 {
				for k, v := range tt.want.SkipFields {
					if annotations.SkipFields[k] != v {
						t.Errorf("SkipFields[%s] = %v, want %v", k, annotations.SkipFields[k], v)
					}
				}
			}

			if len(tt.want.NilMap) > 0 {
				for k, v := range tt.want.NilMap {
					if annotations.NilMap[k] != v {
						t.Errorf("NilMap[%s] = %v, want %v", k, annotations.NilMap[k], v)
					}
				}
			}

			if len(tt.want.ForceMethod) > 0 {
				for k, v := range tt.want.ForceMethod {
					if annotations.ForceMethod[k] != v {
						t.Errorf("ForceMethod[%s] = %v, want %v", k, annotations.ForceMethod[k], v)
					}
				}
			}

			if len(tt.want.Skip) > 0 {
				for k, v := range tt.want.Skip {
					if annotations.Skip[k] != v {
						t.Errorf("Skip[%s] = %v, want %v", k, annotations.Skip[k], v)
					}
				}
			}
		})
	}
}

func TestGenConfigAnnotationIsEmpty(t *testing.T) {
	ann := &GenConfigAnnotation{}
	if !ann.IsEmpty() {
		t.Error("expected IsEmpty() to return true for empty annotation")
	}

	ann.Module = "test"
	if ann.IsEmpty() {
		t.Error("expected IsEmpty() to return false when Module is set")
	}
}

func TestParseCommaList(t *testing.T) {
	result := parseCommaList("a,b,c")
	expected := map[string]bool{"a": true, "b": true, "c": true}

	for k, v := range expected {
		if result[k] != v {
			t.Errorf("parseCommaList()[%s] = %v, want %v", k, result[k], v)
		}
	}

	if len(result) != len(expected) {
		t.Errorf("parseCommaList() length = %d, want %d", len(result), len(expected))
	}

	empty := parseCommaList("")
	if len(empty) != 0 {
		t.Errorf("parseCommaList(\"\") length = %d, want 0", len(empty))
	}
}

func TestAnnotParseSource(t *testing.T) {
	tmpDir := t.TempDir()
	src := `//go:build lua_bindgen
// +lua-bindgen.sh module=test.module skip-fields=x,y skip=z

package foo

type Bar struct {
	X string
	Y int
}

func (b *Bar) Method1() string { return "" }
func (b *Bar) Method2() int { return 0 }
`

	if err := os.WriteFile(filepath.Join(tmpDir, "source.go"), []byte(src), 0644); err != nil {
		t.Fatal(err)
	}

	info, annotations, err := ParseSource(filepath.Join(tmpDir, "source.go"), "Bar")
	if err != nil {
		t.Fatalf("ParseSource() error = %v", err)
	}

	if annotations.Module != "test.module" {
		t.Errorf("Module = %v, want test.module", annotations.Module)
	}

	if !annotations.SkipFields["x"] || !annotations.SkipFields["y"] {
		t.Errorf("SkipFields = %v, want x,y", annotations.SkipFields)
	}

	if !annotations.Skip["z"] {
		t.Errorf("Skip = %v, want z", annotations.Skip)
	}

	if len(info.Methods) != 2 {
		t.Errorf("len(Methods) = %d, want 2", len(info.Methods))
	}
}

func TestCLIOverridesAnnotations(t *testing.T) {
	tmpDir := t.TempDir()
	src := `// +lua-bindgen.sh module=anno.module skip-fields=anno.field

package foo

type Test struct {
	Field string
}
`

	if err := os.WriteFile(filepath.Join(tmpDir, "source.go"), []byte(src), 0644); err != nil {
		t.Fatal(err)
	}

	info, annotations, err := ParseSource(filepath.Join(tmpDir, "source.go"), "Test")
	if err != nil {
		t.Fatalf("ParseSource() error = %v", err)
	}

	if annotations.Module != "anno.module" {
		t.Errorf("annotations.Module = %v, want anno.module", annotations.Module)
	}

	if !annotations.SkipFields["anno.field"] {
		t.Errorf("annotations.SkipFields = %v, want anno.field", annotations.SkipFields)
	}

	_ = info
}
