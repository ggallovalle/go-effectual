## 1. Data Structures

- [x] 1.1 Add `exprLocation` field to `expectObj` struct in `mod_testing.go`
- [x] 1.2 Add `exprString` field to `expectObj` struct (may store parsed expression)

## 2. Capture Expression Location at expect() Time

- [x] 2.1 In `expect()` method, call `lua.Where(l, 1)` to get caller's location
- [x] 2.2 Store the location string in `expectObj.exprLocation`
- [x] 2.3 Remove `lua.Where` call from `expectFail()` - use stored location instead

## 3. Variable Resolution Functions

- [x] 3.1 Implement `findLuaFrame(l *lua.State) *callInfo` - walk `.previous` chain until `isLua()` returns true
- [x] 3.2 Implement `collectLocalVariables(l *lua.State, ci *callInfo) map[string]string` - enumerate locals via `getLocal()`
- [x] 3.3 Implement `collectUpvalues(l *lua.State, ci *callInfo) map[string]string` - enumerate upvalues via `UpValue()`
- [x] 3.4 Implement `collectGlobals(l *lua.State, exclude map[string]bool) map[string]string` - check `_G` for variables not in exclude set
- [x] 3.5 Implement `resolveVariableValue(l *lua.State, v value) string` - push value, use `ToStringMeta()`, pop, return string

## 4. Expression Parsing

- [x] 4.1 Create `extractExpressionVariables(expr string) []string` - parse expression and return variable names
- [x] 4.2 Filter out Lua keywords (`and`, `or`, `not`, etc.)
- [x] 4.3 Filter out method names (after `:` or `.` in expressions like `a:b()` or `a.b`)
- [x] 4.4 Handle binary operators (`+`, `-`, `*`, `/`, `==`, `~=`, `<`, `>`, `<=`, `>=`)

## 5. Format Failure Message with Variables

- [x] 5.1 Modify `expectFail()` to call new variable resolution when no custom msg
- [x] 5.2 Format output: `` expected `<expression>` <expected>, actual <actual> ``
- [x] 5.3 Append variable values: `` - <name> = <value> `` lines
- [x] 5.4 Use `?` for unresolvable variables

## 6. Tests

- [x] 6.1 Add test for `extractExpressionVariables` with various expression patterns
- [x] 6.2 Add test for variable resolution with local variables
- [x] 6.3 Add test for variable resolution with upvalues (closure)
- [x] 6.4 Add test for failure message format with variables
- [x] 6.5 Add test that custom message takes priority (no variable output)