local query = require("std.serde.query")

local Suite = {
	name = "std.serde.query",
	cases = {
		{
			name = "New: empty query",
			fn = function(ctx)
				local q = query.new()
				ctx:expect(q.size):equals(0)
			end,
		},
		{
			name = "Deserialize: parse query string",
			fn = function(ctx)
				local q = query.deserialize("foo=bar&baz=qux")
				ctx:expect(q.size):equals(2)
				ctx:expect(q:has("foo")):is_true()
				ctx:expect(q:has("missing")):is_false()
			end,
		},
		{
			name = "Get: first value",
			fn = function(ctx)
				local q = query.deserialize("foo=bar&foo=baz")
				ctx:expect(q:get("foo")):equals("bar")
				ctx:expect(q:get("missing")):is_nil()
			end,
		},
		{
			name = "GetAll: multiple values",
			fn = function(ctx)
				local q = query.deserialize("tags=lua&tags=stdlib&tags=v2")
				local tags = q:get_all("tags")
				ctx:expect(#tags):equals(3)
				ctx:expect(tags[1]):equals("lua")
				ctx:expect(tags[2]):equals("stdlib")
				ctx:expect(tags[3]):equals("v2")
			end,
		},
		{
			name = "Append: add values",
			fn = function(ctx)
				local q = query.new()
				q:append("foo", "bar")
				ctx:expect(q.size):equals(1)
				q:append("foo", "baz")
				local vals = q:get_all("foo")
				ctx:expect(#vals):equals(2)
				ctx:expect(vals[1]):equals("bar")
				ctx:expect(vals[2]):equals("baz")
			end,
		},
		{
			name = "Set: replace value",
			fn = function(ctx)
				local q = query.deserialize("foo=old&bar=val")
				q:set("foo", "new")
				ctx:expect(q:get("foo")):equals("new")
				ctx:expect(q:get("bar")):equals("val")
			end,
		},
		{
			name = "Delete: remove key",
			fn = function(ctx)
				local q = query.deserialize("foo=bar&baz=qux")
				q:delete("foo")
				ctx:expect(q.size):equals(1)
				ctx:expect(q:get("foo")):is_nil()
				ctx:expect(q:get("baz")):equals("qux")
			end,
		},
		{
			name = "Sort: keys sorted",
			fn = function(ctx)
				local q = query.deserialize("z=last&a=first&m=middle")
				q:sort()
				local str = q:to_string()
				ctx:expect(str):equals("a=first&m=middle&z=last")
			end,
		},
		{
			name = "ToString: sorted keys",
			fn = function(ctx)
				local q = query.deserialize("foo=bar&baz=qux")
				local str = q:to_string()
				ctx:expect(str):equals("baz=qux&foo=bar")
			end,
		},
		{
			name = "Serialize: roundtrip",
			fn = function(ctx)
				local q = query.deserialize("foo=bar")
				ctx:expect(query.serialize(q)):equals("foo=bar")
			end,
		},
		{
			name = "Keys: returns all keys",
			fn = function(ctx)
				local q = query.deserialize("a=1&b=2&c=3")
				local keys = q:keys()
				ctx:expect(#keys):equals(3)
			end,
		},
		{
			name = "Values: returns all values",
			fn = function(ctx)
				local q = query.deserialize("a=1&b=2")
				local vals = q:values()
				ctx:expect(#vals):equals(2)
			end,
		},
		{
			name = "Entries: returns key-value pairs",
			fn = function(ctx)
				local q = query.deserialize("a=1&b=2")
				local entries = q:entries()
				ctx:expect(type(entries)):equals("table")
				ctx:expect(#entries):equals(2)
				ctx:expect(entries[1][1]):not_nil()
				ctx:expect(entries[1][2]):not_nil()
				ctx:expect(entries[2][1]):not_nil()
				ctx:expect(entries[2][2]):not_nil()
				local keys = {}
				for i = 1, 2 do
					keys[entries[i][1]] = true
				end
				ctx:expect(keys["a"]):is_true()
				ctx:expect(keys["b"]):is_true()
			end,
		},
		{
			name = "Bug: malformed query silently ignored",
			fn = function(ctx)
				local q = query.deserialize("foo=%")
				ctx:expect(q.size):equals(0)
			end,
		},
		{
			name = "Bug: pairs iterator yields entries",
			fn = function(ctx)
				local q = query.deserialize("a=1&b=2")
				local count = 0
				for k, v in pairs(q) do
					count = count + 1
				end
				ctx:expect(count > 0):is_true()
			end,
		},
		{
			name = "Bug: entries order consistent after sort",
			fn = function(ctx)
				local q = query.deserialize("z=last&a=first&m=middle")
				local entries = q:entries()
				local keys = {}
				for i = 1, #entries do
					keys[#keys + 1] = entries[i][1]
				end
				local sorted = table.concat(keys, ",")
				q:sort()
				local entries2 = q:entries()
				local keys2 = {}
				for i = 1, #entries2 do
					keys2[#keys2 + 1] = entries2[i][1]
				end
				local sorted2 = table.concat(keys2, ",")
				ctx:expect(sorted):equals(sorted2)
			end,
		},
		{
			name = "Nit: keys order deterministic",
			fn = function(ctx)
				local q = query.deserialize("z=last&a=first&m=middle")
				local keys1 = q:keys()
				local keys2 = q:keys()
				ctx:expect(table.concat(keys1, ",")):equals(table.concat(keys2, ","))
			end,
		},
		{
			name = "Nit: values order deterministic",
			fn = function(ctx)
				local results = {}
				for i = 1, 5 do
					local q = query.deserialize("z=last&a=first&m=middle")
					results[i] = table.concat(q:values(), ",")
				end
				for i = 2, 5 do
					ctx:expect(results[i]):equals(results[1])
				end
			end,
		},
		{
			name = "Nit: keys and entries order consistent",
			fn = function(ctx)
				local q = query.deserialize("z=last&a=first&m=middle")
				q:sort()
				local keys = q:keys()
				local entries = q:entries()
				local entryKeys = {}
				for i = 1, #entries do
					entryKeys[#entryKeys + 1] = entries[i][1]
				end
				ctx:expect(table.concat(keys, ",")):equals(table.concat(entryKeys, ","))
			end,
		},
	},
}

return Suite
