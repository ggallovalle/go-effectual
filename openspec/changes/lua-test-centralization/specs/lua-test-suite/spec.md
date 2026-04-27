## ADDED Requirements

### Requirement: Suite table structure
A Lua test suite SHALL be a table returned from the test file with at least `name` (string) and `cases` (table). The `deps` field is optional.

#### Scenario: Minimal valid suite without deps
- **WHEN** a Lua file returns `{name = "std.semver", cases = {}}`
- **THEN** the runner accepts it as a valid suite with no suite-level dependencies

#### Scenario: Suite with suite-level deps
- **WHEN** a Lua file returns a suite with `deps = {{name = "logger", params = {level = "DEBUG"}}}`
- **THEN** all test cases inherit those dependencies

### Requirement: Suite deps apply to all cases
Suite-level dependencies declared in `suite.deps` SHALL be available to every case in the suite unless overridden.

#### Scenario: Case accesses suite-level dep
- **WHEN** suite has `deps = {{name = "logger"}}`
- **AND** a case accesses `ctx.ext.logger`
- **THEN** the logger extension is available in that case

### Requirement: Case-level deps override suite-level deps
When a case declares a dependency with the same name as a suite-level dependency, the case-level declaration SHALL take precedence.

#### Scenario: Case overrides suite dep params
- **WHEN** suite has `deps = {{name = "logger", params = {level = "DEBUG"}}}`
- **AND** a case has `deps = {{name = "logger", params = {level = "WARN"}}}`
- **THEN** that case receives the WARN level logger

### Requirement: Case table structure
Each case in `cases` SHALL be a table with `name` (string) and `fn` (function). The `deps` field is optional.

#### Scenario: Case without deps
- **WHEN** a case is `{name = "test", fn = function(ctx) end}`
- **THEN** it runs without any extensions

#### Scenario: Case with deps
- **WHEN** a case is `{name = "test", deps = {{name = "logger"}}, fn = function(ctx) end}`
- **THEN** the declared dependencies are available

### Requirement: Dependency table structure
A dependency in `deps` SHALL be either:
- A string, interpreted as `{name = <string>, params = nil}`
- A table with `name` (string) and optional `params` (table)

#### Scenario: String shorthand for dependency
- **WHEN** a case has `deps = {"logger"}`
- **THEN** it is interpreted as `deps = {{name = "logger", params = nil}}`

#### Scenario: Table form for dependency
- **WHEN** a case has `deps = {{name = "logger", params = {level = "DEBUG"}}}`
- **THEN** the extension receives name="logger" and params={level="DEBUG"}

#### Scenario: Dependency with params
- **WHEN** a dep is `{name = "logger", params = {level = "WARN"}}`
- **THEN** the extension receives `params = {level = "WARN"}`

#### Scenario: Dependency without params
- **WHEN** a dep is `{name = "logger"}`
- **THEN** the extension receives `params = nil`

### Requirement: Fresh context per case
Each test case SHALL receive a freshly created context. Extensions are resolved and built for each case independently.

#### Scenario: Different cases get different ext tables
- **WHEN** two cases both declare `deps = {{name = "logger"}}`
- **THEN** each case gets its own `ctx.ext` table with a separate logger instance

### Requirement: Error on unknown dependency
When a case declares a dependency that is not registered, the runner SHALL error with the exact file:line and list of available extensions.

#### Scenario: Unknown dep name
- **WHEN** a case has `deps = {{name = "nonexistent"}}`
- **AND** no extension with that name is registered
- **THEN** error message contains "extension 'nonexistent' not found" and "Available:" list

#### Scenario: Error includes line number
- **WHEN** the error occurs at line 42 of the Lua file
- **THEN** error message contains ":42:" pointing to that line
