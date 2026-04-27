## 1. Parser Updates

- [ ] 1.1 Add `Module`, `Class`, `Method`, `Name` to annotation types in `internal/luagen/types.go`
- [ ] 1.2 Parse `//lua: module <path>` annotation in `internal/luagen/parser.go`
- [ ] 1.3 Parse `//lua: class <TypeName>` annotation in `internal/luagen/parser.go`
- [ ] 1.4 Parse `//lua: method` (forces no-arg method to be method, overriding default field inference) and `//lua: name <name>` annotations in `internal/luagen/parser.go`
- [ ] 1.5 Implement inference: no params → field, has params → method

## 2. Signature Inference

- [ ] 2.1 Add function to detect `(*lua.State) int` signature (go-lua function marker)
- [ ] 2.2 Add function to classify return types for wrapping
- [ ] 2.3 Add logic to emit `@return_overload` for `(value, error)` returns

## 3. Generator Updates

- [ ] 3.1 Generate module type from `//lua: module <path>` declaration — struct with fields (sep, altSep), Name(), Annotations(), Open() method
- [ ] 3.2 Generate class type from `//lua: class <TypeName>` declaration — userdata with FromRaw, Get, Set methods and metatable registration
- [ ] 3.3 Generate `Make<ModuleType>()` factory function — calls nativeSep(), returns module instance
- [ ] 3.4 Generate wrapper for inferred fields (no params, non-lua.State return) — lua wrapper pushes name+value via SetTable
- [ ] 3.5 Generate verbatim call for go-lua signatures — SetupPosix/SetupWin32 called directly in Open(), not wrapped
- [ ] 3.6 Generate field setup in `Open()` method — iterate declared methods in order, call verbatim for go-lua sigs, set fields via SetTable
- [ ] 3.7 Generate class methods with metamethod registration — lua.SetFunctions with method map, register `__index` metamethod
- [ ] 3.8 Include fields and class methods in annotations template — emit `---@field` entries and `---@param`/`---@return` for methods

## 4. EmmyLua Annotation Generation

- [ ] 4.1 Add `@return_overload` emission for Go error-returning functions
- [ ] 4.2 Add auto-generated `@source` annotations pointing to Go source
- [ ] 4.3 Parse `//lua: raises`, `//emmylua: doc`, `//emmylua: deprecated`, etc.
- [ ] 4.4 Generate complete EmmyLua annotation files (`.lua` files)
- [ ] 4.5 Include `@enum` and `@alias` generation when specified
- [ ] 4.6 Generate `@meta` and `@module` markers

## 5. Path Module Refactoring

- [ ] 5.1 Add `//lua: module std.path` (PathModule) and `//lua: class Path` (userdata) to `path.go`
- [ ] 5.2 Add `PathModule` struct with sep/altSep fields
- [ ] 5.3 Add constructor `NewPathModule()` with platform detection
- [ ] 5.4 Add `MainSeparator()` method (inferred as field - no params)
- [ ] 5.5 Add `SetupPosix(l *lua.State)` method (inferred as field - go-lua signature)
- [ ] 5.6 Add `SetupWin32(l *lua.State)` method (inferred as field - go-lua signature)
- [ ] 5.7 Add module methods `//lua: method new`, `join`, `absolute`
- [ ] 5.8 Add `//lua: class Path` with userdata methods (IsAbsolute, Clean, etc.)
- [ ] 5.9 Add `//emmylua: doc` annotations for path methods
- [ ] 5.10 Regenerate bindings and verify tests pass
- [ ] 5.11 Remove `mod_path.go`

## 6. URL Module Refactoring

- [ ] 6.1 Refactor `url.go` to use module-type architecture
- [ ] 6.2 Add `UrlModule` type with state
- [ ] 6.3 Move module functions as methods on UrlModule
- [ ] 6.4 Add field and method annotations
- [ ] 6.5 Add EmmyLua documentation annotations
- [ ] 6.6 Regenerate bindings and verify tests pass
- [ ] 6.7 Remove `mod_url.go`

## 7. Query Module (verify no changes needed)

- [ ] 7.1 Verify `query_bindings.go` already supports field inference pattern
- [ ] 7.2 Add `@return_overload` annotations to query module functions
- [ ] 7.3 Verify tests pass with current implementation

## 8. Build Integration

- [ ] 8.1 Update `lua-bindgen.sh` to generate both Go bindings and EmmyLua files
- [ ] 8.2 Update `effectual lua-defs` to use generated annotations instead of manual files
- [ ] 8.3 Verify all tests pass: `go test ./...`
- [ ] 8.4 Run `go fmt` and `go vet`
