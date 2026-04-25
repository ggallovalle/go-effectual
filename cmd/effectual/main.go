package main

import (
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/ggallovalle/go-effectual"
	"github.com/ggallovalle/go-effectual/fantastic4/vfs4"
	"github.com/ggallovalle/go-effectual/std"
	"github.com/speedata/go-lua"
	"github.com/spf13/cobra"
	"github.com/twpayne/go-vfs"
)

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
	fileDir, err := filepath.Abs(filepath.Dir(file))
	if err != nil {
		log.Fatal(err)
	}

	l := lua.NewStateEx()
	lua.OpenLibraries(l)

	stdPkg := std.NewStdPackage()
	if err := stdPkg.OpenLib(l, slog.Default()); err != nil {
		log.Fatal(err)
	}
	packagePath := []string{
		filepath.Join(fileDir, "lua", "?.lua"),
		filepath.Join(fileDir, "lua", "?", "init.lua"),
	}

	if err := effectual.PackagePathReplace(l, packagePath); err != nil {
		log.Fatal(err)
	}

	if err := effectual.TryRequireLuarocks(l, fileDir); err != nil {
		slog.Default().Debug("luarocks not available", "err", err)
	}

	if err := lua.DoFile(l, file); err != nil {
		log.Fatalf("Error running %s: %v", file, err)
	}
}

func runLuaDefs(cmd *cobra.Command, args []string) {
	folder := filepath.Join(args[0], "definitions")
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

	pkg := std.NewStdPackage()
	if err := pkg.GenerateAnnotations(fs, folder); err != nil {
		log.Fatalf("Error generating annotations: %v", err)
	}
}
