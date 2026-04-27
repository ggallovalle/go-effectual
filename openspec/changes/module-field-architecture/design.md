## Context

The generator currently produces method wrappers and library functions, but module registration happens in separate manual `mod_*.go` files. This creates drift between annotations and actual behavior. The `query_bindings.go` experiment showed it's possible to generate everything.

Additionally, EmmyLua annotation files are manually maintained separately from Go code, creating a second source of truth that can drift.

**Goal**: Generator is sole source of truth - all Go bindings AND EmmyLua annotation files generated from annotations.

## Goals / Non-Goals

**Goals:**
- Generator produces complete module code including `Open()` and module struct
- Signature-based inference for module fields vs methods
- Module type (`//lua: module <path>`) holds module state, enables platform detection
- Generator produces EmmyLua annotation files with `@return_overload` for Go error patterns
- Auto-generated `@source` annotations pointing to Go source

**Non-Goals:**
- Multi-module per file (only one module per source file)
- Runtime module re-registration
- Changes to the Lua VM integration (go-lua library)
- `@async`, `@nodiscard`, `@cast`, `@diagnostic`, `@generic` annotations (don't map from Go semantics)

## Decisions

### Decision 1: Two-Annotation System

**Choice**: Keep `//lua:` for code generation, add `//emmylua:` for supplemental EmmyLua documentation.

```go
//emmylua: doc "The platform path separator"
func (m *PathModule) MainSeparator() string
```

Return types for `//lua:` annotations are inferred from Go AST. `//lua: returns <type>` is not needed.

**Rationale**: Clean separation of concerns. `//lua:` controls what code is generated, with return types inferred from Go. `//emmylua:` provides documentation the generator can't infer (doc strings, deprecation reasons, see references).

### Decision 2: What //emmylua: provides

**Choice**: `//emmylua:` is for explicit-only items. Most EmmyLua annotations are inferred.

```go
//emmylua: doc "..."              # documentation string
//emmylua: deprecated "..."      # deprecation with reason
//emmylua: see OtherSymbol        # cross-reference
//emmylua: version >=2.0.0        # version requirement
//emmylua: alias MyType ...       # type alias definition
//emmylua: enum MyEnum { ... }    # enum definition
```

**Rationale**: Reduces annotation burden. User only specifies what can't be inferred from Go AST.

### Decision 3: Error Handling Modes

**Choice**: Go functions returning `(value, error)` have two modes:

**Mode A: Tuple return (default)**
Without `//lua: raises`, generator emits:
```lua
---@return boolean success
---@return MyType?
---@return_overload true, MyType
---@return_overload false, string
```
Lua usage:
```lua
local ok, result = path.new("/foo")
if not ok then
    print("Error:", result)
end
```

**Mode B: Error raising**
With `//lua: raises SomeType`, generator wraps with error-raising:
```lua
---@raises SomeType
function path.new()
    local result, err = go_call()
    if err ~= nil then error(err) end
    return result
end
```
Lua usage:
```lua
local result = path.new("/foo")  -- raises if error
```

**Rationale**: Some APIs should raise errors (惯用 Lua), others return tuples (for conditional handling). `//lua: raises` triggers error-raising mode AND emits the EmmyLua `@raises` annotation.

### Decision 4: Auto-generated @source

**Choice**: Generator automatically emits `---@source <file>:<line>` for each annotated element.

**Rationale**: No annotation needed - the generator knows the source location from parsing. Maintains traceability from generated Lua back to Go source.

### Decision 5: Module-field inference

**Choice**: Module-type methods are inferred as follows:

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

**Rationale**: Most methods with params are correctly inferred. `//lua: method` only needed to override default field inference when a no-arg method should be callable as a function.

### Decision 6: Signature inference for wrapping

**Choice**: Generator inspects return type to determine behavior:

| Signature | Behavior |
|-----------|----------|
| `func(*lua.State) int` | Verbatim go-lua call in Open() |
| `func() T` where T != lua.State | Generate wrapper pushing value |
| `func() (T, error)` | Generate wrapper with error handling + @return_overload |

**Rationale**: `(*lua.State) int` is distinctive for go-lua functions. Other return types need wrapping to push values to Lua stack.

### Decision 7: Module type vs Class type

**Choice**: Two distinct types serve two purposes:

| Annotation | Type | Purpose |
|------------|------|---------|
| `//lua: class <Name>` | Userdata metatable | Exposes methods to Lua (FromRaw, Get, Set on Query) |
| `//lua: module <path>` | Module state holder | Registers functions/fields to Lua namespace |

```go
// Module type - holds state, registers module functions
//lua: module std.serde.query
type QueryModule struct{}

//lua: method
func (m *QueryModule) New() *Query { return &Query{} }

// Class type - userdata with methods
//lua: class Query
type Query struct {
	params url.Values //lua: skip
}

//lua: nil-map
func (q *Query) FromRaw(raw string) { ... }
```

**Rationale**: Query's methods (FromRaw, Get, Set) operate on query data - they should be class methods. The module functions (new, deserialize) operate on the module namespace - they should be module functions. Separating concerns makes generated code cleaner.

### Decision 8: Annotations template using Go template literals

**Choice**: The generated annotations template uses Go backtick string syntax, not string concatenation.

**Rationale**: Matches existing pattern in manually-written code. Easier to read and maintain.

## Generated Code Structure

### Module type

```go
//lua: module std.path
type ModPath struct {
    name string
    sep  string
    altSep string
}

func (lib *ModPath) Name() string { return lib.name }
func (lib *ModPath) Annotations() string { /* ... */ }

// Inferred as field (no params)
func (lib *ModPath) MainSeparator() string { return lib.sep }

func (lib *ModPath) SetupMainSeparator(l *lua.State) int {
    l.PushString("MAIN_SEPARATOR")
    l.PushString(lib.MainSeparator())
    l.SetTable(-3)
    return 1
}

// Inferred as field (go-lua signature)
func (lib *ModPath) SetupPosix(l *lua.State) int { /* ... */ }

func (lib *ModPath) Open(l *lua.State) int {
    lua.NewLibrary(l, pathLibrary(lib.sep))
    moduleIdx := l.AbsIndex(-1)

    // Module fields - inferred from signature
    lib.SetupMainSeparator(l)  // calls verbatim, sets field
    lib.SetupPosix(l)          // calls verbatim, creates sub-module

    lua.NewMetaTable(l, PATH_HANDLE)
    l.PushValue(-1)
    l.SetField(-2, "__index")
    lua.SetFunctions(l, pathMetatable, 0)
    l.Pop(1)
    return 1
}

func MakeModPath() effectual.LuaModDefinition {
    sep := nativeSep()
    return &ModPath{name: "std.path", sep: sep, altSep: altSep(sep)}
}
```

### Field wrappers (inferred, not annotated)

```go
// For string-returning function:
func luaModPathMainSeparator(l *lua.State) int {
    m := toModPath(l, 1)
    l.PushString("MAIN_SEPARATOR")
    l.PushString(m.MainSeparator())
    l.SetTable(-3)
    return 0
}

// For go-lua function - called verbatim, no wrapper:
func (lib *ModPath) SetupPosix(l *lua.State) int {
    lua.NewLibrary(l, pathLibrary(posixSep))
    posixIdx := l.AbsIndex(-1)
    l.PushString("MAIN_SEPARATOR")
    l.PushString(posixSep)
    l.SetTable(posixIdx)
    l.SetTable(moduleIdx)
    return 1
}
```

## Generated EmmyLua Annotations

### Path module annotations

```lua
---@meta std.path
---@source std/path.go:42

---@class (exact) std.path.Path : userdata
---@operator div(std.path.Path|string): std.path.Path
---@operator concat(std.path.Path|string): string
---@field parent std.path.Path|nil
---@field components string[]
-- ... more fields

local Path = {}

---@param path string
---@return std.path.Path
function Path:join(path) end

local path = {}

---@type string
path.MAIN_SEPARATOR = "/"

---@class std.path.posix : std.path
---@field MAIN_SEPARATOR string
local posix = {}

---@class std.path.win32 : std.path
---@field MAIN_SEPARATOR string
local win32 = {}

return path
```

### Error-returning function annotations

For Go function:
```go
//lua: method deserialize
func Deserialize(raw string) (*Query, error)
```

Generator emits:
```lua
---@param raw string
---@return boolean success
---@return std.serde.query.Query?
---@return_overload true, std.serde.query.Query
---@return_overload false, string
function query.deserialize(raw) end
```

This enables IDE type narrowing:
```lua
local ok, result = query.deserialize("?foo=bar")
if not ok then
    print("Error:", result)  -- IDE knows result is string here
    return
end
print("Query:", result)  -- IDE knows result is Query here
```

## Risks / Trade-offs

[Risk] **Annotation complexity** → Users need to understand both `//lua:` and `//emmylua:`. Mitigation: `//emmylua:` is minimal, most things inferred.

[Risk] **`@return_overload` IDE support** → Depends on EmmyLua IDE properly handling the type narrowing. Mitigation: Test with actual IDE before committing.

[Risk] **Path platform detection timing** → The `MakeModPath()` factory calls `nativeSep()` at init time. This is fine since `runtime.GOOS` doesn't change.

## Migration Plan

1. Add `//lua: module`, `//lua: class`, `//lua: method`, `//lua: name` parsing to `internal/luagen/parser.go`
2. Add field inference logic (no params → field, has params → method) to `internal/luagen/generator.go`
3. Add field wrapper generation and verbatim call generation for go-lua signatures
4. Add `@return_overload` emission for `(value, error)` returns
5. Add auto-generated `@source` annotations
6. Update `path.go` with new module/class annotations, regenerate
7. Verify `go test ./std/...` passes
8. Update `url.go` with new annotations, regenerate
9. Remove `mod_path.go` and `mod_url.go`
10. Update `lua-bindgen.sh` to produce annotation files

## Open Questions

1. **Field ordering in Open()** — Annotation order preserved based on Go method declaration order. Implemented in task 3.6.
