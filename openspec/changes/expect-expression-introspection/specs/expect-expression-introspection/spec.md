## ADDED Requirements

### Requirement: failure messages show resolved variable values
When an assertion fails and no custom message was provided, the failure message SHALL include the values of variables extracted from the expression.

The failure message format SHALL be:
```
expected `<expression>` <expected>, actual <actual>
- <var1> = <value1>
- <var2> = <value2>
```

#### Scenario: method call expression shows receiver and argument values
- **WHEN** caller invokes `ctx:expect(r:contains(v1)):is_false()` where r is a range and v1 is a version
- **THEN** failure message shows:
  - Expression in backticks: `r:contains(v1)`
  - Variable values: `- r = <tostring(r)>`, `- v1 = <tostring(v1)>`

#### Scenario: binary expression shows all operand values
- **WHEN** caller invokes `ctx:expect(a + b):equals(10)` where a is 3 and b is 4
- **THEN** failure message shows:
  - Expression in backticks: `a + b`
  - Variable values: `- a = 3`, `- b = 4`

#### Scenario: unresolved variable shows question mark
- **WHEN** caller invokes `ctx:expect(foo):equals(bar)` where foo cannot be resolved from the call stack
- **THEN** failure message shows:
  - `- foo = ?`

### Requirement: expression location captured at expect() call time
The source location of the expression SHALL be captured when `expect()` is called, not when the assertion fails.

#### Scenario: location captured at expect time
- **WHEN** `ctx:expect(someValue)` is called from line 42 of `test.lua`
- **THEN** the location is stored and used later if assertion fails

### Requirement: variable resolution from call stack
Variables in the expression SHALL be resolved by inspecting the Lua call stack at failure time.

#### Scenario: local variables are resolved
- **WHEN** assertion fails inside a function with local variable `x`
- **THEN** `x` is resolved to its current value

#### Scenario: upvalue variables are resolved
- **WHEN** assertion fails inside a closure that captures outer variable `y` as upvalue
- **THEN** `y` is resolved to its current value

#### Scenario: global variables are resolved as fallback
- **WHEN** a variable is not found in locals or upvalues but exists in the global table
- **THEN** the global value is used

### Requirement: tostring uses Lua __tostring metamethod
Variable values SHALL be formatted using Lua's tostring semantics, including calling `__tostring` metamethods if defined.

#### Scenario: objects with __tostring show custom format
- **WHEN** a semver object with `__tostring` metamethod is resolved
- **THEN** the formatted string (e.g., "1.5.0") is shown, not the default table representation

#### Scenario: userdata without __tostring shows type and address
- **WHEN** a userdata type does not define `__tostring` metamethod
- **THEN** tostring returns the default format: `<type>: 0x<address>`
- **EXAMPLE**: `go/std/semver/Range*: 0x2a356c03a20` for Range objects without `__tostring`
