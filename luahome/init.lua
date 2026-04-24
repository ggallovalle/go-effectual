local log = require("std.log")
--- local dkjson = require("dkjson")
--- local hello_from_dkjson = dkjson.encode(1)

log:info("hello logger", { version = _VERSION, dkjson_says = "ups string.match required" })

print("_VERSION = " .. _VERSION)

---@param path_str string
local function collect_paths(path_str)
    if path_str == nil or path_str == "" then
        return {}
    end
    local t = {}
    local pos = 1
    while true do
        local found = string.find(path_str, ";", pos, true)
        if found then
            table.insert(t, string.sub(path_str, pos, found - 1))
            pos = found + 1
        else
            local last = string.sub(path_str, pos)
            if last ~= "" then
                table.insert(t, last)
            end
            break
        end
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
