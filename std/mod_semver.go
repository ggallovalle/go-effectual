package std

import (
	"strings"
	"text/template"

	"github.com/blang/semver/v4"
	"github.com/speedata/go-lua"
	"github.com/ggallovalle/go-effectual"
)

type ModSemver struct {
	name string
}

type ModSemverApi struct {
	mod *ModSemver
	lua *lua.State
}

const (
	ModSemverName      = "std.semver"
	slugVersionHandle  = "go/std/semver/Version*"
	slugRangeHandle   = "go/std/semver/Range*"
)

func MakeModSemver() effectual.LuaMod[ModSemverApi] {
	return &ModSemver{name: ModSemverName}
}

type RangeWrapper struct {
	r semver.Range
}

func (r *RangeWrapper) Contains(v semver.Version) bool {
	return r.r(v)
}

func (r *RangeWrapper) And(other *RangeWrapper) *RangeWrapper {
	return &RangeWrapper{r: r.r.AND(other.r)}
}

func (r *RangeWrapper) Or(other *RangeWrapper) *RangeWrapper {
	return &RangeWrapper{r: r.r.OR(other.r)}
}

func rangeWrapperToLua(l *lua.State, rw *RangeWrapper) {
	l.PushUserData(rw)
	lua.SetMetaTableNamed(l, slugRangeHandle)
}

func toRangeWrapper(l *lua.State, idx int) *RangeWrapper {
	ud := lua.CheckUserData(l, idx, slugRangeHandle)
	v, ok := ud.(*RangeWrapper)
	if !ok {
		panic("not a *RangeWrapper")
	}
	return v
}

func versionToLua(l *lua.State, v *semver.Version) {
	l.PushUserData(v)
	lua.SetMetaTableNamed(l, slugVersionHandle)
}

func toVersion(l *lua.State, idx int) *semver.Version {
	ud := lua.CheckUserData(l, idx, slugVersionHandle)
	v, ok := ud.(*semver.Version)
	if !ok {
		panic("not a *semver.Version")
	}
	return v
}

var versionGetters = map[string]func(*lua.State){
	"major": func(l *lua.State) {
		v := toVersion(l, 1)
		l.PushInteger(int(v.Major))
	},
	"minor": func(l *lua.State) {
		v := toVersion(l, 1)
		l.PushInteger(int(v.Minor))
	},
	"patch": func(l *lua.State) {
		v := toVersion(l, 1)
		l.PushInteger(int(v.Patch))
	},
}

var versionMetatable = []lua.RegistryFunction{
	{Name: "__tostring", Function: func(l *lua.State) int {
		v := toVersion(l, 1)
		l.PushString(v.String())
		return 1
	}},
	{Name: "__lt", Function: func(l *lua.State) int {
		v1 := toVersion(l, 1)
		v2 := toVersion(l, 2)
		l.PushBoolean(v1.LT(*v2))
		return 1
	}},
	{Name: "__le", Function: func(l *lua.State) int {
		v1 := toVersion(l, 1)
		v2 := toVersion(l, 2)
		l.PushBoolean(v1.LTE(*v2))
		return 1
	}},
	{Name: "__eq", Function: func(l *lua.State) int {
		v1 := toVersion(l, 1)
		v2 := toVersion(l, 2)
		l.PushBoolean(v1.EQ(*v2))
		return 1
	}},
	effectual.LuaMetaIndex(versionGetters, nil),
}

