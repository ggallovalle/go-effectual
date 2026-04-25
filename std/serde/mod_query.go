package serde

import (
	"strings"
	"text/template"

	"github.com/speedata/go-lua"
	"github.com/ggallovalle/go-effectual"
)

type ModQuery struct {
	name string
}

type ModQueryApi struct {
	mod *ModQuery
	lua *lua.State
}

const (
	modQueryName    = "std.serde.query"
	slugQueryHandle = "go/std/serde/query/Query*"
)

func QueryToLua(l *lua.State, q *Query) {
	l.PushUserData(q)
	lua.SetMetaTableNamed(l, slugQueryHandle)
}

func toQuery(l *lua.State, idx int) *Query {
	return lua.CheckUserData(l, idx, slugQueryHandle).(*Query)
}

func queryHas(l *lua.State) int {
	q := toQuery(l, 1)
	key, _ := l.ToString(2)
	l.PushBoolean(q.Has(key))
	return 1
}

func queryGet(l *lua.State) int {
	q := toQuery(l, 1)
	key, _ := l.ToString(2)
	val := q.Get(key)
	if val == "" && !q.Has(key) {
		l.PushNil()
	} else {
		l.PushString(val)
	}
	return 1
}

func queryGetAll(l *lua.State) int {
	q := toQuery(l, 1)
	key, _ := l.ToString(2)
	vals := q.GetAll(key)
	l.NewTable()
	for i, v := range vals {
		l.PushInteger(i + 1)
		l.PushString(v)
		l.SetTable(-3)
	}
	return 1
}

func querySet(l *lua.State) int {
	q := toQuery(l, 1)
	key, _ := l.ToString(2)
	value, _ := l.ToString(3)
	q.Set(key, value)
	return 0
}

func queryAppend(l *lua.State) int {
	q := toQuery(l, 1)
	key, _ := l.ToString(2)
	value, _ := l.ToString(3)
	q.Append(key, value)
	return 0
}

func queryDelete(l *lua.State) int {
	q := toQuery(l, 1)
	key, _ := l.ToString(2)
	q.Delete(key)
	return 0
}

func querySort(l *lua.State) int {
	q := toQuery(l, 1)
	q.Sort()
	return 0
}

func queryToString(l *lua.State) int {
	q := toQuery(l, 1)
	l.PushString(q.ToString())
	return 1
}

func queryKeys(l *lua.State) int {
	q := toQuery(l, 1)
	keys := q.Keys()
	l.NewTable()
	for i, k := range keys {
		l.PushInteger(i + 1)
		l.PushString(k)
		l.SetTable(-3)
	}
	return 1
}

func queryValues(l *lua.State) int {
	q := toQuery(l, 1)
	vals := q.Values()
	l.NewTable()
	for i, v := range vals {
		l.PushInteger(i + 1)
		l.PushString(v)
		l.SetTable(-3)
	}
	return 1
}

func queryEntries(l *lua.State) int {
	q := toQuery(l, 1)
	entries := q.Entries()
	l.NewTable()
	for i, entry := range entries {
		l.PushInteger(i + 1)
		l.NewTable()
		l.PushInteger(1)
		l.PushString(entry[0])
		l.SetTable(-3)
		l.PushInteger(2)
		l.PushString(entry[1])
		l.SetTable(-3)
		l.SetTable(-3)
	}
	return 1
}

var queryMethods = map[string]lua.Function{
	"has":        queryHas,
	"get":        queryGet,
	"get_all":    queryGetAll,
	"set":        querySet,
	"append":     queryAppend,
	"delete":     queryDelete,
	"sort":       querySort,
	"to_string":  queryToString,
	"keys":       queryKeys,
	"values":     queryValues,
	"entries":    queryEntries,
}

var queryGetters = map[string]func(*lua.State){
	"size": func(l *lua.State) {
		q := toQuery(l, 1)
		l.PushInteger(q.Size())
	},
}

