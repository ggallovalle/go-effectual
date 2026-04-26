## Why

Currently, `mod_query.go` contains manual binding code AND a hardcoded annotation template, while `query.go` also contains annotation markers (`//lua:nil-map`, `//lua:force-method`, `//lua:skip-field`). This duplication creates drift between what `query.go` declares and what `mod_query.go` generates. The goal is to make `query.go` the single source of truth for Query bindings, eliminating `mod_query.go` entirely.

## What Changes

- `std/serde/query.go` becomes the sole source of truth for Query bindings
- New lua-bindgen annotations: `//lua: class`, `//lua: module-fn`, `//lua: metamethod`, `//lua: raw`
- `lua-bindgen` learns to generate:
  - Module-level functions (`new`, `deserialize`, `serialize`) from annotated package-level functions
  - Verbatim metamethods (`__tostring`, `__pairs`) from annotated raw implementations
  - Full `ModQuery` struct implementing `LuaModDefinition`
  - `MakeModQuery()` factory function
  - `QueryAnnotations()` with generated template
  - `QueryLibrary()` with proper registration
- `mod_query.go` is deleted entirely
- `luamod.go` updated:
  - `LuaModOpen(l, mod LuaModDefinition)` - non-generic variant
  - `LuaModOpenWithApi[T](l, mod LuaMod[T]) T` - generic variant for modules that need Api
- Existing tests in `mod_query_test.go` pass unchanged
- `url.go` continues using `NewQuery()`, `FromRaw()` etc. directly (no changes)

## Capabilities

### Modified Capabilities

- `lua-bindgen-annotations`: Add new annotation types for class marking, module functions, and verbatim metamethods

### New Capabilities

- `query-module-source-of-truth`: The Query module bindings are generated entirely from `query.go` annotations, with no separate manual binding file

## Impact

- **Deleted**: `std/serde/mod_query.go`
- **Modified**: `internal/luagen/` (parser, generator), `luamod.go`
- **Generated**: `std/serde/query_bindings.go` updated to include full module registration
- **No changes**: Tests, `url.go`, `std/package.go` (consumes annotations via `LuaModDefinition`)
