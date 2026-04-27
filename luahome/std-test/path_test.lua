local path = require("std.path")

local Suite = {
	name = "std.path",
	cases = {
		{
			name = "MAIN_SEPARATOR",
			fn = function(ctx)
				ctx:expect(path.MAIN_SEPARATOR):equals("/")
			end,
		},
		{
			name = "join",
			fn = function(ctx)
				local home = path.new("/home/some-user")
				local config_nvim = path.join(home, ".config", "nvim")
				ctx:expect(tostring(config_nvim)):equals("/home/some-user/.config/nvim")
				ctx:expect(tostring(path.join())):equals("")
				local home_2 = path.join(home)
				home_2:push("bye")
				ctx:expect(tostring(home)):not_equals(tostring(home_2))
			end,
		},
		{
			name = "absolute",
			fn = function(ctx)
				local target = path.new("foo/bar")
				local absolute = path.absolute("foo/bar")
				ctx:expect(absolute:ends_with(target)):is_true()
			end,
		},
		{
			name = "pathbuf_mut_operations",
			fn = function(ctx)
				local home = path.new("/home/some-user")
				home:push("Documents")
				ctx:expect(tostring(home)):equals("/home/some-user/Documents")
				home:pop()
				ctx:expect(tostring(home)):equals("/home/some-user")
			end,
		},
		{
			name = "pathbuf_meta_concat",
			fn = function(ctx)
				local home = path.new("/home/some-user")
				local downloads = home .. " has a downloads folder"
				ctx:expect(tostring(downloads)):equals("/home/some-user has a downloads folder")
				local music = "music at " .. home
				ctx:expect(tostring(music)):equals("music at /home/some-user")
			end,
		},
		{
			name = "pathbuf_join",
			fn = function(ctx)
				local home = path.new("/home/some-user")
				local relative = path.new("Downloads")
				local downloads = home / "Downloads"
				local downloads2 = home:join("Downloads")
				local downloads3 = "/root/" / home
				local downloads4 = "/usr/home/" / relative

				ctx:expect(tostring(downloads)):equals("/home/some-user/Downloads")
				ctx:expect(tostring(home)):equals("/home/some-user")
				ctx:expect(tostring(downloads)):equals(tostring(downloads2))
				ctx:expect(tostring(downloads3)):equals(tostring(home))
				ctx:expect(tostring(downloads4)):equals("/usr/home/Downloads")
			end,
		},
		{
			name = "pathbuf_components",
			fn = function(ctx)
				local home = path.new("/home/some-user")
				home:push("Documents")
				local components = home.components
				ctx:expect(#components):equals(4)
				ctx:expect(components[1]):equals("/")
				ctx:expect(components[2]):equals("home")
				ctx:expect(components[3]):equals("some-user")
				ctx:expect(components[4]):equals("Documents")
			end,
		},
		{
			name = "pathbuf_parent",
			fn = function(ctx)
				local foobar = path.new("/foo/bar")
				local parent = foobar.parent
				ctx:expect(tostring(parent)):equals("/foo")
				local grandparent = parent.parent
				ctx:expect(tostring(grandparent)):equals("/")
				ctx:expect(grandparent.parent):is_nil()
			end,
		},
		{
			name = "pathbuf_ancestors",
			fn = function(ctx)
				local foobar = path.new("/foo/bar")
				local ancestors = foobar.ancestors
				ctx:expect(#ancestors):equals(3)
				ctx:expect(tostring(ancestors[1])):equals("/foo/bar")
				ctx:expect(tostring(ancestors[2])):equals("/foo")
				ctx:expect(tostring(ancestors[3])):equals("/")
				ctx:expect(ancestors[3].parent):is_nil()
			end,
		},
		{
			name = "pathbuf_ends_with",
			fn = function(ctx)
				local source = path.new("/etc/resolv.conf")
				ctx:expect(source:ends_with("resolv.conf")):is_true()
				ctx:expect(source:ends_with(path.new("etc/resolv.conf"))):is_true()
				ctx:expect(source:ends_with("/etc/resolv.conf")):is_true()
				ctx:expect(source:ends_with("/resolv.conf")):is_false()
				ctx:expect(source:ends_with(path.new("conf"))):is_false()
			end,
		},
		{
			name = "pathbuf_starts_with",
			fn = function(ctx)
				local source = path.new("/etc/passwd")
				ctx:expect(source:starts_with("/etc")):is_true()
				ctx:expect(source:starts_with("/etc/")):is_true()
				ctx:expect(source:starts_with("/etc/passwd")):is_true()
				ctx:expect(source:starts_with("/etc/passwd/")):is_true()
				ctx:expect(source:starts_with("/etc/passwd///")):is_true()

				ctx:expect(source:starts_with("/e")):is_false()
				ctx:expect(source:starts_with("/etc/passwd.txt")):is_false()

				ctx:expect(path.new("/etc/foo.rs"):starts_with("/etc/foo")):is_false()
			end,
		},
		{
			name = "pathbuf_extension",
			fn = function(ctx)
				local source = path.new("/etc/resolv.conf")
				ctx:expect(source.extension):equals("conf")

				local noext = path.new("/etc/resolv")
				ctx:expect(noext.extension):is_nil()

				local hiddenfile = path.new("/etc/.resolv")
				ctx:expect(hiddenfile.extension):is_nil()

				local hiddenfilewithdot = path.new("/etc/.resolv.conf")
				ctx:expect(hiddenfilewithdot.extension):equals("conf")

				local multiple_dots = path.new("/etc/archive.tar.gz")
				ctx:expect(multiple_dots.extension):equals("gz")

				local nodot = path.new("/etc/")
				ctx:expect(nodot.extension):is_nil()

				local root = path.new("/")
				ctx:expect(root.extension):is_nil()

				local empty = path.new("")
				ctx:expect(empty.extension):is_nil()
			end,
		},
		{
			name = "pathbuf_file_stem",
			fn = function(ctx)
				local source = path.new("/etc/resolv.conf")
				ctx:expect(source.file_stem):equals("resolv")

				local noext = path.new("/etc/resolv")
				ctx:expect(noext.file_stem):equals("resolv")

				local hiddenfile = path.new("/etc/.resolv")
				ctx:expect(hiddenfile.file_stem):equals(".resolv")

				local hiddenfilewithdot = path.new("/etc/.resolv.conf")
				ctx:expect(hiddenfilewithdot.file_stem):equals(".resolv")

				local multiple_dots = path.new("/etc/archive.tar.gz")
				ctx:expect(multiple_dots.file_stem):equals("archive.tar")

				local nodot = path.new("/etc/")
				ctx:expect(nodot.file_stem):equals("etc")

				local root = path.new("/")
				ctx:expect(root.file_stem):is_nil()

				local empty = path.new("")
				ctx:expect(empty.file_stem):is_nil()
			end,
		},
		{
			name = "pathbuf_file_name",
			fn = function(ctx)
				local source = path.new("/etc/resolv.conf")
				ctx:expect(source.file_name):equals("resolv.conf")

				local noext = path.new("/etc/resolv")
				ctx:expect(noext.file_name):equals("resolv")

				local hiddenfile = path.new("/etc/.resolv")
				ctx:expect(hiddenfile.file_name):equals(".resolv")

				local hiddenfilewithdot = path.new("/etc/.resolv.conf")
				ctx:expect(hiddenfilewithdot.file_name):equals(".resolv.conf")

				local multiple_dots = path.new("/etc/archive.tar.gz")
				ctx:expect(multiple_dots.file_name):equals("archive.tar.gz")

				local nodot = path.new("/etc/")
				ctx:expect(nodot.file_name):equals("etc")

				local root = path.new("/")
				ctx:expect(root.file_name):is_nil()

				local empty = path.new("")
				ctx:expect(empty.file_name):is_nil()

				local parentdir = path.new("/foo/bar/..")
				ctx:expect(parentdir.file_name):is_nil()
			end,
		},
		{
			name = "pathbuf_relativity",
			fn = function(ctx)
				local abs = path.new("/foo/bar")
				ctx:expect(abs.is_absolute):is_true()
				ctx:expect(abs.is_relative):is_false()
				ctx:expect(abs.has_root):is_true()

				local rel = path.new("foo/bar")
				ctx:expect(rel.is_absolute):is_false()
				ctx:expect(rel.is_relative):is_true()
				ctx:expect(rel.has_root):is_false()

				local rel2 = path.new("../foo/bar")
				ctx:expect(rel2.is_absolute):is_false()
				ctx:expect(rel2.is_relative):is_true()
				ctx:expect(rel2.has_root):is_false()

				local root = path.new("/")
				ctx:expect(root.is_absolute):is_true()
				ctx:expect(root.is_relative):is_false()
				ctx:expect(root.has_root):is_true()

				local empty = path.new("")
				ctx:expect(empty.is_absolute):is_false()
				ctx:expect(empty.is_relative):is_true()
				ctx:expect(empty.has_root):is_false()
			end,
		},
		{
			name = "pathbuf_strip_prefix",
			fn = function(ctx)
				local p = path.new("/a/b/c/d")
				local prefix = path.new("/a/b")
				local stripped = p:strip_prefix(prefix)
				ctx:expect(tostring(stripped)):equals("c/d")

				local wrong_prefix = path.new("/a/x")
				local stripped2, err2 = p:strip_prefix(wrong_prefix)
				ctx:expect(stripped2):is_nil()
				ctx:expect(err2):equals("prefix not found")
			end,
		},
		{
			name = "pathbuf_with",
			fn = function(ctx)
				local p = path.new("/a/b/c/d.txt")
				local with_ext = p:with_extension("md")
				ctx:expect(tostring(with_ext)):equals("/a/b/c/d.md")

				local with_name = p:with_file_name("other.conf")
				ctx:expect(tostring(with_name)):equals("/a/b/c/other.conf")
			end,
		},
		{
			name = "MAIN_SEPARATOR constants",
			fn = function(ctx)
				ctx:expect(path.posix.MAIN_SEPARATOR):equals("/")
				ctx:expect(path.win32.MAIN_SEPARATOR):equals("\\")
			end,
		},
		{
			name = "posix path construction",
			fn = function(ctx)
				local p = path.posix.new("foo/bar")
				ctx:expect(tostring(p)):equals("foo/bar")
				ctx:expect(tostring(p:join("baz"))):equals("foo/bar/baz")
			end,
		},
		{
			name = "win32 path construction",
			fn = function(ctx)
				local p = path.win32.new("foo\\bar")
				ctx:expect(tostring(p)):equals("foo\\bar")
				ctx:expect(tostring(p:join("baz"))):equals("foo\\bar\\baz")
			end,
		},
		{
			name = "posix div operator",
			fn = function(ctx)
				local p = path.posix.new("foo")
				local joined = p / "bar"
				ctx:expect(tostring(joined)):equals("foo/bar")
			end,
		},
		{
			name = "win32 div operator",
			fn = function(ctx)
				local p = path.win32.new("foo")
				local joined = p / "bar"
				ctx:expect(tostring(joined)):equals("foo\\bar")
			end,
		},
		{
			name = "posix components",
			fn = function(ctx)
				local p = path.posix.new("/a/b/c")
				local comps = p.components
				ctx:expect(#comps):equals(4)
				ctx:expect(comps[1]):equals("/")
				ctx:expect(comps[2]):equals("a")
				ctx:expect(comps[3]):equals("b")
				ctx:expect(comps[4]):equals("c")
			end,
		},
		{
			name = "win32 components",
			fn = function(ctx)
				local p = path.win32.new("\\a\\b\\c")
				local comps = p.components
				ctx:expect(#comps):equals(4)
				ctx:expect(comps[1]):equals("\\")
				ctx:expect(comps[2]):equals("a")
				ctx:expect(comps[3]):equals("b")
				ctx:expect(comps[4]):equals("c")
			end,
		},
		{
			name = "posix ancestors",
			fn = function(ctx)
				local p = path.posix.new("/a/b/c")
				local ancs = p.ancestors
				ctx:expect(#ancs):equals(4)
				ctx:expect(tostring(ancs[1])):equals("/a/b/c")
				ctx:expect(tostring(ancs[2])):equals("/a/b")
				ctx:expect(tostring(ancs[3])):equals("/a")
				ctx:expect(tostring(ancs[4])):equals("/")
			end,
		},
		{
			name = "win32 ancestors",
			fn = function(ctx)
				local p = path.win32.new("\\a\\b\\c")
				local ancs = p.ancestors
				ctx:expect(#ancs):equals(4)
				ctx:expect(tostring(ancs[1])):equals("\\a\\b\\c")
				ctx:expect(tostring(ancs[2])):equals("\\a\\b")
				ctx:expect(tostring(ancs[3])):equals("\\a")
				ctx:expect(tostring(ancs[4])):equals("\\")
			end,
		},
		{
			name = "posix ends_with",
			fn = function(ctx)
				local p = path.posix.new("/foo/bar")
				ctx:expect(p:ends_with("bar")):is_true()
				ctx:expect(p:ends_with("foo/bar")):is_true()
				ctx:expect(p:ends_with("/bar")):is_false()
				ctx:expect(p:ends_with("\\bar")):is_false()
			end,
		},
		{
			name = "win32 ends_with",
			fn = function(ctx)
				local p = path.win32.new("\\foo\\bar")
				ctx:expect(p:ends_with("bar")):is_true()
				ctx:expect(p:ends_with("foo\\bar")):is_true()
				ctx:expect(p:ends_with("\\bar")):is_false()
				ctx:expect(p:ends_with("/bar")):is_false()
			end,
		},
		{
			name = "posix starts_with",
			fn = function(ctx)
				local p = path.posix.new("/foo/bar")
				ctx:expect(p:starts_with("/foo")):is_true()
				ctx:expect(p:starts_with("/foo/")):is_true()
				ctx:expect(p:starts_with("\\foo")):is_false()
			end,
		},
		{
			name = "win32 starts_with",
			fn = function(ctx)
				local p = path.win32.new("\\foo\\bar")
				ctx:expect(p:starts_with("\\foo")):is_true()
				ctx:expect(p:starts_with("\\foo\\")):is_true()
				ctx:expect(p:starts_with("/foo")):is_false()
			end,
		},
		{
			name = "posix is_absolute",
			fn = function(ctx)
				local abs = path.posix.new("/foo")
				local rel = path.posix.new("foo")
				ctx:expect(abs.is_absolute):is_true()
				ctx:expect(abs.has_root):is_true()
				ctx:expect(rel.is_absolute):is_false()
				ctx:expect(rel.has_root):is_false()
			end,
		},
		{
			name = "win32 is_absolute",
			fn = function(ctx)
				local abs = path.win32.new("\\foo")
				local rel = path.win32.new("foo")
				ctx:expect(abs.is_absolute):is_true()
				ctx:expect(abs.has_root):is_true()
				ctx:expect(rel.is_absolute):is_false()
				ctx:expect(rel.has_root):is_false()
			end,
		},
		{
			name = "posix push and pop",
			fn = function(ctx)
				local p = path.posix.new("/foo")
				p:push("bar")
				ctx:expect(tostring(p)):equals("/foo/bar")
				ctx:expect(p:pop()):is_true()
				ctx:expect(tostring(p)):equals("/foo")
			end,
		},
		{
			name = "win32 push and pop",
			fn = function(ctx)
				local p = path.win32.new("\\foo")
				p:push("bar")
				ctx:expect(tostring(p)):equals("\\foo\\bar")
				ctx:expect(p:pop()):is_true()
				ctx:expect(tostring(p)):equals("\\foo")
			end,
		},
		{
			name = "posix with_extension",
			fn = function(ctx)
				local p = path.posix.new("/foo/bar.txt")
				local newp = p:with_extension("md")
				ctx:expect(tostring(newp)):equals("/foo/bar.md")
			end,
		},
		{
			name = "win32 with_extension",
			fn = function(ctx)
				local p = path.win32.new("\\foo\\bar.txt")
				local newp = p:with_extension("md")
				ctx:expect(tostring(newp)):equals("\\foo\\bar.md")
			end,
		},
		{
			name = "posix with_file_name",
			fn = function(ctx)
				local p = path.posix.new("/foo/bar.txt")
				local newp = p:with_file_name("baz.md")
				ctx:expect(tostring(newp)):equals("/foo/baz.md")
			end,
		},
		{
			name = "win32 with_file_name",
			fn = function(ctx)
				local p = path.win32.new("\\foo\\bar.txt")
				local newp = p:with_file_name("baz.md")
				ctx:expect(tostring(newp)):equals("\\foo\\baz.md")
			end,
		},
		{
			name = "posix file_name",
			fn = function(ctx)
				local p = path.posix.new("/foo/bar.txt")
				ctx:expect(p.file_name):equals("bar.txt")
			end,
		},
		{
			name = "win32 file_name",
			fn = function(ctx)
				local p = path.win32.new("\\foo\\bar.txt")
				ctx:expect(p.file_name):equals("bar.txt")
			end,
		},
		{
			name = "posix strip_prefix",
			fn = function(ctx)
				local p = path.posix.new("/a/b/c/d")
				local stripped, err = p:strip_prefix("/a/b")
				ctx:expect(tostring(stripped)):equals("c/d")
			end,
		},
		{
			name = "win32 strip_prefix",
			fn = function(ctx)
				local p = path.win32.new("\\a\\b\\c\\d")
				local stripped, err = p:strip_prefix("\\a\\b")
				ctx:expect(tostring(stripped)):equals("c\\d")
			end,
		},
	},
}

return Suite
