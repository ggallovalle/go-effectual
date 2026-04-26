## MODIFIED Requirements

### Requirement: failure messages include context, source location, expression, and variable values
When an assertion fails, the error message SHALL include the Lua source file and line number where the assertion was called, the optional msg if provided, and the values formatted as `expected <expected_value>, actual <actual_value>`. The source location MUST be obtained via `lua.Where` or equivalent call stack inspection at the point of failure.

When no custom msg is provided, the error message SHALL include:
1. The Lua expression passed as the first argument to `expect()`, extracted by reading the source file at the call site
2. The values of variables referenced in the expression, resolved from the call stack at failure time

The expression format SHALL be `` expected `<expression>` <expected_value>, actual <actual_value> ``.

When variable values cannot be resolved from the call stack, they SHALL be displayed as `?`.

When a custom msg is provided, msg takes priority over expression extraction and variable resolution.

#### Scenario: failure with custom message shows source line and expected/actual
- **WHEN** caller invokes `expect(nil, "should exist"):not_nil()` from line 10 of `test.lua`
- **THEN** failure message contains "test.lua:10", "should exist", and "expected non-nil, actual nil"

#### Scenario: failure without custom message shows expression and variable values
- **WHEN** caller invokes `expect(v.major):equals(10)` from line 25 of `test.lua` where v.major is 1 and v is a table
- **THEN** failure message contains "test.lua:25", `` `v.major` ``, and variable values: `- v = table: 0x...`

#### Scenario: expression extraction falls back when source unavailable
- **WHEN** caller invokes `expect(nil):not_nil()` from an inline string chunk (e.g. `loadstring`)
- **THEN** failure message contains "expected non-nil, actual nil" without expression
