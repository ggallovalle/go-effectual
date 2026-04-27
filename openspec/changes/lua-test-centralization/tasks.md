## 1. Core Infrastructure

- [ ] 1.1 Define `LuaTestCtxExtension` interface in `std/lua_suite_test.go`
- [ ] 1.2 Refactor `runLuaSuite` signature to `runLuaSuite(t *testing.T, l *lua.State, path string, extensions ...LuaTestCtxExtension)`
- [ ] 1.3 Implement dependency resolution: merge suite.deps + case.deps (case overrides suite)
- [ ] 1.4 Implement extension lookup and Build() call per case
- [ ] 1.5 Add error for unknown dependency with file:line and available extensions
- [ ] 1.6 Ensure fresh ctx.ext per test case

## 2. Logger Extension Implementation

- [ ] 2.1 Create `LuaTestLoggerExtension` struct implementing `LuaTestCtxExtension`
- [ ] 2.2 `Name()` returns `"logger"`
- [ ] 2.3 `Build()` pushes slog.Logger userdata with metatable to stack

## 3. Create Lua Test Files

- [ ] 3.1 Create `luahome/std-test/semver_test.lua` (already exists - verify format)
- [ ] 3.2 Create `luahome/std-test/path_test.lua` from 603 lines of inline Lua in `mod_path_test.go`
- [ ] 3.3 Create `luahome/std-test/url_test.lua` from inline Lua in `mod_url_test.go`
- [ ] 3.4 Create `luahome/std-test/slog_test.lua` from inline Lua in `mod_slog_test.go`
- [ ] 3.5 Create `luahome/std-test/serde/query_test.lua` from inline Lua in `mod_query_test.go`

## 4. Create Single Go Wrapper File

- [ ] 4.1 Create `std/lua_test.go` with `TestLuaSuite_Semver`, `TestLuaSuite_Path`, `TestLuaSuite_Url`, `TestLuaSuite_Slog`, `TestLuaSuite_Query` functions
- [ ] 4.2 Each function sets up required modules and calls `runLuaSuite` with appropriate extensions
- [ ] 4.3 Logger tests pass `&LuaTestLoggerExtension{}` to `runLuaSuite`
- [ ] 4.4 Remove old `mod_path_test.go`, `mod_url_test.go`, `mod_slog_test.go`, `mod_query_test.go`

## 5. Verify mod_testing_test.go

- [ ] 5.1 Confirm `mod_testing_test.go` stays as-is (bootstrapping concerns)
- [ ] 5.2 No changes needed - this tests the test framework itself
