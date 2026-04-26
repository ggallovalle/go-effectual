#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
BINDGEN="$SCRIPT_DIR/lua-bindgen"

if [ ! -f "$BINDGEN" ]; then
    echo "Building lua-bindgen..."
    go build -o "$BINDGEN" ./cmd/lua-bindgen/
fi

echo "Generating Path bindings..."
"$BINDGEN" generate std/path.go \
    --type Path \
    --module std.path \
    --skip StripPrefix,EndsWith,StartsWith \
    --force-method Pop,ToString \
    --nil-map FileName,Extension,FileStem

echo "Generating Query bindings..."
"$BINDGEN" generate std/serde/query.go \
    --type Query

echo "Generating Url bindings..."
"$BINDGEN" generate std/url.go \
    --type Url \
    --module std.url \
    --skip-fields raw,portInferred \
    --force-method ToString \
    --import serde=github.com/ggallovalle/go-effectual/std/serde

echo "Done."
