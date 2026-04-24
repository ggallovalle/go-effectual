package std

import (
	"context"
	"log/slog"
	"strings"
	"text/template"

	"github.com/Shopify/go-lua"
)

type LibGoLogSlug struct {
	Name string
}

const LibNameGoLogSlug = "std.log"
const slugLoggerHandle = "go/std/log/slug/Logger*"

func MakeLibGoLogSlug() LibGoLogSlug {
	return LibGoLogSlug{Name: LibNameGoLogSlug}
}

func toLogger(l *lua.State) *slog.Logger {
	return lua.CheckUserData(l, 1, slugLoggerHandle).(*slog.Logger)
}

func tableToSlogAttrs(l *lua.State, index int) []slog.Attr {
	var attrs []slog.Attr
	if l.Top() >= index && l.IsTable(index) {
		l.PushNil()
		for l.Next(index) {
			if key, ok := l.ToString(-2); ok {
				value := ConvertLuaToAny(l, -1)
				if value != nil {
					attrs = append(attrs, slog.Any(key, value))
				}
			}
			l.Pop(1)
		}
	}
	return attrs
}

func slogMethodToLua(level slog.Level) lua.Function {
	return func(l *lua.State) int {
		logger := toLogger(l)
		msg, _ := l.ToString(2)
		args := tableToSlogAttrs(l, 3)
		logger.LogAttrs(context.Background(), level, msg, args...)
		return 0
	}
}

var slogHandlerMethods = []lua.RegistryFunction{
	{Name: "debug", Function: slogMethodToLua(slog.LevelDebug)},
	{Name: "info", Function: slogMethodToLua(slog.LevelInfo)},
	{Name: "warn", Function: slogMethodToLua(slog.LevelWarn)},
	{Name: "error", Function: slogMethodToLua(slog.LevelError)},
	{Name: "log", Function: func(l *lua.State) int {
		logger := toLogger(l)
		levelStr, _ := l.ToString(2)
		msg, _ := l.ToString(3)

		var level slog.Level
		switch levelStr {
		case "DEBUG":
			level = slog.LevelDebug
		case "INFO":
			level = slog.LevelInfo
		case "WARN":
			level = slog.LevelWarn
		default:
			level = slog.LevelError
		}

		args := tableToSlogAttrs(l, 4)
		logger.LogAttrs(context.Background(), level, msg, args...)
		return 0
	}},
	{Name: "level", Function: func(l *lua.State) int {
		logger := toLogger(l)
		if logger.Enabled(context.Background(), slog.LevelDebug) {
			l.PushString("DEBUG")
		} else if logger.Enabled(context.Background(), slog.LevelInfo) {
			l.PushString("INFO")
		} else if logger.Enabled(context.Background(), slog.LevelWarn) {
			l.PushString("WARN")
		} else {
			l.PushString("ERROR")
		}
		return 1
	}},
}

var slogLoggerLibrary = []lua.RegistryFunction{
	{Name: "new", Function: func(l *lua.State) int {
		logger, ok := l.ToUserData(1).(*slog.Logger)
		if !ok {
			lua.ArgumentError(l, 1, "expected *slog.Logger")
			return 0
		}

		l.PushUserData(logger)
		lua.SetMetaTableNamed(l, slugLoggerHandle)
		return 1
	}},
}

func (lib *LibGoLogSlug) Open(l *lua.State) int {
	lua.NewLibrary(l, slogLoggerLibrary)
	moduleIdx := l.AbsIndex(-1)

	l.PushString("LEVELS")
	l.NewTable()
	levelsIdx := l.AbsIndex(-1)
	l.PushString("DEBUG")
	l.SetField(levelsIdx, "DEBUG")
	l.PushString("INFO")
	l.SetField(levelsIdx, "INFO")
	l.PushString("WARN")
	l.SetField(levelsIdx, "WARN")
	l.PushString("ERROR")
	l.SetField(levelsIdx, "ERROR")
	l.SetTable(moduleIdx)

	lua.NewMetaTable(l, slugLoggerHandle)
	l.PushValue(-1)
	l.SetField(-2, "__index")
	lua.SetFunctions(l, slogHandlerMethods, 0)
	l.Pop(1)

	return 1
}

func (lib *LibGoLogSlug) OpenLib(l *lua.State) {
	lua.Require(l, lib.Name, lib.Open, false)
	l.Pop(1)
}

func (lib *LibGoLogSlug) Require(l *lua.State) {
	l.Global("require")
	l.PushString(lib.Name)
	l.Call(1, 1)
}

// LuaNew is the equivalent to the lua `require(modname).new(logger)`
//
// Example (Go):
//
//	lib.LuaNew(l, myLogger)
//	// Lua stack now has a Logger at top, use l.Pop(1) to consume or store it
func (lib *LibGoLogSlug) LuaNew(l *lua.State, logger *slog.Logger) {
	lib.Require(l)
	l.Field(-1, "new")
	l.PushUserData(logger)
	l.Call(1, 1) // lib.new(logger)
}

var slogLoggerAnnotationsTmpl = template.Must(template.New("SlogLoggerAnnotations").Parse(`
---@meta {{.module}}

---@class {{.module}} : {{.Logger}}
---@field LEVELS {{.LogLevel}}
---@field default {{.Logger}}
local log = {}

---@param logger
---@return {{.Logger}}
function log.new(logger) end

---@enum {{.LogLevel}}
local LogLevel = {
    DEBUG = "DEBUG",
    INFO = "INFO",
    WARN = "WARN",
    ERROR = "ERROR",
}

---@alias {{.Level}}
---| '"DEBUG"' # Debug level
---| '"INFO"'  # Info level
---| '"WARN"'  # Warn level
---| '"ERROR"' # Error level

---@class {{.Logger}}
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

---@param level {{.Level}}
---@param msg string
---@param attrs? table
function Logger:log(level, msg, attrs) end

---@return {{.Level}}
function Logger:level() end

return log
`))

func (lib *LibGoLogSlug) Annotations() string {
	data := map[string]string{
		"module": lib.Name,
		"Logger": lib.Name + ".Logger",
		"LogLevel": lib.Name + ".LogLevel",
		"Level": lib.Name + ".Level",
	}
	var buf strings.Builder
	if err := slogLoggerAnnotationsTmpl.Execute(&buf, data); err != nil {
		return ""
	}
	return buf.String()
}
