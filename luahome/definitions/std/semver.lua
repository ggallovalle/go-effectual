---@meta std.semver

---@class (exact) std.semver.Version : userdata
---@field major integer
---@field minor integer
---@field patch integer
---@operator lt(std.semver.Version): boolean
---@operator le(std.semver.Version): boolean
---@operator eq(std.semver.Version): boolean
local Version = {}

---@class (exact) std.semver.Range : userdata
local Range = {}

--- Performs logical intersection on two ranges
---@param range std.semver.Range
---@return std.semver.Range
function Range:intersect(range) end

--- Performs logical union on two ranges
---@param range std.semver.Range
---@return std.semver.Range
function Range:union(range) end

--- Checks if the range contains the given version
---@param version std.semver.Version
---@return boolean
function Range:contains(version) end

local semver = {}

--- Creates a new Version from a string
---@param version string (e.g., "1.2.3")
---@return std.semver.Version
---@raise if version string is invalid
function semver.new(version) end

--- Creates a new Range from a semver range string
---@param range string (e.g., ">=1.0.0 <2.0.0")
---@return std.semver.Range
---@raise if range string is invalid
function semver.range_new(range) end

return semver
