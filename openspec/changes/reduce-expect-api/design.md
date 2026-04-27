## Context

The `expect` API in `std/mod_testing.go` currently exposes 10 methods for assertions:
- Value checks: `is_nil`, `not_nil`, `is_true`, `is_false`
- Equality: `equals`, `not_equals`
- Ordering: `is_lt`, `not_lt`, `is_le`, `not_le`

Analysis shows `is_lt`, `not_lt`, `is_le`, `not_le` are unused in all Lua test files.

## Goals / Non-Goals

**Goals:**
- Reduce API surface of `expect` from 10 to 6 methods
- Eliminate dead code

**Non-Goals:**
- Changing behavior of remaining 6 methods
- Adding new assertion capabilities

## Decisions

1. **Remove LT/LE methods entirely**
   - Rationale: Unused code is maintenance burden; removing it simplifies the API
   - Alternative: Deprecate first? Unnecessary given zero usage

2. **Keep all other methods unchanged**
   - `is_nil`, `not_nil`, `is_true`, `is_false`, `equals`, `not_equals` remain

## Risks / Trade-offs

[Risk] Existing Lua code uses LT/LE methods
→ Mitigation: Searched codebase - no usages exist

[Risk] Someone adds LT/LE usage after removal
→ Mitigation: Tests will fail at runtime with clear error; easy to restore if genuinely needed
