---@meta std.serde.query

---@class (exact) std.serde.query.Query : userdata
---@field size integer
local Query = {}

---@param key string
---@return boolean
function Query:has(key) end

---@param key string
---@return string|nil
function Query:get(key) end

---@param key string
---@return string[]
function Query:get_all(key) end

---@param key string
---@param value string
function Query:set(key, value) end

---@param key string
---@param value string
function Query:append(key, value) end

---@param key string
function Query:delete(key) end

function Query:sort() end

---@return string
function Query:to_string() end

---@return string[]
function Query:keys() end

---@return string[]
function Query:values() end

---@return {[1]: string, [2]: string}[]
function Query:entries() end

local query = {}

---@return std.serde.query.Query
function query.new() end

---@param raw string
---@return std.serde.query.Query
function query.deserialize(raw) end

---@param q std.serde.query.Query
---@return string
function query.serialize(q) end

return query
