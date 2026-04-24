package std

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/Shopify/go-lua"
	"github.com/ggallovalle/go-effectual"
)

type ModPath struct {
	name string
}

const ModPathName = "std.path"
const slugPathBufHandle = "go/std/path/PathBuf*"

type PathBuf struct {
	raw string
}

func (p *PathBuf) String() string {
	return p.raw
}

func (p *PathBuf) push(path string) {
	if p.raw == "" {
		p.raw = path
		return
	}
	if strings.HasSuffix(p.raw, "/") {
		p.raw += path
	} else {
		p.raw += "/" + path
	}
}

func (p *PathBuf) pop() bool {
	if p.raw == "" || p.raw == "/" {
		p.raw = ""
		return false
	}
	idx := strings.LastIndex(p.raw, "/")
	if idx <= 0 {
		p.raw = ""
		return false
	}
	p.raw = p.raw[:idx]
	return true
}

func (p *PathBuf) join(path string) *PathBuf {
	if strings.HasPrefix(path, "/") {
		return &PathBuf{raw: path}
	}
	newBuf := &PathBuf{raw: p.raw}
	newBuf.push(path)
	return newBuf
}

func (p *PathBuf) ends_with(child string) bool {
	if strings.HasPrefix(child, "/") {
		return p.raw == child
	}
	if strings.HasSuffix(p.raw, "/"+child) {
		return true
	}
	if strings.HasSuffix(p.raw, child) {
		idx := len(p.raw) - len(child) - 1
		if idx < 0 {
			return true
		}
		if p.raw[idx] == '/' {
			return true
		}
	}
	return false
}

