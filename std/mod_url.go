package std

import (
	"net/url"
	"strconv"
	"strings"
	"text/template"

	"github.com/speedata/go-lua"
	"github.com/ggallovalle/go-effectual"
	"github.com/ggallovalle/go-effectual/std/serde"
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

type Url struct {
	raw          string
	scheme      string
	host        string
	port        *int
	portInferred int
	username    *string
	password    *string
	path        *PathBuf
	query       *serde.Query
	fragment    *string
}

func (u *Url) String() string {
	return u.raw
}

func urlNew(l *lua.State) int {
	u := &Url{
		path:  &PathBuf{raw: "", sep: posixSep},
		query: serde.NewQuery(),
	}
	urlToLua(l, u)
	return 1
}

func urlDeserialize(l *lua.State) int {
	raw, _ := l.ToString(1)
	u := parseUrl(raw)
	urlToLua(l, u)
	return 1
}

func parseUrl(raw string) *Url {
	u := &Url{
		path:  &PathBuf{raw: "", sep: posixSep},
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
		u.path = pathBufFromStringSep(parsed.Path, posixSep)
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

func urlSerialize(l *lua.State) int {
	u := toUrl(l, 1)
	l.PushString(u.String())
	return 1
}

func urlToLua(l *lua.State, u *Url) {
	l.PushUserData(u)
	lua.SetMetaTableNamed(l, slugUrlHandle)
}

func toUrl(l *lua.State, idx int) *Url {
	return lua.CheckUserData(l, idx, slugUrlHandle).(*Url)
}

func urlAddQuery(l *lua.State) int {
	u := toUrl(l, 1)
	key, _ := l.ToString(2)
	value, _ := l.ToString(3)
	u.query.Append(key, value)
	rebuildUrl(u)
	return 0
}

func urlRemoveQuery(l *lua.State) int {
	u := toUrl(l, 1)
	key, _ := l.ToString(2)
	u.query.Delete(key)
	rebuildUrl(u)
	return 0
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

var urlMethods = map[string]lua.Function{
	"add_query":    urlAddQuery,
	"remove_query": urlRemoveQuery,
}

var urlGetters = map[string]func(*lua.State){
	"scheme": func(l *lua.State) {
		u := toUrl(l, 1)
		if u.scheme == "" {
			l.PushNil()
		} else {
			l.PushString(u.scheme)
		}
	},
	"host": func(l *lua.State) {
		u := toUrl(l, 1)
		if u.host == "" {
			l.PushNil()
		} else {
			l.PushString(u.host)
		}
	},
	"port": func(l *lua.State) {
		u := toUrl(l, 1)
		if u.port == nil {
			l.PushNil()
		} else {
			l.PushInteger(*u.port)
		}
	},
	"port_inferred": func(l *lua.State) {
		u := toUrl(l, 1)
		l.PushInteger(u.portInferred)
	},
	"username": func(l *lua.State) {
		u := toUrl(l, 1)
		if u.username == nil {
			l.PushNil()
		} else {
			l.PushString(*u.username)
		}
	},
	"password": func(l *lua.State) {
		u := toUrl(l, 1)
		if u.password == nil {
			l.PushNil()
		} else {
			l.PushString(*u.password)
		}
	},
	"path": func(l *lua.State) {
		u := toUrl(l, 1)
		pathBufToLua(l, u.path)
	},
	"query": func(l *lua.State) {
		u := toUrl(l, 1)
		serde.QueryToLua(l, u.query)
	},
	"fragment": func(l *lua.State) {
		u := toUrl(l, 1)
		if u.fragment == nil {
			l.PushNil()
		} else {
			l.PushString(*u.fragment)
		}
	},
}

var urlMetatable = []lua.RegistryFunction{
	{Name: "__tostring", Function: func(l *lua.State) int {
		u := toUrl(l, 1)
		l.PushString(u.String())
		return 1
	}},
	{Name: "__div", Function: func(l *lua.State) int {
		u := toUrl(l, 1)
		arg := toPathBufString(l, 2)
		newPath := u.path.join(arg)
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
		urlToLua(l, newUrl)
		return 1
	}},
	effectual.LuaMetaIndex(urlGetters, urlMethods),
}

func urlLibrary() []lua.RegistryFunction {
	return []lua.RegistryFunction{
		{Name: "new", Function: urlNew},
		{Name: "deserialize", Function: urlDeserialize},
		{Name: "serialize", Function: urlSerialize},
	}
}

var urlAnnotationsTmpl = template.Must(template.New("UrlAnnotations").Parse(`---@meta {{.module}}

---@class (exact) {{.Url}} : userdata
---@field scheme string|nil
---@field host string|nil
---@field port integer|nil
---@field port_inferred integer
---@field username string|nil
---@field password string|nil
---@field path {{.PathBuf}}
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

func (lib *ModUrl) Name() string {
	return lib.name
}

func (lib *ModUrl) Annotations() string {
	data := map[string]string{
		"module":  lib.name,
		"Url":     lib.name + ".Url",
		"PathBuf": "std.path.PathBuf",
		"Query":   "std.serde.query.Query",
	}
	var buf strings.Builder
	if err := urlAnnotationsTmpl.Execute(&buf, data); err != nil {
		return ""
	}
	return buf.String()
}

func (lib *ModUrl) Open(l *lua.State) int {
	lua.NewLibrary(l, urlLibrary())

	lua.NewMetaTable(l, slugUrlHandle)
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
