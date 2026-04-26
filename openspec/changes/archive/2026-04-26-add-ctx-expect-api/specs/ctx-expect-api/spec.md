## ADDED Requirements

### Requirement: ctx.expect creates expect object
The test context SHALL provide an `expect(value, msg?)` method that returns an expect object. The value is the actual value under test. The optional msg is a string used in failure messages.

#### Scenario: expect with value only
- **WHEN** caller invokes `ctx.expect(someValue)`
- **THEN** expect object is returned with value set and no message

#### Scenario: expect with value and message
- **WHEN** caller invokes `ctx.expect(someValue, "should be valid")`
- **THEN** expect object is returned with value and message stored for failure output

### Requirement: is_nil assertion
The expect object SHALL provide `is_nil()` that asserts the value is nil. Fails if value is not nil.

#### Scenario: is_nil passes on nil
- **WHEN** caller invokes `expect(nil):is_nil()`
- **THEN** assertion passes

#### Scenario: is_nil fails on non-nil
- **WHEN** caller invokes `expect(42):is_nil()`
- **THEN** assertion fails with error message including optional msg

### Requirement: not_nil assertion
The expect object SHALL provide `not_nil()` that asserts the value is not nil. Fails if value is nil.

#### Scenario: not_nil passes on non-nil
- **WHEN** caller invokes `expect(42):not_nil()`
- **THEN** assertion passes

#### Scenario: not_nil fails on nil
- **WHEN** caller invokes `expect(nil):not_nil()`
- **THEN** assertion fails with error message

### Requirement: is_true assertion
The expect object SHALL provide `is_true()` that asserts the value is boolean true. Fails if value is not true.

#### Scenario: is_true passes on true
- **WHEN** caller invokes `expect(true):is_true()`
- **THEN** assertion passes

#### Scenario: is_true fails on false
- **WHEN** caller invokes `expect(false):is_true()`
- **THEN** assertion fails

### Requirement: is_false assertion
The expect object SHALL provide `is_false()` that asserts the value is boolean false. Fails if value is not false.

#### Scenario: is_false passes on false
- **WHEN** caller invokes `expect(false):is_false()`
- **THEN** assertion passes

#### Scenario: is_false fails on true
- **WHEN** caller invokes `expect(true):is_false()`
- **THEN** assertion fails

### Requirement: equals assertion
The expect object SHALL provide `equals(other)` that asserts the value equals the other value using Lua `==`. Fails if values are not equal.

#### Scenario: equals passes on equal values
- **WHEN** caller invokes `expect(1):equals(1)`
- **THEN** assertion passes

#### Scenario: equals fails on unequal values
- **WHEN** caller invokes `expect(1):equals(2)`
- **THEN** assertion fails

### Requirement: not_equals assertion
The expect object SHALL provide `not_equals(other)` that asserts the value does not equal the other value. Fails if values are equal.

#### Scenario: not_equals passes on unequal values
- **WHEN** caller invokes `expect(1):not_equals(2)`
- **THEN** assertion passes

#### Scenario: not_equals fails on equal values
- **WHEN** caller invokes `expect(1):not_equals(1)`
- **THEN** assertion fails

### Requirement: is_lt assertion
The expect object SHALL provide `is_lt(other)` that asserts the value is less than the other value using Lua `<`. Fails if value is not less than other.

#### Scenario: is_lt passes when value is less
- **WHEN** caller invokes `expect(1):is_lt(2)`
- **THEN** assertion passes

#### Scenario: is_lt fails when value is not less
- **WHEN** caller invokes `expect(2):is_lt(1)`
- **THEN** assertion fails

### Requirement: not_lt assertion
The expect object SHALL provide `not_lt(other)` that asserts the value is not less than the other value. Fails if value is less than other.

#### Scenario: not_lt passes when value is not less
- **WHEN** caller invokes `expect(2):not_lt(1)`
- **THEN** assertion passes

#### Scenario: not_lt fails when value is less
- **WHEN** caller invokes `expect(1):not_lt(2)`
- **THEN** assertion fails

### Requirement: is_le assertion
The expect object SHALL provide `is_le(other)` that asserts the value is less than or equal to the other value using Lua `<=`. Fails if value is greater than other.

#### Scenario: is_le passes when value is less or equal
- **WHEN** caller invokes `expect(1):is_le(1)`
- **THEN** assertion passes

#### Scenario: is_le fails when value is greater
- **WHEN** caller invokes `expect(2):is_le(1)`
- **THEN** assertion fails

### Requirement: not_le assertion
The expect object SHALL provide `not_le(other)` that asserts the value is not less than or equal to the other value. Fails if value is less than or equal to other.

#### Scenario: not_le passes when value is greater
- **WHEN** caller invokes `expect(2):not_le(1)`
- **THEN** assertion passes

#### Scenario: not_le fails when value is less or equal
- **WHEN** caller invokes `expect(1):not_le(1)`
- **THEN** assertion fails

### Requirement: failure messages include context, source location, and expression
When an assertion fails, the error message SHALL include the Lua source file and line number where the assertion was called, the optional msg if provided, and the values formatted as `expected <expected_value>, actual <actual_value>`. The source location MUST be obtained via `lua.Where` or equivalent call stack inspection at the point of failure. When no custom msg is provided, the error message SHALL include the Lua expression passed as the first argument to `expect()`, extracted by reading the source file at the call site. The expression format SHALL be `expected (expr: <expression>) <expected_value>, actual <actual_value>`. When a custom msg is provided, msg takes priority over expression extraction.

#### Scenario: failure with custom message shows source line and expected/actual
- **WHEN** caller invokes `expect(nil, "should exist"):not_nil()` from line 10 of `test.lua`
- **THEN** failure message contains "test.lua:10", "should exist", and "expected non-nil, actual nil"

#### Scenario: failure without custom message shows expression and expected/actual
- **WHEN** caller invokes `expect(v.major):equals(10)` from line 25 of `test.lua` where v.major is 1
- **THEN** failure message contains "test.lua:25" and "expected (expr: v.major) 10, actual 1"

#### Scenario: equals failure shows both values
- **WHEN** caller invokes `expect(1):equals(2)` from line 5 of `test.lua`
- **THEN** failure message contains "test.lua:5" and "expected (expr: 1) 2, actual 1"

#### Scenario: expression extraction falls back when source unavailable
- **WHEN** caller invokes `expect(nil):not_nil()` from an inline string chunk (e.g. `loadstring`)
- **THEN** failure message contains "expected non-nil, actual nil" without expression
