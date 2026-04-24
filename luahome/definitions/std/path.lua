---@meta std.path

---@class (exact) std.path.PathBuf : userdata
---@operator div(std.path.PathBuf|string): std.path.PathBuf
---@operator concat(std.path.PathBuf|string): string
---@field parent std.path.PathBuf|nil
---@field components string[]
---@field ancestors std.path.PathBuf[]
---@field file_name string|nil
---@field extension string|nil
---@field file_stem string|nil
---@field has_root boolean
---@field is_relative boolean
---@field is_absolute boolean
local PathBuf = {}

--- Appends the given path segments to self, returning a new PathBuf
---@param path string
function PathBuf:push(path) end

--- Removes the last path component from self, returning true on success
---@return boolean
function PathBuf:pop() end

--- Joins self with the given path, returning a new PathBuf. Absolute paths replace self
---@param path string
---@return std.path.PathBuf
function PathBuf:join(path) end

--- Returns true if self ends with the given path segment
---@param child string
---@return boolean
function PathBuf:ends_with(child) end

--- Returns true if self starts with the given path prefix
---@param base string
---@return boolean
function PathBuf:starts_with(base) end

--- Strips the given prefix from self, returning a new PathBuf or an error
---@param prefix string
---@return std.path.PathBuf?
---@return string? Error message if prefix is not found
function PathBuf:strip_prefix(prefix) end

--- Sets the file extension, returning a new PathBuf with the changed extension
---@param ext string
---@return std.path.PathBuf
function PathBuf:with_extension(ext) end

--- Sets the file name component, returning a new PathBuf
---@param name string
---@return std.path.PathBuf
function PathBuf:with_file_name(name) end

local path = {}

---@type string
path.MAIN_SEPARATOR = "/"

--- Creates a new PathBuf from the given path string
---@param value string
---@return std.path.PathBuf
function path.new(value) end

--- Joins multiple path segments together, returning a new PathBuf
---@param ... string|std.path.PathBuf
---@return std.path.PathBuf
function path.join(...) end

--- Converts the given path to an absolute path based on the current working directory
---@param path string|std.path.PathBuf
---@return std.path.PathBuf?
---@return string? Error message if path is empty
function path.absolute(path) end

---@class std.path.posix : std.path
local posix = {}

---@class std.path.win32 : std.path
local win32 = {}

return path
