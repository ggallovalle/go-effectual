local semver = require("std.semver")

local Suite = {
	name = "std.semver",
	cases = {
		{
			name = "Version: new() valid",
			fn = function(ctx)
				local v = semver.new("1.2.3")
				ctx:expect(v):not_nil()
				ctx:expect(v.major):equals(1)
				ctx:expect(v.minor):equals(2)
				ctx:expect(v.patch):equals(3)
			end,
		},
		{
			name = "Version: new() invalid",
			fn = function(ctx)
				local ok, err = pcall(semver.new, "not-a-version")
				ctx:expect(ok):is_false()
				ctx:expect(err):not_nil()
			end,
		},
		{
			name = "Version: __tostring",
			fn = function(ctx)
				local v = semver.new("2.3.4")
				ctx:expect(tostring(v)):equals("2.3.4")
			end,
		},
		{
			name = "Version: __lt",
			fn = function(ctx)
				local v1 = semver.new("1.0.0")
				local v2 = semver.new("2.0.0")
				ctx:expect(v1 < v2):is_true()
				ctx:expect(v2 < v1):is_false()
			end,
		},
		{
			name = "Version: __le",
			fn = function(ctx)
				local v1 = semver.new("1.0.0")
				local v2 = semver.new("2.0.0")
				ctx:expect(v1 <= v2):is_true()
				ctx:expect(v1 <= v1):is_true()
				ctx:expect(v2 <= v1):is_false()
			end,
		},
		{
			name = "Version: __eq",
			fn = function(ctx)
				local v1 = semver.new("1.0.0")
				local v2 = semver.new("1.0.0")
				local v3 = semver.new("2.0.0")
				ctx:expect(v1 == v2):is_true()
				ctx:expect(v1 == v3):is_false()
			end,
		},
		{
			name = "Range: range_new() valid",
			fn = function(ctx)
				local r = semver.range_new(">=1.0.0 <2.0.0")
				ctx:expect(r):not_nil()
			end,
		},
		{
			name = "Range: range_new() invalid",
			fn = function(ctx)
				local ok, err = pcall(semver.range_new, "not-a-range")
				ctx:expect(ok):is_false()
				ctx:expect(err):not_nil()
			end,
		},
		{
			name = "Range: contains",
			fn = function(ctx)
				local r = semver.range_new(">=1.0.0 <2.0.0")
				local v1 = semver.new("1.5.0")
				local v2 = semver.new("2.0.0")
				ctx:expect(r:contains(v1)):is_true()
				ctx:expect(r:contains(v2)):is_false()
			end,
		},
		{
			name = "Range: intersect",
			fn = function(ctx)
				local r1 = semver.range_new(">=1.0.0")
				local r2 = semver.range_new("<2.0.0")
				local r = r1:intersect(r2)
				local v1 = semver.new("1.5.0")
				local v2 = semver.new("2.0.0")
				ctx:expect(r:contains(v1)):is_true()
				ctx:expect(r:contains(v2)):is_false()
			end,
		},
		{
			name = "Range: union",
			fn = function(ctx)
				local r1 = semver.range_new("<1.0.0")
				local r2 = semver.range_new(">=2.0.0")
				local r = r1:union(r2)
				local v1 = semver.new("0.5.0")
				local v2 = semver.new("2.5.0")
				local v3 = semver.new("1.5.0")
				ctx:expect(r:contains(v1)):is_true()
				ctx:expect(r:contains(v2)):is_true()
				ctx:expect(r:contains(v3)):is_false()
			end,
		},
	},
}

return Suite
