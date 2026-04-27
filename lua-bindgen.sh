#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

echo "Generating Path bindings..."
go run "$SCRIPT_DIR/cmd/lua-bindgen" generate std/path.go \
    --type Path \
    --module std.path \
    --skip StripPrefix,EndsWith,StartsWith \
    --force-method Pop,ToString \
    --nil-map FileName,Extension,FileStem

echo "Generating Query bindings..."
go run "$SCRIPT_DIR/cmd/lua-bindgen" generate std/serde/query.go \
    --type Query

echo "Generating Url bindings..."
go run "$SCRIPT_DIR/cmd/lua-bindgen" generate std/url.go \
    --type Url \
    --module std.url \
    --skip-fields raw,portInferred \
    --force-method ToString \
    --import serde=github.com/ggallovalle/go-effectual/std/serde

echo "Done."
