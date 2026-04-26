---@meta std.path

---@class (exact) std.path.Path : userdata
---@operator div(std.path.Path|string): std.path.Path
---@operator concat(std.path.Path|string): string
---@field parent std.path.Path|nil
---@field components string[]
---@field ancestors std.path.Path[]
---@field file_name string|nil
---@field extension string|nil
---@field file_stem string|nil
---@field has_root boolean
---@field is_relative boolean
---@field is_absolute boolean
local Path = {}

--- Appends the given path segments to self, returning a new Path
---@param path string
function Path:push(path) end

--- Removes the last path component from self, returning true on success
---@return boolean
function Path:pop() end

--- Joins self with the given path, returning a new Path. Absolute paths replace self
---@param path string
---@return std.path.Path
function Path:join(path) end

--- Returns true if self ends with the given path segment
---@param child string
---@return boolean
function Path:ends_with(child) end

--- Returns true if self starts with the given path prefix
---@param base string
---@return boolean
function Path:starts_with(base) end

--- Strips the given prefix from self, returning a new Path or an error
---@param prefix string
---@return std.path.Path?
---@return string? Error message if prefix is not found
function Path:strip_prefix(prefix) end

--- Sets the file extension, returning a new Path with the changed extension
---@param ext string
---@return std.path.Path
function Path:with_extension(ext) end

--- Sets the file name component, returning a new Path
---@param name string
---@return std.path.Path
function Path:with_file_name(name) end

local path = {}

---@type string
path.MAIN_SEPARATOR = "/"

--- Creates a new Path from the given path string
---@param value string
---@return std.path.Path
function path.new(value) end

--- Joins multiple path segments together, returning a new Path
---@param ... string|std.path.Path
---@return std.path.Path
function path.join(...) end

--- Converts the given path to an absolute path based on the current working directory
---@param path string|std.path.Path
---@return std.path.Path?
---@return string? Error message if path is empty
function path.absolute(path) end

---@class std.path.posix : std.path
local posix = {}

---@class std.path.win32 : std.path
local win32 = {}

return path
