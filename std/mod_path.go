package std

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/Shopify/go-lua"
	"github.com/ggallovalle/go-effectual"
)

const (
	posixSep          = "/"
	winSep            = "\\"
	slugPathBufHandle = "go/std/path/PathBuf*"
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

type PathBuf struct {
	raw string
	sep string
}

func (p *PathBuf) String() string {
	return p.raw
}

func (p *PathBuf) dir() string {
	alt := altSep(p.sep)
	trimmed := strings.TrimRight(p.raw, p.sep+alt)
	if trimmed == "" {
		return p.sep
	}
	idx := strings.LastIndex(trimmed, p.sep)
	if idx == -1 {
		idx = strings.LastIndex(trimmed, alt)
	}
	if idx == 0 {
		return p.sep
	}
	if idx <= 0 {
		return "."
	}
	return trimmed[:idx]
}

func (p *PathBuf) push(path string) {
	if p.raw == "" {
		p.raw = path
		return
	}
	if strings.HasSuffix(p.raw, p.sep) {
		p.raw += path
	} else {
		p.raw += p.sep + path
	}
}

func (p *PathBuf) pop() bool {
	if p.raw == "" || p.raw == p.sep {
		p.raw = ""
		return false
	}
	idx := strings.LastIndex(p.raw, p.sep)
	if idx <= 0 {
		p.raw = ""
		return false
	}
	p.raw = p.raw[:idx]
	return true
}

func (p *PathBuf) join(path string) *PathBuf {
	if strings.HasPrefix(path, p.sep) || strings.HasPrefix(path, altSep(p.sep)) {
		return &PathBuf{raw: path, sep: p.sep}
	}
	newBuf := &PathBuf{raw: p.raw, sep: p.sep}
	newBuf.push(path)
	return newBuf
}

func (p *PathBuf) ends_with(child string) bool {
	alt := altSep(p.sep)
	if strings.HasPrefix(child, p.sep) || strings.HasPrefix(child, alt) {
		return p.raw == child
	}
	if strings.HasSuffix(p.raw, p.sep+child) || strings.HasSuffix(p.raw, alt+child) {
		return true
	}
	if strings.HasSuffix(p.raw, child) {
		idx := len(p.raw) - len(child) - 1
		if idx < 0 {
			return true
		}
		c := p.raw[idx]
		return c == p.sep[0] || c == alt[0]
	}
	return false
}

func (p *PathBuf) starts_with(base string) bool {
	if strings.HasPrefix(p.raw, base) {
		if len(p.raw) == len(base) {
			return true
		}
		if len(p.raw) > len(base) {
			c := p.raw[len(base)]
			if c == p.sep[0] || c == altSep(p.sep)[0] {
				return true
			}
		}
	}
	normalizedBase := strings.TrimRight(base, "/\\")
	if normalizedBase == "" {
		return false
	}
	if strings.HasPrefix(p.raw, normalizedBase) {
		if len(p.raw) == len(normalizedBase) {
			return true
		}
		if len(p.raw) > len(normalizedBase) {
			c := p.raw[len(normalizedBase)]
			if c == p.sep[0] || c == altSep(p.sep)[0] {
				return true
			}
		}
	}
	return false
}

func (p *PathBuf) strip_prefix(prefix string) (*PathBuf, error) {
	if !strings.HasPrefix(p.raw, prefix) {
		return nil, errors.New("prefix not found")
	}
	result := strings.TrimPrefix(p.raw, prefix)
	result = strings.TrimPrefix(result, "/")
	result = strings.TrimPrefix(result, "\\")
	if result == "" {
		result = p.sep
	}
	return &PathBuf{raw: result, sep: p.sep}, nil
}

func (p *PathBuf) with_extension(ext string) *PathBuf {
	base := filepath.Base(p.raw)
	dotIdx := strings.LastIndex(base, ".")
	if dotIdx <= 0 {
		return &PathBuf{raw: p.raw + "." + ext, sep: p.sep}
	}
	stem := p.raw[:len(p.raw)-len(base)+dotIdx]
	return &PathBuf{raw: stem + "." + ext, sep: p.sep}
}

func (p *PathBuf) with_file_name(name string) *PathBuf {
	dir := p.dir()
	if dir == "." {
		return &PathBuf{raw: name, sep: p.sep}
	}
	if dir != p.sep && !strings.HasSuffix(dir, p.sep) && !strings.HasSuffix(dir, altSep(p.sep)) {
		dir += p.sep
	}
	return &PathBuf{raw: dir + name, sep: p.sep}
}

func (p *PathBuf) components() []string {
	if p.raw == "" {
		return nil
	}
	var parts []string
	trimmed := strings.Trim(p.raw, "/\\")
	if trimmed == "" {
		return []string{p.sep}
	}
	for _, part := range strings.Split(trimmed, p.sep) {
		if part != "" {
			parts = append(parts, part)
		}
	}
	if strings.HasPrefix(p.raw, p.sep) {
		return append([]string{p.sep}, parts...)
	}
	if strings.HasPrefix(p.raw, altSep(p.sep)) {
		return append([]string{altSep(p.sep)}, parts...)
	}
	return parts
}

func (p *PathBuf) ancestors() []*PathBuf {
	var result []*PathBuf
	current := &PathBuf{raw: p.raw, sep: p.sep}
	for current.raw != "" && current.raw != p.sep && current.raw != altSep(p.sep) {
		result = append(result, current)
		parent := &PathBuf{raw: current.dir(), sep: p.sep}
		if parent.raw == current.raw {
			break
		}
		current = parent
	}
	if current.raw == p.sep || current.raw == altSep(p.sep) {
		result = append(result, current)
	}
	return result
}

func (p *PathBuf) parent() *PathBuf {
	if p.raw == "" || p.raw == p.sep || p.raw == altSep(p.sep) {
		return nil
	}
	dir := p.dir()
	if dir == "." {
		return nil
	}
	if dir == p.raw {
		return nil
	}
	if dir == p.sep || dir == altSep(p.sep) {
		return &PathBuf{raw: p.sep, sep: p.sep}
	}
	return &PathBuf{raw: dir, sep: p.sep}
}

func (p *PathBuf) baseName() string {
	trimmed := strings.TrimRight(p.raw, p.sep+altSep(p.sep))
	if trimmed == "" {
		return ""
	}
	idx := strings.LastIndex(trimmed, p.sep)
	if idx == -1 {
		alt := altSep(p.sep)
		idx = strings.LastIndex(trimmed, alt)
	}
	if idx < 0 {
		return trimmed
	}
	return trimmed[idx+1:]
}

func (p *PathBuf) file_name() string {
	name := p.baseName()
	if name == "" || name == ".." {
		return ""
	}
	return name
}

func (p *PathBuf) extension() string {
	name := p.file_name()
	if name == "" {
		return ""
	}
	dotIdx := strings.LastIndex(name, ".")
	if dotIdx <= 0 {
		return ""
	}
	ext := name[dotIdx+1:]
	if ext == "" || strings.Contains(ext, "/") || strings.Contains(ext, "\\") {
		return ""
	}
	return ext
}

func (p *PathBuf) file_stem() string {
	name := p.file_name()
	if name == "" {
		return ""
	}
	dotIdx := strings.LastIndex(name, ".")
	if dotIdx <= 0 {
		return name
	}
	return name[:dotIdx]
}

func (p *PathBuf) has_root() bool {
	return strings.HasPrefix(p.raw, p.sep) || strings.HasPrefix(p.raw, altSep(p.sep))
}

func (p *PathBuf) is_absolute() bool {
	return p.has_root()
}

func (p *PathBuf) is_relative() bool {
	return !p.has_root()
}

func toPathBufString(l *lua.State, idx int) string {
	if l.IsUserData(idx) {
		v := l.ToUserData(idx)
		if v == nil {
			if s, ok := l.ToString(idx); ok {
				return s
			}
			return ""
		}
		switch x := v.(type) {
		case *PathBuf:
			return x.raw
		default:
			if s, ok := l.ToString(idx); ok {
				return s
			}
			return ""
		}
	}
	if l.IsString(idx) {
		s, _ := l.ToString(idx)
		return s
	}
	return ""
}

func pathBufToLua(l *lua.State, p *PathBuf) {
	l.PushUserData(p)
	lua.SetMetaTableNamed(l, slugPathBufHandle)
}

func toPathBuf(l *lua.State, idx int) *PathBuf {
	return lua.CheckUserData(l, idx, slugPathBufHandle).(*PathBuf)
}

func pathBufFromStringSep(s, sep string) *PathBuf {
	alt := altSep(sep)
	trimmed := strings.Trim(s, "/\\")
	if trimmed == "" && strings.ContainsAny(s, "/\\") {
		return &PathBuf{raw: sep, sep: sep}
	}
	if s == "" {
		return &PathBuf{raw: "", sep: sep}
	}
	if !strings.HasPrefix(s, sep) && !strings.HasPrefix(s, alt) {
		return &PathBuf{raw: s, sep: sep}
	}
	return &PathBuf{raw: sep + trimmed, sep: sep}
}

var pathBufMethods = []lua.RegistryFunction{
	{Name: "push", Function: func(l *lua.State) int {
		pb := toPathBuf(l, 1)
		arg := toPathBufString(l, 2)
		pb.push(arg)
		return 0
	}},
	{Name: "pop", Function: func(l *lua.State) int {
		pb := toPathBuf(l, 1)
		l.PushBoolean(pb.pop())
		return 1
	}},
	{Name: "join", Function: func(l *lua.State) int {
		pb := toPathBuf(l, 1)
		arg := toPathBufString(l, 2)
		pathBufToLua(l, pb.join(arg))
		return 1
	}},
	{Name: "ends_with", Function: func(l *lua.State) int {
		pb := toPathBuf(l, 1)
		child := toPathBufString(l, 2)
		l.PushBoolean(pb.ends_with(child))
		return 1
	}},
	{Name: "starts_with", Function: func(l *lua.State) int {
		pb := toPathBuf(l, 1)
		base := toPathBufString(l, 2)
		l.PushBoolean(pb.starts_with(base))
		return 1
	}},
	{Name: "strip_prefix", Function: func(l *lua.State) int {
		pb := toPathBuf(l, 1)
		prefix := toPathBufString(l, 2)
		result, err := pb.strip_prefix(prefix)
		if err != nil {
			l.PushNil()
			l.PushString(err.Error())
			return 2
		}
		pathBufToLua(l, result)
		return 1
	}},
	{Name: "with_extension", Function: func(l *lua.State) int {
		pb := toPathBuf(l, 1)
		ext, _ := l.ToString(2)
		pathBufToLua(l, pb.with_extension(ext))
		return 1
	}},
	{Name: "with_file_name", Function: func(l *lua.State) int {
		pb := toPathBuf(l, 1)
		name, _ := l.ToString(2)
		pathBufToLua(l, pb.with_file_name(name))
		return 1
	}},
}

var pathBufGetters = map[string]func(*lua.State){
	"parent": func(l *lua.State) {
		pb := toPathBuf(l, 1)
		parent := pb.parent()
		if parent == nil {
			l.PushNil()
		} else {
			pathBufToLua(l, parent)
		}
	},
	"components": func(l *lua.State) {
		pb := toPathBuf(l, 1)
		comps := pb.components()
		l.NewTable()
		for i, c := range comps {
			l.PushInteger(i + 1)
			l.PushString(c)
			l.SetTable(-3)
		}
	},
	"ancestors": func(l *lua.State) {
		pb := toPathBuf(l, 1)
		ancs := pb.ancestors()
		l.NewTable()
		for i, a := range ancs {
			l.PushInteger(i + 1)
			pathBufToLua(l, a)
			l.SetTable(-3)
		}
	},
	"file_name": func(l *lua.State) {
		pb := toPathBuf(l, 1)
		name := pb.file_name()
		if name == "" {
			l.PushNil()
		} else {
			l.PushString(name)
		}
	},
	"extension": func(l *lua.State) {
		pb := toPathBuf(l, 1)
		ext := pb.extension()
		if ext == "" {
			l.PushNil()
		} else {
			l.PushString(ext)
		}
	},
	"file_stem": func(l *lua.State) {
		pb := toPathBuf(l, 1)
		stem := pb.file_stem()
		if stem == "" {
			l.PushNil()
		} else {
			l.PushString(stem)
		}
	},
	"has_root": func(l *lua.State) {
		pb := toPathBuf(l, 1)
		l.PushBoolean(pb.has_root())
	},
	"is_absolute": func(l *lua.State) {
		pb := toPathBuf(l, 1)
		l.PushBoolean(pb.is_absolute())
	},
	"is_relative": func(l *lua.State) {
		pb := toPathBuf(l, 1)
		l.PushBoolean(pb.is_relative())
	},
}

var pathBufMetatable = []lua.RegistryFunction{
	{Name: "__tostring", Function: func(l *lua.State) int {
		pb := toPathBuf(l, 1)
		l.PushString(pb.raw)
		return 1
	}},
	{Name: "__div", Function: func(l *lua.State) int {
		var pb *PathBuf
		if ud, ok := l.ToUserData(1).(*PathBuf); ok {
			pb = ud
		} else {
			pb = pathBufFromStringSep(toPathBufString(l, 1), posixSep)
		}
		arg := toPathBufString(l, 2)
		pathBufToLua(l, pb.join(arg))
		return 1
	}},
	{Name: "__concat", Function: func(l *lua.State) int {
		left := toPathBufString(l, 1)
		right := toPathBufString(l, 2)
		l.PushString(left + right)
		return 1
	}},
	{Name: "__index", Function: func(l *lua.State) int {
		key := lua.CheckString(l, 2)
		if l.MetaTable(1) {
			l.Field(-1, key)
			if !l.IsNil(-1) {
				return 1
			}
			l.Pop(1)
		}
		if getter, ok := pathBufGetters[key]; ok {
			getter(l)
			return 1
		}
		l.PushNil()
		return 1
	}},
}

func pathNew(sep string, l *lua.State) int {
	s := toPathBufString(l, 1)
	pathBufToLua(l, pathBufFromStringSep(s, sep))
	return 1
}

func pathJoin(sep string, l *lua.State) int {
	alt := altSep(sep)
	var parts []string
	n := l.Top()
	for i := 1; i <= n; i++ {
		parts = append(parts, toPathBufString(l, i))
	}
	if len(parts) == 0 {
		pathBufToLua(l, &PathBuf{raw: "", sep: sep})
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
		pathBufToLua(l, pathBufFromStringSep(result, sep))
	}
	return 1
}

func pathAbsolute(sep string, l *lua.State) int {
	s := toPathBufString(l, 1)
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
		pathBufToLua(l, pathBufFromStringSep(s, sep))
	} else {
		pathBufToLua(l, pathBufFromStringSep(cwd+sep+s, sep))
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

---@class (exact) {{.PathBuf}} : userdata
---@operator div({{.PathBuf}}|string): {{.PathBuf}}
---@operator concat({{.PathBuf}}|string): string
---@field parent {{.PathBuf}}|nil
---@field components string[]
---@field ancestors {{.PathBuf}}[]
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
---@return {{.PathBuf}}
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
---@return {{.PathBuf}}?
---@return string? Error message if prefix is not found
function PathBuf:strip_prefix(prefix) end

--- Sets the file extension, returning a new PathBuf with the changed extension
---@param ext string
---@return {{.PathBuf}}
function PathBuf:with_extension(ext) end

--- Sets the file name component, returning a new PathBuf
---@param name string
---@return {{.PathBuf}}
function PathBuf:with_file_name(name) end

local path = {}

---@type string
path.MAIN_SEPARATOR = "{{.sep}}"

--- Creates a new PathBuf from the given path string
---@param value string
---@return {{.PathBuf}}
function path.new(value) end

--- Joins multiple path segments together, returning a new PathBuf
---@param ... string|{{.PathBuf}}
---@return {{.PathBuf}}
function path.join(...) end

--- Converts the given path to an absolute path based on the current working directory
---@param path string|{{.PathBuf}}
---@return {{.PathBuf}}?
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
		"module":  lib.name,
		"PathBuf": lib.name + ".PathBuf",
		"sep":     lib.sep,
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

	lua.NewMetaTable(l, slugPathBufHandle)
	l.PushValue(-1)
	l.SetField(-2, "__index")
	lua.SetFunctions(l, pathBufMetatable, 0)
	lua.SetFunctions(l, pathBufMethods, 0)
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