## 1. Parser Annotation Extraction

- [x] 1.1 Add `GenConfigAnnotation` struct in types.go with annotation-sourced fields
- [x] 1.2 Add `extractTypeAnnotations()` function in parser.go to parse `@lua-bindgen.sh` block comments
- [x] 1.3 Handle build tag detection (`//go:build lua_bindgen` or `// +lua-bindgen.sh`)
- [x] 1.4 Parse key=value pairs with comma-separated list support
- [x] 1.5 Return partial GenConfig from annotations for later merge

## 2. Main Integration

- [x] 2.1 Call `extractTypeAnnotations()` in `ParseSource()` flow
- [x] 2.2 Merge annotation config with CLI flags in `runGenerate()` (CLI takes precedence)
- [x] 2.3 Add warning for unknown annotation keys

## 3. Testing

- [x] 3.1 Add test for annotation parsing in parser_test.go
- [x] 3.2 Test CLI precedence (CLI overrides annotation)
- [x] 3.3 Test build tag detection variants
- [x] 3.4 Test multi-value comma-separated parsing

## 4. Apply Query Bindings

- [x] 4.1 Simplify lua-bindgen.sh invocation for std/serde/query.go to just --type Query
- [x] 4.2 Add inline lua annotations to std/serde/query.go (module, skip-field, nil-map, force-method)
