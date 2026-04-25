package std

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/speedata/go-lua"
)

type Path struct {
	raw string
	sep string
}

func (p *Path) String() string {
	return p.raw
}

func (p *Path) dir() string {
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

func (p *Path) Push(path string) {
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

func (p *Path) Pop() bool {
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

func (p *Path) Join(path string) *Path {
	if strings.HasPrefix(path, p.sep) || strings.HasPrefix(path, altSep(p.sep)) {
		return &Path{raw: path, sep: p.sep}
	}
	newBuf := &Path{raw: p.raw, sep: p.sep}
	newBuf.Push(path)
	return newBuf
}

func (p *Path) EndsWith(child string) bool {
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

func (p *Path) StartsWith(base string) bool {
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

func (p *Path) StripPrefix(prefix string) (*Path, error) {
	if !strings.HasPrefix(p.raw, prefix) {
		return nil, errors.New("prefix not found")
	}
	result := strings.TrimPrefix(p.raw, prefix)
	result = strings.TrimPrefix(result, "/")
	result = strings.TrimPrefix(result, "\\")
	if result == "" {
		result = p.sep
	}
	return &Path{raw: result, sep: p.sep}, nil
}

func (p *Path) WithExtension(ext string) *Path {
	base := filepath.Base(p.raw)
	dotIdx := strings.LastIndex(base, ".")
	if dotIdx <= 0 {
		return &Path{raw: p.raw + "." + ext, sep: p.sep}
	}
	stem := p.raw[:len(p.raw)-len(base)+dotIdx]
	return &Path{raw: stem + "." + ext, sep: p.sep}
}

func (p *Path) WithFileName(name string) *Path {
	dir := p.dir()
	if dir == "." {
		return &Path{raw: name, sep: p.sep}
	}
	if dir != p.sep && !strings.HasSuffix(dir, p.sep) && !strings.HasSuffix(dir, altSep(p.sep)) {
		dir += p.sep
	}
	return &Path{raw: dir + name, sep: p.sep}
}

func (p *Path) Components() []string {
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

func (p *Path) Ancestors() []*Path {
	var result []*Path
	current := &Path{raw: p.raw, sep: p.sep}
	for current.raw != "" && current.raw != p.sep && current.raw != altSep(p.sep) {
		result = append(result, current)
		parent := &Path{raw: current.dir(), sep: p.sep}
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

func (p *Path) Parent() *Path {
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
		return &Path{raw: p.sep, sep: p.sep}
	}
	return &Path{raw: dir, sep: p.sep}
}

func (p *Path) baseName() string {
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

func (p *Path) FileName() string {
	name := p.baseName()
	if name == "" || name == ".." {
		return ""
	}
	return name
}

func (p *Path) Extension() string {
	name := p.FileName()
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

func (p *Path) FileStem() string {
	name := p.FileName()
	if name == "" {
		return ""
	}
	dotIdx := strings.LastIndex(name, ".")
	if dotIdx <= 0 {
		return name
	}
	return name[:dotIdx]
}

func (p *Path) HasRoot() bool {
	return strings.HasPrefix(p.raw, p.sep) || strings.HasPrefix(p.raw, altSep(p.sep))
}

func (p *Path) IsAbsolute() bool {
	return p.HasRoot()
}

func (p *Path) IsRelative() bool {
	return !p.HasRoot()
}

func toPathString(l *lua.State, idx int) string {
	if l.IsUserData(idx) {
		v := l.ToUserData(idx)
		if v == nil {
			if s, ok := l.ToString(idx); ok {
				return s
			}
			return ""
		}
		switch x := v.(type) {
		case *Path:
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

func pathFromStringSep(s, sep string) *Path {
	alt := altSep(sep)
	trimmed := strings.Trim(s, "/\\")
	if trimmed == "" && strings.ContainsAny(s, "/\\") {
		return &Path{raw: sep, sep: sep}
	}
	if s == "" {
		return &Path{raw: "", sep: sep}
	}
	if !strings.HasPrefix(s, sep) && !strings.HasPrefix(s, alt) {
		return &Path{raw: s, sep: sep}
	}
	return &Path{raw: sep + trimmed, sep: sep}
}
