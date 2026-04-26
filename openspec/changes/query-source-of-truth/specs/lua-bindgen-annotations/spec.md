## MODIFIED Requirements

### Requirement: Inline annotation syntax
The generator SHALL parse inline `//lua:` annotations that appear directly preceding their target declaration. Each annotation lives on the member it affects.

```go
//lua: module std.serde.query

//lua: class Query
type Query struct {
    params url.Values //lua: skip
}

//lua: nil-map
func (q *Query) Get(key string) string

//lua: force-method
func (q *Query) ToString() string

//lua: module-fn new
func NewQuery() *Query
```

#### Scenario: Class annotation
- **WHEN** a struct has `//lua: class ClassName` annotation
- **THEN** the generator creates Lua class annotations for that type

#### Scenario: Module function annotation
- **WHEN** a function has `//lua: module-fn functionName` annotation
- **THEN** the generator creates a module-level function named `functionName`

#### Scenario: Instance method annotation
- **WHEN** a method has `//lua: nil-map` annotation
- **THEN** the generator marks that method as mapping empty string to nil

#### Scenario: Field annotation
- **WHEN** a struct field has `//lua: skip` annotation
- **THEN** the generator skips that field in bindings

#### Scenario: Raw implementation
- **WHEN** a function has `//lua: raw` annotation
- **THEN** the generator uses the function implementation verbatim without wrapping

### Requirement: Metamethod annotation
The generator SHALL parse `//lua: metamethod` annotations on functions to register them as Lua metamethods.

```go
//lua: metamethod __tostring
//lua: raw
func QueryToString(l *lua.State) int { ... }
```

#### Scenario: Metamethod registration
- **WHEN** a function has `//lua: metamethod __name` annotation (e.g., `__tostring`, `__pairs`)
- **THEN** the generator registers the function as that metamethod

#### Scenario: Raw metamethod
- **WHEN** a function has both `//lua: metamethod __name` and `//lua: raw` annotations
- **THEN** the generator uses the exact implementation for the metamethod

### Requirement: Supported annotation keys
The generator SHALL support the following annotation keys:
- `module`: Lua module name (on package)
- `class`: marks a type as a Lua class within the module
- `module-fn`: marks a function as a module-level function with the given Lua name (e.g., `//lua: module-fn new` on function `NewQuery` creates `query.new()`)
- `nil-map`: marks a method as mapping empty string to nil
- `force-method`: marks a method to be forced as a method (not getter)
- `skip`: marks a method or field to be skipped from bindings
- `metamethod`: marks a function as a metamethod with the given full name (e.g., `__tostring`, `__pairs`)
- `raw`: indicates the function implementation should be used verbatim

#### Scenario: All supported keys present
- **WHEN** various members have the supported annotations
- **THEN** all corresponding configuration fields are set correctly

### Requirement: CLI precedence
CLI flags SHALL take precedence over annotation values. When both CLI flag and annotation specify the same option, the CLI flag value SHALL be used.

#### Scenario: CLI overrides annotation
- **WHEN** CLI provides `--module cli.module` and annotation provides `module=anno.module`
- **THEN** the generator uses `cli.module`

### Requirement: Annotation parsing location
The generator SHALL parse annotations in `internal/luagen/parser.go` during `ParseSource()`. Inline annotations are extracted from the member's preceding comment and stored in `GenConfig`.

#### Scenario: Multiple types in same file
- **WHEN** file defines multiple types with annotations
- **THEN** each type's annotation only configures its own bindings

#### Scenario: Type without annotation
- **WHEN** type has no preceding annotations
- **THEN** generator uses only CLI flag configuration (no error)
