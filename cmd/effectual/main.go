package main

import (
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/Shopify/go-lua"
	"github.com/ggallovalle/go-effectual/fantastic4/vfs4"
	"github.com/ggallovalle/go-effectual/std"
	"github.com/spf13/cobra"
	"github.com/twpayne/go-vfs"
)

type luaMod interface {
	Name() string
	Annotations() string
}

type moduleEntry struct {
	name string
	make func() luaMod
}

var availableModules = []moduleEntry{
	{std.ModSlogName, func() luaMod { return std.MakeModSlog() }},
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "effectual",
		Short: "effectual is a Lua interpreter with additional modules",
	}

	luaCmd := &cobra.Command{
		Use:   "lua [file]",
		Short: "Run a Lua file with std.log module loaded",
		Args:  cobra.ExactArgs(1),
		Run:   runLua,
	}

	luaDefsCmd := &cobra.Command{
		Use:   "lua-defs [folder]",
		Short: "Generate Lua annotations to [folder]/definitions/",
		Args:  cobra.ExactArgs(1),
		Run:   runLuaDefs,
	}
	luaDefsCmd.Flags().StringSliceP("module", "m", nil, "Module to generate definitions for")
	luaDefsCmd.Flags().Bool("dry-run", false, "Log file operations instead of performing them")
	luaDefsCmd.Flags().CountP("verbose", "v", "Increase verbosity (-v=warn, -vv=info, -vvv=debug)")

	rootCmd.PersistentFlags().CountP("verbose", "v", "Increase verbosity (-v=warn, -vv=info, -vvv=debug)")
	rootCmd.AddCommand(luaCmd, luaDefsCmd)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runLua(cmd *cobra.Command, args []string) {
	file := args[0]

	l := lua.NewState()
	lua.OpenLibraries(l)

	stdmod := std.MakeModSlog()
	stdmod.OpenLib(l)

	defaultLogger := slog.Default()
	api := stdmod.Api(l)
	api.SetDefault(defaultLogger)

	if err := lua.DoFile(l, file); err != nil {
		log.Fatalf("Error running %s: %v", file, err)
	}
}

func runLuaDefs(cmd *cobra.Command, args []string) {
	folder := filepath.Join(args[0], "definitions")
	moduleNames, _ := cmd.Flags().GetStringSlice("module")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	verbose, _ := cmd.Flags().GetCount("verbose")

	level := slog.LevelError
	switch verbose {
	case 1:
		level = slog.LevelWarn
	case 2:
		level = slog.LevelInfo
	case 3:
		level = slog.LevelDebug
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))

	var fs vfs.FS
	if dryRun {
		fs = vfs4.NewLogVfs(logger, slog.LevelInfo, nil)
	} else {
		fs = vfs4.NewLogVfs(logger, slog.LevelInfo, vfs.OSFS)
	}

	var mods []luaMod
	if len(moduleNames) == 0 {
		for _, entry := range availableModules {
			mods = append(mods, entry.make())
		}
	} else {
		for _, name := range moduleNames {
			var found bool
			for _, entry := range availableModules {
				if entry.name == name {
					mods = append(mods, entry.make())
					found = true
					break
				}
			}
			if !found {
				available := make([]string, 0, len(availableModules))
				for _, e := range availableModules {
					available = append(available, e.name)
				}
				log.Fatalf("Error: unknown module %q. Available modules: %v", name, available)
			}
		}
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
			log.Fatalf("Error creating directory %s: %v", dir, err)
		}

		if err := fs.WriteFile(path, []byte(annotations), 0644); err != nil {
			log.Fatalf("Error writing %s: %v", path, err)
		}

		log.Printf("Generated %s", path)
	}
}
