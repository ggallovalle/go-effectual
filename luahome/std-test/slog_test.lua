local log = require("std.log")

local Suite = {
	name = "std.slog",
	deps = {"logger"},
	cases = {
		{
			name = "Level: DEBUG",
			deps = {{name = "logger", params = {level = "DEBUG"}}},
			fn = function(ctx)
				ctx:expect(ctx.ext.logger:level()):equals("DEBUG")
			end,
		},
		{
			name = "Level: INFO",
			deps = {{name = "logger", params = {level = "INFO"}}},
			fn = function(ctx)
				ctx:expect(ctx.ext.logger:level()):equals("INFO")
			end,
		},
		{
			name = "Level: WARN",
			deps = {{name = "logger", params = {level = "WARN"}}},
			fn = function(ctx)
				ctx:expect(ctx.ext.logger:level()):equals("WARN")
			end,
		},
		{
			name = "Level: ERROR",
			deps = {{name = "logger", params = {level = "ERROR"}}},
			fn = function(ctx)
				ctx:expect(ctx.ext.logger:level()):equals("ERROR")
			end,
		},
		{
			name = "Default: returns the logger",
			fn = function(ctx)
				ctx:expect(log.default):not_nil()
			end,
		},
		{
			name = "Default: level delegates",
			fn = function(ctx)
				local lv = log:level()
				ctx:expect(lv):equals("DEBUG")
			end,
		},
	},
}

return Suite
