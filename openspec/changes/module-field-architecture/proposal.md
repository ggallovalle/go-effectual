## Why

The current architecture has manual `mod_*.go` files (mod_path.go, mod_url.go) that contain module registration logic separate from the source types. This creates maintenance burden and drift between generated bindings and module code. The generator should be the source of truth - all module code generated from annotations.

## What Changes

- **Two-annotation system**:
  - `//lua:` annotations control Go code generation (bindings, wrappers, module setup)
  - `//emmylua:` annotations provide supplemental EmmyLua documentation (most inferred from Go)
- **New field inference**: Methods with no params inferred as fields (no annotation needed). Methods with params inferred as methods.
- **New `//lua: class <Type>` annotation** - Declares the module state type that holds module-level data
- **Smart signature inference** - Generator determines whether to wrap a function or call verbatim:
  - `(*lua.State) int` → verbatim go-lua function call
  - Other return types → generate wrapper that pushes value to stack
- **Module type replaces module struct** - User defines a Go struct with `//lua: class` whose methods become module functions/fields
- **`@return_overload` for error handling** - Go's `(value, error)` pattern becomes first-class in Lua via EmmyLua's correlated return rows
- **`@source` auto-generation** - The generator automatically includes source location in generated annotations
- **Path and URL refactoring** - These modules refactored to use new module-type architecture
- **BREAKING**: `mod_path.go` and `mod_url.go` removed, replaced by generated code

## Capabilities

### New Capabilities

- **lua-bindgen**: Core generator capabilities including:
  - Field inference for declaring module fields
  - `//lua: class <Type>` annotation for module state types
  - Signature-based wrapping inference (return types inferred from Go)
  - `//lua: metamethod` for metamethod registration
  - `//lua: nil-map` and `//lua: method` for method behavior

- **go-source-of-truth**: The principle that all Lua bindings SHALL be generated from Go source annotations. Manual `mod_*.go` files SHALL NOT exist for modules that can be fully annotated.

- **EmmyLua annotation generation**: Full EmmyLua annotation file generation including:
  - `@return_overload` for Go error-returning functions (enables IDE type narrowing)
  - Auto-generated `@source` pointing to Go source
  - `@deprecated`, `@see`, `@version` via explicit `//emmylua:` annotations
  - `@raises <ErrorType>` custom annotation for error documentation

### Modified Capabilities

- **query-module-source-of-truth**: Already exists. Field inference should apply to Query module generation as well, replacing any remaining manual code in the query bindings.

## Impact

- `cmd/lua-bindgen`: New parsing and generation logic for `module`, `class`, `method`, `name` annotations and field inference
- `internal/luagen`: Generator updated to handle new annotation types, signature inference, and `@return_overload` emission
- `std/mod_path.go`: Removed, replaced by generated code from path.go annotations
- `std/mod_url.go`: Removed, replaced by generated code from url.go annotations
- `std/path_bindings.go`, `std/url_bindings.go`: May need regeneration with new annotations
- `luahome/definitions/`: Generated EmmyLua files produced by generator instead of manually maintained
- Tests in `std/` pass without modification (generated code produces same behavior)