func (p *PathBuf) starts_with(base string) bool {
	// First check exact match or match with trailing slashes in base
	if strings.HasPrefix(p.raw, base) {
		if len(p.raw) == len(base) {
			return true
		}
		if len(p.raw) > len(base) && p.raw[len(base)] == '/' {
			return true
		}
	}
	// Normalize base by removing trailing slashes and check again
	normalizedBase := strings.TrimRight(base, "/")
	if normalizedBase == "" {
		return false
	}
	if strings.HasPrefix(p.raw, normalizedBase) {
		if len(p.raw) == len(normalizedBase) {
			return true
		}
		if len(p.raw) > len(normalizedBase) && p.raw[len(normalizedBase)] == '/' {
			return true
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
	if result == "" {
		result = "/"
	}
	return &PathBuf{raw: result}, nil
}

func (p *PathBuf) with_extension(ext string) *PathBuf {
	dotIdx := strings.LastIndex(filepath.Base(p.raw), ".")
	if dotIdx <= 0 {
		p.raw += "." + ext
		return p
	}
	stem := p.raw[:len(p.raw)-len(filepath.Base(p.raw))+dotIdx]
	return &PathBuf{raw: stem + "." + ext}
}

func (p *PathBuf) with_file_name(name string) *PathBuf {
	dir := filepath.Dir(p.raw)
	if dir == "." {
		return &PathBuf{raw: name}
	}
	if !strings.HasSuffix(dir, "/") && dir != "/" {
		dir += "/"
	}
	return &PathBuf{raw: dir + name}
}

func (p *PathBuf) components() []string {
	if p.raw == "" {
		return nil
	}
	var parts []string
	trimmed := strings.Trim(p.raw, "/")
	if trimmed == "" {
		return []string{"/"}
	}
	for _, part := range strings.Split(trimmed, "/") {
		if part != "" {
			parts = append(parts, part)
		}
	}
	if strings.HasPrefix(p.raw, "/") {
		return append([]string{"/"}, parts...)
	}
	return parts
}

func (p *PathBuf) ancestors() []*PathBuf {
	var result []*PathBuf
	current := &PathBuf{raw: p.raw}
	for current.raw != "" && current.raw != "/" {
		result = append(result, current)
		parent := &PathBuf{raw: filepath.Dir(current.raw)}
		if parent.raw == current.raw {
			break
		}
		current = parent
	}
	if current.raw == "/" {
		result = append(result, current)
	}
	return result
}

func (p *PathBuf) parent() *PathBuf {
	if p.raw == "" || p.raw == "/" {
		return nil
	}
	dir := filepath.Dir(p.raw)
	if dir == "." {
		return nil
	}
	if dir == p.raw {
		return nil
	}
	if dir == "/" {
		return &PathBuf{raw: "/"}
	}
	return &PathBuf{raw: dir}
}

func (p *PathBuf) file_name() string {
	if p.raw == "" || p.raw == "/" {
		return ""
	}
	name := filepath.Base(p.raw)
	if name == "." || name == "/" {
		return ""
	}
	if strings.HasSuffix(p.raw, "/..") || name == ".." {
		return ""
	}
	return name
}

func (p *PathBuf) extension() string {
	name := p.file_name()
	if name == "" || name == "/" {
		return ""
	}
	dotIdx := strings.LastIndex(name, ".")
	if dotIdx <= 0 {
		return ""
	}
	ext := name[dotIdx+1:]
	if ext == "" || strings.Contains(ext, "/") {
		return ""
	}
	return ext
}

func (p *PathBuf) file_stem() string {
	name := p.file_name()
	if name == "" || name == "/" {
		return ""
	}
	dotIdx := strings.LastIndex(name, ".")
	if dotIdx <= 0 {
		return name
	}
	return name[:dotIdx]
}

func (p *PathBuf) has_root() bool {
	return strings.HasPrefix(p.raw, "/")
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

func pathBufFromString(s string) *PathBuf {
	trimmed := strings.Trim(s, "/")
	if trimmed == "" && strings.Contains(s, "/") {
		return &PathBuf{raw: "/"}
	}
	if s == "" {
		return &PathBuf{raw: ""}
	}
	if !strings.HasPrefix(s, "/") {
		return &PathBuf{raw: s}
	}
	return &PathBuf{raw: "/" + trimmed}
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
			pb = pathBufFromString(toPathBufString(l, 1))
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

var pathLibrary = []lua.RegistryFunction{
	{Name: "new", Function: func(l *lua.State) int {
		s := toPathBufString(l, 1)
		pathBufToLua(l, pathBufFromString(s))
		return 1
	}},
	{Name: "join", Function: func(l *lua.State) int {
		var parts []string
		n := l.Top()
		for i := 1; i <= n; i++ {
			parts = append(parts, toPathBufString(l, i))
		}
		if len(parts) == 0 {
			pathBufToLua(l, &PathBuf{raw: ""})
		} else {
			result := parts[0]
			for i := 1; i < len(parts); i++ {
				if strings.HasPrefix(parts[i], "/") {
					result = parts[i]
				} else if strings.HasSuffix(result, "/") {
					result += parts[i]
				} else {
					result += "/" + parts[i]
				}
			}
			pathBufToLua(l, pathBufFromString(result))
		}
		return 1
	}},
	{Name: "absolute", Function: func(l *lua.State) int {
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
		if strings.HasPrefix(s, "/") {
			pathBufToLua(l, pathBufFromString(s))
		} else {
			pathBufToLua(l, pathBufFromString(cwd+"/"+s))
		}
		return 1
	}},
}

func MakeModPath() effectual.LuaMod[struct{}] {
	return &ModPath{name: ModPathName}
}

func (lib *ModPath) Name() string {
	return lib.name
}

func (lib *ModPath) Annotations() string {
	return `---@meta std.path

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
local path = {}

---@type string
path.MAIN_SEPARATOR = "/"

---@param value string
---@return std.path.PathBuf
function path.new(value) end

---@param ... string|std.path.PathBuf
---@return std.path.PathBuf
function path.join(...) end

---@param path string|std.path.PathBuf
---@return std.path.PathBuf?
---@return string? Error message if error
function path.absolute(path) end

return path
`
}

func (lib *ModPath) Open(l *lua.State) int {
	lua.NewLibrary(l, pathLibrary)
	moduleIdx := l.AbsIndex(-1)

	l.PushString("MAIN_SEPARATOR")
	l.PushString("/")
	l.SetTable(moduleIdx)

	lua.NewMetaTable(l, slugPathBufHandle)
	l.PushValue(-1)
	l.SetField(-2, "__index")
	lua.SetFunctions(l, pathBufMetatable, 0)
	lua.SetFunctions(l, pathBufMethods, 0)
	l.Pop(1)

	return 1
}

func (lib *ModPath) OpenLib(l *lua.State) {
	lua.Require(l, lib.name, lib.Open, false)
	l.Pop(1)
}

func (lib *ModPath) Require(l *lua.State) {
	l.Global("require")
	l.PushString(lib.Name())
	l.Call(1, 1)
}

func (lib *ModPath) Api(l *lua.State) struct{} {
	return struct{}{}
}