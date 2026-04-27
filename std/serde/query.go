package serde

import (
	"net/url"
	"slices"
	"strings"

	"github.com/speedata/go-lua"
)

//lua: module std.serde.query
//lua: class Query

type Query struct {
	params url.Values //lua: skip-field
}

//lua: module-fn new
func NewQuery() *Query {
	return &Query{params: url.Values{}}
}

//lua: module-fn deserialize
func Deserialize(raw string) *Query {
	q := NewQuery()
	if raw != "" {
		q.FromRaw(raw)
	}
	return q
}

//lua: module-fn serialize
func Serialize(q *Query) string {
	return q.ToString()
}

//lua: metamethod __tostring
//lua: raw
func QueryToString(l *lua.State) int {
	q := toQuery(l, 1)
	l.PushString(q.ToString())
	return 1
}

//lua: metamethod __pairs
//lua: raw
func QueryPairs(l *lua.State) int {
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
}

func (q *Query) FromRaw(raw string) {
	q.params, _ = url.ParseQuery(strings.TrimPrefix(raw, "?"))
}

func (q *Query) Size() int {
	return len(q.params)
}

func (q *Query) Has(key string) bool {
	_, ok := q.params[key]
	return ok
}

//lua: nil-map
func (q *Query) Get(key string) string {
	vals := q.params[key]
	if len(vals) == 0 {
		return ""
	}
	return vals[0]
}

func (q *Query) GetAll(key string) []string {
	return q.params[key]
}

func (q *Query) Set(key, value string) {
	q.params[key] = []string{value}
}

func (q *Query) Append(key, value string) {
	q.params[key] = append(q.params[key], value)
}

func (q *Query) Delete(key string) {
	delete(q.params, key)
}

func (q *Query) Sort() {
	keys := make([]string, 0, len(q.params))
	for k := range q.params {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	sorted := url.Values{}
	for _, k := range keys {
		sorted[k] = q.params[k]
	}
	q.params = sorted
}

//lua: force-method
func (q *Query) ToString() string {
	return q.params.Encode()
}

//lua: force-method
func (q *Query) Keys() []string {
	keys := make([]string, 0, len(q.params))
	for k := range q.params {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}

//lua: force-method
func (q *Query) Values() []string {
	keys := make([]string, 0, len(q.params))
	for k := range q.params {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	vals := make([]string, 0)
	for _, k := range keys {
		vals = append(vals, q.params[k]...)
	}
	return vals
}

//lua: force-method
func (q *Query) Entries() [][2]string {
	keys := make([]string, 0, len(q.params))
	for k := range q.params {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	entries := make([][2]string, 0)
	for _, k := range keys {
		for _, val := range q.params[k] {
			entries = append(entries, [2]string{k, val})
		}
	}
	return entries
}