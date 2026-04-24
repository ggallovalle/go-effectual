package std_test

import (
	"testing"

	lua "github.com/speedata/go-lua"
	sut "github.com/ggallovalle/go-effectual/std"
	"github.com/ggallovalle/go-effectual"
	"github.com/stretchr/testify/assert"
)

func Test_LibGoUrlBug_UrlParseInvalidPortDoesNotCrash(t *testing.T) {
	l := lua.NewState()
	lua.OpenLibraries(l)
	effectual.LuaModOpen(l, sut.MakeModUrl())

	err := lua.DoString(l, `
		local url = require("std.url")
		-- url.Parse returns error for invalid port "notanumber".
		-- Previously crashed on nil u due to missing portInferred init.
		local u = url.deserialize("http://host:notanumber/path")
		assert(u ~= nil, "should return non-nil url on parse error")
		assert(tostring(u) == "http://host:notanumber/path")
	`)
	assert.NoError(t, err)
}

func Test_LibGoUrlBug_UrlRawStaleAfterPathDiv(t *testing.T) {
	l := lua.NewState()
	lua.OpenLibraries(l)
	effectual.LuaModOpen(l, sut.MakeModUrl())

	err := lua.DoString(l, `
		local url = require("std.url")
		local u = url.deserialize("http://example.com/path")
		assert(tostring(u) == "http://example.com/path")
		-- __div (path join) mutates path but does NOT rebuild u.raw.
		-- u.raw still returns old URL string.
		local u2 = u / "sub"
		assert(tostring(u2) == "http://example.com/path/sub",
			"expected 'http://example.com/path/sub' but got: " .. tostring(u2))
	`)
	assert.NoError(t, err)
}

func Test_LibGoUrlBug_UrlRawStaleAfterAddQuery(t *testing.T) {
	l := lua.NewState()
	lua.OpenLibraries(l)
	effectual.LuaModOpen(l, sut.MakeModUrl())

	err := lua.DoString(l, `
		local url = require("std.url")
		local u = url.deserialize("http://example.com/path")
		assert(tostring(u) == "http://example.com/path")
		-- add_query modifies u.query but does NOT rebuild u.raw.
		u:add_query("foo", "bar")
		assert(tostring(u) == "http://example.com/path?foo=bar",
			"expected updated raw but got: " .. tostring(u))
	`)
	assert.NoError(t, err)
}