var queryMetatable = []lua.RegistryFunction{
	{Name: "__tostring", Function: func(l *lua.State) int {
		q := toQuery(l, 1)
		l.PushString(q.ToString())
		return 1
	}},
	{Name: "__pairs", Function: func(l *lua.State) int {
		q := toQuery(l, 1)
		q.Sort()
		l.PushGoFunction(func(l *lua.State) int {
			idx, _ := l.ToInteger(2)
			q := toQuery(l, 1)
			keys := q.Keys()
			if idx < len(keys) {
				k := keys[idx]
				l.PushInteger(idx + 1)
				l.PushString(k)
				l.PushString(q.Get(k))
				return 3
			}
			return 0
		})
		l.PushValue(1)
		l.PushInteger(0)
		return 3
	}},
	effectual.LuaMetaIndex(queryGetters, queryMethods),
}

func queryNew(l *lua.State) int {
	q := NewQuery()
	QueryToLua(l, q)
	return 1
}

func queryDeserialize(l *lua.State) int {
	raw, _ := l.ToString(1)
	q := NewQuery()
	if raw != "" {
		q.FromRaw(raw)
	}
	QueryToLua(l, q)
	return 1
}

func querySerialize(l *lua.State) int {
	q := toQuery(l, 1)
	l.PushString(q.ToString())
	return 1
}

func queryLibrary() []lua.RegistryFunction {
	return []lua.RegistryFunction{
		{Name: "new", Function: queryNew},
		{Name: "deserialize", Function: queryDeserialize},
		{Name: "serialize", Function: querySerialize},
	}
}

var queryAnnotationsTmpl = template.Must(template.New("QueryAnnotations").Parse(`---@meta {{.module}}

---@class (exact) {{.Query}} : userdata
---@field size integer
local Query = {}

---@param key string
---@return boolean
function Query:has(key) end

---@param key string
---@return string|nil
function Query:get(key) end

---@param key string
---@return string[]
function Query:get_all(key) end

---@param key string
---@param value string
function Query:set(key, value) end

---@param key string
---@param value string
function Query:append(key, value) end

---@param key string
function Query:delete(key) end

function Query:sort() end

---@return string
function Query:to_string() end

---@return string[]
function Query:keys() end

---@return string[]
function Query:values() end

---@return {[1]: string, [2]: string}[]
function Query:entries() end

local {{.name}} = {}

---@return {{.Query}}
function {{.name}}.new() end

---@param raw string
---@return {{.Query}}
function {{.name}}.deserialize(raw) end

---@param q {{.Query}}
---@return string
function {{.name}}.serialize(q) end

return {{.name}}
`))

func (lib *ModQuery) Name() string {
	return lib.name
}

func (lib *ModQuery) Annotations() string {
	data := map[string]string{
		"module": lib.name,
		"name":   "query",
		"Query":  lib.name + ".Query",
	}
	var buf strings.Builder
	if err := queryAnnotationsTmpl.Execute(&buf, data); err != nil {
		return ""
	}
	return buf.String()
}

func (lib *ModQuery) Open(l *lua.State) int {
	lua.NewLibrary(l, queryLibrary())

	lua.NewMetaTable(l, slugQueryHandle)
	l.PushValue(-1)
	l.SetField(-2, "__index")
		lua.SetFunctions(l, queryMetatable, 0)
		for name, fn := range queryMethods {
			l.PushGoFunction(fn)
			l.SetField(-2, name)
		}
		l.Pop(1)

	return 1
}

func (lib *ModQuery) OpenLib(l *lua.State) {
	lua.Require(l, lib.name, lib.Open, false)
	l.Pop(1)
}

func (lib *ModQuery) Require(l *lua.State) {
	l.Global("require")
	l.PushString(lib.Name())
	l.Call(1, 1)
}

func (lib *ModQuery) Api(l *lua.State) ModQueryApi {
	return ModQueryApi{mod: lib, lua: l}
}

func MakeModQuery() effectual.LuaMod[ModQueryApi] {
	return &ModQuery{name: modQueryName}
}