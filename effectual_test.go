package effectual_test

import (
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/speedata/go-lua"
	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
	"github.com/tailscale/hujson"
	"github.com/twpayne/go-vfs"
)

var (
	_ = toml.Decode
	_ = toml.Unmarshal
	_ = toml.Marshal
	_ = lua.NewState
	_ = lua.NewStateEx
	_ = yaml.Marshal
	_ = yaml.Unmarshal
	_ = assert.New
	_ = hujson.Parse
	_ = hujson.Format
	_ = vfs.OSFS
	_ = vfs.Contains
)

func TestStub(t *testing.T) {
}
