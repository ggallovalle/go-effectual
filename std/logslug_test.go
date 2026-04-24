package std_test

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	lua "github.com/Shopify/go-lua"

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
	loggerLib := sut.MakeLibGoLogSlug()
	loggerLib.OpenLib(l)

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
			loggerLib.LuaNew(l, logger)
			l.SetGlobal("logger")

			err := lua.DoString(l, `
				assert(logger:level() == "`+tc.name+`")
			`)
			assert.NoError(t, err)
		})
	}
}
