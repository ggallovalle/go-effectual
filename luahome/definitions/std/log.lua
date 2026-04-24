
---@meta std.log

---@class std.log : std.log.Logger
---@field LEVELS std.log.LogLevel
---@field default std.log.Logger
local log = {}

---@param logger
---@return std.log.Logger
function log.new(logger) end

---@enum std.log.LogLevel
local LogLevel = {
    DEBUG = "DEBUG",
    INFO = "INFO",
    WARN = "WARN",
    ERROR = "ERROR",
}

---@alias std.log.Level
---| '"DEBUG"' # Debug level
---| '"INFO"'  # Info level
---| '"WARN"'  # Warn level
---| '"ERROR"' # Error level

---@class std.log.Logger
local Logger = {}

---@param msg string
---@param attrs? table
function Logger:debug(msg, attrs) end

---@param msg string
---@param attrs? table
function Logger:info(msg, attrs) end

---@param msg string
---@param attrs? table
function Logger:warn(msg, attrs) end

---@param msg string
---@param attrs? table
function Logger:error(msg, attrs) end

---@param level std.log.Level
---@param msg string
---@param attrs? table
function Logger:log(level, msg, attrs) end

---@return std.log.Level
function Logger:level() end

return log
