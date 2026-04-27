## Why

The `expect` testing API in `std.testing` has grown to 10 methods, but only 6 are used in practice. The LT/LE comparison methods (`is_lt`, `not_lt`, `is_le`, `not_le`) are unused across all test files and add unnecessary surface area to the API.

## What Changes

- **Remove** `is_lt`, `not_lt`, `is_le`, `not_le` from `expectMethods` map in `std/mod_testing.go`
- **Keep** `is_nil`, `not_nil`, `is_true`, `is_false`, `equals`, `not_equals`

## Capabilities

### New Capabilities
- `expect-api`: Defines the reduced 6-method expect API for assertions

### Modified Capabilities
- (none)

## Impact

- **Breaking**: Any existing Lua tests using LT/LE methods will fail at runtime. Currently zero usages exist in the codebase.
- **File**: `std/mod_testing.go` - remove 4 methods from `expectMethods` map (lines ~628-683)
