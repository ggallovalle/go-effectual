local function collect_paths(path_str)
    if not path_str or path_str == "" then
        return {}
    end
    local t = {}
    for p in string.gmatch(path_str, "[^;]+") do
        table.insert(t, p)
    end
    return t
end

local function print_numbered_paths(name, path_str)
    print(name .. ":")
    for i, p in ipairs(collect_paths(path_str)) do
        print("  " .. i .. ": " .. p)
    end
end

print_numbered_paths("package.path", package.path)
print_numbered_paths("package.cpath", package.cpath)

local log = require("std.log")
local dkjson = require("dkjson")

local version = dkjson.encode({
    version = _VERSION,
    from = "dkjson",
})

log:info("hello logger", { version = _VERSION, dkjson = version })
