## Context

Lua tests for `std` modules live in two places:
- `luahome/std-test/semver_test.lua` - actual Lua test suite
- `std/mod_*_test.go` - Go files with inline Lua via `lua.DoString`

The inconsistency means tests are harder to write and maintain. `mod_path_test.go` has 603 lines of inline Lua that could live in `luahome/std-test/path_test.lua`.

The current `runLuaSuite` signature is:
```go
func runLuaSuite(t *testing.T, l *lua.State, path string)
```

Extensions are not supported. Any Go-side state needed for tests must be set as globals before calling `runLuaSuite`.

## Goals / Non-Goals

**Goals:**
- `luahome/std-test/` is the single source of truth for all Lua tests
- Extension mechanism is declarative (Lua declares deps) rather than imperative (callback handles all cases)
- Go test files consolidated into single `lua_test.go` with thin wrapper functions (~20 lines total)
- Clear error messages when deps cannot be resolved

**Non-Goals:**
- Do not change the test framework itself (`mod_testing_test.go` stays as-is)
- Do not create new test assertions - `ctx.expect` API remains unchanged
- Do not modify how the Lua modules (`mod_semver.go`, `mod_path.go`, etc.) work

## Decisions

### Decision 1: `runLuaSuite` accepts variadic extensions

**Chosen**: `func runLuaSuite(t *testing.T, l *lua.State, path string, extensions ...LuaTestCtxExtension)`

**Alternatives considered**:
- No extension mechanism (current): requires globals or inline Go code for Go-side state
- Callback approach: requires callback to know about all extensions per-case

**Rationale**: Variadic is simplest API. Callers pass what they need. Runner resolves declared deps against available extensions at runtime.

### Decision 2: `LuaTestCtxExtension` interface

```go
type LuaTestCtxExtension interface {
    Name() string
    Build(l *lua.State, params map[string]any)
}
```

**Alternatives considered**:
- Closure-based: `func(l *lua.State) -> bool` (returns true if handled)
- Named constructor functions per extension

**Rationale**: Interface is self-documenting. Each extension clearly declares its name. `Build` receives params as `map[string]any` for flexibility.

### Decision 3: `suite.deps` and `case.deps` tables (both optional)

Both `deps` fields are optional. A dependency can be:
- A string shorthand: `"logger"` → `{name = "logger", params = nil}`
- A table: `{name = "logger", params = {level = "WARN"}}`

```lua
local Suite = {
    name = "std.semver",
    deps = {"logger"},  -- shorthand for {{name = "logger"}}
    cases = {
        {
            name = "Range: contains",
            deps = {
                {name = "logger", params = {level = "WARN"}},
            },
            fn = function(ctx)
                ctx.ext.logger:info("testing")
            end,
        },
        {
            name = "Version: new() valid",
            -- no deps - runs without extensions
            fn = function(ctx)
                local v = semver.new("1.2.3")
                ctx:expect(v):not_nil()
            end,
        },
    },
}
```

**Alternatives considered**:
- String array only: `deps = {"logger"}` - no params support
- Separate fields: `suite.deps` and `case.requires` - extra field to maintain

**Rationale**: Free-form params table allows extensions to receive typed configuration. Case deps override suite deps for same name - allows refinement without removal.

### Decision 4: Error on unknown dep with available list

When a test case declares `{name = "unknown"}` but no extension registers "unknown":

```
luahome/std-test/slog_test.lua:42: extension 'unknown' not found. Available: logger
```

**Rationale**: Developer immediately knows what went wrong and what's available. Line number points to exact Lua source.

### Decision 5: `ctx.ext` as injection point

Extensions inject into `ctx.ext` table, not directly into `ctx`. This namespaces extensions and avoids polluting the core context object.

### Decision 6: Single Go wrapper file

All Lua suite tests are invoked from a single `std/lua_test.go` file with one function per module:

```go
func TestLuaSuite_Semver(t *testing.T) { runLuaSuite(t, "semver") }
func TestLuaSuite_Path(t *testing.T)  { runLuaSuite(t, "path") }
func TestLuaSuite_Url(t *testing.T)   { runLuaSuite(t, "url") }
func TestLuaSuite_Slog(t *testing.T)  { runLuaSuiteWithExt(t, "slog", &LoggerExtension{Logger: slog.Default()}) }
func TestLuaSuite_Query(t *testing.T) { runLuaSuite(t, "serde/query") }
```

**Alternatives considered**:
- One Go file per module: `mod_semver_test.go`, `mod_path_test.go`, etc.
- Keep existing inline Lua in Go files

**Rationale**: Since all Go wrappers are thin, a single file reduces file count and makes it easy to see all Lua test entry points at once.

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| Params conversion Lua→Go | Use `effectual.ConvertLuaToAny` which handles Lua tables → `map[string]any` |
| Extension panics in Build | Propagates to test failure - extensions should be tested |
| Duplicate extension registration | Last write wins (same name registered twice) |
| Suite and case both declare same dep | Case overrides suite - intentional merge strategy |
