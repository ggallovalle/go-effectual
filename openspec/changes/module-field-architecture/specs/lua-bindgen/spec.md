## ADDED Requirements

### Requirement: Two-annotation system
The generator SHALL support two annotation systems:
- `//lua:` annotations control Go code generation (bindings, wrappers, module setup)
- `//emmylua:` annotations provide supplemental EmmyLua documentation (most can be inferred)

#### Scenario: Both annotations on same member
- **WHEN** a function has both `//lua:` and `//emmylua:` annotations
- **THEN** `//lua:` controls code generation and `//emmylua:` provides IDE documentation

### Requirement: EmmyLua annotations inferred from Go
The generator SHALL infer the following EmmyLua annotations from Go AST and `//lua:` annotations:
- `@class` from Go struct + `//lua: class`
- `@field` from Go struct fields
- `@param` from Go function params
- `@return` from Go function return types
- `@operator` from `//lua: metamethod`
- `@module` and `@meta` generated automatically

### Requirement: Explicit-only //emmylua: annotations
The following EmmyLua annotations SHALL only be provided explicitly via `//emmylua:` (cannot be inferred):
- `@deprecated [reason]` - deprecation with reason
- `@see <symbol>` - cross-reference to other symbols
- `@version <spec>` - version requirements
- `@raises <ErrorType>` - **custom** annotation documenting possible error types
- `@alias` - type alias definitions
- `@enum` - enum definitions

#### Scenario: Explicit only annotations
- **WHEN** a function has `//emmylua: deprecated "Use new API"`
- **THEN** the generator emits `---@deprecated Use new API`

### Requirement: Auto-generated @source
The generator SHALL automatically emit `@source` pointing to the Go source file and line number. No explicit annotation required.

#### Scenario: Source annotation
- **WHEN** lua-bindgen generates bindings for `std/path.go` line 42
- **THEN** the annotation includes `---@source std/path.go:42`

### Requirement: Error handling modes for (value, error) returns
Go functions returning `(value, error)` support two modes:

#### Mode A: Tuple return (default)
Without `//lua: raises`, generator emits `@return_overload`:
```lua
---@return boolean success
---@return std.path.Path?
---@return_overload true, std.path.Path
---@return_overload false, string
```

#### Scenario: Tuple return
- **WHEN** Go function returns `(*Path, error)` without `//lua: raises`
- **THEN** generated Lua returns `(success, result)` tuple

#### Mode B: Error raising
With `//lua: raises SomeType`, generator wraps with `error()`:
```lua
---@raises SomeType
function path.new()
    local result, err = go_call()
    if err ~= nil then error(err) end
    return result
end
```

#### Scenario: Error raising
- **WHEN** Go function returns `(*Path, error)` with `//lua: raises ParseError`
- **THEN** generated Lua calls `error()` on failure
- **AND** EmmyLua annotation includes `@raises ParseError`

### Requirement: Signature-based wrapping inference
The generator SHALL determine wrapping behavior based on the function signature:

| Return Type | Behavior |
|-------------|----------|
| `func(*lua.State) int` | Call verbatim as go-lua function in Open() |
| `func() T` where T is not lua.State | Generate wrapper that pushes value with key from annotation |
| `func() (T, error)` | Generate wrapper with error handling, push value on success |

#### Scenario: Go-lua function verbatim call
- **WHEN** a method on the module type has signature `func(*lua.State) int`
- **THEN** the generator calls the function verbatim in Open(), no wrapper generated

#### Scenario: Value-returning function wrapper
- **WHEN** a method on the module type has no params and returns a non-lua.State type
- **THEN** the generator creates a wrapper that pushes the method name as key and the returned value

### Requirement: Module-field inference
Methods on the module type follow the same inference rules as class methods:

| Method signature | Inferred as | Override |
|------------------|-------------|----------|
| `func(T, ...)` - has parameters | Method | - |
| `func() T` - no params, non-lua.State return | Field | `//lua: method` forces method |
| `func(*lua.State) int` - go-lua signature | Field (verbatim) | `//lua: method` forces method |
| Struct fields | Field | `//lua: skip` excludes |

