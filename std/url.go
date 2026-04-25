package std

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/ggallovalle/go-effectual/std/serde"
	"github.com/speedata/go-lua"
)

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
