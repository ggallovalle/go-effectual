## Context

Test context (`testCtx`) in `std/mod_testing.go` provides `skip` and `log` methods via userdata metatable. Tests currently use raw Lua `assert()` for assertions, which gives poor error messages and inconsistent patterns. The `semver_test.lua` suite demonstrates 10 distinct assertion patterns that need fluent API support.

## Goals / Non-Goals

**Goals:**
- Add `expect(value, msg?)` method to ctx that returns expect object
- Direct assertion methods: is_nil, not_nil, is_true, is_false, equals, not_equals, is_lt, not_lt, is_le, not_le
- Implement all assertions used by semver_test.lua
- Maintain backward compatibility with existing assert usage

**Non-Goals:**
- No matcher combinators (contains, matches, etc.) beyond semver needs
- No async/promise support
- No custom failure message formatting beyond the optional msg parameter
- No AST parsing for expression extraction — simple bracket-depth parser only

## Decisions

**Expect as Lua table with metamethod, not userdata.**
Lua table with `__index` metamethod is simpler than userdata for this case. Methods are shared functions, state (value, msg) stored in table fields. Avoids per-call closure allocation.

**Direct assertion methods, no `.not` modifier.**
`expect(v):is_nil()` and `expect(v):not_nil()` — separate methods for positive and negative assertions. Simpler API, no negation state to track.

**Shared method functions, not closures.**
All assertion methods defined once as `lua.Function` values. Each receives expect table as `self` (arg 1), reads value/msg from it. Zero per-call allocation beyond the expect table itself.

**Failure via `t.FailNow()` through Lua error.**
On assertion failure, push error string to stack and call `l.Error()`. The test runner in `lua_suite_test.go` already handles Lua errors as test failures.

**Source location via `lua.Where`.**
Each assertion failure calls `lua.Where(l, 1)` to get the caller's file:line. This is prepended to the error message. Same pattern already used by `ctx.log` in `mod_testing.go`.

**Expression extraction from source file.**
When no custom msg is provided, the failure handler reads the source file at the call site line and parses the first argument passed to `expect()`. A simple bracket-depth parser extracts the expression between `expect(` and its matching `)`. If source is unavailable (inline string chunks, missing files), expression is omitted from the message. Custom msg always takes priority over expression extraction.

## Risks / Trade-offs

**[Table-based expect vs userdata]** Table is simpler but slightly slower on field access. Negligible for test code. Mitigation: acceptable trade-off for maintainability.

**[No deep equality]** `equals` uses Lua `==` operator. Tables compared by reference, not content. Mitigation: sufficient for semver_test.lua; add deep equal later if needed.

**[Error message quality]** Failure messages use `expected <>, actual <>` format plus optional msg and source location. Clear, consistent, machine-parseable.
