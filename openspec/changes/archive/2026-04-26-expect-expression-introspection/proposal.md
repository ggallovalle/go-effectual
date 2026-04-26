## Why

When test assertions fail, the current error messages show the expression text (e.g., `r:contains(v1)`) but don't show the actual values of the variables in that expression. This makes debugging harder - you can see what was tested but not the values that caused the failure.

## What Changes

- **Enhanced failure messages**: When an assertion fails and no custom message was provided, failure output will include:
  - The expression in backticks: `` `r:contains(v1)` ``
  - Variable values extracted from the expression, formatted with their tostring representation
- **Expression location capture**: The source location is captured at `expect()` call time, not at failure time
- **Variable introspection**: Walk the Lua call stack to resolve local variables, upvalues, and globals referenced in the expression
- **Fallback behavior**: Unresolved variables show as `?`

## Capabilities

### New Capabilities

- `expect-expression-introspection`: Enhanced failure messages that show variable values from the expression

### Modified Capabilities

- `ctx-expect-api`: The existing "failure messages include context, source location, and expression" requirement is modified to include resolved variable values when expression introspection succeeds

## Impact

- **Primary files**: `std/mod_testing.go`, `std/mod_semver.go`
- **Test files**: `std/mod_testing_test.go`, `std/lua_suite_test.go`
- **Implementation approach**: Execute Lua scripts via `lua.LoadString` + `ProtectedCall` to access `debug.getlocal` in correct Lua context
- **Requirement**: Userdata types that should display meaningful values must define `__tostring` metamethod
