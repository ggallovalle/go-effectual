package effectual

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/speedata/go-lua"
)

func TryRequireLuarocks(l *lua.State) error {
	version := fmt.Sprintf("%d.%d", lua.VersionMajor, lua.VersionMinor)

	lrPath, err := exec.Command("luarocks", "--lua-version", version, "--local", "path", "--lr-path").Output()
	if err != nil {
		return fmt.Errorf("luarocks path: %v", err)
	}
	lrCpath, err := exec.Command("luarocks", "--lua-version", version, "--local", "path", "--lr-cpath").Output()
	if err != nil {
		return fmt.Errorf("luarocks cpath: %v", err)
	}

	PackagePathAppend(l, strings.TrimSpace(string(lrPath)))
	PackageCPathAppend(l, strings.TrimSpace(string(lrCpath)))

	return nil
}