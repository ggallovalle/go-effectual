## 1. lua-bindgen Parser Updates

- [ ] 1.1 Add `//lua: class` annotation parsing in `parser.go`
- [ ] 1.2 Add `//lua: module-fn` annotation parsing
- [ ] 1.3 Add `//lua: metamethod` annotation parsing
- [ ] 1.4 Add `//lua: raw` annotation parsing

## 2. lua-bindgen Generator Updates

- [ ] 2.1 Generate `ModQuery` struct with `Name()`, `Annotations()`, `Open()`, `OpenLib()`, `Require()`
- [ ] 2.2 Generate `MakeModQuery()` function
- [ ] 2.3 Generate `QueryAnnotations()` with template
- [ ] 2.4 Generate module-level function wrappers for `//lua: module-fn` annotated functions
- [ ] 2.5 Generate verbatim metamethod registration for `//lua: raw` annotated functions
- [ ] 2.6 Support `//lua: class` to emit class annotations and proper metatable setup

## 3. luamod.go Updates

- [ ] 3.1 Update `LuaModOpen` to accept `LuaModDefinition` (non-generic)
- [ ] 3.2 Add `LuaModOpenWithApi[T]` variant that returns `T`
- [ ] 3.3 Make sure what previously used the generic version to use the non generic one where it makes sens

## 4. query.go Updates

- [ ] 4.1 Add `//lua: class Query` annotation
- [ ] 4.2 Add `//lua: module-fn new` to `NewQuery()`
- [ ] 4.3 Add `//lua: module-fn deserialize` to new `Deserialize()` function
- [ ] 4.4 Add `//lua: module-fn serialize` to new `Serialize()` function
- [ ] 4.5 Add `//lua: metamethod __tostring` + `//lua: raw` with `QueryToString()` implementation
- [ ] 4.6 Add `//lua: metamethod __pairs` + `//lua: raw` with `QueryPairs()` implementation

## 5. Integration

- [ ] 5.1 Delete `mod_query.go`
- [ ] 5.2 Run lua-bindgen to generate new `query_bindings.go`
- [ ] 5.3 Verify generated code compiles
- [ ] 5.4 Run existing tests to verify no changes needed
