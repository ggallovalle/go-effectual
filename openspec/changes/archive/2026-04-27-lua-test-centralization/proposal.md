## Why

Lua tests for `std` modules are currently split: some live in `luahome/std-test/` (like `semver_test.lua`), while others are inline in Go test files (like `mod_path_test.go` with 600+ lines of inline Lua). This inconsistency makes tests harder to maintain and reuse. We want `luahome/std-test/` as the single source of truth for all Lua tests.

Additionally, the current extension mechanism (passing a callback to `runLuaSuite`) is awkward when different test cases need different Go-side state (e.g., a logger). We need a cleaner way to declare test dependencies that can be resolved at runtime.

## What Changes

- Consolidate all Lua tests into `luahome/std-test/` directory
- Replace `mod_path_test.go`, `mod_url_test.go`, `mod_slog_test.go`, `mod_query_test.go` with a single `lua_test.go` containing thin wrapper functions that call `runLuaSuite`
- Keep `mod_testing_test.go` as-is (tests the test framework itself, bootstrapping concerns)
- Refactor `runLuaSuite(t, l, path, extensions...)` to accept variadic `LuaTestCtxExtension` instead of a callback
- Support `suite.deps` and `case.deps` tables in Lua test files for declaring extensions
- Case deps override suite deps for same dependency name
- Error at runtime if a declared dependency is not available, with line number and list of available extensions
- Each test case gets a fresh `ctx` with `ctx.ext` table populated by resolved extensions

## Capabilities

### New Capabilities

- `lua-test-suite`: Test suite format and runner for Lua-based tests in `luahome/std-test/`. Defines the `Suite` table structure with `name`, `deps`, and `cases` fields. Runner resolves declared dependencies against registered extensions and executes each case with a fresh context.
- `lua-test-ctx-extension`: Mechanism for injecting Go-side state into test context via the `LuaTestCtxExtension` interface. Extensions are registered at `runLuaSuite` call time. Each extension implements `Name()` and `Build(l, params)` to inject state into `ctx.ext`.

### Modified Capabilities

- No existing capabilities are modified. This is purely a reorganization and refactoring.

## Impact

- **Code movement**: Tests from `std/mod_path_test.go`, `std/mod_slog_test.go`, `std/mod_url_test.go`, `std/serde/mod_query_test.go` move to `luahome/std-test/`
- **Go test files consolidated**: Four Go test files replaced by single `lua_test.go` with thin wrapper functions
- **New error path**: Unknown dependencies now produce runtime errors with file:line and available extensions list
- **No API changes**: External consumers of the test framework are unaffected (tests are internal)
