package std

import (
	"os"
	"runtime"
	"strings"
	"text/template"

	"github.com/speedata/go-lua"
	"github.com/ggallovalle/go-effectual"
)

const (
	posixSep          = "/"
	winSep            = "\\"
	slugPathHandle = "go/std/path/Path*"
)

func nativeSep() string {
	if runtime.GOOS == "windows" {
		return winSep
	}
	return posixSep
}

func altSep(sep string) string {
	if sep == posixSep {
		return winSep
	}
	return posixSep
}

type pathMod struct {
	name   string
	sep    string
	altSep string
}

type ModPathApi struct {
	mod *pathMod
	lua *lua.State
}

func (api *ModPathApi) New(pathStr string) {
	lib := api.mod
	l := api.lua

	lib.Require(l)
	l.Field(-1, "new")
	l.PushString(pathStr)
	l.Call(1, 1)
}

func (api *ModPathApi) ToPath(index int) (*Path, bool) {
	l := api.lua

	v := lua.CheckUserData(l, index, slugPathHandle)
	if v != nil {
		if pb, ok := v.(*Path); ok {
			return pb, true
		}
	}

	if s, ok := l.ToString(index); ok {
		return pathFromStringSep(s, api.mod.sep), true
	}

	return nil, false
}

func (api *ModPathApi) CheckPath(index int) *Path {
	if pb, ok := api.ToPath(index); ok {
		return pb
	}
	lua.ArgumentError(api.lua, index, "expected std.path.Path or string")
	panic("unreachable")
}

var pathMethods = map[string]lua.Function{
	"push": func(l *lua.State) int {
		pb := toPath(l, 1)
		arg := toPathString(l, 2)
		pb.Push(arg)
		return 0
	},
	"pop": func(l *lua.State) int {
		pb := toPath(l, 1)
		l.PushBoolean(pb.Pop())
		return 1
	},
	"join": func(l *lua.State) int {
		pb := toPath(l, 1)
		arg := toPathString(l, 2)
		pathToLua(l, pb.Join(arg))
		return 1
	},
	"ends_with": func(l *lua.State) int {
		pb := toPath(l, 1)
		child := toPathString(l, 2)
		l.PushBoolean(pb.EndsWith(child))
		return 1
	},
	"starts_with": func(l *lua.State) int {
		pb := toPath(l, 1)
		base := toPathString(l, 2)
		l.PushBoolean(pb.StartsWith(base))
		return 1
	},
	"strip_prefix": func(l *lua.State) int {
		pb := toPath(l, 1)
		prefix := toPathString(l, 2)
		result, err := pb.StripPrefix(prefix)
		if err != nil {
			l.PushNil()
			l.PushString(err.Error())
			return 2
		}
		pathToLua(l, result)
		return 1
	},
	"with_extension": func(l *lua.State) int {
		pb := toPath(l, 1)
		ext, _ := l.ToString(2)
		pathToLua(l, pb.WithExtension(ext))
		return 1
	},
	"with_file_name": func(l *lua.State) int {
		pb := toPath(l, 1)
		name, _ := l.ToString(2)
		pathToLua(l, pb.WithFileName(name))
		return 1
	},
}

var pathGetters = map[string]func(*lua.State){
	"parent": func(l *lua.State) {
		pb := toPath(l, 1)
		parent := pb.Parent()
		if parent == nil {
			l.PushNil()
		} else {
			pathToLua(l, parent)
		}
	},
	"components": func(l *lua.State) {
		pb := toPath(l, 1)
		comps := pb.Components()
		l.NewTable()
		for i, c := range comps {
			l.PushInteger(i + 1)
			l.PushString(c)
			l.SetTable(-3)
		}
	},
	"ancestors": func(l *lua.State) {
		pb := toPath(l, 1)
		ancs := pb.Ancestors()
		l.NewTable()
		for i, a := range ancs {
			l.PushInteger(i + 1)
			pathToLua(l, a)
			l.SetTable(-3)
		}
	},
	"file_name": func(l *lua.State) {
		pb := toPath(l, 1)
		name := pb.FileName()
		if name == "" {
			l.PushNil()
		} else {
			l.PushString(name)
		}
	},
	"extension": func(l *lua.State) {
		pb := toPath(l, 1)
		ext := pb.Extension()
		if ext == "" {
			l.PushNil()
		} else {
			l.PushString(ext)
		}
	},
	"file_stem": func(l *lua.State) {
		pb := toPath(l, 1)
		stem := pb.FileStem()
		if stem == "" {
			l.PushNil()
		} else {
			l.PushString(stem)
		}
	},
	"has_root": func(l *lua.State) {
		pb := toPath(l, 1)
		l.PushBoolean(pb.HasRoot())
	},
	"is_absolute": func(l *lua.State) {
		pb := toPath(l, 1)
		l.PushBoolean(pb.IsAbsolute())
	},
	"is_relative": func(l *lua.State) {
		pb := toPath(l, 1)
		l.PushBoolean(pb.IsRelative())
	},
}

var pathMetatable = []lua.RegistryFunction{
	{Name: "__tostring", Function: func(l *lua.State) int {
		pb := toPath(l, 1)
		l.PushString(pb.raw)
		return 1
	}},
	{Name: "__div", Function: func(l *lua.State) int {
		var pb *Path
		if ud, ok := l.ToUserData(1).(*Path); ok {
			pb = ud
		} else {
			pb = pathFromStringSep(toPathString(l, 1), posixSep)
		}
		arg := toPathString(l, 2)
		pathToLua(l, pb.Join(arg))
		return 1
	}},
	{Name: "__concat", Function: func(l *lua.State) int {
		left := toPathString(l, 1)
		right := toPathString(l, 2)
		l.PushString(left + right)
		return 1
	}},
	effectual.LuaMetaIndex(pathGetters, pathMethods),
}

