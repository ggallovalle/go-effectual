# go-effectual

Effect system for [go-lua](https://github.com/Shopify/go-lua), inspired by [effectual](https://github.com/ggallovalle/effectual).

## Status

**Work in progress** — not yet implemented.

## Goal

Port the [effectual](https://github.com/ggallovalle/effectual) Effect system to Go, exposing it as a Lua library for use in go-lua scripts:

```lua
local EffectGo = require("effectual-go")
```

## Motivation

The [effectual](https://github.com/ggallovalle/effectual) Lua library provides:

- Typed effects with contextual services
- Composable async/await via coroutines
- Dependency injection via `Effect.Service`

This library aims to bring the same patterns to go-lua, enabling Go projects that embed the go-lua VM to offer the Effect pattern to their Lua scripts.

## References

- [go-lua](https://github.com/Shopify/go-lua) — Lua VM in pure Go
- [effectual](https://github.com/ggallovalle/effectual) — The original Lua implementation
- [Effect-TS](https://effect.website) — Inspiration for the Effect pattern in TypeScript
