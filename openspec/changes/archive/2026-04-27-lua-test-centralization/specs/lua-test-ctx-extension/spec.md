## ADDED Requirements

### Requirement: LuaTestCtxExtension interface
Go code SHALL implement the `LuaTestCtxExtension` interface with `Name() string` and `Build(l *lua.State, params map[string]any)`.

#### Scenario: Extension provides name
- **WHEN** an extension returns `"logger"` from `Name()`
- **AND** a test case declares `{name = "logger"}`
- **THEN** that extension's `Build` is called

#### Scenario: Extension receives params
- **WHEN** a case declares `{name = "logger", params = {level = "DEBUG"}}`
- **AND** the matching extension's `Build` is called
- **THEN** `params["level"]` equals `"DEBUG"`

### Requirement: Extensions inject into ctx.ext
Extensions SHALL push their value onto the Lua stack in `Build`, which is then stored in `ctx.ext[name]`.

#### Scenario: Logger extension pushes userdata
- **WHEN** LoggerExtension.Build is called
- **AND** it pushes a logger userdata onto the stack
- **THEN** `ctx.ext.logger` equals that userdata

#### Scenario: Build receives Lua state
- **WHEN** `Build(l, params)` is called
- **THEN** `l` is the current Lua state with the ctx table accessible

### Requirement: runLuaSuite accepts extensions
The `runLuaSuite` function SHALL accept zero or more `LuaTestCtxExtension` arguments.

#### Scenario: No extensions passed
- **WHEN** `runLuaSuite(t, l, "test.lua")` is called with no extensions
- **AND** a test case declares a dependency
- **THEN** an error is raised for the unknown dependency

#### Scenario: Extensions available
- **WHEN** `runLuaSuite(t, l, "test.lua", &LoggerExtension{})` is called
- **AND** a test case declares `{name = "logger"}`
- **THEN** the extension is resolved and available

### Requirement: Extension registry via variadic args
Registered extensions form a map keyed by `Name()`. When multiple extensions share the same name, the last one wins.

#### Scenario: Duplicate extension name
- **WHEN** `runLuaSuite(t, l, "test.lua", ext1, ext2)` is called
- **AND** both have `Name() = "logger"`
- **THEN** ext2's `Build` is used for all "logger" dependencies
