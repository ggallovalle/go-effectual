package effectual

import (
	"fmt"
	"strings"

	"github.com/Shopify/go-lua"
)

const BuiltinPlaceholder = "BUILTIN"

func PackagePathAppend(l *lua.State, path string) {
	packageModify(l, "path", path)
}

func PackageCPathAppend(l *lua.State, cpath string) {
	packageModify(l, "cpath", cpath)
}

func PackagePathPrepend(l *lua.State, path string) {
	packageModify(l, "path", path)
}

func PackageCPathPrepend(l *lua.State, cpath string) {
	packageModify(l, "cpath", cpath)
}

func PackagePathReplace(l *lua.State, paths []string) error {
	return packageReplace(l, "path", paths)
}

func PackageCPathReplace(l *lua.State, cpaths []string) error {
	return packageReplace(l, "cpath", cpaths)
}

func packageModify(l *lua.State, field, value string) {
	l.Global("package")
	l.Field(-1, field)
	old, _ := l.ToString(-1)
	l.Pop(1)

	var new string
	if old == "" {
		new = value
	} else {
		new = old + ";" + value
	}

	l.PushString(field)
	l.PushString(new)
	l.SetTable(-3)
}

func packageReplace(l *lua.State, field string, paths []string) error {
	l.Global("package")
	l.Field(-1, field)
	old, ok := l.ToString(-1)
	l.Pop(1)
	if !ok {
		return fmt.Errorf("could not get current package.%s", field)
	}

	builtinIdx := -1
	for i, p := range paths {
		if p == BuiltinPlaceholder {
			builtinIdx = i
			break
		}
	}

	var result string
	if builtinIdx == -1 {
		result = strings.Join(paths, ";")
	} else {
		before := make([]string, 0, builtinIdx)
		after := make([]string, 0, len(paths)-builtinIdx-1)
		for i, p := range paths {
			if p == BuiltinPlaceholder {
				continue
			}
			if i < builtinIdx {
				before = append(before, p)
			} else {
				after = append(after, p)
			}
		}
		if len(before) > 0 {
			result = strings.Join(before, ";")
		}
		if old != "" {
			result += ";" + old
		}
		if len(after) > 0 {
			result += ";" + strings.Join(after, ";")
		}
	}

	l.Field(-1, field)
	l.PushString(result)
	l.SetTable(-3)
	l.Pop(1)
	return nil
}