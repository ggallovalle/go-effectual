package std

import (
	"net/url"
	"strconv"
	"strings"
	"text/template"

	"github.com/ggallovalle/go-effectual"
	"github.com/ggallovalle/go-effectual/std/serde"
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

type Url struct {
	raw          string
	scheme       string
	host         string
	port         *int
	portInferred int
	username     *string
	password     *string
	path         *Path
	query        *serde.Query
	fragment     *string
}

func (u *Url) String() string {
	return u.raw
}

func (u *Url) AddQuery(key, value string) {
	u.query.Append(key, value)
	rebuildUrl(u)
}

func (u *Url) RemoveQuery(key string) {
	u.query.Delete(key)
	rebuildUrl(u)
}

// lua:metamethod div
func (u *Url) Div(path string) *Url {
	newPath := u.path.Join(path)
	newUrl := &Url{
		raw:          u.raw,
		scheme:       u.scheme,
		host:         u.host,
		port:         u.port,
		portInferred: u.portInferred,
		username:     u.username,
		password:     u.password,
		path:         newPath,
		query:        u.query,
		fragment:     u.fragment,
	}
	rebuildUrl(newUrl)
	return newUrl
}

// lua:module new
func urlNew(l *lua.State) int {
	u := &Url{
		path:  &Path{raw: "", sep: posixSep},
		query: serde.NewQuery(),
	}
	UrlToLua(l, u)
	return 1
}

// lua:module deserialize
func urlDeserialize(l *lua.State) int {
	raw, _ := l.ToString(1)
	u := parseUrl(raw)
	UrlToLua(l, u)
	return 1
}

// lua:module serialize
func urlSerialize(l *lua.State) int {
	u := toUrl(l, 1)
	l.PushString(u.String())
	return 1
}

func parseUrl(raw string) *Url {
	u := &Url{
		path:  &Path{raw: "", sep: posixSep},
		query: serde.NewQuery(),
	}
	if raw == "" {
		return u
	}

	parsed, err := url.Parse(raw)
	if err != nil {
		u.raw = raw
		u.portInferred = defaultPort("")
		return u
	}

	u.raw = raw
	u.scheme = parsed.Scheme
	u.host = parsed.Hostname()

	if parsed.Port() != "" {
		p, err := strconv.Atoi(parsed.Port())
		if err != nil {
			u.portInferred = defaultPort(parsed.Scheme)
		} else {
			u.port = &p
			u.portInferred = p
		}
	} else {
		u.portInferred = defaultPort(parsed.Scheme)
	}

	if parsed.User != nil {
		username := parsed.User.Username()
		u.username = &username
		if p, ok := parsed.User.Password(); ok {
			u.password = &p
		}
	}

	if parsed.Path != "" {
		u.path = pathFromStringSep(parsed.Path, posixSep)
	}

	if parsed.RawQuery != "" {
		u.query.FromRaw(parsed.RawQuery)
	}

	if parsed.Fragment != "" {
		u.fragment = &parsed.Fragment
	}

	return u
}

func defaultPort(scheme string) int {
	switch scheme {
	case "http":
		return 80
	case "https":
		return 443
	default:
		return 0
	}
}

func rebuildUrl(u *Url) {
	var b strings.Builder
	if u.scheme != "" {
		b.WriteString(u.scheme)
		b.WriteString("://")
	}
	if u.username != nil {
		b.WriteString(*u.username)
		if u.password != nil {
			b.WriteByte(':')
			b.WriteString(*u.password)
		}
		b.WriteByte('@')
	}
	if u.host != "" {
		b.WriteString(u.host)
	}
	if u.port != nil && *u.port != u.portInferred {
		b.WriteByte(':')
		b.WriteString(strconv.Itoa(*u.port))
	}
	b.WriteString(u.path.raw)
	if u.query.Size() > 0 {
		b.WriteByte('?')
		b.WriteString(u.query.ToString())
	}
	if u.fragment != nil {
		b.WriteByte('#')
		b.WriteString(*u.fragment)
	}
	u.raw = b.String()
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
