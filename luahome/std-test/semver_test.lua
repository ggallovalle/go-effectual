local semver = require("std.semver")

local Suite = {
	name = "std.semver",
	cases = {
		{
			name = "Version: new() valid",
			fn = function(ctx)
				local v = semver.new("1.2.3")
				ctx:expect(v, "expected version"):not_nil()
				ctx:expect(v.major, "expected major == 1"):equals(1)
				ctx:expect(v.minor, "expected minor == 2"):equals(2)
				ctx:expect(v.patch, "expected patch == 3"):equals(3)
			end,
		},
		{
			name = "Version: new() invalid",
			fn = function(ctx)
				local ok, err = pcall(semver.new, "not-a-version")
				ctx:expect(ok, "expected error"):is_false()
				ctx:expect(err, "expected error message"):not_nil()
			end,
		},
		{
			name = "Version: __tostring",
			fn = function(ctx)
				local v = semver.new("2.3.4")
				ctx:expect(tostring(v), "expected '2.3.4'"):equals("2.3.4")
			end,
		},
		{
			name = "Version: __lt",
			fn = function(ctx)
				local v1 = semver.new("1.0.0")
				local v2 = semver.new("2.0.0")
				ctx:expect(v1 < v2, "expected 1.0.0 < 2.0.0"):is_true()
				ctx:expect(v2 < v1, "expected 2.0.0 not < 1.0.0"):is_false()
			end,
		},
		{
			name = "Version: __le",
			fn = function(ctx)
				local v1 = semver.new("1.0.0")
				local v2 = semver.new("2.0.0")
				ctx:expect(v1 <= v2, "expected 1.0.0 <= 2.0.0"):is_true()
				ctx:expect(v1 <= v1, "expected 1.0.0 <= 1.0.0"):is_true()
				ctx:expect(v2 <= v1, "expected 2.0.0 not <= 1.0.0"):is_false()
			end,
		},
		{
			name = "Version: __eq",
			fn = function(ctx)
				local v1 = semver.new("1.0.0")
				local v2 = semver.new("1.0.0")
				local v3 = semver.new("2.0.0")
				ctx:expect(v1 == v2, "expected 1.0.0 == 1.0.0"):is_true()
				ctx:expect(v1 == v3, "expected 1.0.0 ~= 2.0.0"):is_false()
			end,
		},
		{
			name = "Range: range_new() valid",
			fn = function(ctx)
				local r = semver.range_new(">=1.0.0 <2.0.0")
				ctx:expect(r, "expected range"):not_nil()
			end,
		},
		{
			name = "Range: range_new() invalid",
			fn = function(ctx)
				local ok, err = pcall(semver.range_new, "not-a-range")
				ctx:expect(ok, "expected error"):is_false()
				ctx:expect(err, "expected error message"):not_nil()
			end,
		},
		{
			name = "Range: contains",
			fn = function(ctx)
				local r = semver.range_new(">=1.0.0 <2.0.0")
				local v1 = semver.new("1.5.0")
				local v2 = semver.new("2.0.0")
				ctx:expect(r:contains(v1), "expected 1.5.0 to be contained"):is_true()
				ctx:expect(r:contains(v2), "expected 2.0.0 not to be contained"):is_false()
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
				ctx:expect(r:contains(v1), "expected 1.5.0 to be contained"):is_true()
				ctx:expect(r:contains(v2), "expected 2.0.0 not to be contained"):is_false()
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
				ctx:expect(r:contains(v1), "expected 0.5.0 to be contained"):is_true()
				ctx:expect(r:contains(v2), "expected 2.5.0 to be contained"):is_true()
				ctx:expect(r:contains(v3), "expected 1.5.0 not to be contained"):is_false()
			end,
		},
	},
}

return Suite
