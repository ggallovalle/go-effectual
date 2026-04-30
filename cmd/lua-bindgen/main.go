package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ggallovalle/go-effectual/internal/luagen"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "lua-bindgen",
		Short: "Generate Lua binding code from Go types",
	}

	generateCmd := &cobra.Command{
		Use:   "generate [source.go]",
		Short: "Generate binding functions for a Go type",
		Args:  cobra.ExactArgs(1),
		RunE:  runGenerate,
	}

	generateCmd.Flags().String("package", "", "Go package name")
	generateCmd.Flags().String("type", "", "Go type name to generate bindings for")
	generateCmd.Flags().String("module", "", "Lua module name (e.g., std.serde.query)")
	generateCmd.Flags().StringSlice("skip", nil, "Method names to skip")
	generateCmd.Flags().StringSlice("nil-map", nil, "Method names that map empty string to nil")
	generateCmd.Flags().StringSlice("force-method", nil, "Force 0-arg returning methods to be methods instead of getters")
	generateCmd.Flags().StringSlice("skip-fields", nil, "Struct field names to skip")
	generateCmd.Flags().StringSlice("import", nil, "Import aliases (format: alias=path)")
	generateCmd.Flags().String("output", "", "Output directory (default: same as source)")

	_ = generateCmd.MarkFlagRequired("type")

	rootCmd.AddCommand(generateCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runGenerate(cmd *cobra.Command, args []string) error {
	sourceFile := args[0]

	typeName, _ := cmd.Flags().GetString("type")
	module, _ := cmd.Flags().GetString("module")
	pkg, _ := cmd.Flags().GetString("package")
	skipList, _ := cmd.Flags().GetStringSlice("skip")
	nilMapList, _ := cmd.Flags().GetStringSlice("nil-map")
	forceMethodList, _ := cmd.Flags().GetStringSlice("force-method")
	skipFieldsList, _ := cmd.Flags().GetStringSlice("skip-fields")
	importList, _ := cmd.Flags().GetStringSlice("import")
	outputDir, _ := cmd.Flags().GetString("output")

	info, annotations, err := luagen.ParseSource(sourceFile, typeName)
	if err != nil {
		return fmt.Errorf("parse source: %w", err)
	}

	if pkg == "" {
		pkg = info.Package
	}

	skip := make(map[string]bool)
	for _, s := range skipList {
		skip[s] = true
	}
	for k, v := range annotations.Skip {
		if _, exists := skip[k]; !exists {
			skip[k] = v
		}
	}

	nilMap := make(map[string]bool)
	for _, s := range nilMapList {
		nilMap[s] = true
	}
	for k, v := range annotations.NilMap {
		if _, exists := nilMap[k]; !exists {
			nilMap[k] = v
		}
	}

	forceMethod := make(map[string]bool)
	for _, s := range forceMethodList {
		forceMethod[s] = true
	}
	for k, v := range annotations.ForceMethod {
		if _, exists := forceMethod[k]; !exists {
			forceMethod[k] = v
		}
	}

	skipFields := make(map[string]bool)
	for _, s := range skipFieldsList {
		skipFields[s] = true
	}
	for k, v := range annotations.SkipFields {
		if _, exists := skipFields[k]; !exists {
			skipFields[k] = v
		}
	}

	imports := make(map[string]string)
	for _, imp := range importList {
		parts := strings.SplitN(imp, "=", 2)
		if len(parts) == 2 {
			imports[parts[0]] = parts[1]
		} else {
			imports[""] = imp
		}
	}

	cfg := &luagen.GenConfig{
		Package:     pkg,
		TypeName:    typeName,
		Module:      module,
		Skip:        skip,
		NilMap:      nilMap,
		ForceMethod: forceMethod,
		SkipFields:  skipFields,
		Imports:     imports,
	}

	if cfg.Module == "" {
		cfg.Module = annotations.Module
	}

	luagen.Classify(info, cfg)

	source := luagen.Generate(info, cfg)

	var outputPath string
	if outputDir != "" {
		outputPath, err = luagen.WriteToDir(info, source, outputDir)
	} else {
		sourceDir := filepath.Dir(sourceFile)
		outputPath, err = luagen.WriteToDir(info, source, sourceDir)
	}
	if err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	fmt.Printf("Generated %s\n", outputPath)
	return nil
}
