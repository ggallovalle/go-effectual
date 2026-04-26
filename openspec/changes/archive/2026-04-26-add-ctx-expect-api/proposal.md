## Why

Test suites use raw `assert()` calls with manual error messages. This produces inconsistent assertion patterns, verbose code, and poor failure messages. `ctx.expect()` provides a fluent, chainable assertion API that improves readability and gives structured failure output.

## What Changes

- Add `ctx.expect(value, msg?)` method to test context (`std.testing`)
- Return expect object with direct assertion methods
- Implement assertions needed by `semver_test.lua`: `is_nil`, `not_nil`, `is_true`, `is_false`, `equals`, `not_equals`, `is_lt`, `not_lt`, `is_le`, `not_le`
- Update `semver_test.lua` to use `ctx.expect` instead of `assert`

## Capabilities

### New Capabilities
- `ctx-expect-api`: Direct assertion API on test context with methods like is_nil, not_nil, equals, is_lt, etc.

### Modified Capabilities
<!-- None -->

## Impact

- `std/mod_testing.go`: Add expect method and expect object type to testCtx
- `luahome/std-test/semver_test.lua`: Migrate all assert calls to ctx.expect
- No breaking changes to existing test infrastructure
