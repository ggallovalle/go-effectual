package std_test

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	lua "github.com/speedata/go-lua"

	"github.com/ggallovalle/go-effectual"
	sut "github.com/ggallovalle/go-effectual/std"
	"github.com/stretchr/testify/assert"
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

func Test_LibGoLogSlug_Levels(t *testing.T) {
	l := lua.NewState()
	lua.OpenLibraries(l)
	api := effectual.LuaModOpenWithApi(l, sut.MakeModSlog())

	for _, tc := range []struct {
		name  string
		level slog.Level
	}{
		{"DEBUG", slog.LevelDebug},
		{"INFO", slog.LevelInfo},
		{"WARN", slog.LevelWarn},
		{"ERROR", slog.LevelError},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := slog.New(&allLevelHandler{sink: &buf, level: tc.level})
			api.New(logger)
			l.SetGlobal("logger")

			err := lua.DoString(l, `
				assert(logger:level() == "`+tc.name+`")
			`)
			assert.NoError(t, err)
		})
	}
}

func Test_LibGoLogSlug_Default(t *testing.T) {
	l := lua.NewState()
	lua.OpenLibraries(l)
	api := effectual.LuaModOpenWithApi(l, sut.MakeModSlog())

	var buf bytes.Buffer
	logger := slog.New(&allLevelHandler{sink: &buf, level: slog.LevelDebug})
	api.SetDefault(logger)

	t.Run("default returns the logger", func(t *testing.T) {
		err := lua.DoString(l, `
			local log = require("std.log")
			assert(log.default ~= nil, "expected non-nil default")
		`)
		assert.NoError(t, err)
	})

	t.Run("debug delegates to default", func(t *testing.T) {
		buf.Reset()
		err := lua.DoString(l, `
			local log = require("std.log")
			log:debug("hello from module")`)
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "hello from module")
	})

	t.Run("info delegates to default", func(t *testing.T) {
		buf.Reset()
		err := lua.DoString(l, `
			local log = require("std.log")
			log:info("hello info")`)
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "hello info")
	})

	t.Run("warn delegates to default", func(t *testing.T) {
		buf.Reset()
		err := lua.DoString(l, `
			local log = require("std.log")
			log:warn("hello warn")`)
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "hello warn")
	})

	t.Run("error delegates to default", func(t *testing.T) {
		buf.Reset()
		err := lua.DoString(l, `
			local log = require("std.log")
			log:error("hello error")`)
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "hello error")
	})

	t.Run("log delegates to default", func(t *testing.T) {
		buf.Reset()
		err := lua.DoString(l, `
			local log = require("std.log")
			log:log("DEBUG", "via log")`)
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "via log")
	})

	t.Run("level delegates to default", func(t *testing.T) {
		err := lua.DoString(l, `
			local log = require("std.log")
			local lv = log:level()
			assert(lv == "DEBUG", "expected DEBUG but got " .. lv)
		`)
		assert.NoError(t, err)
	})

	t.Run("attrs passed through", func(t *testing.T) {
		buf.Reset()
		err := lua.DoString(l, `
			local log = require("std.log")
			log:info("msg", {key = "value"})`)
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "key")
		assert.Contains(t, buf.String(), "value")
	})
}

func Test_LibGoLogSlug_Default_Errors(t *testing.T) {
	l := lua.NewState()
	lua.OpenLibraries(l)

	effectual.LuaModOpen(l, sut.MakeModSlog())

	err := lua.DoString(l, `
		local std = require("std.log")
		std:debug("test")
	`)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "std.log: no default logger set")
}
