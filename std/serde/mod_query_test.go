package serde_test

import (
	"testing"

	lua "github.com/Shopify/go-lua"
	sut "github.com/ggallovalle/go-effectual/std/serde"
	"github.com/ggallovalle/go-effectual"
	"github.com/stretchr/testify/assert"
)

func TestQueryNew(t *testing.T) {
	l := lua.NewState()
	lua.OpenLibraries(l)
	effectual.LuaModOpen(l, sut.MakeModQuery())

	err := lua.DoString(l, `
		local query = require("std.serde.query")
		local q = query.new()
		assert(q.size == 0)
	`)
	assert.NoError(t, err)
}

func TestQueryDeserialize(t *testing.T) {
	l := lua.NewState()
	lua.OpenLibraries(l)
	effectual.LuaModOpen(l, sut.MakeModQuery())

	err := lua.DoString(l, `
		local query = require("std.serde.query")
		local q = query.deserialize("foo=bar&baz=qux")
		assert(q.size == 2)
		assert(q:has("foo") == true)
		assert(q:has("missing") == false)
	`)
	assert.NoError(t, err)
}

func TestQueryGet(t *testing.T) {
	l := lua.NewState()
	lua.OpenLibraries(l)
	effectual.LuaModOpen(l, sut.MakeModQuery())

	err := lua.DoString(l, `
		local query = require("std.serde.query")
		local q = query.deserialize("foo=bar&foo=baz")
		assert(q:get("foo") == "bar")
		assert(q:get("missing") == nil)
	`)
	assert.NoError(t, err)
}

func TestQueryGetAll(t *testing.T) {
	l := lua.NewState()
	lua.OpenLibraries(l)
	effectual.LuaModOpen(l, sut.MakeModQuery())

	err := lua.DoString(l, `
		local query = require("std.serde.query")
		local q = query.deserialize("tags=lua&tags=stdlib&tags=v2")
		local tags = q:get_all("tags")
		assert(#tags == 3)
		assert(tags[1] == "lua")
		assert(tags[2] == "stdlib")
		assert(tags[3] == "v2")
	`)
	assert.NoError(t, err)
}

func TestQueryAppend(t *testing.T) {
	l := lua.NewState()
	lua.OpenLibraries(l)
	effectual.LuaModOpen(l, sut.MakeModQuery())

	err := lua.DoString(l, `
		local query = require("std.serde.query")
		local q = query.new()
		q:append("foo", "bar")
		assert(q.size == 1)
		q:append("foo", "baz")
		local vals = q:get_all("foo")
		assert(#vals == 2)
		assert(vals[1] == "bar")
		assert(vals[2] == "baz")
	`)
	assert.NoError(t, err)
}

func TestQuerySet(t *testing.T) {
	l := lua.NewState()
	lua.OpenLibraries(l)
	effectual.LuaModOpen(l, sut.MakeModQuery())

	err := lua.DoString(l, `
		local query = require("std.serde.query")
		local q = query.deserialize("foo=old&bar=val")
		q:set("foo", "new")
		assert(q:get("foo") == "new")
		assert(q:get("bar") == "val")
	`)
	assert.NoError(t, err)
}

func TestQueryDelete(t *testing.T) {
	l := lua.NewState()
	lua.OpenLibraries(l)
	effectual.LuaModOpen(l, sut.MakeModQuery())

	err := lua.DoString(l, `
		local query = require("std.serde.query")
		local q = query.deserialize("foo=bar&baz=qux")
		q:delete("foo")
		assert(q.size == 1)
		assert(q:get("foo") == nil)
		assert(q:get("baz") == "qux")
	`)
	assert.NoError(t, err)
}

func TestQuerySort(t *testing.T) {
	l := lua.NewState()
	lua.OpenLibraries(l)
	effectual.LuaModOpen(l, sut.MakeModQuery())

	err := lua.DoString(l, `
		local query = require("std.serde.query")
		local q = query.deserialize("z=last&a=first&m=middle")
		q:sort()
		local str = q:to_string()
		assert(str == "a=first&m=middle&z=last", "expected sorted but got: " .. str)
	`)
	assert.NoError(t, err)
}

func TestQueryToString(t *testing.T) {
	l := lua.NewState()
	lua.OpenLibraries(l)
	effectual.LuaModOpen(l, sut.MakeModQuery())

	err := lua.DoString(l, `
		local query = require("std.serde.query")
		local q = query.deserialize("foo=bar&baz=qux")
		local str = q:to_string()
		assert(str == "baz=qux&foo=bar", "expected 'baz=qux&foo=bar' (keys sorted) but got '" .. str .. "'")
	`)
	assert.NoError(t, err)
}

func TestQuerySerialize(t *testing.T) {
	l := lua.NewState()
	lua.OpenLibraries(l)
	effectual.LuaModOpen(l, sut.MakeModQuery())

	err := lua.DoString(l, `
		local query = require("std.serde.query")
		local q = query.deserialize("foo=bar")
		assert(query.serialize(q) == "foo=bar")
	`)
	assert.NoError(t, err)
}

func TestQueryKeys(t *testing.T) {
	l := lua.NewState()
	lua.OpenLibraries(l)
	effectual.LuaModOpen(l, sut.MakeModQuery())

	err := lua.DoString(l, `
		local query = require("std.serde.query")
		local q = query.deserialize("a=1&b=2&c=3")
		local keys = q:keys()
		assert(#keys == 3)
	`)
	assert.NoError(t, err)
}

func TestQueryValues(t *testing.T) {
	l := lua.NewState()
	lua.OpenLibraries(l)
	effectual.LuaModOpen(l, sut.MakeModQuery())

	err := lua.DoString(l, `
		local query = require("std.serde.query")
		local q = query.deserialize("a=1&b=2")
		local vals = q:values()
		assert(#vals == 2)
	`)
	assert.NoError(t, err)
}

func TestQueryEntries(t *testing.T) {
	l := lua.NewState()
	lua.OpenLibraries(l)
	effectual.LuaModOpen(l, sut.MakeModQuery())

	err := lua.DoString(l, `
		local query = require("std.serde.query")
		local q = query.deserialize("a=1&b=2")
		local entries = q:entries()
		assert(type(entries) == "table", "expected table, got " .. type(entries))
		assert(#entries == 2, "expected 2 entries, got " .. #entries)
		-- Check that we have both keys (order may vary due to map iteration)
		local keys = {}
		for i = 1, 2 do
			assert(entries[i][1] ~= nil, "entry " .. i .. " has nil key")
			keys[entries[i][1]] = true
			assert(entries[i][2] ~= nil, "entry " .. i .. " has nil value")
		end
		assert(keys["a"] == true, "expected key 'a'")
		assert(keys["b"] == true, "expected key 'b'")
	`)
	assert.NoError(t, err)
}