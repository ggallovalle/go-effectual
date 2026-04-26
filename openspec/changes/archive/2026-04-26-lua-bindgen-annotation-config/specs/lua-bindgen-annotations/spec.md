## ADDED Requirements

### Requirement: Annotation syntax
The generator SHALL parse block comments containing `@lua-bindgen.sh` directly preceding a type declaration. The annotation block uses key=value syntax with comma-separated lists for multi-value options:

```
//go:build lua_bindgen
// +lua-bindgen.sh skip-fields=params nil-map=Get force-method=ToString,Keys,Values,Entries module=std.serde.query

type Query struct {
    // ...
}
```

#### Scenario: Single-value annotation
- **WHEN** annotation contains `module=std.serde.query`
- **THEN** the generator uses `std.serde.query` as the module name

#### Scenario: Multi-value annotation
- **WHEN** annotation contains `force-method=ToString,Keys,Values,Entries`
- **THEN** the generator forces methods ToString, Keys, Values, and Entries to be methods (not getters)

#### Scenario: Mixed single and multi-value annotations
- **WHEN** annotation contains `skip-fields=params module=std.serde.query`
- **THEN** both skip-fields and module are set correctly

### Requirement: Supported annotation keys
The generator SHALL support the following annotation keys:
- `skip-fields`: comma-separated list of struct field names to skip
- `nil-map`: comma-separated list of method names that map empty string to nil
- `force-method`: comma-separated list of method names to force as methods
- `skip`: comma-separated list of method names to skip
- `module`: Lua module name

#### Scenario: All supported keys present
- **WHEN** annotation contains `skip-fields=x nil-map=y force-method=a,b,c skip=d module=e`
- **THEN** all corresponding configuration fields are set

#### Scenario: Unknown annotation key
- **WHEN** annotation contains `unknown-key=value`
- **THEN** the generator SHALL ignore the unknown key with a warning

### Requirement: Build tag requirement
Annotations SHALL only be active when the file has a `//go:build lua_bindgen` build constraint or the annotation block includes `// +lua-bindgen.sh` (which serves as its own constraint marker).

#### Scenario: Build tag present
- **WHEN** source file contains `//go:build lua_bindgen` before annotation
- **THEN** annotations are processed

#### Scenario: Only annotation marker present
- **WHEN** source file contains `// +lua-bindgen.sh` annotation marker
- **THEN** annotations are processed

### Requirement: CLI precedence
CLI flags SHALL take precedence over annotation values. When both CLI flag and annotation specify the same option, the CLI flag value SHALL be used.

#### Scenario: CLI overrides module
- **WHEN** CLI provides `--module cli.module` and annotation provides `module=anno.module`
- **THEN** the generator uses `cli.module`

#### Scenario: CLI overrides skip-fields
- **WHEN** CLI provides `--skip-fields cli_param` and annotation provides `skip-fields=anno_param`
- **THEN** the generator skips `cli_param` not `anno_param`

#### Scenario: Annotation used when no CLI
- **WHEN** CLI does not provide `--skip-fields` and annotation provides `skip-fields=params`
- **THEN** the generator uses the annotation value `params`

### Requirement: Annotation parsing location
The generator SHALL parse annotations in `internal/luagen/parser.go` during `ParseSource()`. Annotations SHALL be extracted from the type's preceding block comment and returned as a partial `GenConfig`.

#### Scenario: Multiple types in same file
- **WHEN** file defines multiple types with annotations
- **THEN** each type's annotation only configures its own bindings

#### Scenario: Type without annotation
- **WHEN** type has no preceding `@lua-bindgen.sh` annotation block
- **THEN** generator uses only CLI flag configuration (no error)
