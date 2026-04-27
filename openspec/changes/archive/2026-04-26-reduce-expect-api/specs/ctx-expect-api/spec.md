## ADDED Requirements

### Requirement: expect API provides 6 assertion methods

The `expect` object returned by `ctx:expect()` SHALL provide exactly 6 methods for assertions:

#### Scenario: nil checks
- **WHEN** value is `nil`
- **THEN** `is_nil()` passes and `not_nil()` fails

#### Scenario: boolean checks
- **WHEN** value is `true`
- **THEN** `is_true()` passes and `is_false()` fails
- **WHEN** value is `false`
- **THEN** `is_false()` passes and `is_true()` fails

#### Scenario: equals assertion
- **WHEN** actual value equals expected value
- **THEN** `equals(expected)` passes
- **WHEN** actual value does not equal expected value
- **THEN** `equals(expected)` fails

#### Scenario: not_equals assertion
- **WHEN** actual value does not equal expected value
- **THEN** `not_equals(expected)` passes
- **WHEN** actual value equals expected value
- **THEN** `not_equals(expected)` fails

### Requirement: LT/LE methods are not available

The `expect` object SHALL NOT have `is_lt`, `not_lt`, `is_le`, or `not_le` methods.

#### Scenario: Accessing removed method
- **WHEN** code calls `expect:is_lt(value)` or similar removed method
- **THEN** Lua error occurs at runtime
