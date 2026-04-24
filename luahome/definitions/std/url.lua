---@meta std.url

---@class (exact) std.url.Url : userdata
---@field scheme string|nil
---@field host string|nil
---@field port integer|nil
---@field port_inferred integer
---@field username string|nil
---@field password string|nil
---@field path std.path.PathBuf
---@field query std.serde.query.Query
---@field fragment string|nil
local Url = {}

---@param path string
function Url:add_query(path) end

---@param key string
function Url:remove_query(key) end

local url = {}

---@return std.url.Url
function url.new() end

---@param raw string
---@return std.url.Url
function url.deserialize(raw) end

---@param u std.url.Url
---@return string
function url.serialize(u) end

return url
