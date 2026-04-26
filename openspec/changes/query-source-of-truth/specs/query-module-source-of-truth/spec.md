## ADDED Requirements

### Requirement: Query module source of truth
The Query module bindings SHALL be generated entirely from annotations in `query.go`. No separate manual binding file (`mod_query.go`) SHALL exist.

#### Scenario: Complete module generation
- **WHEN** lua-bindgen processes `query.go` with proper annotations
- **THEN** the generated bindings include module-level functions, class methods, and metamethods

### Requirement: Module-level functions
The generator SHALL create module-level Lua functions from annotated Go package-level functions.
The function signatures (parameters and return types) should reuse the same euristhic that already
exist for methods

```go
//lua: module-fn new
func NewQuery() *Query { ... }
```

#### Scenario: Constructor function
- **WHEN** a function has `//lua: module-fn new` annotation
- **THEN** the generator creates Lua function `query.new()` that wraps the Go function

#### Scenario: Factory function with parameters
- **WHEN** a function has `//lua: module-fn deserialize` annotation
- **THEN** the generator creates Lua function `query.deserialize(raw)` that calls the Go function with the provided arguments

### Requirement: Verbatim metamethods
The generator SHALL register annotated functions as Lua metamethods when marked with `//lua: raw`.

```go
//lua: metamethod __tostring
//lua: raw
func QueryToString(l *lua.State) int { ... }
```

#### Scenario: Pairs metamethod with custom iteration
- **WHEN** a function has `//lua: metamethod __pairs` and `//lua: raw` annotations
- **THEN** the generated `__pairs` metamethod uses the exact provided implementation

### Requirement: Generated ModQuery struct
The generator SHALL produce a `ModQuery` struct that implements `LuaModDefinition`.

#### Scenario: ModQuery implements LuaModDefinition
- **WHEN** lua-bindgen generates bindings for a module
- **THEN** the generated `ModQuery` struct has `Name()`, `Annotations()`, `Open()`, `OpenLib()`, and `Require()` methods

### Requirement: Generated annotations template
The generator SHALL produce an annotations template and `Annotations()` function.
#### Scenario: Template mirrors existing format
- **WHEN** the generator creates annotations for the Query module
- **THEN** the output format is compatible with existing test expectations.  The annotations template should use go `` string syntax not string concatanation (+).

#### Scenario: Module functions registered
- **WHEN** module functions are defined via `//lua: module-fn`
- **THEN** each is registered with its Lua name in the library
