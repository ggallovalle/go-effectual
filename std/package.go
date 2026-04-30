package std

import (
	"errors"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/ggallovalle/go-effectual"
	"github.com/ggallovalle/go-effectual/std/serde"
	"github.com/speedata/go-lua"
	"github.com/twpayne/go-vfs"
)

type StdPackage struct {
	modSlog  effectual.LuaMod[ModSlogApi]
	modPath  effectual.LuaMod[ModPathApi]
	modUrl   effectual.LuaMod[ModUrlApi]
	modQuery effectual.LuaModDefinition
}

func NewStdPackage() *StdPackage {
	return &StdPackage{
		modSlog:  MakeModSlog(),
		modPath:  MakeModPath(),
		modUrl:   MakeModUrl(),
		modQuery: serde.MakeModQuery(),
	}
}

func (s *StdPackage) OpenLib(l *lua.State, logger *slog.Logger) error {
	if logger == nil {
		return errors.New("logger required")
	}

	s.modSlog.OpenLib(l)
	s.modPath.OpenLib(l)
	s.modUrl.OpenLib(l)
	s.modQuery.OpenLib(l)

	slogApi := s.modSlog.Api(l)
	slogApi.SetDefault(logger)

	return nil
}

func (s *StdPackage) GenerateAnnotations(fs vfs.FS, folder string) error {
	mods := []effectual.LuaModDefinition{
		s.modSlog,
		s.modPath,
		s.modUrl,
		s.modQuery,
	}

	for _, mod := range mods {
		annotations := mod.Annotations()
		if annotations == "" {
			continue
		}

		parts := strings.Split(mod.Name(), ".")
		filename := parts[len(parts)-1] + ".lua"
		dir := filepath.Join(folder, filepath.Join(parts[:len(parts)-1]...))
		path := filepath.Join(dir, filename)

		if err := vfs.MkdirAll(fs, dir, 0755); err != nil {
			return err
		}

		if err := fs.WriteFile(path, []byte(annotations), 0644); err != nil {
			return err
		}
	}

	return nil
}
