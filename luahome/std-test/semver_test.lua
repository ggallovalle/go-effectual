local semver = require("std.semver")

local Suite = {
	name = "std.semver",
	cases = {
		{
			name = "Version: new() valid",
			fn = function(_ctx)
				local v = semver.new("1.2.3")
				assert(v ~= nil, "expected version")
				assert(v.major == 1, "expected major == 1")
				assert(v.minor == 2, "expected minor == 2")
				assert(v.patch == 3, "expected patch == 3")
			end,
		},
		{
			name = "Version: new() invalid",
			fn = function(_ctx)
				local ok, err = pcall(semver.new, "not-a-version")
				assert(not ok, "expected error")
				assert(err ~= nil, "expected error message")
			end,
		},
		{
			name = "Version: __tostring",
			fn = function(_ctx)
				local v = semver.new("2.3.4")
				assert(tostring(v) == "2.3.4", "expected '2.3.4' but got " .. tostring(v))
			end,
		},
		{
			name = "Version: __lt",
			fn = function(_ctx)
				local v1 = semver.new("1.0.0")
				local v2 = semver.new("2.0.0")
				assert(v1 < v2, "expected 1.0.0 < 2.0.0")
				assert(not (v2 < v1), "expected 2.0.0 not < 1.0.0")
			end,
		},
		{
			name = "Version: __le",
			fn = function(_ctx)
				local v1 = semver.new("1.0.0")
				local v2 = semver.new("2.0.0")
				assert(v1 <= v2, "expected 1.0.0 <= 2.0.0")
				assert(v1 <= v1, "expected 1.0.0 <= 1.0.0")
				assert(not (v2 <= v1), "expected 2.0.0 not <= 1.0.0")
			end,
		},
		{
			name = "Version: __eq",
			fn = function(_ctx)
				local v1 = semver.new("1.0.0")
				local v2 = semver.new("1.0.0")
				local v3 = semver.new("2.0.0")
				assert(v1 == v2, "expected 1.0.0 == 1.0.0")
				assert(v1 ~= v3, "expected 1.0.0 ~= 2.0.0")
			end,
		},
		{
			name = "Range: range_new() valid",
			fn = function(_ctx)
				local r = semver.range_new(">=1.0.0 <2.0.0")
				assert(r ~= nil, "expected range")
			end,
		},
		{
			name = "Range: range_new() invalid",
			fn = function(_ctx)
				local ok, err = pcall(semver.range_new, "not-a-range")
				assert(not ok, "expected error")
				assert(err ~= nil, "expected error message")
			end,
		},
		{
			name = "Range: contains",
			fn = function(_ctx)
				local r = semver.range_new(">=1.0.0 <2.0.0")
				local v1 = semver.new("1.5.0")
				local v2 = semver.new("2.0.0")
				assert(r:contains(v1), "expected 1.5.0 to be contained")
				assert(not r:contains(v2), "expected 2.0.0 not to be contained")
			end,
		},
		{
			name = "Range: intersect",
			fn = function(_ctx)
				local r1 = semver.range_new(">=1.0.0")
				local r2 = semver.range_new("<2.0.0")
				local r = r1:intersect(r2)
				local v1 = semver.new("1.5.0")
				local v2 = semver.new("2.0.0")
				assert(r:contains(v1), "expected 1.5.0 to be contained")
				assert(not r:contains(v2), "expected 2.0.0 not to be contained")
			end,
		},
		{
			name = "Range: union",
			fn = function(_ctx)
				local r1 = semver.range_new("<1.0.0")
				local r2 = semver.range_new(">=2.0.0")
				local r = r1:union(r2)
				local v1 = semver.new("0.5.0")
				local v2 = semver.new("2.5.0")
				local v3 = semver.new("1.5.0")
				assert(r:contains(v1), "expected 0.5.0 to be contained")
				assert(r:contains(v2), "expected 2.5.0 to be contained")
				assert(not r:contains(v3), "expected 1.5.0 not to be contained")
			end,
		},
	},
}

return Suite
