package serde_test

import (
	"testing"

	lua "github.com/speedata/go-lua"
	sut "github.com/ggallovalle/go-effectual/std/serde"
	"github.com/ggallovalle/go-effectual"
	"github.com/stretchr/testify/assert"
)

func TestQueryBug_MalformedQueryIgnored(t *testing.T) {
	l := lua.NewState()
	lua.OpenLibraries(l)
	effectual.LuaModOpen(l, sut.MakeModQuery())

	err := lua.DoString(l, `
		local query = require("std.serde.query")
		-- url.ParseQuery error is silently ignored.
		-- Malformed query like "foo=%" should fail but doesn't.
		local q = query.deserialize("foo=%")
		assert(q.size == 0, "malformed query silently ignored, expected size 0")
	`)
	assert.NoError(t, err)
}

func TestQueryBug_PairsIteratorIsNoOp(t *testing.T) {
	l := lua.NewState()
	lua.OpenLibraries(l)
	effectual.LuaModOpen(l, sut.MakeModQuery())

	err := lua.DoString(l, `
		local query = require("std.serde.query")
		local q = query.deserialize("a=1&b=2")
		-- __pairs returns a no-op function that yields nothing.
		-- This for loop should iterate but silently produces nothing.
		local count = 0
		for k, v in pairs(q) do
			count = count + 1
		end
		assert(count > 0, "__pairs iterator should yield entries, got 0")
	`)
	assert.NoError(t, err)
}

func TestQueryBug_EntriesOrderNonDeterministic(t *testing.T) {
	l := lua.NewState()
	lua.OpenLibraries(l)
	effectual.LuaModOpen(l, sut.MakeModQuery())

	err := lua.DoString(l, `
		local query = require("std.serde.query")
		local q = query.deserialize("z=last&a=first&m=middle")
		-- entries() iterates over a map, order is unspecified.
		-- Calling sort() first should produce deterministic order,
		-- but entries() doesn't use it.
		local entries = q:entries()
		local keys = {}
		for i = 1, #entries do
			keys[#keys + 1] = entries[i][1]
		end
		local sorted = table.concat(keys, ",")
		-- Run multiple times to detect non-determinism.
		-- If sort() is called, order should be "a,m,z".
		q:sort()
		local entries2 = q:entries()
		local keys2 = {}
		for i = 1, #entries2 do
			keys2[#keys2 + 1] = entries2[i][1]
		end
		local sorted2 = table.concat(keys2, ",")
		assert(sorted == sorted2, "entries() order non-deterministic: " .. sorted .. " vs " .. sorted2)
	`)
	assert.NoError(t, err)
}