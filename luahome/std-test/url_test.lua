local url = require("std.url")

local Suite = {
	name = "std.url",
	cases = {
		{
			name = "Bug: UrlParseInvalidPortDoesNotCrash",
			fn = function(ctx)
				local u = url.deserialize("http://host:notanumber/path")
				ctx:expect(u):not_nil()
				ctx:expect(tostring(u)):equals("http://host:notanumber/path")
			end,
		},
		{
			name = "Bug: UrlRawStaleAfterPathDiv",
			fn = function(ctx)
				local u = url.deserialize("http://example.com/path")
				ctx:expect(tostring(u)):equals("http://example.com/path")
				local u2 = u / "sub"
				ctx:expect(tostring(u2)):equals("http://example.com/path/sub")
			end,
		},
		{
			name = "Bug: UrlRawStaleAfterAddQuery",
			fn = function(ctx)
				local u = url.deserialize("http://example.com/path")
				ctx:expect(tostring(u)):equals("http://example.com/path")
				u:add_query("foo", "bar")
				ctx:expect(tostring(u)):equals("http://example.com/path?foo=bar")
			end,
		},
	},
}

return Suite
