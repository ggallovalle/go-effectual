package std_test

import (
	"testing"

	lua "github.com/speedata/go-lua"

	"github.com/ggallovalle/go-effectual"
	sut "github.com/ggallovalle/go-effectual/std"
	"github.com/stretchr/testify/assert"
)

func Test_LibGoPath(t *testing.T) {
	l := lua.NewState()
	lua.OpenLibraries(l)
	effectual.LuaModOpen(l, sut.MakeModPath())

	t.Run("MAIN_SEPARATOR", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			assert(path.MAIN_SEPARATOR == "/")
		`)
		assert.NoError(t, err)
	})

	t.Run("join", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			local home = path.new("/home/some-user")
			local config_nvim = path.join(home, ".config", "nvim")
			assert(tostring(config_nvim) == "/home/some-user/.config/nvim")
			assert(tostring(path.join()) == "")
			local home_2 = path.join(home)
			home_2:push("bye")
			assert(tostring(home) ~= tostring(home_2))
		`)
		assert.NoError(t, err)
	})

	t.Run("absolute", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			local target = path.new("foo/bar")
			local absolute = path.absolute("foo/bar")
			assert(absolute:ends_with(target))
		`)
		assert.NoError(t, err)
	})

	t.Run("pathbuf_mut_operations", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			local home = path.new("/home/some-user")
			home:push("Documents")
			assert(tostring(home) == "/home/some-user/Documents")
			home:pop()
			assert(tostring(home) == "/home/some-user")
		`)
		assert.NoError(t, err)
	})

	t.Run("pathbuf_meta_concat", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			local home = path.new("/home/some-user")
			local downloads = home .. " has a downloads folder"
			assert(tostring(downloads) == "/home/some-user has a downloads folder")
			local music = "music at " .. home
			assert(tostring(music) == "music at /home/some-user")
		`)
		assert.NoError(t, err)
	})

	t.Run("pathbuf_join", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			local home = path.new("/home/some-user")
			local relative = path.new("Downloads")
			local downloads = home / "Downloads"
			local downloads2 = home:join("Downloads")
			local downloads3 = "/root/" / home
			local downloads4 = "/usr/home/" / relative

			assert(tostring(downloads) == "/home/some-user/Downloads")
			assert(tostring(home) == "/home/some-user")
			assert(tostring(downloads) == tostring(downloads2))
			assert(tostring(downloads3) == tostring(home), "because home is absolute, so it replaces the lhs")
			assert(tostring(downloads4) == "/usr/home/Downloads")
		`)
		assert.NoError(t, err)
	})

	t.Run("pathbuf_components", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			local home = path.new("/home/some-user")
			home:push("Documents")
			local components = home.components
			assert(#components == 4)
			assert(components[1] == "/")
			assert(components[2] == "home")
			assert(components[3] == "some-user")
			assert(components[4] == "Documents")
		`)
		assert.NoError(t, err)
	})

	t.Run("pathbuf_parent", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			local foobar = path.new("/foo/bar")
			local parent = foobar.parent
			assert(tostring(parent) == "/foo")
			local grandparent = parent.parent
			assert(tostring(grandparent) == "/")
			assert(grandparent.parent == nil)
		`)
		assert.NoError(t, err)
	})

	t.Run("pathbuf_ancestors", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			local foobar = path.new("/foo/bar")
			local ancestors = foobar.ancestors
			assert(#ancestors == 3)
			assert(tostring(ancestors[1]) == "/foo/bar")
			assert(tostring(ancestors[2]) == "/foo")
			assert(tostring(ancestors[3]) == "/")
			assert(ancestors[3].parent == nil)
		`)
		assert.NoError(t, err)
	})

	t.Run("pathbuf_ends_with", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			local source = path.new("/etc/resolv.conf")
			assert(source:ends_with("resolv.conf") == true)
			assert(source:ends_with(path.new("etc/resolv.conf")) == true)
			assert(source:ends_with("/etc/resolv.conf") == true)
			assert(source:ends_with("/resolv.conf") == false)
			assert(source:ends_with(path.new("conf")) == false)
		`)
		assert.NoError(t, err)
	})

	t.Run("pathbuf_starts_with", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			local source = path.new("/etc/passwd")
			assert(source:starts_with("/etc") == true)
			assert(source:starts_with("/etc/") == true)
			assert(source:starts_with("/etc/passwd") == true)
			assert(source:starts_with("/etc/passwd/") == true)
			assert(source:starts_with("/etc/passwd///") == true)

			assert(source:starts_with("/e") == false)
			assert(source:starts_with("/etc/passwd.txt") == false)

			assert(path.new("/etc/foo.rs"):starts_with("/etc/foo") == false)
		`)
		assert.NoError(t, err)
	})

	t.Run("pathbuf_extension", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			local source = path.new("/etc/resolv.conf")
			assert(source.extension == "conf")

			local noext = path.new("/etc/resolv")
			assert(noext.extension == nil)

			local hiddenfile = path.new("/etc/.resolv")
			assert(hiddenfile.extension == nil)

			local hiddenfilewithdot = path.new("/etc/.resolv.conf")
			assert(hiddenfilewithdot.extension == "conf")

			local multiple_dots = path.new("/etc/archive.tar.gz")
			assert(multiple_dots.extension == "gz")

			local nodot = path.new("/etc/")
			assert(nodot.extension == nil)

			local root = path.new("/")
			assert(root.extension == nil)

			local empty = path.new("")
			assert(empty.extension == nil)
		`)
		assert.NoError(t, err)
	})

	t.Run("pathbuf_file_stem", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			local source = path.new("/etc/resolv.conf")
			assert(source.file_stem == "resolv")

			local noext = path.new("/etc/resolv")
			assert(noext.file_stem == "resolv")

			local hiddenfile = path.new("/etc/.resolv")
			assert(hiddenfile.file_stem == ".resolv")

			local hiddenfilewithdot = path.new("/etc/.resolv.conf")
			assert(hiddenfilewithdot.file_stem == ".resolv")

			local multiple_dots = path.new("/etc/archive.tar.gz")
			assert(multiple_dots.file_stem == "archive.tar")

			local nodot = path.new("/etc/")
			assert(nodot.file_stem == "etc")

			local root = path.new("/")
			assert(root.file_stem == nil)

			local empty = path.new("")
			assert(empty.file_stem == nil)
		`)
		assert.NoError(t, err)
	})

	t.Run("pathbuf_file_name", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			local source = path.new("/etc/resolv.conf")
			assert(source.file_name == "resolv.conf")

			local noext = path.new("/etc/resolv")
			assert(noext.file_name == "resolv")

			local hiddenfile = path.new("/etc/.resolv")
			assert(hiddenfile.file_name == ".resolv")

			local hiddenfilewithdot = path.new("/etc/.resolv.conf")
			assert(hiddenfilewithdot.file_name == ".resolv.conf")

			local multiple_dots = path.new("/etc/archive.tar.gz")
			assert(multiple_dots.file_name == "archive.tar.gz")

			local nodot = path.new("/etc/")
			assert(nodot.file_name == "etc")

			local root = path.new("/")
			assert(root.file_name == nil)

			local empty = path.new("")
			assert(empty.file_name == nil)

			local parentdir = path.new("/foo/bar/..")
			assert(parentdir.file_name == nil)
		`)
		assert.NoError(t, err)
	})

	t.Run("pathbuf_relativity", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			local abs = path.new("/foo/bar")
			assert(abs.is_absolute == true)
			assert(abs.is_relative == false)
			assert(abs.has_root == true)

			local rel = path.new("foo/bar")
			assert(rel.is_absolute == false)
			assert(rel.is_relative == true)
			assert(rel.has_root == false)

			local rel2 = path.new("../foo/bar")
			assert(rel2.is_absolute == false)
			assert(rel2.is_relative == true)
			assert(rel2.has_root == false)

			local root = path.new("/")
			assert(root.is_absolute == true)
			assert(root.is_relative == false)
			assert(root.has_root == true)

			local empty = path.new("")
			assert(empty.is_absolute == false)
			assert(empty.is_relative == true)
			assert(empty.has_root == false)
		`)
		assert.NoError(t, err)
	})

	t.Run("pathbuf_strip_prefix", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			local p = path.new("/a/b/c/d")
			local prefix = path.new("/a/b")
			local stripped = p:strip_prefix(prefix)
			assert(tostring(stripped) == "c/d")

			local wrong_prefix = path.new("/a/x")
			local stripped2, err2 = p:strip_prefix(wrong_prefix)
			assert(stripped2 == nil)
			assert(err2 == "prefix not found")
		`)
		assert.NoError(t, err)
	})

	t.Run("pathbuf_with", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			local p = path.new("/a/b/c/d.txt")
			local with_ext = p:with_extension("md")
			assert(tostring(with_ext) == "/a/b/c/d.md")

			local with_name = p:with_file_name("other.conf")
			assert(tostring(with_name) == "/a/b/c/other.conf")
		`)
		assert.NoError(t, err)
	})
}

func Test_LibGoPathPosixWin32(t *testing.T) {
	l := lua.NewState()
	lua.OpenLibraries(l)
	effectual.LuaModOpen(l, sut.MakeModPath())

	t.Run("MAIN_SEPARATOR constants", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			assert(path.posix.MAIN_SEPARATOR == "/")
			assert(path.win32.MAIN_SEPARATOR == "\\")
		`)
		assert.NoError(t, err)
	})

	t.Run("posix path construction", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			local p = path.posix.new("foo/bar")
			assert(tostring(p) == "foo/bar")
			assert(tostring(p:join("baz")) == "foo/bar/baz")
		`)
		assert.NoError(t, err)
	})

	t.Run("win32 path construction", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			local p = path.win32.new("foo\\bar")
			assert(tostring(p) == "foo\\bar")
			assert(tostring(p:join("baz")) == "foo\\bar\\baz")
		`)
		assert.NoError(t, err)
	})

	t.Run("posix div operator", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			local p = path.posix.new("foo")
			local joined = p / "bar"
			assert(tostring(joined) == "foo/bar")
		`)
		assert.NoError(t, err)
	})

	t.Run("win32 div operator", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			local p = path.win32.new("foo")
			local joined = p / "bar"
			assert(tostring(joined) == "foo\\bar")
		`)
		assert.NoError(t, err)
	})

	t.Run("posix components", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			local p = path.posix.new("/a/b/c")
			local comps = p.components
			assert(#comps == 4)
			assert(comps[1] == "/")
			assert(comps[2] == "a")
			assert(comps[3] == "b")
			assert(comps[4] == "c")
		`)
		assert.NoError(t, err)
	})

	t.Run("win32 components", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			local p = path.win32.new("\\a\\b\\c")
			local comps = p.components
			assert(#comps == 4)
			assert(comps[1] == "\\")
			assert(comps[2] == "a")
			assert(comps[3] == "b")
			assert(comps[4] == "c")
		`)
		assert.NoError(t, err)
	})

	t.Run("posix ancestors", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			local p = path.posix.new("/a/b/c")
			local ancs = p.ancestors
			assert(#ancs == 4)
			assert(tostring(ancs[1]) == "/a/b/c")
			assert(tostring(ancs[2]) == "/a/b")
			assert(tostring(ancs[3]) == "/a")
			assert(tostring(ancs[4]) == "/")
		`)
		assert.NoError(t, err)
	})

	t.Run("win32 ancestors", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			local p = path.win32.new("\\a\\b\\c")
			local ancs = p.ancestors
			assert(#ancs == 4)
			assert(tostring(ancs[1]) == "\\a\\b\\c")
			assert(tostring(ancs[2]) == "\\a\\b")
			assert(tostring(ancs[3]) == "\\a")
			assert(tostring(ancs[4]) == "\\")
		`)
		assert.NoError(t, err)
	})

	t.Run("posix ends_with", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			local p = path.posix.new("/foo/bar")
			assert(p:ends_with("bar") == true)
			assert(p:ends_with("foo/bar") == true)
			assert(p:ends_with("/bar") == false)
			assert(p:ends_with("\\bar") == false)
		`)
		assert.NoError(t, err)
	})

	t.Run("win32 ends_with", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			local p = path.win32.new("\\foo\\bar")
			assert(p:ends_with("bar") == true)
			assert(p:ends_with("foo\\bar") == true)
			assert(p:ends_with("\\bar") == false)
			assert(p:ends_with("/bar") == false)
		`)
		assert.NoError(t, err)
	})

	t.Run("posix starts_with", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			local p = path.posix.new("/foo/bar")
			assert(p:starts_with("/foo") == true)
			assert(p:starts_with("/foo/") == true)
			assert(p:starts_with("\\foo") == false)
		`)
		assert.NoError(t, err)
	})

	t.Run("win32 starts_with", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			local p = path.win32.new("\\foo\\bar")
			assert(p:starts_with("\\foo") == true)
			assert(p:starts_with("\\foo\\") == true)
			assert(p:starts_with("/foo") == false)
		`)
		assert.NoError(t, err)
	})

	t.Run("posix is_absolute", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			local abs = path.posix.new("/foo")
			local rel = path.posix.new("foo")
			assert(abs.is_absolute == true)
			assert(abs.has_root == true)
			assert(rel.is_absolute == false)
			assert(rel.has_root == false)
		`)
		assert.NoError(t, err)
	})

	t.Run("win32 is_absolute", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			local abs = path.win32.new("\\foo")
			local rel = path.win32.new("foo")
			assert(abs.is_absolute == true)
			assert(abs.has_root == true)
			assert(rel.is_absolute == false)
			assert(rel.has_root == false)
		`)
		assert.NoError(t, err)
	})

	t.Run("posix push and pop", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			local p = path.posix.new("/foo")
			p:push("bar")
			assert(tostring(p) == "/foo/bar")
			assert(p:pop() == true)
			assert(tostring(p) == "/foo")
		`)
		assert.NoError(t, err)
	})

	t.Run("win32 push and pop", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			local p = path.win32.new("\\foo")
			p:push("bar")
			assert(tostring(p) == "\\foo\\bar")
			assert(p:pop() == true)
			assert(tostring(p) == "\\foo")
		`)
		assert.NoError(t, err)
	})

	t.Run("posix with_extension", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			local p = path.posix.new("/foo/bar.txt")
			local newp = p:with_extension("md")
			assert(tostring(newp) == "/foo/bar.md")
		`)
		assert.NoError(t, err)
	})

	t.Run("win32 with_extension", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			local p = path.win32.new("\\foo\\bar.txt")
			local newp = p:with_extension("md")
			assert(tostring(newp) == "\\foo\\bar.md")
		`)
		assert.NoError(t, err)
	})

	t.Run("posix with_file_name", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			local p = path.posix.new("/foo/bar.txt")
			local newp = p:with_file_name("baz.md")
			assert(tostring(newp) == "/foo/baz.md")
		`)
		assert.NoError(t, err)
	})

	t.Run("win32 with_file_name", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			local p = path.win32.new("\\foo\\bar.txt")
			local newp = p:with_file_name("baz.md")
			assert(tostring(newp) == "\\foo\\baz.md")
		`)
		assert.NoError(t, err)
	})

	t.Run("posix file_name", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			local p = path.posix.new("/foo/bar.txt")
			assert(p.file_name == "bar.txt")
		`)
		assert.NoError(t, err)
	})

	t.Run("win32 file_name", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			local p = path.win32.new("\\foo\\bar.txt")
			assert(p.file_name == "bar.txt")
		`)
		assert.NoError(t, err)
	})

	t.Run("posix strip_prefix", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			local p = path.posix.new("/a/b/c/d")
			local stripped, err = p:strip_prefix("/a/b")
			assert(tostring(stripped) == "c/d")
		`)
		assert.NoError(t, err)
	})

	t.Run("win32 strip_prefix", func(t *testing.T) {
		err := lua.DoString(l, `
			local path = require("std.path")
			local p = path.win32.new("\\a\\b\\c\\d")
			local stripped, err = p:strip_prefix("\\a\\b")
			assert(tostring(stripped) == "c\\d")
		`)
		assert.NoError(t, err)
	})
}
