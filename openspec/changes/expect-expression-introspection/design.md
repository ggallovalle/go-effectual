## Context

The test framework in `std/mod_testing.go` provides `ctx:expect()` assertions for Lua test suites. When assertions fail, error messages include the source location and the expression text, but not the actual values of variables referenced in the expression.

Example current output:
```
lua_suite_test.go:80: Range: intersect: ../luahome/std-test/semver_test.lua:93: expected (expr: r:contains(v1)) false, actual true
```

Desired output:
```
lua_suite_test.go:80: Range: intersect: ../luahome/std-test/semver_test.lua:93: expected `r:contains(v1)` false, actual true
- r = >=1.0.0 AND <2.0.0
- v1 = 1.5.0
```

## Goals / Non-Goals

**Goals:**
- Show variable values in failure messages when assertion fails
- Use go-lua's introspection APIs to resolve variables from the call stack
- Respect Lua's `__tostring` metamethods for value formatting
- Only activate when no custom message is provided (existing behavior)

**Non-Goals:**
- Support for complex expression evaluation (function calls, nested calls) - show `?` for unresolvable
- Modifying Lua state (no adding globals or test data)
- Performance optimization (caching, etc.) - happens only on failure

## Decisions

### 1. Expression location captured at expect() time

**Decision**: Store the source location when `expect()` is called, not when `expectFail()` is called.

**Rationale**: The call site location doesn't change between expect() call and failure. Capturing it at expect() time is cleaner and matches the natural flow.

**Alternatives considered**:
- Parse at failure time (current approach): Works but re-parses on every failure
- Parse once and cache: Adds complexity for marginal benefit

### 2. Variable resolution via Lua script execution

**Decision**: Call `debug.getlocal` via `lua.LoadString` + `ProtectedCall` from within Lua context, not via go-lua's Go API.

**Rationale**: go-lua's internal `debug.getlocal` causes panics when called directly from Go due to call stack inconsistency (Go frames interspersed with Lua frames). Executing a Lua script that calls `debug.getlocal` runs in the correct Lua context.

**Implementation**:
```go
script := fmt.Sprintf(`
    local lvl = %d
    local results = {}
    for i = 1, 100 do
        local name, val = debug.getlocal(lvl, i)
        if name == nil then break end
        if name ~= "(vararg)" and name ~= "(temporary)" and name ~= "(C temporary)" then
            local ok, str = pcall(tostring, val)
            if ok then
                results[name] = str
            else
                results[name] = tostring(val)
            end
        end
    end
    return results
`, testLevel)

lua.LoadString(l, script)
l.ProtectedCall(0, 1, 0)
```

### 3. Stack level search required

**Decision**: Search multiple stack levels (level to level+10) to find the correct Lua frame with local variables.

**Rationale**: When `expectFail()` is called from a Go function assertion method, the Lua call stack has Go frames interspersed. The actual Lua function with local variables may be several levels deeper than expected.

### 4. Value formatting via Lua's tostring

**Decision**: Use Lua's `tostring()` via pcall within the Lua script execution.

**Rationale**: This properly invokes `__tostring` metamethods on userdata objects.

**Critical requirement**: Userdata types that should display meaningful values MUST have a `__tostring` metamethod defined. Without it, `tostring()` returns type+address like `go/std/semver/Range*: 0x...`.

### 5. Expression parsing

**Decision**: Extract identifiers from the expression text using simple tokenization.

**Rationale**: Keep it simple. Parse by splitting on operators and extracting variable names. Filter out Lua keywords and method names.

**Parsing rules**:
- Split on operators: `+`, `-`, `*`, `/`, `==`, `~=`, `<`, `>`, `<=`, `>=`, `and`, `or`, `not`, `:`, `.`
- For `a:b(args)`: `a` is a variable, `b` is a method name
- For `a.b.c`: only `a` (root object) is a variable
- Filter out Lua keywords

### 6. Output format

**Decision**:
```
expected `<expression>` <expected>, actual <actual>
- <var1> = <value1>
- <var2> = <value2>
```

**Rationale**: Clear, readable format. Backticks around expression. Dash-space prefix for variable lines.

## Implementation Summary

**Files modified**:
- `std/mod_testing.go`: Expression capture, parsing, and variable resolution
- `std/mod_semver.go`: Added `__tostring` to `rangeMetatable` for meaningful range display

**Key functions**:
- `extractExpressionVariables(expr string)`: Parses expression and returns variable names
- `collectLocalVariables(l *lua.State, level int)`: Searches levels, runs Lua script, returns map of name->value
- `findLuaFrame(l *lua.State, startLevel int)`: Finds Lua frame by searching stack levels

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| Variable not in scope (optimized away, JIT) | Show `?` for unresolvable variables |
| Closures with upvalues in different frames | Walk full call chain; upvalues may still be resolvable |
| Performance on failure (stack walking, file reading) | Acceptable - only happens on failure path |
| Global variables polluting output | Only use globals as fallback after locals/upvalues |
| Userdata types without `__tostring` | Show type+address; requires adding `__tostring` to types |

## Open Questions

- Should we attempt to resolve function call results (e.g., `foo()` in expression)? **Decision**: No - show `?` for function results, only show simple variable references.
