package std

import (
	"strings"
	"text/template"

	"github.com/ggallovalle/go-effectual"
	"github.com/speedata/go-lua"
)

const (
	modUrlName    = "std.url"
	slugUrlHandle = "go/std/url/Url*"
)

type ModUrl struct {
	name string
}

type ModUrlApi struct {
	mod *ModUrl
	lua *lua.State
}

func (api *ModUrlApi) ToUrl(index int) (*Url, bool) {
	l := api.lua
	v := lua.CheckUserData(l, index, slugUrlHandle)
	if v != nil {
		if u, ok := v.(*Url); ok {
			return u, true
		}
	}
	return nil, false
}

func (api *ModUrlApi) CheckUrl(index int) *Url {
	if u, ok := api.ToUrl(index); ok {
		return u
	}
	lua.ArgumentError(api.lua, index, "expected std.url.Url")
	panic("unreachable")
}

func (lib *ModUrl) Name() string {
	return lib.name
}

func (lib *ModUrl) Open(l *lua.State) int {
	lua.NewLibrary(l, urlLibrary())

	lua.NewMetaTable(l, URL_HANDLE)
	l.PushValue(-1)
	l.SetField(-2, "__index")
	lua.SetFunctions(l, urlMetatable, 0)
	for name, fn := range urlMethods {
		l.PushGoFunction(fn)
		l.SetField(-2, name)
	}
	l.Pop(1)

	return 1
}

func (lib *ModUrl) OpenLib(l *lua.State) {
	lua.Require(l, lib.name, lib.Open, false)
	l.Pop(1)
}

func (lib *ModUrl) Require(l *lua.State) {
	l.Global("require")
	l.PushString(lib.Name())
	l.Call(1, 1)
}

func (lib *ModUrl) Api(l *lua.State) ModUrlApi {
	return ModUrlApi{mod: lib, lua: l}
}

func MakeModUrl() effectual.LuaMod[ModUrlApi] {
	return &ModUrl{name: modUrlName}
}

var urlAnnotationsTmpl = template.Must(template.New("UrlAnnotations").Parse(`---@meta {{.module}}

---@class (exact) {{.Url}} : userdata
---@operator div({{.Url}}|string): {{.Url}}
---@field scheme string|nil
---@field host string|nil
---@field port integer|nil
---@field port_inferred integer
---@field username string|nil
---@field password string|nil
---@field path {{.Path}}
---@field query {{.Query}}
---@field fragment string|nil
local Url = {}

---@param path string
function Url:add_query(path) end

---@param key string
function Url:remove_query(key) end

local url = {}

---@return {{.Url}}
function url.new() end

---@param raw string
---@return {{.Url}}
function url.deserialize(raw) end

---@param u {{.Url}}
---@return string
function url.serialize(u) end

return url
`))

func (lib *ModUrl) Annotations() string {
	data := map[string]string{
		"module": lib.name,
		"Url":    lib.name + ".Url",
		"Path":   "std.path.Path",
		"Query":  "std.serde.query.Query",
	}
	var buf strings.Builder
	if err := urlAnnotationsTmpl.Execute(&buf, data); err != nil {
		return ""
	}
	return buf.String()
}