var rangeMetatable = []lua.RegistryFunction{
	{Name: "__index", Function: func(l *lua.State) int {
		key := lua.CheckString(l, 2)
		switch key {
		case "intersect":
			l.PushGoFunction(func(l *lua.State) int {
				rw1 := toRangeWrapper(l, 1)
				rw2 := toRangeWrapper(l, 2)
				rangeWrapperToLua(l, rw1.And(rw2))
				return 1
			})
			return 1
		case "union":
			l.PushGoFunction(func(l *lua.State) int {
				rw1 := toRangeWrapper(l, 1)
				rw2 := toRangeWrapper(l, 2)
				rangeWrapperToLua(l, rw1.Or(rw2))
				return 1
			})
			return 1
		case "contains":
			l.PushGoFunction(func(l *lua.State) int {
				rw := toRangeWrapper(l, 1)
				v := toVersion(l, 2)
				l.PushBoolean(rw.Contains(*v))
				return 1
			})
			return 1
		}
		l.PushNil()
		return 1
	}},
}

var semverLibrary = []lua.RegistryFunction{
	{Name: "new", Function: func(l *lua.State) int {
		s, _ := l.ToString(1)
		v, err := semver.New(s)
		if err != nil {
			l.PushNil()
			l.PushString(err.Error())
			return 2
		}
		versionToLua(l, v)
		return 1
	}},
	{Name: "range_new", Function: func(l *lua.State) int {
		s, _ := l.ToString(1)
		r, err := semver.ParseRange(s)
		if err != nil {
			l.PushNil()
			l.PushString(err.Error())
			return 2
		}
		rw := &RangeWrapper{r: r}
		rangeWrapperToLua(l, rw)
		return 1
	}},
}

func (lib *ModSemver) Name() string {
	return lib.name
}

func (lib *ModSemver) Open(l *lua.State) int {
	lua.NewLibrary(l, semverLibrary)

	lua.NewMetaTable(l, slugVersionHandle)
	l.PushValue(-1)
	l.SetField(-2, "__index")
	lua.SetFunctions(l, versionMetatable, 0)
	l.Pop(1)

	lua.NewMetaTable(l, slugRangeHandle)
	l.PushValue(-1)
	l.SetField(-2, "__index")
	lua.SetFunctions(l, rangeMetatable, 0)
	l.Pop(1)

	return 1
}

func (lib *ModSemver) OpenLib(l *lua.State) {
	lua.Require(l, lib.name, lib.Open, false)
	l.Pop(1)
}

func (lib *ModSemver) Require(l *lua.State) {
	l.Global("require")
	l.PushString(lib.Name())
	l.Call(1, 1)
}

func (lib *ModSemver) Api(l *lua.State) ModSemverApi {
	return ModSemverApi{mod: lib, lua: l}
}

var semverAnnotationsTmpl = template.Must(template.New("SemverAnnotations").Parse(`---@meta {{.module}}

---@class (exact) {{.Version}} : userdata
---@field major integer
---@field minor integer
---@field patch integer
---@operator lt({{.Version}}): boolean
---@operator le({{.Version}}): boolean
---@operator eq({{.Version}}): boolean
local Version = {}

---@class (exact) {{.Range}} : userdata
local Range = {}

--- Performs logical intersection on two ranges
---@param range {{.Range}}
---@return {{.Range}}
function Range:intersect(range) end

--- Performs logical union on two ranges
---@param range {{.Range}}
---@return {{.Range}}
function Range:union(range) end

--- Checks if the range contains the given version
---@param version {{.Version}}
---@return boolean
function Range:contains(version) end

local semver = {}

--- Creates a new Version from a string
---@param version string (e.g., "1.2.3")
---@return {{.Version}}?
---@return string? Error message if version string is invalid
function semver.new(version) end

--- Creates a new Range from a semver range string
---@param range string (e.g., ">=1.0.0 <2.0.0")
---@return {{.Range}}?
---@return string? Error message if range string is invalid
function semver.range_new(range) end

return semver
`))

func (lib *ModSemver) Annotations() string {
	data := map[string]string{
		"module": lib.name,
		"Version": lib.name + ".Version",
		"Range":  lib.name + ".Range",
	}
	var buf strings.Builder
	if err := semverAnnotationsTmpl.Execute(&buf, data); err != nil {
		return ""
	}
	return buf.String()
}
