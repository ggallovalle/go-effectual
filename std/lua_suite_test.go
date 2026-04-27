package std_test

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
	"unicode"

	lua "github.com/speedata/go-lua"
	"github.com/stretchr/testify/assert"
	"github.com/ggallovalle/go-effectual"
	"github.com/ggallovalle/go-effectual/std"
)

type testLevelHandler struct {
	sink  *bytes.Buffer
	level slog.Level
}

func (h *testLevelHandler) Enabled(_ context.Context, l slog.Level) bool { return l >= h.level }
func (h *testLevelHandler) Handle(_ context.Context, r slog.Record) error {
	return slog.NewTextHandler(h.sink, nil).Handle(context.Background(), r)
}
func (h *testLevelHandler) WithAttrs([]slog.Attr) slog.Handler { return h }
func (h *testLevelHandler) WithGroup(string) slog.Handler      { return h }

type LuaTestLoggerExtension struct {
	Logger *slog.Logger
}

func (e *LuaTestLoggerExtension) Name() string {
	return "logger"
}

func (e *LuaTestLoggerExtension) Build(l *lua.State, params map[string]any) {
	logger := e.Logger
	if logger == nil {
		logger = slog.Default()
	}
	if params != nil {
		if levelStr, ok := params["level"].(string); ok {
			var level slog.Level
			switch levelStr {
			case "DEBUG":
				level = slog.LevelDebug
			case "INFO":
				level = slog.LevelInfo
			case "WARN":
				level = slog.LevelWarn
			case "ERROR":
				level = slog.LevelError
			default:
				level = slog.LevelDebug
			}
			var buf bytes.Buffer
			logger = slog.New(&testLevelHandler{sink: &buf, level: level})
		}
	}
	l.PushUserData(logger)
	lua.SetMetaTableNamed(l, "go/std/log/slug/Logger*")
}

func runLuaSuite(t *testing.T, l *lua.State, path string, extensions ...std.LuaTestCtxExtension) {
	extMap := make(map[string]std.LuaTestCtxExtension)
	for _, ext := range extensions {
		extMap[ext.Name()] = ext
	}

	err := lua.DoFile(l, path)
	if !assert.NoError(t, err) {
		t.Fatalf("failed to execute test file %q: %v", path, err)
	}

	if l.TypeOf(-1) != lua.TypeTable {
		t.Fatalf("expected test file %q to return a table, got %s", path, lua.TypeNameOf(l, -1))
	}

	suiteIdx := l.AbsIndex(-1)

	l.PushString("name")
	l.RawGet(suiteIdx)
	suiteName, _ := l.ToString(-1)
	l.Pop(1)

	l.PushString("deps")
	l.RawGet(suiteIdx)
	suiteDepsIdx := 0
	if !l.IsNil(-1) {
		suiteDepsIdx = l.AbsIndex(-1)
	} else {
		l.Pop(1)
	}

	l.PushString("cases")
	l.RawGet(suiteIdx)
	if l.TypeOf(-1) != lua.TypeTable {
		t.Fatalf("expected suite to have a 'cases' table, got %s", lua.TypeNameOf(l, -1))
	}
	casesIdx := l.AbsIndex(-1)

	l.Global("require")
	l.PushString("std.testing")
	l.Call(1, 1)
	l.PushString("ctx")
	l.RawGet(-2)
	l.SetGlobal("__ctx_func")
	l.Pop(1)

	l.PushNil()
	for l.Next(casesIdx) {
		if !l.IsTable(-1) {
			l.Pop(1)
			continue
		}

		caseIdx := l.AbsIndex(-1)

		l.PushString("name")
		l.RawGet(caseIdx)
		caseName, _ := l.ToString(-1)
		l.Pop(1)

		sanitized := sanitizeTestName(caseName)

		l.PushString("fn")
		l.RawGet(caseIdx)
		fnIdx := l.AbsIndex(-1)
		caseLoc := getFuncLocation(l, fnIdx)
		l.Pop(1)

		l.PushString("deps")
		l.RawGet(caseIdx)
		caseDepsIdx := 0
		if !l.IsNil(-1) {
			caseDepsIdx = l.AbsIndex(-1)
		} else {
			l.Pop(1)
		}

		mergedDeps := mergeDeps(l, suiteDepsIdx, caseDepsIdx)

		t.Run(sanitized, func(t *testing.T) {
			top := l.Top()
			defer l.SetTop(top)

			l.PushString("fn")
			l.RawGet(caseIdx)
			if !l.IsFunction(-1) {
				t.Fatalf("case %q: 'fn' is not a function", caseName)
			}

			l.Global("__ctx_func")
			l.PushUserData(t)
			l.PushString(suiteName)
			l.PushString(caseName)
			l.Call(3, 1)

			l.NewTable()
			extIdx := l.AbsIndex(-1)

			available := make([]string, 0, len(extMap))
			for name := range extMap {
				available = append(available, name)
			}

			var entries []std.CtxExtEntry
			for _, dep := range mergedDeps {
				ext, ok := extMap[dep.name]
				if !ok {
					t.Fatalf("%s (case %q): extension '%s' not found. Available: %s", caseLoc, caseName, dep.name, strings.Join(available, ", "))
				}
				ext.Build(l, dep.params)
				l.SetField(extIdx, dep.name)
				entries = append(entries, std.CtxExtEntry{Ext: ext, Params: dep.params})
			}

			std.SetCtxExt(t, entries)

			l.Pop(1)

			if err := l.ProtectedCall(1, 1, 0); err != nil {
				msg := err.Error()
				if idx := strings.Index(msg, "__SKIP__"); idx != -1 {
					t.Skip(msg[idx+len("__SKIP__"):])
				}
				if l.IsString(-1) {
					s, _ := l.ToString(-1)
					t.Fatalf("%s: %s", caseName, s)
				}
				t.Fatalf("%s: %v", caseName, err)
			}
			if l.IsString(-1) {
				s, _ := l.ToString(-1)
				t.Fatalf("%s: %s", caseName, s)
			}
		})

		if caseDepsIdx != 0 {
			l.Pop(1)
		}
		l.Pop(1)
	}
	if suiteDepsIdx != 0 {
		l.Pop(1)
	}
	l.Pop(1)
}

