## Why

Currently, generator options like `--skip-fields`, `--nil-map`, `--force-method`, `--skip`, and `--module` must be passed as CLI flags. This makes the generator harder to use when generating bindings for multiple types, and the configuration drifts from the source code it relates to. Annotations allow configuration to live alongside the source, improving discoverability and reducing shell script complexity.

## What Changes

- Add support for inline `//lua:` annotations on struct fields and methods in Go source files
- Support annotation-based configuration for: `skip-field`, `nil-map`, `force-method`, `skip`, `module`
- CLI flags take precedence over annotations (allows overriding without modifying source)
- Shell script invocations simplify from multi-flag to minimal flags
- Module can be set via `// lua:module <name>` before the type declaration

## Capabilities

### New Capabilities

- `lua-bindgen-annotations`: Parse inline `//lua:` annotations on struct fields and methods to extract generator configuration. CLI flags override annotation values when both are present.

### Modified Capabilities

- None (new capability only)

## Impact

- `internal/luagen/parser.go`: Enhanced `extractMethodComments()` and `extractStructFields()` to handle inline annotations (`lua:skip`, `lua:nil-map`, `lua:force-method`, `lua:skip-field`)
- `internal/luagen/types.go`: Added `IsForceMethod` field to `MethodInfo`
- `internal/luagen/classifier.go`: Check `m.IsForceMethod` when classifying getters
- `internal/luagen/generator.go`: Check `f.IsSkipped` when generating field getters
- `cmd/lua-bindgen/main.go`: Merge annotation config with CLI flags, CLI takes precedence
- Shell scripts using lua-bindgen can use fewer CLI flags