```go
//lua: module std.path
type PathModule struct {
    sep string  //lua: skip  -- private, not exposed
}

// Inferred as method (has params)
func (m *PathModule) Join(parts ...string) string { ... }

// Inferred as field (no params) but forced to method
//lua: method
func (m *PathModule) New() *Query { return &Query{} }

// Name override for a method
//lua: name parse
func (m *PathModule) NewWithOptions(opts Options) *Query { ... }
```

#### Scenario: Method with params inferred correctly
- **WHEN** a module-type method has parameters
- **THEN** the generator treats it as a method without explicit annotation

#### Scenario: No-arg method inferred as field
- **WHEN** a module-type method has no parameters and returns a non-lua.State type
- **THEN** the generator treats it as a field (inferred)

#### Scenario: Override no-arg to method
- **WHEN** a module-type method has `//lua: method` annotation
- **THEN** the generator treats it as a method regardless of signature (even if no params)

#### Scenario: Field key name from method name
- **WHEN** an inferred field method has no explicit name annotation
- **THEN** the generator uses the Go method name as the Lua field key

### Requirement: Field set via SetTable
The generator SHALL use `l.SetTable(-3)` to set module fields, consuming both key and value from the stack.

### Requirement: Field in annotations template
The generator SHALL include inferred fields in the generated Lua annotations template.

#### Scenario: Template includes fields
- **WHEN** fields are declared on the module type (via inference)
- **THEN** the annotations template includes `---@field <name> <type>` for each field

### Requirement: Module type declaration
The generator SHALL parse `//lua: module <path>` to declare the module state holder type.

```go
//lua: module std.path
type PathModule struct {
    sep    string
    altSep string
}

// Inferred as field (no params)
func (m *PathModule) MainSeparator() string { return m.sep }
```

### Requirement: Class type declaration
The generator SHALL parse `//lua: class <TypeName>` to declare a userdata type with methods.

```go
//lua: class Path
type Path struct {
    path string
}

//lua: method
func (p *Path) IsAbsolute() bool { ... }

//lua: nil-map
func (p *Path) Get(key string) string { ... }
```

Note: Class methods use `//lua: method` to force method treatment. `nil-map` and `skip` still apply.

### Requirement: Module vs Class distinction

| Annotation | Type purpose | Methods become |
|------------|--------------|----------------|
| `//lua: module <path>` | Module state, namespace registration | Inferred as fields by default (no params → field, go-lua signature → field); `//lua: method` overrides to methods |
| `//lua: class <TypeName>` | Userdata metatable | Use `//lua: method` to force method; `nil-map`, `skip` still apply |

### Requirement: Module factory function
The generator SHALL produce a `Make<ModuleType>()` factory function that creates and returns the module type.

### Requirement: Module struct implements LuaModDefinition
The generated module struct SHALL implement `effectual.LuaModDefinition` with all required methods.

### Requirement: Generator as source of truth
All Lua bindings SHALL be generated entirely from Go source annotations. No separate manual binding files (`mod_*.go`) SHALL exist for modules that can be fully annotated.

#### Scenario: Path module generated entirely
- **WHEN** lua-bindgen processes `path.go` with proper annotations
- **THEN** the generated `path_bindings.go` includes all module code; `mod_path.go` SHALL NOT exist

#### Scenario: Url module generated entirely
- **WHEN** lua-bindgen processes `url.go` with proper annotations
- **THEN** the generated `url_bindings.go` includes all module code; `mod_url.go` SHALL NOT exist

### Requirement: Complete module generation
The generator SHALL produce all code necessary for a working Lua module:
- Type converters (ToLua, toType)
- Method wrappers (luaTypeMethod)
- Getter wrappers (luaTypeField)
- Library registration function (typeLibrary)
- Module struct (ModType) with Name(), Annotations(), Open(), OpenLib(), Require()
- Factory function (MakeModType)
- Annotations template (AnnotationsTmpl)

#### Scenario: Complete module structure
- **WHEN** lua-bindgen generates bindings for a type with module annotation
- **THEN** the output includes all components needed for LuaRequire registration