type resolvedDep struct {
	name   string
	params map[string]any
}

func mergeDeps(l *lua.State, suiteDepsIdx, caseDepsIdx int) []resolvedDep {
	result := make([]resolvedDep, 0)
	seen := make(map[string]bool)

	if caseDepsIdx != 0 {
		l.PushNil()
		for l.Next(caseDepsIdx) {
			dep := parseDep(l, l.AbsIndex(-1))
			if dep != nil {
				result = append(result, *dep)
				seen[dep.name] = true
			}
			l.Pop(1)
		}
	}

	if suiteDepsIdx != 0 {
		l.PushNil()
		for l.Next(suiteDepsIdx) {
			dep := parseDep(l, l.AbsIndex(-1))
			if dep != nil && !seen[dep.name] {
				result = append(result, *dep)
				seen[dep.name] = true
			}
			l.Pop(1)
		}
	}

	return result
}

func parseDep(l *lua.State, idx int) *resolvedDep {
	if l.IsString(idx) {
		s, _ := l.ToString(idx)
		return &resolvedDep{name: s}
	}
	if l.IsTable(idx) {
		l.PushString("name")
		l.RawGet(idx)
		name, _ := l.ToString(-1)
		l.Pop(1)

		l.PushString("params")
		l.RawGet(idx)
		var params map[string]any
		if !l.IsNil(-1) {
			params = tableToMap(l, -1)
		}
		l.Pop(1)

		return &resolvedDep{name: name, params: params}
	}
	return nil
}

func tableToMap(l *lua.State, idx int) map[string]any {
	result := make(map[string]any)
	idx = l.AbsIndex(idx)
	l.PushNil()
	for l.Next(idx) {
		if key, ok := l.ToString(-2); ok {
			result[key] = effectual.ConvertLuaToAny(l, -1)
		}
		l.Pop(1)
	}
	return result
}

func sanitizeTestName(name string) string {
	var b strings.Builder
	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			b.WriteRune(r)
		} else {
			b.WriteRune('_')
		}
	}
	s := strings.Trim(b.String(), "_")
	if s == "" {
		return "unnamed"
	}
	return s
}

func getFuncLocation(l *lua.State, fnIdx int) string {
	l.PushValue(fnIdx)
	script := `
		local fn = ...
		local info = debug.getinfo(fn, "Sl")
		if info and info.source and info.currentline and info.currentline > 0 then
			return info.source .. ":" .. info.currentline
		elseif info and info.source and info.linedefined and info.linedefined > 0 then
			return info.source .. ":" .. info.linedefined
		elseif info and info.source then
			return info.source
		end
		return ""
	`
	if err := lua.LoadString(l, script); err != nil {
		l.Pop(2)
		return ""
	}
	l.Insert(-2)
	if err := l.ProtectedCall(1, 1, 0); err != nil {
		l.Pop(1)
		return ""
	}
	loc, _ := l.ToString(-1)
	l.Pop(1)
	if loc == "" {
		return ""
	}
	loc = strings.TrimPrefix(loc, "@")
	return loc
}
