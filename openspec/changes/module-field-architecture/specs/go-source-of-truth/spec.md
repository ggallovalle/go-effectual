## ADDED Requirements

### Requirement: Preference for generated bindings
When a module can be fully annotated, the generator SHALL be used. Manual `mod_*.go` files SHALL NOT be created for new modules or when refactoring existing modules.

#### Scenario: New module uses generator
- **WHEN** a new Lua module is being created
- **THEN** it SHALL be created with proper annotations and generated via lua-bindgen

#### Scenario: Existing module refactored
- **WHEN** a module with `mod_*.go` is being refactored
- **THEN** the refactoring SHALL move to use annotations and generated code, removing the manual file

### Requirement: Generated code replaces manual code
For annotated modules, the generated `*_bindings.go` file SHALL contain all code previously in `mod_*.go`. The generated code SHALL produce identical behavior to the manual code it replaces.

#### Scenario: Path module behavior preserved
- **WHEN** `mod_path.go` is removed and replaced by generated code
- **THEN** all tests in `mod_path_test.go` SHALL pass without modification

#### Scenario: Url module behavior preserved
- **WHEN** `mod_url.go` is removed and replaced by generated code
- **THEN** all tests in `mod_url_test.go` SHALL pass without modification

### Requirement: Single source of truth
For annotated modules, the Go source file with annotations SHALL be the only source of truth. The generated binding file is derived and SHALL be regenerated when the source changes.

#### Scenario: Regeneration after source edit
- **WHEN** annotations in `path.go` are modified
- **THEN** running `lua-bindgen generate` regenerates `path_bindings.go` to reflect the changes

### Requirement: Module type as canonical state holder
The module type (`//lua: class`) SHALL be the canonical holder of module-level state. State SHALL NOT be duplicated in global variables or separate structs.

### Requirement: Platform detection in module type
Platform-specific initialization (such as `runtime.GOOS` detection for path separators) SHALL happen in the module type's constructor or initialization, not in separate global state.

### Requirement: Test compatibility
Generated code for annotated modules SHALL produce behavior compatible with existing tests without modifying the tests.

#### Scenario: Tests unchanged
- **WHEN** a module is refactored from manual `mod_*.go` to generated code
- **THEN** `go test ./std/...` passes for all affected tests

### Requirement: Lua-bindgen as build dependency
For annotated modules, `lua-bindgen` SHALL be run as part of the build process. The `lua-bindgen.sh` script SHALL generate bindings for all annotated modules.

#### Scenario: All modules in lua-bindgen.sh
- **WHEN** a new annotated module is added
- **THEN** its generation command SHALL be added to `lua-bindgen.sh`

### Requirement: Full annotation generation
The generator SHALL produce complete EmmyLua annotation files including:
- `@meta` and `@module` markers
- `@class` declarations with inheritance
- `@field` declarations
- `@param` and `@return` annotations
- `@operator` declarations
- `@enum` and `@alias` definitions where specified
- `@source` pointing to Go source location (auto-generated)
- `@return_overload` for error-returning functions

#### Scenario: Annotations match manual files
- **WHEN** the generator produces annotations for a module
- **THEN** the output SHALL match the format and content of manually-written EmmyLua files in `luahome/definitions/`

### Requirement: Annotations template using Go template literal
The generated annotations template SHALL use Go backtick string syntax, not string concatenation.

#### Scenario: Template format
- **WHEN** the generator creates the annotations template
- **THEN** the template uses ```...``` syntax for template literals
