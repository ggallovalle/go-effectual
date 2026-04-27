## 1. lua-bindgen Parser Updates

- [x] 1.1 Add `//lua: class` annotation parsing in `parser.go`
- [x] 1.2 Add `//lua: module-fn` annotation parsing
- [x] 1.3 Add `//lua: metamethod` annotation parsing
- [x] 1.4 Add `//lua: raw` annotation parsing

## 2. lua-bindgen Generator Updates

- [x] 2.1 Generate `ModQuery` struct with `Name()`, `Annotations()`, `Open()`, `OpenLib()`, `Require()`
- [x] 2.2 Generate `MakeModQuery()` function
- [x] 2.3 Generate `QueryAnnotations()` with template
- [x] 2.4 Generate module-level function wrappers for `//lua: module-fn` annotated functions
- [x] 2.5 Generate verbatim metamethod registration for `//lua: raw` annotated functions
- [x] 2.6 Support `//lua: class` to emit class annotations and proper metatable setup

## 3. luamod.go Updates

- [x] 3.1 Update `LuaModOpen` to accept `LuaModDefinition` (non-generic)
- [x] 3.2 Add `LuaModOpenWithApi[T]` variant that returns `T`
- [x] 3.3 Make sure what previously used the generic version to use the non generic one where it makes sens

## 4. query.go Updates

- [x] 4.1 Add `//lua: class Query` annotation
- [x] 4.2 Add `//lua: module-fn new` to `NewQuery()`
- [x] 4.3 Add `//lua: module-fn deserialize` to new `Deserialize()` function
- [x] 4.4 Add `//lua: module-fn serialize` to new `Serialize()` function
- [x] 4.5 Add `//lua: metamethod __tostring` + `//lua: raw` with `QueryToString()` implementation
- [x] 4.6 Add `//lua: metamethod __pairs` + `//lua: raw` with `QueryPairs()` implementation

## 5. Integration

- [x] 5.1 Delete `mod_query.go`
- [x] 5.2 Run lua-bindgen to generate new `query_bindings.go`
- [x] 5.3 Verify generated code compiles
- [x] 5.4 Run existing tests to verify no changes needed
