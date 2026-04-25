package effectual

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/speedata/go-lua"
)

func TryRequireLuarocks(l *lua.State, dir string) error {
	version := fmt.Sprintf("%d.%d", lua.VersionMajor, lua.VersionMinor)

	lrPath, err := exec.Command("luarocks", "--lua-version", version, "--local", "path", "--lr-path").Output()
	if err != nil {
		return fmt.Errorf("luarocks path: %v", err)
	}
	// luarocks config deploy_lua_dir
	deployLuaDir := filepath.Join(dir, "lua_modules", "share", "lua", version)

	PackagePathAppend(l, filepath.Join(deployLuaDir, "?.lua")+";"+filepath.Join(deployLuaDir, "?", "init.lua"))
	PackagePathAppend(l, strings.TrimSpace(string(lrPath)))

	return nil
}