func pathNew(sep string, l *lua.State) int {
	s := toPathString(l, 1)
	pathToLua(l, pathFromStringSep(s, sep))
	return 1
}

func pathJoin(sep string, l *lua.State) int {
	alt := altSep(sep)
	var parts []string
	n := l.Top()
	for i := 1; i <= n; i++ {
		parts = append(parts, toPathString(l, i))
	}
	if len(parts) == 0 {
		pathToLua(l, &Path{raw: "", sep: sep})
	} else {
		result := parts[0]
		for i := 1; i < len(parts); i++ {
			if strings.HasPrefix(parts[i], sep) || strings.HasPrefix(parts[i], alt) {
				result = parts[i]
			} else if strings.HasSuffix(result, sep) || strings.HasSuffix(result, alt) {
				result += parts[i]
			} else {
				result += sep + parts[i]
			}
		}
		pathToLua(l, pathFromStringSep(result, sep))
	}
	return 1
}

func pathAbsolute(sep string, l *lua.State) int {
	s := toPathString(l, 1)
	if s == "" {
		l.PushNil()
		l.PushString("path is empty")
		return 2
	}
	cwd, err := os.Getwd()
	if err != nil {
		l.PushNil()
		l.PushString("failed to get current directory")
		return 2
	}
	if strings.HasPrefix(s, sep) || strings.HasPrefix(s, altSep(sep)) {
		pathToLua(l, pathFromStringSep(s, sep))
	} else {
		pathToLua(l, pathFromStringSep(cwd+sep+s, sep))
	}
	return 1
}

func pathLibrary(sep string) []lua.RegistryFunction {
	return []lua.RegistryFunction{
		{Name: "new", Function: func(l *lua.State) int {
			return pathNew(sep, l)
		}},
		{Name: "join", Function: func(l *lua.State) int {
			return pathJoin(sep, l)
		}},
		{Name: "absolute", Function: func(l *lua.State) int {
			return pathAbsolute(sep, l)
		}},
	}
}

var pathAnnotationsTmpl = template.Must(template.New("PathAnnotations").Parse(`---@meta {{.module}}

---@class (exact) {{.Path}} : userdata
---@operator div({{.Path}}|string): {{.Path}}
---@operator concat({{.Path}}|string): string
---@field parent {{.Path}}|nil
---@field components string[]
---@field ancestors {{.Path}}[]
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
---@return {{.Path}}
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
---@return {{.Path}}?
---@return string? Error message if prefix is not found
function Path:strip_prefix(prefix) end

--- Sets the file extension, returning a new Path with the changed extension
---@param ext string
---@return {{.Path}}
function Path:with_extension(ext) end

--- Sets the file name component, returning a new Path
---@param name string
---@return {{.Path}}
function Path:with_file_name(name) end

local path = {}

---@type string
path.MAIN_SEPARATOR = "{{.sep}}"

--- Creates a new Path from the given path string
---@param value string
---@return {{.Path}}
function path.new(value) end

--- Joins multiple path segments together, returning a new Path
---@param ... string|{{.Path}}
---@return {{.Path}}
function path.join(...) end

--- Converts the given path to an absolute path based on the current working directory
---@param path string|{{.Path}}
---@return {{.Path}}?
---@return string? Error message if path is empty
function path.absolute(path) end

---@class {{.module}}.posix : {{.module}}
local posix = {}

---@class {{.module}}.win32 : {{.module}}
local win32 = {}

return path
`))

func (lib *pathMod) Name() string {
	return lib.name
}

func (lib *pathMod) Annotations() string {
	data := map[string]string{
		"module": lib.name,
		"Path":   lib.name + ".Path",
		"sep":    lib.sep,
	}
	var buf strings.Builder
	if err := pathAnnotationsTmpl.Execute(&buf, data); err != nil {
		return ""
	}
	return buf.String()
}

func (lib *pathMod) Open(l *lua.State) int {
	lua.NewLibrary(l, pathLibrary(lib.sep))
	moduleIdx := l.AbsIndex(-1)

	l.PushString("MAIN_SEPARATOR")
	l.PushString(lib.sep)
	l.SetTable(moduleIdx)

	l.PushString("posix")
	lua.NewLibrary(l, pathLibrary(posixSep))
	posixIdx := l.AbsIndex(-1)
	l.PushString("MAIN_SEPARATOR")
	l.PushString(posixSep)
	l.SetTable(posixIdx)
	l.SetTable(moduleIdx)

	l.PushString("win32")
	lua.NewLibrary(l, pathLibrary(winSep))
	winIdx := l.AbsIndex(-1)
	l.PushString("MAIN_SEPARATOR")
	l.PushString(winSep)
	l.SetTable(winIdx)
	l.SetTable(moduleIdx)

	lua.NewMetaTable(l, slugPathHandle)
	l.PushValue(-1)
	l.SetField(-2, "__index")
	lua.SetFunctions(l, pathMetatable, 0)
	for name, fn := range pathMethods {
		l.PushGoFunction(fn)
		l.SetField(-2, name)
	}
	l.Pop(1)

	return 1
}

func (lib *pathMod) OpenLib(l *lua.State) {
	lua.Require(l, lib.name, lib.Open, false)
	l.Pop(1)
}

func (lib *pathMod) Require(l *lua.State) {
	l.Global("require")
	l.PushString(lib.Name())
	l.Call(1, 1)
}

func (lib *pathMod) Api(l *lua.State) ModPathApi {
	return ModPathApi{mod: lib, lua: l}
}

func MakeModPath() effectual.LuaMod[ModPathApi] {
	sep := nativeSep()
	return &pathMod{name: "std.path", sep: sep, altSep: altSep(sep)}
}