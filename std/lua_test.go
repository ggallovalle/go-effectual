package std_test

import (
	"bytes"
	"context"
	"log/slog"
	"path/filepath"
	"testing"

	lua "github.com/speedata/go-lua"

	"github.com/ggallovalle/go-effectual"
	sut "github.com/ggallovalle/go-effectual/std"
	serde "github.com/ggallovalle/go-effectual/std/serde"
)

type allLevelHandler struct {
	sink  *bytes.Buffer
	level slog.Level
}

func (h *allLevelHandler) Enabled(_ context.Context, l slog.Level) bool { return l >= h.level }
func (h *allLevelHandler) Handle(_ context.Context, r slog.Record) error {
	if r.Level < h.level {
		return nil
	}
	return slog.NewTextHandler(h.sink, nil).Handle(context.Background(), r)
}
func (h *allLevelHandler) WithAttrs([]slog.Attr) slog.Handler { return h }
func (h *allLevelHandler) WithGroup(string) slog.Handler      { return h }

func TestLuaSuite_Semver(t *testing.T) {
	l := lua.NewState()
	lua.OpenLibraries(l)
	effectual.LuaModOpen(l, sut.MakeModSemver())
	effectual.LuaModOpen(l, sut.MakeModTesting())

	testFile := filepath.Join("..", "luahome", "std-test", "semver_test.lua")
	runLuaSuite(t, l, testFile)
}

func TestLuaSuite_Path(t *testing.T) {
	l := lua.NewState()
	lua.OpenLibraries(l)
	effectual.LuaModOpen(l, sut.MakeModPath())
	effectual.LuaModOpen(l, sut.MakeModTesting())

	testFile := filepath.Join("..", "luahome", "std-test", "path_test.lua")
	runLuaSuite(t, l, testFile)
}

func TestLuaSuite_Url(t *testing.T) {
	l := lua.NewState()
	lua.OpenLibraries(l)
	effectual.LuaModOpen(l, sut.MakeModUrl())
	effectual.LuaModOpen(l, sut.MakeModTesting())

	testFile := filepath.Join("..", "luahome", "std-test", "url_test.lua")
	runLuaSuite(t, l, testFile)
}

func TestLuaSuite_Slog(t *testing.T) {
	l := lua.NewState()
	lua.OpenLibraries(l)
	api := effectual.LuaModOpen(l, sut.MakeModSlog())
	effectual.LuaModOpen(l, sut.MakeModTesting())

	logger := slog.New(&allLevelHandler{sink: nil, level: slog.LevelDebug})
	api.New(logger)
	l.SetGlobal("logger")
	api.SetDefault(logger)

	testFile := filepath.Join("..", "luahome", "std-test", "slog_test.lua")
	runLuaSuite(t, l, testFile, &LuaTestLoggerExtension{Logger: logger})
}

func TestLuaSuite_Query(t *testing.T) {
	l := lua.NewState()
	lua.OpenLibraries(l)
	effectual.LuaModOpen(l, serde.MakeModQuery())
	effectual.LuaModOpen(l, sut.MakeModTesting())

	testFile := filepath.Join("..", "luahome", "std-test", "serde", "query_test.lua")
	runLuaSuite(t, l, testFile)
}

func TestLuaSuite_Slog_NoDefaultErrors(t *testing.T) {
	l := lua.NewState()
	lua.OpenLibraries(l)
	effectual.LuaModOpen(l, sut.MakeModSlog())

	err := lua.DoString(l, `
		local std = require("std.log")
		std:debug("test")
	`)
	if err == nil {
		t.Fatal("expected error when no default logger set")
	}
	if !bytes.Contains([]byte(err.Error()), []byte("std.log: no default logger set")) {
		t.Fatalf("expected 'no default logger set' error, got: %v", err)
	}
}
