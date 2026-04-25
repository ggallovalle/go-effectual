package serde

import (
	"net/url"
	"slices"
	"strings"
)

type Query struct {
	params url.Values
}

func NewQuery() *Query {
	return &Query{params: url.Values{}}
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

func (q *Query) ToString() string {
	return q.params.Encode()
}

func (q *Query) Keys() []string {
	keys := make([]string, 0, len(q.params))
	for k := range q.params {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}

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
