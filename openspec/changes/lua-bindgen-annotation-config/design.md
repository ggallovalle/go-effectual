## Context

Generator configuration currently requires CLI flags for all options. This creates a disconnect between the source type and its binding configuration. The goal is to support annotations directly in Go source files while preserving CLI flag override capability.

Current annotation support includes `//lua:skip`, `//lua:nil-map`, `//lua:metamethod`, `//lua:module` on methods. New `@lua-bindgen.sh` annotations will support type-level configuration.

## Goals / Non-Goals

**Goals:**
- Parse `@lua-bindgen.sh` block comments on type declarations
- Support `skip-fields`, `nil-map`, `force-method`, `skip`, `module` via annotations
- CLI flags take precedence over annotations
- Minimal changes to existing generator flow

**Non-Goals:**
- Breaking changes to existing annotation behavior (`//lua:skip`, etc.)
- Support for all CLI flags via annotations (only the common ones listed above)
- Runtime evaluation of annotations

## Decisions

**1. Annotation format: `@lua-bindgen.sh` block comment**

```go
//go:build lua_bindgen
// +lua-bindgen.sh skip-fields=params nil-map=Get force-method=ToString,Keys,Values,Entries module=std.serde.query

type Query struct {
    // ...
}
```

A block comment directly preceding the type declaration is scanned for `@lua-bindgen.sh`. Inside, key=value pairs (comma-separated for lists) configure the generator.

**2. Annotation parsing in `parser.go`**

New function `extractTypeAnnotations(node *ast.File, typeName string) GenConfigAnnotation` extracts annotation values. Returns a struct holding annotation-based config before CLI merge.

**3. CLI takes precedence (merge in `main.go`)**

`runGenerate()` first extracts annotations, then CLI flags override when present. The `GenConfig` struct receives merged values.

**4. No new types needed**

Reuse existing `GenConfig` fields. Annotations populate the same maps CLI flags do.

## Risks / Trade-offs

- [Risk] Go build constraints (`//go:build`) use similar syntax → Mitigation: Require `@lua-bindgen.sh` prefix, not just any block comment
- [Risk] Complex annotation syntax harder to debug → Mitigation: Validation with clear error messages pointing to source location
- [Trade-off] Annotations increase source coupling → Acceptable: CLI override always available

## Open Questions

- Should we require a build tag like `//go:build lua_bindgen` for annotations to be active? (Proposed: yes, prevents annotations from affecting normal builds)
- How to handle malformed annotations? (Proposed: warn and ignore malformed lines, don't error)
