package serde_test

import (
	"testing"

	lua "github.com/Shopify/go-lua"
	sut "github.com/ggallovalle/go-effectual/std/serde"
	"github.com/ggallovalle/go-effectual"
	"github.com/stretchr/testify/assert"
)

func TestQueryNit_KeysOrderNonDeterministic(t *testing.T) {
	l := lua.NewState()
	lua.OpenLibraries(l)
	effectual.LuaModOpen(l, sut.MakeModQuery())

	err := lua.DoString(l, `
		local query = require("std.serde.query")
		local q = query.deserialize("z=last&a=first&m=middle")
		local keys1 = q:keys()
		local keys2 = q:keys()
		local s1 = table.concat(keys1, ",")
		local s2 = table.concat(keys2, ",")
		assert(s1 == s2, "keys() order non-deterministic: " .. s1 .. " vs " .. s2)
	`)
	assert.NoError(t, err)
}

func TestQueryNit_ValuesOrderNonDeterministic(t *testing.T) {
	l := lua.NewState()
	lua.OpenLibraries(l)
	effectual.LuaModOpen(l, sut.MakeModQuery())

	err := lua.DoString(l, `
		local query = require("std.serde.query")
		-- values() iterates over q.params map, order unspecified.
		-- Run deserialize+values 5 times to detect non-determinism.
		local results = {}
		for i = 1, 5 do
			local q = query.deserialize("z=last&a=first&m=middle")
			results[i] = table.concat(q:values(), ",")
		end
		for i = 2, 5 do
			if results[i] ~= results[1] then
				error("values() non-deterministic: " .. results[1] .. " vs " .. results[i])
			end
		end
	`)
	assert.NoError(t, err)
}

func TestQueryNit_KeysValuesInconsistentWithEntries(t *testing.T) {
	l := lua.NewState()
	lua.OpenLibraries(l)
	effectual.LuaModOpen(l, sut.MakeModQuery())

	err := lua.DoString(l, `
		local query = require("std.serde.query")
		local q = query.deserialize("z=last&a=first&m=middle")
		q:sort()
		local keys = q:keys()
		local entries = q:entries()
		local entryKeys = {}
		for i = 1, #entries do
			entryKeys[#entryKeys + 1] = entries[i][1]
		end
		local keysStr = table.concat(keys, ",")
		local entryKeysStr = table.concat(entryKeys, ",")
		assert(keysStr == entryKeysStr,
			"keys() and entries() order inconsistent: keys=" .. keysStr .. " entries=" .. entryKeysStr)
	`)
	assert.NoError(t, err)
}